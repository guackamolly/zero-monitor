package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/guackamolly/zero-monitor/internal/conn"
	"github.com/guackamolly/zero-monitor/internal/mq"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func main() {
	// 1. Initialize DI.
	sc := createSubscribeContainer()
	ctx := context.Background()
	ctx = mq.InjectSubscribeContainer(ctx, sc)

	// 2. Initialize sub server.
	s := initializeSubServer(ctx)
	defer s.Close()

	// 3. Initialize beacon server.
	uconn := initializeBeaconServer()
	defer uconn.Close()

	// 4. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func initializeSubServer(ctx context.Context) mq.Socket {
	// 1. Find available TCP port.
	tconn := findAvailableTcpPort()
	taddr := tconn.Addr().(*net.TCPAddr)

	// 2. Initialize sub server.
	tconn.Close()

	s := mq.NewSubSocket(ctx)
	s.RegisterSubscriptions()
	err := mq.ConnectSubscribe(s, taddr.IP, taddr.Port)
	if err != nil {
		s.Close()
		log.Fatalf("coudln't open zeromq sub socket, %v\n", err)
	}
	log.Printf("started zeromq sub socket on addr: %s\n", s.Addr())

	return s
}

func initializeBeaconServer() *net.UDPConn {
	// 1. Find available UDP ports.
	uconn := findAvailableUdpPort()

	// 2. Initialize beacon server.
	conn.StartBeaconServer(uconn)
	log.Printf("started udp beacon server on addr: %s\n", uconn.LocalAddr())

	return uconn
}

func findAvailableTcpPort() *net.TCPListener {
	log.Println("finding a TCP port available for incoming requests...")
	tconn, err := conn.FindAvailableTcpPort(conn.NetworkIP)
	if err != nil {
		log.Fatalf("couldn't find a TCP port available for incoming requests, %err\n", err)
	}

	return tconn
}

func findAvailableUdpPort() *net.UDPConn {
	log.Println("finding an UDP port available for incoming requests...")
	uconn, err := conn.FindAvailableUdpPort(conn.NetworkIP)
	if err != nil {
		log.Fatalf("couldn't find an UDP port available for incoming requests, %err\n", err)
	}

	return uconn
}

func createSubscribeContainer() mq.SubscribeContainer {
	return mq.SubscribeContainer{
		NodeManager: service.NewNodeManagerService(),
	}
}
