package main

import (
	"context"
	"embed"
	"fmt"
	"os"

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

	mongoURI := cfg.MongoURI
	client, err := mongo.Create(mongoURI)
	if err != nil {
		return fmt.Errorf("failed to start mongo client: %w", err)
	}
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to mongo: %w", err)
	}
	db := client.Database("ssdev")

	ac := appcontext.AppContext{
		UserService: user.NewService(db.Collection(user.CollectionUsers)),
		SPAContent:  spaContent,
	}

	ac.SetupService = setup.NewSetupService(ac.UserService)

	serverCfg := server.Config{
		Environment: server.Environment(cfg.Env),
	}

	return server.Server(ctx, &ac, &serverCfg).Start(":3000")
}
