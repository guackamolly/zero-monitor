package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	tokenCookie = "token"
)

func NewCookie(
	ectx echo.Context,
	name string,
	value string,
	path string,
	expiry time.Time,
) *http.Cookie {
	c := new(http.Cookie)
	c.Name = name
	c.Value = value
	c.Path = path
	c.Expires = expiry
	c.SameSite = http.SameSiteStrictMode
	c.Secure = ectx.IsTLS()

	return c
}
