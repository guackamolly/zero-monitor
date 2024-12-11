package http

import (
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	// / (public)
	e.GET(rootRoute, rootHandler)

	// /dashboard (admin only)
	e.GET(dashboardRoute, dashboardHandler, adminRouteMiddleware)
	e.POST(dashboardRoute, dashboardFormHandler, adminRouteMiddleware)

	// /network (public)
	e.GET(networkRoute, networkHandler)
	e.GET(networkPublicKeyRoute, networkPublicKeyHandler)
	e.GET(networkConnectionEndpointRoute, networkConnectionEndpointHandler)

	// /network/:id/connections | packages (public)
	e.GET(networkIdRoute, networkIdHandler)
	e.GET(networkIdConnectionsRoute, networkIdConnectionsHandler)
	e.GET(networkIdPackagesRoute, networkIdPackagesHandler)

	// /network/:id/processes (admin only)
	e.GET(networkIdProcessesRoute, networkIdProcessesHandler, adminRouteMiddleware)
	e.POST(networkIdProcessesRoute, networkIdProcessesFormHandler, adminRouteMiddleware)

	// /network/:id/speedtest (POST admin only)
	e.GET(networkIdSpeedtestRoute, networkIdSpeedtestHandler)
	e.POST(networkIdSpeedtestRoute, networkIdSpeedtestFormHandler, adminRouteMiddleware)

	// /network/:id/processes (public)
	e.GET(networkIdSpeedtestHistoryRoute, networkIdSpeedtestHistoryHandler)
	e.GET(networkIdSpeedtestHistoryChartRoute, networkIdSpeedtestHistoryChartHandler)
	e.GET(networkIdSpeedtestIdRoute, networkIdSpeedtestIdHandler)

	// /settings (admin only)
	e.GET(settingsRoute, getSettingsHandler, adminRouteMiddleware)
	e.POST(settingsRoute, updateSettingsHandler, adminRouteMiddleware)

	// user (public if no admin has been registered yet)
	e.GET(userRoute, userHandler)
	e.POST(userRoute, userFormHandler)
	e.GET(userNewRoute, userNewHandler)
	e.POST(userNewRoute, userNewFormHandler)

	e.HTTPErrorHandler = httpErrorHandler()
}

func rootHandler(ectx echo.Context) error {
	v := dashboardView.WithContext(ectx)
	return ectx.Render(200, "homepage", v)
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
		c.Response().WriteHeader(he.Code)
		c.File(fallback)
	}
}
