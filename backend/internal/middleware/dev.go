package middleware

import (
	"github.com/labstack/echo/v4"
)

func ApiDev() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Support receiving the request from a different port, since the frontend dev server will run on a
			// separate port.
			headers := c.Response().Header()
			headers.Add("Access-Control-Allow-Origin", "http://localhost:5173")
			headers.Add("Access-Control-Allow-Credentials", "true")
			return next(c)
		}
	}
}
