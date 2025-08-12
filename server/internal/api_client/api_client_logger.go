package apiclient

import (
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/go-resty/resty/v2"
)

type apiClientLogger struct {
	logger logging.Logger
}

type ApiClientLogger interface {
	resty.Logger
}

func NewApiClientLogger(logger logging.Logger) ApiClientLogger {
	return &apiClientLogger{logger: logger}
}

func (a *apiClientLogger) Errorf(format string, v ...interface{}) {
	a.logger.LogError(format, nil, v)
}

func (a *apiClientLogger) Warnf(format string, v ...interface{}) {
	a.logger.LogWarning(format, nil, v)
}

func (a *apiClientLogger) Debugf(format string, v ...interface{}) {
	a.logger.LogDebug(format, v)
}
