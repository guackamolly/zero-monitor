package mq

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/conn"
)

// Connects a socket for publishing messages to master node.
func ConnectPublish(s Socket, conn conn.Connection) error {
	addr := fmt.Sprintf("tcp://%s:%d", conn.IP, conn.Port)

	return s.Dial(addr)
}

// Connects a socket for subscribing messages from reporting nodes.
func ConnectSubscribe(s Socket, conn conn.Connection) error {
	addr := fmt.Sprintf("tcp://%s:%d", conn.IP, conn.Port)

	return s.Listen(addr)
}

func Close(s Socket) error {
	return s.Close()
}
