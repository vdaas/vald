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

func TestIndexer_Bind(t *testing.T) {
	type fields struct {
		AgentPort              int
		AgentName              string
		AgentNamespace         string
		AgentDNS               string
		Concurrency            int
		AutoIndexDurationLimit string
		AutoIndexCheckDuration string
		AutoIndexLength        uint32
		CreationPoolSize       uint32
		NodeName               string
		Discoverer             *DiscovererClient
	}
	type want struct {
		want *Indexer
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Indexer) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Indexer) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return Indexer when the bind successes",
				fields: fields{
					AgentPort:              8081,
					AgentName:              "vald-agent-ngt",
					AgentNamespace:         "vald",
					AgentDNS:               "vald-agent-ngt.vald.svc.cluster.local",
					Concurrency:            10,
					AutoIndexDurationLimit: "30m",
					AutoIndexCheckDuration: "1m",
					AutoIndexLength:        100,
					CreationPoolSize:       10000,
					NodeName:               "vald-01-worker",
				},
				want: want{
					want: &Indexer{
						AgentPort:              8081,
						AgentName:              "vald-agent-ngt",
						AgentNamespace:         "vald",
						AgentDNS:               "vald-agent-ngt.vald.svc.cluster.local",
						Concurrency:            10,
						AutoIndexDurationLimit: "30m",
						AutoIndexCheckDuration: "1m",
						AutoIndexLength:        100,
						CreationPoolSize:       10000,
						NodeName:               "vald-01-worker",
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return Indexer when the bind successes and Discoverer is not nil",
				fields: fields{
					AgentPort:              8081,
					AgentName:              "vald-agent-ngt",
					AgentNamespace:         "vald",
					AgentDNS:               "vald-agent-ngt.vald.svc.cluster.local",
					Concurrency:            10,
					AutoIndexDurationLimit: "30m",
					AutoIndexCheckDuration: "1m",
					AutoIndexLength:        100,
					CreationPoolSize:       10000,
					NodeName:               "vald-01-worker",
					Discoverer:             new(DiscovererClient),
				},
				want: want{
					want: &Indexer{
						AgentPort:              8081,
						AgentName:              "vald-agent-ngt",
						AgentNamespace:         "vald",
						AgentDNS:               "vald-agent-ngt.vald.svc.cluster.local",
						Concurrency:            10,
						AutoIndexDurationLimit: "30m",
						AutoIndexCheckDuration: "1m",
						AutoIndexLength:        100,
						CreationPoolSize:       10000,
						NodeName:               "vald-01-worker",
						Discoverer: &DiscovererClient{
							Client:             newGRPCClientConfig(),
							AgentClientOptions: newGRPCClientConfig(),
						},
					},
				},
			}
		}(),
		func() test {
			m := map[string]string{
				"AGENT_NAME":                "vald-agent-ngt",
				"AGENT_NAMESPACE":           "vald",
				"AGENT_DNS":                 "vald-agent-ngt.vald.svc.cluster.local",
				"AUTO_INDEX_DURATION_LIMIT": "30m",
				"AUTO_INDEX_CHECK_DURATION": "1m",
				"NODE_NAME":                 "vald-01-worker",
			}

			return test{
				name: "return Indexer when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					AgentPort:              8081,
					AgentName:              "_AGENT_NAME_",
					AgentNamespace:         "_AGENT_NAMESPACE_",
					AgentDNS:               "_AGENT_DNS_",
					Concurrency:            10,
					AutoIndexDurationLimit: "_AUTO_INDEX_DURATION_LIMIT_",
					AutoIndexCheckDuration: "_AUTO_INDEX_CHECK_DURATION_",
					AutoIndexLength:        100,
					CreationPoolSize:       10000,
					NodeName:               "_NODE_NAME_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range m {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &Indexer{
						AgentPort:              8081,
						AgentName:              "vald-agent-ngt",
						AgentNamespace:         "vald",
						AgentDNS:               "vald-agent-ngt.vald.svc.cluster.local",
						Concurrency:            10,
						AutoIndexDurationLimit: "30m",
						AutoIndexCheckDuration: "1m",
						AutoIndexLength:        100,
						CreationPoolSize:       10000,
						NodeName:               "vald-01-worker",
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
			im := &Indexer{
				AgentPort:              test.fields.AgentPort,
				AgentName:              test.fields.AgentName,
				AgentNamespace:         test.fields.AgentNamespace,
				AgentDNS:               test.fields.AgentDNS,
				Concurrency:            test.fields.Concurrency,
				AutoIndexDurationLimit: test.fields.AutoIndexDurationLimit,
				AutoIndexCheckDuration: test.fields.AutoIndexCheckDuration,
				AutoIndexLength:        test.fields.AutoIndexLength,
				CreationPoolSize:       test.fields.CreationPoolSize,
				NodeName:               test.fields.NodeName,
				Discoverer:             test.fields.Discoverer,
			}

			got := im.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
