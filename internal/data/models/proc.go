package models

type Process struct {
	PID    int32
	User   string
	Name   string
	CMD    string
	Memory Memory
}

func NewProcess(
	pid int32,
	user string,
	name string,
	cmd string,
	memory uint64,
) Process {
	return Process{
		PID:    pid,
		User:   user,
		Name:   name,
		CMD:    cmd,
		Memory: Memory(memory),
	}
}
