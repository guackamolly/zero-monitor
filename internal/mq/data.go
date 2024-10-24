package mq

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Msg struct {
	Identity []byte
	Topic    Topic
	Data     any
	Metadata any
}

func (m Msg) WithIdentity(identity []byte) Msg {
	m.Identity = identity
	return m
}

func (m Msg) WithMetadata(metadata any) Msg {
	m.Metadata = metadata
	return m
}

func (m Msg) WithData(data any) Msg {
	m.Data = data
	return m
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
		return Msg{}, fmt.Errorf("couldn't write all bytes to buffer")
	}

	enc := gob.NewDecoder(&buf)
	err = enc.Decode(&m)

	return m, err
}

func Compose(t Topic, d ...any) Msg {
	m := Msg{Topic: t}
	if len(d) > 0 {
		m.Data = d[0]
	}

	return m
}
