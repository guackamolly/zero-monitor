package models

type Info struct {
	CPUArch      string
	CPUCount     int
	TotalRAM     Memory
	OS           string
	Distribution string
}

func NewInfo(
	cpuarch string,
	cpucount int,
	totalram uint64,
	os string,
	distribution string,
) Info {
	return Info{
		CPUArch:      cpuarch,
		CPUCount:     cpucount,
		TotalRAM:     Memory(totalram),
		OS:           os,
		Distribution: distribution,
	}
}
