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
package test

import (
	"context"
	"testing"
)

type Test interface {
	Run(context.Context, *testing.T)
}

type test struct {
	cases    []Caser
	driverFn func(context.Context, DataProvider) []interface{}
}

func New(opts ...Option) Test {
	t := new(test)
	for _, opt := range append(defaultOptions, opts...) {
		opt(t)
	}
	return t
}

func (test *test) Run(ctx context.Context, t *testing.T) {
	t.Helper()
	for _, c := range test.cases {
		t.Run(c.Name(), func(tt *testing.T) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			if fn := c.FieldFunc(); fn != nil {
				c.SetField(fn(tt)...)
			}

			gots := test.driverFn(ctx, c)

			if err := c.AssertFunc()(gots, c.Wants()); err != nil {
				tt.Errorf("assertion error: %v", err)
			}
		})
	}
}
