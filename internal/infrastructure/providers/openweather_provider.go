package providers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"weather-forecast/internal/domain/models"
	infraerror "weather-forecast/internal/infrastructure/errors"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	GetOpenWeatherErrorResponse struct {
		Cod     string `json:"cod"`
		Message string `json:"message"`
	}

	GetOpenWeatherMainResponse struct {
		Temperature float64 `json:"temp"`
		Humidity    int     `json:"humidity"`
	}

	GetOpenWeatherDescriptionResponse struct {
		Description string `json:"description"`
	}

	GetOpenWeatherSuccessResponse struct {
		Weather []GetOpenWeatherDescriptionResponse `json:"weather"`
		Main    GetOpenWeatherMainResponse          `json:"main"`
	}

	OpenWeatherProvider struct {
		apiURL string
		apiKey string
		client *http.Client
		logger logger.Logger
	}
)

func NewOpenWeatherProvider(apiURL, apiKey string, client *http.Client, logger logger.Logger) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		apiURL: apiURL,
		apiKey: apiKey,
		client: client,
		logger: logger,
	}
}

func (p *OpenWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	const notFoundOpenWeatherErrorCode = "404"
	const metricUnits = "metric"

	url, err := url.Parse(p.apiURL)
	if err != nil {
		p.logger.Warnf("Form url: %s", err.Error())
		return nil, infraerror.ErrGetWeather
	}
	queryString := url.Query()
	queryString.Set("q", city)
	queryString.Set("appid", p.apiKey)
	queryString.Set("units", metricUnits)
	url.RawQuery = queryString.Encode()
	stringURL := url.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, stringURL, nil)
	if err != nil {
		p.logger.Warnf("Failed to create get weather request: %s", err.Error())
		return nil, infraerror.ErrGetWeather
	}

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Warnf("Failed make get weather request: %s", err.Error())
		return nil, infraerror.ErrGetWeather
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Warnf("Failed to close response body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.Warnf("Failed to read response body: %s", err.Error())
		return nil, infraerror.ErrGetWeather
	}

	if resp.StatusCode != http.StatusOK {
		var errResponse GetOpenWeatherErrorResponse

		if err := json.Unmarshal(body, &errResponse); err != nil {
			p.logger.Warnf("Failed to unmarshal response body: %s", err.Error())
			return nil, infraerror.ErrGetWeather
		}

		if errResponse.Cod == notFoundOpenWeatherErrorCode {

			p.logger.Warnf("City not found: %s", city)
			return nil, infraerror.ErrCityNotFound
		} else {
			p.logger.Warnf("Error from open weather: %s", errResponse.Message)
			return nil, infraerror.ErrGetWeather
		}

	}

	var weatherResponse GetOpenWeatherSuccessResponse

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		p.logger.Warnf("Failed to unmarshal response body: %s", err.Error())
		return nil, infraerror.ErrGetWeather
	}

	weatherDesc := ""

	if len(weatherResponse.Weather) > 0 {
		weatherDesc = weatherResponse.Weather[0].Description
	} else {
		p.logger.Warnf("OpenWeather did not provide weather description for city: %s", city)
	}

	result := models.Weather{
		Temperature: weatherResponse.Main.Temperature,
		Humidity:    weatherResponse.Main.Humidity,
		Description: weatherDesc,
	}

	return &result, nil

}
