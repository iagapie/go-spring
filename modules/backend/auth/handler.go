package auth

import (
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	signInURL  = "/sign-in"
	refreshURL = "/refresh"
)

type Handler struct {
	Service Service
}

func (h *Handler) Register(b *spring.Backend) {
	mp := []string{echo.POST, echo.OPTIONS}
	b.Match(mp, signInURL, h.signIn)[0].Name = "backend-sign-in"
	b.Match(mp, refreshURL, h.refresh)[0].Name = "backend-refresh"
}

func (h *Handler) signIn(c echo.Context) error {
	c.Logger().Info("BACKEND SIGN IN HANDLER")

	var dto SignInDTO

	c.Logger().Debug("bind SignInDTO")
	if err := c.Bind(&dto); err != nil {
		return err
	}

	c.Logger().Debug("validate SignInDTO")
	if err := c.Validate(&dto); err != nil {
		return err
	}

	r, err := h.Service.Auth(c.Request().Context(), dto)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &r)
}

func (h *Handler) refresh(c echo.Context) error {
	c.Logger().Info("BACKEND REFRESH HANDLER")

	var dto RefreshTokenDTO

	c.Logger().Debug("bind RefreshTokenDTO")
	if err := c.Bind(&dto); err != nil {
		return err
	}

	c.Logger().Debug("validate RefreshTokenDTO")
	if err := c.Validate(&dto); err != nil {
		return err
	}

	r, err := h.Service.RefreshToken(c.Request().Context(), dto)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &r)
}
