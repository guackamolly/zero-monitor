package http

import (
	"context"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
	"github.com/mssola/useragent"
)

const ctxKey = "ctx.key"

func RegisterMiddlewares(e *echo.Echo, ctx context.Context) {
	e.Use(loggingMiddleware())
	e.Use(contextMiddleware(ctx))
}

func loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			req := ectx.Request()

			logging.LogInfo("Host: %s | Method: %s | Path: %s | Client IP: %s", req.Host, req.Method, req.URL.RequestURI(), ectx.RealIP())
			return next(ectx)
		}
	}
}

func contextMiddleware(ctx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			ectx.Set(ctxKey, ctx)

			return next(ectx)
		}
	}
}

func withServiceContainer(ectx echo.Context, with func(*ServiceContainer) error) error {
	ctx, ok := ectx.Get(ctxKey).(context.Context)

	if ok {
		return with(ExtractServiceContainer(ctx))
	}

	return echo.ErrFailedDependency
}

func withPathNode(ectx echo.Context, sc *ServiceContainer, with func(models.Node) error) error {
	id := ectx.Param("id")
	n, ok := sc.NodeManager.Node(id)
	if ok {
		return with(n)
	}

	return echo.ErrNotFound
}

func withSpeedtest(ectx echo.Context, sc *ServiceContainer, with func(models.Speedtest) error) error {
	id := ectx.Param("id2")
	st, ok := sc.NodeSpeedtest.Speedtest(id)
	if ok {
		return with(st)
	}

	return echo.ErrNotFound
}

func extractUserAgent(ectx echo.Context) UserAgent {
	return UserAgent(*useragent.New(ectx.Request().UserAgent()))
}
