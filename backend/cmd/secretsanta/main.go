package main

import (
	"embed"
	"github.com/mdev5000/secretsanta/internal/appcontext"
	"github.com/mdev5000/secretsanta/internal/server"
	"github.com/mdev5000/secretsanta/internal/setup"
)

//go:embed all:embedded/*
var spaContent embed.FS

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ac := appcontext.AppContext{
		SetupService: &setup.Service{},
		SPAContent:   spaContent,
	}

	config := server.Config{}

	return server.Server(&ac, &config).Start(":3000")
}
