package event

import "github.com/guackamolly/zero-monitor/internal/data/models"

type QueryNodeConnectionsEvent struct {
	Event
	NodeID string
}

type QueryNodeConnectionsEventOutput struct {
	EventOutput
	Connections []models.Connection
}

func NewQueryNodeConnectionsEvent(
	nodeID string,
) QueryNodeConnectionsEvent {
	return QueryNodeConnectionsEvent{
		Event:  NewBaseEvent("query-node-connections-event"),
		NodeID: nodeID,
	}
}

func NewQueryNodeConnectionsEventOutput(
	origin Event,
	connections []models.Connection,
	err error,
) QueryNodeConnectionsEventOutput {
	return QueryNodeConnectionsEventOutput{
		EventOutput: NewBaseEventOutput(origin, err),
		Connections: connections,
	}
}
