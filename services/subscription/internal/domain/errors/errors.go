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
	InvalidTokenErrorCode DomainErrorCode = "INVALID_TOKEN_ERROR"
)

var (
	AlreadySubscribedError = apperrors.NewConflict(AlreadySubscribedCode, "email already subscribed")
	TokenNotFoundError     = apperrors.NewNotFound(TokenNotFoundCode, "there is no subscription with such token")
	InvalidTokenError      = apperrors.NewBadRequest(InvalidTokenErrorCode, "invalid token")
)
