package mq

import (
	"fmt"
	"net"
	"os"

	"github.com/guackamolly/zero-monitor/internal/conn"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

const (
	mqSubHostEnvKey = "mq_sub_host"
	mqSubPortEnvKey = "mq_sub_port"
)

var (
	mqSubHost = os.Getenv(mqSubHostEnvKey)
	mqSubPort = os.Getenv(mqSubPortEnvKey)
)

// Connects a socket for publishing messages to master node.
func ConnectPublish(s Socket) error {
	if len(mqSubHost) > 0 && len(mqSubPort) > 0 {
		return s.Dial(fmt.Sprintf("tcp://[%s]:%s", mqSubHost, mqSubPort))
	}

	conn, err := conn.StartBeaconBroadcast()
	if err != nil {
		return err
	}

	return s.Dial(fmt.Sprintf("tcp://[%s]:%d", conn.IP, conn.Port))
}

// Connects a socket for subscribing messages from reporting nodes.
func ConnectSubscribe(s Socket) error {
	ip := subHostIP()

	// if port is unspecified, default to 0 so go internals
	// choose a random available port
	if len(mqSubPort) == 0 {
		mqSubPort = "0"
	}

	return s.Listen(fmt.Sprintf("tcp://[%s]:%s", ip, mqSubPort))
}

func Close(s Socket) error {
	return s.Close()
}

func subHostIP() net.IP {
	if ip := net.ParseIP(mqSubHost); ip != nil {
		return ip
	}

	iaddrs, err := net.InterfaceAddrs()
	if err != nil {
		logging.LogWarning("couldn't lookup interface addresses. defaulting to any address, %v", err)
		return net.IPv4(0, 0, 0, 0)
	}

	for _, iaddr := range iaddrs {
		c, ok := iaddr.(*net.IPNet)
		if ok && c.IP.IsPrivate() {
			return c.IP
		}
	}

	logging.LogWarning("couldn't find the private ip of the closest network interface card. defaulting to any address, %v", err)
	return net.IPv4(0, 0, 0, 0)
}
