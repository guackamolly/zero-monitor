package dbbolt_test

import (
	"testing"

	dbbolt "github.com/guackamolly/zero-monitor/internal/data/db/db-bolt"
)

func TestPathReturnsDefaultIfEnvKeyIsNotSet(t *testing.T) {
	def := "master.db"
	if dbbolt.Path() != def {
		t.Errorf("expected %s but got %s", def, dbbolt.Path())
	}
}
