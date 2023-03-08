package middleware

import "github.com/labstack/echo/v4"

func Dev() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Support receiving the request from a different port, since the frontend dev server will run on a
			// separate port.
			c.Response().Header().Add("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	}
}
