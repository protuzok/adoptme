package restapi

import (
	"adoptme/internal/controller/restapi/middleware"
	"adoptme/internal/controller/restapi/v1"
	"adoptme/internal/usecase"
	"adoptme/pkg/jwt"
	"adoptme/pkg/logger"

	"github.com/labstack/echo/v4"
)

func NewRouter(app *echo.Echo, a usecase.Adoption, c usecase.Catalog, sh usecase.Shelter, vl usecase.Volunteer, jwt *jwt.Manager, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, a, c, sh, vl, jwt, l)
	}
}
