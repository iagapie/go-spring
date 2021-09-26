package helper

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func Ajax(r *http.Request) bool {
	return strings.EqualFold("XMLHttpRequest", r.Header.Get(echo.HeaderXRequestedWith))
}
