package http

import (
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
		code := sc.NodeManager.Code()
		url := URL(ectx, networkRoute, map[string]string{joinQueryParam: code.Code})
		dashboardView = dashboardView.WithInviteLink(NewDashNetworkInviteLinkView(url.String(), code))

		return ectx.Redirect(301, ectx.Request().URL.Path)
	})
}
