package http

import (
	"context"

	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
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

func withSubscriberContainer(ectx echo.Context, with func(*di.SubscribeContainer) error) error {
	ctx, ok := ectx.Get(ctxKey).(context.Context)

	if ok {
		return with(di.ExtractSubscribeContainer(ctx))
	}

	return echo.ErrFailedDependency
}
