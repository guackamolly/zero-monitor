package db

import "github.com/guackamolly/zero-monitor/internal/data/models"

type SpeedtestEntity struct {
	models.Speedtest
	NodeID string
}

type SpeedtestTable CrudTable[SpeedtestEntity, string]

func NewSpeedtestEntity(
	speedtest models.Speedtest,
	nodeid string,
) SpeedtestEntity {
	return SpeedtestEntity{
		Speedtest: speedtest,
		NodeID:    nodeid,
	}
}

func (e SpeedtestEntity) PK() string {
	return e.ID
}
