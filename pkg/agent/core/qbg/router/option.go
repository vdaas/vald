//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/core/qbg/handler/rest"
)

// Option represents the functional option for router.
type Option func(*router)

var defaultOptions = []Option{
	WithTimeout("3s"),
}

// WithHandler returns the option to set the handler for the router.
func WithHandler(h rest.Handler) Option {
	return func(r *router) {
		r.handler = h
	}
}

// WithTimeout returns the option to set the timeout for the router.
func WithTimeout(timeout string) Option {
	return func(r *router) {
		r.timeout = timeout
	}
}

// WithErrGroup returns the option to set the error group for the router.
func WithErrGroup(eg errgroup.Group) Option {
	return func(r *router) {
		r.eg = eg
	}
}
