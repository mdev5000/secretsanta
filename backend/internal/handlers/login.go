package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/requests/login"
	"github.com/mdev5000/secretsanta/internal/user"
	"net/http"
)

type UserHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var data login.Login
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("bad data: %w", err))
	}

	u, err := h.svc.Login(ctx, data.Username, []byte(data.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid credentials")
	}

	c.SetCookie(&http.Cookie{
		Name:  "user.id",
		Value: u.ID.String(),
	})

	return nil
}
