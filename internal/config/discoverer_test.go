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

	"github.com/pkg/errors"
)

func TestDiscoverer_Bind(t *testing.T) {
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
	type fields struct {
		Host        string
		Port        int
		Duration    string
		Client      *GRPCClient
		AgentClient *GRPCClient
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
		           Host: "",
		           Port: 0,
		           Duration: "",
		           Client: GRPCClient{},
		           AgentClient: GRPCClient{},
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
		           Host: "",
		           Port: 0,
		           Duration: "",
		           Client: GRPCClient{},
		           AgentClient: GRPCClient{},
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
			d := &DiscovererClient{
				Host:        test.fields.Host,
				Port:        test.fields.Port,
				Duration:    test.fields.Duration,
				Client:      test.fields.Client,
				AgentClient: test.fields.AgentClient,
			}

			got := d.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
