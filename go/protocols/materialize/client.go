package materialize

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	math "math"
	"sync"

	"github.com/estuary/flow/go/protocols/fdb/tuple"
	pf "github.com/estuary/flow/go/protocols/flow"
	"github.com/sirupsen/logrus"
	"go.gazette.dev/core/broker/client"
)

// TxnClient is a client of a driver's Transactions RPC.
type TxnClient struct {
	client Driver_TransactionsClient
	// Combiners of the materialization, one for each binding.
	combiners []pf.Combiner
	// Guards combiners, which are accessed concurrently from readAcknowledgedAndLoaded().
	combinersMu sync.Mutex
	// Flighted keys of the current transaction for each binding, plus a bounded cache of
	// fully-reduced documents of the last transaction.
	flighted []map[string]json.RawMessage
	// OpFuture that's resolved on completion of a current Loaded phase,
	// or nil if readAcknowledgedAndLoaded is not currently running.
	opLoaded   *client.AsyncOperation
	opened     *TransactionResponse_Opened // Opened response returned by the server while opening.
	rxResponse TransactionResponse         // Response which is received into.
	spec       *pf.MaterializationSpec     // Specification of this Transactions client.
	txRequest  TransactionRequest          // Request which is sent from.
	version    string                      // Version of the client's MaterializationSpec.

	// Temporary storage for a retained DriverCheckpoint. This will be removed.
	deprecatedDriverCP pf.DriverCheckpoint
}

func OpenTransactions(
	ctx context.Context,
	driver DriverClient,
	driverCheckpoint json.RawMessage,
	newCombinerFn func(*pf.MaterializationSpec_Binding) (pf.Combiner, error),
	range_ pf.RangeSpec,
	spec *pf.MaterializationSpec,
	version string,
) (*TxnClient, error) {

	if range_.RClockBegin != 0 || range_.RClockEnd != math.MaxUint32 {
		return nil, fmt.Errorf("materializations cannot split on r-clock: " + range_.String())
	}

	var combiners []pf.Combiner
	var flighted []map[string]json.RawMessage

	for _, b := range spec.Bindings {
		var combiner, err = newCombinerFn(b)
		if err != nil {
			return nil, fmt.Errorf("creating %s combiner: %w", b.Collection.Collection, err)
		}
		combiners = append(combiners, combiner)
		flighted = append(flighted, make(map[string]json.RawMessage))
	}

	rpc, err := driver.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("driver.Transactions: %w", err)
	}
	// Close RPC if remaining initialization fails.
	defer func() {
		if rpc != nil {
			_ = rpc.CloseSend()
		}
	}()

	txRequest, err := WriteOpen(rpc,
		&TransactionRequest_Open{
			Materialization:      spec,
			Version:              version,
			KeyBegin:             range_.KeyBegin,
			KeyEnd:               range_.KeyEnd,
			DriverCheckpointJson: driverCheckpoint,
		})
	if err != nil {
		return nil, err
	}

	// Read Opened response with driver's optional Flow Checkpoint.
	rxResponse, err := ReadOpened(rpc)
	if err != nil {
		return nil, err
	}

	// Write Acknowledge request to re-acknowledge the last commit.
	if err := WriteAcknowledge(rpc, &txRequest); err != nil {
		return nil, err
	}

	var c = &TxnClient{
		combiners: combiners,
		//combinersMu: sync.Mutex{},
		flighted:   flighted,
		opLoaded:   client.NewAsyncOperation(),
		opened:     rxResponse.Opened,
		client:     rpc,
		rxResponse: rxResponse,
		spec:       spec,
		txRequest:  txRequest,
		version:    version,

		// TODO(johnny): Remove me.
		deprecatedDriverCP: pf.DriverCheckpoint{},
	}
	go c.readAcknowledgedAndLoaded(client.NewAsyncOperation())

	rpc = nil // Don't run deferred CloseSend.
	return c, nil
}

// Opened returns the driver's prior Opened response.
func (c *TxnClient) Opened() *TransactionResponse_Opened { return c.opened }

// Close the TxnClient. Close returns an error if the RPC is not in an
// Acknowledged and idle state, or on any other error.
func (c *TxnClient) Close() error {
	c.client.CloseSend()

	var opLoadedErr error
	if c.opLoaded != nil {
		opLoadedErr = c.opLoaded.Err()
	}

	for _, c := range c.combiners {
		c.Destroy()
	}

	// EOF is a graceful shutdown.
	if err := opLoadedErr; err != io.EOF {
		return err
	}
	return nil
}

// AddDocument to the current transaction under the given binding and tuple-encoded key.
func (c *TxnClient) AddDocument(binding int, packedKey []byte, doc json.RawMessage) error {
	// Note that combineRight obtains a lock on `c.combinerMu`, but it's not held
	// while we WriteLoad to the connector (which could block).
	// This allows for a concurrent handling of a Loaded response.

	if load, err := c.combineRight(binding, packedKey, doc); err != nil {
		return err
	} else if !load {
		// No-op.
	} else if err = WriteLoad(c.client, &c.txRequest, binding, packedKey); err != nil {
		return c.writeErr(err)
	}

	return nil
}

func (c *TxnClient) combineRight(binding int, packedKey []byte, doc json.RawMessage) (bool, error) {
	c.combinersMu.Lock()
	defer c.combinersMu.Unlock()

	var flighted = c.flighted[binding]
	var combiner = c.combiners[binding]
	var deltaUpdates = c.spec.Bindings[binding].DeltaUpdates
	var load bool

	if doc, ok := flighted[string(packedKey)]; ok && doc == nil {
		// We've already seen this key within this transaction.
	} else if ok {
		// We retained this document from the last transaction.
		if deltaUpdates {
			panic("we shouldn't have retained if deltaUpdates")
		}
		if err := combiner.ReduceLeft(doc); err != nil {
			return false, fmt.Errorf("combiner.ReduceLeft: %w", err)
		}
		flighted[string(packedKey)] = nil // Clear old value & mark as visited.
	} else {
		// This is a novel key.
		load = !deltaUpdates
		flighted[string(packedKey)] = nil // Mark as visited.
	}

	if err := combiner.CombineRight(doc); err != nil {
		return false, fmt.Errorf("combiner.CombineRight: %w", err)
	}

	return load, nil
}

func (c *TxnClient) Flush(deprecatedRuntimeCP pf.Checkpoint) error {
	if err := WriteFlush(c.client, &c.txRequest, deprecatedRuntimeCP); err != nil {
		return c.writeErr(err)
	} else if err = c.opLoaded.Err(); err != nil {
		return err
	}
	c.opLoaded = nil // readAcknowledgedAndLoaded has completed.

	var err error
	if c.deprecatedDriverCP, err = ReadFlushed(&c.rxResponse); err != nil {
		return err
	}
	return nil
}

func (c *TxnClient) Store() ([]*pf.CombineAPI_Stats, error) {
	// Any remaining flighted keys *not* having `nil` values are retained documents
	// of a prior transaction which were not updated during this one.
	// We garbage collect them here, and achieve the drainBinding() precondition that
	// flighted maps hold only keys of the current transaction with `nil` sentinels.
	for _, flighted := range c.flighted {
		for key, doc := range flighted {
			if doc != nil {
				delete(flighted, key)
			}
		}
	}

	var allStats = make([]*pf.CombineAPI_Stats, 0, len(c.combiners))
	// Drain each binding.
	for i, combiner := range c.combiners {
		if stats, err := c.drainBinding(
			c.flighted[i],
			combiner,
			c.spec.Bindings[i].DeltaUpdates,
			i,
		); err != nil {
			return nil, err
		} else {
			allStats = append(allStats, stats)
		}
	}

	return allStats, nil
}

func (c *TxnClient) StartCommit(runtimeCP pf.Checkpoint) (_ pf.DriverCheckpoint, acknowledged client.OpFuture, _ error) {
	var driverCP pf.DriverCheckpoint

	if err := WriteStartCommit(c.client, &c.txRequest, runtimeCP); err != nil {
		return pf.DriverCheckpoint{}, nil, c.writeErr(err)
	} else if driverCP, err = ReadStartedCommit(c.client, &c.rxResponse); err != nil {
		return pf.DriverCheckpoint{}, nil, err
	}

	// TODO(johnny): remove when Flush can no longer contain a driver checkpoint.
	if len(driverCP.DriverCheckpointJson) == 0 {
		driverCP = c.deprecatedDriverCP
	}

	c.opLoaded = client.NewAsyncOperation()
	var opAcknowledged = client.NewAsyncOperation()
	go c.readAcknowledgedAndLoaded(opAcknowledged)

	return driverCP, opAcknowledged, nil
}

func (c *TxnClient) Acknowledge() error {
	if err := WriteAcknowledge(c.client, &c.txRequest); err != nil {
		return c.writeErr(err)
	}
	return nil
}

func (c *TxnClient) writeErr(err error) error {
	// EOF indicates a stream break, which returns a causal error only with RecvMsg.
	if err != io.EOF {
		return err
	}
	// If opLoaded != nil then readAcknowledgedAndLoaded is running.
	// It will (or has) read an error, and we should wait for it.
	if c.opLoaded != nil {
		return c.opLoaded.Err()
	}
	// Otherwise we must synchronously read the error.
	for {
		if _, err = c.client.Recv(); err != nil {
			return err
		}
	}
}

// drainBinding drains the Combiner of the specified materialization
// binding by sending Store requests for its reduced documents.
func (c *TxnClient) drainBinding(
	flighted map[string]json.RawMessage,
	combiner pf.Combiner,
	deltaUpdates bool,
	binding int,
) (*pf.CombineAPI_Stats, error) {
	// Precondition: |flighted| contains the precise set of keys for this binding in this transaction.
	var remaining = len(flighted)

	// Drain the combiner into materialization Store requests.
	var stats, err = combiner.Drain(func(full bool, docRaw json.RawMessage, packedKey, packedValues []byte) error {
		// Inlined use of string(packedKey) clues compiler escape analysis to avoid allocation.
		if _, ok := flighted[string(packedKey)]; !ok {
			var key, _ = tuple.Unpack(packedKey)
			return fmt.Errorf(
				"driver implementation error: "+
					"loaded key %v (rawKey: %q) was not requested by Flow in this transaction (document %s)",
				key,
				string(packedKey),
				string(docRaw),
			)
		}

		// We're using |full|, an indicator of whether the document was a full
		// reduction or a partial combine, to track whether the document exists
		// in the store. This works because we only issue reduce-left when a
		// document was provided by Loaded or was retained from a previous
		// transaction's Store.

		if err := WriteStore(c.client, &c.txRequest, binding, packedKey, packedValues, docRaw, full); err != nil {
			return c.writeErr(err)
		}

		// We can retain a bounded number of documents from this transaction
		// as a performance optimization, so that they may be directly available
		// to the next transaction without issuing a Load.
		if deltaUpdates || remaining >= cachedDocumentBound {
			delete(flighted, string(packedKey)) // Don't retain.
		} else {
			// We cannot reference |rawDoc| beyond this callback, and must copy.
			// Fortunately, StageStore did just that, appending the document
			// to the staged request Arena, which we can reference here because
			// Arena bytes are write-once.
			var s = c.txRequest.Store
			flighted[string(packedKey)] = s.Arena.Bytes(s.DocsJson[len(s.DocsJson)-1])
		}

		remaining--
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("combine.Finish: %w", err)
	}

	// We should have seen 1:1 combined documents for each flighted key.
	if remaining != 0 {
		logrus.WithFields(logrus.Fields{
			"remaining": remaining,
			"flighted":  len(flighted),
		}).Panic("combiner drained, but expected documents remainder != 0")
	}

	return stats, nil
}

func (c *TxnClient) readAcknowledgedAndLoaded(opAcknowledged *client.AsyncOperation) (__err error) {
	defer func() {
		if opAcknowledged != nil {
			opAcknowledged.Resolve(__err)
		}
		c.opLoaded.Resolve(__err)
	}()

	if err := ReadAcknowledged(c.client, &c.rxResponse); err != nil {
		return err
	}

	opAcknowledged.Resolve(nil)
	opAcknowledged = nil // Don't resolve again.

	c.combinersMu.Lock()
	defer c.combinersMu.Unlock()

	for {
		c.combinersMu.Unlock()
		var loaded, err = ReadLoaded(c.client, &c.rxResponse)
		c.combinersMu.Lock()

		if err != nil {
			return err
		} else if loaded == nil {
			return nil
		}

		if int(loaded.Binding) > len(c.combiners) {
			return fmt.Errorf("driver implementation error (binding %d out of range)", loaded.Binding)
		}

		// Feed documents into the combiner as reduce-left operations.
		var combiner = c.combiners[loaded.Binding]
		for _, slice := range loaded.DocsJson {
			if err := combiner.ReduceLeft(loaded.Arena.Bytes(slice)); err != nil {
				return fmt.Errorf("combiner.ReduceLeft: %w", err)
			}
		}
	}
}

// TODO(johnny): This is an interesting knob we may want expose.
const cachedDocumentBound = 2048
