package conn

import (
	"bytes"
	"encoding/gob"
	"fmt"

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

func encode(d msg) ([]byte, error) {
	var buf bytes.Buffer
	var b []byte

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)

	if err == nil {
		b = buf.Bytes()
	}

	return b, err
}

func decode(b []byte) (msg, error) {
	var buf bytes.Buffer
	var m msg

	n, err := buf.Write(b)

	if n != len(b) || err != nil {
		return m, fmt.Errorf("couldn't write all bytes to buffer")
	}

	enc := gob.NewDecoder(&buf)
	err = enc.Decode(&m)

	return m, err
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
