package http

import (
	"slices"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type NetworkView struct {
	Online             []NodeView
	Offline            []NodeView
	NodeStatsPollBurst time.Duration
}

type NetworkNodeInformationView struct {
	NodeView
}

type NetworkNodeConnectionsView struct {
	NodeView
	Connections []models.Connection
	Err         error
}

type NetworkNodeProcessesView struct {
	NodeView
	Processes []models.Process
	Err       error
}

type StartNetworkNodeSpeedtestView struct {
	NodeView
	Err error
}

type NetworkNodeSpeedtestView struct {
	NodeView
	Speedtest SpeedtestView
	Err       error
}

type NetworkNodeSpeedtestHistoryView struct {
	NodeView
	Speedtests []SpeedtestView
	Err        error
}

func NewNetworkView(
	nodes []models.Node,
	nodeStatsPollBurst time.Duration,
) NetworkView {
	on := []NodeView{}
	off := []NodeView{}
	for _, v := range nodes {
		if v.Online {
			on = append(on, NodeView(v))
		} else {
			off = append(off, NodeView(v))
		}
	}

	return NetworkView{
		Online:             on,
		Offline:            off,
		NodeStatsPollBurst: nodeStatsPollBurst,
	}
}

func NewNetworkNodeInformationView(
	node models.Node,
) NetworkNodeInformationView {
	return NetworkNodeInformationView{
		NodeView: NodeView(node),
	}
}

func NewNetworkNodeConnectionsView(
	node models.Node,
	connections []models.Connection,
	err error,
) NetworkNodeConnectionsView {
	return NetworkNodeConnectionsView{
		NodeView:    NodeView(node),
		Connections: connections,
		Err:         err,
	}
}

func NewNetworkNodeProcessesView(
	node models.Node,
	processes []models.Process,
	err error,
) NetworkNodeProcessesView {
	slices.Reverse(processes)

	return NetworkNodeProcessesView{
		NodeView:  NodeView(node),
		Processes: processes,
		Err:       err,
	}
}

func NewStartNetworkNodeSpeedtestView(
	node models.Node,
	err error,
) StartNetworkNodeSpeedtestView {
	return StartNetworkNodeSpeedtestView{
		NodeView: NodeView(node),
		Err:      err,
	}
}

func NewNetworkNodeSpeedtestView(
	node models.Node,
	speedtest models.Speedtest,
	err error,
) NetworkNodeSpeedtestView {
	return NetworkNodeSpeedtestView{
		NodeView:  NodeView(node),
		Speedtest: NewSpeedtestView(node.ID, speedtest),
		Err:       err,
	}
}

func NewNetworkNodeSpeedtestHistoryView(
	node models.Node,
	speedtests []models.Speedtest,
	err error,
) NetworkNodeSpeedtestHistoryView {
	sts := make([]SpeedtestView, len(speedtests))
	for i := len(speedtests) - 1; i >= 0; i-- {
		sts[i] = NewSpeedtestView(node.ID, speedtests[i])
	}

	// reverse to show most recent speedtests first
	slices.Reverse(sts)

	return NetworkNodeSpeedtestHistoryView{
		NodeView:   NodeView(node),
		Speedtests: sts,
		Err:        err,
	}
}

func (v NetworkNodeProcessesView) CPU() string {
	cpu := models.Percent(0)
	for _, p := range v.Processes {
		cpu += p.CPU
	}

	cpu = cpu / models.Percent(v.Info.CPU.Count)

	return cpu.String()
}

func (v NetworkNodeProcessesView) Memory() string {
	mem := models.Memory(0)
	for _, p := range v.Processes {
		mem += p.Memory
	}

	return mem.String()
}
