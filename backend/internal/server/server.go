package server

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Config struct {
}

func Server(config Config) *echo.Echo {
	e := echo.New()

	e.Static("assets", "../frontend/assets")

	e.GET("app", homePage)
	e.GET("app/*", homePage)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return e
}

func homePage(c echo.Context) error {
	b, err := os.ReadFile("../frontend/build/index.html")
	if err != nil {
		return err
	}
	return c.Blob(200, "text/html", b)
}
