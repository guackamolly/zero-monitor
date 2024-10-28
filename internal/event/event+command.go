package event

type KillNodeProcessEvent struct {
	Event
	NodeID string
	PID    int32
}

type KillNodeProcessEventOutput struct {
	EventOutput
}

func NewKillNodeProcessEvent(
	nodeID string,
	pid int32,
) KillNodeProcessEvent {
	return KillNodeProcessEvent{
		Event:  NewBaseEvent("kill-node-process-event"),
		NodeID: nodeID,
		PID:    pid,
	}
}

func NewKillNodeProcessEventOutput(
	origin Event,
	err error,
) KillNodeProcessEventOutput {
	return KillNodeProcessEventOutput{
		EventOutput: NewBaseEventOutput(origin, err),
	}
}
