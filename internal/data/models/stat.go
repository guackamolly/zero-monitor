package models

type Stats struct {
	CPU     Percent
	RAM     Percent
	CPUTemp Celsius
}

func NewStats(
	cpu float64,
	ram float64,
	cputemp float64,
) Stats {
	return Stats{
		CPU:     Percent(cpu),
		RAM:     Percent(ram),
		CPUTemp: Celsius(cputemp),
	}
}

func UnknownStats() Stats {
	return NewStats(-1, -1, -1)
}
