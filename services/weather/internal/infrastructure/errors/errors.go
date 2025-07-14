package errors

import "weather-forecast/pkg/apperrors"

type WeatherServiceErrorCode string

func (c WeatherServiceErrorCode) String() string {
	return string(c)
}

const (
	GetWeatherErrorCode   WeatherServiceErrorCode = "GET_WEATHER_ERROR"
	CacheMissErrorCode    WeatherServiceErrorCode = "CACHE_MISS_ERROR"
	CityNotFoundErrorCode WeatherServiceErrorCode = "CITY_NOT_FOUND_ERROR"
	CacheErrorCode        WeatherServiceErrorCode = "CACHE_ERROR"
	InternalErrorCode     WeatherServiceErrorCode = "INTERNAL_SERVER_ERROR"
)

var (
	GetWeatherError     = apperrors.NewInternal(GetWeatherErrorCode, "failed to get weather")
	CacheError          = apperrors.NewInternal(CacheErrorCode, "failed to interact with cache")
	CacheMissError      = apperrors.NewNotFound(CacheMissErrorCode, "cache miss")
	InternalServerError = apperrors.NewInternal(InternalErrorCode, "internal server error")
	CityNotFoundError   = apperrors.NewNotFound(CityNotFoundErrorCode, "there is no city with such name")
)
