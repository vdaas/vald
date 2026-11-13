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

package router

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/pkg/gateway/lb/handler/rest"
)

type router struct {
	handler rest.Handler
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
				Name: "Search",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search",
				HandlerFunc: h.Search,
			},

			{
				Name: "Multi Search",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search/multi",
				HandlerFunc: h.MultiSearch,
			},
			{
				Name: "Search By ID",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/search/{id}",
				HandlerFunc: h.SearchByID,
			},
			{
				Name: "Multi Search By ID",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/search/multi/{id}",
				HandlerFunc: h.MultiSearchByID,
			},
			{
				Name: "Linear_Search",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/linearsearch",
				HandlerFunc: h.LinearSearch,
			},
			{
				Name: "Multi Linear_Search",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/linearsearch/multi",
				HandlerFunc: h.MultiLinearSearch,
			},
			{
				Name: "Linear_Search By ID",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/linearsearch/{id}",
				HandlerFunc: h.LinearSearchByID,
			},
			{
				Name: "Multi Linear_Search By ID",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/linearsearch/multi/{id}",
				HandlerFunc: h.MultiLinearSearchByID,
			},
			{
				Name: "Insert",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/insert",
				HandlerFunc: h.Insert,
			},
			{
				Name: "Multiple Insert",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/insert/multi",
				HandlerFunc: h.MultiInsert,
			},
			{
				Name: "Update",
				Methods: []string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				Pattern:     "/update",
				HandlerFunc: h.Update,
			},
			{
				Name: "Multiple Update",
				Methods: []string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				Pattern:     "/update/multi",
				HandlerFunc: h.MultiUpdate,
			},
			{
				Name: "Upsert",
				Methods: []string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				Pattern:     "/upsert",
				HandlerFunc: h.Upsert,
			},
			{
				Name: "Multiple Upsert",
				Methods: []string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				Pattern:     "/upsert/multi",
				HandlerFunc: h.MultiUpsert,
			},
			{
				Name: "Multiple Remove",
				Methods: []string{
					http.MethodDelete,
					http.MethodPost,
				},
				Pattern:     "/delete/multi",
				HandlerFunc: h.MultiRemove,
			},
			{
				Name: "Remove",
				Methods: []string{
					http.MethodDelete,
				},
				Pattern:     "/delete/{id}",
				HandlerFunc: h.Remove,
			},
			{
				Name: "Flush",
				Methods: []string{
					http.MethodDelete,
				},
				Pattern:     "/flush",
				HandlerFunc: h.Flush,
			},
			{
				Name: "GetObject",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/object/{id}",
				HandlerFunc: h.GetObject,
			},
			{
				Name: "Exists",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/exists/{id}",
				HandlerFunc: h.Exists,
			},
			{
				Name: "GetTimestamp",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/object/meta/{id}",
				HandlerFunc: h.GetTimestamp,
			},
			{
				Name: "UpdateTimestamp",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/update/timestamp",
				HandlerFunc: h.UpdateTimestamp,
			},
			{
				Name: "RemoveByTimestamp",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/remove/timestamp",
				HandlerFunc: h.RemoveByTimestamp,
			},
			{
				Name: "IndexInfo",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/index/info",
				HandlerFunc: h.IndexInfo,
			},
			{
				Name: "IndexDetail",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/index/detail",
				HandlerFunc: h.IndexDetail,
			},
			{
				Name: "IndexStatistics",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/index/statistics",
				HandlerFunc: h.IndexStatistics,
			},
			{
				Name: "IndexStatisticsDetail",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/index/statistics/detail",
				HandlerFunc: h.IndexStatisticsDetail,
			},
			{
				Name: "IndexProperty",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/index/property",
				HandlerFunc: h.IndexProperty,
			},
		}...))
}
