package storage

import "errors"

var (
	ErrUsernameTaken = errors.New("username is already taken")
	ErrAliasTaken    = errors.New("alias is already taken")
)
