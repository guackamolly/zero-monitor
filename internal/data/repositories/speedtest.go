package repositories

import "github.com/guackamolly/zero-monitor/internal/data/models"

type SpeedtestRepository interface {
	Start() (chan (models.Speedtest), error)
}

type SpeedtestStoreRepository interface {
	Save(nodeid string, speedtest models.Speedtest) error
	Lookup(id string) (models.Speedtest, bool, error)
	History(nodeid string) ([]models.Speedtest, error)
}
