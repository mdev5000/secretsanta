package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mdev5000/secretsanta/internal/devrunner"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func ignoreIt(i ...interface{}) {

}

func run() error {
	mainCtx, cancel := context.WithCancel(context.Background())
	ignoreIt(cancel)

	uiModel := devrunner.UIModel{
		Shutdown:     cancel,
		ShuttingDown: false,
		BackendData: devrunner.WatcherDetails{
			Status: devrunner.Loading,
		},
		FrontendData: devrunner.WatcherDetails{
			Status: devrunner.Loading,
		},
	}
	prog := tea.NewProgram(uiModel)

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
				prog.Send(devrunner.MessageWatcherUpdate{Output: msg})
			}
		}
	})

	if err := devrunner.WatchBackend(ctx, outputCh, eg, rootPath); err != nil {
		return err
	}

	if err := devrunner.WatchFrontend(ctx, outputCh, eg, rootPath); err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		var err error
		if err = eg.Wait(); err != nil {
			if err == context.Canceled {
				err = nil
			}
		}
		prog.Send(devrunner.Terminate{})
		errCh <- err
	}()

	fmt.Println(devrunner.ScreenClear)
	_, err = prog.Run()
	if err != nil {
		return err
	}
	fmt.Println(devrunner.ScreenClear)

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
