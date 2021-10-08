package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/iagapie/go-spring/modules/backend/auth"
	"github.com/iagapie/go-spring/modules/backend/user"
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/cms/controller"
	"github.com/iagapie/go-spring/modules/cms/theme"
	"github.com/iagapie/go-spring/modules/sys/datasource"
	"github.com/iagapie/go-spring/modules/sys/middleware"
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
	data, err := initData(ctx)
	if err != nil {
		return err
	}
	defer data.db.Close()

	data.log.Infoln("redis initializing")
	rdb := redis.NewClient(&redis.Options{
		Addr:     data.cfg.Redis.Addr,
		Password: data.cfg.Redis.Password,
		DB:       0,
	})
	defer func() {
		if err = rdb.Close(); err != nil {
			log.Error(err)
		}
	}()

	data.log.Infoln("redis cache initializing")
	redisCache := cache.New(&cache.Options{
		Redis: rdb,
	})

	data.log.Infoln("token manager initializing")
	tokenManager := token.New(token.WithJWTKeys(data.cfg.JWT.SigningKeys))

	data.log.Infoln("auth service initializing")
	authService := auth.NewService(data.cfg.JWT.TTL, data.userService, redisCache, tokenManager, data.log.Entry)

	data.log.Infoln("plugin manager initializing")
	plugManager, err := plugin.New(data.cfg.CMS.PluginsPath, data.log)
	if err != nil {
		return err
	}

	data.log.Infoln("spring (echo) framework initializing")
	s := spring.New(data.cfg, data.log)
	view.Add("routeURL", s.Reverse)

	data.log.Infoln("theme initializing")
	theme.SetThemesPath(fmt.Sprintf("%s/frontend", data.cfg.CMS.ThemesPath))
	theme.SetDatasource(datasource.NewFile(data.log))
	theme.SetActiveTheme(data.cfg.CMS.ActiveTheme)

	for _, t := range theme.Themes() {
		s.Frontend.Static(t.Assets())
	}

	data.log.Infoln("plugin manager: RegisterAll")
	plugManager.RegisterAll(s)

	data.log.Infoln("component manager initializing")
	compManager := component.New(plugManager)

	data.log.Infoln("backend authentication handler initializing")
	authHandler := &auth.Handler{Service: authService}
	authHandler.Register(s.Backend)

	userContextKey := "user"
	userTransformFunc := func(ctx context.Context, item interface{}) (interface{}, error) {
		return data.userService.GetByUUID(ctx, item.(string))
	}
	userMiddleware := middleware.Transformer(data.cfg.JWT.ContextKey, userContextKey, userTransformFunc)
	jwtMiddleware := middleware.JWT(data.cfg.JWT, tokenManager)

	data.log.Infoln("backend user handler initializing")
	userHandler := &user.Handler{
		Service:        data.userService,
		JWTMiddleware:  jwtMiddleware,
		UserMiddleware: userMiddleware,
		UserContextKey: userContextKey,
	}
	userHandler.Register(s.Backend)

	data.log.Infoln("cms controller initializing")
	s.HTTPErrorHandler = func(err error, c echo.Context) {
		if errors.Is(err, user.ErrRecordNotFound) {
			err = echo.ErrNotFound.SetInternal(err)
		}

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

	data.log.Infoln("plugin manager: RoutesAll")
	plugManager.RoutesAll(s.Frontend, s.Backend)

	if err = s.Run(); err != nil {
		return fmt.Errorf("error occurred while running HTTP runner: %v", err)
	}

	return nil
}
