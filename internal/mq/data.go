package mq

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

func Compose(t Topic, d ...any) Msg {
	m := Msg{Topic: t}
	if len(d) > 0 {
		m.Data = d[0]
	}

	return m
}
