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
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	type want struct {
		wantCfg *Config
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *Config, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotCfg *Config, err error) error {
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
			name := "/home/vankichi/Documents/vald-read-test.yaml"
			path := name
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
				checkFunc: func(w want, gotCfg *Config, err error) error {
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
					err:     errors.ErrUnsupportedConfigFileType(".yaml"),
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

			gotCfg, err := NewConfig(context.Background(), test.args.path)
			if err := checkFunc(test.want, gotCfg, err); err != nil {
				tt.Errorf("error = %v, got = %#v", err, gotCfg)
			}
		})
	}
}
