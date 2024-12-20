package http

import (
	"github.com/labstack/echo/v4"
)

// GET /user
func userHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		if sc.Authentication.NeedsAdminRegistration() {
			return ectx.Redirect(302, userNewRoute)
		}

		return ectx.Render(200, "user", nil)
	})
}

// POST /user
func userFormHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		username := ectx.FormValue("username")
		password := ectx.FormValue("password")

		t, err := sc.Authentication.Authenticate(username, password)
		if err != nil {
			return ectx.Redirect(302, userRoute)
		}

		ectx.SetCookie(NewCookie(ectx, tokenCookie, t.Value, WithVirtualHost(rootRoute), t.Expiry))
		if p, err := GetLastVisitedPath(ectx); err == nil {
			UnsetLastVisitedPathCookie(ectx)
			return ectx.Redirect(302, p)
		}

		return ectx.Redirect(302, dashboardRoute)
	})
}

// GET /user/new
func userNewHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		if !sc.Authentication.NeedsAdminRegistration() {
			return ectx.Redirect(302, rootRoute)
		}

		return ectx.Render(200, "user/new", nil)
	})
}

// POST /user/new
func userNewFormHandler(ectx echo.Context) error {
	return withServiceContainer(ectx, func(sc *ServiceContainer) error {
		if !sc.Authentication.NeedsAdminRegistration() {
			return ectx.Redirect(302, rootRoute)
		}

		username := ectx.FormValue("username")
		password := ectx.FormValue("password")

		t, err := sc.Authentication.RegisterAdmin(username, password)
		if err != nil {
			return ectx.Redirect(302, userNewRoute)
		}

		ectx.SetCookie(NewCookie(ectx, tokenCookie, t.Value, WithVirtualHost(rootRoute), t.Expiry))
		if p, err := GetLastVisitedPath(ectx); err == nil {
			UnsetLastVisitedPathCookie(ectx)
			return ectx.Redirect(302, p)
		}

		return ectx.Redirect(302, dashboardRoute)
	})
}
