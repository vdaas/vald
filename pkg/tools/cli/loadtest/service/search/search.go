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
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/safety"
	"reflect"
	"sync"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

type search struct {
	eg  errgroup.Group
	r   client.Reader
	c   int
	n   string
	req []*client.SearchRequest
}

func New(opts ...SearchOption) (s *search, err error) {
	s = new(search)
	for _, opt := range append(defaultSearchOpts, opts...) {
		if err = opt(s); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	s.eg = errgroup.Get()
	return s, nil
}

func (s *search) Prepare(ctx context.Context) error {
	dataset, err := assets.Data(s.n)()
	if err != nil {
		return err
	}
	vectors := dataset.Query()
	s.req = make([]*client.SearchRequest, len(vectors))
	for i, v := range vectors {
		s.req[i] = &client.SearchRequest{
			Vector: v,
		}
	}

	return nil
}

func (s *search) Do(ctx context.Context) <-chan error {
	errCh := make(chan error, len(s.req)*10)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(errCh)
		wg := new(sync.WaitGroup)
		sem := make(chan struct{}, s.c)
		for _, req := range s.req {
			wg.Add(1)
			sem <- struct{}{}
			go func(r *client.SearchRequest) {
				defer wg.Done()
				defer func() {
					<-sem
				}()
				_, err := s.r.Search(ctx, r)
				if err != nil {
					errCh <- err
				}
			}(req)
		}
		wg.Wait()
		return nil
	}))
	return errCh
}
