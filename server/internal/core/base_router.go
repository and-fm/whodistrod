package core

import (
	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/logging"
	mw "github.com/and-fm/whodistrod/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

type baseRouter struct {
	dig.In `ignore-unexported:"true"`

	echo *echo.Echo

	Config                   *config.Config
	Logger                   logging.Logger
	RequestLoggingMiddleware mw.RequestLoggingMiddleware
}

type BaseRouter interface {
	BaseEcho() *echo.Echo
}

func NewBaseRouter(b baseRouter) BaseRouter {
	e := echo.New()

	b.echo = e

	b.AddBaseMiddleware()

	return &b
}

func (b *baseRouter) BaseEcho() *echo.Echo {
	return b.echo
}

func (b *baseRouter) AddBaseMiddleware() {
	e := b.echo

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Recover())

	e.Use(b.RequestLoggingMiddleware.MwFunc())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}
