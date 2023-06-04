package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	mw "github.com/mdev5000/secretsanta/internal/middleware"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/resp"
	"github.com/mdev5000/secretsanta/internal/util/session"
	"google.golang.org/protobuf/proto"
)

// wrapAPI wraps an API request. It enforces request contracts at startup and attempts to save session data at the end
// of the request.
func wrapAPI[T proto.Message](s *server, h func(context.Context, echo.Context) resp.Response[T]) echo.HandlerFunc {
	if err := mw.EnsureComplies[T](); err != nil {
		panic(err)
	}
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		rs := h(ctx, c)
		if err := session.TrySaveSession(ctx, s.sessionMgr, c); err != nil {
			return apperror.InternalError(fmt.Errorf("failed to save session data: %w", err))
		}
		if rs.Err != nil {
			return rs.Err
		}
		if rs.Data == nil {
			panic("either Err or Data must be setup on a response")
		}
		return c.JSONBlob(rs.Code, rs.Data)
	}
}
