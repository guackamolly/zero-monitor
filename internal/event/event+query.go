package event

type QueryNodeConnectionsEvent struct {
	Event
	NodeID string
}

func NewQueryNodeConnectionsEvent(
	nodeID string,
) QueryNodeConnectionsEvent {
	return QueryNodeConnectionsEvent{
		Event:  NewBaseEvent("query-node-connections-event"),
		NodeID: nodeID,
	}
}
