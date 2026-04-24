// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

	embedderpb "github.com/vdaas/vald/apis/grpc/v1/embedder"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
)

type Embedder interface {
	Search(context.Context, *embedderpb.SearchRequest) (*payload.Search_Response, error)
	LinearSearch(context.Context, *embedderpb.SearchRequest) (*payload.Search_Response, error)
	Insert(context.Context, *embedderpb.InsertRequest) (*payload.Object_Location, error)
	InsertWithMetadata(context.Context, *embedderpb.InsertWithMetadataRequest) (*payload.Object_Location, error)
	Update(context.Context, *embedderpb.UpdateRequest) (*payload.Object_Location, error)
	UpdateWithMetadata(context.Context, *embedderpb.UpdateWithMetadataRequest) (*payload.Object_Location, error)
	Upsert(context.Context, *embedderpb.UpsertRequest) (*payload.Object_Location, error)
	UpsertWithMetadata(context.Context, *embedderpb.UpsertWithMetadataRequest) (*payload.Object_Location, error)
	Remove(context.Context, *embedderpb.RemoveRequest) (*payload.Object_Location, error)
	RemoveWithMetadata(context.Context, *embedderpb.RemoveRequest) (*payload.Object_Location, error)
	Embedding(context.Context, *embedderpb.Text) (*payload.Object_Vector, error)
}

type embedder struct {
	client  vald.Client
	mclient MetaClient
	llm     LLM
}

func New(opt ...Option) (Embedder, error) {
	e := new(embedder)
	for _, o := range append(defaultOptions, opt...) {
		if err := o(e); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(o))
		}
	}
	if e.client == nil {
		return nil, errors.New("vald client is nil")
	}
	if e.llm == nil {
		return nil, errors.New("llm is nil")
	}
	return e, nil
}

func (e *embedder) Search(
	ctx context.Context, req *embedderpb.SearchRequest,
) (*payload.Search_Response, error) {
	vec, err := e.embeddingFromText(ctx, req.GetText())
	if err != nil {
		return nil, err
	}
	return e.client.Search(ctx, &payload.Search_Request{Vector: vec.GetVector(), Config: req.GetConfig()})
}

func (e *embedder) LinearSearch(
	ctx context.Context, req *embedderpb.SearchRequest,
) (*payload.Search_Response, error) {
	vec, err := e.embeddingFromText(ctx, req.GetText())
	if err != nil {
		return nil, err
	}
	return e.client.LinearSearch(ctx, &payload.Search_Request{Vector: vec.GetVector(), Config: req.GetConfig()})
}

func (e *embedder) Insert(
	ctx context.Context, req *embedderpb.InsertRequest,
) (*payload.Object_Location, error) {
	vec, err := e.vectorFromDocument(ctx, req.GetDocument())
	if err != nil {
		return nil, err
	}
	return e.client.Insert(ctx, &payload.Insert_Request{Vector: vec, Config: req.GetConfig()})
}

func (e *embedder) InsertWithMetadata(
	ctx context.Context, req *embedderpb.InsertWithMetadataRequest,
) (*payload.Object_Location, error) {
	loc, err := e.Insert(ctx, req.GetRequest())
	if err != nil {
		return nil, err
	}
	if err = e.setMetadata(ctx, req.GetRequest().GetDocument().GetId(), req.GetMetadata()); err != nil {
		return nil, err
	}
	return loc, nil
}

func (e *embedder) Update(
	ctx context.Context, req *embedderpb.UpdateRequest,
) (*payload.Object_Location, error) {
	vec, err := e.vectorFromDocument(ctx, req.GetDocument())
	if err != nil {
		return nil, err
	}
	return e.client.Update(ctx, &payload.Update_Request{Vector: vec, Config: req.GetConfig()})
}

func (e *embedder) UpdateWithMetadata(
	ctx context.Context, req *embedderpb.UpdateWithMetadataRequest,
) (*payload.Object_Location, error) {
	loc, err := e.Update(ctx, req.GetRequest())
	if err != nil {
		return nil, err
	}
	if err = e.setMetadata(ctx, req.GetRequest().GetDocument().GetId(), req.GetMetadata()); err != nil {
		return nil, err
	}
	return loc, nil
}

func (e *embedder) Upsert(
	ctx context.Context, req *embedderpb.UpsertRequest,
) (*payload.Object_Location, error) {
	vec, err := e.vectorFromDocument(ctx, req.GetDocument())
	if err != nil {
		return nil, err
	}
	return e.client.Upsert(ctx, &payload.Upsert_Request{Vector: vec, Config: req.GetConfig()})
}

func (e *embedder) UpsertWithMetadata(
	ctx context.Context, req *embedderpb.UpsertWithMetadataRequest,
) (*payload.Object_Location, error) {
	loc, err := e.Upsert(ctx, req.GetRequest())
	if err != nil {
		return nil, err
	}
	if err = e.setMetadata(ctx, req.GetRequest().GetDocument().GetId(), req.GetMetadata()); err != nil {
		return nil, err
	}
	return loc, nil
}

func (e *embedder) Remove(
	ctx context.Context, req *embedderpb.RemoveRequest,
) (*payload.Object_Location, error) {
	return e.client.Remove(ctx, &payload.Remove_Request{
		Id:     &payload.Object_ID{Id: req.GetId()},
		Config: req.GetConfig(),
	})
}

func (e *embedder) RemoveWithMetadata(
	ctx context.Context, req *embedderpb.RemoveRequest,
) (*payload.Object_Location, error) {
	loc, err := e.Remove(ctx, req)
	if err != nil {
		return nil, err
	}
	if err = e.deleteMetadata(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return loc, nil
}

func (e *embedder) Embedding(
	ctx context.Context, req *embedderpb.Text,
) (*payload.Object_Vector, error) {
	return e.embeddingFromText(ctx, req.GetText())
}

func (e *embedder) embeddingFromText(
	ctx context.Context, text string,
) (*payload.Object_Vector, error) {
	if text == "" {
		return nil, errors.New("text is empty")
	}
	vec, err := e.llm.Embed(ctx, text)
	if err != nil {
		return nil, err
	}
	return &payload.Object_Vector{Vector: vec}, nil
}

func (e *embedder) vectorFromDocument(
	ctx context.Context, doc *embedderpb.Document,
) (*payload.Object_Vector, error) {
	if doc == nil {
		return nil, errors.New("document is nil")
	}
	if doc.GetId() == "" {
		return nil, errors.New("document id is empty")
	}
	vec, err := e.embeddingFromText(ctx, doc.GetText())
	if err != nil {
		return nil, err
	}
	vec.Id = doc.GetId()
	vec.Timestamp = doc.GetTimestamp()
	return vec, nil
}

func (e *embedder) setMetadata(ctx context.Context, id string, metadata *payload.Meta_Value) error {
	if metadata == nil {
		return errors.New("metadata is nil")
	}
	if e.mclient == nil {
		return errors.New("meta client is not configured")
	}
	_, err := e.mclient.Set(ctx, &payload.Meta_KeyValue{Key: &payload.Meta_Key{Key: id}, Value: metadata})
	return err
}

func (e *embedder) deleteMetadata(ctx context.Context, id string) error {
	if e.mclient == nil {
		return errors.New("meta client is not configured")
	}
	_, err := e.mclient.Delete(ctx, &payload.Meta_Key{Key: id})
	return err
}
