package domain

import "errors"

var (
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptymessage = errors.New("message cannot be empty")
)