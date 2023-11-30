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

// Package transport provides http transport roundtrip option
package transport

import (
	"net/http"

	"github.com/vdaas/vald/internal/backoff"
)

// Option represents the functional option for transport and backoff.
type Option func(*ert)

var defaultOptions = []Option{
	WithRoundTripper(http.DefaultTransport),
}

// WithRoundTripper returns the Option that set the RoundTripper.
func WithRoundTripper(tr http.RoundTripper) Option {
	return func(e *ert) {
		e.transport = tr
	}
}

// WithBackoff returns the Option that set the backoff.
func WithBackoff(bo backoff.Backoff) Option {
	return func(e *ert) {
		e.bo = bo
	}
}
