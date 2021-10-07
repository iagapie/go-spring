package handler

import (
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/token"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	signInURL  = "/sign-in"
	refreshURL = "/refresh"
)

type (
	SignInDTO struct {
		Email    string `json:"email,omitempty" validate:"required,email,min=3,max=255"`
		Password string `json:"password,omitempty" validate:"required,min=8,max=64"`
	}

	RefreshTokenDTO struct {
		Token string `json:"token,omitempty" validate:"required,uuid4"`
	}

	Tokens struct {
		AccessToken  string `json:"access_token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	AuthHandler struct {
		Access       time.Duration
		Refresh      time.Duration
		TokenManager token.Token
	}
)

func (h *AuthHandler) Register(b *spring.Backend) {
	m := []string{echo.POST, echo.OPTIONS}
	b.Match(m, signInURL, h.signIn)[0].Name = "backend-sign-in"
	b.Match(m, refreshURL, h.refresh)[0].Name = "backend-refresh"
}

func (h *AuthHandler) signIn(c echo.Context) error {
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

	return c.JSON(http.StatusOK, &dto)
}

func (h *AuthHandler) refresh(c echo.Context) error {
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

	return c.JSON(http.StatusOK, &dto)
}
