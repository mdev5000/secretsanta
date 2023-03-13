package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/util/log"
)

type AppHandler = func(ctx context.Context, c echo.Context) error

func APIBase(appCtx context.Context) echo.MiddlewareFunc {
	logger := log.Ctx(appCtx)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Eventually we might add tracing, but for now
			traceId := uuid.New()
			rq := c.Request()
			ctx := rq.Context()
			ctx = log.NewCtx(ctx, logger)
			ctx = attr.CtxPrefix(ctx,
				attr.String("traceId", traceId.String()))
			c.SetRequest(rq.WithContext(ctx))
			return next(c)
		}
	}
}

func Wrap(h AppHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		return h(ctx, c)
	}
}
