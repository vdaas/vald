package runner

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

type Runner interface {
	PreStart() error
	Start(ctx context.Context) <-chan error
	PreStop() error
	Stop(ctx context.Context) error
}

func Run(ctx context.Context, run Runner) (err error) {

	err = run.PreStart()
	if err != nil {
		return err
	}

	rctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ech := run.Start(ctx)
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigCh)

	emap := make(map[string]int)
	for {
		select {
		case <-sigCh:
			log.Warn("daemon is stopping now...")
			cancel()
		case err = <-ech:
			if err != nil {
				log.Error(err)
				emap[err.Error()]++
			}
		case <-rctx.Done():
			err = run.PreStop()
			if err != nil {
				emap[err.Error()]++
			}
			err = run.Stop(ctx)
			if err != nil {
				emap[err.Error()]++
			}
			err = errgroup.Wait()
			if err != nil {
				emap[err.Error()]++
			}
			err = nil
			for msg, count := range emap {
				if msg != "" && strings.HasPrefix(msg, http.ErrServerClosed.Error()) {
					err = errors.Wrapf(err, "error:\t%s\tcount:\t%d", msg, count)
				}
			}
			if err != nil {
				err = errors.ErrDaemonStopFailed(err)
			}
			log.Warn("daemon stopped")
			return err
		}
	}
}
