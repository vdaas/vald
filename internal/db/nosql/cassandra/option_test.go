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

// Package redis provides implementation of Go API for redis interface
package cassandra

import (
	"crypto/tls"
	"reflect"
	"testing"

	"github.com/gocql/gocql"
)

func TestWithHosts(t *testing.T) {
	type args struct {
		hosts []string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHosts(tt.args.hosts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDialer(t *testing.T) {
	type args struct {
		der gocql.Dialer
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDialer(tt.args.der); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDialer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCQLVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCQLVersion(tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCQLVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithProtoVersion(t *testing.T) {
	type args struct {
		version int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithProtoVersion(tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithProtoVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithConnectTimeout(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConnectTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConnectTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPort(tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithKeyspace(t *testing.T) {
	type args struct {
		keyspace string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithKeyspace(tt.args.keyspace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithKeyspace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithNumConns(t *testing.T) {
	type args struct {
		numConns int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithNumConns(tt.args.numConns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNumConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithConsistency(t *testing.T) {
	type args struct {
		consistency string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConsistency(tt.args.consistency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConsistency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCompressor(t *testing.T) {
	type args struct {
		compressor gocql.Compressor
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompressor(tt.args.compressor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCompressor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithUsername(tt.args.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPassword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAuthProvider(t *testing.T) {
	type args struct {
		authProvider func(h *gocql.HostInfo) (gocql.Authenticator, error)
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAuthProvider(tt.args.authProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAuthProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetryPolicyNumRetries(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetryPolicyNumRetries(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRetryPolicyNumRetries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetryPolicyMinDuration(t *testing.T) {
	type args struct {
		minDuration string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetryPolicyMinDuration(tt.args.minDuration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRetryPolicyMinDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetryPolicyMaxDuration(t *testing.T) {
	type args struct {
		maxDuration string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetryPolicyMaxDuration(tt.args.maxDuration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRetryPolicyMaxDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithReconnectionPolicyInitialInterval(t *testing.T) {
	type args struct {
		initialInterval string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithReconnectionPolicyInitialInterval(tt.args.initialInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithReconnectionPolicyInitialInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithReconnectionPolicyMaxRetries(t *testing.T) {
	type args struct {
		maxRetries int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithReconnectionPolicyMaxRetries(tt.args.maxRetries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithReconnectionPolicyMaxRetries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSocketKeepalive(t *testing.T) {
	type args struct {
		socketKeepalive string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSocketKeepalive(tt.args.socketKeepalive); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSocketKeepalive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxPreparedStmts(t *testing.T) {
	type args struct {
		maxPreparedStmts int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMaxPreparedStmts(tt.args.maxPreparedStmts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxPreparedStmts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxRoutingKeyInfo(t *testing.T) {
	type args struct {
		maxRoutingKeyInfo int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMaxRoutingKeyInfo(tt.args.maxRoutingKeyInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxRoutingKeyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPageSize(t *testing.T) {
	type args struct {
		pageSize int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPageSize(tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPageSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSerialConsistency(t *testing.T) {
	type args struct {
		consistency gocql.SerialConsistency
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSerialConsistency(tt.args.consistency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSerialConsistency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLS(t *testing.T) {
	type args struct {
		tls *tls.Config
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTLS(tt.args.tls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLSCertPath(t *testing.T) {
	type args struct {
		certPath string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTLSCertPath(tt.args.certPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTLSCertPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLSKeyPath(t *testing.T) {
	type args struct {
		keyPath string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTLSKeyPath(tt.args.keyPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTLSKeyPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLSCAPath(t *testing.T) {
	type args struct {
		caPath string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTLSCAPath(tt.args.caPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTLSCAPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEnableHostVerification(t *testing.T) {
	type args struct {
		enableHostVerification bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEnableHostVerification(tt.args.enableHostVerification); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnableHostVerification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDefaultTimestamp(t *testing.T) {
	type args struct {
		defaultTimestamp bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDefaultTimestamp(tt.args.defaultTimestamp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDefaultTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDC(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDC(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEnableDCAwareRouting(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEnableDCAwareRouting(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnableDCAwareRouting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableDCAwareRouting(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableDCAwareRouting(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableDCAwareRouting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEnableNonLocalReplicasFallback(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEnableNonLocalReplicasFallback(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnableNonLocalReplicasFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableNonLocalReplicasFallback(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableNonLocalReplicasFallback(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableNonLocalReplicasFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEnableShuffleReplicas(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEnableShuffleReplicas(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnableShuffleReplicas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableShuffleReplicas(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableShuffleReplicas(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableShuffleReplicas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxWaitSchemaAgreement(t *testing.T) {
	type args struct {
		maxWaitSchemaAgreement string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMaxWaitSchemaAgreement(tt.args.maxWaitSchemaAgreement); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxWaitSchemaAgreement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithReconnectInterval(t *testing.T) {
	type args struct {
		reconnectInterval string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithReconnectInterval(tt.args.reconnectInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithReconnectInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIgnorePeerAddr(t *testing.T) {
	type args struct {
		ignorePeerAddr bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithIgnorePeerAddr(tt.args.ignorePeerAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIgnorePeerAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableInitialHostLookup(t *testing.T) {
	type args struct {
		disableInitialHostLookup bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableInitialHostLookup(tt.args.disableInitialHostLookup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableInitialHostLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableNodeStatusEvents(t *testing.T) {
	type args struct {
		disableNodeStatusEvents bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableNodeStatusEvents(tt.args.disableNodeStatusEvents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableNodeStatusEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableTopologyEvents(t *testing.T) {
	type args struct {
		disableTopologyEvents bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableTopologyEvents(tt.args.disableTopologyEvents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableTopologyEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableSchemaEvents(t *testing.T) {
	type args struct {
		disableSchemaEvents bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableSchemaEvents(tt.args.disableSchemaEvents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableSchemaEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableSkipMetadata(t *testing.T) {
	type args struct {
		disableSkipMetadata bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableSkipMetadata(tt.args.disableSkipMetadata); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableSkipMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDefaultIdempotence(t *testing.T) {
	type args struct {
		defaultIdempotence bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDefaultIdempotence(tt.args.defaultIdempotence); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDefaultIdempotence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithWriteCoalesceWaitTime(t *testing.T) {
	type args struct {
		writeCoalesceWaitTime string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithWriteCoalesceWaitTime(tt.args.writeCoalesceWaitTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithWriteCoalesceWaitTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
