package event

import "github.com/guackamolly/zero-monitor/internal/data/models"

type KillNodeProcessEvent struct {
	Event
	NodeID string
	PID    int32
}

type KillNodeProcessEventOutput struct {
	EventOutput
	Processes []models.Process
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
	processes []models.Process,
	err error,
) KillNodeProcessEventOutput {
	return KillNodeProcessEventOutput{
		EventOutput: NewBaseEventOutput(origin, err),
		Processes:   processes,
	}
}
