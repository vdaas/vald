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
package middleware

import (
	"time"

	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type TimeoutOption func(*timeout)

var defaultTimeoutOpts = []TimeoutOption{
	WithTimeout("3s"),
	WithErrorGroup(errgroup.Get()),
}

func WithTimeout(dur string) TimeoutOption {
	return func(t *timeout) {
		var err error
		t.dur, err = timeutil.Parse(dur)
		if err != nil {
			t.dur = time.Second * 3
		}
	}
}

func WithErrorGroup(eg errgroup.Group) TimeoutOption {
	return func(t *timeout) {
		t.eg = eg
	}
}
