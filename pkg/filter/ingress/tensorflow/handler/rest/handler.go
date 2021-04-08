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

	"github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Handler interface {
	GenVector(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	ingress ingress.FilterServer
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (h *handler) GenVector(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Object_Blob
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.ingress.GenVector(r.Context(), req)
	})
}
