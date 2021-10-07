package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/iagapie/go-spring/modules/sys/helper"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateError() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if vErrs, ok := err.(validator.ValidationErrors); ok {
				errs := make([]interface{}, 0, len(vErrs))
				for _, f := range vErrs {
					errs = append(errs, ErrorField{
						Field:   helper.ToSnakeCase(f.Field()),
						Message: fmt.Sprintf("Validation failed on the '%s' tag", f.Tag()),
					})
				}

				data, _ := json.Marshal(errs)

				return echo.NewHTTPError(http.StatusUnprocessableEntity, string(data)).SetInternal(err)
			}

			return err
		}
	}
}
