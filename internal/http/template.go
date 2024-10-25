package http

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

// Global utility template functions.
var funcMap = template.FuncMap{
	"sequence": func(count int) []int {
		s := make([]int, count)
		for i := 0; i < int(count); i++ {
			s[i] = i
		}

		return s
	},
}

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var part string

	tpl, ok := t.templates[name]
	if !ok {
		s := strings.Split(name, "/")
		part = s[len(s)-1]
		tpl, ok = t.templates[s[0]]
	}

	if !ok {
		return fmt.Errorf("failed to find template for name: %s", name)
	}

	if len(part) == 0 {
		return tpl.Execute(w, data)
	}

	return tpl.ExecuteTemplate(w, part, data)
}

func RegisterTemplates(e *echo.Echo) error {
	tps := map[string]*template.Template{}
	for k, v := range templates {
		t, err := template.New(k).Funcs(funcMap).ParseGlob(v)
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
