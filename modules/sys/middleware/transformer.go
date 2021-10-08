package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
)

type TransformFunc func(ctx context.Context, item interface{}) (interface{}, error)

func Transformer(fromContextKey, toContextKey string, transformFunc TransformFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			item, err := transformFunc(c.Request().Context(), c.Get(fromContextKey))
			if err != nil {
				return err
			}
			c.Set(toContextKey, item)
			return next(c)

		}
	}
}
