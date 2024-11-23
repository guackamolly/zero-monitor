package http

import (
	"slices"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/labstack/echo/v4"
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
	Connections           []models.Connection
	ExposedTCPConnections []models.Connection
	ExposedUDPConnections []models.Connection
	Err                   error
}

type NetworkNodePackagesView struct {
	NodeView
	Packages []models.Package
	Err      error
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
	ContextView
	Speedtests []SpeedtestView
	Chart      *SpeedtestHistoryChartView
	Err        error
	Limit      int
}

type NetworkJoinView struct {
	PublicKeyURL          string `json:"public_key_url"`
	ConnectionEndpointURL string `json:"connection_url"`
}

type NetworkConnectionEndpointView struct {
	Host string `json:"host"`
	Port int    `json:"port"`
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
	exposedtcp := []models.Connection{}
	exposedudp := []models.Connection{}
	for _, conn := range connections {
		if !conn.Exposed() {
			continue
		}

		if conn.TCP() {
			exposedtcp = append(exposedtcp, conn)
			continue
		}

		if conn.UDP() {
			exposedudp = append(exposedudp, conn)
			continue
		}
	}

	return NetworkNodeConnectionsView{
		NodeView:              NodeView(node),
		Connections:           connections,
		ExposedTCPConnections: exposedtcp,
		ExposedUDPConnections: exposedudp,
		Err:                   err,
	}
}

func NewNetworkNodePackagesView(
	node models.Node,
	packages []models.Package,
	err error,
) NetworkNodePackagesView {
	return NetworkNodePackagesView{
		NodeView: NodeView(node),
		Packages: packages,
		Err:      err,
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
	ctx echo.Context,
	node models.Node,
	speedtests []models.Speedtest,
	limit int,
	err error,
) NetworkNodeSpeedtestHistoryView {
	sts := make([]SpeedtestView, len(speedtests))
	for i := range speedtests {
		sts[i] = NewSpeedtestView(node.ID, speedtests[i])
	}

	ctxview := NewContextView(ctx)
	chartSpeedtests := EligibleSpeedtestsForChartView(speedtests)
	if len(chartSpeedtests) == 0 {
		return NetworkNodeSpeedtestHistoryView{
			NodeView:    NodeView(node),
			ContextView: ctxview,
			Speedtests:  sts,
			Err:         err,
			Limit:       limit,
		}
	}

	chart := NewSpeedtestHistoryChartView(chartSpeedtests, ctxview.Breakpoint)
	return NetworkNodeSpeedtestHistoryView{
		NodeView:    NodeView(node),
		ContextView: ctxview,
		Chart:       &chart,
		Speedtests:  sts,
		Limit:       limit,
		Err:         err,
	}
}

func NewNetworkJoinView(
	publicKeyUrl string,
	connectionEndpointUrl string,
) NetworkJoinView {
	return NetworkJoinView{
		PublicKeyURL:          publicKeyUrl,
		ConnectionEndpointURL: connectionEndpointUrl,
	}
}

func NewNetworkConnectionEndpointView(
	host string,
	port int,
) NetworkConnectionEndpointView {
	return NetworkConnectionEndpointView{
		Host: host,
		Port: port,
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

func (v NetworkNodeSpeedtestHistoryView) AverageDownloadSpeed() string {
	avg := 0.0
	for _, st := range v.Speedtests {
		avg += float64(st.DownloadSpeed)
	}

	return models.BitRate(avg / float64(len(v.Speedtests))).String()
}

func (v NetworkNodeSpeedtestHistoryView) AverageUploadSpeed() string {
	avg := 0.0
	for _, st := range v.Speedtests {
		avg += float64(st.UploadSpeed)
	}

	return models.BitRate(avg / float64(len(v.Speedtests))).String()
}

func (v NetworkNodeSpeedtestHistoryView) PeakDownloadSpeedtest() SpeedtestView {
	var peak SpeedtestView
	for _, st := range v.Speedtests {
		if st.DownloadSpeed > peak.DownloadSpeed {
			peak = st
		}
	}

	return peak
}
