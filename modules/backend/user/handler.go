package user

import (
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	meURL = "/me"
)

type Handler struct {
	Service        Service
	JWTMiddleware  echo.MiddlewareFunc
	UserMiddleware echo.MiddlewareFunc
	UserContextKey string
}

func (h *Handler) Register(b *spring.Backend) {
	mg := []string{echo.GET, echo.OPTIONS}
	b.Match(mg, meURL, h.me, h.JWTMiddleware, h.UserMiddleware)
}

func (h *Handler) me(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get(h.UserContextKey))
}
