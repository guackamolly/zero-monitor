package mq

import "github.com/guackamolly/zero-monitor/internal/logging"

type EventMsg struct {
	Event
	Msg
}

func (em EventMsg) Topic() Topic {
	return em.Event.Topic()
}

func NewEventMsg(
	event Event,
	msg Msg,
) Msg {
	return EventMsg{
		Msg:   msg,
		Event: event,
	}
}

func (s Socket) Publish(e Event) error {
	return s.PublishMsg(eventToMsg(e))
}

func (s Socket) Subscribe(t Topic) chan (EventOutput) {
	ch := make(chan (EventOutput))
	go func() {
		s.OnMsgReceived(t, func(m Msg) {
			if m.Topic() != t {
				return
			}

			ch <- msgToEventOutput(m)
		})
	}()
	return ch
}

func eventToMsg(e Event) Msg {
	if te, ok := e.(DataEvent); ok {
		return NewEventMsg(
			e,
			compose(e.Topic(), te.Data),
		)
	}

	return NewEventMsg(
		e,
		compose(e.Topic()),
	)
}

func msgToEventOutput(m Msg) EventOutput {
	if tm, ok := m.(EventMsg); ok {
		return NewBaseEventOutput(tm.Event, tm.Data())
	}

	// TODO: handle cast fails!
	logging.LogWarning("received message that couldn't be casted to EventOutput")
	return nil
}
