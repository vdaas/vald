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

// Package metrics provides pprof profiler handler
package metrics

import (
	"net/http"
	"net/http/pprof"

	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/net/http/routing"
)

// NewPProfRoutes returns PProf server route&method information from debug flag
func NewPProfHandler() http.Handler {
	return routing.New(
		routing.WithTimeout("5s"),
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
		}...))
}
