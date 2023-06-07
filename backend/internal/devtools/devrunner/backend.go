package devrunner

import (
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

func runBackend(ctx context.Context, backendPath string, outputCh chan OutputData, updateCh chan struct{}, in io.ReadWriter, errB io.Writer) error {
	outputCh <- OutputData{Source: Backend, Status: Compiling}
	output, exitCode := recompileBackend(backendPath)
	if exitCode != 0 {
		outputCh <- OutputData{Source: Backend, Status: ErrorS, Output: string(output)}
		<-updateCh
		return nil
	}
	outputCh <- OutputData{Source: Backend, Status: Loading}
	in.Write([]byte(ScreenClear))
	exe := createRunBackendCommand(backendPath, in, errB)
	if err := exe.Start(); err != nil {
		return fmt.Errorf("failed to start backend command: %w", err)
	}
	select {
	case <-ctx.Done():
		if err := exe.Process.Kill(); err != nil {
			return fmt.Errorf("failed to restart backend process: %w", err)
		}
		return ctx.Err()
	case <-updateCh:
	}
	if err := exe.Process.Kill(); err != nil {
		return fmt.Errorf("failed to restart backend process: %w", err)
	}
	return nil
}

func WatchBackend(ctx context.Context, outputCh chan OutputData, eg *errgroup.Group, rootPath string) error {
	in := newThreadSafeBuffer()
	errB := newThreadSafeBuffer()
	backendPath := filepath.Join(rootPath, "backend")

	backendWatcher := watcher.New()
	if err := backendWatcher.AddRecursive("internal"); err != nil {
		return err
	}

	updateCh := make(chan struct{})

	eg.Go(func() error {
		for {
			if err := runBackend(ctx, backendPath, outputCh, updateCh, in, errB); err != nil {
				return err
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
		var output string
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				outputB, _ := io.ReadAll(in)
				errOutput, _ := io.ReadAll(errB)
				if len(errOutput) != 0 {
					return fmt.Errorf("error while running backend process:\n\n%s", string(errOutput))
				}
				if len(outputB) == 0 {
					continue
				}
				parts := strings.Split(string(outputB), ScreenClear)
				if len(parts) == 1 {
					output += parts[0]
				} else {
					output = parts[len(parts)-1]
				}
				outputCh <- OutputData{
					Status: Data,
					Source: Backend,
					Err:    nil,
					Output: scrubOutput(output),
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

func recompileBackend(backendPath string) ([]byte, int) {
	cmd := exec.Command("make", "build")
	cmd.Env = os.Environ()
	cmd.Dir = backendPath
	output, _ := cmd.CombinedOutput()
	return output, cmd.ProcessState.ExitCode()
}
