package http

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type SpeedtestPhaseView models.SpeedtestPhase
type SpeedtestView struct {
	models.Speedtest
	NodeID string
}

type SpeedtestHistoryChartView struct {
	ChartView
}

func NewSpeedtestView(
	nodeid string,
	speedtest models.Speedtest,
) SpeedtestView {
	return SpeedtestView{
		Speedtest: speedtest,
		NodeID:    nodeid,
	}
}

func NewSpeedtestStatusElementView(
	status string,
) ElementView {
	return NewElementView("speedtest-status", status)
}

func NewSpeedtestLatencyElementView(
	latency models.Duration,
) ElementView {
	return NewElementView("speedtest-latency", latency.String())
}

func NewSpeedtestDownloadElementView(
	download models.BitRate,
) ElementView {
	return NewElementView("speedtest-download", download.String())
}

func NewSpeedtestUploadElementView(
	upload models.BitRate,
) ElementView {
	return NewElementView("speedtest-upload", upload.String())
}

func NewSpeedtestHistoryChartView(
	speedtests []models.Speedtest,
) SpeedtestHistoryChartView {
	xaxis := NewAxisView(
		float64(speedtests[len(speedtests)-1].TakenAt.UnixMilli())-float64(time.Minute.Milliseconds()),
		float64(speedtests[0].TakenAt.UnixMilli())+float64(time.Minute.Milliseconds()),
		"",
	)

	yaxis := NewAxisView(
		50*models.Mbit,
		120*models.Mbit,
		"Bitrate (Mbps)",
	)

	xvalues := make([]float64, len(speedtests))
	y1values := make([]float64, len(speedtests))
	y2values := make([]float64, len(speedtests))
	for i, st := range speedtests {
		t := st.TakenAt.UnixMilli()

		xvalues[i] = float64(t)
		y1values[i] = float64(st.DownloadSpeed)
		y2values[i] = float64(st.UploadSpeed)
	}

	lines := []LineView{
		NewLineView("Download", xvalues, y1values, TimeFormatter, BitrateFormatter, func(i int) string {
			return speedtests[i].DownloadSpeed.String()
		}),
		NewLineView("Upload", xvalues, y2values, TimeFormatter, BitrateFormatter, func(i int) string {
			return speedtests[i].UploadSpeed.String()
		}),
	}

	return SpeedtestHistoryChartView{
		NewLineChartView(lines, xaxis, yaxis),
	}
}

func (v SpeedtestView) Status() SpeedtestPhaseView {
	return SpeedtestPhaseView(v.Phase)
}

func (v SpeedtestView) TakenAt() string {
	return v.Speedtest.TakenAt.Format(time.DateTime)
}

func (v SpeedtestPhaseView) IsLatencyPhase() bool {
	return v == SpeedtestPhaseView(models.SpeedtestLatency)
}

func (v SpeedtestPhaseView) IsDownloadPhase() bool {
	return v == SpeedtestPhaseView(models.SpeedtestDownload)
}

func (v SpeedtestPhaseView) IsUploadPhase() bool {
	return v == SpeedtestPhaseView(models.SpeedtestUpload)
}
