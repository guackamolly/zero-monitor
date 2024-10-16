package http

import (
	"fmt"
	"os"
	"strings"
)

const (
	serverVirtualHostEnvKey = "server_virtual_host"
)

var (
	serverVirtualHost = os.Getenv(serverVirtualHostEnvKey)
)

func WithVirtualHost(path string) string {
	if !strings.HasPrefix(serverVirtualHost, "/") {
		serverVirtualHost = "/" + serverVirtualHost
	}

	if !strings.HasSuffix(serverVirtualHost, "/") {
		serverVirtualHost = serverVirtualHost + "/"
	}

	if path == "" || path == "/" {
		return serverVirtualHost
	}

	if serverVirtualHost == "/" {
		return path
	}

	if path[0] == '/' {
		return fmt.Sprintf("%s%s", serverVirtualHost, path[1:])
	}

	return fmt.Sprintf("%s%s", serverVirtualHost, path)
}
