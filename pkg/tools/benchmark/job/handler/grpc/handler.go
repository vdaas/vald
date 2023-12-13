//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"

	"github.com/vdaas/vald/pkg/tools/benchmark/job/service"
)

type Benchmark interface {
	Start(context.Context)
}

type server struct {
	job service.Job
}

func New(opts ...Option) (bm Benchmark, err error) {
	b := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		err = opt(b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (s *server) Start(ctx context.Context) {
}
