package main

import "github.com/mdev5000/secretsanta/internal/server"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	config := server.Config{}

	return server.Server(config).Start(":3000")
}
