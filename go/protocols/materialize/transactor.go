package materialize

import (
	"context"
	"encoding/json"
	"fmt"

	pf "github.com/estuary/flow/go/protocols/flow"
	"github.com/sirupsen/logrus"
)

// Transactor is a store-agnostic interface for a materialization driver
// that implements Flow materialization protocol transactions.
type Transactor interface {
	// Load implements the transaction load phase by consuming Load requests
	// from the LoadIterator until a Flush request is read. Requested keys
	// are not known to exist in the store, and very often they won't.
	// Load can ignore keys which are not found in the store. Before Load returns,
	// though, it must ensure loaded() is called for all found documents.
	//
	// Materializations which are delta-updates only are still required to
	// consume the LoadIterator to completion, but the iterator will never
	// yield a document.
	Load(_ *LoadIterator, loaded func(binding int, doc json.RawMessage) error) error
	// Store consumes Store requests from the StoreIterator.
	Store(*StoreIterator) error
	// StartCommit is called upon the Flow runtime's request that the driver
	// begin to commit its transaction.
	// If the driver doesn't use store transactions, StartCommit may be a no-op.
	StartCommit(context.Context, TransactionRequest_StartCommit) (pf.DriverCheckpoint, error)
	// Acknowledge is called upon the Flow runtime's acknowledgement of its
	// recovery log commit.
	//
	// If the driver has an ongoing commit, it must await the completion of
	// that commit before returning.
	//
	// Or, if the driver stages data for idempotent application after commit,
	// it must perform that apply now and return a final error status.
	//
	// Note that Acknowledge may be called multiple times in acknowledgement
	// of a single actual commit. The driver must account for this. If it applies
	// staged data as part of acknowledgement, it must ensure that apply is
	// idempotent.
	Acknowledge(context.Context) error
	// Destroy the Transactor, releasing any held resources.
	Destroy()
}

// RunTransactions processes materialization protocol transactions
// over the established stream against a Transactor.
func RunTransactions(
	stream Driver_TransactionsServer,
	transactor Transactor,
	log *logrus.Entry,
) (_err error) {

	var request *TransactionRequest
	var response *TransactionResponse

	for round := 0; true; round++ {

		// Read Acknowledge => send Acknowledged.
		if err := stream.RecvMsg(request); err != nil {
			return fmt.Errorf("reading expected Acknowledge: %w", err)
		} else if request.Acknowledge == nil {
			return fmt.Errorf("protocol error (expected Acknowledge, got %#v)", request)
		} else if err = transactor.Acknowledge(stream.Context()); err != nil {
			return fmt.Errorf("during Acknowledge phase: %w", err)
		} else if err = WriteAcknowledged(stream, &response); err != nil {
			return err
		}

		log.WithFields(logrus.Fields{
			"round": round,
		}).Debug("Acknowledge finished")

		// Read Load => send Loaded, until Flush is read.
		var loadIt = NewLoadIterator(stream)
		var loaded = 0

		if err := transactor.Load(loadIt, func(binding int, doc json.RawMessage) error {
			loaded++
			return StageLoaded(stream, &response, binding, doc)
		}); loadIt.err != nil {
			// Prefer a LoadIterator.err over `err` as it happened earlier and is
			// likely causal of (or equal to) `err`.
			return loadIt.err
		} else if err != nil {
			return err
		} else if err = WriteFlushed(stream, &response); err != nil {
			return err
		}

		log.WithFields(logrus.Fields{
			"round":  round,
			"total":  loadIt.total,
			"loaded": loaded,
		}).Debug("transactor Load finished")

		var storeIt = NewStoreIterator(stream)

		if err := transactor.Store(storeIt); storeIt.err != nil {
			// Prefer the StoreIterator's error over `err` as it happened first,
			// is likely causal of (or equal to) `err`.
			return storeIt.err
		} else if err != nil {
			return err
		}

		log.WithFields(logrus.Fields{
			"round": round,
			"store": storeIt.total,
		}).Debug("Store finished")

		if checkpoint, err := transactor.StartCommit(stream.Context(), *storeIt.req.StartCommit); err != nil {
			return fmt.Errorf("during StartCommit phase: %w", err)
		} else if err = WriteStartedCommit(stream, &response, checkpoint); err != nil {
			return err
		}
	}
	panic("not reached")
}
