// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/tls"
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
			name: "not set nothing when host is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					host: "host",
				}
				opt(got)

				if got.host != "host" {
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
		port      uint16
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
			name: "not set when port is 0",
			checkFunc: func(opt Option) error {
				got := &server{
					port: 8080,
				}
				opt(got)

				if got.port != 8080 {
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
		srvName   string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name:    "set success",
			srvName: "vald",
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if got.name != "vald" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when srvName is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					name: "name",
				}
				opt(got)

				if got.name != "name" {
					return errors.New("invalid param was set")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithName(tt.srvName)
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

		{
			name: "not set when eg is nil",
			checkFunc: func(opt Option) error {
				eg := errgroup.Get()

				got := &server{
					eg: eg,
				}
				opt(got)

				if !reflect.DeepEqual(got.eg, eg) {
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

		{
			name: "not set when fn is nil",
			checkFunc: func(opt Option) error {
				fn := func() error { return nil }

				got := &server{
					preStartFunc: fn,
				}
				opt(got)

				if reflect.ValueOf(got.preStartFunc).Pointer() != reflect.ValueOf(fn).Pointer() {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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

		{
			name: "not set when fn is nil",
			checkFunc: func(opt Option) error {
				fn := func() error { return nil }

				got := &server{
					preStopFunc: fn,
				}
				opt(got)

				if reflect.ValueOf(got.preStopFunc).Pointer() != reflect.ValueOf(fn).Pointer() {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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
			name: "not set when dur is invalid",
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

		{
			name: "not set when dur is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					pwt: 5 * time.Second,
				}
				opt(got)

				if got.pwt != 0 {
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
			name: "set default when dur is invalid",
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

		{
			name: "set default when dur is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					sddur: 5 * time.Second,
				}
				opt(got)

				if got.sddur != 0 {
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
			name: "set default when dur is invalid",
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

		{
			name: "set default when dur is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					rht: 5 * time.Second,
				}
				opt(got)

				if got.rht != 0 {
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
			name: "set default when dur is invalid",
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

		{
			name: "set default when dur is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					rt: 5 * time.Second,
				}
				opt(got)

				if got.rt != 0 {
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
			name: "set default when dur is invalid",
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

		{
			name: "set default when dur is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					wt: 5 * time.Second,
				}
				opt(got)

				if got.wt != 0 {
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
			name: "set default when dur is invalid",
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

		{
			name: "set default when dur is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					it: 5 * time.Second,
				}
				opt(got)

				if got.it != 0 {
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

		{
			name: "not set when lc is nil",
			checkFunc: func(opt Option) error {
				lc := new(net.ListenConfig)
				got := &server{
					lc: lc,
				}
				opt(got)

				if !reflect.DeepEqual(got.lc, lc) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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
		m         ServerMode
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

		{
			name: "not set when mode is invalid",
			m:    ServerMode(100),
			checkFunc: func(opt Option) error {
				got := &server{
					mode: GRPC,
				}
				opt(got)

				if got.mode != GRPC {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when mode is empty",
			checkFunc: func(opt Option) error {
				got := &server{
					mode: GRPC,
				}
				opt(got)

				if got.mode != GRPC {
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

		{
			name: "not set when cfg is nil",
			checkFunc: func(opt Option) error {
				cfg := new(tls.Config)
				got := &server{
					tcfg: cfg,
				}
				opt(got)

				if !reflect.DeepEqual(got.tcfg, cfg) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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

		{
			name: "not set when hdr is nil",
			checkFunc: func(opt Option) error {
				hdr := new(handler)
				got := new(server)
				got.http.h = hdr
				opt(got)

				if reflect.ValueOf(got.http.h).Pointer() != reflect.ValueOf(hdr).Pointer() {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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

		{
			name: "not set when srv is nil",
			checkFunc: func(opt Option) error {
				srv := new(http.Server)
				got := new(server)
				got.http.srv = srv
				opt(got)

				if !reflect.DeepEqual(got.http.srv, srv) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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

		{
			name: "not set when srv is nil",
			checkFunc: func(opt Option) error {
				srv := new(grpc.Server)
				got := new(server)
				got.grpc.srv = srv
				opt(got)

				if !reflect.DeepEqual(got.grpc.srv, srv) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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

		{
			name: "not set when opts is nil",
			checkFunc: func(opt Option) error {
				opts := []grpc.ServerOption{}
				got := new(server)
				got.grpc.opts = opts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, opts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
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

					for _, reg := range got.grpc.regs {
						if reflect.ValueOf(reg).Pointer() == reflect.ValueOf(fn).Pointer() {
							return nil
						}
					}
					return errors.New("invalid param was set")
				},
			}
		}(),

		{
			name: "not set when f is nil",
			checkFunc: func(opt Option) error {
				fn := func(*grpc.Server) {}
				got := new(server)
				got.grpc.regs = append(got.grpc.regs, fn)
				opt(got)
				for _, reg := range got.grpc.regs {
					if reflect.ValueOf(reg).Pointer() == reflect.ValueOf(fn).Pointer() {
						return nil
					}
				}
				return errors.New("invalid param was set")
			},
		},
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
			name: "set success when size is mode than 0",
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

		{
			name: "set success when size is -1",
			size: -1,
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
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when size is less than -1",
			size: -2,
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
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
			name: "set success when size is more than 0",
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

		{
			name: "set success when size is -1",
			size: -1,
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
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when size is less than -1",
			size: -2,
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
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
			name: "set success when size is more than 0",
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

		{
			name: "set success when size is -1",
			size: -1,
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
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when size is less than -1",
			size: -2,
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
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
			name: "set success when size is more than 0",
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

		{
			name: "set success when size is -1",
			size: -1,
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
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when size is less than -1",
			size: -2,
			checkFunc: func(opt Option) error {
				sopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(100 * time.Second),
				}
				got := new(server)
				got.grpc.opts = sopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, sopts) {
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
			name: "not set when max is empty",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					maxConnIdle: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when max is invalid",
			max:  "vald",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					maxConnIdle: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
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
			name: "not set when max is empty",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					maxConnAge: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when max is invalid",
			max:  "vald",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					maxConnAge: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
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
			name: "not set when max is empty",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					maxConnAgeGrace: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when max is invalid",
			max:  "vald",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					maxConnAgeGrace: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
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
			name: "not set when dur is empty",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					t: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when dur is invalid",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					t: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
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
			name: "returns nothing when dur is empty",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					timeout: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "returns nothing when dur is invalid",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				grpcKeepalive := &grpcKeepalive{
					timeout: 10 * time.Second,
				}
				got := new(server)
				got.grpc.keepAlive = grpcKeepalive
				opt(got)

				if !reflect.DeepEqual(got.grpc.keepAlive, grpcKeepalive) {
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
			name: "set success when size is more than 0",
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

		{
			name: "set success when size is -1",
			size: -1,
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
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when size is less than -1",
			size: -2,
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
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
			name: "set success when size is more than 0",
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

		{
			name: "set success when size is -1",
			size: -1,
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
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when size is less than -1",
			size: -2,
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
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
			name: "not set when to is empty",
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when to is invalid",
			to:   "vald",
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
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
		size      uint32
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

		{
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
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
		size      uint32
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

		{
			name: "not set when size is 0",
			checkFunc: func(opt Option) error {
				gopts := []grpc.ServerOption{
					grpc.ConnectionTimeout(10 * time.Second),
				}
				got := new(server)
				got.grpc.opts = gopts
				opt(got)

				if !reflect.DeepEqual(got.grpc.opts, gopts) {
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
		{
			name:  "Add RecoverInterceptor using 'RecoverInterceptor'",
			names: []string{"RecoverInterceptor"},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 2 {
					return errors.Errorf("Expecting two elements in got.grpc.opts: got = %#v", got)
				}

				return nil
			},
		},
		{
			name:  "Add RecoverInterceptor using 'Recover'",
			names: []string{"Recover"},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 2 {
					return errors.Errorf("Expecting two elements in got.grpc.opts: got = %#v", got)
				}

				return nil
			},
		},
		{
			name:  "Add AccessLogInterceptor using 'AccessLogInterceptor'",
			names: []string{"AccessLogInterceptor"},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 2 {
					return errors.Errorf("Expecting two elements in got.grpc.opts: got = %#v", got)
				}

				return nil
			},
		},
		{
			name:  "Add AccessLogInterceptor using 'AccessLog'",
			names: []string{"AccessLog"},
			checkFunc: func(opt Option) error {
				got := new(server)
				opt(got)

				if len(got.grpc.opts) != 2 {
					return errors.Errorf("Expecting two elements in got.grpc.opts: got = %#v", got)
				}

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
			if err := tt.checkFunc(defaultOptions); err != nil {
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
		port uint16
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
				path: "/path",
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

// NOT IMPLEMENTED BELOW
//
// func TestWithNetwork(t *testing.T) {
// 	type args struct {
// 		network string
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           network:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           network:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithNetwork(test.args.network)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithSocketPath(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithSocketPath(test.args.path)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithSocketFlag(t *testing.T) {
// 	type args struct {
// 		flg control.SocketFlag
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           flg:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           flg:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithSocketFlag(test.args.flg)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithPreStopFunction(t *testing.T) {
// 	type args struct {
// 		f func() error
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           f:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           f:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithPreStopFunction(test.args.f)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithHTTP2Enabled(t *testing.T) {
// 	type args struct {
// 		enabled bool
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           enabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           enabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithHTTP2Enabled(test.args.enabled)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithHandlerLimit(t *testing.T) {
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithHandlerLimit(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithPermitProhibitedCipherSuites(t *testing.T) {
// 	type args struct {
// 		perm bool
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           perm:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           perm:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithPermitProhibitedCipherSuites(test.args.perm)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithMaxUploadBufferPerConnection(t *testing.T) {
// 	type args struct {
// 		size int32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithMaxUploadBufferPerConnection(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithMaxUploadBufferPerStream(t *testing.T) {
// 	type args struct {
// 		size int32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithMaxUploadBufferPerStream(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithMaxConcurrentStreams(t *testing.T) {
// 	type args struct {
// 		size uint32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithMaxConcurrentStreams(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithMaxDecoderHeaderTableSize(t *testing.T) {
// 	type args struct {
// 		size uint32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithMaxDecoderHeaderTableSize(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithMaxEncoderHeaderTableSize(t *testing.T) {
// 	type args struct {
// 		size uint32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithMaxEncoderHeaderTableSize(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithMaxReadFrameSize(t *testing.T) {
// 	type args struct {
// 		size uint32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithMaxReadFrameSize(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithGRPCKeepaliveMinTime(t *testing.T) {
// 	type args struct {
// 		min string
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           min:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           min:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithGRPCKeepaliveMinTime(test.args.min)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithGRPCKeepalivePermitWithoutStream(t *testing.T) {
// 	type args struct {
// 		pws bool
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           pws:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           pws:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithGRPCKeepalivePermitWithoutStream(test.args.pws)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithGRPCMaxConcurrentStreams(t *testing.T) {
// 	type args struct {
// 		size uint32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithGRPCMaxConcurrentStreams(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithGRPCNumStreamWorkers(t *testing.T) {
// 	type args struct {
// 		size uint32
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           size:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithGRPCNumStreamWorkers(test.args.size)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithGRPCSharedWriteBuffer(t *testing.T) {
// 	type args struct {
// 		enable bool
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           enable:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           enable:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithGRPCSharedWriteBuffer(test.args.enable)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithGRPCWaitForHandlers(t *testing.T) {
// 	type args struct {
// 		wait bool
// 	}
// 	type want struct {
// 		want Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Option) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           wait:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           wait:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := WithGRPCWaitForHandlers(test.args.wait)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
