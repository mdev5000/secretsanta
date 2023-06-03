package middleware

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/util/server"
)

func LoadSessionData(s *scs.SessionManager) echo.MiddlewareFunc {
	// LoadAndSave does not work correctly with echo so we can't use it.
	//return echo.WrapMiddleware(sessionMgr.LoadAndSave)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var token string
			cookie, err := c.Cookie(s.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err := s.Load(server.Context(c), token)
			if err != nil {
				return err
			}

			server.SetContext(c, ctx)

			if err := next(c); err != nil {
				return err
			}

			// Not sure what this does it was in LoadAndSave implementation.
			if c.Request().MultipartForm != nil {
				c.Request().MultipartForm.RemoveAll()
			}

			return nil
		}
	}
}
