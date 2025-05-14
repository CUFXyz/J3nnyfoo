package storage

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrNotWritten    = errors.New("data is not written")
	ErrEmptyInput    = errors.New("input empty")
	ErrAlreadyExists = errors.New("link already exists")
)
