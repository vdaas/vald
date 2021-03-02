//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.uber.org/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
}

func TestMain(m *testing.M) {
	log.Init()
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
			env := map[string]string{
				"VERSION": "v1.0.0",
				"TZ":      "UTC",
				"LOGGER":  "glg",
				"LEVEL":   "warn",
				"FORMAT":  "json",
			}
			return test{
				name: "return GlobalConfig when all fields are read from environment variable",
				fields: fields{
					Version: "_VERSION_",
					TZ:      "_TZ_",
					Logging: &Logging{
						Logger: "_LOGGER_",
						Level:  "_LEVEL_",
						Format: "_FORMAT_",
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
						os.Setenv(key, val)
					}
				},
				afterFunc: func() {
					for key := range env {
						os.Unsetenv(key)
					}
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &GlobalConfig{
				Version: test.fields.Version,
				TZ:      test.fields.TZ,
				Logging: test.fields.Logging,
			}

			got := c.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGlobalConfig_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	type fields struct {
		Version string
		TZ      string
		Logging *Logging
	}
	type want struct {
		want *GlobalConfig
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *GlobalConfig, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *GlobalConfig, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			data := []byte(`{
				"version": "v1.0.0",
				"time_zone": "UTC",
				"logging": {
					"logger": "glg",
					"level": "warn",
					"format": "json"
				}}`)
			return test{
				name: "return nil when json unmarshal successes",
				args: args{
					data: data,
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
			data := []byte(`{
				"time_zone": "UTC",
				"logging": {
					"logger": "glg",
					"level": "warn",
					"format": "json"
				}
			}`)
			return test{
				name: "return nil when version key is empty and json unmarshal successes",
				args: args{
					data: data,
				},
				want: want{
					want: &GlobalConfig{
						TZ: "UTC",
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
			data := []byte(`{
				"version": "v1.0.0",
				"logging": {
					"logger": "glg",
					"level": "warn",
					"format": "json"
				}
			}`)
			return test{
				name: "return nil when time_zone key is empty and json unmarshal successes",
				args: args{
					data: data,
				},
				fields: fields{},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
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
			data := []byte(`{
				"version": "v1.0.0",
				"time_zone": "UTC"
			}`)
			return test{
				name: "return nil when logging key is empty and json unmarshal successes",
				args: args{
					data: data,
				},
				fields: fields{},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
					},
					err: nil,
				},
			}
		}(),
		func() test {
			data := []byte(`{
				"version": "v1.0.0",
				"time_zone": "UTC",
				"logging": {
					"level": "warn",
					"format": "json"
				}
			}`)
			return test{
				name: "return nil when logging.logger key is empty and json unmarshal successes",
				args: args{
					data: data,
				},
				fields: fields{},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Level:  "warn",
							Format: "json",
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			data := []byte(`{
				"version": "v1.0.0",
				"time_zone": "UTC",
				"logging": {
					"logger": "glg",
					"format": "json"
				}
			}`)
			return test{
				name: "return nil when logging.level key is empty and json unmarshal successes",
				args: args{
					data: data,
				},
				fields: fields{},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Format: "json",
						},
					},
					err: nil,
				},
			}
		}(),
		func() test {
			data := []byte(`{
				"version": "v1.0.0",
				"time_zone": "UTC",
				"logging": {
					"logger": "glg",
					"level": "warn"
				}
			}`)
			return test{
				name: "return nil when logging.format key is empty and json unmarshal successes",
				args: args{
					data: data,
				},
				fields: fields{},
				want: want{
					want: &GlobalConfig{
						Version: "v1.0.0",
						TZ:      "UTC",
						Logging: &Logging{
							Logger: "glg",
							Level:  "warn",
						},
					},
					err: nil,
				},
			}
		}(),
		// TODO: Add error pattern
		// func() test {
		// 	data := []byte(`{vdaas}`)
		// 	return test{
		// 		name: "return unmarshal error when json data is invalid",
		// 		args: args{
		// 			data: data,
		// 		},
		// 		fields: fields{},
		// 		want: want{
		// 			want: &GlobalConfig{},
		// 			err:  errors.New(""),
		// 		},
		// 	}
		// }(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &GlobalConfig{
				Version: test.fields.Version,
				TZ:      test.fields.TZ,
				Logging: test.fields.Logging,
			}

			err := c.UnmarshalJSON(test.args.data)
			if err := test.checkFunc(test.want, c, err); err != nil {
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
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           path: "",
		           cfg: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           path: "",
		           cfg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := Read(test.args.path, test.args.cfg)
			if err := test.checkFunc(test.want, err); err != nil {
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
		beforeFunc func(args)
		afterFunc  func(args)
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
				name: "return v1.0.0. when val is _VERSION_",
				args: args{
					val: "_VERSION_",
				},
				beforeFunc: func(args) {
					os.Setenv("VERSION", "v1.0.0")
				},
				afterFunc: func(args) {
					os.Unsetenv("VERSION")
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
					val: "$VERSION",
				},
				beforeFunc: func(args) {
					os.Setenv("VERSION", "v1.0.0")
				},
				afterFunc: func(args) {
					os.Unsetenv("VERSION")
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
				name: "return file body contents when val is file://env",
				args: args{
					val: "file://" + fname,
				},
				beforeFunc: func(args) {
					f, err := os.Create(fname)
					if err != nil {
						t.Fatal(err)
					}
					defer f.Close()
					f.WriteString("v1.0.0")
				},
				afterFunc: func(args) {
					if err := os.Remove(fname); err != nil {
						t.Fatal(err)
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
				name: "return file body contents when val is file://env",
				args: args{
					val: "file://" + fname,
				},
				want: want{
					wantRes: "file://" + fname,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotRes := GetActualValue(test.args.val)
			if err := test.checkFunc(test.want, gotRes); err != nil {
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
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			env := map[string]string{
				"VERSION": "v1.0.0",
				"LOGGER":  "glg",
			}
			return test{
				name: "return v1.0.0 and glg when vals are _LOGGER_ and _VERSION_",
				args: args{
					vals: []string{
						"_VERSION_",
						"_LOGGER_",
					},
				},
				beforeFunc: func(args) {
					for key, val := range env {
						os.Setenv(key, val)
					}
				},
				afterFunc: func(args) {
					for key := range env {
						os.Unsetenv(key)
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
						"_VERSION_",
						"LOGGER",
					},
				},
				beforeFunc: func(args) {
					os.Setenv("VERSION", "v1.0.0")
				},
				afterFunc: func(args) {
					os.Unsetenv("VERSION")
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := GetActualValues(test.args.vals)
			if err := test.checkFunc(test.want, got); err != nil {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := checkPrefixAndSuffix(test.args.str, test.args.pref, test.args.suf)
			if err := test.checkFunc(test.want, got); err != nil {
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
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ToRawYaml(test.args.data)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
