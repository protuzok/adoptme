package v1

import (
	"adoptme/internal/controller/restapi/v1/request"
	"adoptme/internal/controller/restapi/v1/response"
	"adoptme/internal/entity"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *V1) registerShelter(c echo.Context) error {
	var body request.RegisterShelter

	if err := c.Bind(&body); err != nil {
		r.l.Error(err, "restapi - v1 - registerShelter")

		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - registerShelter")

		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	user, err := r.sh.Register(c.Request().Context(), body.Name, body.Email, body.Password)
	if err != nil {
		r.l.Error(err, "restapi - v1 - registerShelter")

		if errors.Is(err, entity.ErrUserAlreadyExists) {
			return c.JSON(http.StatusConflict, response.Error{Error: "user already exists"})
		}

		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (r *V1) loginShelter(c echo.Context) error {
	var body request.LoginShelter

	if err := c.Bind(&body); err != nil {
		r.l.Error(err, "restapi - v1 - loginShelter")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - loginShelter")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	token, err := r.sh.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		r.l.Error(err, "restapi - v1 - loginShelter")

		if errors.Is(err, entity.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, response.Error{Error: "invalid credentials"})
		}
		
		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}
	return c.JSON(http.StatusOK, response.Token{Token: token})
}
