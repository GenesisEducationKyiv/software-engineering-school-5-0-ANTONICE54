package roundtrip

import (
	"bytes"
	"io"
	"net/http"
	"weather-forecast/pkg/logger"
)

type (
	LoggingRoundTripper struct {
		transport http.RoundTripper
		logger    logger.Logger
	}
)

func New(logger logger.Logger) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		transport: http.DefaultTransport,
		logger:    logger,
	}
}

func (rt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log := rt.logger.WithContext(req.Context())

	resp, err := rt.transport.RoundTrip(req)
	if err != nil {
		log.Warnf("Request to %s failed: %v", req.URL.String(), err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to read response body from %s: %v", req.URL.String(), err)
		return resp, err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Infof("%s - Response: %s", req.URL.Host, string(bodyBytes))

	return resp, nil

}
