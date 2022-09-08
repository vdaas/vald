//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package cassandra

import (
	"context"
	"crypto/tls"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

// default comparator option for client.
var clientComparatorOpts = []comparator.Option{
	comparator.AllowUnexported(client{}),
	comparator.AllowUnexported(gocql.ClusterConfig{}),
	comparator.Comparer(func(x, y retryPolicy) bool {
		return reflect.DeepEqual(x, y)
	}),
	comparator.Comparer(func(x, y reconnectionPolicy) bool {
		return reflect.DeepEqual(x, y)
	}),
	comparator.Comparer(func(x, y poolConfig) bool {
		return reflect.DeepEqual(x, y)
	}),
	comparator.Comparer(func(x, y hostFilter) bool {
		return reflect.DeepEqual(x, y)
	}),
	comparator.Comparer(func(x, y gocql.PoolConfig) bool {
		return reflect.DeepEqual(x, y)
	}),
	comparator.Comparer(func(x, y gocql.HostSelectionPolicy) bool {
		return reflect.DeepEqual(x, y)
	}),
	comparator.Comparer(func(x, y func(h *gocql.HostInfo) (gocql.Authenticator, error)) bool {
		if (x == nil && y != nil) || (x != nil && y == nil) {
			return false
		}
		if x == nil && y == nil {
			return true
		}
		return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
	}),
	comparator.Comparer(func(x, y gocql.HostFilter) bool {
		if (x == nil && y != nil) || (x != nil && y == nil) {
			return false
		}
		if x == nil && y == nil {
			return true
		}

		switch x.(type) {
		case gocql.HostFilterFunc:
			return true
		}
		return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
	}),

	comparator.Comparer(func(x, y tls.Config) bool {
		return reflect.DeepEqual(x, y)
	}),
}

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Cassandra
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Cassandra, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Cassandra, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(got, w.want, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "New returns default cassandra",
			args: args{
				opts: nil,
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with and username password",
			args: args{
				opts: []Option{
					WithUsername("un"),
					WithPassword("p"),
				},
			},
			want: want{
				want: &client{
					username:                 "un",
					password:                 "p",
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return &gocql.PasswordAuthenticator{
								Username: "un",
								Password: "p",
							}
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with backoff",
			args: args{
				opts: []Option{
					WithRetryPolicyNumRetries(5),
					WithRetryPolicyMinDuration("1s"),
					WithRetryPolicyMaxDuration("5s"),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					retryPolicy: retryPolicy{
						numRetries:  5,
						minDuration: time.Second,
						maxDuration: 5 * time.Second,
					},
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return &gocql.ExponentialBackoffRetryPolicy{
								NumRetries: 5,
								Min:        time.Second,
								Max:        5 * time.Second,
							}
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with DC aware poll config",
			args: args{
				opts: []Option{
					WithDCAwareRouting(true),
					WithDC("dc"),
					WithTokenAwareHostPolicy(false),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.DCAwareRoundRobinPolicy("dc")
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           true,
						dataCenterName:                 "dc",
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     false,
					},
				},
			},
		},
		{
			name: "New returns cassandra with shuffle replica",
			args: args{
				opts: []Option{
					WithShuffleReplicas(true),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy(), gocql.ShuffleReplicas())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          true,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with non local replicas fallback pool",
			args: args{
				opts: []Option{
					WithNonLocalReplicasFallback(true),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy(), gocql.NonLocalReplicasFallback())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: true,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New default cassandra with shuffle replicas and non local replicas fallback",
			args: args{
				opts: []Option{
					WithShuffleReplicas(true),
					WithNonLocalReplicasFallback(true),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy(), gocql.ShuffleReplicas(), gocql.NonLocalReplicasFallback())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: true,
						enableShuffleReplicas:          true,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with host filter enable",
			args: args{
				opts: []Option{
					WithHostFilter(true),
					WithDCHostFilter("dc"),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					hostFilter: hostFilter{
						enable: true,
						dcHost: "dc",
					},
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},
						HostFilter: gocql.DataCentreHostFilter("dc"),

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with whitelist enable",
			args: args{
				opts: []Option{
					WithHostFilter(true),
					WithWhiteListHostFilter([]string{"localhost"}),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					hostFilter: hostFilter{
						enable:    true,
						whiteList: []string{"localhost"},
					},
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},
						HostFilter: gocql.WhiteListHostFilter("localhost"),

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with host filter and whitelist enabled",
			args: args{
				opts: []Option{
					WithHostFilter(true),
					WithDCHostFilter("dc"),
					WithWhiteListHostFilter([]string{"localhost"}),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					hostFilter: hostFilter{
						enable:    true,
						dcHost:    "dc",
						whiteList: []string{"localhost"},
					},
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},
						HostFilter: func() gocql.HostFilter {
							dchf := gocql.DataCentreHostFilter("dc")
							wlhf := gocql.WhiteListHostFilter("localhost")
							return gocql.HostFilterFunc(func(host *gocql.HostInfo) bool {
								return dchf.Accept(host) || wlhf.Accept(host)
							})
						}(),

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
		{
			name: "New returns cassandra with tls",
			args: args{
				opts: []Option{
					WithTLS(&tls.Config{
						MinVersion: tls.VersionTLS13,
					}),
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
						SslOpts: &gocql.SslOptions{
							Config: &tls.Config{
								MinVersion: tls.VersionTLS13,
							},
						},
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
					tls: &tls.Config{
						MinVersion: tls.VersionTLS13,
					},
				},
			},
		},
		{
			name: "New failed to init cassandra and return error",
			args: args{
				opts: []Option{
					func(*client) error {
						return errors.NewErrCriticalOption("n", 1)
					},
				},
			},
			want: want{
				err: errors.NewErrCriticalOption("n", 1),
			},
		},
		{
			name: "New returns default cassandra with option fail",
			args: args{
				opts: []Option{
					func(*client) error {
						return errors.New("err")
					},
				},
			},
			want: want{
				want: &client{
					cqlVersion:               "3.0.0",
					connectTimeout:           600 * time.Millisecond,
					consistency:              gocql.Quorum,
					defaultIdempotence:       false,
					defaultTimestamp:         true,
					disableInitialHostLookup: false,
					disableNodeStatusEvents:  false,
					disableSkipMetadata:      false,
					disableTopologyEvents:    false,
					enableHostVerification:   false,
					ignorePeerAddr:           false,
					maxPreparedStmts:         1000,
					maxRoutingKeyInfo:        1000,
					maxWaitSchemaAgreement:   1 * time.Minute,
					numConns:                 2,
					pageSize:                 5000,
					port:                     9042,
					protoVersion:             0,
					reconnectInterval:        time.Minute,
					serialConsistency:        gocql.LocalSerial,
					timeout:                  600 * time.Millisecond,
					writeCoalesceWaitTime:    200 * time.Microsecond,
					cluster: &gocql.ClusterConfig{
						Authenticator: func() *gocql.PasswordAuthenticator {
							return nil
						}(),
						RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
							return nil
						}(),
						ConvictionPolicy:   NewConvictionPolicy(),
						ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{},
						PoolConfig: gocql.PoolConfig{
							HostSelectionPolicy: func() gocql.HostSelectionPolicy {
								return gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
							}(),
						},

						CQLVersion:               "3.0.0",
						ConnectTimeout:           600 * time.Millisecond,
						Consistency:              gocql.Quorum,
						DefaultIdempotence:       false,
						DefaultTimestamp:         true,
						DisableInitialHostLookup: false,
						Events: events{
							DisableNodeStatusEvents: false,
							DisableTopologyEvents:   false,
						},
						DisableSkipMetadata:    false,
						IgnorePeerAddr:         false,
						MaxPreparedStmts:       1000,
						MaxRoutingKeyInfo:      1000,
						MaxWaitSchemaAgreement: 1 * time.Minute,
						NumConns:               2,
						PageSize:               5000,
						Port:                   9042,
						ProtoVersion:           0,
						ReconnectInterval:      time.Minute,
						SerialConsistency:      gocql.LocalSerial,
						Timeout:                600 * time.Millisecond,
						WriteCoalesceWaitTime:  200 * time.Microsecond,
					},
					poolConfig: poolConfig{
						enableDCAwareRouting:           false,
						enableNonLocalReplicasFallback: false,
						enableShuffleReplicas:          false,
						enableTokenAwareHostPolicy:     true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		hosts          []string
		cqlVersion     string
		protoVersion   int
		timeout        time.Duration
		connectTimeout time.Duration
		port           int
		keyspace       string
		numConns       int
		consistency    gocql.Consistency
		compressor     gocql.Compressor
		username       string
		password       string
		authProvider   func(h *gocql.HostInfo) (gocql.Authenticator, error)
		retryPolicy    struct {
			numRetries  int
			minDuration time.Duration
			maxDuration time.Duration
		}
		reconnectionPolicy struct {
			initialInterval time.Duration
			maxRetries      int
		}
		poolConfig struct {
			dataCenterName                 string
			enableDCAwareRouting           bool
			enableShuffleReplicas          bool
			enableNonLocalReplicasFallback bool
			enableTokenAwareHostPolicy     bool
		}
		hostFilter struct {
			enable    bool
			dcHost    string
			whiteList []string
		}
		socketKeepalive          time.Duration
		maxPreparedStmts         int
		maxRoutingKeyInfo        int
		pageSize                 int
		serialConsistency        gocql.SerialConsistency
		tls                      *tls.Config
		tlsCertPath              string
		tlsKeyPath               string
		tlsCAPath                string
		enableHostVerification   bool
		defaultTimestamp         bool
		reconnectInterval        time.Duration
		maxWaitSchemaAgreement   time.Duration
		ignorePeerAddr           bool
		disableInitialHostLookup bool
		disableNodeStatusEvents  bool
		disableTopologyEvents    bool
		disableSchemaEvents      bool
		disableSkipMetadata      bool
		defaultIdempotence       bool
		dialer                   gocql.Dialer
		writeCoalesceWaitTime    time.Duration
		cluster                  ClusterConfig
		session                  *gocql.Session
	}
	type want struct {
		c   *client
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(*client, want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(c *client, w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(c, w.c) {
			return errors.New("client is not equal")
		}
		return nil
	}
	tests := []test{
		func() test {
			cf := &MockClusterConfig{
				CreateSessionFunc: func() (*gocql.Session, error) {
					return &gocql.Session{}, nil
				},
			}

			return test{
				name: "open create session success",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					cluster: cf,
				},
				want: want{
					c: &client{
						cluster: cf,
						session: &gocql.Session{},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "open create session and return any error if occurred",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					cluster: &gocql.ClusterConfig{},
				},
				want: want{
					err: gocql.ErrNoHosts,
					c: &client{
						cluster: &gocql.ClusterConfig{},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, append(goleakIgnoreOptions, goleak.IgnoreTopFunction("github.com/gocql/gocql.(*eventDebouncer).flusher"))...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &client{
				hosts:                    test.fields.hosts,
				cqlVersion:               test.fields.cqlVersion,
				protoVersion:             test.fields.protoVersion,
				timeout:                  test.fields.timeout,
				connectTimeout:           test.fields.connectTimeout,
				port:                     test.fields.port,
				keyspace:                 test.fields.keyspace,
				numConns:                 test.fields.numConns,
				consistency:              test.fields.consistency,
				compressor:               test.fields.compressor,
				username:                 test.fields.username,
				password:                 test.fields.password,
				authProvider:             test.fields.authProvider,
				retryPolicy:              test.fields.retryPolicy,
				reconnectionPolicy:       test.fields.reconnectionPolicy,
				poolConfig:               test.fields.poolConfig,
				hostFilter:               test.fields.hostFilter,
				socketKeepalive:          test.fields.socketKeepalive,
				maxPreparedStmts:         test.fields.maxPreparedStmts,
				maxRoutingKeyInfo:        test.fields.maxRoutingKeyInfo,
				pageSize:                 test.fields.pageSize,
				serialConsistency:        test.fields.serialConsistency,
				tls:                      test.fields.tls,
				tlsCertPath:              test.fields.tlsCertPath,
				tlsKeyPath:               test.fields.tlsKeyPath,
				tlsCAPath:                test.fields.tlsCAPath,
				enableHostVerification:   test.fields.enableHostVerification,
				defaultTimestamp:         test.fields.defaultTimestamp,
				reconnectInterval:        test.fields.reconnectInterval,
				maxWaitSchemaAgreement:   test.fields.maxWaitSchemaAgreement,
				ignorePeerAddr:           test.fields.ignorePeerAddr,
				disableInitialHostLookup: test.fields.disableInitialHostLookup,
				disableNodeStatusEvents:  test.fields.disableNodeStatusEvents,
				disableTopologyEvents:    test.fields.disableTopologyEvents,
				disableSchemaEvents:      test.fields.disableSchemaEvents,
				disableSkipMetadata:      test.fields.disableSkipMetadata,
				defaultIdempotence:       test.fields.defaultIdempotence,
				dialer:                   test.fields.dialer,
				writeCoalesceWaitTime:    test.fields.writeCoalesceWaitTime,
				cluster:                  test.fields.cluster,
				session:                  test.fields.session,
			}

			err := c.Open(test.args.ctx)
			if err := checkFunc(c, test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Close(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		hosts          []string
		cqlVersion     string
		protoVersion   int
		timeout        time.Duration
		connectTimeout time.Duration
		port           int
		keyspace       string
		numConns       int
		consistency    gocql.Consistency
		compressor     gocql.Compressor
		username       string
		password       string
		authProvider   func(h *gocql.HostInfo) (gocql.Authenticator, error)
		retryPolicy    struct {
			numRetries  int
			minDuration time.Duration
			maxDuration time.Duration
		}
		reconnectionPolicy struct {
			initialInterval time.Duration
			maxRetries      int
		}
		poolConfig struct {
			dataCenterName                 string
			enableDCAwareRouting           bool
			enableShuffleReplicas          bool
			enableNonLocalReplicasFallback bool
			enableTokenAwareHostPolicy     bool
		}
		hostFilter struct {
			enable    bool
			dcHost    string
			whiteList []string
		}
		socketKeepalive          time.Duration
		maxPreparedStmts         int
		maxRoutingKeyInfo        int
		pageSize                 int
		serialConsistency        gocql.SerialConsistency
		tls                      *tls.Config
		tlsCertPath              string
		tlsKeyPath               string
		tlsCAPath                string
		enableHostVerification   bool
		defaultTimestamp         bool
		reconnectInterval        time.Duration
		maxWaitSchemaAgreement   time.Duration
		ignorePeerAddr           bool
		disableInitialHostLookup bool
		disableNodeStatusEvents  bool
		disableTopologyEvents    bool
		disableSchemaEvents      bool
		disableSkipMetadata      bool
		defaultIdempotence       bool
		dialer                   gocql.Dialer
		writeCoalesceWaitTime    time.Duration
		cluster                  *gocql.ClusterConfig
		session                  *gocql.Session
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "close return nil",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				session: &gocql.Session{},
			},
			want: want{},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &client{
				hosts:                    test.fields.hosts,
				cqlVersion:               test.fields.cqlVersion,
				protoVersion:             test.fields.protoVersion,
				timeout:                  test.fields.timeout,
				connectTimeout:           test.fields.connectTimeout,
				port:                     test.fields.port,
				keyspace:                 test.fields.keyspace,
				numConns:                 test.fields.numConns,
				consistency:              test.fields.consistency,
				compressor:               test.fields.compressor,
				username:                 test.fields.username,
				password:                 test.fields.password,
				authProvider:             test.fields.authProvider,
				retryPolicy:              test.fields.retryPolicy,
				reconnectionPolicy:       test.fields.reconnectionPolicy,
				poolConfig:               test.fields.poolConfig,
				hostFilter:               test.fields.hostFilter,
				socketKeepalive:          test.fields.socketKeepalive,
				maxPreparedStmts:         test.fields.maxPreparedStmts,
				maxRoutingKeyInfo:        test.fields.maxRoutingKeyInfo,
				pageSize:                 test.fields.pageSize,
				serialConsistency:        test.fields.serialConsistency,
				tls:                      test.fields.tls,
				tlsCertPath:              test.fields.tlsCertPath,
				tlsKeyPath:               test.fields.tlsKeyPath,
				tlsCAPath:                test.fields.tlsCAPath,
				enableHostVerification:   test.fields.enableHostVerification,
				defaultTimestamp:         test.fields.defaultTimestamp,
				reconnectInterval:        test.fields.reconnectInterval,
				maxWaitSchemaAgreement:   test.fields.maxWaitSchemaAgreement,
				ignorePeerAddr:           test.fields.ignorePeerAddr,
				disableInitialHostLookup: test.fields.disableInitialHostLookup,
				disableNodeStatusEvents:  test.fields.disableNodeStatusEvents,
				disableTopologyEvents:    test.fields.disableTopologyEvents,
				disableSchemaEvents:      test.fields.disableSchemaEvents,
				disableSkipMetadata:      test.fields.disableSkipMetadata,
				defaultIdempotence:       test.fields.defaultIdempotence,
				dialer:                   test.fields.dialer,
				writeCoalesceWaitTime:    test.fields.writeCoalesceWaitTime,
				cluster:                  test.fields.cluster,
				session:                  test.fields.session,
			}

			err := c.Close(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Query(t *testing.T) {
	type args struct {
		stmt  string
		names []string
	}
	type fields struct {
		hosts          []string
		cqlVersion     string
		protoVersion   int
		timeout        time.Duration
		connectTimeout time.Duration
		port           int
		keyspace       string
		numConns       int
		consistency    gocql.Consistency
		compressor     gocql.Compressor
		username       string
		password       string
		authProvider   func(h *gocql.HostInfo) (gocql.Authenticator, error)
		retryPolicy    struct {
			numRetries  int
			minDuration time.Duration
			maxDuration time.Duration
		}
		reconnectionPolicy struct {
			initialInterval time.Duration
			maxRetries      int
		}
		poolConfig struct {
			dataCenterName                 string
			enableDCAwareRouting           bool
			enableShuffleReplicas          bool
			enableNonLocalReplicasFallback bool
			enableTokenAwareHostPolicy     bool
		}
		hostFilter struct {
			enable    bool
			dcHost    string
			whiteList []string
		}
		socketKeepalive          time.Duration
		maxPreparedStmts         int
		maxRoutingKeyInfo        int
		pageSize                 int
		serialConsistency        gocql.SerialConsistency
		tls                      *tls.Config
		tlsCertPath              string
		tlsKeyPath               string
		tlsCAPath                string
		enableHostVerification   bool
		defaultTimestamp         bool
		reconnectInterval        time.Duration
		maxWaitSchemaAgreement   time.Duration
		ignorePeerAddr           bool
		disableInitialHostLookup bool
		disableNodeStatusEvents  bool
		disableTopologyEvents    bool
		disableSchemaEvents      bool
		disableSkipMetadata      bool
		defaultIdempotence       bool
		dialer                   gocql.Dialer
		writeCoalesceWaitTime    time.Duration
		cluster                  *gocql.ClusterConfig
		session                  *gocql.Session
	}
	type want struct {
		want *Queryx
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *Queryx) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *Queryx) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "query return gocqlx.Query",
			args: args{
				stmt:  "stmt",
				names: []string{"n"},
			},
			fields: fields{
				session: &gocql.Session{},
			},
			want: want{
				gocqlx.Query(new(gocql.Session).Query("stmt"), []string{"n"}),
			},
		},
		{
			name: "query return gocqlx.Query with names",
			args: args{
				stmt:  "stmt",
				names: []string{"n", "n1"},
			},
			fields: fields{
				session: &gocql.Session{},
			},
			want: want{
				gocqlx.Query(new(gocql.Session).Query("stmt"), []string{"n", "n1"}),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &client{
				hosts:                    test.fields.hosts,
				cqlVersion:               test.fields.cqlVersion,
				protoVersion:             test.fields.protoVersion,
				timeout:                  test.fields.timeout,
				connectTimeout:           test.fields.connectTimeout,
				port:                     test.fields.port,
				keyspace:                 test.fields.keyspace,
				numConns:                 test.fields.numConns,
				consistency:              test.fields.consistency,
				compressor:               test.fields.compressor,
				username:                 test.fields.username,
				password:                 test.fields.password,
				authProvider:             test.fields.authProvider,
				retryPolicy:              test.fields.retryPolicy,
				reconnectionPolicy:       test.fields.reconnectionPolicy,
				poolConfig:               test.fields.poolConfig,
				hostFilter:               test.fields.hostFilter,
				socketKeepalive:          test.fields.socketKeepalive,
				maxPreparedStmts:         test.fields.maxPreparedStmts,
				maxRoutingKeyInfo:        test.fields.maxRoutingKeyInfo,
				pageSize:                 test.fields.pageSize,
				serialConsistency:        test.fields.serialConsistency,
				tls:                      test.fields.tls,
				tlsCertPath:              test.fields.tlsCertPath,
				tlsKeyPath:               test.fields.tlsKeyPath,
				tlsCAPath:                test.fields.tlsCAPath,
				enableHostVerification:   test.fields.enableHostVerification,
				defaultTimestamp:         test.fields.defaultTimestamp,
				reconnectInterval:        test.fields.reconnectInterval,
				maxWaitSchemaAgreement:   test.fields.maxWaitSchemaAgreement,
				ignorePeerAddr:           test.fields.ignorePeerAddr,
				disableInitialHostLookup: test.fields.disableInitialHostLookup,
				disableNodeStatusEvents:  test.fields.disableNodeStatusEvents,
				disableTopologyEvents:    test.fields.disableTopologyEvents,
				disableSchemaEvents:      test.fields.disableSchemaEvents,
				disableSkipMetadata:      test.fields.disableSkipMetadata,
				defaultIdempotence:       test.fields.defaultIdempotence,
				dialer:                   test.fields.dialer,
				writeCoalesceWaitTime:    test.fields.writeCoalesceWaitTime,
				cluster:                  test.fields.cluster,
				session:                  test.fields.session,
			}

			got := c.Query(test.args.stmt, test.args.names)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	type args struct {
		table   string
		columns []string
		cmps    []Cmp
	}
	type want struct {
		wantStmt  string
		wantNames []string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string, []string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotStmt string, gotNames []string) error {
		if !reflect.DeepEqual(gotStmt, w.wantStmt) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStmt, w.wantStmt)
		}
		if !reflect.DeepEqual(gotNames, w.wantNames) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNames, w.wantNames)
		}
		return nil
	}
	tests := []test{
		func() test {
			stmt, names := qb.Select("t").Columns("col").Where(Eq("col")).ToCql()
			return test{
				name: "selete return qb.select",
				args: args{
					table:   "t",
					columns: []string{"col"},
					cmps:    []Cmp{Eq("col")},
				},
				want: want{
					wantStmt:  stmt,
					wantNames: names,
				},
			}
		}(),
		func() test {
			stmt, names := qb.Select("t").Columns("col", "col1").Where(Eq("cmp")).Where(Eq("cmp1")).ToCql()
			return test{
				name: "selete return qb.select with cols and cmps",
				args: args{
					table:   "t",
					columns: []string{"col", "col1"},
					cmps:    []Cmp{Eq("cmp"), Eq("cmp1")},
				},
				want: want{
					wantStmt:  stmt,
					wantNames: names,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotStmt, gotNames := Select(test.args.table, test.args.columns, test.args.cmps...)
			if err := checkFunc(test.want, gotStmt, gotNames); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		table string
		cmps  []Cmp
	}
	type want struct {
		want *DeleteBuilder
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *DeleteBuilder) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *DeleteBuilder) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "delete returns qb.delete",
			args: args{
				table: "t",
				cmps: []Cmp{
					Eq("col"),
				},
			},
			want: want{
				want: qb.Delete("t").Where(qb.Eq("col")),
			},
		},
		{
			name: "delete returns qb.delete with cmps",
			args: args{
				table: "t",
				cmps: []Cmp{
					Eq("col"),
					Eq("col1"),
				},
			},
			want: want{
				want: qb.Delete("t").Where(qb.Eq("col")).Where(qb.Eq("col1")),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Delete(test.args.table, test.args.cmps...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	type args struct {
		table   string
		columns []string
	}
	type want struct {
		want *InsertBuilder
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *InsertBuilder) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *InsertBuilder) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "insert return qb.insert",
			args: args{
				table:   "t",
				columns: []string{"col"},
			},
			want: want{
				want: qb.Insert("t").Columns("col"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Insert(test.args.table, test.args.columns...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		table string
	}
	type want struct {
		want *UpdateBuilder
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *UpdateBuilder) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *UpdateBuilder) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "update return qb.update",
			args: args{
				table: "t",
			},
			want: want{
				want: qb.Update("t"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Update(test.args.table)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestBatch(t *testing.T) {
	type want struct {
		want *BatchBuilder
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *BatchBuilder) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *BatchBuilder) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "bath return qb.batch",
			want: want{
				want: qb.Batch(),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Batch()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestEq(t *testing.T) {
	type args struct {
		column string
	}
	type want struct {
		want Cmp
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Cmp) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Cmp) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "eq return qb.eq",
			args: args{
				column: "col",
			},
			want: want{
				want: qb.Eq("col"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Eq(test.args.column)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIn(t *testing.T) {
	type args struct {
		column string
	}
	type want struct {
		want Cmp
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Cmp) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Cmp) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "in return qb.in",
			args: args{
				column: "col",
			},
			want: want{
				want: qb.In("col"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := In(test.args.column)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		column string
	}
	type want struct {
		want Cmp
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Cmp) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Cmp) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "contains return qb.contains",
			args: args{
				column: "col",
			},
			want: want{
				want: qb.Contains("col"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Contains(test.args.column)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWrapErrorWithKeys(t *testing.T) {
	type args struct {
		err  error
		keys []string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return error not found",
			args: args{
				err:  ErrNotFound,
				keys: []string{"k1"},
			},
			want: want{
				err: errors.ErrCassandraNotFound("k1"),
			},
		},
		{
			name: "return unavilable error",
			args: args{
				err: ErrUnavailable,
			},
			want: want{
				err: errors.ErrCassandraUnavailable(),
			},
		},
		{
			name: "return unsupported error",
			args: args{
				err: ErrUnsupported,
			},
			want: want{
				err: ErrUnsupported,
			},
		},
		{
			name: "return too many stmts error",
			args: args{
				err: ErrTooManyStmts,
			},
			want: want{
				err: ErrTooManyStmts,
			},
		},
		{
			name: "return use stmt error",
			args: args{
				err: ErrUseStmt,
			},
			want: want{
				err: ErrUseStmt,
			},
		},
		{
			name: "return session closed error",
			args: args{
				err: ErrSessionClosed,
			},
			want: want{
				err: ErrSessionClosed,
			},
		},
		{
			name: "return no connection error",
			args: args{
				err: ErrNoConnections,
			},
			want: want{
				err: ErrNoConnections,
			},
		},
		{
			name: "return no keyspace error",
			args: args{
				err: ErrNoKeyspace,
			},
			want: want{
				err: ErrNoKeyspace,
			},
		},
		{
			name: "return keyspace does not exist error",
			args: args{
				err: ErrKeyspaceDoesNotExist,
			},
			want: want{
				err: ErrKeyspaceDoesNotExist,
			},
		},
		{
			name: "return no metadata error",
			args: args{
				err: ErrNoMetadata,
			},
			want: want{
				err: ErrNoMetadata,
			},
		},
		{
			name: "return no hosts error",
			args: args{
				err: ErrNoHosts,
			},
			want: want{
				err: ErrNoHosts,
			},
		},
		{
			name: "return no connection started error",
			args: args{
				err: ErrNoConnectionsStarted,
			},
			want: want{
				err: ErrNoConnectionsStarted,
			},
		},
		{
			name: "return host query failed error",
			args: args{
				err: ErrHostQueryFailed,
			},
			want: want{
				err: ErrHostQueryFailed,
			},
		},
		{
			name: "return other error",
			args: args{
				err: errors.New("err"),
			},
			want: want{
				err: errors.New("err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			err := WrapErrorWithKeys(test.args.err, test.args.keys...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
