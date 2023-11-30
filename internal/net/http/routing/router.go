//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package routing provides implementation of Go API for routing http Handler wrapped by rest.Func
package routing

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/json"
	"github.com/vdaas/vald/internal/net/http/middleware"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/strings"
)

type router struct {
	middlewares []middleware.Wrapper
	routes      []Route
}

// New returns Routed http.Handler.
func New(opts ...Option) http.Handler {
	r := new(router)
	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 32

	rt := mux.NewRouter().StrictSlash(true)
	for _, route := range r.routes {
		for _, mw := range r.middlewares {
			route.HandlerFunc = mw.Wrap(route.HandlerFunc)
		}

		rt.Handle(route.Pattern,
			r.routing(route.Name, route.Pattern,
				route.Methods, route.HandlerFunc)).Name(route.Name)
	}

	return rt
}

// routing wraps the handler.Func and returns a new http.Handler.
// routing helps to handle unsupported HTTP method, timeout,
// and the error returned from the handler.Func.
func (*router) routing(
	name, _ string, m []string, h rest.Func,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var (
				err  error
				code int
			)

			for _, method := range m {
				if strings.EqualFold(r.Method, method) {
					// execute only if the request method is inside the method list
					code, err = h(w, r)
					if err != nil && code != http.StatusServiceUnavailable {
						err = json.ErrorHandler(w, r,
							err.Error()+" at handler "+name,
							code,
							err)
						if err != nil {
							log.Error(err)
						}
					}
					return
				}
			}

			err = json.ErrorHandler(w, r,
				"Invalid Request Method",
				http.StatusMethodNotAllowed,
				errors.ErrInvalidRequest)
			if err != nil {
				log.Error(err)
			}
		})
}
