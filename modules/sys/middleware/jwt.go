package middleware

import (
	"github.com/iagapie/go-spring/modules/sys/config"
	"github.com/iagapie/go-spring/modules/sys/token"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWT(cfg config.JWT, t token.Token) echo.MiddlewareFunc {
	jwtCfg := middleware.JWTConfig{
		Skipper: middleware.DefaultSkipper,
		ContextKey: cfg.ContextKey,
		TokenLookup: cfg.TokenLookup,
		AuthScheme: cfg.AuthScheme,
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			return t.Validate(auth)
		},
	}
	return middleware.JWTWithConfig(jwtCfg)
}
