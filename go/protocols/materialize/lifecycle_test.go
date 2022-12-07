package materialize

import (
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/estuary/flow/go/protocols/fdb/tuple"
	pf "github.com/estuary/flow/go/protocols/flow"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

//go:generate flowctl-go api build --build-id temp.db --directory testdata/ --source testdata/flow.yaml
//go:generate sqlite3 file:testdata/temp.db "SELECT WRITEFILE('testdata/materialization.proto', spec) FROM built_materializations WHERE materialization = 'test/sqlite';"

func TestStreamLifecycle(t *testing.T) {
	var specBytes, err = ioutil.ReadFile("testdata/materialization.proto")
	require.NoError(t, err)
	var spec pf.MaterializationSpec
	require.NoError(t, spec.Unmarshal(specBytes))

	var stream = new(stream)
	var srvRPC = &srvStream{stream: stream}
	var cliRPC = &clientStream{stream: stream}

	txRequest, err := WriteOpen(cliRPC, &TransactionRequest_Open{
		Materialization: &spec,
		Version:         "someVersion",
	})
	require.NoError(t, err)

	rxRequest, err := ReadOpen(srvRPC)
	require.NoError(t, err)

	txResponse, err := WriteOpened(srvRPC, &TransactionResponse_Opened{
		RuntimeCheckpoint: []byte(`recovered-runtime-checkpoint`),
	})
	require.NoError(t, err)

	rxResponse, err := ReadOpened(cliRPC)
	require.NoError(t, err)

	// Write Acknowledge and read Acknowledged.
	require.NoError(t, WriteAcknowledge(cliRPC, &txRequest))
	require.NoError(t, ReadAcknowledge(srvRPC, &rxRequest))

	require.NoError(t, WriteAcknowledged(srvRPC, &txResponse))
	require.NoError(t, ReadAcknowledged(cliRPC, &rxResponse))

	// Runtime sends multiple Loads, then Flush.
	require.NoError(t, WriteLoad(cliRPC, &txRequest, 0, tuple.Tuple{"key-1"}.Pack()))
	require.NoError(t, WriteLoad(cliRPC, &txRequest, 1, tuple.Tuple{2}.Pack()))
	require.NoError(t, WriteLoad(cliRPC, &txRequest, 1, tuple.Tuple{-3}.Pack()))
	require.NoError(t, WriteLoad(cliRPC, &txRequest, 1, tuple.Tuple{"four"}.Pack()))
	require.NoError(t, WriteLoad(cliRPC, &txRequest, 0, tuple.Tuple{[]byte("five")}.Pack()))
	require.NoError(t, WriteFlush(cliRPC, &txRequest,
		// Deprecated checkpoint, to be removed.
		pf.Checkpoint{AckIntents: map[pf.Journal][]byte{"deprecated": nil}}))

	var requireLoad = func(load Load, ok bool, err error, binding int, expect tuple.Tuple) {
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, binding, load.Binding)
		require.Equal(t, expect, load.Key)
	}

	// Driver reads Loads.
	load, ok, err := ReadLoad(srvRPC, &rxRequest)
	requireLoad(load, ok, err, 0, tuple.Tuple{"key-1"})
	load, ok, err = ReadLoad(srvRPC, &rxRequest)
	requireLoad(load, ok, err, 1, tuple.Tuple{int64(2)})
	load, ok, err = ReadLoad(srvRPC, &rxRequest)
	requireLoad(load, ok, err, 1, tuple.Tuple{int64(-3)})
	load, ok, err = ReadLoad(srvRPC, &rxRequest)
	requireLoad(load, ok, err, 1, tuple.Tuple{"four"})
	load, ok, err = ReadLoad(srvRPC, &rxRequest)
	requireLoad(load, ok, err, 0, tuple.Tuple{[]byte("five")})
	_, ok, err = ReadLoad(srvRPC, &rxRequest)
	require.False(t, ok)
	require.Nil(t, err)

	// Driver responds with Loaded, then Flushed.
	require.NoError(t, WriteLoaded(srvRPC, &txResponse, 0, []byte(`loaded-1`)))
	require.NoError(t, WriteLoaded(srvRPC, &txResponse, 0, []byte(`loaded-2`)))
	require.NoError(t, WriteLoaded(srvRPC, &txResponse, 2, []byte(`loaded-3`)))
	require.NoError(t, WriteFlushed(srvRPC, &txResponse))

	// Runtime reads Loaded.
	loaded, err := ReadLoaded(cliRPC, &rxResponse)
	require.NoError(t, err)
	require.Equal(t, 0, int(loaded.Binding))
	require.Equal(t, 2, len(loaded.DocsJson))
	loaded, err = ReadLoaded(cliRPC, &rxResponse)
	require.NoError(t, err)
	require.Equal(t, 2, int(loaded.Binding))
	require.Equal(t, 1, len(loaded.DocsJson))
	loaded, err = ReadLoaded(cliRPC, &rxResponse)
	require.NoError(t, err)
	require.Nil(t, loaded) // Indicates end of Loaded responses.

	// Runtime reads Flushed.
	_, err = ReadFlushed(&rxResponse)
	require.NoError(t, err)

	// Runtime sends Store, then StartCommit with runtime checkpoint.
	require.NoError(t, WriteStore(cliRPC, &txRequest,
		0, tuple.Tuple{"key-1"}.Pack(), tuple.Tuple{false}.Pack(), []byte(`doc-1`), true))
	require.NoError(t, WriteStore(cliRPC, &txRequest,
		0, tuple.Tuple{"key", 2}.Pack(), tuple.Tuple{"two"}.Pack(), []byte(`doc-2`), false))
	require.NoError(t, WriteStore(cliRPC, &txRequest,
		1, tuple.Tuple{"three"}.Pack(), tuple.Tuple{true}.Pack(), []byte(`doc-3`), true))
	require.NoError(t, WriteStartCommit(cliRPC, &txRequest, pf.Checkpoint{
		AckIntents: map[pf.Journal][]byte{"a-checkpoint": nil}}))

	var requireStore = func(store Store, ok bool, err error, binding int, key, values tuple.Tuple, raw string, exists bool) {
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, binding, store.Binding)
		require.Equal(t, key, store.Key)
		require.Equal(t, values, store.Values)
		require.Equal(t, raw, string(store.RawJSON))
		require.Equal(t, exists, store.Exists)
	}

	// Driver reads stores.
	store, ok, err := ReadStore(srvRPC, &rxRequest)
	requireStore(store, ok, err, 0, tuple.Tuple{"key-1"}, tuple.Tuple{false}, "doc-1", true)
	store, ok, err = ReadStore(srvRPC, &rxRequest)
	requireStore(store, ok, err, 0, tuple.Tuple{"key", int64(2)}, tuple.Tuple{"two"}, "doc-2", false)
	store, ok, err = ReadStore(srvRPC, &rxRequest)
	requireStore(store, ok, err, 1, tuple.Tuple{"three"}, tuple.Tuple{true}, "doc-3", true)
	_, ok, err = ReadStore(srvRPC, &rxRequest)
	require.False(t, ok)
	require.Nil(t, err)

	// Driver reads StartCommit.
	runtimeCP, err := ReadStartCommit(&rxRequest)
	require.NoError(t, err)
	require.NotEmpty(t, runtimeCP)

	// Driver sends StartedCommit.
	require.NoError(t, WriteStartedCommit(srvRPC, &txResponse,
		pf.DriverCheckpoint{DriverCheckpointJson: []byte(`checkpoint`)}))

	// Runtime reads StartedCommit.
	driverCP, err := ReadStartedCommit(cliRPC, &rxResponse)
	require.NoError(t, err)
	require.Equal(t, "checkpoint", string(driverCP.DriverCheckpointJson))

	// Write Acknowledge and read Acknowledged.
	require.NoError(t, WriteAcknowledge(cliRPC, &txRequest))
	require.NoError(t, ReadAcknowledge(srvRPC, &rxRequest))

	require.NoError(t, WriteAcknowledged(srvRPC, &txResponse))
	require.NoError(t, ReadAcknowledged(cliRPC, &rxResponse))

	// Snapshot to verify driver responses.
	cupaloy.SnapshotT(t, stream.req, stream.resp)
}

type stream struct {
	reqInd  int
	req     []TransactionRequest
	respInd int
	resp    []TransactionResponse
}

func (s stream) Context() context.Context { return context.Background() }

type clientStream struct{ *stream }
type srvStream struct{ *stream }

func (s *clientStream) Send(r *TransactionRequest) error {
	s.req = append(s.req, *r)
	return nil
}

func (s *srvStream) Send(r *TransactionResponse) error {
	s.resp = append(s.resp, *r)
	return nil
}

func (s *clientStream) RecvMsg(out interface{}) error {
	if len(s.resp) == s.respInd {
		return io.EOF
	}

	*out.(*TransactionResponse) = *proto.Clone(&s.resp[s.respInd]).(*TransactionResponse)
	s.respInd += 1
	return nil
}

func (s *srvStream) RecvMsg(out interface{}) error {
	if len(s.req) == s.reqInd {
		return io.EOF
	}

	*out.(*TransactionRequest) = *proto.Clone(&s.req[s.reqInd]).(*TransactionRequest)
	s.reqInd += 1
	return nil
}
