package http

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	files = map[string]string{
		WithVirtualHost("/index.css"):     "index.css",
		WithVirtualHost("/manifest.json"): "manifest.json",
	}

	dirs = map[string]string{
		WithVirtualHost("/static"): "static/",
	}

	templates = map[string]string{
		"homepage":                      "index.gohtml",
		"dashboard":                     "tpl/dashboard/*.gohtml",
		"network":                       "tpl/network/*.gohtml",
		"network/:id":                   "tpl/network/id/*.gohtml",
		"network/:id/actions":           "tpl/network/id/actions/*.gohtml",
		"network/:id/connections":       "tpl/network/id/connections/*.gohtml",
		"network/:id/packages":          "tpl/network/id/packages/*.gohtml",
		"network/:id/processes":         "tpl/network/id/processes/*.gohtml",
		"network/:id/speedtest":         "tpl/network/id/speedtest/*.gohtml",
		"network/:id/speedtest/history": "tpl/network/id/speedtest/history/*.gohtml",
		"network/:id/speedtest/:id":     "tpl/network/id/speedtest/id/*.gohtml",
		"settings":                      "tpl/settings/*.gohtml",
		"user":                          "tpl/user/*.gohtml",
		"user/new":                      "tpl/user/new/*.gohtml",
	}

	httpErrors = map[int]string{
		404: "404/index.html",
		500: "500/index.html",
	}

	fallback = httpErrors[500]

	modtime = time.Now()
)

func RegisterStaticFiles(e *echo.Echo, fs fs.FS) error {
	e.Filesystem = fs

	for k, v := range files {
		e.FileFS(k, v, fs, cacheStaticFileMiddleware(modtime))
	}

	for k, v := range dirs {
		fs = echo.MustSubFS(fs, v)

		// same thing as e.StaticFS, execept [staticFileMiddleware] is included.
		e.Add(
			http.MethodGet,
			k+"*",
			echo.StaticDirectoryHandler(fs, false),
			cacheStaticFileMiddleware(modtime),
		)
	}

	return nil
}
