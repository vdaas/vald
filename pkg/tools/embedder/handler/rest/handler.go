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

package rest

import (
	"net/http"

	embedderpb "github.com/vdaas/vald/apis/grpc/v1/embedder"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
	"github.com/vdaas/vald/pkg/tools/embedder/handler/grpc"
)

type Handler interface {
	Index(http.ResponseWriter, *http.Request) (int, error)
	Exists(http.ResponseWriter, *http.Request) (int, error)
	Search(http.ResponseWriter, *http.Request) (int, error)
	LinearSearch(http.ResponseWriter, *http.Request) (int, error)
	Insert(http.ResponseWriter, *http.Request) (int, error)
	InsertWithMetadata(http.ResponseWriter, *http.Request) (int, error)
	Update(http.ResponseWriter, *http.Request) (int, error)
	UpdateWithMetadata(http.ResponseWriter, *http.Request) (int, error)
	Upsert(http.ResponseWriter, *http.Request) (int, error)
	UpsertWithMetadata(http.ResponseWriter, *http.Request) (int, error)
	Remove(http.ResponseWriter, *http.Request) (int, error)
	RemoveWithMetadata(http.ResponseWriter, *http.Request) (int, error)
	Embedding(http.ResponseWriter, *http.Request) (int, error)
}

type handler struct {
	embedder grpc.Server
}

func New(opts ...Option) Handler {
	h := new(handler)
	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (*handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	data := make(map[string]any)
	return json.Handler(w, r, &data, func() (any, error) {
		return dump.Request(nil, data, r)
	})
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.SearchRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.Search(r.Context(), req) })
}

func (h *handler) LinearSearch(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.SearchRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.LinearSearch(r.Context(), req) })
}

func (h *handler) Insert(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.InsertRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.Insert(r.Context(), req) })
}

func (h *handler) InsertWithMetadata(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.InsertWithMetadataRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.InsertWithMetadata(r.Context(), req) })
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.UpdateRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.Update(r.Context(), req) })
}

func (h *handler) UpdateWithMetadata(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.UpdateWithMetadataRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.UpdateWithMetadata(r.Context(), req) })
}

func (h *handler) Upsert(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.UpsertRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.Upsert(r.Context(), req) })
}

func (h *handler) UpsertWithMetadata(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.UpsertWithMetadataRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.UpsertWithMetadata(r.Context(), req) })
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.RemoveRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.Remove(r.Context(), req) })
}

func (h *handler) RemoveWithMetadata(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.RemoveRequest
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.RemoveWithMetadata(r.Context(), req) })
}

func (h *handler) Embedding(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *embedderpb.Text
	return json.Handler(w, r, &req, func() (any, error) { return h.embedder.Embedding(r.Context(), req) })
}

func (*handler) Exists(w http.ResponseWriter, r *http.Request) (int, error) {
	return json.Handler(w, r, &payload.Empty{}, func() (any, error) { return map[string]any{"ok": true}, nil })
}
