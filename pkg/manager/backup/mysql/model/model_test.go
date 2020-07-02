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

// Package grpc provides grpc server logic
package model

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestMetaVector_GetUUID(t *testing.T) {
	type fields struct {
		UUID   string
		Vector []byte
		Meta   string
		IPs    []string
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MetaVector{
				UUID:   test.fields.UUID,
				Vector: test.fields.Vector,
				Meta:   test.fields.Meta,
				IPs:    test.fields.IPs,
			}

			got := m.GetUUID()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestMetaVector_GetVector(t *testing.T) {
	type fields struct {
		UUID   string
		Vector []byte
		Meta   string
		IPs    []string
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MetaVector{
				UUID:   test.fields.UUID,
				Vector: test.fields.Vector,
				Meta:   test.fields.Meta,
				IPs:    test.fields.IPs,
			}

			got := m.GetVector()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestMetaVector_GetMeta(t *testing.T) {
	type fields struct {
		UUID   string
		Vector []byte
		Meta   string
		IPs    []string
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MetaVector{
				UUID:   test.fields.UUID,
				Vector: test.fields.Vector,
				Meta:   test.fields.Meta,
				IPs:    test.fields.IPs,
			}

			got := m.GetMeta()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestMetaVector_GetIPs(t *testing.T) {
	type fields struct {
		UUID   string
		Vector []byte
		Meta   string
		IPs    []string
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
		           UUID: "",
		           Vector: nil,
		           Meta: "",
		           IPs: nil,
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MetaVector{
				UUID:   test.fields.UUID,
				Vector: test.fields.Vector,
				Meta:   test.fields.Meta,
				IPs:    test.fields.IPs,
			}

			got := m.GetIPs()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
