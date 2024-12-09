package mq

import (
	"context"
	"fmt"
	"log"

	"github.com/go-zeromq/zmq4"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Wraps a callback that is executed when a message is received, in order to notify a listener.
type messageListener struct {
	handler func(Msg)
}

// Wraps [zmq4.Socket] so context can be accessed for extracting dependencies.
type Socket struct {
	zmq4.Socket
	ctx          context.Context
	framesLength int
	// Listeners of specific topics, beside the actual socket message handler.
	// This field only makes sense for sub sockets.
	listeners map[Topic][]*messageListener
	// Clients of the socket, identified by their machine ID and their socket identity
	// This field only makes sense for sub sockets.
	Clients map[string][]byte
	// ZeroMQ identity that is bundled on message frames
	Identity []byte
}

// Creates a new sub [zmq4-Socket] wrapper with a custom context.
// The context must contain all the dependencies required by the socket.
func NewSubSocket(ctx context.Context) Socket {
	id := []byte(models.UUID())
	return Socket{
		Socket:       zmq4.NewRouter(ctx, zmq4.WithID(id)),
		ctx:          ctx,
		framesLength: 2,
		listeners:    map[Topic][]*messageListener{},
		Clients:      map[string][]byte{},
		Identity:     id,
	}
}

// Creates a new sub [zmq4-Socket] wrapper with a custom context.
// The context must contain all the dependencies required by the socket.
func NewPubSocket(ctx context.Context) Socket {
	id := []byte(models.UUID())
	return Socket{
		Socket:       zmq4.NewDealer(ctx, zmq4.WithAutomaticReconnect(true), zmq4.WithID(id)),
		ctx:          ctx,
		framesLength: 1,
		Identity:     id,
	}
}

func (s Socket) Context() context.Context {
	return s.ctx
}

// Registers a listener on this socket for a particular topic.
// Returns a callback that the listener should use to cancel the subscription to new messages of the topic.
func (s Socket) OnMsgReceived(t Topic, h func(m Msg)) func() {
	ls, ok := s.listeners[t]
	if !ok {
		ls = []*messageListener{}
	}

	l := &messageListener{handler: h}
	s.listeners[t] = append(ls, l)

	return func() {
		ls = []*messageListener{}
		for _, lh := range s.listeners[t] {
			if l == lh {
				continue
			}

			ls = append(ls, lh)
		}

		s.listeners[t] = ls
	}
}

// Publishes a message from a pub socket to the sub socket and waits for a reply.
// Receiver must be a pub socket.
func (s Socket) PublishMsg(m Msg) error {
	var err error
	logging.LogDebug("publishing Msg with topic: %d", m.Topic)

	m = m.WithIdentity(s.Identity)
	if m.Topic == HelloNetwork {
		return s.publishHelloNetworkMsg(m)
	}

	if m.Topic.Sensitive() {
		m, err = m.Encrypt()
	}

	if err != nil {
		return err
	}

	b, err := models.Encode(m)
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
		return models.Decode[Msg](zm.Frames[0])
	}

	m, err := models.Decode[Msg](zm.Frames[1])
	if err != nil {
		return Msg{}, err
	}

	if m.Topic == HelloNetwork {
		return s.interceptHelloNetworkMsg(m)
	}

	if m.Topic.Sensitive() {
		m, err = m.Decrypt()
	}

	if err != nil {
		return Msg{}, err
	}

	s.onMsgReceived(m)
	return m, nil
}

// Replies to a pub socket from the sub socket.
// Receiver must be a sub socket.
func (s Socket) ReplyMsg(id []byte, m Msg) {
	b, err := models.Encode(m)
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
	ls, ok := s.listeners[msg.Topic]
	if !ok {
		return
	}

	logging.LogDebug("calling %d handlers", len(ls))
	for _, l := range ls {
		l.handler(msg)
	}
}

func (s Socket) publishHelloNetworkMsg(
	m Msg,
) error {
	key, err := GenerateCipherKey()
	if err != nil {
		return err
	}

	err = RegisterCipherKey(s.Identity, key)
	if err != nil {
		return err
	}

	key, err = EncryptAsymmetric(key)
	if err != nil {
		return err
	}

	b, err := models.Encode(m)
	if err != nil {
		return err
	}

	em := Msg{
		Identity: s.Identity,
		Topic:    m.Topic,
		Data:     b,
		Metadata: key,
	}

	b, err = models.Encode(em)
	if err != nil {
		return err
	}

	return s.Send(zmq4.NewMsg(b))
}

func (s Socket) interceptHelloNetworkMsg(
	m Msg,
) (Msg, error) {
	var err error
	key, ok := m.Metadata.([]byte)
	if !ok {
		return Msg{}, fmt.Errorf("key is not a bitstream")
	}

	bs, ok := m.Data.([]byte)
	if !ok {
		return Msg{}, fmt.Errorf("data is not a bitstream")
	}

	key, err = DecryptAsymmetric(key)
	if err != nil {
		return Msg{}, err
	}

	err = RegisterCipherKey(m.Identity, key)
	if err != nil {
		return Msg{}, err
	}

	m, err = models.Decode[Msg](bs)
	if err != nil {
		return Msg{}, err
	}

	return m, nil
}
