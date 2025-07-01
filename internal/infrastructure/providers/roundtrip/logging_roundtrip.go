package roundtrip

import (
	"bytes"
	"io"
	"net/http"
	"weather-forecast/internal/infrastructure/logger"
)

type (
	LoggingRoundTripper struct {
		transport     http.RoundTripper
		fileLogger    logger.Logger
		consoleLogger logger.Logger
	}
)

func New(fileLog, consoleLog logger.Logger) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		transport:     http.DefaultTransport,
		fileLogger:    fileLog,
		consoleLogger: consoleLog,
	}
}

func (rt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	resp, err := rt.transport.RoundTrip(req)
	if err != nil {
		rt.consoleLogger.Warnf("Request to %s failed: %v", req.URL.String(), err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		rt.consoleLogger.Warnf("Failed to read response body from %s: %v", req.URL.String(), err)
		return resp, err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	rt.fileLogger.Infof("%s - Response: %s", req.URL.Host, string(bodyBytes))

	return resp, nil

}
