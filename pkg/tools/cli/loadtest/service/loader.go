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
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
)

// Loader is representation of load test
type Loader interface {
	Prepare(context.Context) error
	Do(context.Context) <-chan error
}

type (
	requestFunc func(assets.Dataset) ([]interface{}, error)
	loaderFunc  func(context.Context, vald.ValdClient, interface{}, ...grpc.CallOption) error
)

type loader struct {
	eg               errgroup.Group
	client           grpc.Client
	addr             string
	concurrency      int
	dataset          string
	requests         []interface{}
	progressDuration time.Duration
	requestsFunc     requestFunc
	loaderFunc       loaderFunc
	operation        config.Operation
}

// NewLoader returns Loader implementation.
func NewLoader(opts ...Option) (Loader, error) {
	l := new(loader)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(l); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	switch l.operation {
	case config.Insert:
		l.requestsFunc, l.loaderFunc = newInsert()
	case config.Search:
		l.requestsFunc, l.loaderFunc = newSearch()
	default:
		return nil, errors.Errorf("undefined method: %v", l.operation)
	}

	return l, nil
}

// Prepare generate request data.
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
	if err != nil {
		return err
	}

	if len(l.requests) == 0 {
		return errors.New("prepare data is empty")
	}

	return nil
}

// Do operates load test.
func (l *loader) Do(ctx context.Context) <-chan error {
	ech := make(chan error, len(l.requests))

	// TODO: related to #403.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		select {
		case <-ctx.Done():
			ech <- errors.Wrap(err, ctx.Err().Error())
		case ech <- err:
		}
		return ech
	}

	var pgCnt int32 = 0
	var start time.Time
	progress := func() {
		log.Infof("progress %d items, %f[qps]", pgCnt, float64(pgCnt)/time.Now().Sub(start).Seconds())
	}
	l.eg.Go(safety.RecoverFunc(func() error {
		ticker := time.NewTicker(l.progressDuration)
		defer ticker.Stop()
		for pgCnt != int32(len(l.requests)) {
			if err := func() error {
				defer progress()
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					return nil
				}
			}(); err != nil {
				return err
			}
		}
		return nil
	}))
	start = time.Now()
	var errCnt int32 = 0
	l.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		eg, egctx := errgroup.New(ctx)
		eg.Limitation(l.concurrency)
		for _, req := range l.requests {
			r := req
			eg.Go(safety.RecoverFunc(
				func() error {
					_, err := l.client.Do(egctx, l.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
						err := l.loaderFunc(ctx, vald.NewValdClient(conn), r, copts...)
						atomic.AddInt32(&pgCnt, 1)
						if err != nil {
							log.Warn(err)
							atomic.AddInt32(&errCnt, 1)
						}
						return nil, err
					})
					if err != nil {
						select {
						case <-ctx.Done():
							ech <- errors.Wrap(err, ctx.Err().Error())
						case ech <- err:
						}
					}
					return nil
				}))
		}
		if err := eg.Wait(); err != nil {
			log.Warn(err)
			select {
			case <-ctx.Done():
				ech <- errors.Wrap(err, ctx.Err().Error())
			case ech <- err:
			}
			return p.Signal(syscall.SIGKILL) // TODO: #403
		}
		log.Infof("Error ratio: %.2f%%", float64(errCnt)/float64(pgCnt)*100)
		return p.Signal(syscall.SIGTERM) // TODO: #403
	}))
	return ech
}
