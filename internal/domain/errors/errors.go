package errors

import "errors"

var (
	ErrAlreadySubscribed = errors.New("email already subscribed")
	ErrTokenNotFound     = errors.New("there is no subscription with such token")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidFrequency  = errors.New("unexpected frequency value")
)
