package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func Encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	var b []byte

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)

	if err == nil {
		b = buf.Bytes()
	}

	return b, err
}

func Decode[T any](b []byte) (T, error) {
	var buf bytes.Buffer
	var m T

	n, err := buf.Write(b)

	if n != len(b) || err != nil {
		return m, fmt.Errorf("couldn't write all bytes to buffer")
	}

	enc := gob.NewDecoder(&buf)
	err = enc.Decode(&m)

	return m, err
}
