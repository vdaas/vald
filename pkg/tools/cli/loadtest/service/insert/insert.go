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
	"reflect"
	"sync"

	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errors"
)

type insert struct {
	w client.Writer
	p int
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
	return i, nil
}

func (i *insert) PreStart(ctx context.Context) error {
	dataset, err := assets.Data(i.n)()
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

func (i *insert) Start(ctx context.Context) (<-chan error, error) {
	errCh := make(chan error, len(i.req)*10)
	go func() {
		wg := new(sync.WaitGroup)
		limCh := make(chan struct{}, i.p)
		for _, req := range i.req {
			wg.Add(1)
			limCh <- struct{}{}
			go func(r *client.ObjectVector) {
				defer wg.Done()
				defer func() {
					<-limCh
				}()
				err := i.w.Insert(ctx, r)
				if err != nil {
					errCh <- err
				}
			}(req)
		}
		wg.Wait()
		close(errCh)
	}()
	return errCh, nil
}
