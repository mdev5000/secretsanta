package appctx

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/util/log"
	fz "github.com/mdev5000/secretsanta/internal/util/log/flog-zero"
)

func Init(c echo.Context) context.Context {
	ctx := c.Request().Context()
	ctx = log.NewCtx(ctx, fz.New(os.Stdout))
	return ctx
}
