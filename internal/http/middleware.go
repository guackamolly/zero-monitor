package http

import (
	"context"
	"net/http"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
	"github.com/mssola/useragent"
)

const ctxKey = "ctx.key"
const isAdminKey = "is.admin.key"

const staticCacheControlHeader = "public, max-age=604800"

func RegisterMiddlewares(e *echo.Echo, ctx context.Context) {
	e.Use(loggingMiddleware())
	e.Use(contextMiddleware(ctx))
	e.Use(cookieMiddleware())
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

// Middleware that adds cache-control and last-modified headers.
// only apply it for static files. NEVER apply this middleware to template routes.
func cacheStaticFileMiddleware(modtime time.Time) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			h := ectx.Response().Header()
			h.Set(echo.HeaderCacheControl, staticCacheControlHeader)
			h.Set(echo.HeaderLastModified, modtime.Format(time.RFC1123))

			return next(ectx)
		}
	}
}

// todo: add integration tests to make sure isAdminKey is NEVER SET if token is not from admin user (probabbly need to fuzzy test token)
func cookieMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			var cookie *http.Cookie
			var err error

			if cookie, err = ectx.Cookie(tokenCookie); err != nil || cookie == nil {
				return next(ectx)
			}

			return withServiceContainer(ectx, func(sc *ServiceContainer) error {
				ectx.Set(isAdminKey, sc.Authorization.HasAdminRights(cookie.Value))
				return next(ectx)
			})
		}
	}
}

// Use this middleware to guard routes that can only be accessed by admin users.
// todo: add integration tests to make sure only admin token cookies can access protected pages.
var adminRouteMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		return withServiceContainer(ectx, func(sc *ServiceContainer) error {
			var cookie *http.Cookie
			var err error

			if cookie, err = ectx.Cookie(tokenCookie); err != nil || cookie == nil {
				return ectx.Redirect(302, userRoute)
			}

			if !sc.Authorization.HasAdminRights(cookie.Value) {
				return ectx.Redirect(302, userRoute)
			}

			return next(ectx)
		})
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

func withJoinCode(ectx echo.Context, sc *ServiceContainer, with func(code string) error) error {
	c := ectx.QueryParam(joinQueryParam)
	if !sc.NodeManager.Valid(c) {
		return echo.ErrUnauthorized
	}

	return with(c)
}

func extractQuery(ectx echo.Context, param string) (string, bool) {
	if !ectx.QueryParams().Has(param) {
		return "", false
	}

	return ectx.QueryParam(param), true
}

func extractUserAgent(ectx echo.Context) UserAgent {
	return UserAgent(*useragent.New(ectx.Request().UserAgent()))
}

func hasAdminRights(ectx echo.Context) bool {
	return ectx.Get(isAdminKey) == true
}
