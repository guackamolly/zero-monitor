package http

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

const (
	serverHostEnvKey         = "server_host"
	serverPortEnvKey         = "server_port"
	serverTLSCertFilePathKey = "server_tls_crt_fp"
	serverTLSKeyFilePathKey  = "server_tls_key_fp"
)

var (
	serverHost            = os.Getenv(serverHostEnvKey)
	serverPort            = os.Getenv(serverPortEnvKey)
	serverTLSCertFilePath = os.Getenv(serverTLSCertFilePathKey)
	serverTLSKeyFilePath  = os.Getenv(serverTLSKeyFilePathKey)
)

func Start(e *echo.Echo) error {
	addr := ServerAddress()

	if len(serverTLSCertFilePath) == 0 {
		return e.Start(addr)
	}

	return e.StartTLS(addr, serverTLSCertFilePath, serverTLSKeyFilePath)
}

func ServerAddress() string {
	return fmt.Sprintf("%s:%s", serverHost, serverPort)
}
