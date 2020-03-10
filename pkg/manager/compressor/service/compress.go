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
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/compress"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type Compressor interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	Compress(ctx context.Context, vector []float32) ([]byte, error)
	Decompress(ctx context.Context, bytes []byte) ([]float32, error)
	MultiCompress(ctx context.Context, vectors [][]float32) ([][]byte, error)
	MultiDecompress(ctx context.Context, bytess [][]byte) ([][]float32, error)
}

type compressor struct {
	algorithm        string
	compressionLevel int
	compressor       compress.Compressor
	limitation       int
	buffer           int
	running          atomic.Value
	eg               errgroup.Group
	jobCh            chan func() error
}

func NewCompressor(opts ...CompressorOption) (Compressor, error) {
	c := new(compressor)
	for _, opt := range append(defaultCompressorOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	c.running.Store(false)

	return c, nil
}

func (c *compressor) PreStart(ctx context.Context) (err error) {
	var compressor compress.Compressor

	switch config.CompressAlgorithm(c.algorithm) {
	case config.GOB:
		compressor, err = compress.NewGob()
	case config.GZIP:
		compressor, err = compress.NewGzip(
			compress.WithGzipCompressionLevel(c.compressionLevel),
		)
	case config.LZ4:
		compressor, err = compress.NewLZ4(
			compress.WithLZ4CompressionLevel(c.compressionLevel),
		)
	case config.ZSTD:
		compressor, err = compress.NewZstd(
			compress.WithZstdCompressionLevel(c.compressionLevel),
		)
	default:
		return errors.ErrCompressorNameNotFound(c.algorithm)
	}

	if err != nil {
		return err
	}

	c.compressor = compressor

	return nil
}

func (c *compressor) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	eg, ctx := errgroup.New(ctx)
	eg.Limitation(c.limitation)

	c.jobCh = make(chan func() error, c.buffer)

	c.running.Store(true)
	c.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(c.jobCh)
		for {
			select {
			case <-ctx.Done():
				c.running.Store(false)
				if err = ctx.Err(); err != nil {
					return errors.Wrap(eg.Wait(), err.Error())
				}
				return eg.Wait()
			case job := <-c.jobCh:
				eg.Go(safety.RecoverFunc(func() (err error) {
					err = job()
					if err != nil {
						log.Debug(err)
						runtime.Gosched()
						ech <- err
						err = nil
					}
					return nil
				}))
			}
		}
	}))

	return ech
}

func (c *compressor) dispatchCompress(ctx context.Context, vectors ...[]float32) (results [][]byte, errs error) {

	results = make([][]byte, len(vectors))

	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	c.eg.Go(safety.RecoverFunc(func() error {
		defer wg.Done()

		for iter, vector := range vectors {
			if c.running.Load().(bool) {
				wg.Add(1)
				c.jobCh <- func(i int, v []float32) func() error {
					return func() error {
						defer wg.Done()

						select {
						case <-ctx.Done():
							return ctx.Err()
						default:
						}

						res, err := c.compressor.CompressVector(v)
						if err != nil {
							mu.Lock()
							errs = errors.Wrap(errs, err.Error())
							mu.Unlock()
							return err
						}

						mu.Lock()
						results[i] = res
						mu.Unlock()

						return nil
					}
				}(iter, vector)
			}
		}

		return nil
	}))

	wg.Wait()

	for _, result := range results {
		if result == nil {
			errs = errors.Wrap(errs, errors.ErrCompressFailed.Error())
		}
	}

	return results, errs
}

func (c *compressor) dispatchDecompress(ctx context.Context, bytess ...[]byte) (results [][]float32, errs error) {
	results = make([][]float32, len(bytess))

	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	c.eg.Go(safety.RecoverFunc(func() error {
		defer wg.Done()

		for iter, bytes := range bytess {
			if c.running.Load().(bool) {
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
							mu.Lock()
							errs = errors.Wrap(errs, err.Error())
							mu.Unlock()
							return err
						}

						mu.Lock()
						results[i] = res
						mu.Unlock()

						return nil
					}
				}(iter, bytes)
			}
		}

		return nil
	}))

	wg.Wait()

	for _, result := range results {
		if result == nil {
			errs = errors.Wrap(errs, errors.ErrDecompressFailed.Error())
		}
	}

	return results, errs
}

func (c *compressor) Compress(ctx context.Context, vector []float32) ([]byte, error) {
	ress, err := c.dispatchCompress(ctx, vector)
	if err != nil {
		return nil, err
	}
	if len(ress) != 1 {
		return nil, errors.ErrCompressedDataNotFound
	}

	return ress[0], nil
}

func (c *compressor) Decompress(ctx context.Context, bytes []byte) ([]float32, error) {
	ress, err := c.dispatchDecompress(ctx, bytes)
	if err != nil {
		return nil, err
	}
	if len(ress) != 1 {
		return nil, errors.ErrDecompressedDataNotFound
	}

	return ress[0], nil
}

func (c *compressor) MultiCompress(ctx context.Context, vectors [][]float32) ([][]byte, error) {
	return c.dispatchCompress(ctx, vectors...)
}

func (c *compressor) MultiDecompress(ctx context.Context, bytess [][]byte) ([][]float32, error) {
	return c.dispatchDecompress(ctx, bytess...)
}
