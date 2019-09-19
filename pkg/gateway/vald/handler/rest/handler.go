//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"net/http"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	Exists(w http.ResponseWriter, r *http.Request) (int, error)
	Search(w http.ResponseWriter, r *http.Request) (int, error)
	SearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	Insert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiInsert(w http.ResponseWriter, r *http.Request) (int, error)
	Update(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpdate(w http.ResponseWriter, r *http.Request) (int, error)
	Remove(w http.ResponseWriter, r *http.Request) (int, error)
	MultiRemove(w http.ResponseWriter, r *http.Request) (int, error)
	GetObject(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	agent agent.AgentServer
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOpts, opts...) {
		opt(h)
	}
	return h
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprint(w, r.URL.String())
	return http.StatusOK, nil
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.Search(r.Context(), req)
	})
}

func (h *handler) SearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.SearchByID(r.Context(), req)
	})
}

func (h *handler) Insert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_Vector
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.Insert(r.Context(), req)
	})
}

func (h *handler) MultiInsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_Vectors
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.MultiInsert(r.Context(), req)
	})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_Vector
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.Update(r.Context(), req)
	})
}

func (h *handler) MultiUpdate(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_Vectors
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.MultiUpdate(r.Context(), req)
	})
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_ID
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.Remove(r.Context(), req)
	})
}

func (h *handler) MultiRemove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_IDs
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.MultiRemove(r.Context(), req)
	})
}

func (h *handler) GetObject(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_ID
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.GetObject(r.Context(), req)
	})
}

func (h *handler) Exists(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_ID
	return json.Handler(w, r, req, func() (interface{}, error) {
		return h.agent.Exists(r.Context(), req)
	})
}
