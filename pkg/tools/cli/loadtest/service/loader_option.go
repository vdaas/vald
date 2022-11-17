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
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
)

// Option is load test configuration.
type Option func(*loader) error

var defaultOptions = []Option{
	WithConcurrency(100),
	WithBatchSize(1),
	WithErrGroup(errgroup.Get()),
	WithProgressDuration("5s"),
}

// WithAddr sets load test server address.
func WithAddr(a string) Option {
	return func(l *loader) error {
		l.addr = a
		return nil
	}
}

// WithClient sets grpc Client.
func WithClient(c grpc.Client) Option {
	return func(l *loader) error {
		if c != nil {
			l.client = c
			return nil
		}
		return errors.New("client must not be nil")
	}
}

// WithConcurrency sets load test concurrency.
func WithConcurrency(c int) Option {
	return func(l *loader) error {
		if c <= 0 {
			return errors.New("concurrency must be natural number")
		}
		l.concurrency = c
		return nil
	}
}

// WithBatchSize sets load test batch size.
func WithBatchSize(b int) Option {
	return func(l *loader) error {
		if b <= 0 {
			return errors.New("batch size must be natural number")
		}
		l.batchSize = b
		return nil
	}
}

// WithDataset sets dataset name.
func WithDataset(n string) Option {
	return func(l *loader) error {
		l.dataset = n
		return nil
	}
}

// WithErrGroup sets user specified error group.
func WithErrGroup(eg errgroup.Group) Option {
	return func(l *loader) error {
		if eg != nil {
			l.eg = eg
		}
		return nil
	}
}

// WithProgressDuration sets duration of progress show.
func WithProgressDuration(pd string) Option {
	return func(l *loader) error {
		if pd == "" {
			return nil
		}
		d, err := timeutil.Parse(pd)
		if err != nil {
			return err
		}
		l.progressDuration = d
		return nil
	}
}

// WithOperation sets operation of load test.
func WithOperation(op string) Option {
	return func(l *loader) error {
		l.operation = config.OperationMethod(op)
		return nil
	}
}
