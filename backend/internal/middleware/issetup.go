package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/setup"
	"net/http"
)

func IsSetup(svc *setup.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// @todo eventually remove this
			if path == "/example" || path == "/app/example" {
				return next(c)
			}

			isSetupPath := path == "/app/setup"
			isSetup, err := svc.IsSetup(c.Request().Context())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			setSetupCookie(c, isSetup)

			if isSetupPath {
				if isSetup {
					return c.Redirect(http.StatusTemporaryRedirect, "/app")
				}
				return next(c)
			}

			if isSetup {
				return next(c)
			}

			return c.Redirect(http.StatusTemporaryRedirect, "/app/setup")
		}
	}
}

func setSetupCookie(c echo.Context, isSetup bool) {
	value := "false"
	if isSetup {
		value = "true"
	}
	c.SetCookie(&http.Cookie{
		Name:  "site.isSetup",
		Value: value,
		Path:  "/",
	})
}
