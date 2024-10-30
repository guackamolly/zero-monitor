package models_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func TestSetOfflineReturnsNodeValueWithOnlinePropertyAsFalse(t *testing.T) {
	input := models.Node{Online: true}
	output := input.SetOffline()
	if output.Online {
		t.FailNow()
	}
}
