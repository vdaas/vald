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
package search

import (
	"context"
	"fmt"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
	"reflect"
)

type search struct {
	eg          errgroup.Group
	client      grpc.Client
	addr        string
	concurrency int
	dataset     string
	req         []*payload.Search_Request
}

func New(opts ...Option) (s *search, err error) {
	s = new(search)
	for _, opt := range append(defaultOpts, opts...) {
		if err = opt(s); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	s.eg = errgroup.Get()
	return s, nil
}

func (s *search) Prepare(ctx context.Context) error {
	fn := assets.Data(s.dataset)
	if fn == nil {
		return fmt.Errorf("dataset load function is nil: %s", s.dataset)
	}
	dataset, err := assets.Data(s.dataset)()
	if err != nil {
		return err
	}
	vectors := dataset.Query()
	s.req = make([]*payload.Search_Request, len(vectors))
	for i, v := range vectors {
		s.req[i] = &payload.Search_Request{
			Vector: v,
		}
	}

	return nil
}

func (s *search) Do(ctx context.Context) <-chan error {
	errCh := make(chan error, len(s.req))
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(errCh)
		eg, egctx := errgroup.New(ctx)
		eg.Limitation(s.concurrency)
		for _, req := range s.req {
			r := req
			eg.Go(func() error {
				_, err := s.client.Do(egctx, s.addr, func(Ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
					return vald.NewValdClient(conn).Search(ctx, r, copts...)
				})
				if err != nil {
					errCh <- err
				}
				return nil
			})
		}
		return eg.Wait()
	}))
	return errCh
}
