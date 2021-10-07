package spring

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/iagapie/go-spring/modules/sys/config"
	middleware2 "github.com/iagapie/go-spring/modules/sys/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	Version = "1.0.0"
	timeout = 5 * time.Second
)

type (
	Spring struct {
		*echo.Echo
		Frontend *Frontend
		Backend *Backend
		Cfg     config.Cfg
	}

	Frontend struct {
		*echo.Group
	}

	Backend struct {
		*echo.Group
	}

	Context interface {
		echo.Context
		Spring() *Spring
	}

	spCtx struct {
		echo.Context
		spring *Spring
	}

	valid struct {
		v *validator.Validate
	}
)

func New(cfg config.Cfg, l echo.Logger) *Spring {
	e := echo.New()
	e.HideBanner = true

	s := &Spring{
		Echo: e,
		Cfg:  cfg,
	}
	s.Debug = s.Cfg.App.Debug
	s.Logger = l
	s.StdLogger = log.New(s.Logger.Output(), s.Logger.Prefix()+": ", 0)
	s.Validator = &valid{validator.New()}

	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&spCtx{c, s})
		}
	})
	s.Use(middleware.Logger(), middleware.Recover(), middleware.CORS(), middleware2.ValidateError())

	s.Frontend = &Frontend{
		Group: s.Group(""),
	}

	s.Backend = &Backend{
		Group: s.Group(s.Cfg.CMS.BackendURI),
	}

	return s
}

func (s *Spring) Run() error {
	s.Logger.Printf("Spring v%s, Echo v%s", Version, echo.Version)

	go s.graceful(syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)

	addr := fmt.Sprintf(":%d", s.Cfg.App.Port)
	if err := s.Start(addr); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.Logger.Info("Spring shutdown")
		} else {
			return err
		}
	}

	return nil
}

func (s *Spring) graceful(signals ...os.Signal) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, signals...)
	sig := <-quit

	s.Logger.Infof("Caught signal %s. Shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Fatal(err)
	}
}

func (c *spCtx) Spring() *Spring {
	return c.spring
}

func (v *valid) Validate(i interface{}) error {
	return v.v.Struct(i)
}

func ToSpringContext(c echo.Context) Context {
	return c.(Context)
}

func IsBackend(c echo.Context) bool {
	sc := ToSpringContext(c)
	return strings.HasPrefix(sc.Request().RequestURI, sc.Spring().Cfg.CMS.BackendURI)
}
