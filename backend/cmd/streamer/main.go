package main

import (
	"context"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/devrunner"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cancelled := false
	mainCtx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if cancelled {
				os.Exit(1)
			}
			// sig is a ^C, handle it
			fmt.Println("\nShutting down...")
			cancel()
			cancelled = true
		}
	}()

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	rootPath, err := filepath.Abs(filepath.Join(wd, ".."))
	if err != nil {
		return fmt.Errorf("failed to get root path: %w", err)
	}

	eg, ctx := errgroup.WithContext(mainCtx)

	outputCh := make(chan devrunner.OutputData)

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-outputCh:
				if msg.Err != nil {
					fmt.Println("Error: ", msg.Err.Error())
					continue
				}
				fmt.Println("--", msg.Source, "---------------------------------------------")
				fmt.Println(msg.Output)
			}
		}
	})

	if err := devrunner.WatchBackend(ctx, outputCh, eg, rootPath); err != nil {
		return err
	}

	if err := devrunner.WatchFrontend(ctx, outputCh, eg, rootPath); err != nil {
		return err
	}

	if err := eg.Wait(); err != nil {
		if err != context.Canceled {
			return err
		}
	}

	return nil
}
