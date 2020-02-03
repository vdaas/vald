//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

// Package runner provides implementation of process runner
package runner

import (
	"errors"
	"reflect"
	"testing"
)

func TestWithName(t *testing.T) {
	type test struct {
		name      string
		str       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when str is not empty",
			str:  "name",
			checkFunc: func(o Option) error {
				got := new(runner)
				o(got)

				if got.name != "str" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when str is empty",
			checkFunc: func(o Option) error {
				got := &runner{
					name: "name",
				}
				o(got)

				if got.name != "name" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithName(tt.str)
			if err := tt.checkFunc(got); err != nil {
				t.Errorf("WithName() error = %v", err)
			}
		})
	}
}

func TestWithVersion(t *testing.T) {
	type args struct {
		ver string
		max string
		min string
	}
	tests := []struct {
		name      string
		args      args
		checkFunc func(Option) error
	}{
		{
			name: "set ver success when ver is not empty",
			args: args{
				ver: "ver",
			},
			checkFunc: func(o Option) error {
				got := new(runner)
				o(got)

				if got.version != "ver" {
					return errors.New("invalid param was set")
				}

				if got.maxVersion != "" {
					return errors.New("maxVersion should be empty")
				}

				if got.minVersion != "" {
					return errors.New("minVersion should be empty")
				}
				return nil
			},
		},

		{
			name: "not set ver when ver is empty",
			checkFunc: func(o Option) error {
				got := &runner{
					version: "ver",
				}
				o(got)

				if got.version != "ver" {
					return errors.New("invalid param was set")
				}

				if got.maxVersion != "" {
					return errors.New("maxVersion should be empty")
				}

				if got.minVersion != "" {
					return errors.New("minVersion should be empty")
				}
				return nil
			},
		},

		{
			name: "set min success when min is not empty",
			args: args{
				min: "min",
			},
			checkFunc: func(o Option) error {
				got := new(runner)
				o(got)

				if got.version != "" {
					return errors.New("version should be empty")
				}

				if got.maxVersion != "" {
					return errors.New("maxVersion should be empty")
				}

				if got.minVersion != "min" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set min when min is empty",
			checkFunc: func(o Option) error {
				got := &runner{
					minVersion: "min",
				}
				o(got)

				if got.version != "" {
					return errors.New("version should be empty")
				}

				if got.maxVersion != "" {
					return errors.New("maxVersion should be empty")
				}

				if got.minVersion != "min" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set max success when max is not empty",
			args: args{
				max: "max",
			},
			checkFunc: func(o Option) error {
				got := new(runner)
				o(got)

				if got.version != "" {
					return errors.New("version should be empty")
				}

				if got.maxVersion != "max" {
					return errors.New("invalid param was set")
				}

				if got.minVersion != "" {
					return errors.New("minVersion should be empty")
				}
				return nil
			},
		},

		{
			name: "not set max when max is empty",
			checkFunc: func(o Option) error {
				got := &runner{
					maxVersion: "max",
				}
				o(got)

				if got.version != "" {
					return errors.New("version should be empty")
				}

				if got.maxVersion != "max" {
					return errors.New("invalid param was set")
				}

				if got.minVersion != "" {
					return errors.New("minVersion should be empty")
				}
				return nil
			},
		},

		{
			name: "set all success",
			args: args{
				ver: "ver",
				min: "min",
				max: "max",
			},
			checkFunc: func(o Option) error {
				got := new(runner)
				o(got)

				if got.version != "ver" {
					return errors.New("invalid ver param was set")
				}

				if got.maxVersion != "max" {
					return errors.New("invalid max param was set")
				}

				if got.minVersion != "min" {
					return errors.New("invalid min param was set")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithVersion(tt.args.ver, tt.args.max, tt.args.min)
			if err := tt.checkFunc(got); err != nil {
				t.Errorf("WithVersion() error = %v", err)
			}
		})
	}
}

func TestWithConfigLoader(t *testing.T) {
	type test struct {
		name      string
		f         func(string) (interface{}, string, string, error)
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			f := func(string) (interface{}, string, string, error) {
				return nil, "", "", nil
			}
			return test{
				name: "set success when f is not nil",
				f:    f,
				checkFunc: func(o Option) error {
					got := new(runner)
					o(got)

					if reflect.ValueOf(got.loadConfig).Pointer() != reflect.ValueOf(f).Pointer() {
						return errors.New("invalid min param was set")
					}
					return nil
				},
			}
		}(),

		{
			name: "not set when f is nil",
			checkFunc: func(o Option) error {
				f := func(string) (interface{}, string, string, error) {
					return nil, "", "", nil
				}
				got := &runner{
					loadConfig: f,
				}
				o(got)

				if reflect.ValueOf(got.loadConfig).Pointer() != reflect.ValueOf(f).Pointer() {
					return errors.New("invalid min param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithConfigLoader(tt.f)
			if err := tt.checkFunc(got); err != nil {
				t.Errorf("WithConfigLoader() error = %v", err)
			}
		})
	}
}

func TestWithDaemonInitializer(t *testing.T) {
	type test struct {
		name      string
		f         func(interface{}) (Runner, error)
		checkFunc func(Option) error
	}
	tests := []test{
		func() test {
			f := func(interface{}) (Runner, error) {
				return nil, nil
			}
			return test{
				name: "set success when f is not nil",
				f:    f,
				checkFunc: func(o Option) error {
					got := new(runner)
					o(got)

					if reflect.ValueOf(got.initializeDaemon).Pointer() != reflect.ValueOf(f).Pointer() {
						return errors.New("invalid min param was set")
					}
					return nil
				},
			}
		}(),

		{
			name: "not set when f is nil",
			checkFunc: func(o Option) error {
				f := func(interface{}) (Runner, error) {
					return nil, nil
				}
				got := &runner{
					initializeDaemon: f,
				}
				o(got)

				if reflect.ValueOf(got.initializeDaemon).Pointer() != reflect.ValueOf(f).Pointer() {
					return errors.New("invalid min param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithDaemonInitializer(tt.f)
			if err := tt.checkFunc(got); err != nil {
				t.Errorf("WithDaemonInitializer() error = %v", err)
			}
		})
	}
}
