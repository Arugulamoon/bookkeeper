package web

import (
	"errors"
	"html/template"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func (t *TemplateRenderer) Render(
	c *echo.Context, w io.Writer, name string, data any,
) error {
	ts, ok := t.templates[name]
	if !ok {
		return errors.New("template not found: " + name)
	}
	return ts.ExecuteTemplate(w, name, data)
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func add(a, b int) int {
	return a + b
}

var functions = template.FuncMap{
	"formatDate": formatDate,
	"add":        add,
	"join":       strings.Join,
	"split":      strings.Split,
}

func NewTemplateCache(
	dir string,
) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
