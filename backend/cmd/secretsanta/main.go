package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/server"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/user"
)

//go:embed all:embedded/*
var spaContent embed.FS

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()

	mongoURI := "mongodb://root:rootPW@localhost"
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

	config := server.Config{}

	return server.Server(ctx, &ac, &config).Start(":3000")
}
