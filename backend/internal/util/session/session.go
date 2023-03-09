package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/user"
)

const (
	KeyUserID = "UserID"
)

type CookieStore = sessions.CookieStore

type Session struct {
	sess *sessions.Session
}

func (s Session) SetUserID(userID user.ID) {
	s.Set(KeyUserID, userID)
}

func (s Session) UserID() (user.ID, error) {
	id, ok := s.Get(KeyUserID).(user.ID)
	if !ok {
		return id, fmt.Errorf("invalid userID '%t'", s.Get(KeyUserID))
	}
	return id, nil
}

func (s Session) Set(k interface{}, v interface{}) {
	s.sess.Values[k] = v
}

func (s Session) Get(k interface{}) interface{} {
	return s.sess.Values[k]
}

func (s Session) Save(c echo.Context) error {
	return s.sess.Save(c.Request(), c.Response())
}

func Store(secret []byte) *CookieStore {
	return sessions.NewCookieStore(secret)
}

func Get(c echo.Context) (Session, error) {
	s, err := session.Get("session", c)
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	if err != nil {
		return Session{}, fmt.Errorf("failed to get session: %w", err)
	}
	return Session{sess: s}, nil
}
