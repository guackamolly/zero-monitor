package models

type Process struct {
	PID  int32
	Name string
}

func NewProcess(
	pid int32,
	name string,
) Process {
	return Process{
		PID:  pid,
		Name: name,
	}
}
