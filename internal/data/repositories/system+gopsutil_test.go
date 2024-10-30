//go:build integration
// +build integration

package repositories_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/repositories"
)

func TestGopsUtilSystemRepositoryExtractsRxTxEverySecond(t *testing.T) {
	repo := repositories.NewGopsUtilSystemRepository()
	before := repo.TotalRx + repo.TotalTx

	// force Rx/Tx change
	_, err := http.Get("https://github.com/guackamolly/zero-monitor")
	if err != nil {
		t.Logf("failed to GET before sleeping, Rx/Tx might not change..., %v", err)
	}
	time.Sleep(time.Second)
	after := repo.TotalRx + repo.TotalTx

	if before == after {
		t.FailNow()
	}
}
