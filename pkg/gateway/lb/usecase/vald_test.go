//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package usecase represents gateways usecase layer
package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/config"

	"github.com/vdaas/vald/internal/net/grpc"
)

func Test_discovererOpts(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.Data
		dopts  []grpc.Option
		aopts  []grpc.Option
		assert func(*testing.T, []discoverer.Option, error)
	}{
		{
			name: "Not create read replica client when read replica client option is not set",
			cfg: &config.Data{
				Gateway: &iconfig.LB{
					AgentName:      "agent",
					AgentNamespace: "agent-ns",
					AgentPort:      8081,
					AgentDNS:       "agent-dns",
					Discoverer: &iconfig.DiscovererClient{
						Duration: "1m",
					},
					NodeName: "node",
				},
			},
			dopts: []grpc.Option{},
			aopts: []grpc.Option{},
			assert: func(t *testing.T, opts []discoverer.Option, err error) {
				require.NoError(t, err)

				client, err := discoverer.New(opts...)
				require.NoError(t, err)

				// check multiple times to ensure that the client is not a read replica client
				require.Equal(t, client.GetClient(), client.GetReadClient())
				require.Equal(t, client.GetClient(), client.GetReadClient())
				require.Equal(t, client.GetClient(), client.GetReadClient())
			},
		},
		{
			name: "create read replica client when read replica client option is set",
			cfg: &config.Data{
				Gateway: &iconfig.LB{
					AgentName:      "agent",
					AgentNamespace: "agent-ns",
					AgentPort:      8081,
					AgentDNS:       "agent-dns",
					Discoverer: &iconfig.DiscovererClient{
						Duration: "1m",
					},
					NodeName: "node",
					ReadReplicaClient: iconfig.ReadReplicaClient{
						Client: &iconfig.GRPCClient{},
					},
					// set this to big enough value to ensure that the round robin counter won't reset to 0
					ReadReplicaReplicas: 100,
				},
			},
			dopts: []grpc.Option{},
			aopts: []grpc.Option{},
			assert: func(t *testing.T, opts []discoverer.Option, err error) {
				require.NoError(t, err)

				client, err := discoverer.New(opts...)
				require.NoError(t, err)

				// ensure that GetReadClient() returns a read replica client by calling it multiple times beforehand
				// and increments the round robin counter
				client.GetReadClient()
				client.GetReadClient()
				client.GetReadClient()

				require.NotEqual(t, client.GetClient(), client.GetReadClient())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := discovererOpts(tt.cfg, tt.dopts, tt.aopts, errgroup.Get())
			tt.assert(t, opts, err)
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		cfg *config.Data
// 	}
// 	type want struct {
// 		wantR runner.Runner
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, runner.Runner, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotR runner.Runner, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotR, w.wantR) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotR, w.wantR)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           cfg:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotR, err := New(test.args.cfg)
// 			if err := checkFunc(test.want, gotR, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
