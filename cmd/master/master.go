package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guackamolly/zero-monitor/internal/config"
	dbb "github.com/guackamolly/zero-monitor/internal/data/db"
	dbbolt "github.com/guackamolly/zero-monitor/internal/data/db/db-bolt"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
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
	mq.LoadAsymmetricBlock(false)
	s := initializeSubServer(ctx)
	defer s.Close()

	// 4. Initialize database.
	db := initializeDatabase()
	defer db.Close()

	// 5. Initialize http server.
	sc = updateServiceContainer(sc, &s, db)
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
	s := mq.NewSubSocket(ctx)
	s.RegisterSubscriptions()
	err := mq.ConnectSubscribe(s)
	if err != nil {
		s.Close()
		log.Fatalf("coudln't open zeromq sub socket, %v\n", err)
	}
	log.Printf("started zeromq sub socket on addr: %s\n", s.Addr())

	return s
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

func initializeDatabase() dbb.Database {
	db := dbbolt.NewBoltDatabase(dbbolt.Path())
	err := db.Open()
	if err != nil {
		logging.LogFatal("couldn't initialize database, %v", db)
	}

	return db
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
	db dbb.Database,
) http.ServiceContainer {
	zps := event.NewZeroMQEventPubSub(s)
	stt, ok := db.Table(dbb.TableSpeedtest)
	if !ok {
		logging.LogFatal("table %s hasn't been initialized", dbb.TableSpeedtest)
	}
	sps := repositories.NewDatabaseSpeedtestStoreRepository(stt.(dbb.SpeedtestTable))

	sc.NodeCommander = service.NewNodeCommanderService(zps, zps)
	sc.NodeSpeedtest = service.NewNodeSpeedtestService(zps, zps, sps)

	return sc
}
