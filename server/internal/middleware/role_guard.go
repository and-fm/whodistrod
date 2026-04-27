package middleware

import (
	"net/http"
	"slices"

	"github.com/and-fm/whodistrod/internal/session"
	"github.com/labstack/echo/v4"
)

func RoleGuard(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := session.GetSession(c)

			if slices.Contains(session.Roles, role) {
				return next(c)
			} else {
				c.Error(c.NoContent(
					http.StatusForbidden,
				))
			}
			return nil
		}
	}
}
