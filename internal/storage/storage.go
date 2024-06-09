package storage

import "errors"

var (
	ErrURLNotExists = errors.New("url not found")
	ErrURLExists    = errors.New("url exists")
	ErrURLNotFound  = errors.New("url not found")
)
