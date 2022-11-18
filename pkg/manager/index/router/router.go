//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/pkg/manager/index/handler/rest"
)

type router struct {
	handler rest.Handler
	timeout string
}

// New returns REST route&method information from handler interface.
func New(opts ...Option) http.Handler {
	r := new(router)

	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

	h := r.handler

	return routing.New(
		routing.WithRoutes([]routing.Route{
			{
				Name: "Index",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/",
				HandlerFunc: h.Index,
			},
			{
				Name: "IndexInfo",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/index",
				HandlerFunc: h.IndexInfo,
			},
		}...))
}
