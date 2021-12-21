package runtime

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/estuary/flow/go/bindings"
	"github.com/estuary/flow/go/flow"
	"github.com/estuary/flow/go/shuffle"
	"github.com/estuary/protocols/fdb/tuple"
	pf "github.com/estuary/protocols/flow"
	log "github.com/sirupsen/logrus"
	"go.gazette.dev/core/broker/client"
	"go.gazette.dev/core/consumer"
	"go.gazette.dev/core/consumer/recoverylog"
	"go.gazette.dev/core/message"
)

// Derive is a top-level Application which implements the derivation workflow.
type Derive struct {
	// Derive binding that's used for the life of the derivation shard.
	binding *bindings.Derive
	// Coordinator of shuffled reads for this derivation shard.
	coordinator *shuffle.Coordinator
	// FlowConsumer which owns this Derive shard.
	host *FlowConsumer
	// Instrumented RocksDB recorder.
	recorder *recoverylog.Recorder
	// Active derivation specification, updated in RestoreCheckpoint.
	// This is duplicated from the task term to avoid needing type assertions on each usage.
	derivation *pf.DerivationSpec
	// Timestamp corresponding to the beginning of the current transaction
	txnOpened time.Time
	// Embedded processing state scoped to a current task version.
	// Updated in RestoreCheckpoint.
	taskTerm
	// Embedded task reader scoped to current task revision.
	// Also updated in RestoreCheckpoint.
	taskReader
}

var _ Application = (*Derive)(nil)

// NewDeriveApp builds and returns a *Derive Application.
func NewDeriveApp(host *FlowConsumer, shard consumer.Shard, recorder *recoverylog.Recorder) (*Derive, error) {
	var coordinator = shuffle.NewCoordinator(shard.Context(), shard.JournalClient(), host.Builds)

	var derive = &Derive{
		binding:     nil, // Lazily initialized.
		coordinator: coordinator,
		host:        host,
		recorder:    recorder,
		taskTerm:    taskTerm{},
		taskReader:  taskReader{},
	}
	return derive, nil
}

// RestoreCheckpoint initializes a processing term for the derivation,
// configures the API binding delegate, and restores the last checkpoint.
// It implements the consumer.Store interface.
func (d *Derive) RestoreCheckpoint(shard consumer.Shard) (cp pf.Checkpoint, err error) {
	if err = d.initTerm(shard, d.host); err != nil {
		return pf.Checkpoint{}, err
	}
	d.derivation = d.taskTerm.task.(*pf.DerivationSpec)

	defer func() {
		if err == nil {
			d.Log(log.DebugLevel, log.Fields{
				"derivation": d.labels.TaskName,
				"shard":      d.shardSpec.Id,
				"build":      d.labels.Build,
				"checkpoint": cp,
			}, "initialized processing term")

		} else {
			d.Log(log.ErrorLevel, log.Fields{
				"error": err.Error(),
			}, "failed to initialize processing term")
		}
	}()

	if err = d.initReader(&d.taskTerm, shard, d.derivation.TaskShuffles(), d.host); err != nil {
		return pf.Checkpoint{}, err
	}

	tsClient, err := d.build.TypeScriptClient()
	if err != nil {
		return pf.Checkpoint{}, fmt.Errorf("building TypeScript client: %w", err)
	}

	if d.binding != nil {
		// No-op.
	} else if d.binding, err = bindings.NewDerive(d.recorder, d.recorder.Dir(), d.LogPublisher); err != nil {
		return pf.Checkpoint{}, fmt.Errorf("creating derive service: %w", err)
	}

	err = d.binding.Configure(shard.FQN(), d.schemaIndex, d.derivation, tsClient)
	if err != nil {
		return pf.Checkpoint{}, fmt.Errorf("configuring derive API: %w", err)
	}

	cp, err = d.binding.RestoreCheckpoint()
	return cp, err
}

// Destroy releases the API binding delegate, which also cleans up the associated
// Rust-held RocksDB and its files.
func (d *Derive) Destroy() {
	d.taskTerm.destroy()
	// binding could be nil if there was a failure during initialization
	if d.binding != nil {
		d.binding.Destroy()
	}
}

// BeginTxn begins a derive transaction.
func (d *Derive) BeginTxn(shard consumer.Shard) error {
	d.txnOpened = time.Now().UTC()
	d.binding.BeginTxn()
	return nil
}

// ConsumeMessage passes the message to the derive worker.
func (d *Derive) ConsumeMessage(_ consumer.Shard, env message.Envelope, _ *message.Publisher) error {
	var doc = env.Message.(pf.IndexedShuffleResponse)
	var uuid = doc.UuidParts[doc.Index]

	for i := range d.derivation.Transforms {
		// Find *Shuffle with equal pointer.
		if &d.derivation.Transforms[i].Shuffle == doc.Shuffle {
			return d.binding.Add(
				uuid,
				doc.Arena.Bytes(doc.PackedKey[doc.Index]),
				uint32(i),
				doc.Arena.Bytes(doc.DocsJson[doc.Index]),
			)
		}
	}
	panic("matching shuffle not found")
}

// FinalizeTxn finishes and drains the derive worker transaction,
// and publishes each combined document to the derived collection.
func (d *Derive) FinalizeTxn(shard consumer.Shard, pub *message.Publisher) error {
	var mapper = flow.NewMapper(shard.Context(), d.host.Service.Etcd, d.host.Journals, shard.FQN())
	var collection = &d.derivation.Collection

	var stats, err = d.binding.Drain(func(full bool, doc json.RawMessage, packedKey, packedPartitions []byte) error {
		if full {
			panic("derivation produces only partially combined documents")
		}

		partitions, err := tuple.Unpack(packedPartitions)
		if err != nil {
			return fmt.Errorf("unpacking partitions: %w", err)
		}
		_, err = pub.PublishUncommitted(mapper.Map, flow.Mappable{
			Spec:       collection,
			Doc:        doc,
			PackedKey:  packedKey,
			Partitions: partitions,
		})
		return err
	})
	if err != nil {
		return err
	}
	var statsEvent = d.deriveStats(stats)
	var statsMessage = d.StatsFormatter.FormatEvent(statsEvent)
	if _, err := pub.PublishUncommitted(mapper.Map, statsMessage); err != nil {
		return fmt.Errorf("publishing stats document: %w", err)
	}
	return nil
}

func (d *Derive) deriveStats(txnStats *pf.DeriveAPI_Stats) StatsEvent {
	// assert that our task is a derivation and panic if not.
	var tfStats = make(map[string]DeriveTransformStats, len(txnStats.Transforms))
	// Only output register stats if at least one participating transform has an update lambda. This
	// allows for distinguishing between transforms where no update was invoked (Register stats will
	// be omitted) and transforms where the update lambda happened to only update existing registers
	// (Created will be 0).
	var includesUpdate = false
	for i, tf := range txnStats.Transforms {
		// Don't include transforms that didn't participate in this transaction.
		if tf != nil && tf.Input != nil {
			var tfSpec = d.derivation.Transforms[i]
			var stats = DeriveTransformStats{
				Input: docsAndBytesFromProto(tf.Input),
			}
			if tfSpec.UpdateLambda != nil {
				includesUpdate = true
				stats.Update = &InvokeStats{
					Out:          docsAndBytesFromProto(tf.Update.Output),
					SecondsTotal: tf.Update.TotalSeconds,
				}
			}
			if tfSpec.PublishLambda != nil {
				stats.Publish = &InvokeStats{
					Out:          docsAndBytesFromProto(tf.Publish.Output),
					SecondsTotal: tf.Publish.TotalSeconds,
				}
			}
			tfStats[tfSpec.Transform.String()] = stats
		}
	}
	var event = d.taskTerm.StatsFormatter.NewEvent(d.txnOpened)
	event.Derive = &DeriveStats{
		Transforms: tfStats,
		Out:        docsAndBytesFromProto(txnStats.Output),
	}
	if includesUpdate {
		event.Derive.Registers = &DeriveRegisterStats{
			CreatedTotal: txnStats.Registers.Created,
		}
	}
	return event
}

// StartCommit implements the Store interface, and writes the current transaction
// as an atomic RocksDB WriteBatch, guarded by a write barrier.
func (d *Derive) StartCommit(_ consumer.Shard, cp pf.Checkpoint, waitFor client.OpFutures) client.OpFuture {
	d.Log(log.DebugLevel, log.Fields{
		"derivation": d.labels.TaskName,
		"shard":      d.shardSpec.Id,
		"build":      d.labels.Build,
		"checkpoint": cp,
	}, "StartCommit")

	// Install a barrier such that we don't begin writing until |waitFor| has resolved.
	_ = d.recorder.Barrier(waitFor)

	// Ask the worker to apply its rocks WriteBatch, with our marshalled Checkpoint.
	if err := d.binding.PrepareCommit(cp); err != nil {
		return client.FinishedOperation(err)
	}
	// Another barrier which notifies when the WriteBatch
	// has been durably recorded to the recovery log.
	return d.recorder.Barrier(nil)
}

// FinishedTxn logs if an error occurred.
func (d *Derive) FinishedTxn(_ consumer.Shard, op consumer.OpFuture) {
	logTxnFinished(d.LogPublisher, op)
}

// Coordinator returns the shard's *shuffle.Coordinator.
func (d *Derive) Coordinator() *shuffle.Coordinator { return d.coordinator }

// ClearRegistersForTest delegates the request to its worker.
func (d *Derive) ClearRegistersForTest() error {
	return d.binding.ClearRegisters()
}
