package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/mq"
	"github.com/guackamolly/zero-monitor/internal/service"
	"github.com/showwin/speedtest-go/speedtest"

	_ "github.com/guackamolly/zero-monitor/internal/build"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env
	loadEnv()

	// 1. Initialize DI & logging.
	pc := createPublishContainer()
	ctx := context.Background()
	ctx = mq.InjectPublishContainer(ctx, pc)

	logging.AddLogger(logging.NewConsoleLogger())

	// 2. Initialize pub server.
	loadCrypto()
	s := mq.NewPubSocket(ctx)
	defer s.Close()

	err := mq.ConnectPublish(s)
	if err != nil {
		log.Fatalf("coudln't open zeromq pub socket, %v\n", err)
	}
	log.Printf("started zeromq pub socket on addr: %s\n", s.Addr())
	s.RegisterPublishers()

	// 3. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func loadEnv() {
	if err := godotenv.Load(".env"); err == nil {
		return
	}

	if d, err := config.Dir(); err == nil && godotenv.Load(filepath.Join(d, "node.env")) == nil {
		return
	}

	log.Fatalf("couldn't load .env or $CFG_DIR/node.env!")
}

func loadCrypto() {
	err := mq.LoadAsymmetricBlock(true)
	if err != nil {
		logging.LogFatal("failed to load message queue public key, %v", err)
	}
}

func createPublishContainer() mq.PublishContainer {
	system := repositories.NewGopsUtilSystemRepository()
	speedtest := repositories.NewNetSpeedtestRepository(speedtest.New())
	nrs := service.NewNodeReporterService(system, speedtest)

	return mq.PublishContainer{
		GetCurrentNode:            nrs.Node,
		GetCurrentNodeConnections: nrs.Connections,
		GetCurrentNodeProcesses:   nrs.Processes,
		GetCurrentNodePackages:    nrs.Packages,
		StartNodeStatsPolling:     nrs.Start,
		UpdateNodeStatsPolling:    nrs.Update,
		KillNodeProcess:           nrs.KillProcess,
		StartNodeSpeedtest:        nrs.Speedtest,
	}
}
