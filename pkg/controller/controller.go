package controller

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/router"
	"github.com/iagapie/go-spring/pkg/spring"
	"github.com/iagapie/go-spring/pkg/theme"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type (
	Controller interface {
		Run(c echo.Context) error
	}

	ctr struct {
		s            *spring.Spring
		t            theme.Theme
		r            router.Router
		funcs        theme.FuncMap
		p            theme.Page
		l            theme.Layout
		vars         map[string]interface{}
		pageContents string
	}
)

func New(s *spring.Spring, t theme.Theme) Controller {
	ct := &ctr{
		s:     s,
		t:     t,
		r:     router.New(),
		funcs: make(theme.FuncMap),
	}
	ct.init()
	return ct
}

func (ct *ctr) Run(c echo.Context) error {
	if err := ct.reset(true); err != nil {
		return err
	}

	ct.findPage(c.Request().RequestURI)
	if ct.p == nil {
		return echo.ErrNotFound
	}

	ct.findLayout()
	if ct.l == nil {
		err := echo.ErrInternalServerError
		err.Internal = fmt.Errorf("layout %s not found", ct.p.Layout())
		return err
	}

	ct.initVars(c)

	page, err := ct.p.Render(ct.vars)
	if err != nil {
		return fmt.Errorf("page render: %w", err)
	}
	ct.pageContents = page

	layout, err := ct.l.Render(ct.vars)
	if err != nil {
		return fmt.Errorf("layout render: %w", err)
	}

	return c.HTML(http.StatusOK, layout)
}

func (ct *ctr) init() {
	ct.funcs["now"] = time.Now
	ct.funcs["time"] = func(t time.Time, layout string) string {
		return t.Format(layout)
	}
	ct.funcs["is_page"] = func(name string) bool {
		return strings.EqualFold(ct.p.Name(), name)
	}
	ct.funcs["page_url"] = func(name string, params ...interface{}) string {
		if url := ct.r.URI(name, params...); url != "" {
			return url
		}
		return ct.r.URIFromPattern(name, params...)
	}
	ct.funcs["route_url"] = ct.s.Reverse
	ct.funcs["assets"] = func(name string) string {
		uri, _ := ct.t.Assets()
		return fmt.Sprintf("%s/%s", uri, name)
	}
	ct.funcs["partial"] = func(name string, data ...interface{}) template.HTML {
		var vars interface{}
		if len(data) > 0 {
			vars = data[0]
		}
		if p := ct.t.Partial(name); p != nil {
			if !p.Parsed() {
				p.AddFuncs(ct.funcs)
				if err := p.Load(); err != nil {
					return ""
				}
			}
			if partial, err := p.Render(vars); err == nil {
				return template.HTML(partial)
			}
		}
		return ""
	}
	ct.funcs["page"] = func() template.HTML {
		return template.HTML(ct.pageContents)
	}
}

func (ct *ctr) reset(load bool) error {
	if load {
		ct.t.Reset()
		ct.r.Reset()

		for _, p := range ct.t.Pages() {
			if !p.Parsed() {
				p.AddFuncs(ct.funcs)
				if err := p.Load(); err != nil {
					return err
				}
			}
			ct.r.Route(p.Name(), p.URL())
		}

		ct.r.Sort()

		for _, l := range ct.t.Layouts() {
			if !l.Parsed() {
				l.AddFuncs(ct.funcs)
				if err := l.Load(); err != nil {
					return err
				}
			}
		}
	}

	ct.p = nil
	ct.l = nil
	ct.vars = nil
	ct.pageContents = ""

	return nil
}

func (ct *ctr) findPage(url string) {
	if ct.r.Match(url) {
		ct.p = ct.t.Pages()[ct.r.Matched()]

		if ct.p.IsHidden() {
			ct.p = nil
		}
	}

	if ct.p == nil && ct.r.Match("/404") {
		ct.p = ct.t.Pages()[ct.r.Matched()]
	}
}

func (ct *ctr) findLayout() {
	if layout, ok := ct.t.Layouts()[ct.p.Layout()]; ok {
		ct.l = layout
	}
}

func (ct *ctr) initVars(c echo.Context) {
	ct.vars = make(map[string]interface{})
	ct.vars["request"] = c.Request()
	ct.vars["param"] = ct.r.Parameters()
	ct.vars["layout"] = ct.l.CfgProps()

	props := ct.p.CfgProps()
	props["name"] = ct.p.Name()
	props["file"] = ct.p.File()
	ct.vars["page"] = props
}
