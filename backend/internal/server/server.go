package server

import (
	"context"
	"embed"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/mdev5000/secretsanta/internal/requests/gen/core"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/appjson"
	"github.com/mdev5000/secretsanta/internal/util/resp"
	"github.com/mdev5000/secretsanta/internal/util/session"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/handlers"
	mw "github.com/mdev5000/secretsanta/internal/middleware"
	"github.com/mdev5000/secretsanta/internal/util/log"
)

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

func (e Environment) IsDev() bool {
	return e == Dev || e == "development"
}

type Config struct {
	Environment Environment
	SetupCh     chan struct{}
}

var readIndexFile = sync.Once{}
var indexFile []byte

type commonHandlers struct {
	setup *handlers.SetupHandler
}

type server struct {
	appCtx     context.Context
	appContext *appcontext.AppContext
	config     *Config
	e          *echo.Echo
	handlers   commonHandlers
	sessionMgr *scs.SessionManager
}

func wrap[T proto.Message](s *server, h func(context.Context, echo.Context) resp.Response[T]) echo.HandlerFunc {
	if err := mw.EnsureComplies[T](); err != nil {
		panic(err)
	}
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		rs := h(ctx, c)
		if err := session.TrySaveSession(ctx, s.sessionMgr, c); err != nil {
			return apperror.InternalError(fmt.Errorf("failed to save session data: %w", err))
		}
		if rs.Err != nil {
			return rs.Err
		}
		if rs.Data == nil {
			panic("either Err or Data must be setup on a response")
		}
		return c.JSONBlob(rs.Code, rs.Data)
	}
}

func (s *server) wrap(h mw.AppHandler) echo.HandlerFunc {
	return mw.Wrap(h)
}

func apiOptions(methods ...string) echo.HandlerFunc {
	allow := strings.Join(methods, ",")
	return func(c echo.Context) error {
		headers := c.Response().Header()
		headers.Add("Allow", allow)
		headers.Add("Accept", "application/json")
		headers.Add("Access-Control-Request-Method", allow)
		return c.Blob(http.StatusNoContent, "application/json", nil)
	}
}

func (s *server) setupServer() {
	apiGroup := s.e.Group("/api")
	s.apiBase(apiGroup)
	apiGroup.GET("/setup/status", wrap(s, s.handlers.setup.Status))
	apiGroup.GET("/setup/leader-status", s.wrap(s.handlers.setup.LeaderStatus))
	apiGroup.POST("/setup/finalize", wrap(s, s.handlers.setup.FinalizeSetup))
	apiGroup.OPTIONS("/setup/finalize", apiOptions("POST"))
	apiGroup.POST("/setup/finalize-quick", s.wrap(s.handlers.setup.FinalizeSetupQuick))
	apiGroup.OPTIONS("/setup/finalize-quick", apiOptions("POST"))
	s.exampleAPIRoute(apiGroup)

	appGroup := s.e.Group("")
	appGroup.Use(mw.IsSetup(s.appContext.SetupService))
	s.appRoutes(appGroup)
}

func (s *server) appRoutes(appGroup *echo.Group) {
	// Set up the assets file server.
	contentHandler := echo.WrapHandler(http.FileServer(http.FS(s.appContext.SPAContent)))
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/embedded/$1"})

	s.e.GET("assets/*", contentHandler, contentRewrite)
	homePage := mkHomePageHandler(s.sessionMgr, s.appContext.SPAContent)
	appGroup.GET("app", homePage)
	appGroup.GET("app/*", homePage)

	appGroup.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/app")
	})
}

func (s *server) apiBase(apiGroup *echo.Group) {
	//apiGroup.Use(mw.APIBase(s.appCtx, s.config.Environment.IsDev()))

	// If environment is development, setup development middlewares
	if s.config.Environment.IsDev() {
		apiGroup.Use(mw.ApiDev())
	}
}

func (s *server) apiRoutes(apiGroup *echo.Group) {
	s.apiBase(apiGroup)

	apiGroup.GET("/setup/status", wrap(s, s.handlers.setup.Status))

	s.exampleAPIRoute(apiGroup)

	userHandler := handlers.NewAuthHandler(s.appContext.UserService, s.sessionMgr)
	apiGroup.POST("/login", wrap(s, userHandler.Login))
	apiGroup.OPTIONS("/login", apiOptions("POST"))
	apiGroup.POST("/logout", wrap(s, userHandler.Logout))
	apiGroup.OPTIONS("/logout", apiOptions("POST"))

	// Authenticated  API routes

	authGroup := apiGroup.Group("",
		mw.EnsureLoggedIn(s.appCtx, s.sessionMgr),
	)
	authGroup.GET("/auth-test", func(c echo.Context) error {
		return c.JSONBlob(202, []byte(`{"something": "yay"}`))
	})
}

func (s *server) exampleAPIRoute(apiGroup *echo.Group) {
	apiGroup.GET("/example", s.wrap(func(ctx context.Context, c echo.Context) error {
		log.Ctx(ctx).Info("example log", attr.String("first", "value"))
		login := core.Login{
			Username: "username",
			Password: "password",
		}
		return appjson.JSONOk(c, &login)
	}))
}

func Server(appCtx context.Context, ac *appcontext.AppContext, config *Config) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = mw.ErrorHandler(appCtx)

	sessionStore := session.New(ac.Db, config.Environment.IsDev())

	s := server{
		appCtx:     appCtx,
		appContext: ac,
		config:     config,
		e:          e,
		sessionMgr: sessionStore,
		handlers: commonHandlers{
			setup: handlers.NewSetupHandler(ac.SetupService, appCtx, config.SetupCh),
		},
	}

	// Setup application logging.
	e.Use(mw.Logging(s.appCtx, s.config.Environment.IsDev()))

	// Load and save sessions.
	e.Use(mw.LoadSessionData(s.sessionMgr))

	isSetup, _ := ac.SetupService.IsSetup(appCtx)

	if !isSetup {
		log.Ctx(appCtx).Info("app is not setup, starting in setup mode")
		// Server starts with a different set of routes when setting up, then panics and restarts.
		s.setupServer()
		return s.e
	}

	apiGroup := s.e.Group("/api")
	s.apiRoutes(apiGroup)

	appGroup := s.e.Group("")
	s.appRoutes(appGroup)

	return s.e
}

func mkHomePageHandler(sm *session.Manager, spaContent embed.FS) echo.HandlerFunc {
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
		ctx := c.Request().Context()
		if err != nil {
			log.Ctx(ctx).Error("could not read SPA index", attr.Err(err))
			return apperror.InternalError(err)
		}
		if !isNonAuthPage(c) {
			isLoggedIn, err := session.IsLoggedIn(ctx, sm)
			if err != nil {
				return apperror.InternalError(err)
			}
			if !isLoggedIn {
				return c.Redirect(http.StatusTemporaryRedirect, "/app/login")
			}
		}
		return c.Blob(200, "text/html", indexFile)
	}
}

func isNonAuthPage(c echo.Context) bool {
	path := c.Request().URL.Path
	return path == "/app/login" ||
		path == "/app/logout" ||
		strings.HasPrefix(path, "/app/setup") ||
		// @todo should probably remove this
		strings.HasPrefix(path, "/app/example")
}
