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
package service

import (
	"context"
	"os"
	"reflect"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
)

// Loader is representation of load test.
type Loader interface {
	Prepare(context.Context) error
	Do(context.Context) <-chan error
}

type (
	loadFunc func(context.Context, *grpc.ClientConn, interface{}, ...grpc.CallOption) (interface{}, error)
)

type loader struct {
	eg               errgroup.Group
	client           grpc.Client
	addr             string
	concurrency      int
	batchSize        int
	dataset          string
	progressDuration time.Duration
	loaderFunc       loadFunc
	dataProvider     func() interface{}
	dataSize         int
	operation        config.Operation
}

// NewLoader returns Loader implementation.
func NewLoader(opts ...Option) (Loader, error) {
	l := new(loader)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(l); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	var err error
	switch l.operation {
	case config.Insert:
		l.loaderFunc, err = l.newInsert()
	case config.StreamInsert:
		l.loaderFunc, err = l.newStreamInsert()
	case config.Search:
		l.loaderFunc, err = l.newSearch()
	case config.StreamSearch:
		l.loaderFunc, err = l.newStreamSearch()
	default:
		err = errors.Errorf("undefined operation: %s", l.operation.String())
	}
	if err != nil {
		return nil, err
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

	switch l.operation {
	case config.Insert, config.StreamInsert:
		l.dataProvider, l.dataSize, err = insertRequestProvider(dataset, l.batchSize)
	case config.Search, config.StreamSearch:
		l.dataProvider, l.dataSize, err = searchRequestProvider(dataset)
	}
	if err != nil {
		return err
	}

	return nil
}

// Do operates load test.
func (l *loader) Do(ctx context.Context) <-chan error {
	ech := make(chan error, l.dataSize)
	finalize := func(ctx context.Context, err error) {
		select {
		case <-ctx.Done():
			ech <- errors.Wrap(err, ctx.Err().Error())
		case ech <- err:
		}
	}

	// TODO: related to #403.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		finalize(ctx, err)
		return ech
	}

	var pgCnt, errCnt int32 = 0, 0
	var start, end time.Time
	vps := func(count int, t1, t2 time.Time) float64 {
		return float64(count) / t2.Sub(t1).Seconds()
	}
	progress := func() {
		log.Infof("progress %d requests, %f[vps], error: %d", pgCnt, vps(int(pgCnt)*l.batchSize, start, time.Now()), errCnt)
	}

	f := func(i interface{}, err error) {
		atomic.AddInt32(&pgCnt, 1)
		if err != nil {
			atomic.AddInt32(&errCnt, 1)
		}
	}

	ticker := time.NewTicker(l.progressDuration)
	l.eg.Go(safety.RecoverFunc(func() error {
		defer ticker.Stop()
		for int(pgCnt) < l.dataSize {
			if err := func() error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					progress()
					return nil
				}
			}(); err != nil {
				return err
			}
		}
		return nil
	}))

	l.eg.Go(safety.RecoverFunc(func() error {
		log.Infof("start load test(%s)", l.operation.String())
		defer close(ech)
		defer ticker.Stop()
		start = time.Now()
		err := l.do(ctx, f, finalize)
		end = time.Now()

		if errCnt > 0 {
			log.Warnf("Error ratio: %.2f%%", float64(errCnt)/float64(pgCnt)*100)
		}
		if err != nil {
			finalize(ctx, err)
			return p.Signal(syscall.SIGKILL) // TODO: #403
		}
		log.Infof("result:%d\t%d\t%f", l.concurrency, l.batchSize, vps(int(pgCnt)*l.batchSize, start, end))

		return p.Signal(syscall.SIGTERM) // TODO: #403
	}))
	return ech
}

func (l *loader) do(ctx context.Context, f func(interface{}, error), notify func(context.Context, error)) (err error) {
	eg, egctx := errgroup.New(ctx)

	switch l.operation {
	case config.StreamInsert, config.StreamSearch:
		var newData func() interface{}
		switch l.operation {
		case config.StreamInsert:
			newData = func() interface{} {
				return new(payload.Empty)
			}
		case config.StreamSearch:
			newData = func() interface{} {
				return new(payload.Search_Response)
			}
		}
		eg.Go(safety.RecoverFunc(func() (err error) {
			defer func() {
				if err != nil {
					notify(egctx, err)
					err = nil
				}
			}()
			_, err = l.client.Do(egctx, l.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
				st, err := l.loaderFunc(ctx, conn, nil, copts...)
				if err != nil {
					return nil, err
				}
				return nil, grpc.BidirectionalStreamClient(st.(grpc.ClientStream), l.dataProvider, newData, f)
			})
			return err
		}))
		err = eg.Wait()
	case config.Insert, config.Search:
		eg.Limitation(l.concurrency)

		for {
			r := l.dataProvider()
			if r == nil {
				break
			}

			eg.Go(safety.RecoverFunc(func() (err error) {
				defer func() {
					notify(egctx, err)
					err = nil
				}()
				_, err = l.client.Do(egctx, l.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
					res, err := l.loaderFunc(egctx, conn, r)
					f(res, err)
					return res, err
				})

				return err
			}))
		}
		err = eg.Wait()
	default:
		err = errors.Errorf("undefined type: %s", l.operation.String())
	}
	return
}
