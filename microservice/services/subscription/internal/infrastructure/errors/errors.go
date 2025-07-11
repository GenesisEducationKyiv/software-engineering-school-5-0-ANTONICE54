package errors

import "weather-forecast/pkg/apperrors"

type InfrastructureErrorCode string

func (c InfrastructureErrorCode) String() string {
	return string(c)
}

const (
	DatabaseErrorCode InfrastructureErrorCode = "DATABASE_ERROR"
	InternalErrorCode InfrastructureErrorCode = "INTERNAL_ERROR"
)

var (
	DatabaseError = apperrors.NewInternal(DatabaseErrorCode, "database raised an error")
	InternalError = apperrors.NewInternal(InternalErrorCode, "internal server error")
)
