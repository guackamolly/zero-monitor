package mq

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Msg interface {
	Identity() []byte
	Topic() Topic
	Data() any
}

type BaseMsg struct {
	BIdentity []byte
	BTopic    Topic
	BData     any
}

func (m BaseMsg) Identity() []byte {
	return m.BIdentity
}

func (m BaseMsg) Topic() Topic {
	return m.BTopic
}

func (m BaseMsg) Data() any {
	return m.BData
}

func encode(m Msg) ([]byte, error) {
	var buf bytes.Buffer
	var b []byte

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)

	if err == nil {
		b = buf.Bytes()
	}

	return b, err
}

func decode(b []byte) (Msg, error) {
	var buf bytes.Buffer
	var m Msg

	n, err := buf.Write(b)

	if n != len(b) || err != nil {
		return m, fmt.Errorf("couldn't write all bytes to buffer")
	}

	enc := gob.NewDecoder(&buf)
	err = enc.Decode(&m)

	return m, err
}

func compose(t Topic, d ...any) Msg {
	m := BaseMsg{BTopic: t}
	if len(d) > 0 {
		m.BData = d[0]
	}

	return m
}
