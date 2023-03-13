package cookie

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SiteSetupCookie(isSetup bool) *http.Cookie {
	value := "false"
	if isSetup {
		value = "true"
	}
	return MakeCookie(Cookie{
		Name:  "site.isSetup",
		Value: value,
	})
}

const setupLeaderCookieName = "setup.uuid"

func GetSetupLeaderCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(setupLeaderCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetupLeaderCookie(uuid string) *http.Cookie {
	return MakeCookie(Cookie{
		Name:  setupLeaderCookieName,
		Value: uuid,
	})
}
