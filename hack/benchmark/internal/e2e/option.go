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

// Package e2e provides e2e testing framework functions
package e2e

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client/v1/client"
)

type Option func(*e2e)

var defaultOptions = []Option{
	WithServerStarter(
		func(context.Context, testing.TB, assets.Dataset) func() {
			return func() {}
		},
	),
}

func WithName(name string) Option {
	return func(e *e2e) {
		if len(name) != 0 {
			e.name = name
		}
	}
}

func WithClient(c client.Client) Option {
	return func(e *e2e) {
		if c != nil {
			e.client = c
		}
	}
}

func WithStrategy(strategis ...Strategy) Option {
	return func(e *e2e) {
		if len(strategis) != 0 {
			e.strategies = strategis
		}
	}
}

func WithServerStarter(f func(context.Context, testing.TB, assets.Dataset) func()) Option {
	return func(e *e2e) {
		if f != nil {
			e.srvStarter = f
		}
	}
}
