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
	"github.com/vdaas/vald/internal/net/grpc"
)

type Option func(*insert) error

var (
	defaultOpts = []Option{
		WithConcurrency(100),
	}
)

func WithAddr(a string) Option {
	return func(i *insert) error {
		i.addr = a
		return nil
	}
}

func WithClient(c grpc.Client) Option {
	return func(i *insert) error {
		i.client = c
		return nil
	}
}

func WithConcurrency(c int) Option {
	return func(i *insert) error {
		i.concurrency = c
		return nil
	}
}

func WithDataset(n string) Option {
	return func(i *insert) error {
		i.dataset = n
		return nil
	}
}
