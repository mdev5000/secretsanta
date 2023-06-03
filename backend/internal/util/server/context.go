package server

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

func Context(e echo.Context) context.Context {
	return e.Request().Context()
}

func SetContext(e echo.Context, ctx context.Context) {
	e.SetRequest(e.Request().WithContext(ctx))
}
