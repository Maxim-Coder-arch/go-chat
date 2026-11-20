package domain

import "errors"

var (
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptyMessage  = errors.New("message cannot be empty")
)
