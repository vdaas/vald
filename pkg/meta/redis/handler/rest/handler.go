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

// Package rest provides rest api logic
package rest

import (
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/meta"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	GetMeta(w http.ResponseWriter, r *http.Request) (int, error)
	GetMetas(w http.ResponseWriter, r *http.Request) (int, error)
	GetMetaInverse(w http.ResponseWriter, r *http.Request) (int, error)
	GetMetasInverse(w http.ResponseWriter, r *http.Request) (int, error)
	SetMeta(w http.ResponseWriter, r *http.Request) (int, error)
	SetMetas(w http.ResponseWriter, r *http.Request) (int, error)
	DeleteMeta(w http.ResponseWriter, r *http.Request) (int, error)
	DeleteMetas(w http.ResponseWriter, r *http.Request) (int, error)
	DeleteMetaInverse(w http.ResponseWriter, r *http.Request) (int, error)
	DeleteMetasInverse(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	meta meta.MetaServer
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	data := make(map[string]interface{})
	return json.Handler(w, r, &data, func() (interface{}, error) {
		return dump.Request(nil, data, r)
	})
}

func (h *handler) GetMeta(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Key)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.GetMeta(r.Context(), req)
	})
}

func (h *handler) GetMetas(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Keys)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.GetMetas(r.Context(), req)
	})
}

func (h *handler) GetMetaInverse(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Val)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.GetMetaInverse(r.Context(), req)
	})
}

func (h *handler) GetMetasInverse(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Vals)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.GetMetasInverse(r.Context(), req)
	})
}

func (h *handler) SetMeta(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_KeyVal)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.SetMeta(r.Context(), req)
	})
}

func (h *handler) SetMetas(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_KeyVals)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.SetMetas(r.Context(), req)
	})
}

func (h *handler) DeleteMeta(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Key)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.DeleteMeta(r.Context(), req)
	})
}

func (h *handler) DeleteMetas(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Keys)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.DeleteMetas(r.Context(), req)
	})
}

func (h *handler) DeleteMetaInverse(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Val)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.DeleteMetaInverse(r.Context(), req)
	})
}

func (h *handler) DeleteMetasInverse(w http.ResponseWriter, r *http.Request) (int, error) {
	req := new(payload.Meta_Vals)
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.meta.DeleteMetasInverse(r.Context(), req)
	})
}
