package cookie

import (
	"net/http"
	"os"
)

var domain = ""

func init() {
	// @todo review this
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
