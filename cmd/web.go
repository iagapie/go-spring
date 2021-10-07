package cmd

import (
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/iagapie/go-spring/modules/backend/auth"
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/cms/controller"
	"github.com/iagapie/go-spring/modules/cms/theme"
	"github.com/iagapie/go-spring/modules/sys/datasource"
	"github.com/iagapie/go-spring/modules/sys/plugin"
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/token"
	"github.com/iagapie/go-spring/modules/sys/view"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
	global, err := initGlobalData(ctx)
	if err != nil {
		return err
	}
	defer global.db.Close()

	global.log.Infoln("redis initializing")
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.cfg.Redis.Addr,
		Password: global.cfg.Redis.Password,
		DB:       0,
	})
	defer func() {
		if err = rdb.Close(); err != nil {
			log.Error(err)
		}
	}()

	global.log.Infoln("redis cache initializing")
	redisCache := cache.New(&cache.Options{
		Redis: rdb,
	})

	global.log.Infoln("token manager initializing")
	tokenManager := token.New(token.WithJWTKeys(global.cfg.JWT.SigningKeys))

	global.log.Infoln("auth service initializing")
	authService := auth.NewService(global.cfg.JWT.TTL, global.userService, redisCache, tokenManager, global.log.Entry)

	global.log.Infoln("plugin manager initializing")
	plugManager, err := plugin.New(global.cfg.CMS.PluginsPath, global.log)
	if err != nil {
		return err
	}

	global.log.Infoln("spring (echo) framework initializing")
	s := spring.New(global.cfg, global.log)
	view.Add("routeURL", s.Reverse)

	theme.SetThemesPath(fmt.Sprintf("%s/frontend", global.cfg.CMS.ThemesPath))
	theme.SetDatasource(datasource.NewFile(global.log))
	theme.SetActiveTheme(global.cfg.CMS.ActiveTheme)

	for _, t := range theme.Themes() {
		s.Frontend.Static(t.Assets())
	}

	plugManager.RegisterAll(s)

	global.log.Infoln("component manager initializing")
	compManager := component.New(plugManager)

	authHandler := &auth.Handler{Service: authService}
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
