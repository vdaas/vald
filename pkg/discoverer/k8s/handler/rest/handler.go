//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package rest provides rest api logic
package rest

import (
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	Pods(w http.ResponseWriter, r *http.Request) (int, error)
	Nodes(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	dsc discoverer.DiscovererServer
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (*handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	data := make(map[string]interface{})
	return json.Handler(w, r, &data, func() (interface{}, error) {
		return dump.Request(nil, data, r)
	})
}

func (h *handler) Pods(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Discoverer_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.dsc.Pods(r.Context(), req)
	})
}

func (h *handler) Nodes(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Discoverer_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.dsc.Nodes(r.Context(), req)
	})
}
