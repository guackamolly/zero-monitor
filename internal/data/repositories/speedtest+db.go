package repositories

import (
	"github.com/guackamolly/zero-monitor/internal/data/db"
	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type DatabaseSpeedtestStoreRepository struct {
	tbl db.SpeedtestTable
}

func NewDatabaseSpeedtestStoreRepository(
	tbl db.SpeedtestTable,
) *DatabaseSpeedtestStoreRepository {
	return &DatabaseSpeedtestStoreRepository{
		tbl: tbl,
	}
}

func (r DatabaseSpeedtestStoreRepository) Save(nodeid string, speedtest models.Speedtest) error {
	entity := db.NewSpeedtestEntity(speedtest, nodeid)
	return r.tbl.Insert(entity)
}

func (r DatabaseSpeedtestStoreRepository) History(nodeid string) ([]models.Speedtest, error) {
	entities, err := r.tbl.All()
	if err != nil {
		return nil, err
	}

	sts := make([]models.Speedtest, len(entities))
	for i, entity := range entities {
		sts[i] = entity.Speedtest
	}

	return sts, nil
}
