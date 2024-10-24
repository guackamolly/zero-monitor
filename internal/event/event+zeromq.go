package event

import (
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/mq"
)

type ZeroMQEventPublisher struct {
	mq.Socket
}

func (s ZeroMQEventPublisher) Publish(e Event) error {
	return s.PublishMsg(eventToMsg(e))
}

func (s ZeroMQEventPublisher) Subscribe(t mq.Topic) chan (EventOutput) {
	ch := make(chan (EventOutput))
	go func() {
		s.OnMsgReceived(t, func(m mq.Msg) {
			if m.Topic != t {
				return
			}

			ch <- msgToEventOutput(m)
		})
	}()
	return ch
}

func eventToMsg(e Event) mq.Msg {
	if te, ok := e.(DataEvent); ok {
		return mq.Compose(e.Topic(), te.Data).WithMetadata(e)
	}

	return mq.Compose(e.Topic()).WithMetadata(e)
}

func msgToEventOutput(m mq.Msg) EventOutput {
	if e, ok := m.Metadata.(Event); ok {
		return NewBaseEventOutput(e, m.Data)
	}

	// TODO: handle cast fails!
	logging.LogWarning("received message that couldn't be casted to EventOutput")
	return nil
}
