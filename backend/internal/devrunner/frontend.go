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
	"time"
)

func WatchFrontend(ctx context.Context, outputCh chan OutputData, eg *errgroup.Group, rootPath string) error {
	in := bytes.NewBuffer(nil)

	frontEndPath := filepath.Join(rootPath, "frontend")

	cmd := exec.CommandContext(ctx, "npm", "run", "dev")
	errOut := bytes.NewBuffer(nil)
	cmd.Stderr = errOut
	cmd.Stdout = in
	cmd.Env = os.Environ()
	cmd.Dir = frontEndPath

	eg.Go(func() error {
		ticker := time.NewTicker(200 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				output, _ := io.ReadAll(in)
				if len(output) == 0 {
					continue
				}
				fmt.Println("output")
				outputCh <- OutputData{
					Source: Frontend,
					Output: string(output),
				}
			}
		}
	})

	eg.Go(func() error {
		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("%w: failed with output:\n%s", err, errOut.String())
		}
		return nil
	})

	return nil
}
