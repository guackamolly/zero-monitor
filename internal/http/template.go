package http

import (
	"fmt"
	"html/template"
	"io"

	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("failed to find template for name: %s", name)
	}

	return tpl.Execute(w, data)
}

func RegisterTemplates(e *echo.Echo) error {
	tps := map[string]*template.Template{}
	for k, v := range templates {
		t, err := template.New(k).ParseFiles(v)
		if err != nil {
			logging.LogFatal("failed to parse template file, %v", err)
		}

		tps[k] = t
	}

	t := &Template{
		templates: tps,
	}

	e.Renderer = t

	return nil
}
