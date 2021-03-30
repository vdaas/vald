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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestDebug_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Profile struct {
			Enable bool    `yaml:"enable" json:"enable"`
			Server *Server `yaml:"server" json:"server"`
		}
		Log struct {
			Level string `yaml:"level" json:"level"`
			Mode  string `yaml:"mode" json:"mode"`
		}
	}
	type want struct {
		want *Debug
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Debug) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Debug) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			p := struct {
				Enable bool    `yaml:"enable" json:"enable"`
				Server *Server `yaml:"server" json:"server"`
			}{
				Enable: true,
				Server: &Server{
					Name:          "",
					Network:       "tcp",
					Host:          "",
					Port:          8081,
					SocketPath:    "",
					Mode:          "GRPC",
					ProbeWaitTime: "3s",
					HTTP:          &HTTP{},
					GRPC:          &GRPC{},
					SocketOption:  &SocketOption{},
					Restart:       false,
				},
			}
			log := struct {
				Level string `yaml:"level" json:"level"`
				Mode  string `yaml:"mode" json:"mode"`
			}{
				Level: "Error",
				Mode:  "raw",
			}
			return test{
				name: "return the Debug when all variables are not empty",
				fields: fields{
					Profile: p,
					Log:     log,
				},
				want: want{
					want: &Debug{
						Profile: p,
						Log:     log,
					},
				},
			}
		}(),
		func() test {
			p := struct {
				Enable bool    `yaml:"enable" json:"enable"`
				Server *Server `yaml:"server" json:"server"`
			}{
				Enable: false,
			}
			log := struct {
				Level string `yaml:"level" json:"level"`
				Mode  string `yaml:"mode" json:"mode"`
			}{
				Level: "Error",
				Mode:  "raw",
			}
			return test{
				name: "return the Debug when all variable are not empty except server's paramter",
				fields: fields{
					Profile: p,
					Log:     log,
				},
				want: want{
					want: &Debug{
						Profile: p,
						Log:     log,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return the Debug when all variables are empty",
				fields: fields{},
				want: want{
					want: &Debug{},
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			d := &Debug{
				Profile: test.fields.Profile,
				Log:     test.fields.Log,
			}

			got := d.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
