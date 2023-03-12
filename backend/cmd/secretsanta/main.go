package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/config"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/server"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/log"
	fz "github.com/mdev5000/secretsanta/internal/util/log/flog-zero"
)

//go:embed all:embedded/*
var spaContent embed.FS

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var logger log.Logger = fz.New(os.Stdout)

	ctx := context.Background()
	ctx = log.NewCtx(ctx, logger)

	cfg := config.Config{
		MongoURI: "mongodb://root:rootPW@localhost",
		Env:      string(server.Prod),
	}
	if err := config.LoadConfig(&cfg); err != nil {
		return err
	}

	log.Ctx(ctx).Info("started with config", attr.Interface("config", cfg))

	log.Ctx(ctx).Info("connecting to mongo")
	mongoURI := cfg.MongoURI
	client, err := mongo.Create(mongoURI)
	if err != nil {
		log.Ctx(ctx).Error("failed to start mongo", attr.Err(err))
		return fmt.Errorf("failed to start mongo client: %w", err)
	}
	if err := client.Connect(ctx); err != nil {
		log.Ctx(ctx).Error("failed to connect to mongo", attr.Err(err))
		return fmt.Errorf("failed to connect to mongo: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Ctx(ctx).Error("failed to ping mongo", attr.Err(err))
		// @todo make this a read endpoint thing at some point.
	}
	db := client.Database("ssdev")

	ac := appcontext.AppContext{
		UserService: user.NewService(db.Collection(user.CollectionUsers)),
		SPAContent:  spaContent,
	}

	ac.SetupService = setup.NewSetupService(ac.UserService)

	return runServer(ctx, ac, cfg)
}

func runServer(ctx context.Context, ac appcontext.AppContext, cfg config.Config) (err error) {
	setupCh := make(chan struct{})
	defer close(setupCh)

	serverCfg := server.Config{
		Environment: server.Environment(cfg.Env),
		SetupCh:     setupCh,
	}
	address := ":3000"

	var e *echo.Echo
	go func() {
		<-setupCh
		log.Ctx(ctx).Info("setup has been completed, restarting server")
		if err := e.Shutdown(ctx); err != nil {
			log.Ctx(ctx).Error("error occurred at shutdown during setup restart", attr.Err(err))
		}
	}()

	appIsSetup, err := ac.SetupService.IsSetup(ctx)
	if err != nil {
		return fmt.Errorf("failed to determine if app was setup: %w", err)
	}

	e = server.Server(ctx, &ac, &serverCfg)
	if err := e.Start(address); err != nil && err != http.ErrServerClosed {
		return err
	}

	if !appIsSetup {
		log.Ctx(ctx).Info("app started in non-setup state, so restarting server")
		e = server.Server(ctx, &ac, &serverCfg)
		return e.Start(address)
	}

	return nil
}
