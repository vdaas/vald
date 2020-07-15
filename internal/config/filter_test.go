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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestEgressFilter_Bind(t *testing.T) {
	type fields struct {
		Client *GRPCClient
	}
	type want struct {
		want *EgressFilter
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *EgressFilter) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *EgressFilter) error {
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
		           Client: GRPCClient{},
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
		           Client: GRPCClient{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			e := &EgressFilter{
				Client: test.fields.Client,
			}

			got := e.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestIngressFilter_Bind(t *testing.T) {
	type fields struct {
		Client *GRPCClient
		Search []string
		Insert []string
		Update []string
		Upsert []string
	}
	type want struct {
		want *IngressFilter
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *IngressFilter) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *IngressFilter) error {
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
		           Client: GRPCClient{},
		           Search: nil,
		           Insert: nil,
		           Update: nil,
		           Upsert: nil,
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
		           Client: GRPCClient{},
		           Search: nil,
		           Insert: nil,
		           Update: nil,
		           Upsert: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			i := &IngressFilter{
				Client: test.fields.Client,
				Search: test.fields.Search,
				Insert: test.fields.Insert,
				Update: test.fields.Update,
				Upsert: test.fields.Upsert,
			}

			got := i.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
