package v1

import (
	"adoptme/internal/controller/restapi/v1/request"
	"adoptme/internal/controller/restapi/v1/response"
	"adoptme/internal/entity"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *V1) volunteerRegister(c echo.Context) error {
	var body request.RegisterVolunteer

	if err := c.Bind(&body); err != nil {
		r.l.Error(err, "restapi - v1 - volunteerRegister")

		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - volunteerRegister")

		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	user, err := r.vl.Register(c.Request().Context(), body.Name, body.Email, body.Password)
	if err != nil {
		r.l.Error(err, "restapi - v1 - volunteerRegister")

		if errors.Is(err, entity.ErrUserAlreadyExists) {
			return c.JSON(http.StatusConflict, response.Error{Error: "user already exists"})
		}

		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (r *V1) volunteerLogin(c echo.Context) error {
	var body request.LoginVolunteer

	if err := c.Bind(&body); err != nil {
		r.l.Error(err, "restapi - v1 - volunteerLogin")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - volunteerLogin")
		return c.JSON(http.StatusBadRequest, response.Error{Error: "invalid request body"})
	}

	token, err := r.vl.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		r.l.Error(err, "restapi - v1 - volunteerLogin")

		if errors.Is(err, entity.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, response.Error{Error: "invalid credentials"})
		}
		
		return c.JSON(http.StatusInternalServerError, response.Error{Error: "internal server error"})
	}
	return c.JSON(http.StatusOK, response.Token{Token: token})
}
