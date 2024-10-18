package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/conn"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/http"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/mq"
	"github.com/guackamolly/zero-monitor/internal/service"
	"github.com/labstack/echo/v4"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 1. Load config
	cfg := loadConfig()

	// 2. Initialize DI.
	sc := createSubscribeContainer(cfg)
	ctx := context.Background()
	ctx = di.InjectSubscribeContainer(ctx, sc)

	// 3. Initialize sub server.
	s := initializeSubServer(ctx)
	defer s.Close()

	// 4. Initialize beacon server.
	uconn := initializeBeaconServer()
	defer uconn.Close()

	// 5. Initialize http server.
	e := initializeHttpServer(ctx)
	defer e.Close()

	// 6. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	// 7. Try to save config
	saveConfig(cfg)
}

func loadConfig() config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("failed to load config, %v", err)
	}

	return cfg
}

func saveConfig(cfg config.Config) {
	err := config.Save(cfg)
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

func initializeBeaconServer() *net.UDPConn {
	// Find available UDP ports.
	uconn := findAvailableUdpPort()

	// Initialize beacon server.
	conn.StartBeaconServer(uconn)
	log.Printf("started udp beacon server on addr: %s\n", uconn.LocalAddr())

	return uconn
}

func initializeHttpServer(ctx context.Context) *echo.Echo {
	// Initialize echo framework.
	e := echo.New()

	// Initialize logging.
	logging.AddLogger(logging.NewConsoleLogger())
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

func createSubscribeContainer(cfg config.Config) di.SubscribeContainer {
	ns := make([]models.Node, len(cfg.TrustedNetwork))
	i := 0
	for _, n := range cfg.TrustedNetwork {
		ns[i] = n
		i++
	}

	return di.SubscribeContainer{
		NodeManager: service.NewNodeManagerService(ns...),
	}
}
