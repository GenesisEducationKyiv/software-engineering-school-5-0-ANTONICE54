package errors

import (
	"weather-forecast/pkg/apperrors"
)

type InfrastructureErrorCode string

func (c InfrastructureErrorCode) String() string {
	return string(c)
}

const (
	DatabaseErrorCode     InfrastructureErrorCode = "DATABASE_ERROR"
	GetWeatherErrorCode   InfrastructureErrorCode = "GET_WEATHER_ERROR"
	CacheErrorCode        InfrastructureErrorCode = "CACHE_ERROR"
	CacheMissErrorCode    InfrastructureErrorCode = "CACHE_MISS_ERROR"
	CityNotFoundErrorCode InfrastructureErrorCode = "CITY_NOT_FOUND_ERROR"
	InternalErrorCode     InfrastructureErrorCode = "INTERNAL_ERROR"
	InvalidTokenErrorCode InfrastructureErrorCode = "INVALID_TOKEN_ERROR"
)

var (
	DatabaseError     = apperrors.NewInternal(DatabaseErrorCode, "database raised an error")
	GetWeatherError   = apperrors.NewInternal(GetWeatherErrorCode, "failed to get weather")
	CacheError        = apperrors.NewInternal(CacheErrorCode, "failed to interact with cache")
	CacheMissError    = apperrors.NewNotFound(CacheMissErrorCode, "cache miss")
	CityNotFoundError = apperrors.NewNotFound(CityNotFoundErrorCode, "there is no city with such name")
	InternalError     = apperrors.NewInternal(InternalErrorCode, "internal server error")
	InvalidTokenError = apperrors.NewBadRequest(InvalidTokenErrorCode, "invalid token")
)
