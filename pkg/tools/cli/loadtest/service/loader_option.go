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
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Option func(*loader) error

var (
	defaultOpts = []Option{
		WithConcurrency(100),
		WithErrGroup(errgroup.Get()),
		WithProgressDuration(5 * time.Second),
	}
)

func WithAddr(a string) Option {
	return func(l *loader) error {
		l.addr = a
		return nil
	}
}

func WithClient(c grpc.Client) Option {
	return func(l *loader) error {
		if c != nil {
			l.client = c
			return nil
		}
		return errors.Errorf("client must not be nil")
	}
}

func WithConcurrency(c int) Option {
	return func(l *loader) error {
		if c > 0 {
			l.concurrency = c
		}
		return nil
	}
}

func WithDataset(n string) Option {
	return func(l *loader) error {
		l.dataset = n
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(l *loader) error {
		if eg != nil {
			l.eg = eg
		}
		return nil
	}
}

func WithProgressDuration(pd time.Duration) Option {
	return func(l *loader) error {
		l.progressDuration = pd
		return nil
	}
}
