package materialize

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/estuary/flow/go/flow/ops"
	"github.com/estuary/flow/go/materialize/driver/image"
	pf "github.com/estuary/flow/go/protocols/flow"
	pm "github.com/estuary/flow/go/protocols/materialize"
)

// NewDriver returns a new driver implementation for the given EndpointType.
func NewDriver(
	ctx context.Context,
	endpointType pf.EndpointType,
	endpointSpec json.RawMessage,
	connectorNetwork string,
	logPublisher ops.Logger,
) (pm.DriverClient, error) {

	switch endpointType {
	case pf.EndpointType_FLOW_SINK:
		return pm.AdaptServerToClient(image.NewDriver(connectorNetwork, logPublisher)), nil
	default:
		return nil, fmt.Errorf("unknown endpoint %v", endpointType)
	}
}
