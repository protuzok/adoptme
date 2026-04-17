package v1

import (
	"adoptme/internal/usecase"
	"adoptme/pkg/logger"

	"github.com/go-playground/validator/v10"
)

type V1 struct {
	a  usecase.Adoption
	c  usecase.Catalog
	sh usecase.Shelter
	vl usecase.Volunteer

	l logger.Interface
	v *validator.Validate
}
