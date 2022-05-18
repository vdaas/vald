//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package servers provides implementation of Go API for managing server flow
package servers

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/strings"
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
	sddur   time.Duration
}

func New(opts ...Option) Listener {
	l := new(listener)
	for _, opt := range append(defaultOptions, opts...) {
		opt(l)
	}

	if l.sus != nil && len(l.sus) != 0 && (l.sds == nil || len(l.sds) == 0) {
		s := make([]string, len(l.sus))
		copy(s, l.sus)
		l.sds = s
		sort.Sort(sort.Reverse(sort.StringSlice(l.sds)))
	}
	return l
}

func (l *listener) ListenAndServe(ctx context.Context) <-chan error {
	ech := make(chan error, len(l.servers)*10)

	for _, name := range l.sus {
		srv, ok := l.servers[name]

		if !ok || srv == nil {
			select {
			case <-ctx.Done():
			case ech <- errors.ErrServerNotFound(name):
			}
			continue
		}

		if !l.servers[name].IsRunning() {
			err := l.servers[name].ListenAndServe(ctx, ech)
			if err != nil {
				select {
				case <-ctx.Done():
				case ech <- err:
				}
			}
		}
	}

	for name := range l.servers {
		if !l.servers[name].IsRunning() {
			err := l.servers[name].ListenAndServe(ctx, ech)
			if err != nil {
				select {
				case <-ctx.Done():
				case ech <- err:
				}
			}
		}
	}

	l.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		<-ctx.Done()
		err = ctx.Err()
		if err != nil &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return nil
	}))
	return ech
}

func (l *listener) Shutdown(ctx context.Context) (err error) {
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
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				emap[err.Error()] = struct{}{}
			}
		}
	}

	for name := range l.servers {
		if l.servers[name].IsRunning() {
			err = l.servers[name].Shutdown(ctx)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				emap[err.Error()] = struct{}{}
			}
		}
	}

	err = nil
	for msg := range emap {
		if msg != "" && !strings.HasPrefix(msg, http.ErrServerClosed.Error()) {
			err = errors.Wrap(err, msg)
		}
	}

	return err
}
