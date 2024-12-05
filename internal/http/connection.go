package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Start(e *echo.Echo, host string, port string) error {
	return e.Start(fmt.Sprintf("[%s]:%s", host, port))
}

func StartTLS(e *echo.Echo, host string, port string, certfilepath string, keyfilepath string) error {
	return e.StartTLS(fmt.Sprintf("[%s]:%s", host, port), certfilepath, keyfilepath)
}
