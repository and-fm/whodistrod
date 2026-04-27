package middleware

import (
	"fmt"

	"github.com/and-fm/whodistrod/internal/session"
	"github.com/and-fm/whodistrod/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func SessionMiddleware(pg *pgxpool.Pool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionToken, err := c.Cookie("session")

			if err != nil {
				return c.NoContent(401)
			}

			sessionId := utils.HashSessionIdFromToken(string(sessionToken.Value))

			sessionContext := session.SessionContext{
				SessionId: sessionId,
				UserId:    0,
			}

			err = pg.QueryRow(utils.Ctb(), "select s.user_id from public.session s where s.id = $1", sessionId).Scan(&sessionContext.UserId)
			if err != nil {
				return echo.NewHTTPError(401, fmt.Errorf("SessionMiddleware - select public.session: %w", err).Error())
			}

			rows, err := pg.Query(utils.Ctb(), "select ur.role from public.user_role ur where ur.user_id = $1", sessionContext.UserId)
			if err != nil {
				return echo.NewHTTPError(401, fmt.Errorf("SessionMiddleware - select public.user_role: %w", err).Error())
			}

			roles, err := pgx.CollectRows(rows, pgx.RowTo[string])
			if err != nil {
				return echo.NewHTTPError(401, fmt.Errorf("SessionMiddleware - collect public.user_role: %w", err).Error())
			}

			sessionContext.Roles = roles

			c.Set("session", sessionContext)

			return next(c)
		}
	}
}
