package thstrings

import "errors"

var (
	ErrCompile         = errors.New("could not be compiled")
	ErrFailedExecution = errors.New("failed execution")
	ErrEmptyTemplate   = errors.New("empty template, nothing to execute")
)
