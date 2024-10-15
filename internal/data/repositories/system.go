package repositories

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemRepository interface {
	Info() (models.Info, error)
	Stats() (models.Stats, error)
}

type GopsUtilSystemRepository struct{}

func (r GopsUtilSystemRepository) Info() (models.Info, error) {
	cc, err := cpu.Counts(false)
	if err != nil {
		return models.Info{}, err
	}

	rs, err := mem.VirtualMemory()
	if err != nil {
		return models.Info{}, err
	}

	os, err := host.Info()
	if err != nil {
		return models.Info{}, err
	}

	return models.Info{
		CPUCount:     cc,
		CPUArch:      os.KernelArch,
		TotalRAM:     rs.Total,
		OS:           os.OS,
		Distribution: os.Platform,
	}, nil
}

func (r GopsUtilSystemRepository) Stats() (models.Stats, error) {
	cp, err := cpu.Percent(time.Millisecond*150, false)
	if err != nil {
		return models.Stats{}, err
	}

	rs, err := mem.VirtualMemory()
	if err != nil {
		return models.Stats{}, err
	}

	st, _ := host.SensorsTemperatures()

	temp := float64(0)
	for _, t := range st {
		temp += t.Temperature
	}

	if l := len(st); l > 0 {
		temp = temp / float64(l)
	} else {
		temp = -1
	}

	return models.Stats{
		CPU:     cp[0],
		RAM:     rs.UsedPercent,
		CPUTemp: temp,
	}, nil
}
