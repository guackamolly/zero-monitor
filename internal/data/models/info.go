package models

type Info struct {
	CPUArch      string
	CPUCount     int
	TotalRAM     Memory
	Hostname     string
	OS           string
	Distribution string
}

func NewInfo(
	cpuarch string,
	cpucount int,
	totalram uint64,
	hostname string,
	os string,
	distribution string,
) Info {
	return Info{
		CPUArch:      cpuarch,
		CPUCount:     cpucount,
		TotalRAM:     Memory(totalram),
		Hostname:     hostname,
		OS:           os,
		Distribution: distribution,
	}
}
