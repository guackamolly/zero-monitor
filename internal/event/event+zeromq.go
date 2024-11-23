package event

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
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

func (p ZeroMQEventPubSub) Subscribe(e Event) (chan (EventOutput), CloseSubscription) {
	msg, err := p.eventToMsg(e)
	if err != nil {
		return nil, nil
	}
	t := msg.Topic

	ch := make(chan (EventOutput))
	close := p.OnMsgReceived(t, func(m mq.Msg) {
		if m.Topic != t {
			return
		}

		go func() {
			o, _ := p.msgToEventOutput(m)
			ch <- o
		}()
	})

	return ch, close
}

func (p ZeroMQEventPubSub) PublicKey() ([]byte, error) {
	return mq.DerivePublicKey()
}

func (p ZeroMQEventPubSub) Address() models.Address {
	addr, err := models.NewNetAddress(p.Addr())
	if err != nil {
		// TODO: there's no way this errors since zeromq uses tcp sockets, maybe think of a better way to create the address
		panic("zeromq is not using a tcp socket")
	}

	return addr
}

func (p ZeroMQEventPubSub) eventToMsg(e Event) (mq.Msg, error) {
	switch te := e.(type) {
	case QueryNodeConnectionsEvent:
		sid, ok := p.Clients[te.NodeID]
		if !ok {
			return mq.Msg{}, fmt.Errorf("no pub client associated with id, %v", te.NodeID)
		}

		return mq.Compose(mq.NodeConnections).WithIdentity(sid).WithMetadata(e), nil
	case QueryNodeProcessesEvent:
		sid, ok := p.Clients[te.NodeID]
		if !ok {
			return mq.Msg{}, fmt.Errorf("no pub client associated with id, %v", te.NodeID)
		}

		return mq.Compose(mq.NodeProcesses).WithIdentity(sid).WithMetadata(e), nil
	case QueryNodePackagesEvent:
		sid, ok := p.Clients[te.NodeID]
		if !ok {
			return mq.Msg{}, fmt.Errorf("no pub client associated with id, %v", te.NodeID)
		}

		return mq.Compose(mq.NodePackages).WithIdentity(sid).WithMetadata(e), nil
	case KillNodeProcessEvent:
		sid, ok := p.Clients[te.NodeID]
		if !ok {
			return mq.Msg{}, fmt.Errorf("no pub client associated with id, %v", te.NodeID)
		}

		return mq.Compose(mq.KillNodeProcess, mq.KillNodeProcessRequest{PID: te.PID}).WithIdentity(sid).WithMetadata(e), nil
	case StartNodeSpeedtestEvent:
		sid, ok := p.Clients[te.NodeID]
		if !ok {
			return mq.Msg{}, fmt.Errorf("no pub client associated with id, %v", te.NodeID)
		}

		return mq.Compose(mq.StartNodeSpeedtest).WithIdentity(sid).WithMetadata(e), nil
	default:
		return mq.Msg{}, fmt.Errorf("couldn't match event with a topic, %v", e)
	}
}

func (p ZeroMQEventPubSub) msgToEventOutput(m mq.Msg) (EventOutput, error) {
	switch m.Topic {
	case mq.NodeConnections:
		d, err := typedMsgData[mq.NodeConnectionsResponse](m)
		return NewQueryNodeConnectionsEventOutput(
			m.Metadata.(Event),
			d.Connections,
			err,
		), nil
	case mq.NodePackages:
		d, err := typedMsgData[mq.NodePackagesResponse](m)
		return NewQueryNodePackagesEventOutput(
			m.Metadata.(Event),
			d.Packages,
			err,
		), nil
	case mq.NodeProcesses:
		d, err := typedMsgData[mq.NodeProcessesResponse](m)
		return NewQueryNodeProcessesEventOutput(
			m.Metadata.(Event),
			d.Processes,
			err,
		), nil
	case mq.KillNodeProcess:
		err := errorMsgData(m)
		return NewKillNodeProcessEventOutput(
			m.Metadata.(Event),
			err,
		), nil
	case mq.StartNodeSpeedtest:
		d, err := typedMsgData[mq.NodeSpeedtestResponse](m)
		return NewNodeSpeedtestEventOutput(
			m.Metadata.(Event),
			d.Speedtest,
			err,
		), nil
	default:
		return nil, fmt.Errorf("couldn't match message with a topic, %v", m)
	}
}

func typedMsgData[T any](
	msg mq.Msg,
) (T, error) {
	if t, ok := msg.Data.(T); ok {
		return t, nil
	}

	var t T
	return t, errorMsgData(msg)
}

func errorMsgData(
	msg mq.Msg,
) error {
	if msg.Data == nil {
		return nil
	}

	if te, ok := msg.Data.(mq.OPError); ok {
		return errors.New(te.Error())
	}

	return fmt.Errorf("failed to understand data type, %v", msg.Data)
}

func init() {
	gob.Register(BaseEvent{})
	gob.Register(QueryNodeConnectionsEvent{})
	gob.Register(QueryNodePackagesEvent{})
	gob.Register(QueryNodeProcessesEvent{})
	gob.Register(KillNodeProcessEvent{})
	gob.Register(StartNodeSpeedtestEvent{})
}
