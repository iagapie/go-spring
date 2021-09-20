package main

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/manager"
	"github.com/iagapie/go-spring/pkg/spring"
	"github.com/iagapie/go-spring/plugins/welcome/components"
	"github.com/labstack/echo/v4"
	"net/http"
)

var Plugin manager.Plugin

type plug struct {
	details manager.PluginDetails
	m       *manager.Manager
}

func (p *plug) Details() manager.PluginDetails {
	return p.details
}

func (p *plug) Register(m *manager.Manager) {
	p.m = m
}

func (p *plug) Boot(s *spring.Spring) {
	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			s.Logger.Infof("Info request from plugin %s - %s - %s", p.details.Code, p.details.Name, p.details.Version)
			if spring.IsBackend(c) {
				s.Logger.Info("Request is for backend")
			} else {
				s.Logger.Info("Request is for frontend")
			}
			return next(c)
		}
	})
}

func (p *plug) RegisterComponents() map[string]manager.ComponentFactory {
	return map[string]manager.ComponentFactory{
		"todo": func(props manager.ComponentProps) (manager.Component, error) {
			return components.NewTodo(props), nil
		},
	}
}

func (p *plug) Routes(f *spring.Frontend, b *spring.Backend) {
	f.GET("/welcome", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Frontend: %s - %s", p.details.Name, p.m.ConfigManager.Config().App.Name))
	}).Name = "welcome"

	b.GET("/welcome", func(c echo.Context) error {
		return c.String(http.StatusOK, "Backend welcome - "+p.m.ConfigManager.Config().App.Name)
	})
}

func init() {
	Plugin = &plug{
		details: manager.PluginDetails{
			Code:        "welcome",
			Name:        "Welcome Plugin",
			Version:     "1.0.0",
			Description: "This is a demo plugin",
			Author:      "Igor Agapie",
		},
	}
}
