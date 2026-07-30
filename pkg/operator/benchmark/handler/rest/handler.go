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

package rest

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
	"github.com/vdaas/vald/pkg/operator/benchmark/handler/grpc"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	ScenarioStatus(w http.ResponseWriter, r *http.Request) (int, error)
	JobStatus(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	bm grpc.Benchmark
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOpts, opts...) {
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

func (h *handler) ScenarioStatus(w http.ResponseWriter, r *http.Request) (int, error) {
	var req any
	return json.Handler(w, r, &req, func() (any, error) {
		return h.bm.GetScenarioStatus(), nil
	})
}

func (h *handler) JobStatus(w http.ResponseWriter, r *http.Request) (int, error) {
	var req any
	return json.Handler(w, r, &req, func() (any, error) {
		return h.bm.GetBenchmarkJobStatus(), nil
	})
}
