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

type QueryNodeProcessesEvent struct {
	Event
	NodeID string
}

type QueryNodePackagesEvent struct {
	Event
	NodeID string
}

type QueryNodeProcessesEventOutput struct {
	EventOutput
	Processes []models.Process
}

type QueryNodePackagesEventOutput struct {
	EventOutput
	Packages []models.Package
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

func NewQueryNodeProcessesEvent(
	nodeID string,
) QueryNodeProcessesEvent {
	return QueryNodeProcessesEvent{
		Event:  NewBaseEvent("query-node-processes-event"),
		NodeID: nodeID,
	}
}

func NewQueryNodeProcessesEventOutput(
	origin Event,
	processes []models.Process,
	err error,
) QueryNodeProcessesEventOutput {
	return QueryNodeProcessesEventOutput{
		EventOutput: NewBaseEventOutput(origin, err),
		Processes:   processes,
	}
}

func NewQueryNodePackagesEvent(
	nodeID string,
) QueryNodePackagesEvent {
	return QueryNodePackagesEvent{
		Event:  NewBaseEvent("query-node-packages-event"),
		NodeID: nodeID,
	}
}

func NewQueryNodePackagesEventOutput(
	origin Event,
	packages []models.Package,
	err error,
) QueryNodePackagesEventOutput {
	return QueryNodePackagesEventOutput{
		EventOutput: NewBaseEventOutput(origin, err),
		Packages:    packages,
	}
}
