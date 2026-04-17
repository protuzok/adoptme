package v1

import (
	"adoptme/internal/controller/restapi/v1/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *V1) listAnimals(c echo.Context) error {
	animals, err := r.c.ListAnimal(c.Request().Context())
	if err != nil {
		r.l.Error(err, "restapi - v1 - listAnimals")

		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}

	return c.JSON(http.StatusOK, animals)
}

func (r *V1) listShelters(c echo.Context) error {
	shelters, err := r.c.ListShelters(c.Request().Context())
	if err != nil {
		r.l.Error(err, "restapi - v1 - listShelters")

		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}

	return c.JSON(http.StatusOK, shelters)
}

func (r *V1) listVolunteers(c echo.Context) error {
	volunteers, err := r.c.ListVolunteer(c.Request().Context())
	if err != nil {
		r.l.Error(err, "restapi - v1 - listVolunteers")

		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}

	return c.JSON(http.StatusOK, volunteers)
}
