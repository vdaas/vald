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

func TestGateway_Bind(t *testing.T) {
	type fields struct {
		AgentPort      int
		AgentName      string
		AgentNamespace string
		AgentDNS       string
		NodeName       string
		IndexReplica   int
		Discoverer     *DiscovererClient
		Meta           *Meta
		BackupManager  *BackupManager
		EgressFilter   *EgressFilter
	}
	type want struct {
		want *Gateway
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Gateway) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Gateway) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			port := 8081
			name := "vald-agent-ngt-0"
			ns := "vald"
			dns := "vald-agent-ngt.vald.svc.local"
			node := "vald-prod"
			ireplica := 3
			return test{
				name: "return Gateway when only params related agent and dns are set",
				fields: fields{
					AgentPort:      port,
					AgentName:      name,
					AgentNamespace: ns,
					AgentDNS:       dns,
					NodeName:       node,
					IndexReplica:   ireplica,
				},
				want: want{
					want: &Gateway{
						AgentPort:      port,
						AgentName:      name,
						AgentNamespace: ns,
						AgentDNS:       dns,
						NodeName:       node,
						IndexReplica:   ireplica,
						Meta:           &Meta{},
					},
				},
			}
		}(),
		func() test {
			port := 8081
			name := "vald-agent-ngt-0"
			ns := "vald"
			dns := "vald-agent-ngt.vald.svc.local"
			node := "vald-prod"
			ireplica := 3
			disc := &DiscovererClient{
				Duration: "10m",
			}
			meta := &Meta{
				Host:                      "vald-meta.svc.local",
				Port:                      8081,
				EnableCache:               true,
				CacheExpiration:           "3m",
				ExpiredCacheCheckDuration: "10m",
			}
			bmanager := &BackupManager{}
			efilter := &EgressFilter{}
			return test{
				name: "return Gateway when all of params are set",
				fields: fields{
					AgentPort:      port,
					AgentName:      name,
					AgentNamespace: ns,
					AgentDNS:       dns,
					NodeName:       node,
					IndexReplica:   ireplica,
					Discoverer:     disc,
					Meta:           meta,
					BackupManager:  bmanager,
					EgressFilter:   efilter,
				},
				want: want{
					want: &Gateway{
						AgentPort:      port,
						AgentName:      name,
						AgentNamespace: ns,
						AgentDNS:       dns,
						NodeName:       node,
						IndexReplica:   ireplica,
						Discoverer: &DiscovererClient{
							Duration: "10m",
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
						Meta: &Meta{
							Host: "vald-meta.svc.local",
							Port: 8081,
							Client: &GRPCClient{
								Addrs: []string{
									"vald-meta.svc.local:8081",
								},
								DialOption: &DialOption{
									Insecure: true,
								},
							},
							EnableCache:               true,
							CacheExpiration:           "3m",
							ExpiredCacheCheckDuration: "10m",
						},
						BackupManager: &BackupManager{
							Client: &GRPCClient{
								DialOption: &DialOption{
									Insecure: true,
								},
							},
						},
						EgressFilter: &EgressFilter{},
					},
				},
			}
		}(),
		func() test {
			envPrefix := "GATEWAY_BIND_"
			p := map[string]string{
				envPrefix + "AGENT_NAME":      "vald-agent-ngt-0",
				envPrefix + "AGENT_NAMESPACE": "vald",
				envPrefix + "AGENT_DNS":       "vald-agent-ngt.svc.local",
				envPrefix + "NODE_NAME":       "vald-prod",
			}
			port := 8081
			ireplica := 3
			return test{
				name: "return Gateway when params set as environment value",
				fields: fields{
					AgentPort:      port,
					AgentName:      "_" + envPrefix + "AGENT_NAME_",
					AgentNamespace: "_" + envPrefix + "AGENT_NAMESPACE_",
					AgentDNS:       "_" + envPrefix + "AGENT_DNS_",
					NodeName:       "_" + envPrefix + "NODE_NAME_",
					IndexReplica:   ireplica,
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &Gateway{
						AgentPort:      8081,
						AgentName:      "vald-agent-ngt-0",
						AgentNamespace: "vald",
						AgentDNS:       "vald-agent-ngt.svc.local",
						NodeName:       "vald-prod",
						IndexReplica:   3,
						Meta:           &Meta{},
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Gateway when all params are not set",
				fields: fields{},
				want: want{
					want: &Gateway{
						Meta: &Meta{},
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
			g := &Gateway{
				AgentPort:      test.fields.AgentPort,
				AgentName:      test.fields.AgentName,
				AgentNamespace: test.fields.AgentNamespace,
				AgentDNS:       test.fields.AgentDNS,
				NodeName:       test.fields.NodeName,
				IndexReplica:   test.fields.IndexReplica,
				Discoverer:     test.fields.Discoverer,
				Meta:           test.fields.Meta,
				BackupManager:  test.fields.BackupManager,
				EgressFilter:   test.fields.EgressFilter,
			}

			got := g.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
