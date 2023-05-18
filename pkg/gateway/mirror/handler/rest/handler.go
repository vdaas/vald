// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package rest

import (
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Handler interface {
	Register(w http.ResponseWriter, r *http.Request) (int, error)
	Advertise(w http.ResponseWriter, r *http.Request) (int, error)
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	Exists(w http.ResponseWriter, r *http.Request) (int, error)
	Search(w http.ResponseWriter, r *http.Request) (int, error)
	SearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearch(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearch(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearch(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	Insert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiInsert(w http.ResponseWriter, r *http.Request) (int, error)
	Update(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpdate(w http.ResponseWriter, r *http.Request) (int, error)
	Upsert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpsert(w http.ResponseWriter, r *http.Request) (int, error)
	Remove(w http.ResponseWriter, r *http.Request) (int, error)
	MultiRemove(w http.ResponseWriter, r *http.Request) (int, error)
	GetObject(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	vald vald.ServerWithMirror
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Mirror_Targets
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Register(r.Context(), req)
	})
}

func (h *handler) Advertise(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Mirror_Targets
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Advertise(r.Context(), req)
	})
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	data := make(map[string]interface{})
	return json.Handler(w, r, &data, func() (interface{}, error) {
		return dump.Request(nil, data, r)
	})
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Search(r.Context(), req)
	})
}

func (h *handler) SearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.SearchByID(r.Context(), req)
	})
}

func (h *handler) MultiSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiSearch(r.Context(), req)
	})
}

func (h *handler) MultiSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiIDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiSearchByID(r.Context(), req)
	})
}

func (h *handler) LinearSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.LinearSearch(r.Context(), req)
	})
}

func (h *handler) LinearSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.LinearSearchByID(r.Context(), req)
	})
}

func (h *handler) MultiLinearSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiLinearSearch(r.Context(), req)
	})
}

func (h *handler) MultiLinearSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiIDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiLinearSearchByID(r.Context(), req)
	})
}

func (h *handler) Insert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Insert(r.Context(), req)
	})
}

func (h *handler) MultiInsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiInsert(r.Context(), req)
	})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Update(r.Context(), req)
	})
}

func (h *handler) MultiUpdate(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiUpdate(r.Context(), req)
	})
}

func (h *handler) Upsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Upsert_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Upsert(r.Context(), req)
	})
}

func (h *handler) MultiUpsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Upsert_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiUpsert(r.Context(), req)
	})
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Remove(r.Context(), req)
	})
}

func (h *handler) MultiRemove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiRemove(r.Context(), req)
	})
}

func (h *handler) GetObject(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_VectorRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.GetObject(r.Context(), req)
	})
}

func (h *handler) Exists(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_ID
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Exists(r.Context(), req)
	})
}
