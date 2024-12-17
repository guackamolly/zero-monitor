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

type DisconnectNodeEvent struct {
	Event
	NodeID string
}

type DisconnectNodeEventOutput struct {
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

func NewDisconnectNodeEvent(
	nodeID string,
) DisconnectNodeEvent {
	return DisconnectNodeEvent{
		Event:  NewBaseEvent("Disconnect-node-event"),
		NodeID: nodeID,
	}
}

func NewDisconnectNodeEventOutput(
	origin Event,
	err error,
) DisconnectNodeEventOutput {
	return DisconnectNodeEventOutput{
		EventOutput: NewBaseEventOutput(origin, err),
	}
}
