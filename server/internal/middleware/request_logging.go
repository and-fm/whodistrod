package middleware

import (
	"log/slog"

	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

type requestLoggingMiddleware struct {
	dig.In `ignore-unexported:"true"`

	Config *config.Config
	Logger logging.Logger

	mwFunc echo.MiddlewareFunc // created in New func
}

type RequestLoggingMiddleware interface {
	MwFunc() echo.MiddlewareFunc
}

func NewRequestLoggingMiddleware(r requestLoggingMiddleware) RequestLoggingMiddleware {
	r.mwFunc = middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogLatency:  true,
		LogError:    true,
		LogMethod:   true,
		LogRemoteIP: true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logModel := &logging.LogModel{
				AppName:   r.Config.AppName,
				Env:       r.Config.Env,
				Method:    v.Method,
				URI:       v.URI,
				RemoteIP:  v.RemoteIP,
				Status:    v.Status,
				Latency:   v.Latency.String(),
			}

			if v.Error != nil {
				logModel.Error = v.Error.Error()
				r.Logger.Log(slog.LevelError, "REQUEST_ERROR", logModel)
			} else {
				r.Logger.Log(slog.LevelInfo, "REQUEST", logModel)
			}

			return nil
		},
	})

	return &r
}

func (r *requestLoggingMiddleware) MwFunc() echo.MiddlewareFunc {
	return r.mwFunc
}
