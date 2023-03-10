package devrunner

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func ensureKilled(cmd *exec.Cmd) error {
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Release()
		}
	}()
	if cmd.Process == nil {
		return fmt.Errorf("could not kill front end: process does not exist")
	}
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

func WatchFrontend(ctx context.Context, outputCh chan OutputData, eg *errgroup.Group, rootPath string) error {
	in := newThreadSafeBuffer()
	errOut := newThreadSafeBuffer()

	frontEndPath := filepath.Join(rootPath, "frontend")

	cmd := exec.Command("npm", "run", "dev")
	cmd.Stderr = errOut
	cmd.Stdout = in
	cmd.Env = os.Environ()
	cmd.Dir = frontEndPath
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Setpgid = true

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%w: failed to start command with output:\n%s", err, errOut.String())
	}

	eg.Go(func() error {
		ticker := time.NewTicker(200 * time.Millisecond)
		var output string
		for {
			select {
			case <-ctx.Done():
				return ensureKilled(cmd)
			case <-ticker.C:
				if cmd.ProcessState != nil {
					if cmd.ProcessState.ExitCode() != -1 {
						b, _ := cmd.CombinedOutput()
						return fmt.Errorf("frontend crashed unexpectedly:\n%s", string(b))
					}
				}

				outputB, _ := io.ReadAll(in)
				if len(outputB) == 0 {
					continue
				}

				errB, _ := io.ReadAll(errOut)

				newOutput := string(outputB)
				// > secretsanta is the first line of output so this indicates that the screen has been cleared.
				// Haven't figured out a better way to detect this yet. :/
				parts := strings.Split(strings.TrimSpace(newOutput), "> secretsanta")
				if len(parts) == 1 {
					output += parts[0]
				} else {
					output = "> secretsanta" + parts[len(parts)-1]
				}

				outputAll := output + string(errB)

				outputCh <- OutputData{
					Status: Data,
					Source: Frontend,
					Output: scrubOutput(outputAll),
				}
			}
		}
	})

	return nil
}
