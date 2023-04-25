package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/requests/gen/core"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/appjson"
	"github.com/mdev5000/secretsanta/internal/util/resp"
	"github.com/mdev5000/secretsanta/internal/util/session"
)

type Sessions interface {
	Put(ctx context.Context, key string, value interface{})
}

type UserHandler struct {
	sessions Sessions
	svc      *user.Service
}

func NewUserHandler(svc *user.Service, sessions Sessions) *UserHandler {
	return &UserHandler{svc: svc, sessions: sessions}
}

func (h *UserHandler) Login(ctx context.Context, c echo.Context) resp.ResponseEmpty {
	var data core.Login

	if err := appjson.UnmarshalJSON(c, &data); err != nil {
		// log error
		return resp.EmptyErr(apperror.Error(apperror.BadRequest, err))
	}

	u, err := h.svc.Login(ctx, data.Username, []byte(data.Password))
	if err != nil {
		// log error
		return resp.EmptyErr(apperror.Error(apperror.InvalidLogin, err, attr.String("user", u.ID)))
	}

	h.sessions.Put(ctx, session.UserID, u.ID)

	return resp.Empty(200)
}
