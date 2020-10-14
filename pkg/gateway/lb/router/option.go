//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"github.com/vdaas/vald/pkg/gateway/lb/handler/rest"
)

type Option func(*router)

var (
	defaultOpts = []Option{
		WithTimeout("3s"),
	}
)

func WithHandler(h rest.Handler) Option {
	return func(r *router) {
		r.handler = h
	}
}

func WithTimeout(timeout string) Option {
	return func(r *router) {
		r.timeout = timeout
	}
}
