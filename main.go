package main

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/controller"
	"github.com/iagapie/go-spring/pkg/manager"
	"github.com/iagapie/go-spring/pkg/spring"
	"github.com/iagapie/go-spring/pkg/template"
	"github.com/iagapie/go-spring/pkg/theme"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	template2 "html/template"
	"net/http"
	"strings"
)

func main() {
	configManger, err := manager.NewConfigManager("./configs/app", "./configs/cms")
	if err != nil {
		panic(err)
	}
	cfg := configManger.Config()

	lvl := logrus.InfoLevel
	if cfg.App.Debug {
		lvl = logrus.DebugLevel
	}
	logManager := manager.NewLogManager(manager.WithLogLevel(lvl))
	log := logManager.Logrus()

	theme.SetThemesPath(cfg.CMS.ThemesPath + "/frontend")
	theme.SetActiveTheme(cfg.CMS.ActiveTheme)

	t := theme.ActiveTheme()

	pluginManager, err := manager.NewPluginManager(cfg.CMS.PluginsPath, logManager.Logrus())
	if err != nil {
		log.Fatalln(err)
	}

	navManager := manager.NewNavigationManager(pluginManager)

	m := &manager.Manager{
		ConfigManager:     configManger,
		LogManager:        logManager,
		PluginManager:     pluginManager,
		NavigationManager: navManager,
	}

	pluginManager.RegisterAll(m)

	s := spring.New(configManger.Config(), logManager.Echo())
	ctr := controller.New(s, t)

	backendRenderer := template.New(
		"backend",
		fmt.Sprintf("%s/backend/views", cfg.CMS.ThemesPath),
		template.WithGlobal("backend_uri", cfg.CMS.BackendURI),
		template.WithFunc("url", s.Reverse),
		template.WithFunc("backend_url", func(uri string) string {
			return fmt.Sprintf("%s/%s", cfg.CMS.BackendURI, strings.TrimLeft(uri, "/"))
		}),
	)

	s.Backend.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Echo().Renderer = backendRenderer
			return next(c)
		}
	})

	pluginManager.BootAll(s)

	s.Backend.Static("", fmt.Sprintf("%s/backend/assets", cfg.CMS.ThemesPath))
	s.Backend.GET("/dashboard", func(c echo.Context) error {
		return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
			"name": "Dolly!",
			"foo":  template2.HTML("<h2>Fooooo</h2>"),
			"boo": map[string]interface{}{
				"boo": "BOOOOOO",
			},
		})
	}).Name = "dashboard"

	s.Frontend.Static(t.Assets())
	s.Frontend.Any("/", ctr.Run)
	s.Frontend.Any("/*", ctr.Run)

	pluginManager.RoutesAll(s.Frontend, s.Backend)

	if err = s.Run(); err != nil {
		s.Logger.Fatalf("Error occurred while running HTTP runner: %v", err)
	}
}
