package devrunner

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unicode"
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
	in := bytes.NewBuffer(nil)

	frontEndPath := filepath.Join(rootPath, "frontend")

	cmd := exec.Command("npm", "run", "dev")
	errOut := bytes.NewBuffer(nil)
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
				//output = strings.Replace(output, ScreenClear, "", -1)
				//re := regexp.MustCompile("[^a-zA-Z. \n]+")
				//output = re.ReplaceAllString(strings.TrimSpace(output), "")

				chars := []rune(string(outputB) + string(errB))
				outChars := make([]rune, 0, len(chars))
				for _, char := range chars {
					if unicode.IsPrint(char) || char == '\n' {
						outChars = append(outChars, char)
					}
				}

				outputCh <- OutputData{
					Status: Data,
					Source: Frontend,
					Output: string(outChars),
				}
			}
		}
	})

	return nil
}
