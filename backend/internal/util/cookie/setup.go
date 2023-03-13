package cookie

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SiteSetupCookie(ctx context.Context, isSetup bool) *http.Cookie {
	value := "false"
	if isSetup {
		value = "true"
	}
	c := NewCookie(ctx)
	c.Name = "site.isSetup"
	c.Value = value
	return c
}

const setupLeaderCookieName = "setup.uuid"

func GetSetupLeaderCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(setupLeaderCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetupLeaderCookie(ctx context.Context, uuid string) *http.Cookie {
	c := NewCookie(ctx)
	c.Name = setupLeaderCookieName
	c.Value = uuid
	c.HttpOnly = true
	return c
}
