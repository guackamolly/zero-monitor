package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/conn"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/event"
	"github.com/guackamolly/zero-monitor/internal/http"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/mq"
	"github.com/guackamolly/zero-monitor/internal/service"
	"github.com/labstack/echo/v4"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 1. Load config
	logging.AddLogger(logging.NewConsoleLogger())
	cfg := loadConfig()

	// 2. Initialize DI.
	sc := createServiceContainer(cfg)
	suc := createSubContainer(sc)
	ctx := context.Background()
	ctx = mq.InjectSubscribeContainer(ctx, suc)

	// 3. Initialize sub server.
	s := initializeSubServer(ctx)
	defer s.Close()

	// 4. Initialize beacon server.
	uconn := initializeBeaconServer(addrToConn(s.Addr()))
	defer uconn.Close()

	// 5. Initialize http server.
	sc = updateServiceContainer(sc, &s)
	ctx = http.InjectServiceContainer(ctx, sc)
	e := initializeHttpServer(ctx)
	defer e.Close()

	// 6. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	// 7. Try to save config
	saveConfig(sc.MasterConfiguration)
}

func loadConfig() config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("failed to load config, %v", err)
	}

	return cfg
}

func saveConfig(s *service.MasterConfigurationService) {
	err := s.Save()
	if err != nil {
		log.Printf("failed to save config, %v", err)
	}
}

func initializeSubServer(ctx context.Context) mq.Socket {
	// Find available TCP port.
	tconn := findAvailableTcpPort()
	taddr := tconn.Addr().(*net.TCPAddr)

	// Initialize sub server.
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

func initializeBeaconServer(subConn conn.Connection) *net.UDPConn {
	// Find available UDP ports.
	uconn := findAvailableUdpPort()

	// Initialize beacon server.
	conn.StartBeaconServer(uconn, subConn)
	log.Printf("started udp beacon server on addr: %s\n", uconn.LocalAddr())

	return uconn
}

func initializeHttpServer(ctx context.Context) *echo.Echo {
	// Initialize echo framework.
	e := echo.New()

	// Initialize logging.
	logging.AddLogger(logging.NewEchoLogger(e.Logger))

	// Register server dependencies.
	http.RegisterHandlers(e)
	http.RegisterMiddlewares(e, ctx)
	http.RegisterStaticFiles(e)
	http.RegisterTemplates(e)

	// Start server.
	go func() {
		logging.LogFatal("server exit %v", http.Start(e))
	}()

	return e
}

func addrToConn(addr net.Addr) conn.Connection {
	switch c := addr.(type) {
	case *net.TCPAddr:
		return conn.Connection{Port: c.Port, IP: c.IP}
	case *net.UDPAddr:
		return conn.Connection{Port: c.Port, IP: c.IP}
	}

	logging.LogFatal("couldn't convert %v to conn.Connection", addr)
	return conn.Connection{}
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

func createServiceContainer(
	cfg config.Config,
) http.ServiceContainer {
	ns := make([]models.Node, len(cfg.TrustedNetwork))
	i := 0
	for _, n := range cfg.TrustedNetwork {
		ns[i] = n
		i++
	}

	mcs := service.NewMasterConfigurationService(&cfg)
	nms := service.NewNodeManagerService(ns...)
	nss := service.NewNodeSchedulerService(
		mcs.Current,
		mcs.Stream,
		mcs.Save,
		mcs.UpdateTrustedNetwork,
		nms.Update,
		nms.Network,
		nms.Stream,
	)

	return http.ServiceContainer{
		NodeManager:         nms,
		NodeScheduler:       nss,
		MasterConfiguration: mcs,
	}
}

func createSubContainer(sc http.ServiceContainer) mq.SubscribeContainer {
	cfg := sc.MasterConfiguration.Current()
	return mq.SubscribeContainer{
		JoinNodesNetwork:            sc.NodeManager.Join,
		UpdateNodesNetwork:          sc.NodeManager.Update,
		GetNodeStatsPollingDuration: cfg.NodeStatsPolling.Duration,
		GetNodeStatsPollingDurationUpdates: func() chan (time.Duration) {
			ch := make(chan (time.Duration))
			sp := cfg.NodeStatsPolling.Duration()
			cfgs := sc.MasterConfiguration.Stream()
			go func() {
				for cfg = range cfgs {
					usp := cfg.NodeStatsPolling.Duration()
					if usp == sp {
						continue
					}
					sp = usp
					ch <- sp
				}
			}()

			return ch
		},
	}
}

func updateServiceContainer(
	sc http.ServiceContainer,
	s *mq.Socket,
) http.ServiceContainer {
	zps := event.NewZeroMQEventPubSub(s)
	sc.NodeCommander = service.NewNodeCommanderService(zps, zps)
	return sc
}
