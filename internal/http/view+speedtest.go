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
