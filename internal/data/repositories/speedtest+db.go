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
	// TODO: request only those that match the foreign key (nodeid)
	entities, err := r.tbl.All()
	if err != nil {
		return nil, err
	}

	sts := []models.Speedtest{}
	for _, entity := range entities {
		if entity.NodeID == nodeid {
			sts = append(sts, entity.Speedtest)
		}
	}

	return sts, nil
}
