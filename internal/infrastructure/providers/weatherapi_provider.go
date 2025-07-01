package providers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	GetWeatherConditionResponse struct {
		Text string `json:"text"`
	}

	GetWeatherCurrentResponse struct {
		TempC     float64                     `json:"temp_c"`
		Condition GetWeatherConditionResponse `json:"condition"`
		Humidity  int                         `json:"humidity"`
	}

	GetWeatherSuccessResponse struct {
		Current GetWeatherCurrentResponse `json:"current"`
	}

	GetWeatherErrorDetails struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	GetWeatherErrorResponse struct {
		Error GetWeatherErrorDetails `json:"error"`
	}

	WeatherAPIProvider struct {
		apiURL string
		apiKey string
		client *http.Client
		logger logger.Logger
	}
)

func NewWeatherAPIProvider(apiURL, apiKey string, httpClient *http.Client, logger logger.Logger) *WeatherAPIProvider {
	return &WeatherAPIProvider{
		apiURL: apiURL,
		apiKey: apiKey,
		client: httpClient,
		logger: logger,
	}
}

func (p *WeatherAPIProvider) GetWeatherByCity(ctx context.Context, city string) (*models.Weather, error) {

	const notFoundWeatherAPIErrorCode = 1006

	url, err := url.Parse(p.apiURL)
	if err != nil {
		p.logger.Warnf("Form url: %s", err.Error())
		return nil, apperrors.GetWeatherError
	}
	queryString := url.Query()
	queryString.Set("key", p.apiKey)
	queryString.Set("q", city)
	url.RawQuery = queryString.Encode()
	stringURL := url.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, stringURL, nil)
	if err != nil {
		p.logger.Warnf("Failed to create get weather request: %s", err.Error())
		return nil, apperrors.GetWeatherError
	}

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Warnf("Failed make get weather request: %s", err.Error())
		return nil, apperrors.GetWeatherError
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Warnf("Failed to close response body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.Warnf("Failed to read response body: %s", err.Error())
		return nil, apperrors.GetWeatherError
	}

	if resp.StatusCode != http.StatusOK {

		var errResponse GetWeatherErrorResponse

		if err := json.Unmarshal(body, &errResponse); err != nil {
			p.logger.Warnf("Failed to unmarshal response body: %s", err.Error())
			return nil, apperrors.GetWeatherError
		}

		if errResponse.Error.Code == notFoundWeatherAPIErrorCode {
			p.logger.Warnf("City not found: %s", city)
			return nil, apperrors.CityNotFoundError
		} else {
			p.logger.Warnf("Error from weather api: %s", errResponse.Error.Message)
			return nil, apperrors.GetWeatherError
		}

	}

	var weather GetWeatherSuccessResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		p.logger.Warnf("Failed to unmarshal response body: %s", err.Error())
		return nil, apperrors.GetWeatherError
	}

	result := models.Weather{
		Temperature: weather.Current.TempC,
		Humidity:    weather.Current.Humidity,
		Description: weather.Current.Condition.Text,
	}

	return &result, nil
}
