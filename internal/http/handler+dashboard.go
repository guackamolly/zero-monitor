package http

import (
	"net"

	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

// Holds the current view of the dashboard page.
var dashboardView = NewDashboardView()

// GET /dashboard
func dashboardHandler(ectx echo.Context) error {
	return ectx.Render(200, "dashboard", dashboardView)
}

// POST /dashboard
func dashboardFormHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		if dashboardView.ShowInviteLink() {
			return ectx.Redirect(301, ectx.Request().URL.Path)
		}

		host := ServerAddress(ectx)
		if IsReverseProxyRequest(ectx) {
			host = net.JoinHostPort(ExtractReverseProxyIP(ectx), ExtractPort(ectx))
		} else if IsBindToUnspecified(ectx) {
			ip, err := sc.Networking.PrivateIP()
			if err != nil {
				logging.LogError("server is bind on unspecified address, couldn't get a private interface IP to build url")
				return echo.ErrInternalServerError
			}

			_, port, _ := net.SplitHostPort(host)
			host = net.JoinHostPort(ip.String(), port)
		}

		code := sc.NodeManager.Code()
		url := RawURL(ectx, host, networkRoute, map[string]string{joinQueryParam: code.Code})
		dashboardView = dashboardView.WithInviteLink(NewDashNetworkInviteLinkView(url.String(), code))

		return ectx.Redirect(301, ectx.Request().URL.Path)
	})
}
