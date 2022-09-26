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

// Package setting stores all server application settings
package config

import (
	"io/fs"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	type want struct {
		wantCfg *Data
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *Data, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotCfg *Data, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(gotCfg, w.wantCfg,
			comparator.IgnoreTypes(config.Observability{})); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	tests := []test{
		func() test {
			data := `{
				"version": "v1.0.0",
				"server_config": {
					"full_shutdown_duration": "10ms"
				},
				"observability": {
					"enabled": true
				},
				"ngt": {
					"index_path": "/var/index"
				}
			}`
			return test{
				name: "return Data and nil when the json bind successes",
				args: args{
					path: "bind_success.json",
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					f, err := os.Create(a.path)
					if err != nil {
						t.Fatal(err)
					}
					if _, err := f.Write([]byte(data)); err != nil {
						t.Fatal(err)
					}
					if err := f.Close(); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := os.Remove(a.path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantCfg: &Data{
						GlobalConfig: config.GlobalConfig{
							Version: "v1.0.0",
						},
						Server: &config.Servers{
							FullShutdownDuration: "10ms",
							ShutdownStrategy:     make([]string, 0),
							StartUpStrategy:      make([]string, 0),
							TLS: &config.TLS{
								Enabled: false,
							},
						},
						Observability: &config.Observability{
							Enabled:    true,
							Metrics:    new(config.Metrics),
							Trace:      new(config.Trace),
							Prometheus: new(config.Prometheus),
							Jaeger:     new(config.Jaeger),
						},
						NGT: &config.NGT{
							IndexPath: "/var/index",
							VQueue:    new(config.VQueue),
							KVSDB:     new(config.KVSDB),
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			data := `{
				"version": "v1.0.0",
				"server_config": {
					"full_shutdown_duration": "10ms"
				},
				"ngt": {
					"index_path": "/var/index"
				}
			}`
			return test{
				name: "return Data and nil when the json bind successes but the input json value of observability is empty",
				args: args{
					path: "bind_success_but_observability_is_empty.json",
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					f, err := os.Create(a.path)
					if err != nil {
						t.Fatal(err)
					}
					if _, err := f.Write([]byte(data)); err != nil {
						t.Fatal(err)
					}
					if err := f.Close(); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := os.Remove(a.path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantCfg: &Data{
						GlobalConfig: config.GlobalConfig{
							Version: "v1.0.0",
						},
						Server: &config.Servers{
							FullShutdownDuration: "10ms",
							ShutdownStrategy:     make([]string, 0),
							StartUpStrategy:      make([]string, 0),
							TLS: &config.TLS{
								Enabled: false,
							},
						},
						Observability: new(config.Observability),
						NGT: &config.NGT{
							IndexPath: "/var/index",
							VQueue:    new(config.VQueue),
							KVSDB:     new(config.KVSDB),
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			data := `
                version: v1.0.0
                server_config:
                  full_shutdown_duration: 10ms
                observability:
                  enabled: true
                ngt:
                  index_path: /var/index
                `
			return test{
				name: "return Data and nil when the yaml bind successes",
				args: args{
					path: "bind_success.yaml",
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					f, err := os.Create(a.path)
					if err != nil {
						t.Fatal(err)
					}
					if _, err := f.Write([]byte(data)); err != nil {
						t.Fatal(err)
					}
					if err := f.Close(); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := os.Remove(a.path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantCfg: &Data{
						GlobalConfig: config.GlobalConfig{
							Version: "v1.0.0",
						},
						Server: &config.Servers{
							FullShutdownDuration: "10ms",
							ShutdownStrategy:     make([]string, 0),
							StartUpStrategy:      make([]string, 0),
							TLS: &config.TLS{
								Enabled: false,
							},
						},
						Observability: &config.Observability{
							Enabled:    true,
							Metrics:    new(config.Metrics),
							Trace:      new(config.Trace),
							Prometheus: new(config.Prometheus),
							Jaeger:     new(config.Jaeger),
						},
						NGT: &config.NGT{
							IndexPath: "/var/index",
							VQueue:    new(config.VQueue),
							KVSDB:     new(config.KVSDB),
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			data := `
                version: v1.0.0
                server_config:
                  full_shutdown_duration: 10ms
                ngt:
                  index_path: /var/index
                `
			return test{
				name: "return Data and nil when the yaml bind successes but the input yaml value of observability is empty",
				args: args{
					path: "bind_success_but_observability_is_empty.yaml",
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					f, err := os.Create(a.path)
					if err != nil {
						t.Fatal(err)
					}
					if _, err := f.Write([]byte(data)); err != nil {
						t.Fatal(err)
					}
					if err := f.Close(); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := os.Remove(a.path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantCfg: &Data{
						GlobalConfig: config.GlobalConfig{
							Version: "v1.0.0",
						},
						Server: &config.Servers{
							FullShutdownDuration: "10ms",
							ShutdownStrategy:     make([]string, 0),
							StartUpStrategy:      make([]string, 0),
							TLS: &config.TLS{
								Enabled: false,
							},
						},
						Observability: new(config.Observability),
						NGT: &config.NGT{
							IndexPath: "/var/index",
							VQueue:    new(config.VQueue),
							KVSDB:     new(config.KVSDB),
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "not_found.txt"
			return test{
				name: "return error when the file does not exist",
				args: args{
					path: path,
				},
				want: want{
					wantCfg: nil,
					err: &fs.PathError{
						Op:   "open",
						Path: path,
						Err:  syscall.Errno(0x2),
					},
				},
			}
		}(),
		func() test {
			path := "empty.json"
			return test{
				name: "return error when the json file is empty",
				args: args{
					path: path,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					f, err := os.Create(a.path)
					if err != nil {
						t.Fatal(err)
					}
					if err := f.Close(); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := os.Remove(a.path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantCfg: nil,
					err:     io.EOF,
				},
			}
		}(),
		func() test {
			path := "empty.yaml"
			return test{
				name: "return error when the yaml file is empty",
				args: args{
					path: path,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					f, err := os.Create(a.path)
					if err != nil {
						t.Fatal(err)
					}
					if err := f.Close(); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := os.Remove(a.path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantCfg: nil,
					err:     io.EOF,
				},
			}
		}(),
		func() test {
			name := "unreadable.txt"
			dir, err := file.MkdirTemp("")
			if err != nil {
				dir = file.Join("/tmp", strconv.FormatInt(time.Now().Unix(), 10))
			}
			path := file.Join(dir, name)
			return test{
				name: "return error when can't read file",
				args: args{
					path: path,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					f, err := file.Open(a.path, os.O_CREATE, fs.ModeIrregular)
					if err != nil {
						if errors.Is(err, fs.ErrPermission) {
							return
						}
						t.Error(err)
					}
					if err := f.Close(); err != nil {
						t.Error(err)
					}
				},
				checkFunc: func(w want, gotCfg *Data, err error) error {
					if errors.Is(err, fs.ErrPermission) {
						return nil
					}
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					return nil
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := os.Remove(a.path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantCfg: nil,
					err:     errors.ErrUnsupportedConfigFileType(".txt"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotCfg, err := NewConfig(test.args.path)
			if err := checkFunc(test.want, gotCfg, err); err != nil {
				tt.Errorf("error = %v, got = %#v", err, gotCfg)
			}
		})
	}
}
