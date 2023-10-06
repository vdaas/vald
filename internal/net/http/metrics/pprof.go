//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package metrics provides pprof profiler handler
package metrics

import (
	"net/http"
	"net/http/pprof"

	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/net/http/routing"
)

// NewPProfRoutes returns PProf server route&method information from debug flag.
func NewPProfHandler() http.Handler {
	return routing.New(
		routing.WithRoutes([]routing.Route{
			{
				Name: "Debug pprof",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Index),
			},
			{
				Name: "Debug cmdline",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/cmdline",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Cmdline),
			},
			{
				Name: "Debug profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/profile",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Profile),
			},
			{
				Name: "Debug symbol profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/symbol",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Symbol),
			},
			{
				Name: "Debug trace profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/trace",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Trace),
			},
			{
				Name: "Debug allocs profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/allocs",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Handler("allocs").ServeHTTP),
			},
			{
				Name: "Debug heap profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/heap",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Handler("heap").ServeHTTP),
			},
			{
				Name: "Debug goroutine profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/goroutine",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Handler("goroutine").ServeHTTP),
			},
			{
				Name: "Debug thread profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/threadcreate",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Handler("threadcreate").ServeHTTP),
			},
			{
				Name: "Debug block profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/block",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Handler("block").ServeHTTP),
			},

			{
				Name: "Debug mutex profile",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/debug/pprof/mutex",
				HandlerFunc: rest.HandlerToRestFunc(pprof.Handler("mutex").ServeHTTP),
			},
		}...))
}
