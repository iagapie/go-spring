package controller

import (
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/cms/theme"
	sysRouter "github.com/iagapie/go-spring/modules/sys/router"
	"net/http"
)

type current struct {
	Debug      bool
	Request    *http.Request
	Theme      theme.Theme
	Page       theme.View
	Layout     theme.View
	RouteParam sysRouter.Params
	Param      map[string]interface{}
	Component  map[string]component.Component
	Self       component.Component
}

func (cur *current) Copy() *current {
	rp := make(sysRouter.Params, len(cur.RouteParam))
	for k, v := range cur.RouteParam {
		rp[k] = v
	}
	p := make(map[string]interface{}, len(cur.Param))
	for k, v := range cur.Param {
		p[k] = v
	}
	c := make(map[string]component.Component, len(cur.Component))
	for k, v := range cur.Component {
		c[k] = v
	}
	return &current{
		Request:    cur.Request,
		Theme:      cur.Theme,
		Page:       cur.Page,
		Layout:     cur.Layout,
		RouteParam: rp,
		Param:      p,
		Component:  c,
	}
}
