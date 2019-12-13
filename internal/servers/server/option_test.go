package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"google.golang.org/grpc"
)

func TestWithHost(t *testing.T) {
	type args struct {
		host string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				host: "host",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.host != "host" {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHost(tt.args.host)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type args struct {
		port uint
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				port: 8080,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.port != 8080 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPort(tt.args.port)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithName(t *testing.T) {
	type args struct {
		name string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				name: "vald",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.name != "vald" {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithName(tt.args.name)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithErrorGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				eg: errgroup.Get(),
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if !reflect.DeepEqual(got.eg, errgroup.Get()) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithErrorGroup(tt.args.eg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStartFunc(t *testing.T) {
	type args struct {
		fn func() error
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			fn := func() error { return nil }

			return test{
				name: "set success",
				args: args{
					fn: fn,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.preStartFunc).Pointer() != reflect.ValueOf(fn).Pointer() {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStartFunc(tt.args.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStopFunc(t *testing.T) {
	type args struct {
		fn func() error
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			fn := func() error { return nil }

			return test{
				name: "set success",
				args: args{
					fn: fn,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.preStopFunc).Pointer() != reflect.ValueOf(fn).Pointer() {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStopFunction(tt.args.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithProbeWaitTime(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "1s",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.pwt != 1*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithProbeWaitTime(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithShutdownDuration(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "1s",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.sddur != 1*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithShutdownDuration(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithReadHeaderTimeout(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "1s",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.rht != 1*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithReadHeaderTimeout(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithReadTimeout(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "1s",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.rt != 1*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithReadTimeout(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithWriteTimeout(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "1s",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.wt != 1*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithWriteTimeout(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithIdleTimeout(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "1s",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.it != 1*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithIdleTimeout(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithListenConfig(t *testing.T) {
	type args struct {
		lc *net.ListenConfig
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			lc := new(net.ListenConfig)

			return test{
				name: "set success",
				args: args{
					lc: lc,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.lc, lc) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithListenConfig(tt.args.lc)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithServerMode(t *testing.T) {
	type args struct {
		m mode
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				m: REST,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.mode != REST {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithServerMode(tt.args.m)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithTLSConfig(t *testing.T) {
	type args struct {
		cfg *tls.Config
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			cfg := new(tls.Config)

			return test{
				name: "set success",
				args: args{
					cfg: cfg,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.tcfg, cfg) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTLSConfig(tt.args.cfg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithHTTPHandler(t *testing.T) {
	type handler struct {
		http.Handler
	}

	type args struct {
		handler http.Handler
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			hdr := new(handler)

			return test{
				name: "set success",
				args: args{
					handler: hdr,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.http.h).Pointer() != reflect.ValueOf(hdr).Pointer() {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHTTPHandler(tt.args.handler)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithHTTPServer(t *testing.T) {
	type args struct {
		srv *http.Server
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			srv := new(http.Server)

			return test{
				name: "set success",
				args: args{
					srv: srv,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.http.srv, srv) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHTTPServer(tt.args.srv)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCServer(t *testing.T) {
	type args struct {
		srv *grpc.Server
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			srv := new(grpc.Server)

			return test{
				name: "set success",
				args: args{
					srv: srv,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.grpc.srv, srv) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCServer(tt.args.srv)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCOption(t *testing.T) {
	type args struct {
		opts []grpc.ServerOption
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			opts := []grpc.ServerOption{}

			return test{
				name: "set success",
				args: args{
					opts: opts,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.grpc.opts, opts) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCOption(tt.args.opts...)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCRegistFunc(t *testing.T) {
	type args struct {
		fn func(*grpc.Server)
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			fn := func(*grpc.Server) {}

			return test{
				name: "set success",
				args: args{
					fn: fn,
				},
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.grpc.reg).Pointer() != reflect.ValueOf(fn).Pointer() {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCRegistFunc(tt.args.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithEnableRestart(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.enableRestart != true {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithEnableRestart()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDisableRestart(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.enableRestart != false {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDisableRestart()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCMaxReceiveMessageSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCMaxReceiveMessageSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCMaxSendMessageSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCMaxSendMessageSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCInitialWindowSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCInitialWindowSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCInitialConnWindowSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCInitialConnWindowSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveMaxConnIdle(t *testing.T) {
	type args struct {
		max string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				max: "10m",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.maxConnIdle != 10*time.Minute {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveMaxConnIdle(tt.args.max)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveMaxConnAge(t *testing.T) {
	type args struct {
		max string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				max: "20m",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.maxConnAge != 20*time.Minute {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveMaxConnAge(tt.args.max)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveMaxConnAgeGrace(t *testing.T) {
	type args struct {
		max string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				max: "30m",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.maxConnAgeGrace != 30*time.Minute {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveMaxConnAgeGrace(tt.args.max)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveTime(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "40m",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.t != 40*time.Minute {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveTime(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveTimeout(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "50m",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.timeout != 50*time.Minute {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveTimeout(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCWriteBufferSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCWriteBufferSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCReadBufferSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCReadBufferSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCConnectionTimeout(t *testing.T) {
	type args struct {
		to string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				to: "60m",
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCConnectionTimeout(tt.args.to)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCMaxHeaderListSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCMaxHeaderListSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCHeaderTableSize(t *testing.T) {
	type args struct {
		size int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				size: 1024,
			},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCHeaderTableSize(tt.args.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCInterceptors(t *testing.T) {
	type args struct {
		names []string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{},
			checkFunc: func(opt Option) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCInterceptors(tt.args.names...)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDefaultOption(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(opts []Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opts []Option) error {
				got := new(server)

				for _, opt := range opts {
					opt(got)
				}

				if got.enableRestart != false {
					return fmt.Errorf("invalid param (enableRestart) was set")
				}

				if got.mode != REST {
					return fmt.Errorf("invalid param (mode) was set")
				}

				if !reflect.DeepEqual(got.eg, errgroup.Get()) {
					return fmt.Errorf("invalid param (eg) was set")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.checkFunc(defaultOpts); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDefaultHealthServerOption(t *testing.T) {
	type args struct {
		name string
		host string
		path string
		port uint
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opts []Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				name: "name",
				host: "host",
				path: "path",
				port: 8080,
			},
			checkFunc: func(opts []Option) error {
				got := new(server)

				for _, opt := range opts {
					opt(got)
				}

				if got.name != "name" {
					return fmt.Errorf("invalid param (name) was set")
				}

				if !reflect.DeepEqual(got.eg, errgroup.Get()) {
					return fmt.Errorf("invalid param (eg) was set")
				}

				if got.http.h == nil {
					return fmt.Errorf("invalid param (http.h) was set")
				}

				if got.host != "host" {
					return fmt.Errorf("invalid param (host) was set")
				}

				if got.it != 3*time.Second {
					return fmt.Errorf("invalid param (it) was set")
				}

				if got.port != 8080 {
					return fmt.Errorf("invalid param (port) was set")
				}

				if got.pwt != 2*time.Second {
					return fmt.Errorf("invalid param (pwt) was set")
				}

				if got.rht != 3*time.Second {
					return fmt.Errorf("invalid param (rht) was set")
				}

				if got.rt != 2*time.Second {
					return fmt.Errorf("invalid param (rt) was set")
				}

				if got.mode != REST {
					return fmt.Errorf("invalid param (mode) was set")
				}

				if got.sddur != 4*time.Second {
					return fmt.Errorf("invalid param (sddur) was set")
				}

				if got.wt != 3*time.Second {
					return fmt.Errorf("invalid param (wt) was set")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := HealthServerOpts(tt.args.name, tt.args.host, tt.args.path, tt.args.port)
			if err := tt.checkFunc(opts); err != nil {
				t.Error(err)
			}
		})
	}
}
