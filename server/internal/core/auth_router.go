package core

import (
	"github.com/and-fm/whodistrod/internal/databases"
	mw "github.com/and-fm/whodistrod/internal/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type authenticatedRouter struct {
	dig.In `ignore-unexported:"true"`

	echo *echo.Group // created in the New func

	Base BaseRouter
	Pg   databases.Postgres
}

type AuthenticatedRouter interface {
	Echo() *echo.Group
}

func NewAuthenticatedRouter(a authenticatedRouter) AuthenticatedRouter {
	g := a.Base.Echo().Group("")

	a.echo = g

	a.AddAuthMiddleware()

	return &a
}

func (b *authenticatedRouter) Echo() *echo.Group {
	return b.echo
}

func (a *authenticatedRouter) AddAuthMiddleware() {
	e := a.echo

	e.Use(mw.SessionMiddleware(a.Pg.GetPool()))
}
