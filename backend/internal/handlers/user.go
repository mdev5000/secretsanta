package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/resp"
)

type UserHandler struct {
	svc *user.Service
}

func (h *SetupHandler) FindAll(ctx context.Context, c echo.Context) resp.ResponseEmpty {
	return resp.Empty(200)
}
