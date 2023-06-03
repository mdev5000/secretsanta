package cookie

import (
	"context"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"net/http"
	"time"
)

type Cookie = http.Cookie

func NewCookie(ctx context.Context) *http.Cookie {
	c := http.Cookie{}
	c.Expires = time.Now().Add(10 * 365 * 24 * time.Hour)
	c.Path = "/"

	isDev, _ := ctx.Value(appcontext.KeyIsDev).(bool)
	if isDev {
		// @todo figure out how to separate dev cookies
		c.SameSite = http.SameSiteLaxMode
	}

	return &c
}
