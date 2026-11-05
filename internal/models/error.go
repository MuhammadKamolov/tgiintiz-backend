package models

import "errors"

var (
	ErrEmailTaken         = errors.New("email is already taken")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrForbidden          = errors.New("access denied")
)
