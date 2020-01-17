package errgroup

import (
	"context"
	"reflect"
	sync "sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func clearGlobalObject() {
	instance = nil
	once = sync.Once{}
}

func TestNew(t *testing.T) {
	type test struct {
		name string
		ctx  context.Context
	}

	tests := []test{
		{
			name: "returns eg and context",
			ctx:  context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eg, ctx := New(tt.ctx)
			if eg == nil || ctx == nil {
				t.Errorf("eg or ctx is nil. eg: %v, ctx: %v", eg, ctx)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type test struct {
		name       string
		ctx        context.Context
		beforeFunc func()
		checkFunc  func(ctx context.Context) error
	}

	tests := []test{
		{
			name: "returns egctx when New function is called",
			ctx:  context.Background(),
			beforeFunc: func() {
				clearGlobalObject()
			},
			checkFunc: func(egctx context.Context) error {
				if egctx == nil || instance == nil {
					return errors.Errorf("egctx or instance is nil. egctx: %v, instance: %v", egctx, instance)
				}
				return nil
			},
		},

		func() test {
			ctx := context.Background()

			return test{
				name: "returns ctx of argument when instance is already initialized",
				ctx:  ctx,
				beforeFunc: func() {
					instance = new(group)
					once = sync.Once{}
					once.Do(func() {})
				},
				checkFunc: func(egctx context.Context) error {
					if !reflect.DeepEqual(egctx, ctx) {
						return errors.Errorf("egctx is not equals. want: %v, got: %v", ctx, egctx)
					}

					if instance == nil {
						return errors.New("instance is nil")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				clearGlobalObject()
			}()

			if tt.beforeFunc != nil {
				tt.beforeFunc()
			}

			egctx := Init(tt.ctx)
			if err := tt.checkFunc(egctx); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type test struct {
		name       string
		beforeFunc func()
	}

	tests := []test{
		{
			name: "returns new instance when instance is nil",
			beforeFunc: func() {
				clearGlobalObject()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				clearGlobalObject()
			}()

			if tt.beforeFunc != nil {
				tt.beforeFunc()
			}

			if eg := Get(); eg == nil {
				t.Error("eg is nil")
			}
		})
	}
}

func TestLimitation(t *testing.T) {
	type args struct {
		limit int
	}

	type field struct {
		limitation       chan struct{}
		enableLimitation atomic.Value
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(eg *group) error
	}

	tests := []test{
		{
			name: "processing is successes when limitation is greater than 0",
			args: args{
				limit: 10,
			},
			field: field{
				limitation: make(chan struct{}),
			},
			checkFunc: func(eg *group) error {
				if ok := eg.enableLimitation.Load().(bool); !ok {
					return errors.Errorf("enableLimitation is wrong. want: %v, got: %v", true, !ok)
				}
				return nil
			},
		},

		{
			name: "processing is successes when limitation is 0 or less",
			args: args{
				limit: 0,
			},
			checkFunc: func(eg *group) error {
				if ok := eg.enableLimitation.Load().(bool); ok {
					return errors.Errorf("enableLimitation is wrong. want: %v, got: %v", false, ok)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eg := &group{
				limitation:       tt.field.limitation,
				enableLimitation: tt.field.enableLimitation,
			}
			eg.Limitation(tt.args.limit)

			if err := tt.checkFunc(eg); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEgWait(t *testing.T) {
	type receiver struct {
		g Group
	}

	type test struct {
		name      string
		receiver  receiver
		checkFunc func() error
		want      error
	}

	tests := []test{
		func() test {
			var enableLimitation atomic.Value
			enableLimitation.Store(true)

			err1 := errors.New("fail_1")
			err2 := errors.New("fail_2")

			g := &group{
				enableLimitation: enableLimitation,
				egctx:            context.Background(),
				limitation:       make(chan struct{}),
				errs: []error{
					err1,
					err2,
				},
			}

			return test{
				name: "processing is successes when function is not nil",
				receiver: receiver{
					g: g,
				},
				checkFunc: func() error {
					if ok := g.enableLimitation.Load().(bool); ok {
						return errors.Errorf("enableLimitation is not equals. want: %v, got: %v", false, ok)
					}
					return nil
				},
				want: errors.Wrap(
					errors.Wrap(
						g.err,
						err1.Error(),
					),
					err2.Error(),
				),
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.receiver.g.Wait()
			if !errors.Is(tt.want, got) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}

			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWait(t *testing.T) {
	type global struct {
		instance Group
	}

	type test struct {
		name      string
		global    global
		checkFunc func() error
		want      error
	}

	tests := []test{
		func() test {
			var enableLimitation atomic.Value
			enableLimitation.Store(true)

			g := &group{
				enableLimitation: enableLimitation,
				egctx:            context.Background(),
				limitation:       make(chan struct{}),
			}

			return test{
				name: "processing is successes when errors is nil",
				global: global{
					instance: g,
				},
				checkFunc: func() error {
					if ok := g.enableLimitation.Load().(bool); ok {
						return errors.Errorf("enableLimitation is not equals. want: %v, got: %v", false, ok)
					}
					return nil
				},
				want: nil,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				clearGlobalObject()
			}()

			instance = tt.global.instance

			got := Wait()
			if !errors.Is(tt.want, got) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}

			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDoCancel(t *testing.T) {
	type receiver struct {
		g *group
	}

	type test struct {
		name      string
		receiver  receiver
		afterFunc func()
		checkFunc func() error
	}

	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())

			g := &group{
				cancel: cancel,
			}

			return test{
				name: "processing is successes when cancel is not nil",
				receiver: receiver{
					g: g,
				},
				checkFunc: func() error {
					if !errors.Is(context.Canceled, ctx.Err()) {
						return errors.Errorf("context error is wrong. want: %v, got: %v", context.Canceled, ctx.Err())
					}
					return nil
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())

			return test{
				name: "do nothing when cancel is nil",
				receiver: receiver{
					g: new(group),
				},
				checkFunc: func() error {
					if err := ctx.Err(); err != nil {
						return errors.Errorf("context error is not nil: %v", err)
					}
					return nil
				},
				afterFunc: func() {
					cancel()
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.afterFunc != nil {
					tt.afterFunc()
				}
			}()

			tt.receiver.g.doCancel()
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEgGo(t *testing.T) {
	type args struct {
		f func() error
	}

	type receiver struct {
		g Group
	}

	type test struct {
		name      string
		args      args
		receiver  receiver
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var enableLimitation atomic.Value
			enableLimitation.Store(true)

			err := errors.New("fail")
			f := func() error {
				return err
			}

			g := &group{
				enableLimitation: enableLimitation,
				egctx:            context.Background(),
				emap:             make(map[string]struct{}),
				limitation:       make(chan struct{}, 1),
			}

			return test{
				name: "processing is successes when function is not nil",
				args: args{
					f: f,
				},
				receiver: receiver{
					g: g,
				},
				checkFunc: func() error {
					g.Wait()

					if len(g.emap) != 1 {
						return errors.Errorf("emap count is wrong. want: %d, got: %d", 1, len(g.emap))
					}

					if _, ok := g.emap[err.Error()]; !ok {
						return errors.Errorf("%s is not contains into the emap", err.Error())
					}

					if len(g.errs) != 1 {
						return errors.Errorf("errs count is wrong. want: %d, got: %d", 1, len(g.errs))
					}

					if !errors.Is(g.errs[0], err) {
						return errors.Errorf("errs[0] is not equals. want: %v, got: %v", g.errs[0], err)
					}

					return nil
				},
			}
		}(),

		{
			name: "do nothing when function is nil",
			receiver: receiver{
				g: new(group),
			},
			checkFunc: func() error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.receiver.g.Go(tt.args.f)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGo(t *testing.T) {
	type args struct {
		f func() error
	}

	type global struct {
		instance Group
	}

	type test struct {
		name   string
		args   args
		global global
	}

	tests := []test{
		{
			name: "do nothing when function is nil",
			global: global{
				instance: new(group),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				clearGlobalObject()
			}()

			instance = tt.global.instance
			Go(tt.args.f)
		})
	}
}
