package main

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/sys/plugin"
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/view"
	"github.com/iagapie/go-spring/plugins/spring/demo/components"
	"github.com/labstack/echo/v4"
	"net/http"
)

var Plugin plugin.Plugin

type plug struct {
	details plugin.Details
	s       *spring.Spring
}

func (p *plug) Details() plugin.Details {
	return p.details
}

func (p *plug) Register(s *spring.Spring) {
	p.s = s
}

func (p *plug) Routes(f *spring.Frontend, b *spring.Backend) {
	f.GET("/welcome", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Frontend: %s - %s", p.details.Name, p.s.Cfg.App.Name))
	}).Name = "welcome"

	b.GET("/welcome", func(c echo.Context) error {
		return c.String(http.StatusOK, "Backend welcome - "+p.s.Cfg.App.Name)
	})
}

func (p *plug) RegisterComponents() component.FactoryMap {
	return component.FactoryMap{
		"todo": func(v view.View, props component.Props) (component.Component, error) {
			return components.NewTodo(props), nil
		},
	}
}

func init() {
	Plugin = &plug{
		details: plugin.Details{
			Code:        "spring_demo",
			Name:        "Spring Demo Plugin",
			Version:     "1.0.0",
			Description: "This is a spring demo plugin",
			Author:      "Igor Agapie",
		},
	}
}
