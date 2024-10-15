package models

type Stats struct {
	CPU     float64
	RAM     float64
	CPUTemp float64
}

func UnknownStats() Stats {
	return Stats{
		CPU:     -1,
		RAM:     -1,
		CPUTemp: -1,
	}
}
