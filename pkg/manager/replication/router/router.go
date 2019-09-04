// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/pkg/agent/ngt/handler/rest"
)

type router struct {
	handler rest.Handler
	timeout string
}

// New returns REST route&method information from handler interface
func New(opts ...Option) http.Handler {

	r := new(router)

	for _, opt := range append(defaultOpts, opts...) {
		opt(r)
	}

	h := r.handler

	return routing.New(
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
					http.MethodGet,
				},
				"/search/{id}",
				h.SearchByID,
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
				"/delete/{id}",
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
					http.MethodGet,
				},
				"/index/create/{pool}",
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
		}...),
		routing.WithTimeout(r.timeout))
}
