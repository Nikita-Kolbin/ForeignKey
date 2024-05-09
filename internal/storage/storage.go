package storage

import "errors"

var (
	ErrLoginTaken = errors.New("login is already taken")
	ErrAliasTaken = errors.New("alias is already taken")
)
