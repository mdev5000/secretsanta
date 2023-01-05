package server

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/handlers"
	mw "github.com/mdev5000/secretsanta/internal/middleware"
	"net/http"
	"sync"
)

type Config struct {
}

var readIndexFile = sync.Once{}
var indexFile []byte

func Server(ac *appcontext.AppContext, config *Config) *echo.Echo {

	e := echo.New()

	setupHandler := handlers.NewSetupHandler(ac.SetupService)
	e.POST("/setup", setupHandler.FinalizeSetup)

	// Set up the assets file server.
	contentHandler := echo.WrapHandler(http.FileServer(http.FS(ac.SPAContent)))
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/embedded/$1"})
	e.GET("assets/*", contentHandler, contentRewrite)

	appGroup := e.Group("/")

	if !ac.SetupService.IsSetup() {
		appGroup.Use(mw.IsSetup(ac.SetupService))
	}

	homePage := mkHomePageHandler(ac.SPAContent)
	appGroup.GET("app", homePage)
	appGroup.GET("app/*", homePage)

	appGroup.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/app")
	})

	return e
}

func mkHomePageHandler(spaContent embed.FS) echo.HandlerFunc {
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
