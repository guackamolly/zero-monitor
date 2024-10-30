package http

import (
	"strconv"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

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

func networkIdHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			return ectx.Render(200, "network/:id", NewNetworkNodeInformationView(n))
		})
	})
}

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

func networkIdProcessesHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			if !n.Online {
				return ectx.Render(200, "network/:id/processes", NewNetworkNodeProcessesView(n, []models.Process{}, nil))
			}

			logging.LogInfo("fetching node processes")
			procs, err := sc.NodeCommander.Processes(n.ID)
			if err != nil {
				logging.LogError("failed to fetch node processes, %v", procs)
			}

			return ectx.Render(200, "network/:id/processes", NewNetworkNodeProcessesView(n, procs, err))
		})
	})
}

func networkIdProcessesFormHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		return withPathNode(ectx, sc, func(n models.Node) error {
			if !n.Online {
				return ectx.Redirect(301, ectx.Request().URL.Path)
			}

			pid, err := strconv.Atoi(ectx.FormValue("kill"))
			if err != nil {
				logging.LogError("failed to convert pid %s to int, %v", pid, err)
				// todo: handle failed conversion
			}

			logging.LogInfo("killing node process")
			err = sc.NodeCommander.KillProcess(n.ID, int32(pid))
			if err != nil {
				logging.LogError("failed to kill node process, %v", err)
				// todo: handle
			}

			return ectx.Redirect(301, ectx.Request().URL.Path)
		})
	})
}

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
