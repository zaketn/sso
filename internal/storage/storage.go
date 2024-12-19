package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user now found")
	ErrAppNotFound  = errors.New("app not found")
)
