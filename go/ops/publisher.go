package ops

import (
	"fmt"

	"github.com/estuary/flow/go/labels"
)

// Publisher of operation Logs and Stats.
// TODO(johnny): Publisher covers ops.Logs, but does not yet
// cover ops.Stats.
type Publisher interface {
	// PublishLog publishes a Log instance.
	PublishLog(Log)
	// Labels which are the context of this Publisher.
	Labels() labels.ShardLabeling
}

// ShardRef is a reference to a specific task shard that produced logs and stats.
// * ops-catalog/ops-task-schema.json
// * crate/ops/lib.rs
type ShardRef struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	KeyBegin    string `json:"keyBegin"`
	RClockBegin string `json:"rClockBegin"`
}

func NewShardRef(labeling labels.ShardLabeling) ShardRef {
	return ShardRef{
		Name:        labeling.TaskName,
		Kind:        labeling.TaskType,
		KeyBegin:    fmt.Sprintf("%08x", labeling.Range.KeyBegin),
		RClockBegin: fmt.Sprintf("%08x", labeling.Range.RClockBegin),
	}
}
