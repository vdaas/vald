// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package servers provides implementation of Go API for managing server flow
package servers

import (
	"context"
	"net/http"
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
	sus     []string
	sds     []string
	cancel  context.CancelFunc
	ech     chan error
	sddur   time.Duration
}

func New(opts ...Option) Listener {
	l := new(listener)
	for _, opt := range append(defaultOpts, opts...) {
		opt(l)
	}
	l.ech = make(chan error, len(l.servers)*10)
	return l
}

func (l *listener) ListenAndServe(ctx context.Context) <-chan error {

	type sinfo struct {
		ech <-chan error
		srv server.Server
	}

	var (
		echs = make([]sinfo, len(l.servers))
		rctx context.Context
	)

	for _, name := range l.sus {
		srv, ok := l.servers[name]

		if !ok || srv == nil {
			l.ech <- errors.ErrServerNotFound(name)
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

	errgroup.Go(safety.RecoverFunc(func() (err error) {
		for {
			for i := range echs {
				select {
				case <-rctx.Done():
					err = rctx.Err()
					if err != nil && err != context.Canceled {
						return err
					}
					return nil
				case err = <-echs[i].ech:
					if err != nil && err != http.ErrServerClosed {
						l.ech <- err
					}
					err = nil
					echs[i] = sinfo{
						ech: echs[i].srv.ListenAndServe(),
						srv: echs[i].srv,
					}
				}
			}
		}
	}))
	return l.ech
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
