package cmd

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/backend/handler"
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/cms/controller"
	"github.com/iagapie/go-spring/modules/cms/theme"
	"github.com/iagapie/go-spring/modules/sys/config"
	"github.com/iagapie/go-spring/modules/sys/datasource"
	"github.com/iagapie/go-spring/modules/sys/helper"
	"github.com/iagapie/go-spring/modules/sys/logger"
	"github.com/iagapie/go-spring/modules/sys/plugin"
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/token"
	"github.com/iagapie/go-spring/modules/sys/view"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

var Web = &cli.Command{
	Name:  "web",
	Usage: "Start Spring CMS web server",
	Description: `Spring CMS web server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
}

func runWeb(ctx *cli.Context) error {
	var cfg config.Cfg
	if err := helper.ReadConfigWithEnv(&cfg, ctx.StringSlice("config")...); err != nil {
		return err
	}

	log := logger.New(logger.WithDebug(cfg.App.Debug))

	tokenManager := token.New(token.WithJWTKeys(cfg.JWT.SigningKeys))

	plugManager, err := plugin.New(cfg.CMS.PluginsPath, log)
	if err != nil {
		return err
	}

	s := spring.New(cfg, log)
	view.Add("routeURL", s.Reverse)

	theme.SetThemesPath(fmt.Sprintf("%s/frontend", cfg.CMS.ThemesPath))
	theme.SetDatasource(datasource.NewFile(log))
	theme.SetActiveTheme(cfg.CMS.ActiveTheme)

	for _, t := range theme.Themes() {
		s.Frontend.Static(t.Assets())
	}

	plugManager.RegisterAll(s)

	compManager := component.New(plugManager)

	authHandler := &handler.AuthHandler{
		TokenManager: tokenManager,
		Access:       cfg.JWT.TTL.Access,
		Refresh:      cfg.JWT.TTL.Refresh,
	}
	authHandler.Register(s.Backend)

	s.HTTPErrorHandler = func(err error, c echo.Context) {
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
		return fmt.Errorf("error occurred while running HTTP runner: %v", err)
	}

	return nil
}
