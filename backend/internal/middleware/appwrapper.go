package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/util/log"
)

type AppHandler = func(ctx context.Context, c echo.Context) error

func Wrap(serverCtx context.Context, h AppHandler) echo.HandlerFunc {
	logger := log.Ctx(serverCtx)
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ctx = log.NewCtx(ctx, logger)
		return h(ctx, c)
	}
}
