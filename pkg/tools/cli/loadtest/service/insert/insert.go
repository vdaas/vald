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
	"fmt"
	"github.com/vdaas/vald/internal/log"
	"reflect"
	"sync"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

type insert struct {
	eg errgroup.Group

	w client.Writer
	c int
	n string

	req []*client.ObjectVector
}

func New(opts ...InsertOption) (i *insert, err error) {
	i = new(insert)
	for _, opt := range append(defaultInsertOpts, opts...) {
		if err = opt(i); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	i.eg = errgroup.Get()
	return i, nil
}

func (i *insert) Prepare(ctx context.Context) error {
	fn := assets.Data(i.n)
	if fn == nil {
		return fmt.Errorf("dataset load funciton is nil: %s", i.n)
	}
	dataset, err := fn()
	if err != nil {
		return err
	}
	vectors := dataset.Train()
	ids := dataset.IDs()
	i.req = make([]*client.ObjectVector, len(vectors))
	for j, v := range vectors {
		i.req[j] = &client.ObjectVector{
			Id:     ids[j],
			Vector: v,
		}
	}
	return nil
}

func (i *insert) Do(ctx context.Context) <-chan error {
	errCh := make(chan error, len(i.req)*10)
	i.eg.Go(safety.RecoverFunc(func() error {
		defer close(errCh)
		wg := new(sync.WaitGroup)
		sem := make(chan struct{}, i.c)
		for _, req := range i.req {
			wg.Add(1)
			sem <- struct{}{}
			go func(r *client.ObjectVector) {
				defer wg.Done()
				defer func() {
					<-sem
				}()
				err := i.w.Insert(ctx, r)
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
