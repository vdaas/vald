//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/params"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/timeutil/location"
	ver "github.com/vdaas/vald/internal/version"
	"go.uber.org/automaxprocs/maxprocs"
)

type Runner interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) (<-chan error, error)
	PreStop(ctx context.Context) error
	Stop(ctx context.Context) error
	PostStop(ctx context.Context) error
}

type runner struct {
	version          string
	maxVersion       string
	minVersion       string
	name             string
	loadConfig       func(string) (interface{}, *config.GlobalConfig, error)
	initializeDaemon func(interface{}) (Runner, error)
}

func Do(ctx context.Context, opts ...Option) error {
	r := new(runner)

	for _, opt := range append(defaultOpts, opts...) {
		opt(r)
	}

	info.Init(r.name)

	p, isHelp, err := params.New(
		params.WithConfigFileDescription(fmt.Sprintf("%s config file path", r.name)),
	).Parse()

	if isHelp || err != nil {
		log.Init(log.WithLevel(level.FATAL.String()))
		return err
	}

	if p.ShowVersion() {
		log.Init(log.WithLevel(level.INFO.String()))
		log.Info(info.String())
		return nil
	}

	cfg, ccfg, err := r.loadConfig(p.ConfigFilePath())
	if err != nil {
		log.Init(log.WithLevel(level.FATAL.String()))
		return err
	}

	if lcfg := ccfg.Logging; lcfg != nil {
		log.Init(
			log.WithLoggerType(lcfg.Logger),
			log.WithLevel(lcfg.Level),
			log.WithFormat(lcfg.Format),
		)
	} else {
		log.Init()
	}

	// set location temporary for initialization logging
	location.Set(ccfg.TZ)

	err = ver.Check(ccfg.Version, r.maxVersion, r.minVersion)
	if err != nil {
		return err
	}

	mfunc, err := maxprocs.Set(maxprocs.Logger(log.Infof))
	if err != nil {
		mfunc()
		return err
	}

	daemon, err := r.initializeDaemon(cfg)
	if err != nil {
		return err
	}

	log.Infof("service %s %s starting...", r.name, ccfg.Version)

	// reset timelocation to override external libs & running logging
	location.Set(ccfg.TZ)
	return Run(ctx, daemon, r.name)
}

func Run(ctx context.Context, run Runner, name string) (err error) {
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	rctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rctx = errgroup.Init(rctx)

	log.Info("executing daemon pre-start function")
	err = run.PreStart(rctx)
	if err != nil {
		return err
	}

	log.Info("executing daemon start function")
	ech, err := run.Start(rctx)
	if err != nil {
		return errors.ErrDaemonStartFailed(err)
	}

	emap := make(map[string]int)
	errs := make([]error, 0, 10)

	for {
		select {
		case sig := <-sigCh:
			log.Warnf("%s signal received daemon will stopping soon...", sig)
			cancel()
		case err = <-ech:
			if err != nil {
				log.Error(errors.ErrStartFunc(name, err))
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				emap[err.Error()]++
			}
		case <-rctx.Done():
			log.Info("executing daemon pre-stop function")
			err = safety.RecoverFunc(func() error {
				return run.PreStop(ctx)
			})()
			if err != nil {
				log.Error(errors.ErrPreStopFunc(name, err))
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				emap[err.Error()]++
			}

			log.Info("executing daemon stop function")
			err = safety.RecoverFunc(func() error {
				return run.Stop(ctx)
			})()
			if err != nil {
				log.Error(errors.ErrStopFunc(name, err))
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				emap[err.Error()]++
			}

			log.Info("executing daemon post-stop function")
			err = safety.RecoverFunc(func() error {
				return run.PostStop(ctx)
			})()
			if err != nil {
				log.Error(errors.ErrPostStopFunc(name, err))
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				emap[err.Error()]++
			}

			err = errgroup.Wait()
			if err != nil {
				log.Error(errors.ErrRunnerWait(name, err))
				if _, ok := emap[err.Error()]; !ok {
					errs = append(errs, err)
				}
				emap[err.Error()]++
			}

			err = nil
			for _, ierr := range errs {
				if ierr != nil {
					msg := ierr.Error()
					if msg != "" &&
						!strings.Contains(msg, http.ErrServerClosed.Error()) &&
						!strings.Contains(msg, context.Canceled.Error()) {
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
