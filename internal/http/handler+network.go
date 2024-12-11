package http

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/service"
	"github.com/labstack/echo/v4"
)

// GET /network
func networkHandler(ectx echo.Context) error {
	if upgrader.WantsToUpgrade(*ectx.Request()) {
		return networkWebsocketHandler(ectx)
	}

	if _, ok := extractQuery(ectx, joinQueryParam); ok {
		return networkJoinHandler(ectx)
	}

	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		view := NewNetworkView(
			sc.NodeManager.Network(),
			sc.MasterConfiguration.Current().NodeStatsPolling.Duration(),
		)

		return ectx.Render(200, "network", view)
	})
}

// GET /network?join=...
func networkJoinHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withJoinCode(ectx, sc, func(code string) error {
			q := map[string]string{joinQueryParam: code}
			v := NewNetworkJoinView(
				URL(ectx, networkPublicKeyRoute, q).String(),
				URL(ectx, networkConnectionEndpointRoute, q).String(),
			)

			return ectx.JSON(200, v)
		})
	})
}

// GET /network/public-key?join=...
func networkPublicKeyHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withJoinCode(ectx, sc, func(code string) error {
			key, err := sc.Network.PublicKey()
			if err != nil {
				return echo.ErrNotFound
			}

			return ectx.String(200, string(key))
		})
	})
}

// GET /network/connection-endpoint?join=...
func networkConnectionEndpointHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withJoinCode(ectx, sc, func(code string) error {
			addr := sc.Network.Address()

			if !IsLocalRequest(ectx) {
				return ectx.JSON(200, NewNetworkConnectionEndpointView(ExtractHost(ectx), int(addr.Port)))
			}

			var ip net.IP
			var err error
			if !IsLocalRequest(ectx) {
				ip, err = sc.Networking.PublicIP()
			} else if addr.Network() {
				ip, err = sc.Networking.PrivateIP()
			} else {
				ip = net.IP(addr.IP)
			}

			if err != nil {
				return echo.ErrInternalServerError
			}

			return ectx.JSON(200, NewNetworkConnectionEndpointView(ip.String(), int(addr.Port)))
		})
	})
}

// GET /network/:id
func networkIdHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			return ectx.Render(200, "network/:id", NewNetworkNodeInformationView(ectx, n))
		})
	})
}

// GET /network/:id/connections
func networkIdConnectionsHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			if !n.Online {
				return ectx.Render(200, "network/:id/connections", NewNetworkNodeConnectionsView(n, []models.Connection{}, nil))
			}

			logging.LogDebug("fetching node connections")
			conns, err := sc.NodeCommander.Connections(n.ID)
			if err != nil {
				logging.LogError("failed to fetch node connections, %v", conns)
			}

			return ectx.Render(200, "network/:id/connections", NewNetworkNodeConnectionsView(n, conns, err))
		})
	})
}

// GET /network/:id/packages
func networkIdPackagesHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			rerr, rok := FromRedirectWithError(ectx)

			if !n.Online {
				return ectx.Render(200, "network/:id/packages", NewNetworkNodePackagesView(n, []models.Package{}, rerr))
			}

			logging.LogDebug("fetching node packages")
			packages, err := sc.NodeCommander.Packages(n.ID)
			if err != nil {
				logging.LogError("failed to fetch node packages, %v", packages)
			}

			if rok {
				err = errors.Join(err, rerr)
			}

			return ectx.Render(200, "network/:id/packages", NewNetworkNodePackagesView(n, packages, err))
		})
	})
}

// GET /network/:id/processes
func networkIdProcessesHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			rerr, rok := FromRedirectWithError(ectx)

			if !n.Online {
				return ectx.Render(200, "network/:id/processes", NewNetworkNodeProcessesView(n, []models.Process{}, rerr))
			}

			logging.LogDebug("fetching node processes")
			procs, err := sc.NodeCommander.Processes(n.ID)
			if err != nil {
				logging.LogError("failed to fetch node processes, %v", procs)
			}

			if rok {
				err = errors.Join(err, rerr)
			}

			return ectx.Render(200, "network/:id/processes", NewNetworkNodeProcessesView(n, procs, err))
		})
	})
}

// POST /network/:id/processes
func networkIdProcessesFormHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			if !n.Online {
				return ectx.Redirect(301, ectx.Request().URL.Path)
			}

			killPid := ectx.FormValue("kill")
			pid, err := strconv.Atoi(killPid)
			if err != nil {
				logging.LogError("failed to convert pid %s to int, %v", pid, err)
				return RedirectWithError(ectx, fmt.Errorf("process with PID %s couldn't be found", killPid))
			}

			logging.LogDebug("killing node process")
			err = sc.NodeCommander.KillProcess(n.ID, int32(pid))
			if err != nil {
				logging.LogError("failed to kill node process, %v", err)
				return RedirectWithError(ectx, err)
			}

			return ectx.Redirect(301, ectx.Request().URL.Path)
		})
	})
}

// GET /network/:id/speedtest
func networkIdSpeedtestHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			if !n.Online {
				return ectx.Render(200, "network/:id/speedtest", NewStartNetworkNodeSpeedtestView(n, nil))
			}

			return ectx.Render(200, "network/:id/speedtest", NewStartNetworkNodeSpeedtestView(n, nil))
		})
	})
}

// POST /network/:id/speedtest
func networkIdSpeedtestFormHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			if !n.Online {
				return ectx.Redirect(301, ectx.Request().URL.Path)
			}

			logging.LogDebug("starting speedtest")
			st, err := sc.NodeSpeedtest.Start(n.ID)
			if err != nil {
				logging.LogError("failed to start speedtest, %v", err)
				return RedirectWithError(ectx, err)
			}

			return ectx.Redirect(301, ectx.Request().URL.JoinPath(st.ID).Path)
		})
	})
}

// GET /network/:id/speedtest/history
func networkIdSpeedtestHistoryHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			sts, ok := sc.NodeSpeedtest.History(n.ID)
			if !ok {
				logging.LogWarning("no history found for node %s", n.ID)
				sts = []models.Speedtest{}
			}
			return ectx.Render(200, "network/:id/speedtest/history", NewNetworkNodeSpeedtestHistoryView(
				ectx, n, sts, service.SpeedtestHistoryLimit, nil),
			)
		})
	})
}

// GET /network/:id/speedtest/history/chart
func networkIdSpeedtestHistoryChartHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			bp, ok := ExtractBreakpoint(ectx)
			if !ok {
				bp = NewContextView(ectx).Breakpoint()
			}

			sts, ok := sc.NodeSpeedtest.History(n.ID)
			if !ok {
				logging.LogWarning("no history found for node %s", n.ID)
				sts = []models.Speedtest{}
			}
			return ectx.Render(200, "network/:id/speedtest/history/chart",
				NewSpeedtestHistoryChartView(EligibleSpeedtestsForChartView(sts), bp),
			)
		})
	})
}

// GET /network/:id/speedtest/:id
func networkIdSpeedtestIdHandler(ectx echo.Context) error {
	if upgrader.WantsToUpgrade(*ectx.Request()) {
		return networkIdSpeedtestWebsocketHandler(ectx)
	}

	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			return withSpeedtest(ectx, sc, func(st models.Speedtest) error {
				if st.Finished() {
					return ectx.Render(200, "network/:id/speedtest/:id", NewNetworkNodeSpeedtestView(n, st, nil))
				}

				if !n.Online {
					return ectx.Redirect(301, ectx.Request().URL.JoinPath("..").Path)
				}

				return ectx.Render(200, "network/:id/speedtest/:id", NewNetworkNodeSpeedtestView(n, st, nil))
			})
		})
	})
}

// GET /network
func networkWebsocketHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		ws, err := upgrader.Upgrade(ectx.Response(), ectx.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		s := sc.NodeManager.Stream()
		mcs := sc.MasterConfiguration
		for cn := range s {
			if ws.IsClosed() {
				break
			}

			nsp := mcs.Current().NodeStatsPolling
			view := NewNetworkView(cn, nsp.Duration())
			err = ws.WriteTemplate(ectx, "network/nodes", view)
			if err != nil {
				logging.LogError("failed to write template in ws %v, %v", ws, err)
				continue
			}
		}
		sc.NodeManager.Release(s)
		return nil
	})
}

// GET /network/:id/speedtest/:id
func networkIdSpeedtestWebsocketHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			return withSpeedtest(ectx, sc, func(st models.Speedtest) error {
				ws, err := upgrader.Upgrade(ectx.Response(), ectx.Request(), nil)
				if err != nil {
					return err
				}
				defer ws.Close()

				s, ok := sc.NodeSpeedtest.Updates(st.ID)
				if !ok {
					ws.WriteTemplate(ectx, "network/:id/speedtest/:id", NewNetworkNodeSpeedtestView(n, models.Speedtest{}, err))
					return ws.Close()
				}

				ph := st.Phase
				for st = range s {
					// todo: if client breaks, stream is hanging on commander side, need to notify it!
					if ws.IsClosed() {
						break
					}

					if st.Phase != ph {
						ph = st.Phase
						status, err := RenderString(ectx, "network/:id/speedtest/:id/status", SpeedtestPhaseView(ph))
						if err != nil {
							logging.LogError("failed to write template in ws %v, %v", ws, err)
							continue
						}

						err = ws.WriteJSON(NewSpeedtestStatusElementView(status))
						if err != nil {
							logging.LogError("failed to write template in ws %v, %v", ws, err)
							continue
						}
					}

					switch st.Phase {
					case models.SpeedtestLatency:
						err = ws.WriteJSON(NewSpeedtestLatencyElementView(st.Latency))
					case models.SpeedtestDownload:
						err = ws.WriteJSON(NewSpeedtestDownloadElementView(st.DownloadSpeed))
					case models.SpeedtestUpload:
						err = ws.WriteJSON(NewSpeedtestUploadElementView(st.UploadSpeed))
					default:
						err = nil
					}

					if err != nil {
						logging.LogError("failed to write template in ws %v, %v", ws, err)
						continue
					}
				}

				return nil
			})
		})
	})
}
