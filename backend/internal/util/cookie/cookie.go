package cookie

import (
	"net/http"
	"os"
)

var domain = ""

func init() {
	domain = os.Getenv("DOMAIN")
}

type Cookie = http.Cookie

func MakeCookie(c Cookie) *http.Cookie {
	if c.Domain == "" {
		c.Domain = domain
	}
	if c.Path == "" {
		c.Path = "/"
	}
	return &c
}

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
