package http

import "github.com/guackamolly/zero-monitor/internal/data/models"

type SpeedtestPhaseView models.SpeedtestPhase

type SpeedtestView struct {
	models.Speedtest
}

func NewSpeedtestView(
	speedtest models.Speedtest,
) SpeedtestView {
	return SpeedtestView{
		Speedtest: speedtest,
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

func (v SpeedtestPhaseView) IsLatencyPhase() bool {
	return v == SpeedtestPhaseView(models.SpeedtestLatency)
}

func (v SpeedtestPhaseView) IsDownloadPhase() bool {
	return v == SpeedtestPhaseView(models.SpeedtestDownload)
}

func (v SpeedtestPhaseView) IsUploadPhase() bool {
	return v == SpeedtestPhaseView(models.SpeedtestUpload)
}
