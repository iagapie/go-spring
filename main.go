package main

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/cms/controller"
	"github.com/iagapie/go-spring/modules/cms/theme"
	"github.com/iagapie/go-spring/modules/sys/config"
	"github.com/iagapie/go-spring/modules/sys/datasource"
	"github.com/iagapie/go-spring/modules/sys/helper"
	"github.com/iagapie/go-spring/modules/sys/logger"
	"github.com/iagapie/go-spring/modules/sys/plugin"
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/view"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	var cfg config.Cfg
	if err := helper.ReadConfig(&cfg, "./configs/app", "./configs/cms"); err != nil {
		panic(err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	lvl := logrus.InfoLevel
	if cfg.App.Debug {
		lvl = logrus.DebugLevel
	}
	log := logger.New(logger.WithLevel(lvl))

	plugManager, err := plugin.New(cfg.CMS.PluginsPath, log)
	if err != nil {
		panic(err)
	}

	s := spring.New(cfg, log)
	view.Add("route_url", s.Reverse)

	theme.SetThemesPath(fmt.Sprintf("%s/frontend", cfg.CMS.ThemesPath))
	theme.SetDatasource(datasource.NewFile(log))
	theme.SetActiveTheme(cfg.CMS.ActiveTheme)

	for _, t := range theme.Themes() {
		s.Frontend.Static(t.Assets())
	}

	plugManager.RegisterAll(s)

	compManager := component.New(plugManager)

	s.HTTPErrorHandler = func(err error, c echo.Context) {
		// TODO: add backend check
		ctr := controller.New(s, compManager)
		ctr.Error(err, c)
	}

	s.Frontend.Any("/", func(c echo.Context) error {
		ctr := controller.New(s, compManager)
		return ctr.Run(c)
	})

	s.Frontend.Any("/*", func(c echo.Context) error {
		ctr := controller.New(s, compManager)
		return ctr.Run(c)
	})

	plugManager.RoutesAll(s.Frontend, s.Backend)

	if err = s.Run(); err != nil {
		s.Logger.Fatalf("Error occurred while running HTTP runner: %v", err)
	}
}
