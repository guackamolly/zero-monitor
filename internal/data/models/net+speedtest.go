package models

import "time"

// Init > Latency > Download > Upload > Finish
const (
	SpeedtestInit SpeedtestPhase = iota + 1
	SpeedtestLatency
	SpeedtestDownload
	SpeedtestUpload
	SpeedtestFinish
)

type SpeedtestPhase byte

type Speedtest struct {
	ID             string
	TakenAt        time.Time
	ServerProvider string
	ServerLocation string
	ServerMediator string
	ServerDistance Distance
	DownloadSpeed  BitRate
	UploadSpeed    BitRate
	Latency        Duration
	Phase          SpeedtestPhase
}

func NewSpeedtest(
	id string,
	serverprovider string,
	serverlocation string,
	servermediator string,
	serverdistance float64,
) Speedtest {
	return Speedtest{
		ID:             id,
		TakenAt:        time.Now(),
		ServerProvider: serverprovider,
		ServerLocation: serverlocation,
		ServerMediator: servermediator,
		ServerDistance: Distance(serverdistance),
		DownloadSpeed:  BitRate(0),
		UploadSpeed:    BitRate(0),
		Latency:        Duration(0),
		Phase:          SpeedtestInit,
	}
}

func (t Speedtest) WithUpdatedLatency(
	latency int64,
) Speedtest {
	t.Latency = Duration(latency)
	return t
}

func (t Speedtest) WithUpdatedDownloadSpeed(
	speed float64,
) Speedtest {
	t.DownloadSpeed = BitRate(speed)
	return t
}

func (t Speedtest) WithUpdatedUploadSpeed(
	speed float64,
) Speedtest {
	t.UploadSpeed = BitRate(speed)
	return t
}

func (t Speedtest) NextPhase() Speedtest {
	if t.Phase == SpeedtestFinish {
		return t
	}

	t.Phase = t.Phase + 1
	return t
}

func (t Speedtest) Finished() bool {
	return t.Phase == SpeedtestFinish
}

func (t Speedtest) FinishedLatency() bool {
	return t.Phase > SpeedtestLatency
}

func (t Speedtest) FinishedDownload() bool {
	return t.Phase > SpeedtestDownload
}

func (t Speedtest) FinishedUpload() bool {
	return t.Phase > SpeedtestUpload
}
