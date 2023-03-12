package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/util/appctx"
)

type AppHandler = func(ctx context.Context, c echo.Context) error

func Wrap(h AppHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := appctx.Init(c)
		return h(ctx, c)
	}
}
