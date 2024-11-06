package http

import (
	"github.com/labstack/echo/v4"
)

const (
	serverPublicRootEnvKey = "server_public_root"
)

var (
	serverPublicRoot = GetEnv(serverPublicRootEnvKey, func() string {
		return "public/"
	})

	files = map[string]string{
		WithVirtualHost("/index.html"):    serverPublicRoot + "index.html",
		WithVirtualHost("/index.css"):     serverPublicRoot + "index.css",
		WithVirtualHost("/manifest.json"): serverPublicRoot + "manifest.json",
	}

	dirs = map[string]string{
		WithVirtualHost("/static"):  serverPublicRoot + "static/",
		WithVirtualHost("/content"): serverPublicRoot + "content/",
		WithVirtualHost("/about"):   serverPublicRoot + "about/",
		WithVirtualHost("/contact"): serverPublicRoot + "contact/",
	}

	templates = map[string]string{
		"network":                       serverPublicRoot + "tpl/network/*.gohtml",
		"network/:id":                   serverPublicRoot + "tpl/network/:id/*.gohtml",
		"network/:id/connections":       serverPublicRoot + "tpl/network/:id/connections/*.gohtml",
		"network/:id/processes":         serverPublicRoot + "tpl/network/:id/processes/*.gohtml",
		"network/:id/speedtest":         serverPublicRoot + "tpl/network/:id/speedtest/*.gohtml",
		"network/:id/speedtest/history": serverPublicRoot + "tpl/network/:id/speedtest/history/*.gohtml",
		"network/:id/speedtest/:id":     serverPublicRoot + "tpl/network/:id/speedtest/:id/*.gohtml",
		"settings":                      serverPublicRoot + "tpl/settings/*.gohtml",
	}

	httpErrors = map[int]string{
		404: serverPublicRoot + "404/index.html",
		500: serverPublicRoot + "500/index.html",
	}

	root     = files[WithVirtualHost("/index.html")]
	fallback = root
)

func RegisterStaticFiles(e *echo.Echo) error {
	for k, v := range files {
		e.File(k, v)
	}

	for k, v := range dirs {
		e.Static(k, v)
	}

	return nil
}

// Returns a tuples that identifies:
// [0] - Content directory that is accessible through the file system (physical)
// [1] - Content directory that is accessible through the network (virtual)
func ContentDir() [2]string {
	return [2]string{
		dirs[WithVirtualHost("/content")],
		WithVirtualHost("/content/"),
	}
}
