package middleware

import (
	"adoptme/pkg/logger"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func buildRequestMessage(c echo.Context) string {
	var result strings.Builder

	result.WriteString(c.RealIP())
	result.WriteString(" - ")
	result.WriteString(c.Request().Method)
	result.WriteString(" ")
	result.WriteString(c.Request().URL.Path)
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(c.Response().Status))
	result.WriteString(" ")
	result.WriteString(strconv.FormatInt(c.Response().Size, 10))

	return result.String()
}

func Logger(l logger.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			l.Info(buildRequestMessage(c))
			return err
		}
	}
}
