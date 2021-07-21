//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/rebalancer"
	"github.com/vdaas/vald/internal/singleflight"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/service/job"
)

type Rebalancer interface {
	rebalancer.JobServer
	Start(context.Context)
}

type server struct {
	rb    job.Rebalancer
	group singleflight.Group
}

func New(opts ...Option) (rb Rebalancer, err error) {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		err = opt(s)
		if err != nil {
			return nil, err
		}
	}

	s.group = singleflight.New()

	return s, nil
}

func (s *server) Start(ctx context.Context) {
}
