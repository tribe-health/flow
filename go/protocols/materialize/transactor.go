package materialize

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	pf "github.com/estuary/flow/go/protocols/flow"
	"github.com/sirupsen/logrus"
	"go.gazette.dev/core/broker/client"
)

// Transactor is a store-agnostic interface for a materialization driver
// that implements Flow materialization protocol transactions.
type Transactor interface {
	// Load implements the transaction load phase by consuming Load requests
	// from the LoadIterator and calling the provided `loaded` callback.
	// Load can ignore keys which are not found in the store, and it may
	// defer calls to `loaded` for as long as it wishes, so long as `loaded`
	// is called for every found document prior to returning.
	//
	// If this Transactor chooses to uses concurrency in StartCommit, note
	// that Load may be called while its returned OpFuture is still running.
	// However, absent an error, LoadIterator.Next() will not return false
	// until that OpFuture has resolved.
	//
	// Typically a Transactor that chooses to use concurrency should "stage"
	// loads for later evaluation, and then evaluate all loads upon that
	// commit resolving, or even wait until Next() returns false.
	//
	// Waiting for the prior commit ensures that evaluated loads reflect its
	// updates, and thus meet the formal "read-committed"
	// guarantee required by the runtime.
	Load(_ *LoadIterator, loaded func(binding int, doc json.RawMessage) error) error
	// Store consumes Store requests from the StoreIterator.
	Store(*StoreIterator) error
	// StartCommit begins to commit the transaction. Upon its return a commit
	// operation may still be running in the background, and the returned
	// OpFuture must resolve with its completion.
	// (Note that upon its resolution, Acknowledged will be sent to the Runtime).
	//
	// # When using the "Remote Store is Authoritative" pattern:
	//
	// StartCommit must include `runtimeCheckpoint` within its endpoint
	// transaction and either immediately or asynchronously commit.
	// If the Transactor commits synchronously, it may return a nil OpFuture.
	//
	// # When using the "Recovery Log is Authoritative with Idempotent Apply" pattern:
	//
	// StartCommit must return a DriverCheckpoint which encodes the staged
	// application. It must begin an asynchronous application of this staged
	// update, returning its OpFuture.
	//
	// That async application MUST await a call to Acknowledge before taking action,
	// however, to ensure that the DriverCheckpoint returned by StartCommit
	// has been durably committed to the runtime recovery log.
	//
	// Note it's possible that the DriverCheckpoint may commit to the log,
	// but then the runtime or this Transactor may crash before the application
	// is able to complete. For this reason, on startup a Transactor must take
	// care to (re-)apply a staged update encoded in the opened DriverCheckpoint.
	StartCommit(_ context.Context, runtimeCheckpoint []byte) (*pf.DriverCheckpoint, pf.OpFuture, error)
	// RuntimeCommitted is called after StartCommit, upon the runtime completing
	// its commit to its recovery log.
	//
	// Most Transactors can ignore this signal, but those using the
	// "Recovery Log is Authoritative with Idempotent Apply" should use it
	// to unblock a background apply operation, which may (only) now proceed.
	RuntimeCommitted(context.Context) error
	// Destroy the Transactor, releasing any held resources.
	Destroy()
}

// RunTransactions processes materialization protocol transactions
// over the established stream against a Driver.
func RunTransactions(
	stream Driver_TransactionsServer,
	newTransactor func(context.Context, TransactionRequest_Open) (Transactor, *TransactionResponse_Opened, error),
) (_err error) {

	var rxRequest, err = ReadOpen(stream)
	if err != nil {
		return err
	}
	transactor, opened, err := newTransactor(stream.Context(), *rxRequest.Open)
	if err != nil {
		return err
	}

	defer func() {
		if _err != nil {
			logrus.WithError(_err).Error("RunTransactions failed")
		} else {
			logrus.Debug("RunTransactions finished")
		}
		transactor.Destroy()
	}()

	txResponse, err := WriteOpened(stream, opened)
	if err != nil {
		return err
	}

	var (
		// doAckErr is the last doCommitAck() result,
		// and is readable upon its close of its parameter `doAckDone`.
		doAckErr error
		// doLoadErr is the last doLoad() result,
		// and is readable upon its close of its parameter `doLoadDone`.
		doLoadErr error
	)

	// doAck is a closure for the async await of a previously
	// started commit, upon which it writes Acknowledged to the runtime.
	// It has an exclusive ability to write to `stream` until it returns.
	var doAck = func(
		round int,
		opCommit pf.OpFuture, // Resolves when the prior commit completes.
		doAckDone chan<- struct{}, // To be closed upon return.
		doLoadDone <-chan struct{}, // Signaled when doLoad() has completed.
	) (__out error) {

		defer func() {
			logrus.WithFields(logrus.Fields{
				"round": round,
				"error": __out,
			}).Debug("doAck finished")

			doAckErr = __out
			close(doAckDone)
		}()

		// Wait for commit to complete, with cancellation checks.
		select {
		case <-opCommit.Done():
			if err := opCommit.Err(); err != nil {
				return err
			}
		case <-doLoadDone:
			// doLoad() must have errored, as it otherwise cannot
			// complete until we send Acknowledged.
			return nil
		}
		logrus.Debug("Commit finished")

		return WriteAcknowledged(stream, &txResponse)
	}

	// doLoad is a closure for async execution of Transactor.Load.
	var doLoad = func(
		round int,
		it *LoadIterator,
		doAckDone <-chan struct{}, // Signaled when doAck() has completed.
		doLoadDone chan<- struct{}, // To be closed upon return.
	) (__out error) {

		var loaded int
		defer func() {
			logrus.WithFields(logrus.Fields{
				"round":  round,
				"total":  it.total,
				"loaded": loaded,
				"error":  __out,
			}).Debug("doLoad finished")

			doLoadErr = __out
			close(doLoadDone)
		}()

		var err = transactor.Load(it, func(binding int, doc json.RawMessage) error {
			if doAckDone != nil {
				// Wait for doAck() to complete and then clear our local copy of its channel.
				_, doAckDone = <-doAckDone, nil
			}
			if doAckErr != nil {
				// We cannot write Loaded responses if doAck() failed as it would
				// violate the protocol (out-of-order responses). Bail out.
				return context.Canceled
			}

			loaded++
			return WriteLoaded(stream, &txResponse, binding, doc)
		})

		if doAckDone == nil && doAckErr != nil {
			return nil // Cancelled by doAck().
		} else if it.err != nil {
			// Prefer the iterator's error over `err` as it's earlier in the chain
			// of dependency and is likely causal of (or equal to) `err`.
			return it.err
		}
		return err
	}

	// opCommit is a future for the most-recent started commit.
	var opCommit pf.OpFuture = client.FinishedOperation(nil)

	for round := 0; true; round++ {
		var (
			doAckDone  = make(chan struct{}) // Signals doCommitAck() is done.
			doLoadDone = make(chan struct{}) // Signals doLoad() is done.
			loadIt     = LoadIterator{stream: stream, request: &rxRequest}
		)

		if err = ReadAcknowledge(stream, &rxRequest); err != nil {
			return err
		} else if round == 0 {
			// Suppress explicit Acknowledge of the opened commit.
			// newTransactor() is expected to have already taken any required
			// action to apply this commit to the store (where applicable).
		} else if err = transactor.RuntimeCommitted(stream.Context()); err != nil {
			return fmt.Errorf("transactor.RuntimeCommitted: %w", err)
		}

		// Begin async acknowledgement of the prior transaction.
		// On completion, Acknowledged has been written to the Runtime,
		// which allows the concurrent Load phase to begin to close.
		go doAck(round, opCommit, doAckDone, doLoadDone)

		// Begin an async load of the current transaction.
		// On completion, |loadCh| is closed and |loadErr| is its status.
		go doLoad(round, &loadIt, doAckDone, doLoadDone)

		// Join over the previous commit and current load.
		// We deliberately don't join over `doAckDone`, as it can only fail
		// because of opCommit.Err() or because `doLoad()` failed.
		for doAckDone != nil || doLoadDone != nil {
			select {
			case <-doAckDone:
				if doAckErr != nil {
					return fmt.Errorf("commit failed: %w", doAckErr)
				}
				doAckDone = nil
			case <-doLoadDone:
				if doLoadErr != nil && doLoadErr != io.EOF {
					return fmt.Errorf("transactor.Load: %w", doLoadErr)
				}
				doLoadDone = nil
			}
		}

		if doLoadErr == io.EOF {
			return nil // Graceful shutdown.
		}

		if err = ReadFlush(&rxRequest); err != nil {
			return err
		} else if err = WriteFlushed(stream, &txResponse); err != nil {
			return err
		}
		logrus.WithField("round", round).Debug("wrote Flushed")

		// Process all Store requests until StartCommit is read.
		var storeIt = StoreIterator{stream: stream, request: &rxRequest}
		if err = transactor.Store(&storeIt); storeIt.err != nil {
			return storeIt.err // Prefer an iterator error as it's more directly causal.
		} else if err != nil {
			return fmt.Errorf("transactor.Store: %w", err)
		}
		logrus.WithFields(logrus.Fields{"round": round, "stored": storeIt.total}).Debug("Store finished")

		var runtimeCheckpoint []byte
		var driverCheckpoint *pf.DriverCheckpoint

		if runtimeCheckpoint, err = ReadStartCommit(&rxRequest); err != nil {
			return err
		} else if driverCheckpoint, opCommit, err = transactor.StartCommit(stream.Context(), runtimeCheckpoint); err != nil {
			return fmt.Errorf("transactor.StartCommit: %w", err)
		} else if err = WriteStartedCommit(stream, &txResponse, driverCheckpoint); err != nil {
			return err
		}

		// As a convenience, map a nil OpFuture to a pre-resolved one so the
		// rest of our handling can ignore the nil case.
		if opCommit == nil {
			opCommit = client.FinishedOperation(nil)
		}
	}
	panic("not reached")
}
