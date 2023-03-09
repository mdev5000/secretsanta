package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/requests/gen"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/requests"
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

	var data gen.Login

	if err := requests.UnmarshalJSON(c, &data); err != nil {
		// log error
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	u, err := h.svc.Login(ctx, data.Username, []byte(data.Password))
	if err != nil {
		// log error
		return echo.NewHTTPError(http.StatusBadRequest, "invalid credentials")
	}

	c.SetCookie(&http.Cookie{
		Name:  "user.id",
		Value: u.ID.String(),
	})

	return nil
}
