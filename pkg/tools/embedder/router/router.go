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

package router

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/middleware"
	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/tools/embedder/handler/rest"
)

type router struct {
	handler rest.Handler
	eg      errgroup.Group
	timeout string
}

func New(opts ...Option) http.Handler {
	r := new(router)
	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}
	h := r.handler
	return routing.New(
		routing.WithMiddleware(middleware.NewTimeout(middleware.WithTimeout(r.timeout), middleware.WithErrorGroup(r.eg))),
		routing.WithRoutes([]routing.Route{
			{Name: "Index", Methods: []string{http.MethodGet}, Pattern: "/", HandlerFunc: h.Index},
			{Name: "Search", Methods: []string{http.MethodPost}, Pattern: "/search", HandlerFunc: h.Search},
			{Name: "LinearSearch", Methods: []string{http.MethodPost}, Pattern: "/linearsearch", HandlerFunc: h.LinearSearch},
			{Name: "Insert", Methods: []string{http.MethodPost}, Pattern: "/insert", HandlerFunc: h.Insert},
			{Name: "Insert With Metadata", Methods: []string{http.MethodPost}, Pattern: "/insert/with-metadata", HandlerFunc: h.InsertWithMetadata},
			{Name: "Update", Methods: []string{http.MethodPost}, Pattern: "/update", HandlerFunc: h.Update},
			{Name: "Update With Metadata", Methods: []string{http.MethodPost}, Pattern: "/update/with-metadata", HandlerFunc: h.UpdateWithMetadata},
			{Name: "Upsert", Methods: []string{http.MethodPost}, Pattern: "/upsert", HandlerFunc: h.Upsert},
			{Name: "Upsert With Metadata", Methods: []string{http.MethodPost}, Pattern: "/upsert/with-metadata", HandlerFunc: h.UpsertWithMetadata},
			{Name: "Remove", Methods: []string{http.MethodPost}, Pattern: "/remove", HandlerFunc: h.Remove},
			{Name: "Remove With Metadata", Methods: []string{http.MethodPost}, Pattern: "/remove/with-metadata", HandlerFunc: h.RemoveWithMetadata},
			{Name: "Embedding", Methods: []string{http.MethodPost}, Pattern: "/embedding", HandlerFunc: h.Embedding},
		}...))
}
