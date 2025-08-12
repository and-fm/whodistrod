package logging

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/and-fm/whodistrod/internal/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type logger struct {
	dig.In `ignore-unexported:"true"`

	Config *config.Config

	logger *slog.Logger // created in the New func
}

type Logger interface {
	LogDebug(msg string, args ...any)
	LogInfo(msg string, ctx echo.Context, args ...any)
	LogWarning(msg string, ctx echo.Context, args ...any)
	LogError(msg string, ctx echo.Context, args ...any)
	LogFatal(msg string, ctx echo.Context, args ...any)
	Log(level slog.Level, msg string, logModel *LogModel)
}

func NewLogger(l logger) Logger {
	l.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &l
}

func (l *logger) LogDebug(msg string, args ...any) {
	logModel := &LogModel{
		AppName: l.Config.AppName,
		Env:     l.Config.Env,
	}

	l.Log(slog.LevelDebug, fmt.Sprintf(msg, args...), logModel)
}

func (l *logger) LogInfo(msg string, ctx echo.Context, args ...any) {
	logModel := &LogModel{
		AppName: l.Config.AppName,
		Env:     l.Config.Env,
	}

	if ctx != nil {
		logModel.Method = ctx.Request().Method
		logModel.URI = ctx.Request().RequestURI
	}

	l.Log(slog.LevelInfo, fmt.Sprintf(msg, args...), logModel)
}

func (l *logger) LogWarning(msg string, ctx echo.Context, args ...any) {
	logModel := &LogModel{
		AppName: l.Config.AppName,
		Env:     l.Config.Env,
	}

	if ctx != nil {
		logModel.Method = ctx.Request().Method
		logModel.URI = ctx.Request().RequestURI
	}

	l.Log(slog.LevelWarn, fmt.Sprintf(msg, args...), logModel)
}

func (l *logger) LogError(msg string, ctx echo.Context, args ...any) {
	logModel := &LogModel{
		AppName: l.Config.AppName,
		Env:     l.Config.Env,
		Error:   "error",
	}

	if ctx != nil {
		logModel.Method = ctx.Request().Method
		logModel.URI = ctx.Request().RequestURI
	}

	l.Log(slog.LevelError, fmt.Sprintf(msg, args...), logModel)
}

func (l *logger) LogFatal(msg string, ctx echo.Context, args ...any) {
	l.LogError(msg, ctx, args...)
	log.Fatalf(msg, args...)
}

func (l *logger) Log(level slog.Level, msg string, logModel *LogModel) {
	l.logger.LogAttrs(context.Background(), level, msg,
		slog.String("appName", logModel.AppName),
		slog.String("env", logModel.Env),
		slog.String("method", logModel.Method),
		slog.String("uri", logModel.URI),
		slog.String("ip", logModel.RemoteIP),
		slog.Int("status", logModel.Status),
		slog.String("latency", logModel.Latency),
		slog.String("error", logModel.Error),
		slog.String("sessionId", logModel.SessionId),
		slog.Int("userId", logModel.UserId),
		slog.String("userEmail", logModel.UserEmail),
		slog.String("username", logModel.Username),
	)
}
