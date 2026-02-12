//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package grpc

import (
	"time"

	stats "github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type Server interface {
	vald.Server
	stats.StatsDetailServer
}

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway
	timeout           time.Duration
	replica           int
	streamConcurrency int
	multiConcurrency  int
	name              string
	ip                string
	vald.UnimplementedValdServer
	stats.UnimplementedStatsDetailServer
}

const apiName = "vald/gateway/lb"

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}
