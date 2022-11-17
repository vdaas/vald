// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package routing

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/middleware"
)

func TestWithMiddleware(t *testing.T) {
	t.Parallel()
	type test struct {
		name      string
		mw        middleware.Wrapper
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			mw := new(middlewareMock)

			return test{
				name: "set success",
				mw:   mw,
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.middlewares) != 1 {
						return errors.New("invalid params count was set")
					}

					if got, want := got.middlewares[0], mw; !reflect.DeepEqual(got, want) {
						return errors.New("invalid params was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			opt := WithMiddleware(test.mw)
			if err := test.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithMiddlewares(t *testing.T) {
	t.Parallel()
	type test struct {
		name      string
		mws       []middleware.Wrapper
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			mws := []middleware.Wrapper{
				new(middlewareMock),
				new(middlewareMock),
			}

			return test{
				name: "set success when middlewares field is nil",
				mws:  mws,
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.middlewares) != 2 {
						return errors.New("invalid params count was set")
					}

					for i := range got.middlewares {
						if got, want := got.middlewares[i], mws[i]; !reflect.DeepEqual(got, want) {
							return errors.New("invalid params was set")
						}
					}

					return nil
				},
			}
		}(),

		func() test {
			mw := new(middlewareMock)

			mws := []middleware.Wrapper{
				new(middlewareMock),
				new(middlewareMock),
			}

			return test{
				name: "set success when middlewares field is not nil",
				mws:  mws,
				checkFunc: func(opt Option) error {
					got := &router{
						middlewares: []middleware.Wrapper{
							mw,
						},
					}
					opt(got)

					if len(got.middlewares) != 3 {
						return errors.New("invalid params count was set")
					}

					mws := append([]middleware.Wrapper{mw}, mws...)
					for i := range got.middlewares {
						if got, want := got.middlewares[i], mws[i]; !reflect.DeepEqual(got, want) {
							return errors.New("invalid params was set")
						}
					}

					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			opt := WithMiddlewares(test.mws...)
			if err := test.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRoute(t *testing.T) {
	t.Parallel()
	type test struct {
		name      string
		route     Route
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			r := Route{}

			return test{
				name:  "set success",
				route: r,
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.routes) != 1 {
						return errors.New("invalid params count was set")
					}

					if got, want := got.routes[0], r; !reflect.DeepEqual(got, want) {
						return errors.New("invalid params was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			opt := WithRoute(test.route)
			if err := test.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRoutes(t *testing.T) {
	t.Parallel()
	type test struct {
		name      string
		routes    []Route
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			rs := []Route{
				{},
				{},
			}

			return test{
				name:   "set success when routes field is nil",
				routes: rs,
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.routes) != 2 {
						return errors.New("invalid params count was set")
					}

					for i := range got.routes {
						if got, want := got.routes[i], rs[i]; !reflect.DeepEqual(got, want) {
							return errors.New("invalid params was set")
						}
					}

					return nil
				},
			}
		}(),

		func() test {
			r := Route{}

			rs := []Route{
				{},
				{},
			}

			return test{
				name:   "set success when routes field is not nil",
				routes: rs,
				checkFunc: func(opt Option) error {
					got := &router{
						routes: []Route{
							r,
						},
					}
					opt(got)

					if len(got.routes) != 3 {
						return errors.New("invalid params count was set")
					}

					rs := append([]Route{r}, rs...)
					for i := range got.routes {
						if got, want := got.routes[i], rs[i]; !reflect.DeepEqual(got, want) {
							return errors.New("invalid params was set")
						}
					}

					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			opt := WithRoutes(test.routes...)
			if err := test.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
