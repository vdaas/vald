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

	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/pkg/gateway/meta/handler/rest"
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
				Name: "Search",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search",
				HandlerFunc: h.Search,
			},
			{
				Name: "Search With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search/metadata",
				HandlerFunc: h.SearchWithMetadata,
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
				Name: "Search By ID With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search/id/metadata",
				HandlerFunc: h.SearchByIDWithMetadata,
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
				Name: "Multi Search With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search/multiple/metadata",
				HandlerFunc: h.MultiSearchWithMetadata,
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
				Name: "Multi Search By ID With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search/id/multiple/metadata",
				HandlerFunc: h.MultiSearchByIDWithMetadata,
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
				Name: "Linear Search With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/linearsearch/metadata",
				HandlerFunc: h.LinearSearchWithMetadata,
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
				Name: "Linear Search By ID With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/linearsearch/id/metadata",
				HandlerFunc: h.LinearSearchByIDWithMetadata,
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
				Name: "Multi Linear Search With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/linearsearch/multiple/metadata",
				HandlerFunc: h.MultiLinearSearchWithMetadata,
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
				Name: "Multi Linear Search By ID With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/linearsearch/id/multiple/metadata",
				HandlerFunc: h.MultiLinearSearchByIDWithMetadata,
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
				Name: "Insert With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/insert/metadata",
				HandlerFunc: h.InsertWithMetadata,
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
				Name: "Multiple Insert With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/insert/multiple/metadata",
				HandlerFunc: h.MultiInsertWithMetadata,
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
				Name: "Update With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/update/metadata",
				HandlerFunc: h.UpdateWithMetadata,
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
				Name: "Multiple Update With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/update/multiple/metadata",
				HandlerFunc: h.MultiUpdateWithMetadata,
			},
			{
				Name: "Update Timestamp With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/update/timestamp/metadata",
				HandlerFunc: h.UpdateTimestampWithMetadata,
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
				Name: "Upsert With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/upsert/metadata",
				HandlerFunc: h.UpsertWithMetadata,
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
				Name: "Multiple Upsert With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/upsert/multiple/metadata",
				HandlerFunc: h.MultiUpsertWithMetadata,
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
				Name: "Remove With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/remove/metadata",
				HandlerFunc: h.RemoveWithMetadata,
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
				Name: "Multiple Remove With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/remove/multiple/metadata",
				HandlerFunc: h.MultiRemoveWithMetadata,
			},
			{
				Name: "Remove By Timestamp With Metadata",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/remove/timestamp/metadata",
				HandlerFunc: h.RemoveByTimestampWithMetadata,
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
				Name: "StreamListObject With Metadata",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/object/list/metadata",
				HandlerFunc: h.StreamListObjectWithMetadata,
			},
			{
				Name: "GetObject With Metadata",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/object/{id}/metadata",
				HandlerFunc: h.GetObjectWithMetadata,
			},
		}...))
}
