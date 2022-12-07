package materialize

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/estuary/flow/go/protocols/fdb/tuple"
	"github.com/estuary/flow/go/protocols/flow"
	pf "github.com/estuary/flow/go/protocols/flow"
	pc "go.gazette.dev/core/consumer/protocol"
)

// Protocol routines for sending TransactionRequest follow:

type TransactionRequestTx interface {
	Send(*TransactionRequest) error
}

func WriteOpen(stream TransactionRequestTx, open *TransactionRequest_Open) (TransactionRequest, error) {
	var request = TransactionRequest{Open: open}

	if err := stream.Send(&request); err != nil {
		return TransactionRequest{}, fmt.Errorf("sending Open: %w", err)
	}
	return request, nil
}

func WriteAcknowledge(stream TransactionRequestTx, request *TransactionRequest) error {
	if request.Open == nil && request.StartCommit == nil {
		panic(fmt.Sprintf("expected prior request is Open or StartCommit, got %#v", request))
	}
	*request = TransactionRequest{
		Acknowledge: &TransactionRequest_Acknowledge{},
	}
	if err := stream.Send(request); err != nil {
		return fmt.Errorf("sending Acknowledge request: %w", err)
	}
	return nil
}

func WriteLoad(
	stream TransactionRequestTx,
	request *TransactionRequest,
	binding int,
	packedKey []byte,
) error {
	if request.Acknowledge == nil && request.Load == nil {
		panic(fmt.Sprintf("expected prior request is Acknowledge or Load, got %#v", request))
	}

	// Send current `request` if it uses a different binding or would re-allocate.
	if request.Load != nil {
		var rem int
		if l := (*request).Load; int(l.Binding) != binding {
			rem = -1 // Must flush this request.
		} else if cap(l.PackedKeys) != len(l.PackedKeys) {
			rem = cap(l.Arena) - len(l.Arena)
		}
		if rem < len(packedKey) {
			if err := stream.Send(request); err != nil {
				return fmt.Errorf("sending Load request: %w", err)
			}
			request.Load = nil
		}
	}

	if request.Load == nil {
		*request = TransactionRequest{
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

func WriteFlush(
	stream TransactionRequestTx,
	request *TransactionRequest,
	deprecatedCheckpoint pc.Checkpoint, // Will be removed.
) error {
	if request.Acknowledge == nil && request.Load == nil {
		panic(fmt.Sprintf("expected prior request is Acknowledge or Load, got %#v", request))
	}
	// Flush partial Load request, if required.
	if request.Load != nil {
		if err := stream.Send(request); err != nil {
			return fmt.Errorf("flushing final Load request: %w", err)
		}
		*request = TransactionRequest{}
	}

	var checkpointBytes, err = deprecatedCheckpoint.Marshal()
	if err != nil {
		panic(err) // Cannot fail.
	}
	*request = TransactionRequest{
		Flush: &TransactionRequest_Flush{
			DeprecatedRuntimeCheckpoint: checkpointBytes,
		},
	}
	if err := stream.Send(request); err != nil {
		return fmt.Errorf("sending Flush request: %w", err)
	}
	return nil
}

func WriteStore(
	stream TransactionRequestTx,
	request *TransactionRequest,
	binding int,
	packedKey []byte,
	packedValues []byte,
	doc json.RawMessage,
	exists bool,
) error {
	if request.Flush == nil && request.Store == nil {
		panic(fmt.Sprintf("expected prior request is Flush or Store, got %#v", request))
	}

	// Send current |request| if we would re-allocate.
	if request.Store != nil {
		var rem int
		if s := (*request).Store; int(s.Binding) != binding {
			rem = -1 // Must flush this request.
		} else if cap(s.PackedKeys) != len(s.PackedKeys) {
			rem = cap(s.Arena) - len(s.Arena)
		}
		if need := len(packedKey) + len(packedValues) + len(doc); need > rem {
			if err := stream.Send(request); err != nil {
				return fmt.Errorf("sending Store request: %w", err)
			}
			*request = TransactionRequest{}
		}
	}

	if request.Store == nil {
		*request = TransactionRequest{
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

func WriteStartCommit(
	stream TransactionRequestTx,
	request *TransactionRequest,
	checkpoint pc.Checkpoint,
) error {
	if request.Flush == nil && request.Store == nil {
		panic(fmt.Sprintf("expected prior request is Flush or Store, got %#v", request))
	}
	// Flush partial Store request, if required.
	if request.Store != nil {
		if err := stream.Send(request); err != nil {
			return fmt.Errorf("flushing final Store request: %w", err)
		}
		*request = TransactionRequest{}
	}

	var checkpointBytes, err = checkpoint.Marshal()
	if err != nil {
		panic(err) // Cannot fail.
	}
	*request = TransactionRequest{
		StartCommit: &TransactionRequest_StartCommit{
			RuntimeCheckpoint: checkpointBytes,
		},
	}
	if err := stream.Send(request); err != nil {
		return fmt.Errorf("sending StartCommit request: %w", err)
	}
	return nil
}

// Protocol routines for receiving TransactionRequest follow:

type TransactionRequestRx interface {
	RecvMsg(interface{}) error
}

func ReadOpen(stream TransactionRequestRx) (TransactionRequest, error) {
	var request TransactionRequest

	if err := stream.RecvMsg(&request); err != nil {
		return TransactionRequest{}, fmt.Errorf("reading Open: %w", err)
	} else if request.Open == nil {
		return TransactionRequest{}, fmt.Errorf("protocol error (expected Open, got %#v)", request)
	} else if err = request.Validate(); err != nil {
		return TransactionRequest{}, fmt.Errorf("validating read Open: %w", err)
	}
	return request, nil
}

func ReadAcknowledge(stream TransactionRequestRx, request *TransactionRequest) error {
	if request.Open == nil && request.StartCommit == nil {
		panic(fmt.Sprintf("expected prior request is Open or StartCommit, got %#v", request))
	} else if err := stream.RecvMsg(request); err != nil {
		return fmt.Errorf("reading expected Acknowledge: %w", err)
	} else if request.Acknowledge == nil {
		return fmt.Errorf("protocol error (expected Acknowledge, got %#v)", request)
	} else if err = request.Validate(); err != nil {
		return fmt.Errorf("validating read Acknowledge: %w", err)
	}
	return nil
}

type Load struct {
	Binding int         // Binding index of this document to load.
	Key     tuple.Tuple // Key of the next document to load.
}

func ReadLoad(stream TransactionRequestRx, request *TransactionRequest) (_ Load, ok bool, _ error) {
	if request.Acknowledge == nil && request.Load == nil {
		panic(fmt.Sprintf("expected prior request is Acknowledge or Load, got %#v", request))
	}
	// Read next `Load` request from `stream`?
	if request.Acknowledge != nil || len(request.Load.PackedKeys) == 0 {
		if err := stream.RecvMsg(request); err == io.EOF && request.Acknowledge != nil {
			return Load{}, false, io.EOF // Clean shutdown.
		} else if err != nil {
			return Load{}, false, fmt.Errorf("reading Load: %w", err)
		} else if request.Load == nil {
			return Load{}, false, nil // No loads remain.
		} else if err = request.Validate(); err != nil {
			return Load{}, false, fmt.Errorf("validating read Load: %w", err)
		}
	}
	var l = request.Load

	var key, err = tuple.Unpack(l.Arena.Bytes(l.PackedKeys[0]))
	if err != nil {
		return Load{}, false, fmt.Errorf("unpacking Load key: %w", err)
	}

	l.PackedKeys = l.PackedKeys[1:]

	return Load{
		Binding: int(l.Binding),
		Key:     key,
	}, true, nil
}

func ReadFlush(request *TransactionRequest) error {
	if request.Flush == nil {
		return fmt.Errorf("protocol error (expected Flush, got %#v)", request)
	} else if err := request.Validate(); err != nil {
		return fmt.Errorf("validating read Flush: %w", err)
	}
	return nil
}

type Store struct {
	Binding int             // Binding index of this stored document.
	Key     tuple.Tuple     // Key of the next document to store.
	Values  tuple.Tuple     // Values of the next document to store.
	RawJSON json.RawMessage // Document to store.
	Exists  bool            // Does this document exist in the store already?
}

func ReadStore(stream TransactionRequestRx, request *TransactionRequest) (_ Store, ok bool, _ error) {
	if request.Flush == nil && request.Store == nil {
		panic(fmt.Sprintf("expected prior request is Flush or Store, got %#v", request))
	}
	// Read next `Store` request from `stream`?
	if request.Flush != nil || len(request.Store.PackedKeys) == 0 {
		if err := stream.RecvMsg(request); err != nil {
			return Store{}, false, fmt.Errorf("reading Store: %w", err)
		} else if request.Store == nil {
			return Store{}, false, nil // No stores remain.
		} else if err = request.Validate(); err != nil {
			return Store{}, false, fmt.Errorf("validating read Store: %w", err)
		}
	}
	var s = request.Store

	var key, err = tuple.Unpack(s.Arena.Bytes(s.PackedKeys[0]))
	if err != nil {
		return Store{}, false, fmt.Errorf("unpacking Store key: %w", err)
	}
	values, err := tuple.Unpack(s.Arena.Bytes(s.PackedValues[0]))
	if err != nil {
		return Store{}, false, fmt.Errorf("unpacking Store values: %w", err)
	}
	var rawJSON = s.Arena.Bytes(s.DocsJson[0])
	var exists = s.Exists[0]

	s.PackedKeys = s.PackedKeys[1:]
	s.PackedValues = s.PackedValues[1:]
	s.DocsJson = s.DocsJson[1:]
	s.Exists = s.Exists[1:]

	return Store{
		Binding: int(s.Binding),
		Key:     key,
		Values:  values,
		RawJSON: rawJSON,
		Exists:  exists,
	}, true, nil
}

func ReadStartCommit(request *TransactionRequest) (runtimeCheckpoint []byte, _ error) {
	if request.StartCommit == nil {
		return nil, fmt.Errorf("protocol error (expected StartCommit, got %#v)", request)
	} else if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("validating read StartCommit: %w", err)
	}
	return request.StartCommit.RuntimeCheckpoint, nil
}

// Protocol routines for sending TransactionResponse follow:

type TransactionResponseTx interface {
	Send(*TransactionResponse) error
}

func WriteOpened(stream TransactionResponseTx, opened *TransactionResponse_Opened) (TransactionResponse, error) {
	var response = TransactionResponse{Opened: opened}

	if err := stream.Send(&response); err != nil {
		return TransactionResponse{}, fmt.Errorf("sending Opened: %w", err)
	}
	return response, nil
}

func WriteAcknowledged(stream TransactionResponseTx, response *TransactionResponse) error {
	if response.Opened == nil && response.StartedCommit == nil {
		panic(fmt.Sprintf("expected prior response is Opened or StartedCommit, got %#v", response))
	}
	*response = TransactionResponse{
		Acknowledged: &TransactionResponse_Acknowledged{},
	}
	if err := stream.Send(response); err != nil {
		return fmt.Errorf("sending Acknowledged response: %w", err)
	}
	return nil
}

func WriteLoaded(
	stream TransactionResponseTx,
	response *TransactionResponse,
	binding int,
	document json.RawMessage,
) error {
	if response.Acknowledged == nil && response.Loaded == nil {
		panic(fmt.Sprintf("expected prior response is Acknowledged or Loaded, got %#v", response))
	}

	// Send current `response` if it uses a different binding or would re-allocate.
	if response.Loaded != nil {
		var rem int
		if l := (*response).Loaded; int(l.Binding) != binding {
			rem = -1 // Must flush this response.
		} else if cap(l.DocsJson) != len(l.DocsJson) {
			rem = cap(l.Arena) - len(l.Arena)
		}
		if rem < len(document) {
			if err := stream.Send(response); err != nil {
				return fmt.Errorf("sending Loaded response: %w", err)
			}
			response.Loaded = nil
		}
	}

	if response.Loaded == nil {
		*response = TransactionResponse{
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

func WriteFlushed(stream TransactionResponseTx, response *TransactionResponse) error {
	if response.Acknowledged == nil && response.Loaded == nil {
		panic(fmt.Sprintf("expected prior response is Acknowledged or Loaded, got %#v", response))
	}
	// Flush partial Loaded response, if required.
	if response.Loaded != nil {
		if err := stream.Send(response); err != nil {
			return fmt.Errorf("flushing final Loaded response: %w", err)
		}
		*response = TransactionResponse{}
	}

	*response = TransactionResponse{
		// Flushed as-a DriverCheckpoint is deprecated and will be removed.
		Flushed: &flow.DriverCheckpoint{
			DriverCheckpointJson: []byte("{}"),
			Rfc7396MergePatch:    true,
		},
	}
	if err := stream.Send(response); err != nil {
		return fmt.Errorf("sending Flushed response: %w", err)
	}
	return nil
}

func WriteStartedCommit(
	stream TransactionResponseTx,
	response *TransactionResponse,
	checkpoint pf.DriverCheckpoint,
) error {
	if response.Flushed == nil {
		panic(fmt.Sprintf("expected prior response is Flushed, got %#v", response))
	}
	*response = TransactionResponse{
		StartedCommit: &TransactionResponse_StartedCommit{
			DriverCheckpoint: &checkpoint,
		},
	}
	if err := stream.Send(response); err != nil {
		return fmt.Errorf("sending StartedCommit: %w", err)
	}
	return nil
}

// Protocol routines for reading TransactionResponse follow:

type TransactionResponseRx interface {
	RecvMsg(interface{}) error
}

func ReadOpened(stream TransactionResponseRx) (TransactionResponse, error) {
	var response TransactionResponse

	if err := stream.RecvMsg(&response); err != nil {
		return TransactionResponse{}, fmt.Errorf("reading Opened: %w", err)
	} else if response.Opened == nil {
		return TransactionResponse{}, fmt.Errorf("protocol error (expected Opened, got %#v)", response)
	} else if err = response.Validate(); err != nil {
		return TransactionResponse{}, fmt.Errorf("validating read Opened: %w", err)
	}
	return response, nil
}

func ReadAcknowledged(stream TransactionResponseRx, response *TransactionResponse) error {
	if response.Opened == nil && response.StartedCommit == nil {
		panic(fmt.Sprintf("expected prior response is Opened or StartedCommit, got %#v", response))
	} else if err := stream.RecvMsg(response); err != nil {
		return fmt.Errorf("reading expected Acknowledge: %w", err)
	} else if response.Acknowledged == nil {
		return fmt.Errorf("protocol error (expected Acknowledged, got %#v)", response)
	} else if err = response.Validate(); err != nil {
		return fmt.Errorf("validating read Acknowledged: %w", err)
	}
	return nil
}

func ReadLoaded(stream TransactionResponseRx, response *TransactionResponse) (*TransactionResponse_Loaded, error) {
	if response.Acknowledged == nil && response.Loaded == nil {
		panic(fmt.Sprintf("expected prior response is Acknowledged or Loaded, got %#v", response))
	} else if err := stream.RecvMsg(response); err != nil {
		return nil, fmt.Errorf("reading Loaded: %w", err)
	} else if response.Loaded == nil {
		return nil, nil // No loads remain.
	} else if err = response.Validate(); err != nil {
		return nil, fmt.Errorf("validating read Loaded: %w", err)
	}
	return response.Loaded, nil
}

func ReadFlushed(response *TransactionResponse) (deprecatedDriverCP pf.DriverCheckpoint, _ error) {
	if response.Flushed == nil {
		return pf.DriverCheckpoint{}, fmt.Errorf("protocol error (expected Flushed, got %#v)", response)
	} else if err := response.Validate(); err != nil {
		return pf.DriverCheckpoint{}, fmt.Errorf("validating read Flushed: %w", err)
	}
	return *response.Flushed, nil
}

func ReadStartedCommit(stream TransactionResponseRx, response *TransactionResponse) (pf.DriverCheckpoint, error) {
	if response.Flushed == nil {
		panic(fmt.Sprintf("expected prior response is Flushed, got %#v", response))
	} else if err := stream.RecvMsg(response); err != nil {
		return pf.DriverCheckpoint{}, fmt.Errorf("reading expected StartedCommit: %w", err)
	} else if response.StartedCommit == nil {
		return pf.DriverCheckpoint{}, fmt.Errorf("protocol error (expected StartedCommit, got %#v)", response)
	} else if err = response.Validate(); err != nil {
		return pf.DriverCheckpoint{}, fmt.Errorf("validating read StartedCommit: %w", err)
	}
	return *response.StartedCommit.DriverCheckpoint, nil
}

const (
	arenaSize = 16 * 1024
	sliceSize = 32
)
