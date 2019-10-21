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

// Package servers provides implementation of Go API for managing server flow
package servers

import (
	"context"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
)

type Listener interface {
	ListenAndServe(context.Context) <-chan error
	Shutdown(context.Context) error
}

type listener struct {
	servers map[string]server.Server
	eg      errgroup.Group
	sus     []string
	sds     []string
	cancel  context.CancelFunc
	sddur   time.Duration
}

func New(opts ...Option) Listener {
	l := new(listener)
	for _, opt := range append(defaultOpts, opts...) {
		opt(l)
	}
	return l
}

func (l *listener) ListenAndServe(ctx context.Context) <-chan error {

	type sinfo struct {
		ech <-chan error
		srv server.Server
	}

	ech := make(chan error, len(l.servers)*10)
	var rctx context.Context
	echs := make([]sinfo, len(l.servers))

	for _, name := range l.sus {
		srv, ok := l.servers[name]

		if !ok || srv == nil {
			ech <- errors.ErrServerNotFound(name)
			continue
		}

		if !l.servers[name].IsRunning() {
			echs = append(echs, sinfo{
				ech: l.servers[name].ListenAndServe(),
				srv: l.servers[name],
			})
		}
	}

	for name := range l.servers {
		if !l.servers[name].IsRunning() {
			echs = append(echs, sinfo{
				ech: l.servers[name].ListenAndServe(),
				srv: l.servers[name],
			})
		}
	}

	rctx, l.cancel = context.WithCancel(ctx)

	l.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			for i := range echs[:len(echs)] {
				select {
				case <-rctx.Done():
					err = rctx.Err()
					if err != nil && err != context.Canceled {
						return err
					}
					return nil
				case err = <-echs[i].ech:
					if err != nil && err != http.ErrServerClosed {
						ech <- err
					}
					err = nil
					echs[i] = sinfo{
						ech: echs[i].srv.ListenAndServe(),
						srv: echs[i].srv,
					}
				}
			}
			runtime.Gosched()
		}
	}))
	return ech
}

func (l *listener) Shutdown(ctx context.Context) (err error) {

	defer l.cancel()

	ctx, cancel := context.WithTimeout(ctx, l.sddur)
	defer cancel()

	emap := make(map[string]struct{})

	for _, name := range l.sds {
		srv, ok := l.servers[name]

		if !ok || srv == nil {
			return errors.ErrServerNotFound(name)
		}

		if l.servers[name].IsRunning() {
			err = l.servers[name].Shutdown(ctx)
			if err != nil && err != http.ErrServerClosed {
				emap[err.Error()] = struct{}{}
			}
		}
	}

	for name := range l.servers {
		if l.servers[name].IsRunning() {
			err = l.servers[name].Shutdown(ctx)
			if err != nil && err != http.ErrServerClosed {
				emap[err.Error()] = struct{}{}
			}
		}
	}

	err = nil
	for msg := range emap {
		if msg != "" && strings.HasPrefix(msg, http.ErrServerClosed.Error()) {
			err = errors.Wrap(err, msg)
		}
	}

	return err
}
