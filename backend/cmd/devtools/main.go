package main

import (
	"bytes"
	"context"
	"fmt"
	devrunner2 "github.com/mdev5000/secretsanta/internal/devtools/devrunner"
	scw "github.com/mdev5000/secretsanta/internal/devtools/svelte_check_wrapper"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func ignoreIt(i ...interface{}) {

}

func run() error {
	rootCmd := &cobra.Command{
		Use:   "devtools",
		Short: "Set of development tools",
		Args:  cobra.MinimumNArgs(1),
	}

	watchCmd := &cobra.Command{
		Use:   "watcher",
		Short: "Start the development watcher",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatcher()
		},
	}

	svelteCheck := &cobra.Command{
		Use:   "svelte-check",
		Short: "Run svelte-check but filter useless checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSvelteCheck(cmd.OutOrStdout())
		},
	}

	rootCmd.AddCommand(
		watchCmd,
		svelteCheck,
	)

	return rootCmd.Execute()
}

func runWatcher() error {
	mainCtx, cancel := context.WithCancel(context.Background())
	ignoreIt(cancel)

	uiModel := devrunner2.UIModel{
		Shutdown:     cancel,
		ShuttingDown: false,
		BackendData: devrunner2.WatcherDetails{
			Status: devrunner2.Loading,
		},
		FrontendData: devrunner2.WatcherDetails{
			Status: devrunner2.Loading,
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

	outputCh := make(chan devrunner2.OutputData)

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-outputCh:
				prog.Send(devrunner2.MessageWatcherUpdate{Output: msg})
			}
		}
	})

	if err := devrunner2.WatchBackend(ctx, outputCh, eg, rootPath); err != nil {
		return err
	}

	if err := devrunner2.WatchFrontend(ctx, outputCh, eg, rootPath); err != nil {
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
		prog.Send(devrunner2.Terminate{})
		errCh <- err
	}()

	fmt.Println(devrunner2.ScreenClear)
	_, err = prog.Run()
	if err != nil {
		return err
	}
	fmt.Println(devrunner2.ScreenClear)

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}

func runSvelteCheck(out io.Writer) error {
	cmd := exec.Command("npm", "run", "check:machine")
	frontendPath, err := filepath.Abs("../frontend")
	if err != nil {
		return fmt.Errorf("frontend path error: %w", err)
	}
	cmd.Dir = frontendPath
	o, _ := cmd.CombinedOutput()
	checks := scw.ParseLines(bytes.NewBuffer(o), scw.And(
		scw.IgnoreDataTestIdMessage,
		scw.IgnoreNodeModules,
	))
	success := scw.Passed(checks)
	for _, check := range checks {
		check.Print(out)
	}
	if success {
		return nil
	}
	os.Exit(1)
	return nil
}
