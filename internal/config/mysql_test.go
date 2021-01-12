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

// Package config providers configuration type and load configuration logic
package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/db/rdb/mysql"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestMySQL_Bind(t *testing.T) {
	type fields struct {
		DB                   string
		Host                 string
		Port                 int
		User                 string
		Pass                 string
		Name                 string
		Charset              string
		Timezone             string
		InitialPingTimeLimit string
		InitialPingDuration  string
		ConnMaxLifeTime      string
		MaxOpenConns         int
		MaxIdleConns         int
		TLS                  *TLS
		TCP                  *TCP
	}
	type want struct {
		want *MySQL
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *MySQL) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *MySQL) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return MySQL when all fields contain no prefix/suffix symbol and tls and tcp configuration is not set",
			fields: fields{
				DB:                   "db",
				Host:                 "host",
				Port:                 80,
				User:                 "user",
				Pass:                 "pass",
				Name:                 "name",
				Charset:              "charset",
				Timezone:             "timezone",
				InitialPingTimeLimit: "initialPingTimeLimit",
				InitialPingDuration:  "initialPingDuration",
				ConnMaxLifeTime:      "connMaxLifeTime",
				MaxOpenConns:         10,
				MaxIdleConns:         100,
			},
			want: want{
				want: &MySQL{
					DB:                   "db",
					Host:                 "host",
					Port:                 80,
					User:                 "user",
					Pass:                 "pass",
					Name:                 "name",
					Charset:              "charset",
					Timezone:             "timezone",
					InitialPingTimeLimit: "initialPingTimeLimit",
					InitialPingDuration:  "initialPingDuration",
					ConnMaxLifeTime:      "connMaxLifeTime",
					MaxOpenConns:         10,
					MaxIdleConns:         100,
					TLS:                  new(TLS),
					TCP:                  new(TCP),
				},
			},
		},

		{
			name: "return MySQL when all fields contain no prefix/suffix symbol and tls and tcp configuration is set",
			fields: fields{
				DB:                   "db",
				Host:                 "host",
				Port:                 80,
				User:                 "user",
				Pass:                 "pass",
				Name:                 "name",
				Charset:              "charset",
				Timezone:             "timezone",
				InitialPingTimeLimit: "initialPingTimeLimit",
				InitialPingDuration:  "initialPingDuration",
				ConnMaxLifeTime:      "connMaxLifeTime",
				MaxOpenConns:         10,
				MaxIdleConns:         100,
				TLS: &TLS{
					Enabled: true,
				},
				TCP: &TCP{
					DNS: new(DNS),
				},
			},
			want: want{
				want: &MySQL{
					DB:                   "db",
					Host:                 "host",
					Port:                 80,
					User:                 "user",
					Pass:                 "pass",
					Name:                 "name",
					Charset:              "charset",
					Timezone:             "timezone",
					InitialPingTimeLimit: "initialPingTimeLimit",
					InitialPingDuration:  "initialPingDuration",
					ConnMaxLifeTime:      "connMaxLifeTime",
					MaxOpenConns:         10,
					MaxIdleConns:         100,
					TLS: &TLS{
						Enabled: true,
					},
					TCP: &TCP{
						DNS: new(DNS),
					},
				},
			},
		},

		{
			name: "return MySQL with environment variable when it contains `_` as prefix and suffix",
			fields: fields{
				DB:                   "_db_",
				Host:                 "_host_",
				Port:                 80,
				User:                 "_user_",
				Pass:                 "_pass_",
				Name:                 "_name_",
				Charset:              "_charset_",
				Timezone:             "_timezone_",
				InitialPingTimeLimit: "_initialPingTimeLimit_",
				InitialPingDuration:  "_initialPingDuration_",
				ConnMaxLifeTime:      "_connMaxLifeTime_",
				MaxOpenConns:         10,
				MaxIdleConns:         100,
				TLS:                  new(TLS),
				TCP:                  new(TCP),
			},
			want: want{
				want: &MySQL{
					DB:                   "db",
					Host:                 "host",
					Port:                 80,
					User:                 "user",
					Pass:                 "pass",
					Name:                 "name",
					Charset:              "charset",
					Timezone:             "timezone",
					InitialPingTimeLimit: "initialPingTimeLimit",
					InitialPingDuration:  "initialPingDuration",
					ConnMaxLifeTime:      "connMaxLifeTime",
					MaxOpenConns:         10,
					MaxIdleConns:         100,
					TLS:                  new(TLS),
					TCP:                  new(TCP),
				},
			},
			beforeFunc: func() {
				_ = os.Setenv("db", "db")
				_ = os.Setenv("host", "host")
				_ = os.Setenv("user", "user")
				_ = os.Setenv("pass", "pass")
				_ = os.Setenv("name", "name")
				_ = os.Setenv("charset", "charset")
				_ = os.Setenv("timezone", "timezone")
				_ = os.Setenv("initialPingTimeLimit", "initialPingTimeLimit")
				_ = os.Setenv("initialPingDuration", "initialPingDuration")
				_ = os.Setenv("connMaxLifeTime", "connMaxLifeTime")
			},
			afterFunc: func() {
				_ = os.Unsetenv("db")
				_ = os.Unsetenv("host")
				_ = os.Unsetenv("user")
				_ = os.Unsetenv("pass")
				_ = os.Unsetenv("name")
				_ = os.Unsetenv("charset")
				_ = os.Unsetenv("timezone")
				_ = os.Unsetenv("initialPingTimeLimit")
				_ = os.Unsetenv("initialPingDuration")
				_ = os.Unsetenv("connMaxLifeTime")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MySQL{
				DB:                   test.fields.DB,
				Host:                 test.fields.Host,
				Port:                 test.fields.Port,
				User:                 test.fields.User,
				Pass:                 test.fields.Pass,
				Name:                 test.fields.Name,
				Charset:              test.fields.Charset,
				Timezone:             test.fields.Timezone,
				InitialPingTimeLimit: test.fields.InitialPingTimeLimit,
				InitialPingDuration:  test.fields.InitialPingDuration,
				ConnMaxLifeTime:      test.fields.ConnMaxLifeTime,
				MaxOpenConns:         test.fields.MaxOpenConns,
				MaxIdleConns:         test.fields.MaxIdleConns,
				TLS:                  test.fields.TLS,
				TCP:                  test.fields.TCP,
			}

			got := m.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMySQL_Opts(t *testing.T) {
	t.Parallel()
	type fields struct {
		DB                   string
		Host                 string
		Port                 int
		User                 string
		Pass                 string
		Name                 string
		Charset              string
		Timezone             string
		InitialPingTimeLimit string
		InitialPingDuration  string
		ConnMaxLifeTime      string
		MaxOpenConns         int
		MaxIdleConns         int
		TLS                  *TLS
		TCP                  *TCP
	}
	type want struct {
		want []mysql.Option
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []mysql.Option, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []mysql.Option, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           DB: "",
		           Host: "",
		           Port: 0,
		           User: "",
		           Pass: "",
		           Name: "",
		           Charset: "",
		           Timezone: "",
		           InitialPingTimeLimit: "",
		           InitialPingDuration: "",
		           ConnMaxLifeTime: "",
		           MaxOpenConns: 0,
		           MaxIdleConns: 0,
		           TLS: TLS{},
		           TCP: TCP{},
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
		           fields: fields {
		           DB: "",
		           Host: "",
		           Port: 0,
		           User: "",
		           Pass: "",
		           Name: "",
		           Charset: "",
		           Timezone: "",
		           InitialPingTimeLimit: "",
		           InitialPingDuration: "",
		           ConnMaxLifeTime: "",
		           MaxOpenConns: 0,
		           MaxIdleConns: 0,
		           TLS: TLS{},
		           TCP: TCP{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MySQL{
				DB:                   test.fields.DB,
				Host:                 test.fields.Host,
				Port:                 test.fields.Port,
				User:                 test.fields.User,
				Pass:                 test.fields.Pass,
				Name:                 test.fields.Name,
				Charset:              test.fields.Charset,
				Timezone:             test.fields.Timezone,
				InitialPingTimeLimit: test.fields.InitialPingTimeLimit,
				InitialPingDuration:  test.fields.InitialPingDuration,
				ConnMaxLifeTime:      test.fields.ConnMaxLifeTime,
				MaxOpenConns:         test.fields.MaxOpenConns,
				MaxIdleConns:         test.fields.MaxIdleConns,
				TLS:                  test.fields.TLS,
				TCP:                  test.fields.TCP,
			}

			got, err := m.Opts()
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
