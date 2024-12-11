package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guackamolly/zero-monitor/internal/banner"
	"github.com/guackamolly/zero-monitor/internal/bootstrap"
	"github.com/guackamolly/zero-monitor/internal/config"
	dbb "github.com/guackamolly/zero-monitor/internal/data/db"
	dbbolt "github.com/guackamolly/zero-monitor/internal/data/db/db-bolt"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/env"
	"github.com/guackamolly/zero-monitor/internal/event"
	"github.com/guackamolly/zero-monitor/internal/http"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/mq"
	"github.com/guackamolly/zero-monitor/internal/service"
	"github.com/guackamolly/zero-monitor/public"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"

	build "github.com/guackamolly/zero-monitor/internal/build"
	flags "github.com/guackamolly/zero-monitor/internal/build/flags"
)

func init() {
	flags.WithMasterFlags()

	if build.Release() && !flags.Verbose() {
		logging.DisableDebugLogs()
	}
	logging.AddLogger(logging.NewConsoleLogger())
	banner.Print()
}

func main() {
	// 1. Load env + config
	env := loadEnv()
	cfg := loadConfig()

	// 2. Initialize DI.
	sc := createServiceContainer(cfg)
	suc := createSubContainer(sc)
	ctx := context.Background()
	ctx = mq.InjectSubscribeContainer(ctx, suc)

	// 3. Initialize sub server.
	loadCrypto(env.MessageQueueTransportPemKey)

	s := initializeSubServer(ctx, env.MessageQueueHost, env.MessageQueuePort)
	defer s.Close()

	// 4. Initialize database.
	db := initializeDatabase(env.BoltDBPath)
	defer db.Close()

	// 5. Initialize http server.
	sc = updateServiceContainer(sc, &s, db)
	ctx = http.InjectServiceContainer(ctx, sc)
	e := initializeHttpServer(ctx, env.ServerHost, env.ServerPort, env.ServerTLSCert, env.ServerTLSKey, env.ServerVirtualHost)
	defer e.Close()

	go logIfAdminNeedsRegistration(sc.Authentication)

	// 6. Await termination...
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	// 7. Try to save config
	saveConfig(sc.MasterConfiguration)
}

func loadEnv() env.MasterEnv {
	if env, err := env.Master(); err == nil {
		return env
	}

	logging.LogDebug("couldn't lookup .env, bootstrapping configuration values...")
	return bootstrap.Master()
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

func loadCrypto(keyfilepath string) {
	err := mq.LoadAsymmetricBlock(keyfilepath)
	if err != nil {
		logging.LogError("failed to load message queue private key, %v", err)
		logging.LogWarning("message queue sensitive data will not be encrypted!")
	}
}

func initializeSubServer(ctx context.Context, host, port string) mq.Socket {
	s := mq.NewSubSocket(ctx)
	s.RegisterSubscriptions()
	err := mq.ConnectSubscribe(s, host, port)
	if err != nil {
		s.Close()
		log.Fatalf("coudln't open zeromq sub socket, %v\n", err)
	}

	logging.LogInfo("⇨ ZeroMQ server started on %s", color.Green(fmt.Sprintf("tcp://%s", s.Addr())))

	return s
}

func initializeHttpServer(
	ctx context.Context,
	host, port, certfilepath, keyfilepath string,
	virtualhost string,
) *echo.Echo {
	// Initialize echo framework.
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	// Register server dependencies.
	http.RegisterHandlers(e)
	http.RegisterMiddlewares(e, ctx)
	http.RegisterStaticFiles(e, public.FS)
	http.RegisterTemplates(e, public.FS)
	http.SetVirtualHost(virtualhost)

	https := len(certfilepath) > 0 && len(keyfilepath) > 0

	// Start server.
	go func() {
		if https {
			logging.LogFatal("server exit %v", http.StartTLS(e, host, port, certfilepath, keyfilepath))
		}

		logging.LogFatal("server exit %v", http.Start(e, host, port))
	}()

	go func() {
		for {
			var addr net.Addr
			if https {
				addr = e.TLSListenerAddr()
			} else {
				addr = e.ListenerAddr()
			}

			if addr == nil {
				time.Sleep(200 * time.Millisecond)
				continue
			}

			if https {
				logging.LogInfo("⇨ https server started on %s", color.Green(fmt.Sprintf("https://%s", addr)))
			} else {
				logging.LogInfo("⇨ http server started on %s", color.Green(fmt.Sprintf("http://%s", addr)))
			}

			return
		}
	}()

	return e
}

func initializeDatabase(dbpath string) dbb.Database {
	db := dbbolt.NewBoltDatabase(dbpath)
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
		JoinNodesNetwork:                   sc.NodeManager.Join,
		UpdateNodesNetwork:                 sc.NodeManager.Update,
		GetNodeStatsPollingDuration:        cfg.NodeStatsPolling.Duration,
		AuthenticateNodesNetwork:           sc.NodeManager.Authenticate,
		RequiresNodesNetworkAuthentication: func(n models.Node) bool { return !sc.NodeManager.IsAuthenticated(n) },
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
	crt, _ := db.Table(dbb.TableCredentials)
	ust, _ := db.Table(dbb.TableUser)

	if !ok {
		logging.LogFatal("table %s hasn't been initialized", dbb.TableSpeedtest)
	}
	sps := repositories.NewDatabaseSpeedtestStoreRepository(stt.(dbb.SpeedtestTable))
	authRepo := repositories.NewDatabaseAuthenticationRepository(crt.(dbb.CredentialsTable), ust.(dbb.UserTable))
	userRepo := repositories.NewDatabaseUserRepository(ust.(dbb.UserTable))

	tokens := service.TokenBucket{}

	sc.NodeCommander = service.NewNodeCommanderService(zps, zps)
	sc.NodeSpeedtest = service.NewNodeSpeedtestService(zps, zps, sps)
	sc.Network = service.NewNetworkService(zps)
	sc.Networking = service.NewNetworkingService()
	sc.Authentication = service.NewAuthenticationService(authRepo, userRepo, &tokens)
	sc.Authorization = service.NewAuthorizationService(&tokens)

	return sc
}

func logIfAdminNeedsRegistration(
	as *service.AuthenticationService,
) {
	if !as.NeedsAdminRegistration() {
		return
	}

	logging.LogWarning("no admin account has been registered yet! register one now at /dashboard")
}
