package server

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/handlers"
	mw "github.com/mdev5000/secretsanta/internal/middleware"
	"github.com/mdev5000/secretsanta/internal/requests/gen/setup"
	"github.com/mdev5000/secretsanta/internal/util/env"
	"github.com/mdev5000/secretsanta/internal/util/log"
	"github.com/mdev5000/secretsanta/internal/util/resp"
	"github.com/mdev5000/secretsanta/internal/util/session"
	"net/http"
)

type Config struct {
	Environment env.Environment
	TermCh      chan struct{}
}

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
	rb         *routesBuilder
}

func Server(appCtx context.Context, ac *appcontext.AppContext, config *Config) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = mw.ErrorHandler(appCtx)

	sessionStore := session.New(ac.Db, config.Environment.IsDev())

	s := server{
		appCtx:     appCtx,
		appContext: ac,
		config:     config,
		rb:         &routesBuilder{},
		e:          e,
		sessionMgr: sessionStore,
		handlers: commonHandlers{
			setup: handlers.NewSetupHandler(ac.SetupService, appCtx, config.TermCh),
		},
	}

	// Setup application logging.
	e.Use(mw.Logging(s.appCtx, s.config.Environment.IsDev()))

	s.healthRoutes()

	// Load sessions (saved separately due to issues setting headers after body is returned).
	e.Use(mw.LoadSessionData(s.sessionMgr))

	if s.config.Environment.NonProd() {
		s.testRoutes()
	}

	// If the application as not been set up, start in setup mode.
	isSetup, _ := ac.SetupService.IsSetup(appCtx)
	log.Ctx(appCtx).Info("checking setup status", attr.Bool("isSetup", isSetup))
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

func (s *server) setupServer() {
	apiGroup := s.e.Group("/api")
	s.apiCommonBase(apiGroup)
	apiGroup.GET("/setup/status", wrapAPI(s, s.handlers.setup.Status))
	apiGroup.GET("/setup/leader-status", wrapAPI(s, s.handlers.setup.LeaderStatus))

	s.rb.Group(apiGroup, "/setup/finalize").
		POST(wrapAPI(s, s.handlers.setup.FinalizeSetup)).
		Build()

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
	homePageHandler := handlers.MkHomePageHandler(s.sessionMgr, s.appContext.SPAContent)
	appGroup.GET("app", homePageHandler)
	appGroup.GET("app/*", homePageHandler)

	appGroup.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/app")
	})
}

// apiCommonBase sets up api middleware that is common in both setup and general apis
func (s *server) apiCommonBase(apiGroup *echo.Group) {
	//apiGroup.Use(mw.APIBase(s.appCtx, s.config.Environment.IsDev()))

	// If environment is development, setup development middlewares
	if s.config.Environment.IsDev() {
		apiGroup.Use(mw.ApiDev())
	}
}

func (s *server) apiRoutes(apiGroup *echo.Group) {
	s.apiCommonBase(apiGroup)

	apiGroup.GET("/setup/status", wrapAPI(s, s.handlers.setup.Status))

	s.exampleAPIRoute(apiGroup)

	userHandler := handlers.NewAuthHandler(s.appContext.UserService, s.sessionMgr)
	apiGroup.POST("/login", wrapAPI(s, userHandler.Login))
	apiGroup.OPTIONS("/login", handlers.ApiOptions("POST"))
	apiGroup.POST("/logout", wrapAPI(s, userHandler.Logout))
	apiGroup.OPTIONS("/logout", handlers.ApiOptions("POST"))

	// Authenticated  API routes

	authGroup := apiGroup.Group("",
		mw.EnsureLoggedIn(s.appCtx, s.sessionMgr),
	)
	authGroup.GET("/auth-test", func(c echo.Context) error {
		return c.JSONBlob(202, []byte(`{"something": "yay"}`))
	})
}

// testRoutes adds routes useful for testing.
func (s *server) testRoutes() {
	testGroup := s.e.Group("/test")

	testHandler := handlers.Test{
		Db:     s.appContext.Db,
		TermCh: s.config.TermCh,
	}

	testGroup.POST("/delete-all", testHandler.DeleteAll)
	testGroup.POST("/delete-all-restart", testHandler.DeleteAllAndRestart)
}

func (s *server) exampleAPIRoute(apiGroup *echo.Group) {
	apiGroup.GET("/example", wrapAPI(s, func(ctx context.Context, c echo.Context) resp.Response[*setup.Status] {
		log.Ctx(ctx).Info("example log", attr.String("first", "value"))
		status := setup.Status{
			Status: "some status",
		}
		return resp.Ok(200, &status)
	}))
}

func (s *server) healthRoutes() {
	h := handlers.Health{
		Db: s.appContext.Db,
	}
	s.e.GET("ready", h.Ready)
}
