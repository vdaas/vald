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

package mysql

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_vector_GetUUID(t *testing.T) {
	type fields struct {
		data data
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns UUID when UUID of vector is not empty",
			fields: fields{
				data: data{
					UUID: "vald-vector-01",
				},
			},
			want: want{
				want: "vald-vector-01",
			},
		},
		{
			name: "returns UUID when UUID of vector is empty string",
			fields: fields{
				data: data{
					UUID: "",
				},
			},
			want: want{
				want: "",
			},
		},
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
			m := &vector{
				data: test.fields.data,
			}

			got := m.GetUUID()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vector_GetVector(t *testing.T) {
	type fields struct {
		data data
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			v := []byte("vdaas/vald")
			return test{
				name: "returns Vector when Vector of vector is not empty",
				fields: fields{
					data: data{
						Vector: v,
					},
				},
				want: want{
					want: v,
				},
			}
		}(),
		func() test {
			return test{
				name: "returns Vector when Vector of vector is empty",
				want: want{
					want: nil,
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
			m := &vector{
				data: test.fields.data,
			}

			got := m.GetVector()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vector_GetIPs(t *testing.T) {
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
		},
		{
			name: "returns empty array when podIP is nil",
			want: want{
				want: []string{},
			},
		},
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
			m := &vector{
				podIPs: test.fields.podIPs,
			}

			got := m.GetIPs()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
