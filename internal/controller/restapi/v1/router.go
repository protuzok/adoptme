package v1

import (
	"adoptme/internal/controller/restapi/middleware"
	"adoptme/internal/usecase"
	"adoptme/pkg/jwt"
	"adoptme/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func NewRoutes(apiV1Group *echo.Group, a usecase.Adoption, c usecase.Catalog, sh usecase.Shelter, vl usecase.Volunteer, jwt *jwt.Manager, l logger.Interface) {
	r := &V1{
		a:  a,
		c:  c,
		sh: sh,
		vl: vl,
		l:  l,
		v:  validator.New(validator.WithRequiredStructEnabled()),
	}

	// todo для кожного хенделера треба буде зробити swag-документацію
	// Public routes
	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.POST("/shelter/register", r.registerShelter)
		authGroup.POST("/shelter/login", r.loginShelter)

		authGroup.POST("/volunteer/register", r.volunteerRegister)
		authGroup.POST("/volunteer/login", r.volunteerLogin)
	}

	// Protected routes
	protected := apiV1Group.Group("", middleware.Auth(jwt))

	adoptionGroup := protected.Group("/adoption")
	{
		adoptionGroup.POST("/", r.registerAnimal)
		adoptionGroup.POST("/:id/transfer", r.transfer)
	}

	catalogGroup := protected.Group("/catalog")
	{
		catalogGroup.GET("/animals", r.listAnimals)
		catalogGroup.GET("/shelters", r.listShelters)
		catalogGroup.GET("/volunteers", r.listVolunteers)
	}
}
