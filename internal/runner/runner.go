//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
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
