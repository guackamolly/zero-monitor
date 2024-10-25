package event

import (
	"encoding/gob"
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/mq"
)

type ZeroMQEventPubSub struct {
	*mq.Socket
}

func NewZeroMQEventPubSub(
	socket *mq.Socket,
) ZeroMQEventPubSub {
	return ZeroMQEventPubSub{
		Socket: socket,
	}
}

func (p ZeroMQEventPubSub) Publish(e Event) error {
	msg, err := p.eventToMsg(e)
	if err != nil {
		return err
	}
	p.ReplyMsg(msg.Identity, msg)
	return nil
}

func (p ZeroMQEventPubSub) Subscribe(e Event) chan (EventOutput) {
	msg, err := p.eventToMsg(e)
	if err != nil {
		return nil
	}
	t := msg.Topic

	ch := make(chan (EventOutput))
	p.OnMsgReceived(t, func(m mq.Msg) {
		if m.Topic != t {
			return
		}

		go func() {
			o, _ := p.msgToEventOutput(m)
			ch <- o
		}()
	})
	return ch
}

func (p ZeroMQEventPubSub) eventToMsg(e Event) (mq.Msg, error) {
	switch te := e.(type) {
	case QueryNodeConnectionsEvent:
		sid, ok := p.Clients[te.NodeID]
		if !ok {
			return mq.Msg{}, fmt.Errorf("no pub client associated with id, %v", te.NodeID)
		}

		return mq.Compose(mq.NodeConnections).WithIdentity(sid).WithMetadata(e), nil
	default:
		return mq.Msg{}, fmt.Errorf("couldn't match event with a topic, %v", e)
	}
}

func (p ZeroMQEventPubSub) msgToEventOutput(m mq.Msg) (EventOutput, error) {
	switch m.Topic {
	case mq.NodeConnections:
		if resp, ok := m.Data.(mq.NodeConnectionsResponse); ok {
			return NewQueryNodeConnectionsEventOutput(
				m.Metadata.(Event),
				resp.Connections,
				nil,
			), nil
		}

		return NewQueryNodeConnectionsEventOutput(
			m.Metadata.(Event),
			nil,
			m.Data.(error),
		), nil
	default:
		return nil, fmt.Errorf("couldn't match message with a topic, %v", m)
	}
}

func init() {
	gob.Register(BaseEvent{})
	gob.Register(QueryNodeConnectionsEvent{})
}
