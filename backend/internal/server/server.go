package server

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"sync"
)

type Config struct {
	SpaContent embed.FS
}

var readIndexFile = sync.Once{}
var indexFile []byte

func Server(config *Config) *echo.Echo {
	e := echo.New()

	contentHandler := echo.WrapHandler(http.FileServer(http.FS(config.SpaContent)))
	// The embedded files will all be in the '/static' folder so need to rewrite the request (could also do this with fs.Sub)
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/embedded/$1"})
	e.GET("assets/*", contentHandler, contentRewrite)

	homePage := mkHomePageHandler(config)
	e.GET("app", homePage)
	e.GET("app/*", homePage)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/app")
	})

	return e
}

func mkHomePageHandler(config *Config) echo.HandlerFunc {
	spaContent := config.SpaContent
	return func(c echo.Context) error {
		var err error
		readIndexFile.Do(func() {
			b, readErr := spaContent.ReadFile("embedded/spa/index.html")
			if readErr != nil {
				err = readErr
				return
			}
			indexFile = b
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return c.Blob(200, "text/html", indexFile)
	}
}
