package entity

import "errors"

var (
	ErrEmptyValue   = errors.New("can not be empty")
	ErrInvalidValue = errors.New("invalid value")
	ErrNotFound     = errors.New("not found")
	ErrCompile      = errors.New("could not be compiled")
)
