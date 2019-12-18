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
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		checkFunc func(Option) error
	}{
		{
			name: "set name success",
			args: args{
				name: "dummyName",
			},
			checkFunc: func(o Option) error {
				r := &runner{}
				o(r)
				if r.name != "dummyName" {
					return errors.New("cannot set value")
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithName(tt.args.name)
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
			name: "set ver success",
			args: args{
				ver: "dummyVer",
			},
			checkFunc: func(o Option) error {
				r := &runner{}
				o(r)
				if r.version != "dummyVer" {
					return errors.New("cannot set value")
				}
				if r.maxVersion != "" {
					return errors.New("maxVersion should be empty")
				}
				if r.minVersion != "" {
					return errors.New("minVersion should be empty")
				}
				return nil
			},
		},
		{
			name: "set min ver success",
			args: args{
				min: "dummyMinVer",
			},
			checkFunc: func(o Option) error {
				r := &runner{}
				o(r)
				if r.version != "" {
					return errors.New("version should be empty")
				}
				if r.maxVersion != "" {
					return errors.New("maxVersion should be empty")
				}
				if r.minVersion != "dummyMinVer" {
					return errors.New("cannot set value")
				}
				return nil
			},
		},
		{
			name: "set max ver success",
			args: args{
				max: "dummyMaxVer",
			},
			checkFunc: func(o Option) error {
				r := &runner{}
				o(r)
				if r.version != "" {
					return errors.New("version should be empty")
				}
				if r.maxVersion != "dummyMaxVer" {
					return errors.New("cannot set value")
				}
				if r.minVersion != "" {
					return errors.New("minVersion should be empty")
				}
				return nil
			},
		},
		{
			name: "set all success",
			args: args{
				ver: "dummyVer",
				min: "dummyMinVer",
				max: "dummyMaxVer",
			},
			checkFunc: func(o Option) error {
				r := &runner{}
				o(r)
				if r.version != "dummyVer" {
					return errors.New("cannot set value")
				}
				if r.maxVersion != "dummyMaxVer" {
					return errors.New("cannot set value")
				}
				if r.minVersion != "dummyMinVer" {
					return errors.New("cannot set value")
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
	type args struct {
		f func(string) (interface{}, string, error)
	}
	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}
	tests := []test{
		func() test {
			f := func(string) (interface{}, string, error) {
				return nil, "", nil
			}
			return test{
				name: "set config load success",
				args: args{f: f},
				checkFunc: func(o Option) error {
					r := &runner{}
					o(r)
					if reflect.ValueOf(r.loadConfig) != reflect.ValueOf(f) {
						return errors.New("cannot set value")
					}
					return nil
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithConfigLoader(tt.args.f)
			if err := tt.checkFunc(got); err != nil {
				t.Errorf("WithConfigLoader() error = %v", err)
			}
		})
	}
}

func TestWithDaemonInitializer(t *testing.T) {
	type args struct {
		f func(interface{}) (Runner, error)
	}
	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}
	tests := []test{
		func() test {
			f := func(interface{}) (Runner, error) {
				return nil, nil
			}
			return test{
				name: "set config load success",
				args: args{f: f},
				checkFunc: func(o Option) error {
					r := &runner{}
					o(r)
					if reflect.ValueOf(r.initializeDaemon) != reflect.ValueOf(f) {
						return errors.New("cannot set value")
					}
					return nil
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithDaemonInitializer(tt.args.f)
			if err := tt.checkFunc(got); err != nil {
				t.Errorf("WithDaemonInitializer() error = %v", err)
			}
		})
	}
}
