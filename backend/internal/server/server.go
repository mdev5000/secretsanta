package server

import (
	"context"
	"embed"
	"fmt"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/util/log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/handlers"
	mw "github.com/mdev5000/secretsanta/internal/middleware"
	"github.com/mdev5000/secretsanta/internal/requests/gen"
	"github.com/mdev5000/secretsanta/internal/util/requests"
)

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

type Config struct {
	Environment Environment
}

var readIndexFile = sync.Once{}
var indexFile []byte

func Server(ctx context.Context, ac *appcontext.AppContext, config *Config) *echo.Echo {

	e := echo.New()

	isSetup, _ := ac.SetupService.IsSetup(ctx)

	if !isSetup {
		setupHandler := handlers.NewSetupHandler(ac.SetupService)
		e.POST("api/setup", setupHandler.FinalizeSetup)
	}

	// Set up the assets file server.
	contentHandler := echo.WrapHandler(http.FileServer(http.FS(ac.SPAContent)))
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/embedded/$1"})
	e.GET("assets/*", contentHandler, contentRewrite)

	appGroup := e.Group("")

	// If environment is development, setup development middlewares
	if config.Environment == Dev {
		appGroup.Use(mw.Dev())
	}

	if !isSetup {
		appGroup.Use(mw.IsSetup(ac.SetupService))
	}

	homePage := mkHomePageHandler(ac.SPAContent)
	appGroup.GET("app", homePage)
	appGroup.GET("app/*", homePage)

	appGroup.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/app")
	})

	appGroup.GET("/example", mw.Wrap(func(ctx context.Context, c echo.Context) error {
		log.Ctx(ctx).Info("example log", attr.String("first", "value"))
		login := gen.Login{
			Username: "username",
			Password: "password",
		}
		return requests.JSON(c, &login)
	}))

	apiGroup := e.Group("api")

	userHandler := handlers.NewUserHandler(ac.UserService)
	apiGroup.POST("login", userHandler.Login)

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
