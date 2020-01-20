package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"google.golang.org/grpc"
)

func TestWithHost(t *testing.T) {
	type test struct {
		name      string
		host      string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			host: "host",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.host != "host" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing",
			host: "",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.host != "" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHost(tt.host)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type test struct {
		name      string
		port      uint
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			port: 8080,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.port != 8080 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.port != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPort(tt.port)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithName(t *testing.T) {
	type test struct {
		name      string
		arg       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			arg:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.name != "vald" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithName(tt.arg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithErrorGroup(t *testing.T) {
	type test struct {
		name      string
		eg        errgroup.Group
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			eg:   errgroup.Get(),
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if !reflect.DeepEqual(got.eg, errgroup.Get()) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithErrorGroup(tt.eg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStartFunc(t *testing.T) {
	type test struct {
		name      string
		fn        func() error
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			fn := func() error { return nil }

			return test{
				name: "set success",
				fn:   fn,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.preStartFunc).Pointer() != reflect.ValueOf(fn).Pointer() {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStartFunc(tt.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStopFunc(t *testing.T) {
	type test struct {
		name      string
		fn        func() error
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			fn := func() error { return nil }

			return test{
				name: "set success",
				fn:   fn,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.preStopFunc).Pointer() != reflect.ValueOf(fn).Pointer() {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStopFunction(tt.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithProbeWaitTime(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "1s",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.pwt != 1*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.pwt != 5*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithProbeWaitTime(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithShutdownDuration(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "1s",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.sddur != 1*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.sddur != 5*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithShutdownDuration(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithReadHeaderTimeout(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "1s",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.rht != 1*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.rht != 5*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithReadHeaderTimeout(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithReadTimeout(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "1s",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.rt != 1*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.rt != 5*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithReadTimeout(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithWriteTimeout(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "1s",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.wt != 1*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.wt != 5*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithWriteTimeout(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithIdleTimeout(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "1s",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.it != 1*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.it != 5*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithIdleTimeout(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithListenConfig(t *testing.T) {
	type test struct {
		name      string
		lc        *net.ListenConfig
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			lc := new(net.ListenConfig)

			return test{
				name: "set success",
				lc:   lc,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.lc, lc) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithListenConfig(tt.lc)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithServerMode(t *testing.T) {
	type test struct {
		name      string
		m         mode
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			m:    REST,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.mode != REST {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithServerMode(tt.m)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithTLSConfig(t *testing.T) {
	type test struct {
		name      string
		cfg       *tls.Config
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			cfg := new(tls.Config)

			return test{
				name: "set success",
				cfg:  cfg,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.tcfg, cfg) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTLSConfig(tt.cfg)
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

	type test struct {
		name      string
		handler   http.Handler
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			hdr := new(handler)

			return test{
				name:    "set success",
				handler: hdr,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.http.h).Pointer() != reflect.ValueOf(hdr).Pointer() {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHTTPHandler(tt.handler)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithHTTPServer(t *testing.T) {
	type test struct {
		name      string
		srv       *http.Server
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			srv := new(http.Server)

			return test{
				name: "set success",
				srv:  srv,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.http.srv, srv) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHTTPServer(tt.srv)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCServer(t *testing.T) {
	type test struct {
		name      string
		srv       *grpc.Server
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			srv := new(grpc.Server)

			return test{
				name: "set success",
				srv:  srv,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.grpc.srv, srv) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCServer(tt.srv)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCOption(t *testing.T) {
	type test struct {
		name      string
		opts      []grpc.ServerOption
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			opts := []grpc.ServerOption{}

			return test{
				name: "set success",
				opts: opts,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if !reflect.DeepEqual(got.grpc.opts, opts) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCOption(tt.opts...)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCRegistFunc(t *testing.T) {
	type test struct {
		name      string
		fn        func(*grpc.Server)
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			fn := func(*grpc.Server) {}

			return test{
				name: "set success",
				fn:   fn,
				checkFunc: func(opt Option) error {
					got := new(server)
					opt(got)

					if reflect.ValueOf(got.grpc.reg).Pointer() != reflect.ValueOf(fn).Pointer() {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCRegistFunc(tt.fn)
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
					return errors.New("invalid param was set")
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
					return errors.New("invalid param was set")
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
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCMaxReceiveMessageSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCMaxSendMessageSize(t *testing.T) {
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCMaxSendMessageSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCInitialWindowSize(t *testing.T) {
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCInitialWindowSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCInitialConnWindowSize(t *testing.T) {
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCInitialConnWindowSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveMaxConnIdle(t *testing.T) {
	type test struct {
		name      string
		max       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			max:  "10m",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.maxConnIdle != 10*time.Minute {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when max is empty",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when max is invalid",
			max:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveMaxConnIdle(tt.max)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveMaxConnAge(t *testing.T) {
	type test struct {
		name      string
		max       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			max:  "20m",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.maxConnAge != 20*time.Minute {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when max is empty",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when max is invalid",
			max:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveMaxConnAge(tt.max)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveMaxConnAgeGrace(t *testing.T) {
	type test struct {
		name      string
		max       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			max:  "30m",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.maxConnAgeGrace != 30*time.Minute {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when max is empty",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when max is invalid",
			max:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveMaxConnAgeGrace(tt.max)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveTime(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "40m",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.t != 40*time.Minute {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when dur is empty",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when dur is invalid",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveTime(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCKeepaliveTimeout(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "50m",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive.timeout != 50*time.Minute {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when dur is empty",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when dur is invalid",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.grpc.keepAlive != nil {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCKeepaliveTimeout(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCWriteBufferSize(t *testing.T) {
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCWriteBufferSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCReadBufferSize(t *testing.T) {
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCReadBufferSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCConnectionTimeout(t *testing.T) {
	type test struct {
		name      string
		to        string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			to:   "60m",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when to is empty",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing when to is invalid",
			to:   "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCConnectionTimeout(tt.to)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCMaxHeaderListSize(t *testing.T) {
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCMaxHeaderListSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCHeaderTableSize(t *testing.T) {
	type test struct {
		name      string
		size      int
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			size: 1024,
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 1 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCHeaderTableSize(tt.size)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPCInterceptors(t *testing.T) {
	type test struct {
		name      string
		names     []string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCInterceptors(tt.names...)
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
					return errors.New("invalid param (enableRestart) was set")
				}

				if got.mode != REST {
					return errors.New("invalid param (mode) was set")
				}

				if !reflect.DeepEqual(got.eg, errgroup.Get()) {
					return errors.New("invalid param (eg) was set")
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
					return errors.New("invalid param (name) was set")
				}

				if !reflect.DeepEqual(got.eg, errgroup.Get()) {
					return errors.New("invalid param (eg) was set")
				}

				if got.http.h == nil {
					return errors.New("invalid param (http.h) was set")
				}

				if got.host != "host" {
					return errors.New("invalid param (host) was set")
				}

				if got.it != 3*time.Second {
					return errors.New("invalid param (it) was set")
				}

				if got.port != 8080 {
					return errors.New("invalid param (port) was set")
				}

				if got.pwt != 2*time.Second {
					return errors.New("invalid param (pwt) was set")
				}

				if got.rht != 3*time.Second {
					return errors.New("invalid param (rht) was set")
				}

				if got.rt != 2*time.Second {
					return errors.New("invalid param (rt) was set")
				}

				if got.mode != REST {
					return errors.New("invalid param (mode) was set")
				}

				if got.sddur != 4*time.Second {
					return errors.New("invalid param (sddur) was set")
				}

				if got.wt != 3*time.Second {
					return errors.New("invalid param (wt) was set")
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
