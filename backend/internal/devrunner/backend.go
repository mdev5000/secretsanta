package devrunner

import (
	"bytes"
	"context"
	"fmt"
	"github.com/radovskyb/watcher"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func WatchBackend(ctx context.Context, outputCh chan OutputData, eg *errgroup.Group, rootPath string) error {
	in := bytes.NewBuffer(nil)
	errB := bytes.NewBuffer(nil)
	backendPath := filepath.Join(rootPath, "backend")

	backendWatcher := watcher.New()
	if err := backendWatcher.AddRecursive("internal"); err != nil {
		return err
	}

	updateCh := make(chan struct{})

	eg.Go(func() error {
		for {
			outputCh <- OutputData{Source: Backend, Status: Compiling}
			output, err := recompileBackend(backendPath)
			if err != nil {
				fmt.Errorf("failed to recompile backend: %w:\n%s", err, string(output))
			}
			outputCh <- OutputData{Source: Backend, Status: Loading}
			in.WriteString(ScreenClear)
			exe := createRunBackendCommand(backendPath, in, errB)
			if err := exe.Start(); err != nil {
				return fmt.Errorf("failed to start backend command: %w", err)
			}
			select {
			case <-ctx.Done():
				if err := exe.Process.Kill(); err != nil {
					return fmt.Errorf("failed to restart backend process: %w", err)
				}
				return nil
			case <-updateCh:
			}
			if err := exe.Process.Kill(); err != nil {
				return fmt.Errorf("failed to restart backend process: %w", err)
			}
		}
	})

	eg.Go(func() error {
		backendWatcher.Wait()
		for {
			select {
			case <-ctx.Done():
				backendWatcher.Close()
				return ctx.Err()
			case <-backendWatcher.Event:
				updateCh <- struct{}{}
			case err := <-backendWatcher.Error:
				return err
			}
		}
	})

	eg.Go(func() error {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				output, _ := io.ReadAll(in)
				errOutput, _ := io.ReadAll(errB)
				if len(errOutput) != 0 {
					return fmt.Errorf("error while running backend process:\n\n%s", string(errOutput))
				}
				if len(output) == 0 {
					continue
				}
				outputCh <- OutputData{
					Status: Data,
					Source: Backend,
					Err:    nil,
					Output: strings.Replace(string(output), ScreenClear, "", -1),
				}
			}
		}
	})

	go func() {
		if err := backendWatcher.Start(time.Millisecond * 100); err != nil {
			panic(fmt.Errorf("failed to start backend watcher: %w", err))
		}
	}()

	return nil
}

func createRunBackendCommand(backendPath string, in io.Writer, errB io.Writer) *exec.Cmd {
	cmd := exec.Command("./_build/secretsanta")
	cmd.Stderr = errB
	cmd.Stdout = in
	cmd.Env = append(os.Environ(), "ENV=development")
	cmd.Dir = backendPath
	return cmd
}

func recompileBackend(backendPath string) ([]byte, error) {
	cmd := exec.Command("make", "build")
	cmd.Env = os.Environ()
	cmd.Dir = backendPath
	return cmd.CombinedOutput()
}
