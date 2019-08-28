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

// Package routing provides implementation of Go API for routing http Handler wrapped by rest.Func
package routing

import (
	"time"

	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*router)

var (
	defaultOpts = []Option{
		WithTimeout("3s"),
	}
)

func WithTimeout(timeout string) Option {
	return func(r *router) {
		var err error
		r.timeout, err = timeutil.Parse(timeout)
		if err != nil {
			r.timeout = time.Second * 3
		}
	}
}

func WithRoute(route Route) Option {
	return func(r *router) {
		r.routes = append(r.routes, route)
	}
}

func WithRoutes(routes ...Route) Option {
	return func(r *router) {
		r.routes = routes
	}
}
