package mq

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
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

func (m Msg) WithError(err error) Msg {
	m.Data = &OPError{Err: err.Error()}
	return m
}

func (m Msg) Encrypt() (Msg, error) {
	bs, err := models.Encode(m)
	if err != nil {
		return Msg{}, err
	}

	bs, nonce, err := EncryptCipher(m.Identity, bs)
	if err != nil {
		return Msg{}, err
	}

	return Msg{
		Identity: m.Identity,
		Topic:    m.Topic,
		Data:     bs,
		Metadata: nonce,
	}, nil
}

func (m Msg) Decrypt() (Msg, error) {
	bs, ok := m.Data.([]byte)
	if !ok {
		return Msg{}, fmt.Errorf("data is not a bitstream")
	}

	nonce, ok := m.Metadata.([]byte)
	if !ok {
		return Msg{}, fmt.Errorf("nonce is not a bitstream")
	}

	bs, err := DecryptCipher(m.Identity, bs, nonce)
	if err != nil {
		return Msg{}, err
	}

	return models.Decode[Msg](bs)
}

func Compose(t Topic, d ...any) Msg {
	m := Msg{Topic: t}
	if len(d) > 0 {
		m.Data = d[0]
	}

	return m
}
