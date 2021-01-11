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
				"Debug pprof",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/",
				rest.HandlerToRestFunc(pprof.Index),
			},
			{
				"Debug cmdline",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/cmdline",
				rest.HandlerToRestFunc(pprof.Cmdline),
			},
			{
				"Debug profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/profile",
				rest.HandlerToRestFunc(pprof.Profile),
			},
			{
				"Debug symbol profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/symbol",
				rest.HandlerToRestFunc(pprof.Symbol),
			},
			{
				"Debug trace profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/trace",
				rest.HandlerToRestFunc(pprof.Trace),
			},
			{
				"Debug allocs profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/allocs",
				rest.HandlerToRestFunc(pprof.Handler("allocs").ServeHTTP),
			},
			{
				"Debug heap profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/heap",
				rest.HandlerToRestFunc(pprof.Handler("heap").ServeHTTP),
			},
			{
				"Debug goroutine profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/goroutine",
				rest.HandlerToRestFunc(pprof.Handler("goroutine").ServeHTTP),
			},
			{
				"Debug thread profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/threadcreate",
				rest.HandlerToRestFunc(pprof.Handler("threadcreate").ServeHTTP),
			},
			{
				"Debug block profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/block",
				rest.HandlerToRestFunc(pprof.Handler("block").ServeHTTP),
			},

			{
				"Debug mutex profile",
				[]string{
					http.MethodGet,
				},
				"/debug/pprof/mutex",
				rest.HandlerToRestFunc(pprof.Handler("mutex").ServeHTTP),
			},
		}...))
}
