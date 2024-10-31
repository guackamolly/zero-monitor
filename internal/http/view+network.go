package http

import (
	"slices"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type SpeedtestView models.Speedtest

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

type NetworkNodeSpeedtestView struct {
	NodeView
	Speedtest SpeedtestView
	Err       error
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

func NewNetworkNodeSpeedtestView(
	node models.Node,
	speedtest models.Speedtest,
	err error,
) NetworkNodeSpeedtestView {
	return NetworkNodeSpeedtestView{
		NodeView:  NodeView(node),
		Speedtest: SpeedtestView(speedtest),
		Err:       err,
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
