package models

type Info struct {
	CPUArch      string
	CPUModel     string
	CPUCache     Memory
	CPUCount     int
	TotalRAM     Memory
	TotalDisk    Memory
	Hostname     string
	OS           string
	Distribution string
	Kernel       string
}

func NewInfo(
	cpuarch string,
	cpumodel string,
	cpucache int32,
	cpucount int,
	totalram uint64,
	totaldisk uint64,
	hostname string,
	os string,
	distribution string,
	kernel string,
) Info {
	return Info{
		CPUArch:      cpuarch,
		CPUModel:     cpumodel,
		CPUCache:     Memory(cpucache),
		CPUCount:     cpucount,
		TotalRAM:     Memory(totalram),
		TotalDisk:    Memory(totaldisk),
		Hostname:     hostname,
		OS:           os,
		Distribution: distribution,
		Kernel:       kernel,
	}
}
