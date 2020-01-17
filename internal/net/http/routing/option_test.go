package routing

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/middleware"
)

func TestWithMiddleware(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMiddleware(tt.mw)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithMiddlewares(t *testing.T) {
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

					mws = append([]middleware.Wrapper{mw}, mws[:]...)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMiddlewares(tt.mws...)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRoute(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRoute(tt.route)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRoutes(t *testing.T) {
	type test struct {
		name      string
		routes    []Route
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			rs := []Route{
				Route{},
				Route{},
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
				Route{},
				Route{},
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

					rs = append([]Route{r}, rs[:]...)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRoutes(tt.routes...)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
