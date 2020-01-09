package servers

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/servers/server"
)

func TestNew(t *testing.T) {
	type test struct {
		name      string
		opts      []Option
		checkFunc func(got, want *listener) error
		want      *listener
	}

	tests := []test{
		{
			name: "initialize with default options",
			want: &listener{
				eg: errgroup.Get(),
			},
			checkFunc: func(got *listener, want *listener) error {
				if !reflect.DeepEqual(got, want) {
					return errors.Errorf("not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},
		{
			name: "initialize with custom options",
			opts: []Option{
				WithStartUpStrategy([]string{
					"strg_1",
					"strg_2",
				}),
			},
			want: &listener{
				eg: errgroup.Get(),
				sus: []string{
					"strg_1",
					"strg_2",
				},
				sds: []string{
					"strg_2",
					"strg_1",
				},
			},
			checkFunc: func(got *listener, want *listener) error {
				if !reflect.DeepEqual(got.sus, want.sus) {
					return errors.Errorf("sus is not equals. want: %v, got: %v", want.sus, got.sus)
				}

				if !reflect.DeepEqual(got.sds, want.sds) {
					return errors.Errorf("sds is not equals. want: %v, got: %v", want.sds, got.sds)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.opts...)
			if err := tt.checkFunc(got.(*listener), tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestListenAndServe(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type field struct {
		eg      errgroup.Group
		servers map[string]server.Server
		sus     []string
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(got, want <-chan error) error
		want      <-chan error
	}

	tests := []test{
		func() test {
			srv1 := &mockServer{
				IsRunningFunc: func() bool {
					return false
				},
				ListenAndServeFunc: func(context.Context, chan<- error) error {
					return nil
				},
			}

			srv2 := &mockServer{
				IsRunningFunc: func() bool {
					return false
				},
				ListenAndServeFunc: func(context.Context, chan<- error) error {
					return nil
				},
			}

			srv3Err := errors.New("srv3 error")
			srv3 := &mockServer{
				IsRunningFunc: func() bool {
					return false
				},
				ListenAndServeFunc: func(context.Context, chan<- error) error {
					return srv3Err
				},
			}

			servers := map[string]server.Server{
				"srv1": srv1,
				"srv2": srv2,
				"srv3": srv3,
			}

			sus := []string{
				"srv1",
				"srv2",
				"srv3",
				"srv4",
			}

			errCh := make(chan error, len(servers))
			errCh <- srv3Err
			errCh <- errors.ErrServerNotFound("srv4")
			close(errCh)

			return test{
				name: "ListenAndServe is success",
				args: args{
					ctx: func() context.Context {
						ctx, cancel := context.WithCancel(context.Background())
						defer cancel()
						return ctx
					}(),
				},
				field: field{
					eg:      errgroup.Get(),
					servers: servers,
					sus:     sus,
				},
				want: errCh,
				checkFunc: func(got <-chan error, want <-chan error) error {
					gerrs := make([]error, 0, len(servers))
					for err := range got {
						gerrs = append(gerrs, err)
					}

					werrs := make([]error, 0, len(servers))
					for err := range want {
						werrs = append(werrs, err)
					}

					fmt.Println(gerrs)
					fmt.Println(werrs)

					if len(werrs) != len(gerrs) {
						return errors.Errorf("errors count is not equals: want: %v, got: %v", len(werrs), len(gerrs))
					}

					for i := range werrs {
						if gerrs[i].Error() != werrs[i].Error() {
							return errors.Errorf("errors[%d] is not equals: want: %v, got: %v", i, werrs[i], gerrs[i])
						}
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()

			l := &listener{
				eg:      tt.field.eg,
				servers: tt.field.servers,
				sus:     tt.field.sus,
			}

			errCh := l.ListenAndServe(ctx)
			if err := tt.checkFunc(errCh, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestShutdown(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type field struct {
		eg      errgroup.Group
		servers map[string]server.Server
		sdr     []string
		sddur   time.Duration
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(got, want error) error
		want      error
	}

	tests := []test{
		func() test {
			srv1 := NewMockServer()
			srv1.IsRunningFunc = func() bool {
				return true
			}
			srv1.ShutdownFunc = func(context.Context) error {
				return nil
			}

			srv2 := NewMockServer()
			srv2.IsRunningFunc = func() bool {
				return true
			}
			srv2.ShutdownFunc = func(context.Context) error {
				return nil
			}

			servers := map[string]server.Server{
				"srv1": srv1,
				"srv2": srv2,
			}

			sdr := []string{
				"srv1", "srv2",
			}

			return test{
				name: "Shutdown is success",
				args: args{
					ctx: context.Background(),
				},
				field: field{
					eg:      errgroup.Get(),
					servers: servers,
					sdr:     sdr,
				},

				checkFunc: func(got error, want error) error {
					if got != nil {
						return fmt.Errorf("return error: %v", got)
					}
					return nil
				},
				want: nil,
			}
		}(),
		{
			name: "server not found error",
			args: args{
				ctx: context.Background(),
			},
			field: field{
				eg:      errgroup.Get(),
				servers: map[string]server.Server{},
				sdr: []string{
					"srv1",
				},
			},
			checkFunc: func(got error, want error) error {
				if got.Error() != want.Error() {
					return fmt.Errorf("not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
			want: errors.ErrServerNotFound("srv1"),
		},
		func() test {
			want := errors.Wrap(fmt.Errorf("unexpected error"), "faild to shutdown")

			srv1 := NewMockServer()
			srv1.IsRunningFunc = func() bool {
				return true
			}
			srv1.ShutdownFunc = func(context.Context) error {
				return want
			}

			servers := map[string]server.Server{
				"srv1": srv1,
			}

			sdr := []string{
				"srv1",
			}

			return test{
				name: "unexpected error",
				args: args{
					ctx: context.Background(),
				},
				field: field{
					eg:      errgroup.Get(),
					servers: servers,
					sdr:     sdr,
				},
				checkFunc: func(got error, want error) error {
					if got.Error() != want.Error() {
						return fmt.Errorf("not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
				want: want,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()

			l := &listener{
				eg:      tt.field.eg,
				servers: tt.field.servers,
				sds:     tt.field.sdr,
				sddur:   tt.field.sddur,
			}

			err := l.Shutdown(ctx)
			if err := tt.checkFunc(err, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}
