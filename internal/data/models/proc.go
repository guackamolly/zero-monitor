package models

type Process struct {
	PID    int32
	User   string
	Name   string
	CMD    string
	Memory Memory
	CPU    Percent
}

func NewProcess(
	pid int32,
	user string,
	name string,
	cmd string,
) Process {
	return Process{
		PID:  pid,
		User: user,
		Name: name,
		CMD:  cmd,
	}
}

func (p Process) WithUpdatedMemory(
	memory uint64,
) Process {
	p.Memory = Memory(memory)
	return p
}

func (p Process) WithUpdatedCPU(
	cpu float64,
) Process {
	p.CPU = Percent(cpu)
	return p
}
