package middleware

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/server"
	"github.com/mdev5000/secretsanta/internal/util/session"
)

func EnsureLoggedIn(ctx context.Context, sm *session.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, err := session.UserId(server.Context(c), sm)
			if err != nil {
				return err
			}
			if userId == "" {
				// not logged in
				return apperror.Error(apperror.NotAuthenticated, errors.New("not authenticated"))
			}
			return next(c)
		}
	}
}
