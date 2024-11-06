package repositories_test

import (
	"slices"
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/db"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
)

type TestSpeedtestTable struct {
	db.SpeedtestTable
	Speedtests map[string][]models.Speedtest
}

func NewTestSpeedtestTable(
	speedtests map[string][]models.Speedtest,
) *TestSpeedtestTable {
	return &TestSpeedtestTable{
		Speedtests: speedtests,
	}
}

func (t TestSpeedtestTable) All() ([]db.SpeedtestEntity, error) {
	entities := []db.SpeedtestEntity{}
	for nid, sts := range t.Speedtests {
		for _, st := range sts {
			entities = append(entities, db.NewSpeedtestEntity(st, nid))
		}
	}

	return entities, nil
}

func TestHistoryReturnsOnlyTheSpeedtestsDoneOnSpecificNode(t *testing.T) {
	nodeid := "node.id.1"
	nodeid2 := "node.id.2"
	nodeidSpeedtests := []models.Speedtest{
		models.NewSpeedtest("id.1", "-", "-", "-", 0),
		models.NewSpeedtest("id.2", "-", "-", "-", 0),
		models.NewSpeedtest("id.3", "-", "-", "-", 0),
	}
	nodeid2Speedtests := []models.Speedtest{
		models.NewSpeedtest("id.4", "-", "-", "-", 0),
	}

	tbl := NewTestSpeedtestTable(map[string][]models.Speedtest{
		nodeid:  nodeidSpeedtests,
		nodeid2: nodeid2Speedtests,
	})
	repo := repositories.NewDatabaseSpeedtestStoreRepository(tbl)

	hs, err := repo.History(nodeid)
	if err != nil {
		t.Fatalf("was not expecting History() to error, but got %v", err)
	}

	if !slices.Equal(hs, nodeidSpeedtests) {
		t.Errorf("expected %v to equal to %v", hs, nodeidSpeedtests)
	}
}
