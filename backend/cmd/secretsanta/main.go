package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/util/env"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/config"
	"github.com/mdev5000/secretsanta/internal/family"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/server"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/log"
	fz "github.com/mdev5000/secretsanta/internal/util/log/flog-zero"
	"github.com/mdev5000/secretsanta/internal/util/transactions"
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
		Env:      string(env.Prod),
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
		Db:             db,
		TransactionMgr: transactions.NoTransactions(),
		UserService:    user.NewService(db.Collection(user.CollectionUsers)),
		FamilyService:  family.NewService(db.Collection(family.CollectionFamilies)),
		SPAContent:     spaContent,
	}

	ac.SetupService = setup.NewService(ac.TransactionMgr, ac.UserService, ac.FamilyService)

	restart := true
	for restart {
		if err := resetCaches(ctx, ac); err != nil {
			return err
		}
		var err error
		restart, err = runServer(ctx, ac, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func resetCaches(ctx context.Context, ac appcontext.AppContext) error {
	ac.SetupService.ClearCache()
	return nil
}

func runServer(ctx context.Context, ac appcontext.AppContext, cfg config.Config) (restart bool, err error) {
	termCh := make(chan struct{})
	restartLock := sync.RWMutex{}

	serverCfg := server.Config{
		Environment: env.Environment(cfg.Env),
		TermCh:      termCh,
	}
	address := ":3000"

	var e *echo.Echo
	go func() {
		<-termCh
		restartLock.Lock()
		restart = true
		restartLock.Unlock()
		log.Ctx(ctx).Info("setup or delete-all has been completed, restarting server")
		if err := e.Shutdown(ctx); err != nil {
			log.Ctx(ctx).Error("error occurred at shutdown during setup restart", attr.Err(err))
		}
	}()

	//appIsSetup, err := ac.SetupService.IsSetup(ctx)
	//if err != nil {
	//	return false, fmt.Errorf("failed to determine if app was setup: %w", err)
	//}

	getRestart := func() bool {
		restartLock.RLock()
		defer restartLock.RUnlock()
		return restart
	}

	e = server.Server(ctx, &ac, &serverCfg)
	if err := e.Start(address); err != nil && err != http.ErrServerClosed {
		return getRestart(), err
	}

	//if !appIsSetup {
	//	log.Ctx(ctx).Info("app started in non-setup state, so restarting server")
	//	e = server.Server(ctx, &ac, &serverCfg)
	//	return e.Start(address)
	//}

	return getRestart(), nil
}
