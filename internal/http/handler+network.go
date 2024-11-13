package http

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
	"github.com/mssola/useragent"
)

// GET /network
func networkHandler(ectx echo.Context) error {
	if upgrader.WantsToUpgrade(*ectx.Request()) {
		return networkWebsocketHandler(ectx)
	}

	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		view := NewNetworkView(
			sc.NodeManager.Network(),
			sc.MasterConfiguration.Current().NodeStatsPolling.Duration(),
		)

		return ectx.Render(200, "network", view)
	})
}

// GET /network/:id
func networkIdHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			return ectx.Render(200, "network/:id", NewNetworkNodeInformationView(n))
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

			logging.LogInfo("fetching node connections")
			conns, err := sc.NodeCommander.Connections(n.ID)
			if err != nil {
				logging.LogError("failed to fetch node connections, %v", conns)
			}

			return ectx.Render(200, "network/:id/connections", NewNetworkNodeConnectionsView(n, conns, err))
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

			logging.LogInfo("fetching node processes")
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

			logging.LogInfo("killing node process")
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

			logging.LogInfo("starting speedtest")
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
			return withUserAgent(ectx, func(ua *useragent.UserAgent) error {
				sts, ok := sc.NodeSpeedtest.History(n.ID)
				if !ok {
					logging.LogWarning("no history found for node %s", n.ID)
					sts = []models.Speedtest{}
				}
				return ectx.Render(200, "network/:id/speedtest/history", NewNetworkNodeSpeedtestHistoryView(n, sts, ua.Mobile(), nil))
			})
		})
	})
}

// GET /network/:id/speedtest/history/chart
func networkIdSpeedtestHistoryChartHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			return withUserAgent(ectx, func(ua *useragent.UserAgent) error {
				bp, ok := ExtractBreakpoint(ectx)
				if !ok {
					bp = DesktopBreakpoint
				}

				sts, ok := sc.NodeSpeedtest.History(n.ID)
				if !ok {
					logging.LogWarning("no history found for node %s", n.ID)
					sts = []models.Speedtest{}
				}
				return ectx.Render(200, "network/:id/speedtest/history/chart", NewSpeedtestHistoryChartView(sts, bp))
			})
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
