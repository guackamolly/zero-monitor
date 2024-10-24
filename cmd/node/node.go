package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/guackamolly/zero-monitor/internal/conn"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/mq"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func main() {
	// 1. Initialize DI & logging.
	pc := createPublishContainer()
	ctx := context.Background()
	ctx = mq.InjectPublishContainer(ctx, pc)

	logging.AddLogger(logging.NewConsoleLogger())

	// 2. Find master node in local network.
	conn, err := conn.StartBeaconBroadcast()
	if err != nil {
		log.Fatalf("failed to probe master node, %v\n", err)
	}

	// 3. Initialize pub server.
	s := mq.NewPubSocket(ctx)
	defer s.Close()

	err = mq.ConnectPublish(s, conn.IP, conn.Port)
	if err != nil {
		log.Fatalf("coudln't open zeromq pub socket, %v\n", err)
	}
	log.Printf("started zeromq pub socket on addr: %s\n", s.Addr())
	s.RegisterPublishers()

	// 4. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func createPublishContainer() mq.PublishContainer {
	system := repositories.GopsUtilSystemRepository{}
	nrs := service.NewNodeReporterService(system)

	return mq.PublishContainer{
		GetCurrentNode:            nrs.Node,
		GetCurrentNodeConnections: nrs.Temp,
		StartNodeStatsPolling:     nrs.Start,
		UpdateNodeStatsPolling:    nrs.Update,
	}
}
