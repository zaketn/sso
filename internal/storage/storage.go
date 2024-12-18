package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user now found")
	ErrAppNowFound  = errors.New("app not found")
)
