package controller

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/helper"
	"github.com/iagapie/go-spring/pkg/manager"
	"github.com/iagapie/go-spring/pkg/router"
	"github.com/iagapie/go-spring/pkg/spring"
	"github.com/iagapie/go-spring/pkg/theme"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
		compManager manager.ComponentManager
		log         *logrus.Entry
		s           *spring.Spring
		t           theme.Theme
		r           router.Router
		funcs       theme.FuncMap
		c           echo.Context
		p           theme.Page
		l           theme.Layout
		vars        map[string]interface{}
		current     theme.View
	}
)

func New(s *spring.Spring, t theme.Theme, compManager manager.ComponentManager, log *logrus.Entry) Controller {
	ct := &ctr{
		compManager: compManager,
		log:         log,
		s:           s,
		t:           t,
		r:           router.New(),
		funcs:       make(theme.FuncMap),
	}
	ct.init()
	return ct
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

	ct.funcs["page"] = func() template.HTML {
		page, err := ct.render(ct.p, nil)
		if err != nil {
			ct.log.Warnf("page render: %v", err)
			return ""
		}
		return template.HTML(page)
	}

	render := func(v theme.View, vars map[string]interface{}) template.HTML {
		if !v.Parsed() {
			v.AddFuncs(ct.funcs)
			if err := v.Load(); err != nil {
				ct.log.Warnf("view load: %v", err)
				return ""
			}
		}
		content, err := ct.render(v, vars)
		if err != nil {
			ct.log.Warnf("view render: %v", err)
			return ""
		}
		return template.HTML(content)
	}

	ct.funcs["partial"] = func(name string, data ...interface{}) template.HTML {
		if p := ct.t.Partial(name); p != nil {
			return render(p, map[string]interface{}{"data": data})
		}
		ct.log.Warnf("partial %s not found", name)
		return ""
	}

	ct.funcs["component"] = func(code string) template.HTML {
		if cfgComp, ok := ct.current.CfgComps()[code]; ok {
			if ct.compManager.Has(cfgComp.Name) {
				compFn := ct.compManager.Resolve(cfgComp.Name)
				comp, err := compFn(manager.ComponentProps(cfgComp.Props))
				if err != nil {
					ct.log.Warnf("component %s init error: %v", code, err)
					return ""
				}
				file := comp.Details().ViewFile
				var v theme.View
				if v = ct.t.Partial(fmt.Sprintf("%s/%s", code, file)); v == nil {
					file = fmt.Sprintf("%s/components/%s/%s.html", ct.compManager.FindPlugin(comp).Dir(), code, file)
					if helper.FileExists(file) {
						v = theme.NewView(theme.WithViewFile(file))
					}
				}
				if v == nil {
					ct.log.Warnf("component %s view not found", code)
					return ""
				}
				return render(v, map[string]interface{}{"prop": comp.Props()})
			}
		}
		ct.log.Warnf("component %s not found", code)
		return ""
	}
}

func (ct *ctr) Run(c echo.Context) error {
	ct.c = c

	if err := ct.reset(true); err != nil {
		return err
	}

	if !ct.findPage(c.Request().RequestURI) {
		return echo.ErrNotFound
	}

	if err := ct.findLayout(); err != nil {
		return err
	}

	layout, err := ct.render(ct.l, nil)
	if err != nil {
		return fmt.Errorf("layout render: %w", err)
	}

	return c.HTML(http.StatusOK, layout)
}

func (ct *ctr) render(v theme.View, vars map[string]interface{}) (string, error) {
	oldView := ct.current
	defer func() {
		ct.current = oldView
	}()
	ct.current = v
	ct.initVars()

	for k, val := range vars {
		ct.vars[k] = val
	}

	return v.Render(ct.vars)
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
	}

	ct.p = nil
	ct.l = nil
	ct.vars = nil

	return nil
}

func (ct *ctr) findPage(url string) bool {
	if ct.r.Match(url) {
		ct.p = ct.t.Pages()[ct.r.Matched()]

		if ct.p.IsHidden() {
			ct.p = nil
		}
	}

	if ct.p == nil && ct.r.Match("/404") {
		ct.p = ct.t.Pages()[ct.r.Matched()]
	}

	return ct.p != nil
}

func (ct *ctr) findLayout() error {
	if layout, ok := ct.t.Layouts()[ct.p.Layout()]; ok {
		if !layout.Parsed() {
			layout.AddFuncs(ct.funcs)
			if err := layout.Load(); err != nil {
				return err
			}
		}
		ct.l = layout
		return nil
	}

	err := echo.ErrInternalServerError
	err.Internal = fmt.Errorf("layout %s not found", ct.p.Layout())
	return err
}

func (ct *ctr) initVars() {
	ct.vars = make(map[string]interface{})
	ct.vars["request"] = ct.c.Request()
	ct.vars["param"] = ct.r.Parameters()
	ct.vars["layout"] = ct.l.CfgProps()
	ct.vars["comps"] = ct.current.CfgComps()

	props := ct.p.CfgProps()
	props["name"] = ct.p.Name()
	props["file"] = ct.p.File()
	ct.vars["page"] = props

}
