package storage

import "errors"

var (
	ErrUserExists          = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserSessionNotFound = errors.New("user session not found")
	ErrAppNotFound         = errors.New("app not found")
)
