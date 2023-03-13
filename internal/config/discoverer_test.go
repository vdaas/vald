//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestDiscoverer_Bind(t *testing.T) {
	type fields struct {
		Name              string
		Namespace         string
		DiscoveryDuration string
		Net               *Net
	}
	type want struct {
		want *Discoverer
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Discoverer) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Discoverer) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return Discoverer when the bind successes",
				fields: fields{
					Name:              "discoverer",
					Namespace:         "vald",
					DiscoveryDuration: "10ms",
				},
				want: want{
					want: &Discoverer{
						Name:              "discoverer",
						Namespace:         "vald",
						DiscoveryDuration: "10ms",
						Net:               new(Net),
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return Discoverer when the bind successes and Net is not nil",
				fields: fields{
					Name:              "discoverer",
					Namespace:         "vald",
					DiscoveryDuration: "10ms",
					Net:               new(Net),
				},
				want: want{
					want: &Discoverer{
						Name:              "discoverer",
						Namespace:         "vald",
						DiscoveryDuration: "10ms",
						Net:               new(Net),
					},
				},
			}
		}(),
		func() test {
			suffix := "_FOR_TEST_DISCOVERER_BIND"
			m := map[string]string{
				"NAME" + suffix:               "discoverer",
				"NAMESPACE" + suffix:          "vald",
				"DISCOVERY_DURATION" + suffix: "10ms",
			}
			return test{
				name: "return Discoverer when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					Name:              "_NAME" + suffix + "_",
					Namespace:         "_NAMESPACE" + suffix + "_",
					DiscoveryDuration: "_DISCOVERY_DURATION" + suffix + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for key, val := range m {
						t.Setenv(key, val)
					}
				},
				want: want{
					want: &Discoverer{
						Name:              "discoverer",
						Namespace:         "vald",
						DiscoveryDuration: "10ms",
						Net:               new(Net),
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			d := &Discoverer{
				Name:              test.fields.Name,
				Namespace:         test.fields.Namespace,
				DiscoveryDuration: test.fields.DiscoveryDuration,
			}

			got := d.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDiscovererClient_Bind(t *testing.T) {
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
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *DiscovererClient) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return DiscovererClient when the bind successes",
				fields: fields{
					Duration: "10ms",
				},
				want: want{
					want: &DiscovererClient{
						Duration: "10ms",
						Client: &GRPCClient{
							DialOption: &DialOption{
								Insecure: true,
							},
						},
						AgentClientOptions: &GRPCClient{
							DialOption: &DialOption{
								Insecure: true,
							},
						},
					},
				},
			}
		}(),
		func() test {
			c := new(GRPCClient)
			ac := new(GRPCClient)

			return test{
				name: "return DiscovererClient when the bind successes and Client and AgentClientOptions is not nil",
				fields: fields{
					Duration:           "10ms",
					Client:             c,
					AgentClientOptions: ac,
				},
				want: want{
					want: &DiscovererClient{
						Duration: "10ms",
						Client: &GRPCClient{
							ConnectionPool: new(ConnectionPool),
							DialOption: &DialOption{
								Insecure: true,
							},
							TLS: new(TLS),
						},
						AgentClientOptions: &GRPCClient{
							ConnectionPool: new(ConnectionPool),
							DialOption: &DialOption{
								Insecure: true,
							},
							TLS: new(TLS),
						},
					},
				},
			}
		}(),
		func() test {
			key := "DURATION_FOR_TEST_DISCOVERER_CLIENT_BIND"
			val := "10ms"

			return test{
				name: "return DiscovererClient when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					Duration: "_" + key + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					t.Setenv(key, val)
				},
				want: want{
					want: &DiscovererClient{
						Duration: "10ms",
						Client: &GRPCClient{
							DialOption: &DialOption{
								Insecure: true,
							},
						},
						AgentClientOptions: &GRPCClient{
							DialOption: &DialOption{
								Insecure: true,
							},
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			d := &DiscovererClient{
				Duration:           test.fields.Duration,
				Client:             test.fields.Client,
				AgentClientOptions: test.fields.AgentClientOptions,
			}

			got := d.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
