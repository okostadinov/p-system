package models

import "errors"

var (
	ErrNoRecord           = errors.New("no record")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
