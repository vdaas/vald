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
package insert

import (
	"context"
	"os"
	"reflect"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

type insert struct {
	eg          errgroup.Group
	client      grpc.Client
	addr        string
	concurrency int
	dataset     string
	req         []*payload.Object_Vector
}

func New(opts ...Option) (i *insert, err error) {
	i = new(insert)
	for _, opt := range append(defaultOpts, opts...) {
		if err = opt(i); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	i.eg = errgroup.Get()
	return i, nil
}

func (i *insert) Prepare(ctx context.Context) error {
	fn := assets.Data(i.dataset)
	if fn == nil {
		return errors.Errorf("dataset load funciton is nil: %s", i.dataset)
	}
	dataset, err := fn()
	if err != nil {
		return err
	}
	vectors := dataset.Train()
	ids := dataset.IDs()
	i.req = make([]*payload.Object_Vector, len(vectors))
	for j, v := range vectors {
		i.req[j] = &payload.Object_Vector{
			Id:     ids[j],
			Vector: v,
		}
	}
	return nil
}

func (i *insert) Do(ctx context.Context) <-chan error {
	ech := make(chan error, len(i.req))

	// TODO: related to #403.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		ech <- err
		return ech
	}

	var progress int32 = 0
	ticker := time.NewTicker(5 * time.Second)
	i.eg.Go(safety.RecoverFunc(func() error {
		for progress != int32(len(i.req)) {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				log.Debugf("insert %d items", progress)
			}
		}
		return nil
	}))
	i.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		eg, egctx := errgroup.New(ctx)
		eg.Limitation(i.concurrency)
		for _, req := range i.req {
			r := req
			eg.Go(func() error {
				_, err := i.client.Do(egctx, i.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
					_, err := vald.NewValdClient(conn).Insert(ctx, r, copts...)
					atomic.AddInt32(&progress, 1)
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
		if err != nil {
			log.Warn(err)
			ech <- err
			return p.Signal(syscall.SIGKILL) // TODO: #403
		}
		return p.Signal(syscall.SIGTERM) // TODO: #403
	}))
	return ech
}
