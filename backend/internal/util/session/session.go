package session

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"net/http"
	"time"
)

const UserID = "UserID"

type Manager = scs.SessionManager

func New(db *mongo.Database, isDev bool) *Manager {
	sessionStore := scs.New()
	sessionStore.Store = mongodbstore.New(db)
	sessionStore.Cookie.Path = "/"
	if isDev {
		sessionStore.Cookie.SameSite = http.SameSiteLaxMode
		sessionStore.Cookie.Secure = false
	}
	return sessionStore
}

func Get[T any](ctx context.Context, sm *Manager, key string) (T, error) {
	v := sm.Get(ctx, key)
	vv, ok := v.(T)
	if !ok {
		if v == nil {
			return vv, nil
		}
		return vv, fmt.Errorf("invalid session key type for '%s' exptected %T was %T", key, vv, v)
	}
	return vv, nil
}

func Put[T any](ctx context.Context, sm *Manager, key string, val T) {
	sm.Put(ctx, key, val)
}

func Destroy(ctx context.Context, sm *Manager) error {
	return sm.Destroy(ctx)
}

func IsLoggedIn(ctx context.Context, sm *Manager) (bool, error) {
	userId, err := UserId(ctx, sm)
	return userId != "", err
}

func UserId(ctx context.Context, sm *Manager) (string, error) {
	return Get[string](ctx, sm, UserID)
}

func TrySaveSession(ctx context.Context, sm *Manager, e echo.Context) error {
	w := e.Response()

	switch sm.Status(ctx) {
	case scs.Modified:
		token, expiry, err := sm.Commit(ctx)
		if err != nil {
			return err
		}
		sm.WriteSessionCookie(ctx, w, token, expiry)
	case scs.Destroyed:
		sm.WriteSessionCookie(ctx, w, "", time.Time{})
	}

	w.Header().Add("Vary", "Cookie")
	return nil
}
