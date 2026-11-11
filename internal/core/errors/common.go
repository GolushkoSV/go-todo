package core_errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("bad arguments")
	ErrConflict        = errors.New("conflict")
)
