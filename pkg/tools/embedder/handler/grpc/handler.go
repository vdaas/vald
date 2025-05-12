//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/embedder"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/tools/embedder/service"
)

type Server interface {
	// Insert(context.Context, *payload.Object_Blob) (*payload.Object_Location, error)
	// Search(context.Context, *payload.Object_Blob) (*payload.Search_Response, error)
	// Embedding(context.Context, *payload.Object_Blob) (*payload.Object_Vector, error)
	// Commit(context.Context, *payload.Empty) (*payload.Empty, error)
	embedder.EmbedderServer
}

type server struct {
	name     string
	ip       string
	embedder service.Embedder
	embedder.UnimplementedEmbedderServer
}

const (
	apiName = "vald/tools/embedder"
)

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
