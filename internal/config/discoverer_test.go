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
			m := map[string]string{
				"NAME":               "discoverer",
				"NAMESPACE":          "vald",
				"DISCOVERY_DURATION": "10ms",
			}
			return test{
				name: "return Discoverer when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					Name:              "_NAME_",
					Namespace:         "_NAMESPACE_",
					DiscoveryDuration: "_DISCOVERY_DURATION_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for key, val := range m {
						if err := os.Setenv(key, val); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for key := range m {
						if err := os.Unsetenv(key); err != nil {
							t.Fatal(err)
						}
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
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
						Duration:           "10ms",
						Client:             newGRPCClientConfig(),
						AgentClientOptions: newGRPCClientConfig(),
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
			key := "DURATION"
			val := "10ms"

			return test{
				name: "return DiscovererClient when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					Duration: "_" + key + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					if err := os.Setenv(key, val); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					if err := os.Unsetenv(key); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &DiscovererClient{
						Duration:           "10ms",
						Client:             newGRPCClientConfig(),
						AgentClientOptions: newGRPCClientConfig(),
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
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
