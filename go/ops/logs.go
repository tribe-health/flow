package ops

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	pf "github.com/estuary/flow/go/protocols/flow"
)

// Log is the canonical shape of a Flow operations Log document.
// See also:
// * ops-catalog/ops-log-schema.json
// * crate/ops/lib.rs
type Log struct {
	Meta struct {
		UUID string `json:"uuid"`
	} `json:"_meta"`
	Timestamp time.Time       `json:"ts"`
	Level     pf.LogLevel     `json:"level"`
	Message   string          `json:"message"`
	Fields    json.RawMessage `json:"fields,omitempty"`
	Shard     ShardRef        `json:"shard,omitempty"`
	Spans     []Log           `json:"spans,omitempty"`
}

// LogCollection returns the collection to which logs of the given shard are written.
func LogCollection(taskName string) pf.Collection {
	return pf.Collection(fmt.Sprintf("ops/%s/logs", strings.Split(taskName, "/")[0]))
}

// ValidateLogsCollection sanity-checks that the given CollectionSpec is appropriate
// for storing instances of Log documents.
func ValidateLogsCollection(spec *pf.CollectionSpec) error {
	if !reflect.DeepEqual(
		spec.KeyPtrs,
		[]string{"/shard/name", "/shard/keyBegin", "/shard/rClockBegin", "/ts"},
	) {
		return fmt.Errorf("CollectionSpec doesn't have expected key: %v", spec.KeyPtrs)
	}

	if !reflect.DeepEqual(spec.PartitionFields, []string{"kind", "name"}) {
		return fmt.Errorf(
			"CollectionSpec doesn't have expected partitions 'kind' & 'name': %v",
			spec.PartitionFields)
	}

	return nil
}

// PublishLog constructs and publishes a Log using the given Publisher.
// Fields must be pairs of a string key followed by a JSON-encodable interface{} value.
// PublishLog panics if `fields` are odd, or if a field isn't a string,
// or if it cannot be encoded as JSON.
func PublishLog(publisher Publisher, level pf.LogLevel, message string, fields ...interface{}) {
	if publisher.Labels().LogLevel < level {
		return
	}

	// NOTE(johnny): We panic because incorrect fields are a developer
	// implementation error, and not a user or input error.
	if len(fields)%2 != 0 {
		panic(fmt.Sprintf("fields must be of even length: %#v", fields))
	}

	var m = make(map[string]interface{}, len(fields)/2)
	for i := 0; i != len(fields); i += 2 {
		var key = fields[i].(string)
		var value = fields[i+1]

		// Errors typically have JSON struct marshalling behavior and appear as '{}',
		// so explicitly cast them to their displayed string.
		if err, ok := value.(error); ok {
			value = err.Error()
		}

		m[key] = value
	}

	var fieldsRaw, err = json.Marshal(m)
	if err != nil {
		panic(err)
	}

	publisher.PublishLog(Log{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    json.RawMessage(fieldsRaw),
		Shard:     NewShardRef(publisher.Labels()),
		Spans:     nil, // Not supported from Go.
	})
}
