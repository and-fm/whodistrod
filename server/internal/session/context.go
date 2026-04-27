package session

import (
	"github.com/labstack/echo/v4"
)

func GetSession(c echo.Context) *SessionContext {
	sessionCtx := c.Get("session")

	session, _ := sessionCtx.(SessionContext)

	return &session
}
