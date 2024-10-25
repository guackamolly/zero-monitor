package models

import "time"

type Stats struct {
	CPU     Percent
	RAM     Percent
	Disk    Percent
	CPUTemp Celsius
	Uptime  Duration
	Rx      IORate
	Tx      IORate
	TotalRx Memory
	TotalTx Memory
}

func NewStats(
	cpu float64,
	ram float64,
	disk float64,
	cputemp float64,
	uptime uint64,
	rx uint64,
	tx uint64,
	totalRx uint64,
	totalTx uint64,
) Stats {
	return Stats{
		CPU:     Percent(cpu),
		RAM:     Percent(ram),
		Disk:    Percent(disk),
		CPUTemp: Celsius(cputemp),
		Uptime:  Duration(time.Duration(uptime) * time.Second),
		Rx:      IORate(rx),
		Tx:      IORate(tx),
		TotalRx: Memory(totalRx),
		TotalTx: Memory(totalTx),
	}
}

func UnknownStats() Stats {
	return NewStats(-1, -1, -1, -1, 0, 0, 0, 0, 0)
}
