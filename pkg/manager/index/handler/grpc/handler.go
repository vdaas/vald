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

// Package grpc provides grpc server logic
package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/manager/index"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/manager/index/service"
)

type server struct {
	indexer service.Indexer
}

func New(opts ...Option) index.IndexServer {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) IndexInfo(ctx context.Context, _ *payload.Empty) (res *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-index.IndexInfo")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return &payload.Info_Index_Count{
		Stored:      s.indexer.NumberOfUUIDs(),
		Uncommitted: s.indexer.NumberOfUncommittedUUIDs(),
		Indexing:    s.indexer.IsIndexing(),
	}, nil
}
