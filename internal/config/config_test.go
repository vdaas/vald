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

// package config providers configuration type and load configuration logic
package config

import (
	"io/fs"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	os.Exit(m.Run())
}

func TestGlobalConfig_Bind(t *testing.T) {
	type fields struct {
		Version string
		TZ      string
		Logging *Logging
	}
	type want struct {
		want *GlobalConfig
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *GlobalConfig) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *GlobalConfig) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return GlobalConfig when all fields are embedded",
				fields: fields{
					Version: "v1.0.0",
					TZ:      "UTC",
					Logging: &Logging{
						Logger: "glg",
						Level:  "warn",
						Format: "json",
					},
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Level:  "warn",
							Format: "json",
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return GlobalConfig when version and time_zone are embedded but logging is nil",
				fields: fields{
					Version: "v1.0.0",
					TZ:      "UTC",
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return GlobalConfig when version is empty and time_zone is embedded but logging is nil",
				fields: fields{
					TZ: "UTC",
				},
				want: want{
					want: &GlobalConfig{
						TZ: "UTC",
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return GlobalConfig when version is embedded and time_zone is empty but logging is nil",
				fields: fields{
					Version: "v1.0.0",
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return GlobalConfig when Logging.Logger is an empty",
				fields: fields{
					Version: "v1.0.0",
					TZ:      "UTC",
					Logging: &Logging{
						Level:  "warn",
						Format: "json",
					},
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Level:  "warn",
							Format: "json",
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return GlobalConfig when Logging.Level is an empty",
				fields: fields{
					Version: "v1.0.0",
					TZ:      "UTC",
					Logging: &Logging{
						Logger: "glg",
						Format: "json",
					},
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Format: "json",
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return GlobalConfig when Logging.Format is an empty",
				fields: fields{
					Version: "v1.0.0",
					TZ:      "UTC",
					Logging: &Logging{
						Logger: "glg",
						Level:  "warn",
					},
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Level:  "warn",
						},
					},
				},
			}
		}(),
		func() test {
			envPrefix := "GLOBALCONFIG_BIND_"
			env := map[string]string{
				envPrefix + "VERSION": "v1.0.0",
				envPrefix + "TZ":      "UTC",
				envPrefix + "LOGGER":  "glg",
				envPrefix + "LEVEL":   "warn",
				envPrefix + "FORMAT":  "json",
			}

			return test{
				name: "return GlobalConfig when all fields are read from environment variable",
				fields: fields{
					Version: "_" + envPrefix + "VERSION_",
					TZ:      "_" + envPrefix + "TZ_",
					Logging: &Logging{
						Logger: "_" + envPrefix + "LOGGER_",
						Level:  "_" + envPrefix + "LEVEL_",
						Format: "_" + envPrefix + "FORMAT_",
					},
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Level:  "warn",
							Format: "json",
						},
					},
				},
				beforeFunc: func() {
					for key, val := range env {
						t.Setenv(key, val)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &GlobalConfig{
				Version: test.fields.Version,
				TZ:      test.fields.TZ,
				Logging: test.fields.Logging,
			}

			got := c.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		path string
		cfg  interface{}
	}
	type want struct {
		want interface{}
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			path := "read_config_test.json"
			data := `{
				"version": "v1.0.0",
				"time_zone": "UTC",
				"logging": {
					"logger": "glg",
					"level": "warn",
					"format": "json"
				}}`
			cfg := new(GlobalConfig)

			return test{
				name: "return nil when read json file and input data type is struct",
				args: args{
					path: path,
					cfg:  cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Level:  "warn",
							Format: "json",
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_config_test.json"
			data := `{
				"version": "v1.0.0",
				"time_zone": "UTC"
				}`
			cfg := make(map[string]string)

			return test{
				name: "return nil when read json file successes and input data type is map",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &map[string]string{
						"version":   "v1.0.0",
						"time_zone": "UTC",
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_config_test.json"
			data := `{
				"version": "v1.0.0",
				"time_zone": "UTC",
				"logging": {
					"logger": "glg"
				}
			}`
			cfg := make(map[string]interface{})

			return test{
				name: "return nil when read json file successes and input data type is nested map",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &map[string]interface{}{
						"version":   "v1.0.0",
						"time_zone": "UTC",
						"logging": map[string]interface{}{
							"logger": "glg",
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_config_test.json"
			data := `[
				{
					"addr": "0.0.0.0",
					"port": "8080"
				},
				{
					"addr": "0.0.0.0",
					"port": "3001"
				}
			]`
			cfg := make([]map[string]interface{}, 0)

			return test{
				name: "return nil when read json file successes and input data type is map slice",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &[]map[string]interface{}{
						{
							"addr": "0.0.0.0",
							"port": "8080",
						},
						{
							"addr": "0.0.0.0",
							"port": "3001",
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_config_test.json"
			data := `"vdaas"`
			var cfg string

			return test{
				name: "return nil when read json file successes and input data type is string",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: func() (str *string) {
						s := "vdaas"
						return &s
					}(),
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_test_config.yaml"
			data := "time_zone: UTC\nversion: v1.0.0\nlogging:\n  format: json\n  level: warn\n  logger: glg"
			cfg := new(GlobalConfig)

			return test{
				name: "return nil when read yaml file and input data type is struct",
				args: args{
					path: path,
					cfg:  cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Level:  "warn",
							Format: "json",
						},
					},
					err: nil,
				},
			}
		}(),

		func() test {
			path := "read_config_test.yaml"
			data := "version: v1.0.0\ntime_zone: UTC"
			cfg := make(map[string]string)

			return test{
				name: "return nil when read yaml file successes and input data type is map",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &map[string]string{
						"version":   "v1.0.0",
						"time_zone": "UTC",
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_config_test.yaml"
			data := "version: v1.0.0\ntime_zone: UTC\nlogging:\n  logger: glg"
			cfg := make(map[string]interface{})

			return test{
				name: "return nil when read yaml file successes and input data type is nested map",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &map[string]interface{}{
						"version":   "v1.0.0",
						"time_zone": "UTC",
						"logging": map[interface{}]interface{}{
							"logger": "glg",
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_config_test.yaml"
			data := "- \n  addr: 0.0.0.0\n  port: \"8080\"\n- \n  addr: 0.0.0.0\n  port: \"3001\""
			cfg := make([]map[string]interface{}, 0)

			return test{
				name: "return nil when read yaml file successes and input data type is map slice",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &[]map[string]interface{}{
						{
							"addr": "0.0.0.0",
							"port": "8080",
						},
						{
							"addr": "0.0.0.0",
							"port": "3001",
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_config_test.yaml"
			data := `"vdaas"`
			var cfg string

			return test{
				name: "return nil when read yaml file successes and input data type is string",
				args: args{
					path: path,
					cfg:  &cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: func() (str *string) {
						s := "vdaas"
						return &s
					}(),
					err: nil,
				},
			}
		}(),
		func() test {
			path := "read_test_config.yaml"
			cfg := new(GlobalConfig)

			return test{
				name: "return no entry error when the file open fails",
				args: args{
					path: path,
					cfg:  cfg,
				},
				want: want{
					want: cfg,
					err: &fs.PathError{
						Op:   "open",
						Path: "read_test_config.yaml",
						Err:  syscall.ENOENT,
					},
				},
			}
		}(),
		func() test {
			path := "read_test_config.yaml"
			data := "timezone\n:"
			cfg := new(GlobalConfig)

			return test{
				name: "return yaml decode error when the contents of yaml is invalid",
				args: args{
					path: path,
					cfg:  cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: cfg,
					err:  errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `timezone` into config.GlobalConfig"),
				},
			}
		}(),
		func() test {
			path := "read_test_config.json"
			data := "timezone\n:"
			cfg := new(GlobalConfig)

			return test{
				name: "return json decode error when the contents of json file is invalid",
				args: args{
					path: path,
					cfg:  cfg,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString(data); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: cfg,
					err:  errors.New("invalid character 't' looking for beginning of value"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			err := Read(test.args.path, test.args.cfg)
			if err := checkFunc(test.want, test.args.cfg, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGetActualValue(t *testing.T) {
	type args struct {
		val string
	}
	type want struct {
		wantRes string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes string) error {
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return v1.0.0 when val is set in environment variable",
				args: args{
					val: "_GETACTUALVALUE_VERSION_",
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					t.Setenv("GETACTUALVALUE_VERSION", "v1.0.0")
				},
				want: want{
					wantRes: "v1.0.0",
				},
			}
		}(),
		func() test {
			return test{
				name: "return v1.0.0 when val is $VERSION",
				args: args{
					val: "$GETACTUALVALUE_1_VERSION",
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					t.Setenv("GETACTUALVALUE_1_VERSION", "v1.0.0")
				},
				want: want{
					wantRes: "v1.0.0",
				},
			}
		}(),
		func() test {
			return test{
				name: "return VERSION version when val is VERSION",
				args: args{
					val: "VERSION",
				},
				want: want{
					wantRes: "VERSION",
				},
			}
		}(),
		func() test {
			fname := "version"

			return test{
				name: "return file contents when val is file://env",
				args: args{
					val: "file://" + fname,
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					f, err := os.Create(fname)
					if err != nil {
						t.Error(err)
						return
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()

					if _, err := f.WriteString("v1.0.0"); err != nil {
						t.Error(err)
					}
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := os.Remove(fname); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: "v1.0.0",
				},
			}
		}(),
		func() test {
			fname := "version"
			return test{
				name: "return file contents when val is file://env",
				args: args{
					val: "file://" + fname,
				},
				want: want{
					wantRes: "file://" + fname,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			gotRes := GetActualValue(test.args.val)
			if err := checkFunc(test.want, gotRes); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGetActualValues(t *testing.T) {
	type args struct {
		vals []string
	}
	type want struct {
		want []string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got []string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			envPrefix := "GETACTUALVALUES_"
			env := map[string]string{
				envPrefix + "VERSION": "v1.0.0",
				envPrefix + "LOGGER":  "glg",
			}

			return test{
				name: "return v1.0.0 and glg when vals are _LOGGER_ and _VERSION_",
				args: args{
					vals: []string{
						"_" + envPrefix + "VERSION_",
						"_" + envPrefix + "LOGGER_",
					},
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					for key, val := range env {
						t.Setenv(key, val)
					}
				},
				want: want{
					want: []string{
						"v1.0.0",
						"glg",
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return v1.0.0 and LOGGER when vals are _VERSION_ and LOGGER",
				args: args{
					vals: []string{
						"_GETACTUALVALUES_1_VERSION_",
						"LOGGER",
					},
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					t.Setenv("GETACTUALVALUES_1_VERSION", "v1.0.0")
				},
				want: want{
					want: []string{
						"v1.0.0",
						"LOGGER",
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return empty when vals is empty",
				args: args{
					vals: []string{},
				},
				want: want{
					want: []string{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := GetActualValues(test.args.vals)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_checkPrefixAndSuffix(t *testing.T) {
	type args struct {
		str  string
		pref string
		suf  string
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return true when prefix and suffix are _ and str is _POD_NAME_",
			args: args{
				str:  "_POD_NAME_",
				pref: "_",
				suf:  "_",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when prefix and suffix are _ and str is __POD_NAME__",
			args: args{
				str:  "__POD_NAME__",
				pref: "_",
				suf:  "_",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when prefix and suffix are __ and str is __POD_NAME__",
			args: args{
				str:  "__POD_NAME__",
				pref: "__",
				suf:  "__",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when prefix is $ and suffix is # and str is $POD_NAME#",
			args: args{
				str:  "$POD_NAME#",
				pref: "$",
				suf:  "#",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when prefix is $# and suffix is #$ and str is $#POD_NAME#$",
			args: args{
				str:  "$#POD_NAME#$",
				pref: "$#",
				suf:  "#$",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when prefix is _ and suffix is empty and str is _POD_NAME_",
			args: args{
				str:  "_POD_NAME_",
				pref: "_",
				suf:  "",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when prefix is empty and suffix is _ and str is _POD_NAME_",
			args: args{
				str:  "_POD_NAME_",
				pref: "",
				suf:  "_",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false when prefix and suffix are _ and str is empty",
			args: args{
				str:  "",
				pref: "_",
				suf:  "_",
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when prefix and suffix are _ and str is _POD_NAME",
			args: args{
				str:  "_POD_NAME",
				pref: "_",
				suf:  "_",
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when prefix and suffix are _ and str is POD_NAME_",
			args: args{
				str:  "POD_NAME_",
				pref: "_",
				suf:  "_",
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when prefix and suffix are _ and str is POD_NAME&",
			args: args{
				str:  "POD_NAME&",
				pref: "_",
				suf:  "_",
			},
			want: want{
				want: false,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := checkPrefixAndSuffix(test.args.str, test.args.pref, test.args.suf)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestToRawYaml(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return row string when data is an int type",
			args: args{
				data: 1,
			},
			want: want{
				want: "1\n",
			},
		},
		{
			name: "return row string when data is a string type",
			args: args{
				data: "vdaas.vald",
			},
			want: want{
				want: "vdaas.vald\n",
			},
		},
		{
			name: "return row string when data is a map string type",
			args: args{
				data: map[string]string{
					"time_zone": "UTC",
				},
			},
			want: want{
				want: "time_zone: UTC\n",
			},
		},
		{
			name: "return row string when data is a nested map type",
			args: args{
				data: map[string]interface{}{
					"logging": map[string]interface{}{
						"logger": "glg",
					},
				},
			},
			want: want{
				want: "logging:\n  logger: glg\n",
			},
		},
		{
			name: "return row string when data is a empty string",
			args: args{
				data: "",
			},
			want: want{
				want: "\"\"\n",
			},
		},
		{
			name: "return row string when data is a GlobalConfig type",
			args: args{
				data: GlobalConfig{
					Version: "v1.0.0",
					TZ:      "UTC",
					Logging: &Logging{
						Logger: "glg",
						Level:  "warn",
						Format: "json",
					},
				},
			},
			want: want{
				want: "version: v1.0.0\ntime_zone: UTC\nlogging:\n  logger: glg\n  level: warn\n  format: json\n",
			},
		},
		{
			name: "return row string when data is a nil",
			args: args{
				data: nil,
			},
			want: want{
				want: "null\n",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := ToRawYaml(test.args.data)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
