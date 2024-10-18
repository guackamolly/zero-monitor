package mq

import (
	"context"
	"log"

	"github.com/go-zeromq/zmq4"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Wraps [zmq4.Socket] so context can be accessed for extracting dependencies.
type Socket struct {
	zmq4.Socket
	ctx context.Context
}

// Creates a new sub [zmq4-Socket] wrapper with a custom context.
// The context must contain all the dependencies required by the socket.
func NewSubSocket(ctx context.Context) Socket {
	return Socket{
		Socket: zmq4.NewRouter(ctx),
		ctx:    ctx,
	}
}

// Creates a new sub [zmq4-Socket] wrapper with a custom context.
// The context must contain all the dependencies required by the socket.
func NewPubSocket(ctx context.Context) Socket {
	return Socket{
		Socket: zmq4.NewReq(ctx, zmq4.WithAutomaticReconnect(true)),
		ctx:    ctx,
	}
}

func (s Socket) Context() context.Context {
	return s.ctx
}

func (s Socket) PublishAndForget(m msg) {
	logging.LogInfo("publishing msg with topic: %d", m.Topic)

	b, err := encode(m)
	if err != nil {
		log.Printf("failed to encode message, %v\n", err)
	}

	err = s.Send(zmq4.NewMsg(b))
	if err != nil {
		log.Printf("failed to publish message, %v\n", err)
	}
}

func (s Socket) Reply(m msg) {
	b, err := encode(m)
	if err != nil {
		log.Printf("failed to encode reply message, %v\n", err)
	}

	err = s.Send(zmq4.NewMsg(b))
	if err != nil {
		log.Printf("failed to reploy message, %v\n", err)
	}
}
