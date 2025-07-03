package errors

import (
	"weather-forecast/pkg/apperrors"
)

type DomainErrorCode string

func (c DomainErrorCode) String() string {
	return string(c)
}

const (
	AlreadySubscribedCode DomainErrorCode = "ALREADY_SUBSCRIBED"
	TokenNotFoundCode     DomainErrorCode = "TOKEN_NOT_FOUND"
)

var (
	AlreadySubscribedError = apperrors.NewConflict(AlreadySubscribedCode, "email already subscribed")
	TokenNotFoundError     = apperrors.NewNotFound(TokenNotFoundCode, "there is no subscription with such token")
)
