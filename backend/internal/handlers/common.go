package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// ApiOptions creates an options handler for a given set of methods.
func ApiOptions(methods ...string) echo.HandlerFunc {
	allow := strings.Join(methods, ",")
	return func(c echo.Context) error {
		headers := c.Response().Header()
		headers.Add("Allow", allow)
		headers.Add("Accept", "application/json")
		headers.Add("Access-Control-Request-Method", allow)
		return c.Blob(http.StatusNoContent, "application/json", nil)
	}
}
