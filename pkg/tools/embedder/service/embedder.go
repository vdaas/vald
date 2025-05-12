// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	agent "github.com/vdaas/vald/internal/client/v1/client/agent/core"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	grpc "google.golang.org/grpc"
)

type Embedder interface {
	Insert(ctx context.Context, id, doc string) (*payload.Object_Location, error)
	Search(ctx context.Context, doc string) (*payload.Search_Response, error)
	Commit(ctx context.Context, in *payload.Control_CreateIndexRequest, opts ...grpc.CallOption) error
	Embed(ctx context.Context, doc string) (*payload.Object_Vector, error)
}

type embedder struct {
	client  vald.Client
	aclient agent.Client
	llm     LLM
}

func New(opt ...Option) (Embedder, error) {
	e := new(embedder)
	for _, o := range append(defaultOptions, opt...) {
		if err := o(e); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(o))
		}
	}
	// TODO: Validate fields
	return e, nil
}

func (e *embedder) Insert(ctx context.Context, id, doc string) (*payload.Object_Location, error) {
	vec, err := e.Embed(ctx, doc)
	if err != nil {
		return nil, err
	}
	vec.Id = id
	return e.client.Insert(ctx, &payload.Insert_Request{
		Vector: vec,
	})
}

func (e *embedder) Search(ctx context.Context, doc string) (*payload.Search_Response, error) {
	vec, err := e.Embed(ctx, doc)
	if err != nil {
		return nil, err
	}
	return e.client.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
	})
}

func (e *embedder) Commit(
	ctx context.Context, in *payload.Control_CreateIndexRequest, opts ...grpc.CallOption,
) error {
	_, err := e.aclient.CreateAndSaveIndex(ctx, in, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (e *embedder) Embed(ctx context.Context, doc string) (*payload.Object_Vector, error) {
	vec, err := e.llm.Embed(ctx, doc)
	if err != nil {
		return nil, err
	}
	return &payload.Object_Vector{
		Vector: vec,
	}, nil
}
