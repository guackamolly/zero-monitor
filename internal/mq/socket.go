package mq

import (
	"context"
	"log"

	"github.com/go-zeromq/zmq4"
)

// Wraps [zmq4.Socket] so context can be accessed for extracting dependencies.
type Socket struct {
	zmq4.Socket
	ctx context.Context
}

// Creates a new [zmq4-Socket] wrapper with a custom context.
// The context must contain all the dependencies required by the socket.
func NewSocket(ctx context.Context) Socket {
	return Socket{
		Socket: zmq4.NewRep(ctx),
		ctx:    ctx,
	}
}

func (s Socket) Context() context.Context {
	return s.ctx
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
