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
	Breakpoint
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
	breakpoint Breakpoint,
) SpeedtestHistoryChartView {
	xaxis := NewAxisView(
		float64(speedtests[len(speedtests)-1].TakenAt.UnixMilli())-float64(time.Minute.Milliseconds()),
		float64(speedtests[0].TakenAt.UnixMilli())+float64(time.Minute.Milliseconds()),
		"",
	)

	min := speedtests[0].DownloadSpeed
	max := speedtests[0].UploadSpeed

	for _, st := range speedtests {
		if st.DownloadSpeed < min {
			min = st.DownloadSpeed
		}

		if st.UploadSpeed < min {
			min = st.UploadSpeed
		}

		if st.DownloadSpeed > max {
			max = st.DownloadSpeed
		}

		if st.UploadSpeed > max {
			max = st.UploadSpeed
		}
	}

	min -= min * 0.15
	max += max * 0.15

	yaxis := NewAxisView(
		float64(min),
		float64(max),
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
		breakpoint,
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

func (v SpeedtestHistoryChartView) SVG() string {
	return v.ChartView.SVG(v.Breakpoint.ChartSize())
}

func EligibleSpeedtestsForChartView(speedtests []models.Speedtest) []models.Speedtest {
	if len(speedtests) == 0 {
		return speedtests
	}

	sts := []models.Speedtest{
		speedtests[0],
	}
	for i := 1; i < len(speedtests); i++ {
		if sts[i-1].TakenAt.Sub(speedtests[i].TakenAt) > 20*time.Minute {
			break
		}

		sts = append(sts, speedtests[i])
	}

	if len(sts) < 3 {
		return []models.Speedtest{}
	}

	return sts
}
