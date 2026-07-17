// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package config

import (
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

func writeTestConfigFile(t *testing.T, path, data string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(data); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func removeTestConfigFile(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	type want struct {
		wantCfg *Data
		err     error
	}
	type test struct {
		want       want
		checkFunc  func(want, *Data, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
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
			data := `
version: v1.0.0
server_config:
  full_shutdown_duration: 10ms
observability:
  enabled: true
exporter:
  index_path: /var/export/index
  concurrency: 10
`
			return test{
				name: "return Data and nil when the yaml bind succeeds",
				args: args{
					path: "bind_success.yaml",
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					writeTestConfigFile(t, a.path, data)
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					removeTestConfigFile(t, a.path)
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
							Enabled: true,
							OTLP: &config.OTLP{
								Attribute: new(config.OTLPAttribute),
							},
							Metrics: new(config.Metrics),
							Trace:   new(config.Trace),
						},
						Exporter: &config.IndexExporter{
							IndexPath:   "/var/export/index",
							Concurrency: 10,
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			data := `
version: v1.0.0
observability:
  enabled: true
exporter:
  index_path: /var/export/index
`
			return test{
				name: "return error when server_config is not specified",
				args: args{
					path: "server_nil.yaml",
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					writeTestConfigFile(t, a.path, data)
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					removeTestConfigFile(t, a.path)
				},
				want: want{
					wantCfg: nil,
					err:     errors.ErrInvalidConfig,
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
`
			return test{
				name: "return error when exporter is not specified",
				args: args{
					path: "exporter_nil.yaml",
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					writeTestConfigFile(t, a.path, data)
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					removeTestConfigFile(t, a.path)
				},
				want: want{
					wantCfg: nil,
					err:     errors.ErrInvalidConfig,
				},
			}
		}(),
		func() test {
			name := "not_found.yaml"
			dir, err := file.MkdirTemp("")
			if err != nil {
				t.Fatal(err)
			}
			path := file.Join(dir, name)
			return test{
				name: "return error when the file does not exist",
				args: args{
					path: path,
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.RemoveAll(dir); err != nil {
						t.Errorf("failed to remove temp dir %s: %v", dir, err)
					}
				},
				checkFunc: func(w want, gotCfg *Data, err error) error {
					if err == nil {
						return errors.New("expected error but got nil")
					}
					if gotCfg != nil {
						return errors.Errorf("got cfg: \"%#v\",\n\t\t\t\twant: nil", gotCfg)
					}
					return nil
				},
				want: want{
					wantCfg: nil,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			// intentionally not t.Parallel(): concurrent subtests race the
			// first-time encoding/json type-encoder cache build for *Data
			// against a sibling's goleak.VerifyNone, which can panic
			// goleak's stack parser on the resulting elided-frame stack.
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

			gotCfg, err := New(test.args.path)
			if cerr := checkFunc(test.want, gotCfg, err); cerr != nil {
				tt.Errorf("error = %v, got = %#v", cerr, gotCfg)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
