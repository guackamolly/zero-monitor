package http

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.GET(rootRoute, rootHandler)

	e.GET(dashboardRoute, dashboardHandler)
	e.POST(dashboardRoute, dashboardFormHandler)

	e.GET(networkRoute, networkHandler)
	e.GET(networkPublicKeyRoute, networkPublicKeyHandler)
	e.GET(networkConnectionEndpointRoute, networkConnectionEndpointHandler)

	e.GET(networkIdRoute, networkIdHandler)
	e.GET(networkIdConnectionsRoute, networkIdConnectionsHandler)
	e.GET(networkIdPackagesRoute, networkIdPackagesHandler)
	e.GET(networkIdProcessesRoute, networkIdProcessesHandler)
	e.POST(networkIdProcessesRoute, networkIdProcessesFormHandler)
	e.GET(networkIdSpeedtestRoute, networkIdSpeedtestHandler)
	e.POST(networkIdSpeedtestRoute, networkIdSpeedtestFormHandler)
	e.GET(networkIdSpeedtestHistoryRoute, networkIdSpeedtestHistoryHandler)
	e.GET(networkIdSpeedtestHistoryChartRoute, networkIdSpeedtestHistoryChartHandler)
	e.GET(networkIdSpeedtestIdRoute, networkIdSpeedtestIdHandler)

	e.GET(settingsRoute, getSettingsHandler)
	e.POST(settingsRoute, updateSettingsHandler)

	e.HTTPErrorHandler = httpErrorHandler()
}

func rootHandler(ectx echo.Context) error {
	return fmt.Errorf("not implemented yet")
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
			he = echo.NewHTTPError(500, err)
		}

		// if error page available, serve it
		if f, eok := httpErrors[he.Code]; eok {
			c.Response().WriteHeader(he.Code)
			c.File(f)
			return
		}

		// if no match, resort to fallback
		c.File(fallback)
	}
}
