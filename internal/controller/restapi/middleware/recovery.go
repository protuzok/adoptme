package middleware

import (
	"adoptme/pkg/logger"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/labstack/echo/v4"
)

func buildPanicMessage(c echo.Context, err any) string {
	var result strings.Builder

	result.WriteString(c.RealIP())
	result.WriteString(" - ")
	result.WriteString(c.Request().Method)
	result.WriteString(" ")
	result.WriteString(c.Request().RequestURI)
	result.WriteString(" PANIC DETECTED: ")
	_, _ = fmt.Fprintf(&result, "%v\n%s\n", err, debug.Stack())

	return result.String()
}

func Recovery(l logger.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					l.Error(buildPanicMessage(c, r))

					c.Error(echo.ErrInternalServerError)
				}
			}()

			return next(c)
		}
	}
}
