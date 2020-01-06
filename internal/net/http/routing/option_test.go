package routing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/net/http/middleware"
)

func TestWithMiddleware(t *testing.T) {
	type args struct {
		mw middleware.Wrapper
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			mw := new(middlewareMock)

			return test{
				name: "set success",
				args: args{
					mw: mw,
				},
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.middlewares) != 1 {
						return fmt.Errorf("invalid params count was set")
					}

					if got, want := got.middlewares[0], mw; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("invalid params was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMiddleware(tt.args.mw)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithMiddlewares(t *testing.T) {
	type args struct {
		mws []middleware.Wrapper
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			mws := []middleware.Wrapper{
				new(middlewareMock),
				new(middlewareMock),
			}

			return test{
				name: "set success",
				args: args{
					mws: mws,
				},
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.middlewares) != 2 {
						return fmt.Errorf("invalid params count was set")
					}

					for i := range got.middlewares {
						if got, want := got.middlewares[i], mws[i]; !reflect.DeepEqual(got, want) {
							return fmt.Errorf("invalid params was set")
						}
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMiddlewares(tt.args.mws...)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRoute(t *testing.T) {
	type args struct {
		route Route
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			r := Route{}

			return test{
				name: "set success",
				args: args{
					route: r,
				},
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.routes) != 1 {
						return fmt.Errorf("invalid params count was set")
					}

					if got, want := got.routes[0], r; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("invalid params was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRoute(tt.args.route)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRoutes(t *testing.T) {
	type args struct {
		routes []Route
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			rs := []Route{
				Route{},
				Route{},
			}

			return test{
				name: "set success",
				args: args{
					routes: rs,
				},
				checkFunc: func(opt Option) error {
					got := new(router)
					opt(got)

					if len(got.routes) != 2 {
						return fmt.Errorf("invalid params count was set")
					}

					for i := range got.routes {
						if got, want := got.routes[i], rs[i]; !reflect.DeepEqual(got, want) {
							return fmt.Errorf("invalid params was set")
						}
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRoutes(tt.args.routes...)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
