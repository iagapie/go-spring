package template

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

var _ echo.Renderer = &templateRenderer{}
var _ Renderer = &templateRenderer{}

type (
	Renderer interface {
		Render(io.Writer, string, interface{}, echo.Context) error
		Write(io.Writer, string, interface{}) error
		ParseAllFiles() error
	}

	Option interface {
		apply(t *templateRenderer)
	}

	option func(t *templateRenderer)

	templateRenderer struct {
		*template.Template
		name      string
		viewsPath string
		suffix    string
		reload    bool
		global    map[string]interface{}
		funcs     template.FuncMap
	}

	tmplData struct {
		Global map[string]interface{}
		Local  interface{}
	}
)

func (fn option) apply(t *templateRenderer) {
	fn(t)
}

func WithViewSuffix(suffix string) Option {
	return option(func(t *templateRenderer) {
		t.suffix = suffix
	})
}

func WithReload(reload bool) Option {
	return option(func(t *templateRenderer) {
		t.reload = reload
	})
}

func WithGlobal(name string, data interface{}) Option {
	return option(func(t *templateRenderer) {
		t.global[name] = data
	})
}

func WithFunc(name string, f interface{}) Option {
	return option(func(t *templateRenderer) {
		t.funcs[name] = f
	})
}

func New(name string, path string, opts ...Option) Renderer {
	t := &templateRenderer{
		Template:  nil,
		name:      name,
		viewsPath: path,
		suffix:    ".html",
		reload:    true,
		global:    map[string]interface{}{},
		funcs:     template.FuncMap{},
	}

	for _, opt := range opts {
		opt.apply(t)
	}

	return t
}

func (t *templateRenderer) ParseAllFiles() error {
	t.Template = template.New(t.name).Funcs(t.funcs)

	return filepath.Walk(t.viewsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, t.suffix) {
			if _, err = t.ParseFiles(path); err != nil {
				return err
			}
		}

		return nil
	})
}

func (t *templateRenderer) Write(w io.Writer, name string, data interface{}) error {
	if t.Template == nil || t.reload {
		if err := t.ParseAllFiles(); err != nil {
			return err
		}
	}

	return t.ExecuteTemplate(w, name, tmplData{
		Global: t.global,
		Local:  data,
	})
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.Write(w, name, data)
}
