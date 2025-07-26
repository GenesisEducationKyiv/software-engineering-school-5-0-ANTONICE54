package errors

import (
	"errors"
)

var (
	ErrDatabase          = errors.New("database raised an error")
	ErrInternal          = errors.New("internal server error")
	ErrUnknownEventRoute = errors.New("unknown event route")
)
