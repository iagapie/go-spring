package controller

import (
	"encoding/json"
	"fmt"
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/cms/router"
	"github.com/iagapie/go-spring/modules/cms/theme"
	"github.com/iagapie/go-spring/modules/sys/helper"
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/view"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type (
	Controller interface {
		Error(err error, c echo.Context)
		Run(c echo.Context) error
		RunPage(c echo.Context, page theme.View, useAjax bool) (string, error)
	}

	controller struct {
		s            *spring.Spring
		t            theme.Theme
		compManager  *component.Manager
		compParamRe  *regexp.Regexp
		handlerRe    *regexp.Regexp
		partialRe    *regexp.Regexp
		router       *router.Router
		partialStack *component.PartialStack
		cur          *current
		pageContents string
		componentCtx component.Component
		out          func(code int, b []byte) error
	}
)

const (
	HeaderRequestHandler  = "X_SPRING_REQUEST_HANDLER"
	HeaderRequestPartials = "X_SPRING_REQUEST_PARTIALS"

	ErrHTML = "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Spring CMS - Error</title></head><body><div class=\"container\"><h1>Error</h1><p>We're sorry, but something went wrong and the page cannot be displayed.</p></div></body></html>"
)

func New(s *spring.Spring, compManager *component.Manager) Controller {
	t := theme.ActiveTheme()
	r := router.NewRouter(t)
	stack := component.NewPartialStack()

	ctr := &controller{
		s:            s,
		t:            t,
		router:       r,
		partialStack: stack,
		compManager:  compManager,
		compParamRe:  regexp.MustCompile("^{{([^}]+)}}$"),
		handlerRe:    regexp.MustCompile("^(?:\\w+\\:{2})?On[A-Z][\\w+]*$"),
		partialRe:    regexp.MustCompile("^(?:\\w+\\:{2})?[\\w\\_\\-\\.\\/]+$"),
	}

	t.Funcs(funcs(ctr))

	return ctr
}

func (ctr *controller) Error(err error, c echo.Context) {
	if helper.Ajax(c.Request()) {
		c.Echo().DefaultHTTPErrorHandler(err, c)
		return
	}

	for _, header := range c.Request().Header.Values(echo.HeaderAccept) {
		if strings.HasPrefix(header, echo.MIMEApplicationJSON) {
			c.Echo().DefaultHTTPErrorHandler(err, c)
			return
		}
	}

	if c.Response().Committed {
		return
	}

	status := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
		status = he.Code
	}

	page := ctr.router.FindByURL(fmt.Sprintf("/%d", status))
	if page == nil {
		status = http.StatusInternalServerError
		page = ctr.router.FindByURL("/error")
	}

	var result string

	if page != nil {
		result, err = ctr.RunPage(c, page, true)
		if err != nil {
			ctr.s.Logger.Error(err)
		}
	}

	if len(result) == 0 {
		result = ErrHTML
	}

	if err = c.HTML(status, result); err != nil {
		ctr.s.Logger.Error(err)
	}
}

func (ctr *controller) Run(c echo.Context) error {
	if ctr.s.Cfg.App.Debug {
		ctr.router.Reset()
	}

	page := ctr.router.FindByURL(c.Request().RequestURI)
	if page == nil || page.Prop("is_hidden") == "1" {
		return echo.ErrNotFound
	}

	ctr.out = c.HTMLBlob
	result, err := ctr.RunPage(c, page, true)
	if err != nil {
		return err
	}

	return ctr.out(http.StatusOK, []byte(result))
}

func (ctr *controller) RunPage(c echo.Context, page theme.View, useAjax bool) (string, error) {
	layout := ctr.t.Layout(page.Prop("layout"))

	ctr.pageContents = ""
	ctr.componentCtx = nil
	ctr.cur = &current{
		Debug:      ctr.s.Cfg.App.Debug,
		Request:    c.Request(),
		Theme:      ctr.t,
		Page:       page,
		Layout:     layout,
		RouteParam: ctr.router.Params(),
		Param:      make(map[string]interface{}),
		Component:  make(map[string]component.Component),
	}

	if err := ctr.initComponents(); err != nil {
		return "", err
	}

	if useAjax && c.Request().Method == echo.POST {
		if ajaxResponse, err := ctr.execAjaxHandlers(c); err != nil || len(ajaxResponse) > 0 {
			ctr.out = c.JSONBlob
			return ajaxResponse, err
		}

		if handler := c.Request().PostFormValue("_handler"); len(handler) > 0 {
			// TODO: verify csrf token
			handlerResponse, err := ctr.runAjaxHandler(handler, c)
			if err != nil {
				return "", err
			}
			if res, ok := handlerResponse.(bool); !ok || res == false {
				data, err := json.Marshal(handlerResponse)
				if err != nil {
					return "", err
				}
				ctr.out = c.JSONBlob
				return string(data), nil
			}
		}
	}

	if cycleResponse := ctr.execPageCycle(); len(cycleResponse) > 0 {
		return cycleResponse, nil
	}

	contents, err := page.Render(ctr.cur)
	if layout == nil || err != nil {
		return contents, err
	}
	ctr.pageContents = contents

	return layout.Render(ctr.cur)
}

func (ctr *controller) execPageCycle() string {
	if ctr.cur.Layout != nil {
		if result := ctr.cur.Layout.RunComps(ctr.cur.Request); len(result) > 0 {
			return result
		}
	}

	return ctr.cur.Page.RunComps(ctr.cur.Request)
}

func (ctr *controller) initComponents() error {
	if ctr.cur.Layout != nil {
		ctr.cur.Layout.ClearComps()
		for _, c := range ctr.cur.Layout.CfgComps() {
			if _, err := ctr.addComponent(c.Name, c.Alias, c.Props, true); err != nil {
				return err
			}
		}
	}

	ctr.cur.Page.ClearComps()
	for _, c := range ctr.cur.Page.CfgComps() {
		if _, err := ctr.addComponent(c.Name, c.Alias, c.Props, false); err != nil {
			return err
		}
	}

	return nil
}

func (ctr *controller) addComponent(name, alias string, props view.Props, addToLayout bool) (component.Component, error) {
	var v view.View
	v = ctr.cur.Page

	if addToLayout {
		v = ctr.cur.Layout
	}

	comp, err := ctr.compManager.MakeComponent(name, v, component.Props(props))
	if err != nil {
		return nil, err
	}

	comp.SetAlias(alias)
	ctr.cur.Component[alias] = comp

	if addToLayout {
		ctr.cur.Layout.AddComp(alias, comp)
	} else {
		ctr.cur.Page.AddComp(alias, comp)
	}

	ctr.setComponentPropertiesFromParams(comp, nil)
	comp.Init(ctr.s)

	return comp, nil
}

func (ctr *controller) setComponentPropertiesFromParams(comp component.Component, params map[string]interface{}) {
	for k, value := range comp.Props() {
		if matches := ctr.compParamRe.FindStringSubmatch(value); len(matches) == 2 {
			paramName := strings.TrimSpace(matches[1])
			newValue := ""

			if strings.HasPrefix(paramName, ":") {
				if paramValue, ok := ctr.router.Params()[paramName[1:]]; ok {
					newValue = paramValue
				}
			} else if params != nil {
				if v, ok := params[paramName]; ok {
					if vv, ok := v.(string); ok {
						newValue = vv
					}
				}
			}

			comp.SetProp(k, newValue)
			comp.SetExternalPropName(k, paramName)
		}
	}
}

func (ctr *controller) renderPage() template.HTML {
	return template.HTML(ctr.pageContents)
}

func (ctr *controller) renderPartial(name string, params ...view.Param) template.HTML {
	cur := ctr.cur
	ctr.cur = cur.Copy()
	for _, p := range params {
		ctr.cur.Param[p.Name] = p.Value
	}

	var partial theme.View

	if index := strings.Index(name, "::"); index != -1 {
		alias, partialName := name[:index], name[index+2:]
		var comp component.Component

		if len(alias) == 0 {
			if ctr.componentCtx != nil {
				comp = ctr.componentCtx
			} else if comp = ctr.findComponentByPartial(partialName); comp == nil {
				ctr.s.Logger.Warnf("component %s not found", partialName)
				ctr.cur = cur
				return ""
			}
		} else if comp = ctr.findComponentByName(alias); comp == nil {
			ctr.s.Logger.Warnf("component %s not found", alias)
			ctr.cur = cur
			return ""
		}

		ctr.componentCtx = comp

		if partial = ctr.getComponentPartial(comp, partialName); partial == nil {
			ctr.s.Logger.Warnf("partial %s not found", name)
			ctr.cur = cur
			return ""
		}

		ctr.cur.Self = comp
	} else if partial = ctr.t.Partial(name); partial == nil {
		ctr.s.Logger.Warnf("partial %s not found", name)
		ctr.cur = cur
		return ""
	}

	ctr.partialStack.StackPartial()

	partial.ClearComps()
	for _, c := range partial.CfgComps() {
		comp, err := ctr.compManager.MakeComponent(c.Name, ctr.cur.Page, component.Props(c.Props))
		if err != nil {
			ctr.partialStack.UnstackPartial()
			ctr.cur = cur
			ctr.s.Logger.Warn(err)
			return ""
		}
		comp.SetAlias(c.Alias)
		ctr.cur.Component[c.Alias] = comp
		partial.AddComp(c.Alias, comp)
		ctr.partialStack.AddComponent(c.Alias, comp)

		ctr.setComponentPropertiesFromParams(comp, ctr.cur.Param)
		comp.Init(ctr.s)
	}

	partial.RunComps(ctr.cur.Request)

	partialContent, err := partial.Render(ctr.cur)

	ctr.partialStack.UnstackPartial()
	ctr.cur = cur

	if err != nil {
		ctr.s.Logger.Warn(err)
		return ""
	}

	return template.HTML(partialContent)
}

func (ctr *controller) renderComponent(name string, params ...view.Param) template.HTML {
	result := ""
	prevCtx := ctr.componentCtx

	if c := ctr.findComponentByName(name); c != nil {
		for _, p := range params {
			c.SetProp(p.Name, fmt.Sprintf("%v", p.Value))
		}

		ctr.componentCtx = c
		result = c.OnRender()

		if len(result) == 0 {
			f := c.Details().ViewFile
			if len(f) == 0 {
				f = "default"
			}
			result = string(ctr.renderPartial(fmt.Sprintf("%s::%s", name, f)))
		}
	}

	ctr.componentCtx = prevCtx

	return template.HTML(result)
}

func (ctr *controller) findComponentByName(name string) component.Component {
	if c := ctr.cur.Page.Comp(name); c != nil {
		return c
	}

	if ctr.cur.Layout != nil {
		if c := ctr.cur.Layout.Comp(name); c != nil {
			return c
		}
	}

	if partialComp := ctr.partialStack.Component(name); partialComp != nil {
		return partialComp
	}

	return nil
}

func (ctr *controller) findComponentByPartial(partial string) component.Component {
	for _, c := range ctr.cur.Page.AllComps() {
		if ctr.getComponentPartial(c, partial) != nil {
			return c
		}
	}
	if ctr.cur.Layout != nil {
		for _, c := range ctr.cur.Layout.AllComps() {
			if ctr.getComponentPartial(c, partial) != nil {
				return c
			}
		}
	}
	return nil
}

func (ctr *controller) getComponentPartial(comp component.Component, name string) theme.View {
	name = fmt.Sprintf("%s/%s", comp.Alias(), name)
	if v := ctr.t.Partial(name); v != nil {
		return v
	}
	if p := ctr.compManager.FindPlugin(comp); p != nil {
		if v := ctr.t.ComponentPartial(p.Dir(), name); v != nil {
			return v
		}
	}
	return nil
}

func (ctr *controller) findComponentByHandler(handler string) component.Component {
	if c := ctr.cur.Page.FindByHandler(handler); c != nil {
		return c
	}
	if ctr.cur.Layout != nil {
		if c := ctr.cur.Layout.FindByHandler(handler); c != nil {
			return c
		}
	}
	return nil
}

func (ctr *controller) getAjaxHandler(c echo.Context) string {
	if !helper.Ajax(c.Request()) || c.Request().Method != echo.POST {
		return ""
	}
	if handler := c.Request().Header.Get(HeaderRequestHandler); len(handler) > 0 {
		return handler
	}
	return ""
}

func (ctr *controller) execAjaxHandlers(c echo.Context) (string, error) {
	handler := ctr.getAjaxHandler(c)
	if len(handler) == 0 {
		return "", nil
	}

	if !ctr.handlerRe.MatchString(handler) {
		return "", fmt.Errorf("ajax handler invalid name: %s", handler)
	}

	partialList := strings.Split(strings.TrimSpace(c.Request().Header.Get(HeaderRequestPartials)), "&")
	for _, p := range partialList {
		if !ctr.partialRe.MatchString(p) {
			return "", fmt.Errorf("partial invalid name: %s", p)
		}
	}

	result, err := ctr.runAjaxHandler(handler, c)
	if err != nil {
		return "", err
	}
	if rAjax, ok := result.(bool); ok && rAjax == false {
		return "", fmt.Errorf("ajax handler %s not found", handler)
	}

	responseContents := make(map[string]interface{})
	for _, p := range partialList {
		responseContents[p] = string(ctr.renderPartial(p))
	}

	rv := reflect.ValueOf(result)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.IsValid() {
		switch rv.Kind() {
		case reflect.Map:
			for _, key := range rv.MapKeys() {
				responseContents[fmt.Sprintf("%v", key.Interface())] = rv.MapIndex(key).Interface()
			}
		case reflect.String, reflect.Slice:
			responseContents["result"] = rv.Interface()
		case reflect.Interface, reflect.Struct:
			data, err := json.Marshal(rv.Interface())
			if err != nil {
				return "", err
			}
			return string(data), nil
		}
	}

	data, err := json.Marshal(responseContents)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (ctr *controller) runAjaxHandler(handler string, c echo.Context) (interface{}, error) {
	if index := strings.Index(handler, "::"); index != -1 {
		componentName, handlerName := handler[:index], handler[index+2:]
		if comp := ctr.findComponentByName(componentName); comp != nil {
			if r, err, ok := callCompMethod(comp, handlerName, c); ok {
				ctr.componentCtx = comp
				if r == nil {
					return true, err
				}
				return r, err
			}
		}
	} else if comp := ctr.findComponentByHandler(handler); comp != nil {
		if r, err, ok := callCompMethod(comp, handler, c); ok {
			ctr.componentCtx = comp
			if r == nil {
				return true, err
			}
			return r, err
		}
	}
	if handler == "OnAjax" {
		return true, nil
	}
	return false, nil
}

func callCompMethod(comp component.Component, method string, c echo.Context) (interface{}, error, bool) {
	if m := reflect.ValueOf(comp).MethodByName(method); m.IsValid() {
		t := m.Type()
		if t.NumIn() == 1 && reflect.TypeOf(c).Implements(t.In(0)) {
			var r interface{}
			var err error
			if response := m.Call([]reflect.Value{reflect.ValueOf(c)}); len(response) > 0 && len(response) <= 2 {
				for _, i := range response {
					if e, ok := i.Interface().(error); ok {
						err = e
					} else {
						r = i.Interface()
					}
				}
			}
			return r, err, true
		}
	}
	return nil, nil, false
}
