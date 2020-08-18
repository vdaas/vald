//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package mysql

import (
	"database/sql"
	"reflect"
	"testing"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func Test_metaVector_GetUUID(t *testing.T) {
	type fields struct {
		meta meta
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns UUID when UUID of meta is not empty",
			fields: fields{
				meta: meta{
					UUID: "vald-vector-01",
				},
			},
			want: want{
				want: "vald-vector-01",
			},
		},
		{
			name: "returns UUID when UUID of meta is empty string",
			fields: fields{
				meta: meta{
					UUID: "",
				},
			},
			want: want{
				want: "",
			},
			checkFunc: defaultCheckFunc,
		},
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
			m := &metaVector{
				meta: test.fields.meta,
			}

			got := m.GetUUID()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_metaVector_GetVector(t *testing.T) {
	type fields struct {
		meta meta
	}
	type want struct {
		want []byte
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []byte) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []byte) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			v := []byte("vdaas/vald")
			return test{
				name: "returns Vector when Vector of meta is not empty",
				fields: fields{
					meta: meta{
						Vector: v,
					},
				},
				want: want{
					want: v,
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			return test{
				name: "returns Vector when Vector of meta is empty",
				want: want{
					want: nil,
				},
				checkFunc: defaultCheckFunc,
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
			m := &metaVector{
				meta: test.fields.meta,
			}

			got := m.GetVector()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_metaVector_GetMeta(t *testing.T) {
	type fields struct {
		meta meta
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns MetaString when MetaString is not empty",
			fields: fields{
				meta: meta{
					Meta: dbr.NullString{
						sql.NullString{
							String: "vdaas/vald",
							Valid:  false,
						},
					},
				},
			},
			want: want{
				want: "vdaas/vald",
			},
		},
		{
			name: "returns MetaString when MetaString is empty",
			want: want{
				want: "",
			},
		},
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
			m := &metaVector{
				meta: test.fields.meta,
			}

			got := m.GetMeta()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_metaVector_GetIPs(t *testing.T) {
	type fields struct {
		podIPs []podIP
	}
	type want struct {
		want []string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns ips when podIP is not nil",
			fields: fields{
				podIPs: []podIP{
					{
						ID: 1,
						IP: "192.168.1.1",
					},
					{
						ID: 2,
						IP: "192.168.1.2",
					},
				},
			},
			want: want{
				want: []string{
					"192.168.1.1",
					"192.168.1.2",
				},
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "returns empty array when podIP is nil",
			want: want{
				want: []string{},
			},
			checkFunc: defaultCheckFunc,
		},
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
			m := &metaVector{
				podIPs: test.fields.podIPs,
			}

			got := m.GetIPs()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
