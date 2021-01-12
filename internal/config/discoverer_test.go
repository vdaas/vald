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

func TestDiscoverer_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name              string
		Namespace         string
		DiscoveryDuration string
	}
	type want struct {
		want *Discoverer
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Discoverer) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Discoverer) error {
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
		           Name: "",
		           Namespace: "",
		           DiscoveryDuration: "",
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
		           Name: "",
		           Namespace: "",
		           DiscoveryDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
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
			d := &Discoverer{
				Name:              test.fields.Name,
				Namespace:         test.fields.Namespace,
				DiscoveryDuration: test.fields.DiscoveryDuration,
			}

			got := d.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDiscovererClient_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Duration           string
		Client             *GRPCClient
		AgentClientOptions *GRPCClient
	}
	type want struct {
		want *DiscovererClient
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *DiscovererClient) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *DiscovererClient) error {
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
		           Duration: "",
		           Client: GRPCClient{},
		           AgentClientOptions: GRPCClient{},
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
		           Duration: "",
		           Client: GRPCClient{},
		           AgentClientOptions: GRPCClient{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
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
			d := &DiscovererClient{
				Duration:           test.fields.Duration,
				Client:             test.fields.Client,
				AgentClientOptions: test.fields.AgentClientOptions,
			}

			got := d.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
