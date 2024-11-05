package dbbolt

import (
	"encoding/gob"

	"github.com/guackamolly/zero-monitor/internal/data/db"
)

// Register here all entities that will be persisted in a bolt database.
func init() {
	gob.Register(db.SpeedtestEntity{})
}
