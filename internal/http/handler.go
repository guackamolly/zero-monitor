package http

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.GET(rootRoute, rootHandler)
	e.GET(networkRoute, networkHandler)
	e.GET(networkIdRoute, networkIdHandler)
	e.GET(networkIdConnectionsRoute, networkIdConnectionsHandler)
	e.GET(networkIdProcessesRoute, networkIdProcessesHandler)
	e.POST(networkIdProcessesRoute, networkIdProcessesFormHandler)
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
