package mq

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type msg struct {
	topic Topic
	data  any
}

func encode(m msg) ([]byte, error) {
	var buf bytes.Buffer
	var b []byte

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)

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

func compose(t Topic, d ...any) msg {
	m := msg{topic: t}
	if len(d) > 0 {
		m.data = d[0]
	}

	return m
}
