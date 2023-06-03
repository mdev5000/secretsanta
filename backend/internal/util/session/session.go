package session

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/mdev5000/secretsanta/internal/mongo"
)

const UserID = "UserID"

type Manager = scs.SessionManager

func New(db *mongo.Database) *Manager {
	sessionStore := scs.New()
	sessionStore.Store = mongodbstore.New(db)
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

func UserId(ctx context.Context, sm *Manager) (int, error) {
	return Get[int](ctx, sm, UserID)
}
