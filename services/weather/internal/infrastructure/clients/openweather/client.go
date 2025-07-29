package openweather

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"weather-forecast/pkg/logger"
	"weather-service/internal/config"
	"weather-service/internal/infrastructure/errors"
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
	log := c.logger.WithContext(ctx)

	log.Infof("Calling OpenWeather API for city: %s", city)

	const notFoundOpenWeatherErrorCode = "404"
	const metricUnits = "metric"

	url, err := url.Parse(c.apiURL)
	if err != nil {
		log.Warnf("Form url: %s", err.Error())
		return nil, errors.GetWeatherError
	}
	queryString := url.Query()
	queryString.Set("q", city)
	queryString.Set("appid", c.apiKey)
	queryString.Set("units", metricUnits)
	url.RawQuery = queryString.Encode()
	stringURL := url.String()

	log.Debugf("Making request to OpenWeather: %s", stringURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, stringURL, nil)
	if err != nil {
		log.Warnf("Failed to create get weather request: %s", err.Error())
		return nil, errors.GetWeatherError
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Warnf("Failed make get weather request: %s", err.Error())
		return nil, errors.GetWeatherError
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("Failed to close response body: %s", err.Error())
		}
	}()

	log.Debugf("OpenWeather API responded with status: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to read response body: %s", err.Error())
		return nil, errors.GetWeatherError
	}

	if resp.StatusCode != http.StatusOK {
		var errResponse OpenWeatherErrorResponse

		if err := json.Unmarshal(body, &errResponse); err != nil {
			log.Warnf("Failed to unmarshal response body: %s", err.Error())
			return nil, errors.GetWeatherError
		}

		if errResponse.Cod == notFoundOpenWeatherErrorCode {

			log.Warnf("City not found: %s", city)
			return nil, errors.CityNotFoundError
		} else {
			log.Warnf("Error from open weather: %s", errResponse.Message)
			return nil, errors.GetWeatherError
		}

	}

	var weatherResponse OpenWeatherSuccessResponse

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		log.Warnf("Failed to unmarshal response body: %s", err.Error())
		return nil, errors.GetWeatherError
	}

	log.Infof("Successfully received weather from OpenWeather for city: %s", city)

	return &weatherResponse, nil

}
