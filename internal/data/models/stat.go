package models

import "time"

type Stats struct {
	CPU     Percent
	RAM     Percent
	CPUTemp Celsius
	Uptime  Duration
}

func NewStats(
	cpu float64,
	ram float64,
	cputemp float64,
	uptime uint64,
) Stats {
	return Stats{
		CPU:     Percent(cpu),
		RAM:     Percent(ram),
		CPUTemp: Celsius(cputemp),
		Uptime:  Duration(time.Duration(uptime) * time.Second),
	}
}

func UnknownStats() Stats {
	return NewStats(-1, -1, -1, 0)
}
