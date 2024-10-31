package models

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
	Provider      string
	Server        string
	DownloadSpeed IORate
	UploadSpeed   IORate
	Latency       Duration
	Phase         SpeedtestPhase
}

func NewSpeedtest(
	provider string,
	server string,
) Speedtest {
	return Speedtest{
		Provider:      provider,
		Server:        server,
		DownloadSpeed: IORate(0),
		UploadSpeed:   IORate(0),
		Latency:       Duration(0),
		Phase:         SpeedtestInit,
	}
}

func (t Speedtest) WithUpdatedLatency(
	latency int64,
) Speedtest {
	t.Latency = Duration(latency)
	return t
}

func (t Speedtest) WithUpdatedDownloadSpeed(
	speed uint64,
) Speedtest {
	t.DownloadSpeed = IORate(speed)
	return t
}

func (t Speedtest) WithUpdatedUploadSpeed(
	speed uint64,
) Speedtest {
	t.UploadSpeed = IORate(speed)
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
