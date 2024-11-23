package repositories

import (
	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type SystemRepository interface {
	Info() (models.MachineInfo, error)
	Stats() (models.Stats, error)
	Conns() ([]models.Connection, error)
	Procs() ([]models.Process, error)
	Pkgs() ([]models.Package, error)
	KillProc(pid int32) error
}
