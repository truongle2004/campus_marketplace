package user

import "errors"

var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user profile already exists")
	ErrInvalidCampus = errors.New("invalid campus")
	ErrEmailRequired = errors.New("email claim missing from token")
)
