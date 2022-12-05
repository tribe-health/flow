package materialize

import (
	"context"
	"io"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/estuary/flow/go/protocols/fdb/tuple"
	pf "github.com/estuary/flow/go/protocols/flow"
	"github.com/stretchr/testify/require"
)

func TestStreamLifecycle(t *testing.T) {
	var stream = new(stream)
	var recvFn = &srvStream{stream: stream}
	var sendFn = &clientStream{stream: stream}

	var staged *TransactionRequest
	var staged2 *TransactionResponse

	// Runtime sends Acknowledge.
	require.NoError(t, WriteAcknowledge(sendFn, &staged))
	// Driver responds with Acknowledged.
	require.NoError(t, WriteAcknowledged(recvFn, &staged2))

	// Runtime sends multiple Loads, then Flush.
	require.NoError(t, WriteAcknowledge(sendFn, &staged))
	require.NoError(t, StageLoad(sendFn, &staged, 0, tuple.Tuple{"key-1"}.Pack()))
	require.NoError(t, StageLoad(sendFn, &staged, 1, tuple.Tuple{2}.Pack()))
	require.NoError(t, StageLoad(sendFn, &staged, 1, tuple.Tuple{-3}.Pack()))
	require.NoError(t, StageLoad(sendFn, &staged, 1, tuple.Tuple{"four"}.Pack()))
	require.NoError(t, StageLoad(sendFn, &staged, 3, tuple.Tuple{[]byte("five")}.Pack()))
	require.NoError(t, WriteFlush(sendFn, &staged,
		// Deprecated checkpoint, to be removed.
		pf.Checkpoint{AckIntents: map[pf.Journal][]byte{"deprecated": nil}}))

	// Driver reads Loads.
	var it = &LoadIterator{stream: recvFn}
	require.True(t, it.Next())
	require.Equal(t, tuple.Tuple{"key-1"}, it.Key)
	require.True(t, it.Next())
	require.Equal(t, tuple.Tuple{int64(2)}, it.Key)
	require.True(t, it.Next())
	require.Equal(t, tuple.Tuple{int64(-3)}, it.Key)
	require.True(t, it.Next())
	require.Equal(t, tuple.Tuple{"four"}, it.Key)

	require.True(t, it.Next())
	require.Equal(t, tuple.Tuple{[]byte("five")}, it.Key)

	require.False(t, it.Next())
	require.Nil(t, it.Err())

	// Driver responds with Loaded, then Flushed.
	require.NoError(t, StageLoaded(recvFn, &staged2, 0, []byte(`loaded-1`)))
	require.NoError(t, StageLoaded(recvFn, &staged2, 0, []byte(`loaded-2`)))
	require.NoError(t, StageLoaded(recvFn, &staged2, 2, []byte(`loaded-3`)))
	require.NoError(t, WriteFlushed(recvFn, &staged2))

	// Runtime sends Store, then StartCommit with runtime checkpoint.
	require.NoError(t, StageStore(sendFn, &staged,
		0, tuple.Tuple{"key-1"}.Pack(), tuple.Tuple{false}.Pack(), []byte(`doc-1`), true))
	require.NoError(t, StageStore(sendFn, &staged,
		0, tuple.Tuple{"key", 2}.Pack(), tuple.Tuple{"two"}.Pack(), []byte(`doc-2`), false))
	require.NoError(t, StageStore(sendFn, &staged,
		1, tuple.Tuple{"three"}.Pack(), tuple.Tuple{true}.Pack(), []byte(`doc-3`), true))
	require.NoError(t, WriteStartCommit(sendFn, &staged, pf.Checkpoint{
		AckIntents: map[pf.Journal][]byte{"foo": nil}}))

	// Driver reads stores.
	var sit = &StoreIterator{stream: recvFn}
	require.True(t, sit.Next())
	require.Equal(t, 0, sit.Binding)
	require.Equal(t, tuple.Tuple{"key-1"}, sit.Key)
	require.Equal(t, tuple.Tuple{false}, sit.Values)
	require.Equal(t, []byte(`doc-1`), []byte(sit.RawJSON))
	require.Equal(t, true, sit.Exists)

	require.True(t, sit.Next())
	require.Equal(t, 0, sit.Binding)
	require.Equal(t, tuple.Tuple{"key", int64(2)}, sit.Key)
	require.Equal(t, tuple.Tuple{"two"}, sit.Values)
	require.Equal(t, []byte(`doc-2`), []byte(sit.RawJSON))
	require.Equal(t, false, sit.Exists)

	require.True(t, sit.Next())
	require.Equal(t, 1, sit.Binding)
	require.Equal(t, tuple.Tuple{"three"}, sit.Key)
	require.Equal(t, tuple.Tuple{true}, sit.Values)
	require.Equal(t, []byte(`doc-3`), []byte(sit.RawJSON))
	require.Equal(t, true, sit.Exists)

	require.False(t, sit.Next())
	require.Nil(t, sit.Err())
	require.NotEmpty(t, sit.StartCommit().FlowCheckpoint)

	// Driver sends StartedCommit.
	require.NoError(t, WriteStartedCommit(recvFn, &staged2,
		pf.DriverCheckpoint{DriverCheckpointJson: []byte(`checkpoint`)}))

	// Snapshot to verify driver responses.
	cupaloy.SnapshotT(t, stream.resp)
}

type stream struct {
	req  []*TransactionRequest
	resp []*TransactionResponse
}

func (s stream) Context() context.Context { return context.Background() }

type clientStream struct{ *stream }
type srvStream struct{ *stream }

func (s *clientStream) Send(r *TransactionRequest) error {
	s.req = append(s.req, r)
	return nil
}

func (s *srvStream) Send(r *TransactionResponse) error {
	s.resp = append(s.resp, r)
	return nil
}

func (s *clientStream) Recv() (*TransactionResponse, error) {
	if len(s.resp) == 0 {
		return nil, io.EOF
	}

	var r = s.resp[0]
	s.resp = s.resp[1:]
	return r, nil
}

func (s *srvStream) Recv() (*TransactionRequest, error) {
	if len(s.req) == 0 {
		return nil, io.EOF
	}

	var r = s.req[0]
	s.req = s.req[1:]
	return r, nil
}

func (s *srvStream) RecvMsg(out interface{}) error {
	if r, err := s.Recv(); err != nil {
		return err
	} else {
		*out.(*TransactionRequest) = *r
		return nil
	}
}
