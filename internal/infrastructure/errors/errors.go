package errors

import "errors"

var (
	ErrGetWeather   = errors.New("failed to get weather")
	ErrCityNotFound = errors.New("there is no city with such name")
	ErrDatabase     = errors.New("database raised an error")
	ErrCache        = errors.New("failed to interact with cache")
	ErrCacheMiss    = errors.New("cache miss")
	ErrInternal     = errors.New("internal server error")
)
