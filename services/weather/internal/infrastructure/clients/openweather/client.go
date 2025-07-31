package openweather

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"weather-forecast/pkg/logger"
	"weather-service/internal/config"
	infraerrors "weather-service/internal/infrastructure/errors"
)

type (
	OpenWeatherErrorResponse struct {
		Cod     string `json:"cod"`
		Message string `json:"message"`
	}

	OpenWeatherMainResponse struct {
		Temperature float64 `json:"temp"`
		Humidity    int     `json:"humidity"`
	}

	OpenWeatherDescriptionResponse struct {
		Description string `json:"description"`
	}

	OpenWeatherSuccessResponse struct {
		Weather []OpenWeatherDescriptionResponse `json:"weather"`
		Main    OpenWeatherMainResponse          `json:"main"`
	}

	OpenWeatherClient struct {
		apiURL string
		apiKey string
		client *http.Client
		logger logger.Logger
	}
)

func NewClient(cfg *config.Config, client *http.Client, logger logger.Logger) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiURL: cfg.OpenWeatherURL,
		apiKey: cfg.OpenWeatherKey,
		client: client,
		logger: logger}
}

func (c *OpenWeatherClient) GetWeather(ctx context.Context, city string) (*OpenWeatherSuccessResponse, error) {
	const notFoundOpenWeatherErrorCode = "404"
	const metricUnits = "metric"

	url, err := url.Parse(c.apiURL)
	if err != nil {
		c.logger.Warnf("Form url: %s", err.Error())
		return nil, infraerrors.ErrGetWeather
	}
	queryString := url.Query()
	queryString.Set("q", city)
	queryString.Set("appid", c.apiKey)
	queryString.Set("units", metricUnits)
	url.RawQuery = queryString.Encode()
	stringURL := url.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, stringURL, nil)
	if err != nil {
		c.logger.Warnf("Failed to create get weather request: %s", err.Error())
		return nil, infraerrors.ErrGetWeather
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Warnf("Failed make get weather request: %s", err.Error())
		return nil, infraerrors.ErrGetWeather
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Warnf("Failed to close response body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Warnf("Failed to read response body: %s", err.Error())
		return nil, infraerrors.ErrGetWeather
	}

	if resp.StatusCode != http.StatusOK {
		var errResponse OpenWeatherErrorResponse

		if err := json.Unmarshal(body, &errResponse); err != nil {
			c.logger.Warnf("Failed to unmarshal response body: %s", err.Error())
			return nil, infraerrors.ErrGetWeather
		}

		if errResponse.Cod == notFoundOpenWeatherErrorCode {

			c.logger.Warnf("City not found: %s", city)
			return nil, infraerrors.ErrCityNotFound
		} else {
			c.logger.Warnf("Error from open weather: %s", errResponse.Message)
			return nil, infraerrors.ErrGetWeather
		}

	}

	var weatherResponse OpenWeatherSuccessResponse

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		c.logger.Warnf("Failed to unmarshal response body: %s", err.Error())
		return nil, infraerrors.ErrGetWeather
	}

	return &weatherResponse, nil

}
