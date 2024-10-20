package http

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.GET(rootRoute, rootHandler)
	e.GET(networkRoute, networkHandler)
	e.GET(settingsRoute, getSettingsHandler)
	e.POST(settingsRoute, updateSettingsHandler)

	e.HTTPErrorHandler = httpErrorHandler()
}

func rootHandler(ectx echo.Context) error {
	return fmt.Errorf("not implemented yet")
}

func networkHandler(ectx echo.Context) error {
	if upgrader.WantsToUpgrade(*ectx.Request()) {
		return networkWebsocketHandler(ectx)
	}

	return withSubscriberContainer(ectx, func(sc *di.SubscribeContainer) error {
		view := NewServerStatsView(sc.NodeManager.Network())

		return ectx.Render(200, "network", view)
	})
}

func networkWebsocketHandler(ectx echo.Context) error {
	return withSubscriberContainer(ectx, func(sc *di.SubscribeContainer) error {
		ws, err := upgrader.Upgrade(ectx.Response(), ectx.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		s := sc.NodeManager.Stream()
		for cn := range s {
			if ws.IsClosed() {
				break
			}

			view := NewServerStatsView(cn)
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

func httpErrorHandler() func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		// make sure to not process any false positives
		if err == nil {
			return
		}

		logging.LogError("handling error... %v", err)
		he, ok := err.(*echo.HTTPError)

		// If all cast fail, serve fallback
		if !ok {
			c.File(fallback)
			return
		}

		// if error page available, serve it
		if f, eok := errors[he.Code]; eok {
			c.File(f)
			return
		}

		// if no match, resort to fallback
		c.File(fallback)
	}
}
