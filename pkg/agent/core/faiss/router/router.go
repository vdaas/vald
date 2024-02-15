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

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"net/http"

<<<<<<< HEAD
	"github.com/vdaas/vald/internal/net/http/middleware"
	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/internal/sync/errgroup"
=======
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/http/middleware"
	"github.com/vdaas/vald/internal/net/http/routing"
>>>>>>> feature/agent/qbg
	"github.com/vdaas/vald/pkg/agent/core/faiss/handler/rest"
)

type router struct {
	handler rest.Handler
	eg      errgroup.Group
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
		routing.WithMiddleware(
			middleware.NewTimeout(
				middleware.WithTimeout(r.timeout),
				middleware.WithErrorGroup(r.eg),
			)),
		routing.WithRoutes([]routing.Route{
			{
				"Index",
				[]string{
					http.MethodGet,
				},
				"/",
				h.Index,
			},
			{
				"Search",
				[]string{
					http.MethodPost,
				},
				"/search",
				h.Search,
			},
			{
				"Search By ID",
				[]string{
					http.MethodPost,
				},
				"/id/search",
				h.SearchByID,
			},
			{
				"LinearSearch",
				[]string{
					http.MethodPost,
				},
				"/linearsearch",
				h.LinearSearch,
			},
			{
				"LinearSearch By ID",
				[]string{
					http.MethodPost,
				},
				"/id/linearsearch",
				h.LinearSearchByID,
			},
			{
				"Insert",
				[]string{
					http.MethodPost,
				},
				"/insert",
				h.Insert,
			},
			{
				"Multiple Insert",
				[]string{
					http.MethodPost,
				},
				"/insert/multi",
				h.MultiInsert,
			},
			{
				"Update",
				[]string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				"/update",
				h.Update,
			},
			{
				"Multiple Update",
				[]string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				"/update/multi",
				h.MultiUpdate,
			},
			{
				"Remove",
				[]string{
					http.MethodDelete,
				},
				"/delete",
				h.Remove,
			},
			{
				"Multiple Remove",
				[]string{
					http.MethodDelete,
					http.MethodPost,
				},
				"/delete/multi",
				h.MultiRemove,
			},
			{
				"Create Index",
				[]string{
					http.MethodPost,
				},
				"/index/create",
				h.CreateIndex,
			},
			{
				"Save Index",
				[]string{
					http.MethodGet,
				},
				"/index/save",
				h.SaveIndex,
			},
			{
				"GetObject",
				[]string{
					http.MethodGet,
				},
				"/object/{id}",
				h.GetObject,
			},
		}...))
}
