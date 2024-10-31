package repositories

import "github.com/guackamolly/zero-monitor/internal/data/models"

type SpeedtestRepository interface {
	Start() (chan (models.Speedtest), error)
}
