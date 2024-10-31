package event

import "github.com/guackamolly/zero-monitor/internal/data/models"

type KillNodeProcessEvent struct {
	Event
	NodeID string
	PID    int32
}

type KillNodeProcessEventOutput struct {
	EventOutput
}

type StartNodeSpeedtestEvent struct {
	Event
	NodeID string
}

type NodeSpeedtestEventOutput struct {
	EventOutput
	Speedtest models.Speedtest
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

func NewStartNodeSpeedtestEvent(
	nodeID string,
) StartNodeSpeedtestEvent {
	return StartNodeSpeedtestEvent{
		Event:  NewBaseEvent("kill-node-process-event"),
		NodeID: nodeID,
	}
}

func NewNodeSpeedtestEventOutput(
	origin Event,
	speedtest models.Speedtest,
	err error,
) NodeSpeedtestEventOutput {
	return NodeSpeedtestEventOutput{
		EventOutput: NewBaseEventOutput(origin, err),
		Speedtest:   speedtest,
	}
}
