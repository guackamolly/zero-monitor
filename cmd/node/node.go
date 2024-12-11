package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/guackamolly/zero-monitor/internal/banner"
	"github.com/guackamolly/zero-monitor/internal/bootstrap"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/env"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/mq"
	"github.com/guackamolly/zero-monitor/internal/service"
	"github.com/showwin/speedtest-go/speedtest"

	build "github.com/guackamolly/zero-monitor/internal/build"
	flags "github.com/guackamolly/zero-monitor/internal/build/flags"
)

func init() {
	flags.WithNodeFlags()

	if build.Release() && !flags.Verbose() {
		logging.DisableDebugLogs()
	}
	logging.AddLogger(logging.NewConsoleLogger())
	banner.Print()
}

func main() {
	// 1. Load env
	env := loadEnv()

	// 2. Initialize DI.
	pc := createPublishContainer()
	ctx := context.Background()
	ctx = mq.InjectPublishContainer(ctx, pc)

	// 3. Initialize pub server.
	loadCrypto(env.MessageQueueTransportPubKey)
	s := mq.NewPubSocket(ctx)
	defer s.Close()

	err := mq.ConnectPublish(s, env.MessageQueueHost, env.MessageQueuePort)
	if err != nil {
		log.Fatalf("coudln't open zeromq pub socket, %v\n", err)
	}
	log.Printf("started zeromq pub socket on addr: %s\n", s.Addr())
	s.RegisterPublishers(env.MessageQueueInviteCode)

	// 4. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func loadEnv() env.NodeEnv {
	if env, err := env.Node(); err == nil && len(flags.InviteLink()) == 0 {
		return env
	}

	logging.LogDebug("couldn't lookup .env, bootstrapping configuration values...")
	return bootstrap.Node(flags.InviteLink())
}

func loadCrypto(keyfilepath string) {
	err := mq.LoadAsymmetricBlock(keyfilepath)
	if err != nil {
		logging.LogError("failed to load message queue public key, %v", err)
		logging.LogWarning("message queue sensitive data will not be encrypted!")
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
