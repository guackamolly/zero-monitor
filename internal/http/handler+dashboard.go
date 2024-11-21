package http

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
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
	url := URL(ectx, rootRoute, map[string]string{joinQueryParam: models.UUID()})
	expiry := time.Now().Add(5 * time.Minute)

	dashboardView = dashboardView.WithInviteLink(NewDashNetworkInviteLinkView(url.String(), expiry))
	return ectx.Redirect(301, ectx.Request().URL.Path)
}
