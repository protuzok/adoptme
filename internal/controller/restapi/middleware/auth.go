package middleware

import (
	"adoptme/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const _bearerParts = 2

type errorResponse struct {
	Error string `json:"error"`
}

// Auth returns a JWT authentication middleware for Echo.
func Auth(jwtManager *jwt.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.JSON(http.StatusUnauthorized, errorResponse{Error: "missing authorization header"})
			}

			parts := strings.SplitN(header, " ", _bearerParts)
			if len(parts) != _bearerParts || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid authorization header format"})
			}

			userID, err := jwtManager.ParseToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid or expired token"})
			}

			c.Set("userID", userID)

			return next(c)
		}
	}
}
