package v1

import (
	"adoptme/internal/controller/restapi/v1/request"
	"adoptme/internal/controller/restapi/v1/response"
	"adoptme/internal/entity"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *V1) registerAnimal(c echo.Context) error {
	var body request.RegisterAnimal

	if err := c.Bind(&body); err != nil {
		r.l.Error(err, "restapi - v1 - registerAnimal")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - registerAnimal")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	animal, err := r.a.RegisterAnimal(c.Request().Context(), body.Name, body.OwnerID, body.OwnerType)
	if err != nil {
		r.l.Error(err, "restapi - v1 - registerAnimal")
		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}

	return c.JSON(http.StatusCreated, animal)
}

func (r *V1) transfer(c echo.Context) error {
	var body request.Transfer

	if err := c.Bind(&body); err != nil {
		r.l.Error(err, "restapi - v1 - transfer")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - transfer")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	animalID := c.Param("id")
	if err := r.v.Var(animalID, "required,uuid"); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid animal id"})
	}

	err := r.a.TransferAnimal(c.Request().Context(), animalID, body.NewOwnerID, body.NewOwnerType)
	if err != nil {
		r.l.Error(err, "restapi - v1 - transfer")

		if errors.Is(err, entity.ErrAnimalNotFound) {
			return c.JSON(http.StatusNotFound, response.Error{Error: "animal not found"})
		}

		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}

	return c.NoContent(http.StatusNoContent)
}
