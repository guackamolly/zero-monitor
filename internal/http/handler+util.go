package http

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

// Common headers to lookup in [IsReverseProxyRequest].
var checkReverseProxyHeaders = []string{
	echo.HeaderXForwardedFor,
	echo.HeaderXForwardedProto,
	echo.HeaderXRealIP,
	"X-Forwarded-Host",
}

// Serves as a global bucket for storing errors that occur in handlers.
// Each error is associated to an UUID.
var handlerErrorsBucket = map[string]error{}

// Self redirects and appends an error uuid to the url.
func RedirectWithError(ectx echo.Context, err error) error {
	u := ectx.Request().URL
	q := u.Query()

	uuid := StoreHandlerError(err)
	q.Set(redirectErrorQueryParam, uuid)
	u.RawQuery = q.Encode()

	return ectx.Redirect(302, u.RequestURI())
}

// Tries to extract an error that may have originated from a self redirect.
func FromRedirectWithError(ectx echo.Context) (error, bool) {
	uuid := ectx.Request().URL.Query().Get(redirectErrorQueryParam)
	err, ok := handlerErrorsBucket[uuid]
	if ok {
		delete(handlerErrorsBucket, uuid)
	}

	return err, ok
}

func ExtractBreakpoint(ectx echo.Context) (Breakpoint, bool) {
	bp := ectx.Request().URL.Query().Get(breakpointQueryParam)
	if bp == "" {
		return 0, false
	}

	bpconv, err := strconv.ParseInt(bp, 10, 32)
	if err != nil {
		logging.LogWarning("failed to convert breakpoint %s to int, %v", bpconv, err)
		return 0, false
	}
	breakpoint := Breakpoint(bpconv)
	if breakpoint <= MobileBreakpoint {
		return MobileBreakpoint, true
	}

	if breakpoint <= TabletBreakpoint {
		return TabletBreakpoint, true
	}

	return DesktopBreakpoint, true
}

func StoreHandlerError(err error) string {
	uuid := models.UUID()
	handlerErrorsBucket[uuid] = err

	return uuid
}

func RenderString(ectx echo.Context, tpl string, v any) (string, error) {
	var buf strings.Builder

	err := ectx.Echo().Renderer.Render(&buf, tpl, v, ectx)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Composes an URL relative to the server host.
func URL(ectx echo.Context, path string, query map[string]string) *url.URL {
	host := ectx.Request().Host
	return RawURL(ectx, host, path, query)
}

// Composes a raw url. Host must be in the address format (<host>:<port>).
func RawURL(ectx echo.Context, host string, path string, query map[string]string) *url.URL {
	var qb strings.Builder

	c := 0
	for k, v := range query {
		if c > 0 {
			qb.WriteByte('&')
		}

		qb.WriteString(k)
		qb.WriteByte('=')
		qb.WriteString(v)

		c++
	}

	return &url.URL{
		Scheme:   ectx.Scheme(),
		Host:     host,
		Path:     path,
		RawQuery: qb.String(),
	}
}

func ServerAddress(ectx echo.Context) string {
	addr := extractServerAddr(ectx)
	return fmt.Sprintf("[%s]:%d", addr.IP, addr.Port)
}

func IsLocalRequest(ectx echo.Context) bool {
	ip := net.ParseIP(ectx.RealIP())

	return ip.IsLoopback() || ip.IsPrivate()
}

func IsReverseProxyRequest(ectx echo.Context) bool {
	h := ectx.Request().Header
	for _, c := range checkReverseProxyHeaders {
		if _, ok := h[c]; ok {
			return true
		}
	}

	return false
}

// Checks if the server is bind to an unspecified address (e.g., 0.0.0.0)
func IsBindToUnspecified(ectx echo.Context) bool {
	addr := extractServerAddr(ectx)
	return addr.IP.IsUnspecified()
}

func ExtractReverseProxyIP(ectx echo.Context) string {
	xff, ok := ectx.Request().Header[echo.HeaderXForwardedFor]
	if !ok || len(xff) < 2 {
		return echo.ExtractIPDirect()(ectx.Request())
	}

	return xff[1]
}

func ExtractHost(ectx echo.Context) string {
	addr, _, err := net.SplitHostPort(ectx.Request().Host)
	if err != nil {
		return ectx.Request().Host
	}

	return addr
}

func ExtractPort(ectx echo.Context) string {
	_, port, _ := net.SplitHostPort(ectx.Request().Host)
	return port
}

func GetLastVisitedPath(ectx echo.Context) (string, error) {
	cookie, err := ectx.Cookie(lastVisitedPageCookie)
	if err != nil {
		return "", err
	}

	if cookie == nil {
		return "", errors.New("cookie does not exist")
	}

	return cookie.Value, nil
}

func SetLastVisitedPathCookie(ectx echo.Context) {
	ectx.SetCookie(
		NewCookie(
			ectx,
			lastVisitedPageCookie,
			ectx.Request().URL.String(),
			WithVirtualHost(rootRoute),
			time.Now().Add(5*time.Minute),
		),
	)
}

func UnsetLastVisitedPathCookie(ectx echo.Context) {
	cookie, err := ectx.Cookie(lastVisitedPageCookie)
	if err != nil {
		logging.LogWarning("trying to unset cookie that does not exist")
		return
	}

	if cookie == nil {
		logging.LogWarning("trying to unset cookie that is nil")
		return
	}

	cookie.MaxAge = -1
	ectx.SetCookie(cookie)
}

func extractServerAddr(ectx echo.Context) *net.TCPAddr {
	var addr net.Addr

	if ectx.IsTLS() {
		addr = ectx.Echo().TLSListenerAddr()
	} else {
		addr = ectx.Echo().ListenerAddr()
	}

	if c, ok := addr.(*net.TCPAddr); ok {
		return c
	}

	panic("extractServerAddr: unexpected unreachable statement")
}
