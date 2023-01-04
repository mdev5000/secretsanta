package main

import (
	"embed"
	"github.com/mdev5000/secretsanta/internal/server"
)

//go:embed all:embedded/*
var spaContent embed.FS

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	config := server.Config{
		SpaContent: spaContent,
	}

	return server.Server(&config).Start(":3000")
}
