package http

import (
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo"
)

func RegisterHandlers(e *echo.Echo) {
	e.GET(rootRoute, anyRouteHandler)
	e.HTTPErrorHandler = httpErrorHandler()
}
func anyRouteHandler(ectx echo.Context) error {
	return ectx.File(root)
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
