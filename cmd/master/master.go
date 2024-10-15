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

	// 2. Find available TCP & UDP ports.
	tconn, uconn := findAvailableTcpAndUdpPorts()
	taddr := tconn.Addr().(*net.TCPAddr)

	// 3. Initialize sub server.
	tconn.Close()

	s := mq.NewSubSocket(ctx)
	s.RegisterSubscriptions()
	err := mq.ConnectSubscribe(s, taddr.IP, taddr.Port)
	if err != nil {
		uconn.Close()
		log.Fatalf("coudln't open zeromq sub socket, %v\n", err)
	}
	defer s.Close()
	log.Printf("started zeromq sub socket on addr: %s\n", s.Addr())

	// 4. Initialize beacon server.
	conn.StartBeaconServer(uconn)
	defer uconn.Close()
	log.Printf("started udp beacon server on addr: %s\n", uconn.LocalAddr())

	// 5. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func findAvailableTcpAndUdpPorts() (*net.TCPListener, *net.UDPConn) {
	log.Println("finding a TCP port available for incoming requests...")
	tconn, err := conn.FindAvailableTcpPort(conn.NetworkIP)
	if err != nil {
		log.Fatalf("couldn't find a TCP port available for incoming requests, %err\n", err)
	}

	log.Println("finding an UDP port available for incoming requests...")
	uconn, err := conn.FindAvailableUdpPort(conn.NetworkIP)
	if err != nil {
		log.Fatalf("couldn't find an UDP port available for incoming requests, %err\n", err)
	}

	return tconn, uconn
}

func createSubscribeContainer() mq.SubscribeContainer {
	return mq.SubscribeContainer{
		NodeManager: service.NewNodeManagerService(),
	}
}
