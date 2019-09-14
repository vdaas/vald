//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
