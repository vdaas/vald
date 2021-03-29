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

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestBackupManager_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Client *GRPCClient
	}
	type want struct {
		want *BackupManager
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *BackupManager) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *BackupManager) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return Backup when client is nil",
			fields: fields{
				Client: nil,
			},
			want: want{
				want: &BackupManager{
					Client: newGRPCClientConfig(),
				},
			},
		},
		func() test {
			c := &GRPCClient{
				HealthCheckDuration: "1s",
			}
			return test{
				name: "return Backup when client is not nil",
				fields: fields{
					Client: c,
				},
				want: want{
					want: &BackupManager{
						Client: c.Bind(),
					},
				},
			}
		}(),
		func() test {
			k := "ADDR"
			v := "http://backupmanager.com"
			c := &GRPCClient{
				Addrs: []string{v},
			}
			return test{
				name: "return Backup when addrs is set via environment variable",
				fields: fields{
					Client: &GRPCClient{
						Addrs: []string{
							"_" + k + "_",
						},
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					if err := os.Setenv(k, v); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					if err := os.Unsetenv(k); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &BackupManager{
						Client: c.Bind(),
					},
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &BackupManager{
				Client: test.fields.Client,
			}

			got := b.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
