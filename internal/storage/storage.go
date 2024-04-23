package storage

import "errors"

var (
	ErrUsernameTaken = errors.New("username is already taken")
)
