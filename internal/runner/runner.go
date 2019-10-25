//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	PreStop(ctx context.Context) error
	Stop(ctx context.Context) error
	PostStop(ctx context.Context) error
}

func Run(ctx context.Context, run Runner) (err error) {

	err = run.PreStart(ctx)
	if err != nil {
		return err
	}

	rctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ech := run.Start(ctx)

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	emap := make(map[string]int)
	errs := make([]error, 0, 10)

	for {
		select {
		case <-sigCh:
			log.Warn("daemon is stopping now...")
			cancel()
		case err = <-ech:
			if err != nil {
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				log.Error(err)
				emap[err.Error()]++
			}
		case <-rctx.Done():
			err = run.PreStop(ctx)
			if err != nil {
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				log.Error(err)
				emap[err.Error()]++

			}
			err = run.Stop(ctx)
			if err != nil {
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				log.Error(err)
				emap[err.Error()]++
			}
			err = run.PostStop(ctx)
			if err != nil {
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				log.Error(err)
				emap[err.Error()]++
			}
			err = errgroup.Wait()
			if err != nil {
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				log.Error(err)
				emap[err.Error()]++
			}
			err = nil
			for _, ierr := range errs {
				if ierr != nil {
					msg := ierr.Error()
					if msg != "" && !strings.Contains(msg, http.ErrServerClosed.Error()) {
						err = errors.Wrapf(err, "error:\t%s\tcount:\t%d", msg, emap[msg])
					}
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
