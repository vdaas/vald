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
package servers

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/servers/server"
)

func TestWithServer(t *testing.T) {
	type test struct {
		name      string
		srv       server.Server
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			srv := &mockServer{
				NameFunc: func() string {
					return "srv"
				},
			}

			return test{
				name: "set success",
				srv:  srv,
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if len(got.servers) != 1 {
						return errors.Errorf("servers count is wrong. want: %v, got: %v", 1, len(got.servers))
					}

					gsrv, ok := got.servers["srv"]
					if !ok {
						return errors.New("servers['srv'] is nothing")
					}

					if !reflect.DeepEqual(gsrv, srv) {
						return errors.Errorf("servers['srv'] is not equals. want: %v, got: %b", srv, gsrv)
					}

					return nil
				},
			}
		}(),
		{
			name: "do nothing",
			checkFunc: func(opt Option) error {
				got := new(listener)
				opt(got)

				if got.servers != nil {
					return errors.Errorf("server is not nil: %v", got.servers)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithServer(tt.srv)
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
		func() test {
			eg, _ := errgroup.New(context.Background())

			return test{
				name: "set success",
				eg:   eg,
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if !reflect.DeepEqual(got.eg, eg) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
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

func TestWithShutdownDuration(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt Option) error {
				got := new(listener)
				opt(got)

				if !reflect.DeepEqual(got.sddur, 10*time.Second) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			dur:  "vald",
			checkFunc: func(opt Option) error {
				got := new(listener)
				opt(got)

				if !reflect.DeepEqual(got.sddur, 20*time.Second) {
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

func TestWithStartUpStrategy(t *testing.T) {
	type test struct {
		name      string
		strg      []string
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			strg := []string{
				"strg_1",
				"strg_2",
			}

			return test{
				name: "set success",
				strg: strg,
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if !reflect.DeepEqual(got.sus, strg) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithStartUpStrategy(tt.strg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithShutdownStrategy(t *testing.T) {
	type test struct {
		name      string
		strg      []string
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			strg := []string{
				"strg_1",
				"strg_2",
			}

			return test{
				name: "set success",
				strg: strg,
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if !reflect.DeepEqual(got.sds, strg) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithShutdownStrategy(tt.strg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
