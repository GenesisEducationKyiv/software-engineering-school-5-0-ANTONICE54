package weatherapi

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
	WeatherConditionResponse struct {
		Text string `json:"text"`
	}

	WeatherCurrentResponse struct {
		TempC     float64                  `json:"temp_c"`
		Condition WeatherConditionResponse `json:"condition"`
		Humidity  int                      `json:"humidity"`
	}

	WeatherSuccessResponse struct {
		Current WeatherCurrentResponse `json:"current"`
	}

	WeatherErrorDetails struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	WeatherErrorResponse struct {
		Error WeatherErrorDetails `json:"error"`
	}

	WeatherAPIClient struct {
		apiURL string
		apiKey string
		client *http.Client
		logger logger.Logger
	}
)

func NewClient(cfg *config.Config, httpClient *http.Client, logger logger.Logger) *WeatherAPIClient {
	return &WeatherAPIClient{
		apiURL: cfg.WeatherAPIURL,
		apiKey: cfg.WeatherAPIKey,
		client: httpClient,
		logger: logger,
	}
}

func (c *WeatherAPIClient) GetWeather(ctx context.Context, city string) (*WeatherSuccessResponse, error) {
	log := c.logger.WithContext(ctx)

	log.Infof("Calling WeatherAPI for city: %s", city)

	const notFoundWeatherAPIErrorCode = 1006

	url, err := url.Parse(c.apiURL)
	if err != nil {
		log.Warnf("Form url: %s", err.Error())
		return nil, errors.GetWeatherError
	}
	queryString := url.Query()
	queryString.Set("key", c.apiKey)
	queryString.Set("q", city)
	url.RawQuery = queryString.Encode()
	stringURL := url.String()

	log.Debugf("Making request to WeatherAPI: %s", stringURL)

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

	log.Debugf("WeatherAPI responded with status: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to read response body: %s", err.Error())
		return nil, errors.GetWeatherError
	}

	if resp.StatusCode != http.StatusOK {

		var errResponse WeatherErrorResponse

		if err := json.Unmarshal(body, &errResponse); err != nil {
			log.Warnf("Failed to unmarshal response body: %s", err.Error())
			return nil, errors.GetWeatherError
		}

		if errResponse.Error.Code == notFoundWeatherAPIErrorCode {
			log.Warnf("City not found: %s", city)
			return nil, errors.CityNotFoundError
		} else {
			log.Warnf("Error from weather api: %s", errResponse.Error.Message)
			return nil, errors.GetWeatherError
		}

	}

	var weather WeatherSuccessResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		log.Warnf("Failed to unmarshal response body: %s", err.Error())
		return nil, errors.GetWeatherError
	}

	log.Infof("Successfully received weather from WeatherAPI for city: %s", city)

	return &weather, nil
}
