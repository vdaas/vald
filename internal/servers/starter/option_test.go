//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package starter

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/servers/server"
)

func TestWithConfig(t *testing.T) {
	type test struct {
		name      string
		cfg       *config.Servers
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			cfg := new(config.Servers)

			return test{
				name: "set success",
				cfg:  cfg,
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if !reflect.DeepEqual(got.cfg, cfg) {
						return errors.New("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithConfig(tt.cfg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPC(t *testing.T) {
	type test struct {
		name      string
		opts      func(cfg *config.Server) []server.Option
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			srvOpts := func(cfg *config.Server) []server.Option {
				return nil
			}

			return test{
				name: "set success",
				opts: srvOpts,
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if reflect.ValueOf(got.grpc).Pointer() != reflect.ValueOf(srvOpts).Pointer() {
						return errors.New("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPC(tt.opts)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithREST(t *testing.T) {
	type test struct {
		name      string
		opts      func(cfg *config.Server) []server.Option
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			srvOpts := func(cfg *config.Server) []server.Option {
				return nil
			}

			return test{
				name: "set success",
				opts: srvOpts,
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if reflect.ValueOf(got.rest).Pointer() != reflect.ValueOf(srvOpts).Pointer() {
						return errors.New("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithREST(tt.opts)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGQL(t *testing.T) {
	type test struct {
		name      string
		opts      func(cfg *config.Server) []server.Option
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			srvOpts := func(cfg *config.Server) []server.Option {
				return nil
			}

			return test{
				name: "set success",
				opts: srvOpts,
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if reflect.ValueOf(got.gql).Pointer() != reflect.ValueOf(srvOpts).Pointer() {
						return errors.New("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGQL(tt.opts)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStartFunc(t *testing.T) {
	type args struct {
		name string
		fn   func() error
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			name := "name"
			fn := func() error { return nil }
			return test{
				name: "set success",
				args: args{
					name: name,
					fn:   fn,
				},
				checkFunc: func(opt Option) error {
					got := &srvs{
						pstartf: make(map[string]func() error, 1),
					}
					opt(got)

					if gfn, ok := got.pstartf[name]; ok {
						if reflect.ValueOf(gfn).Pointer() != reflect.ValueOf(fn).Pointer() {
							return errors.New("invalid param was set")
						}
					} else {
						return errors.New("param was not set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStartFunc(tt.args.name, tt.args.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStopFunc(t *testing.T) {
	type args struct {
		name string
		fn   func() error
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			name := "name"
			fn := func() error { return nil }
			return test{
				name: "set success",
				args: args{
					name: name,
					fn:   fn,
				},
				checkFunc: func(opt Option) error {
					got := &srvs{
						pstopf: make(map[string]func() error, 1),
					}
					opt(got)

					if gfn, ok := got.pstopf[name]; ok {
						if reflect.ValueOf(gfn).Pointer() != reflect.ValueOf(fn).Pointer() {
							return errors.New("invalid param was set")
						}
					} else {
						return errors.New("param was not set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStopFunc(tt.args.name, tt.args.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
