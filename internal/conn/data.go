package conn

import (
	"encoding/gob"

	"github.com/guackamolly/zero-monitor/internal"
)

type msg struct {
	Key  byte
	Id   string
	Data any
}

func init() {
	gob.Register(Connection{})
}

func compose(key byte, data ...any) msg {
	switch len(data) {
	case 0:
		return msg{Id: internal.MachineId, Key: key}
	case 1:
		return msg{Id: internal.MachineId, Key: key, Data: data[0]}
	default:
		return msg{Id: internal.MachineId, Key: key, Data: data}
	}
}
