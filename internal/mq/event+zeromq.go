package mq

import "github.com/guackamolly/zero-monitor/internal/logging"

func (s Socket) Publish(e Event) error {
	return s.PublishMsg(eventToMsg(e))
}

func (s Socket) Subscribe(t Topic) chan (EventOutput) {
	ch := make(chan (EventOutput))
	go func() {
		s.OnMsgReceived(t, func(m Msg) {
			if m.Topic != t {
				return
			}

			ch <- msgToEventOutput(m)
		})
	}()
	return ch
}

func eventToMsg(e Event) Msg {
	if te, ok := e.(DataEvent); ok {
		return compose(e.Topic(), te.Data).WithMetadata(e)
	}

	return compose(e.Topic()).WithMetadata(e)
}

func msgToEventOutput(m Msg) EventOutput {
	if e, ok := m.Metadata.(Event); ok {
		return NewBaseEventOutput(e, m.Data)
	}

	// TODO: handle cast fails!
	logging.LogWarning("received message that couldn't be casted to EventOutput")
	return nil
}
