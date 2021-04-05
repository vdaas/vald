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

func TestLB_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		AgentPort      int
		AgentName      string
		AgentNamespace string
		AgentDNS       string
		NodeName       string
		IndexReplica   int
		Discoverer     *DiscovererClient
	}
	type want struct {
		want *LB
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *LB) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *LB) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			agentPort := 8081
			agentName := "vald-agent-ngt"
			agentNamespace := "vald"
			agentDNS := "vald-agent-ngt.vald.svc.cluster.local"
			nodeName := "vald-01-worker"
			indexReplica := 3

			return test{
				name: "return LB when the bind successes and the Discoverer is nil",
				fields: fields{
					AgentPort:      agentPort,
					AgentName:      agentName,
					AgentNamespace: agentNamespace,
					AgentDNS:       agentDNS,
					NodeName:       nodeName,
					IndexReplica:   indexReplica,
				},
				want: want{
					want: &LB{
						AgentPort:      agentPort,
						AgentName:      agentName,
						AgentNamespace: agentNamespace,
						AgentDNS:       agentDNS,
						NodeName:       nodeName,
						IndexReplica:   indexReplica,
					},
				},
			}
		}(),
		func() test {
			agentPort := 8081
			agentName := "vald-agent-ngt"
			agentNamespace := "vald"
			agentDNS := "vald-agent-ngt.vald.svc.cluster.local"
			nodeName := "vald-01-worker"
			indexReplica := 3

			return test{
				name: "return LB when the bind successes and the Discoverer is not nil",
				fields: fields{
					AgentPort:      agentPort,
					AgentName:      agentName,
					AgentNamespace: agentNamespace,
					AgentDNS:       agentDNS,
					NodeName:       nodeName,
					IndexReplica:   indexReplica,
					Discoverer:     new(DiscovererClient),
				},
				want: want{
					want: &LB{
						AgentPort:      agentPort,
						AgentName:      agentName,
						AgentNamespace: agentNamespace,
						AgentDNS:       agentDNS,
						NodeName:       nodeName,
						IndexReplica:   indexReplica,
						Discoverer: &DiscovererClient{
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
				},
			}
		}(),
		func() test {
			agentName := "vald-agent-ngt"
			agentNamespace := "vald"
			agentDNS := "vald-agent-ngt.vald.svc.cluster.local"
			nodeName := "vald-01-worker"

			m := map[string]string{
				"AGENT_NAME":      agentName,
				"AGENT_NAMESPACE": agentNamespace,
				"AGENT_DNS":       agentDNS,
				"NODE_NAME":       nodeName,
			}
			return test{
				name: "return LB when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					AgentPort:      8081,
					AgentName:      "_AGENT_NAME_",
					AgentNamespace: "_AGENT_NAMESPACE_",
					AgentDNS:       "_AGENT_DNS_",
					NodeName:       "_NODE_NAME_",
					IndexReplica:   3,
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
					want: &LB{
						AgentPort:      8081,
						AgentName:      agentName,
						AgentNamespace: agentNamespace,
						AgentDNS:       agentDNS,
						NodeName:       nodeName,
						IndexReplica:   3,
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
			g := &LB{
				AgentPort:      test.fields.AgentPort,
				AgentName:      test.fields.AgentName,
				AgentNamespace: test.fields.AgentNamespace,
				AgentDNS:       test.fields.AgentDNS,
				NodeName:       test.fields.NodeName,
				IndexReplica:   test.fields.IndexReplica,
				Discoverer:     test.fields.Discoverer,
			}

			got := g.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
