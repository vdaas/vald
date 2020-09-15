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

// Package service provides meta service
package service

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
)

type Filter interface {
	Start(ctx context.Context) (<-chan error, error)
	IngressSearchVector(ctx context.Context, vec []float32) ([]float32, error)
	IngressInsertVector(ctx context.Context, vec []float32) ([]float32, error)
	IngressUpdateVector(ctx context.Context, vec []float32) ([]float32, error)
	IngressUpsertVector(ctx context.Context, vec []float32) ([]float32, error)
	IngressObject(ctx context.Context, vec []float32) ([]float32, error)
	EgressVectors(ctx context.Context, vecs []*payload.Object_Distance) ([]*payload.Object_Distance, error)
}

type filter struct {
	ingress grpc.Client
	egress  grpc.Client
	eg      errgroup.Group
}

func New(opts ...Option) (fi Filter, err error) {
	f := new(filter)
	for _, opt := range append(defaultOpts, opts...) {
		if err = opt(f); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return f, nil
}

func (f *filter) Start(ctx context.Context) (<-chan error, error) {
	if f.ingress == nil && f.egress == nil {
		return nil, nil
	}
	var (
		inch, egch <-chan error
		err        error
	)
	if f.ingress != nil {
		inch, err = f.ingress.StartConnectionMonitor(ctx)
		if err != nil {
			return nil, errors.Wrap(f.ingress.Close(), err.Error())
		}
	}
	if f.egress != nil {
		egch, err = f.egress.StartConnectionMonitor(ctx)
		if err != nil {
			if f.ingress != nil {
				err = errors.Wrap(f.ingress.Close(), err.Error())
			}
			return nil, errors.Wrap(f.egress.Close(), err.Error())
		}
	}
	ech := make(chan error, 2)
	f.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-inch:
			case err = <-egch:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))

	return ech, nil
}

func (f *filter) IngressSearchVector(ctx context.Context, vec []float32) ([]float32, error) {
	return nil, nil
}
func (f *filter) IngressInsertVector(ctx context.Context, vec []float32) ([]float32, error) {
	return nil, nil
}
func (f *filter) IngressUpdateVector(ctx context.Context, vec []float32) ([]float32, error) {
	return nil, nil
}
func (f *filter) IngressUpsertVector(ctx context.Context, vec []float32) ([]float32, error) {
	return nil, nil
}
func (f *filter) IngressObject(ctx context.Context, vec []float32) ([]float32, error) {
	return nil, nil
}
func (f *filter) EgressVectors(ctx context.Context, vecs []*payload.Object_Distance) ([]*payload.Object_Distance, error) {
	return nil, nil
}
