package http

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)

	if err != nil {
		c.Logger().Error(err)
	}

	return err
}

func RegisterTemplates(e *echo.Echo) error {
	t := &Template{
		templates: template.Must(template.ParseFiles(templates...)),
	}

	e.Renderer = t

	return nil
}
