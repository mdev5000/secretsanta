package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/requests/gen/core"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/appjson"
	"github.com/mdev5000/secretsanta/internal/util/cookie"
	"github.com/mdev5000/secretsanta/internal/util/log"
	"github.com/mdev5000/secretsanta/internal/util/resp"
	"github.com/mdev5000/secretsanta/internal/util/session"
	"net/http"
)

type AuthHandler struct {
	sessions *session.Manager
	svc      *user.Service
}

func NewAuthHandler(svc *user.Service, sessions *session.Manager) *AuthHandler {
	return &AuthHandler{svc: svc, sessions: sessions}
}

func (h *AuthHandler) Login(ctx context.Context, c echo.Context) resp.ResponseEmpty {
	var data core.Login

	log.Ctx(ctx).Debug("authenticating user")

	userId, err := session.UserId(ctx, h.sessions)
	if err != nil {
		return resp.EmptyErr(apperror.Error(apperror.BadRequest, err))
	}
	if userId != "" {
		return resp.EmptyErr(apperror.Error(apperror.AlreadyLoggedIn, err))
	}

	if err := appjson.UnmarshalJSON(c, &data); err != nil {
		// log error
		return resp.EmptyErr(apperror.Error(apperror.BadRequest, err))
	}

	u, err := h.svc.Login(ctx, data.Username, []byte(data.Password))
	if err != nil {
		// log error
		return resp.EmptyErr(apperror.Error(apperror.InvalidLogin, err, attr.String("user", u.ID)))
	}

	session.Put[string](ctx, h.sessions, session.UserID, u.ID)
	c.SetCookie(loginCookie(ctx))

	return resp.Empty(200)
}

func loginCookie(ctx context.Context) *http.Cookie {
	loggedIn := cookie.NewCookie(ctx)
	loggedIn.Name = "loggedIn"
	loggedIn.Value = "true"
	return loggedIn
}

func (h *AuthHandler) Logout(ctx context.Context, e echo.Context) resp.ResponseEmpty {
	loginCook := loginCookie(ctx)
	loginCook.Value = ""
	loginCook.MaxAge = 0
	e.SetCookie(loginCook)
	err := session.Destroy(ctx, h.sessions)
	if err != nil {
		return resp.EmptyErr(err)
	}
	return resp.Empty(200)
}
