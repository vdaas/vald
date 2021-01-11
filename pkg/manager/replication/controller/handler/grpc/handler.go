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

	"github.com/vdaas/vald/apis/grpc/v1/manager/replication/controller"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/pkg/manager/replication/controller/service"
)

type server struct {
	controller service.Replicator
}

func New(opts ...Option) controller.ReplicationControllerServer {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) ReplicationInfo(ctx context.Context, _ *payload.Empty) (res *payload.Replication_Agents, err error) {
	return nil, nil
}
