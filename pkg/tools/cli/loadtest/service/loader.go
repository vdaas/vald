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
package service

import (
	"context"
	"os"
	"reflect"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

type Loader interface {
	Prepare(context.Context) error
	Do(context.Context) <-chan error
}

type loader struct {
	eg           errgroup.Group
	client       grpc.Client
	addr         string
	concurrency  int
	dataset      string
	requests     []interface{}
	progressDuration time.Duration
	requestsFunc func(assets.Dataset) ([]interface{}, error)
	loaderFunc   func(context.Context, vald.ValdClient, interface{}, ...grpc.CallOption) error
}

func newLoader(opts... Option) (l *loader, err error) {
	l = new(loader)
	for _, opt := range append(defaultOpts, opts...) {
		if err = opt(l); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return l, nil
}

func (l *loader) Prepare(context.Context) (err error) {
	fn := assets.Data(l.dataset)
	if fn == nil {
		return errors.Errorf("dataset load function is nil: %s", l.dataset)
	}
	dataset, err := fn()
	if err != nil {
		return err
	}
	l.requests, err = l.requestsFunc(dataset)
	return err
}

func (l *loader) Do(ctx context.Context) <-chan error {
	ech := make(chan error, len(l.requests))

	// TODO: related to #403.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		ech <- err
		return ech
	}

	var pg int32 = 0
	progress := func() {
		log.Debugf("progress %d items", pg)
	}
	ticker := time.NewTicker(l.progressDuration)
	l.eg.Go(safety.RecoverFunc(func() error {
		for pg != int32(len(l.requests)) {
			select {
			case <-ctx.Done():
				progress()
				return nil
			case <-ticker.C:
				progress()
			}
		}
		return nil
	}))
	l.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		eg, egctx := errgroup.New(ctx)
		eg.Limitation(l.concurrency)
		for _, req := range l.requests {
			r := req
			eg.Go(func() error {
				_, err := l.client.Do(egctx, l.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
					err := l.loaderFunc(ctx, vald.NewValdClient(conn), r, copts...)
					atomic.AddInt32(&pg, 1)
					if err != nil {
						log.Warn(err)
					}
					return nil, err
				})
				if err != nil {
					ech <- err
				}
				return nil
			})
		}
		err := eg.Wait()
		time.Sleep(5 * time.Second) // prevent too early shutdown
		if err != nil {
			log.Warn(err)
			ech <- err
			return p.Signal(syscall.SIGKILL) // TODO: #403
		}
		return p.Signal(syscall.SIGTERM) // TODO: #403
	}))
	return ech
}
