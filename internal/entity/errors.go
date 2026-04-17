package entity

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrAnimalNotFound     = errors.New("animal not found")
	ErrNotFound           = errors.New("not found") // todo delete from catalog_test
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)
