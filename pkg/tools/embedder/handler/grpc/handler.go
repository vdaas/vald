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
	"context"
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/embedder"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/tools/embedder/service"
)

type Server interface {
	embedder.EmbedderServer
}

type server struct {
	embedder service.Embedder
	embedder.UnimplementedEmbedderServer
}

func New(opts ...Option) (Server, error) {
	s := new(server)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(s); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}
	return s, nil
}

func (s *server) Search(
	ctx context.Context, req *embedder.SearchRequest,
) (*payload.Search_Response, error) {
	return s.embedder.Search(ctx, req)
}

func (s *server) LinearSearch(
	ctx context.Context, req *embedder.SearchRequest,
) (*payload.Search_Response, error) {
	return s.embedder.LinearSearch(ctx, req)
}

func (s *server) Insert(
	ctx context.Context, req *embedder.InsertRequest,
) (*payload.Object_Location, error) {
	return s.embedder.Insert(ctx, req)
}

func (s *server) InsertWithMetadata(
	ctx context.Context, req *embedder.InsertWithMetadataRequest,
) (*payload.Object_Location, error) {
	return s.embedder.InsertWithMetadata(ctx, req)
}

func (s *server) Update(
	ctx context.Context, req *embedder.UpdateRequest,
) (*payload.Object_Location, error) {
	return s.embedder.Update(ctx, req)
}

func (s *server) UpdateWithMetadata(
	ctx context.Context, req *embedder.UpdateWithMetadataRequest,
) (*payload.Object_Location, error) {
	return s.embedder.UpdateWithMetadata(ctx, req)
}

func (s *server) Upsert(
	ctx context.Context, req *embedder.UpsertRequest,
) (*payload.Object_Location, error) {
	return s.embedder.Upsert(ctx, req)
}

func (s *server) UpsertWithMetadata(
	ctx context.Context, req *embedder.UpsertWithMetadataRequest,
) (*payload.Object_Location, error) {
	return s.embedder.UpsertWithMetadata(ctx, req)
}

func (s *server) Remove(
	ctx context.Context, req *embedder.RemoveRequest,
) (*payload.Object_Location, error) {
	return s.embedder.Remove(ctx, req)
}

func (s *server) RemoveWithMetadata(
	ctx context.Context, req *embedder.RemoveRequest,
) (*payload.Object_Location, error) {
	return s.embedder.RemoveWithMetadata(ctx, req)
}

func (s *server) Embedding(
	ctx context.Context, req *embedder.Text,
) (*payload.Object_Vector, error) {
	return s.embedder.Embedding(ctx, req)
}
