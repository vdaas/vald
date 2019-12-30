//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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
	"reflect"
	"sync"

	"github.com/vdaas/vald/internal/compress"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
)

type Compressor interface {
	Start(ctx context.Context) <-chan error
	Compress(ctx context.Context, vector []float64) ([]byte, error)
	Decompress(ctx context.Context, bytes []byte) ([]float64, error)
	MultiCompress(ctx context.Context, vectors [][]float64) ([][]byte, error)
	MultiDecompress(ctx context.Context, bytess [][]byte) ([][]float64, error)
}

type compressor struct {
	compressor compress.Compressor
	limitation int
	buffer     int
	eg         errgroup.Group
	jobCh      chan func() error
}

func NewCompressor(opts ...CompressorOption) (Compressor, error) {
	c := new(compressor)
	for _, opt := range append(defaultCompressorOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

func (c *compressor) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	eg, ctx := errgroup.New(ctx)
	eg.Limitation(c.limitation)

	c.jobCh = make(chan func() error, c.buffer)

	c.eg.Go(safety.RecoverFunc(func() (err error) {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case j := <-c.jobCh:
				eg.Go(safety.RecoverFunc(j))
			}
		}
	}))

	return ech
}

func (c *compressor) dispatchCompress(ctx context.Context, vectors ...[]float64) (results [][]byte, errs error) {
	results = make([][]byte, len(vectors))

	wg := new(sync.WaitGroup)
	wg.Add(1)
	c.eg.Go(safety.RecoverFunc(func() error {
		defer wg.Done()

		for iter, vector := range vectors {
			wg.Add(1)
			c.jobCh <- func(i int, v []float64) func() error {
				return func() error {
					defer wg.Done()

					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}

					res, err := c.compressor.CompressVector(v)
					if err != nil {
						errs = errors.Wrap(errs, err.Error())
						return err
					}

					results[i] = res

					return nil
				}
			}(iter, vector)
		}

		return nil
	}))

	wg.Wait()

	for _, result := range results {
		if result == nil {
			errs = errors.Wrap(errs, errors.ErrCompressFailed().Error())
		}
	}

	return results, errs
}

func (c *compressor) dispatchDecompress(ctx context.Context, bytess ...[]byte) (results [][]float64, errs error) {
	results = make([][]float64, len(bytess))

	wg := new(sync.WaitGroup)
	wg.Add(1)
	c.eg.Go(safety.RecoverFunc(func() error {
		defer wg.Done()

		for iter, bytes := range bytess {
			wg.Add(1)
			c.jobCh <- func(i int, b []byte) func() error {
				return func() error {
					defer wg.Done()

					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}

					res, err := c.compressor.DecompressVector(b)
					if err != nil {
						errs = errors.Wrap(errs, err.Error())
						return err
					}

					results[i] = res

					return nil
				}
			}(iter, bytes)
		}

		return nil
	}))

	wg.Wait()

	for _, result := range results {
		if result == nil {
			errs = errors.Wrap(errs, errors.ErrDecompressFailed().Error())
		}
	}

	return results, errs
}

func (c *compressor) Compress(ctx context.Context, vector []float64) ([]byte, error) {
	ress, err := c.dispatchCompress(ctx, vector)
	if len(ress) != 1 {
		return nil, err
	}

	return ress[0], err
}

func (c *compressor) Decompress(ctx context.Context, bytes []byte) ([]float64, error) {
	ress, err := c.dispatchDecompress(ctx, bytes)
	if len(ress) != 1 {
		return nil, err
	}

	return ress[0], err
}

func (c *compressor) MultiCompress(ctx context.Context, vectors [][]float64) ([][]byte, error) {
	return c.dispatchCompress(ctx, vectors...)
}

func (c *compressor) MultiDecompress(ctx context.Context, bytess [][]byte) ([][]float64, error) {
	return c.dispatchDecompress(ctx, bytess...)
}
