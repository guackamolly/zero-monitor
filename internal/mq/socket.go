package mq

import (
	"context"
	"fmt"
	"log"

	"github.com/go-zeromq/zmq4"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Wraps [zmq4.Socket] so context can be accessed for extracting dependencies.
type Socket struct {
	zmq4.Socket
	ctx          context.Context
	framesLength int
	listeners    map[Topic][]func(Msg)
	// Clients of the socket, identified by their machine ID and their socket identity
	// This field only makes sense for sub sockets.
	Clients map[string][]byte
}

// Creates a new sub [zmq4-Socket] wrapper with a custom context.
// The context must contain all the dependencies required by the socket.
func NewSubSocket(ctx context.Context) Socket {
	return Socket{
		Socket:       zmq4.NewRouter(ctx),
		ctx:          ctx,
		framesLength: 2,
		listeners:    map[Topic][]func(Msg){},
		Clients:      map[string][]byte{},
	}
}

// Creates a new sub [zmq4-Socket] wrapper with a custom context.
// The context must contain all the dependencies required by the socket.
func NewPubSocket(ctx context.Context) Socket {
	return Socket{
		Socket:       zmq4.NewDealer(ctx, zmq4.WithAutomaticReconnect(true)),
		ctx:          ctx,
		framesLength: 1,
	}
}

func (s Socket) Context() context.Context {
	return s.ctx
}

func (s Socket) OnMsgReceived(t Topic, h func(m Msg)) {
	hs, ok := s.listeners[t]
	if !ok {
		hs = []func(Msg){}
	}

	s.listeners[t] = append(hs, h)
}

// Publishes a message from a pub socket to the sub socket and waits for a reply.
// Receiver must be a pub socket.
func (s Socket) PublishMsg(m Msg) error {
	logging.LogInfo("publishing Msg with topic: %d", m.Topic)

	b, err := encode(m)
	if err != nil {
		return err
	}

	err = s.Send(zmq4.NewMsg(b))
	if err != nil {
		return err
	}

	return nil
}

func (s Socket) ReceiveMsg() (Msg, error) {
	zm, err := s.Recv()
	if err != nil {
		return Msg{}, err
	}

	if l := len(zm.Frames); l != s.framesLength {
		err = fmt.Errorf("received corrupted message, expected %d frames but got: %d", s.framesLength, l)
		return Msg{}, err
	}

	if s.framesLength == 1 {
		return decode(zm.Frames[0])
	}

	m, err := decode(zm.Frames[1])
	if err != nil {
		return Msg{}, err
	}

	m = m.WithIdentity(zm.Frames[0])
	s.onMsgReceived(m)

	return m, nil
}

// Replies to a pub socket from the sub socket.
// Receiver must be a sub socket.
func (s Socket) ReplyMsg(id []byte, m Msg) {
	b, err := encode(m)
	if err != nil {
		log.Printf("failed to encode reply message, %v\n", err)
	}

	err = s.Send(zmq4.NewMsgFrom(id, b))
	if err != nil {
		log.Printf("failed to reply message, %v\n", err)
	}
}

func (s Socket) onMsgReceived(
	msg Msg,
) {
	hs, ok := s.listeners[msg.Topic]
	if !ok {
		return
	}

	logging.LogInfo("calling %d handlers", len(hs))
	for _, h := range hs {
		h(msg)
	}
}
