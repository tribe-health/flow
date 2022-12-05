package materialize

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/estuary/flow/go/protocols/fdb/tuple"
	"github.com/estuary/flow/go/protocols/flow"
	pf "github.com/estuary/flow/go/protocols/flow"
	pc "go.gazette.dev/core/consumer/protocol"
)

// StageLoad potentially sends a previously staged Load into the stream,
// and then stages its arguments into request.Load.
func StageLoad(
	stream interface {
		Send(*TransactionRequest) error
	},
	request **TransactionRequest,
	binding int,
	packedKey []byte,
) error {
	// Send current |request| if it uses a different binding or would re-allocate.
	if *request != nil {
		var rem int
		if l := (*request).Load; int(l.Binding) != binding {
			rem = -1 // Must flush this request.
		} else if cap(l.PackedKeys) != len(l.PackedKeys) {
			rem = cap(l.Arena) - len(l.Arena)
		}
		if rem < len(packedKey) {
			if err := stream.Send(*request); err != nil {
				return fmt.Errorf("sending Load request: %w", err)
			}
			*request = nil
		}
	}

	if *request == nil {
		*request = &TransactionRequest{
			Load: &TransactionRequest_Load{
				Binding:    uint32(binding),
				Arena:      make(pf.Arena, 0, arenaSize),
				PackedKeys: make([]pf.Slice, 0, sliceSize),
			},
		}
	}

	var l = (*request).Load
	l.PackedKeys = append(l.PackedKeys, l.Arena.Add(packedKey))

	return nil
}

// StageLoaded potentially sends a previously staged Loaded into the stream,
// and then stages its arguments into response.Loaded.
func StageLoaded(
	stream interface {
		Send(*TransactionResponse) error
	},
	response **TransactionResponse,
	binding int,
	document json.RawMessage,
) error {
	// Send current |response| if we would re-allocate.
	if *response != nil {
		var rem int
		if l := (*response).Loaded; int(l.Binding) != binding {
			rem = -1 // Must flush this response.
		} else if cap(l.DocsJson) != len(l.DocsJson) {
			rem = cap(l.Arena) - len(l.Arena)
		}
		if rem < len(document) {
			if err := stream.Send(*response); err != nil {
				return fmt.Errorf("sending Loaded response: %w", err)
			}
			*response = nil
		}
	}

	if *response == nil {
		*response = &TransactionResponse{
			Loaded: &TransactionResponse_Loaded{
				Binding:  uint32(binding),
				Arena:    make(pf.Arena, 0, arenaSize),
				DocsJson: make([]pf.Slice, 0, sliceSize),
			},
		}
	}

	var l = (*response).Loaded
	l.DocsJson = append(l.DocsJson, l.Arena.Add(document))

	return nil
}

// WriteFlush flushes a pending Load request, and sends a Flush request
// with the provided (deprecated) Flow checkpoint.
func WriteFlush(
	stream interface {
		Send(*TransactionRequest) error
	},
	request **TransactionRequest,
	deprecatedCheckpoint pc.Checkpoint,
) error {
	// Flush partial Load request, if required.
	if *request != nil {
		if err := stream.Send(*request); err != nil {
			return fmt.Errorf("flushing final Load request: %w", err)
		}
		*request = nil
	}

	var checkpointBytes, err = deprecatedCheckpoint.Marshal()
	if err != nil {
		panic(err) // Cannot fail.
	}

	if err := stream.Send(&TransactionRequest{
		Flush: &TransactionRequest_Flush{
			DeprecatedFlowCheckpoint: checkpointBytes,
		},
	}); err != nil {
		return fmt.Errorf("sending Flush request: %w", err)
	}

	return nil
}

// WriteFlushed flushes a pending Loaded response, and sends a Flushed response.
func WriteFlushed(
	stream interface {
		Send(*TransactionResponse) error
	},
	response **TransactionResponse,
) error {
	// Flush partial Loaded response, if required.
	if *response != nil {
		if err := stream.Send(*response); err != nil {
			return fmt.Errorf("flushing final Loaded response: %w", err)
		}
		*response = nil
	}

	if err := stream.Send(&TransactionResponse{
		Flushed: &flow.DriverCheckpoint{},
	}); err != nil {
		return fmt.Errorf("sending Flushed response: %w", err)
	}

	return nil
}

// StageStore potentially sends a previously staged Store into the stream,
// and then stages its arguments into request.Store.
func StageStore(
	stream interface {
		Send(*TransactionRequest) error
	},
	request **TransactionRequest,
	binding int,
	packedKey []byte,
	packedValues []byte,
	doc json.RawMessage,
	exists bool,
) error {
	// Send current |request| if we would re-allocate.
	if *request != nil {
		var rem int
		if s := (*request).Store; int(s.Binding) != binding {
			rem = -1 // Must flush this request.
		} else if cap(s.PackedKeys) != len(s.PackedKeys) {
			rem = cap(s.Arena) - len(s.Arena)
		}
		if need := len(packedKey) + len(packedValues) + len(doc); need > rem {
			if err := stream.Send(*request); err != nil {
				return fmt.Errorf("sending Store request: %w", err)
			}
			*request = nil
		}
	}

	if *request == nil {
		*request = &TransactionRequest{
			Store: &TransactionRequest_Store{
				Binding:      uint32(binding),
				Arena:        make(pf.Arena, 0, arenaSize),
				PackedKeys:   make([]pf.Slice, 0, sliceSize),
				PackedValues: make([]pf.Slice, 0, sliceSize),
				DocsJson:     make([]pf.Slice, 0, sliceSize),
			},
		}
	}

	var s = (*request).Store
	s.PackedKeys = append(s.PackedKeys, s.Arena.Add(packedKey))
	s.PackedValues = append(s.PackedValues, s.Arena.Add(packedValues))
	s.DocsJson = append(s.DocsJson, s.Arena.Add(doc))
	s.Exists = append(s.Exists, exists)

	return nil
}

// WriteStartCommit flushes a pending Store request, and sends a StartCommit request.
func WriteStartCommit(
	stream interface {
		Send(*TransactionRequest) error
	},
	request **TransactionRequest,
	checkpoint pc.Checkpoint,
) error {
	// Flush partial Store request, if required.
	if *request != nil {
		if err := stream.Send(*request); err != nil {
			return fmt.Errorf("flushing final Store request: %w", err)
		}
		*request = nil
	}

	var checkpointBytes, err = checkpoint.Marshal()
	if err != nil {
		panic(err) // Cannot fail.
	}

	if err := stream.Send(&TransactionRequest{
		StartCommit: &TransactionRequest_StartCommit{
			FlowCheckpoint: checkpointBytes,
		},
	}); err != nil {
		return fmt.Errorf("sending StartCommit request: %w", err)
	}

	return nil
}

// WriteStartedCommit writes a StartedCommit response to the stream.
func WriteStartedCommit(
	stream interface {
		Send(*TransactionResponse) error
	},
	response **TransactionResponse,
	checkpoint pf.DriverCheckpoint,
) error {
	// We must have sent Prepared prior to DriverCommitted.
	if *response != nil {
		panic("expected nil response")
	}

	if err := stream.Send(&TransactionResponse{
		StartedCommit: &TransactionResponse_StartedCommit{
			DriverCheckpoint: &checkpoint,
		},
	}); err != nil {
		return fmt.Errorf("sending StartedCommit response: %w", err)
	}

	return nil
}

// WriteAcknowledge writes an Acknowledge request into the stream.
func WriteAcknowledge(
	stream interface {
		Send(*TransactionRequest) error
	},
	request **TransactionRequest,
) error {
	if *request != nil {
		panic("expected nil request")
	}

	if err := stream.Send(&TransactionRequest{
		Acknowledge: &TransactionRequest_Acknowledge{},
	}); err != nil {
		return fmt.Errorf("sending Acknowledge request: %w", err)
	}

	return nil
}

// WriteAcknowledged writes an Acknowledged response to the stream.
func WriteAcknowledged(
	stream interface {
		Send(*TransactionResponse) error
	},
	response **TransactionResponse,
) error {
	// We must have sent StartedCommit prior to Acknowledged.
	if *response != nil {
		panic("expected nil response")
	}

	if err := stream.Send(&TransactionResponse{
		Acknowledged: &TransactionResponse_Acknowledged{},
	}); err != nil {
		return fmt.Errorf("sending Acknowledged response: %w", err)
	}

	return nil
}

// LoadIterator is an iterator over Load requests.
type LoadIterator struct {
	Binding int         // Binding index of this document to load.
	Key     tuple.Tuple // Key of the next document to load.

	stream interface {
		Context() context.Context
		RecvMsg(m interface{}) error // Receives Load, Acknowledge, and Prepare.
	}
	req   TransactionRequest // Last read request.
	index int                // Last-returned document index within |req|.
	total int                // Total number of iterated keys.
	err   error              // Final error.
}

// NewLoadIterator returns a *LoadIterator of the stream.
func NewLoadIterator(stream Driver_TransactionsServer) *LoadIterator {
	return &LoadIterator{stream: stream}
}

// poll returns true if there is at least one LoadRequest message with at least one key
// remaining to be read. This will read the next message from the stream if required. This will not
// advance the iterator, so it can be used to check whether the LoadIterator contains at least one
// key to load without actually consuming the next key. If false is returned, then there are no
// remaining keys and poll must not be called again. Note that if poll returns
// true, then Next may still return false if the LoadRequest message is malformed.
func (it *LoadIterator) poll() bool {
	if it.err != nil || it.req.Flush != nil {
		panic("Poll called again after having returned false")
	}

	// Must we read another request?
	if it.req.Load == nil || it.index == len(it.req.Load.PackedKeys) {
		// Use RecvMsg to re-use |it.req| without allocation.
		// Safe because we fully process |it.req| between calls to RecvMsg.
		if err := it.stream.RecvMsg(&it.req); err == io.EOF {
			if it.total != 0 {
				it.err = fmt.Errorf("unexpected EOF when there are loaded keys")
			} else {
				it.err = io.EOF // Clean shutdown
			}
			return false
		} else if err != nil {
			it.err = fmt.Errorf("reading Load: %w", err)
			return false
		}

		if it.req.Flush != nil {
			return false // Flush ends the Load phase.
		} else if it.req.Load == nil || len(it.req.Load.PackedKeys) == 0 {
			it.err = fmt.Errorf("protocol error (expected non-empty Load, got %#v)", it.req)
			return false
		}
		it.index = 0
		it.Binding = int(it.req.Load.Binding)
	}
	return true
}

// Context returns the Context of this LoadIterator.
func (it *LoadIterator) Context() context.Context { return it.stream.Context() }

// Next returns true if there is another Load and makes it available via Key.
// When a Prepare is read, or if an error is encountered, it returns false
// and must not be called again.
func (it *LoadIterator) Next() bool {
	if !it.poll() {
		return false
	}

	var l = it.req.Load
	it.Key, it.err = tuple.Unpack(it.req.Load.Arena.Bytes(l.PackedKeys[it.index]))

	if it.err != nil {
		it.err = fmt.Errorf("unpacking Load key: %w", it.err)
		return false
	}

	it.index++
	it.total++
	return true
}

// Err returns an encountered error.
func (it *LoadIterator) Err() error {
	return it.err
}

// StoreIterator is an iterator over Store requests.
type StoreIterator struct {
	Binding int             // Binding index of this stored document.
	Key     tuple.Tuple     // Key of the next document to store.
	Values  tuple.Tuple     // Values of the next document to store.
	RawJSON json.RawMessage // Document to store.
	Exists  bool            // Does this document exist in the store already?

	stream interface {
		Context() context.Context
		RecvMsg(m interface{}) error // Receives Store and Commit.
	}
	req   TransactionRequest
	index int
	total int
	err   error
}

// NewStoreIterator returns a *StoreIterator of the stream.
func NewStoreIterator(stream Driver_TransactionsServer) *StoreIterator {
	return &StoreIterator{stream: stream}
}

// poll returns true if there is at least one StoreRequest message with at least one document
// remaining to be read. This will read the next message from the stream if required. This will not
// advance the iterator, so it can be used to check whether the StoreIterator contains at least one
// document to store without actually consuming the next document. If false is returned, then there
// are no remaining documents and poll must not be called again. Note that if poll returns true,
// then Next may still return false if the StoreRequest message is malformed.
func (it *StoreIterator) poll() bool {
	if it.err != nil || it.req.StartCommit != nil {
		panic("Poll called again after having returned false")
	}
	// Must we read another request?
	if it.req.Store == nil || it.index == len(it.req.Store.PackedKeys) {
		// Use RecvMsg to re-use |it.req| without allocation.
		// Safe because we fully process |it.req| between calls to RecvMsg.
		if err := it.stream.RecvMsg(&it.req); err != nil {
			it.err = fmt.Errorf("reading Store: %w", err)
			return false
		}

		if it.req.StartCommit != nil {
			return false // StartCommit ends the Store phase.
		} else if it.req.Store == nil || len(it.req.Store.PackedKeys) == 0 {
			it.err = fmt.Errorf("expected non-empty Store, got %#v", it.req)
			return false
		}
		it.index = 0
		it.Binding = int(it.req.Store.Binding)
	}
	return true
}

// Context returns the Context of this StoreIterator.
func (it *StoreIterator) Context() context.Context { return it.stream.Context() }

// Next returns true if there is another Store and makes it available.
// When a Commit is read, or if an error is encountered, it returns false
// and must not be called again.
func (it *StoreIterator) Next() bool {
	if !it.poll() {
		return false
	}

	var s = it.req.Store

	it.Key, it.err = tuple.Unpack(s.Arena.Bytes(s.PackedKeys[it.index]))
	if it.err != nil {
		it.err = fmt.Errorf("unpacking Store key: %w", it.err)
		return false
	}
	it.Values, it.err = tuple.Unpack(s.Arena.Bytes(s.PackedValues[it.index]))
	if it.err != nil {
		it.err = fmt.Errorf("unpacking Store values: %w", it.err)
		return false
	}
	it.RawJSON = s.Arena.Bytes(s.DocsJson[it.index])
	it.Exists = s.Exists[it.index]

	it.index++
	it.total++
	return true
}

// Err returns an encountered error.
func (it *StoreIterator) Err() error {
	return it.err
}

// StartCommit returns a TransactionRequest_StartCommit which caused this StoreIterator
// to terminate. It's valid only after Next returns false and Err is nil.
func (it *StoreIterator) StartCommit() *TransactionRequest_StartCommit {
	return it.req.StartCommit
}

const (
	arenaSize = 16 * 1024
	sliceSize = 32
)
