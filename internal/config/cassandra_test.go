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
	"time"

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	testdata "github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestCassandra_Bind(t *testing.T) {
	type fields struct {
		Hosts                    []string
		CQLVersion               string
		ProtoVersion             int
		Timeout                  string
		ConnectTimeout           string
		Port                     int
		Keyspace                 string
		NumConns                 int
		Consistency              string
		SerialConsistency        string
		Username                 string
		Password                 string
		PoolConfig               *PoolConfig
		RetryPolicy              *RetryPolicy
		ReconnectionPolicy       *ReconnectionPolicy
		HostFilter               *HostFilter
		SocketKeepalive          string
		MaxPreparedStmts         int
		MaxRoutingKeyInfo        int
		PageSize                 int
		TLS                      *TLS
		Net                      *Net
		EnableHostVerification   bool
		DefaultTimestamp         bool
		ReconnectInterval        string
		MaxWaitSchemaAgreement   string
		IgnorePeerAddr           bool
		DisableInitialHostLookup bool
		DisableNodeStatusEvents  bool
		DisableTopologyEvents    bool
		DisableSchemaEvents      bool
		DisableSkipMetadata      bool
		DefaultIdempotence       bool
		WriteCoalesceWaitTime    string
		KVTable                  string
		VKTable                  string
		VectorBackupTable        string
	}
	type want struct {
		want *Cassandra
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Cassandra) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Cassandra) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return Cassandra that is variable set when parameters are not empty",
				fields: fields{
					Hosts: []string{
						"cassandra-0.cassandra.default.svc.cluster.local",
						"cassandra-1.cassandra.default.svc.cluster.local",
						"cassandra-2.cassandra.default.svc.cluster.local",
					},
					CQLVersion:        "3.0.0",
					ProtoVersion:      0,
					Timeout:           "600ms",
					ConnectTimeout:    "3s",
					Port:              9042,
					Keyspace:          "vald",
					NumConns:          2,
					Consistency:       "quorum",
					SerialConsistency: "localserial",
					Username:          "root",
					Password:          "password",
					PoolConfig: &PoolConfig{
						DataCenter:               "",
						DCAwareRouting:           false,
						NonLocalReplicasFallback: false,
						ShuffleReplicas:          false,
						TokenAwareHostPolicy:     false,
					},
					RetryPolicy: &RetryPolicy{
						NumRetries:  3,
						MinDuration: "10ms",
						MaxDuration: "1s",
					},
					ReconnectionPolicy: &ReconnectionPolicy{
						MaxRetries:      3,
						InitialInterval: "100ms",
					},
					HostFilter: &HostFilter{
						Enabled:    false,
						DataCenter: "",
						WhiteList:  []string{},
					},
					SocketKeepalive:   "0s",
					MaxPreparedStmts:  1000,
					MaxRoutingKeyInfo: 1000,
					PageSize:          5000,
					TLS: &TLS{
						Enabled:            false,
						Cert:               "/path/ro/cert",
						Key:                "/path/to/key",
						CA:                 "/path/to/ca",
						InsecureSkipVerify: false,
					},
					Net: &Net{
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "5m",
							CacheExpiration: "24h",
						},
						Dialer: &Dialer{
							Timeout:          "30s",
							Keepalive:        "10m",
							DualStackEnabled: false,
						},
						TLS: &TLS{
							Enabled:            false,
							Cert:               "/path/ro/cert",
							Key:                "/path/to/key",
							CA:                 "/path/to/ca",
							InsecureSkipVerify: false,
						},
						SocketOption: &SocketOption{
							ReusePort:                true,
							ReuseAddr:                true,
							TCPFastOpen:              true,
							TCPCork:                  false,
							TCPDeferAccept:           true,
							IPTransparent:            false,
							IPRecoverDestinationAddr: false,
						},
					},
					EnableHostVerification:   false,
					DefaultTimestamp:         true,
					ReconnectInterval:        "",
					MaxWaitSchemaAgreement:   "",
					IgnorePeerAddr:           false,
					DisableInitialHostLookup: false,
					DisableNodeStatusEvents:  false,
					DisableTopologyEvents:    false,
					DisableSchemaEvents:      false,
					DisableSkipMetadata:      false,
					DefaultIdempotence:       false,
					WriteCoalesceWaitTime:    "200ms",
					KVTable:                  "kv",
					VKTable:                  "vk",
					VectorBackupTable:        "backup_vector",
				},
				want: want{
					want: &Cassandra{
						Hosts: []string{
							"cassandra-0.cassandra.default.svc.cluster.local",
							"cassandra-1.cassandra.default.svc.cluster.local",
							"cassandra-2.cassandra.default.svc.cluster.local",
						},
						CQLVersion:        "3.0.0",
						ProtoVersion:      0,
						Timeout:           "600ms",
						ConnectTimeout:    "3s",
						Port:              9042,
						Keyspace:          "vald",
						NumConns:          2,
						Consistency:       "quorum",
						SerialConsistency: "localserial",
						Username:          "root",
						Password:          "password",
						PoolConfig: &PoolConfig{
							DataCenter:               "",
							DCAwareRouting:           false,
							NonLocalReplicasFallback: false,
							ShuffleReplicas:          false,
							TokenAwareHostPolicy:     false,
						},
						RetryPolicy: &RetryPolicy{
							NumRetries:  3,
							MinDuration: "10ms",
							MaxDuration: "1s",
						},
						ReconnectionPolicy: &ReconnectionPolicy{
							MaxRetries:      3,
							InitialInterval: "100ms",
						},
						HostFilter: &HostFilter{
							Enabled:    false,
							DataCenter: "",
							WhiteList:  []string{},
						},
						SocketKeepalive:   "0s",
						MaxPreparedStmts:  1000,
						MaxRoutingKeyInfo: 1000,
						PageSize:          5000,
						TLS: &TLS{
							Enabled:            false,
							Cert:               "/path/ro/cert",
							Key:                "/path/to/key",
							CA:                 "/path/to/ca",
							InsecureSkipVerify: false,
						},
						Net: &Net{
							DNS: &DNS{
								CacheEnabled:    true,
								RefreshDuration: "5m",
								CacheExpiration: "24h",
							},
							Dialer: &Dialer{
								Timeout:          "30s",
								Keepalive:        "10m",
								DualStackEnabled: false,
							},
							TLS: &TLS{
								Enabled:            false,
								Cert:               "/path/ro/cert",
								Key:                "/path/to/key",
								CA:                 "/path/to/ca",
								InsecureSkipVerify: false,
							},
							SocketOption: &SocketOption{
								ReusePort:                true,
								ReuseAddr:                true,
								TCPFastOpen:              true,
								TCPCork:                  false,
								TCPDeferAccept:           true,
								IPTransparent:            false,
								IPRecoverDestinationAddr: false,
							},
						},
						EnableHostVerification:   false,
						DefaultTimestamp:         true,
						ReconnectInterval:        "",
						MaxWaitSchemaAgreement:   "",
						IgnorePeerAddr:           false,
						DisableInitialHostLookup: false,
						DisableNodeStatusEvents:  false,
						DisableTopologyEvents:    false,
						DisableSchemaEvents:      false,
						DisableSkipMetadata:      false,
						DefaultIdempotence:       false,
						WriteCoalesceWaitTime:    "200ms",
						KVTable:                  "kv",
						VKTable:                  "vk",
						VectorBackupTable:        "backup_vector",
					},
				},
			}
		}(),
		func() test {
			key := "CASSANDRA_BIND_PASSWORD"
			val := "cassandra_password"
			return test{
				name: "return Cassandra struct when Password is set via the envirionment value",
				fields: fields{
					Password: "_" + key + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					t.Setenv(key, val)
				},
				want: want{
					want: &Cassandra{
						Password: val,
						TLS:      new(TLS),
						Net:      new(Net),
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Cassandra that is the default variables set when all parameters are empty or nil",
				fields: fields{},
				want: want{
					want: &Cassandra{
						TLS: new(TLS),
						Net: new(Net),
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
			c := &Cassandra{
				Hosts:                    test.fields.Hosts,
				CQLVersion:               test.fields.CQLVersion,
				ProtoVersion:             test.fields.ProtoVersion,
				Timeout:                  test.fields.Timeout,
				ConnectTimeout:           test.fields.ConnectTimeout,
				Port:                     test.fields.Port,
				Keyspace:                 test.fields.Keyspace,
				NumConns:                 test.fields.NumConns,
				Consistency:              test.fields.Consistency,
				SerialConsistency:        test.fields.SerialConsistency,
				Username:                 test.fields.Username,
				Password:                 test.fields.Password,
				PoolConfig:               test.fields.PoolConfig,
				RetryPolicy:              test.fields.RetryPolicy,
				ReconnectionPolicy:       test.fields.ReconnectionPolicy,
				HostFilter:               test.fields.HostFilter,
				SocketKeepalive:          test.fields.SocketKeepalive,
				MaxPreparedStmts:         test.fields.MaxPreparedStmts,
				MaxRoutingKeyInfo:        test.fields.MaxRoutingKeyInfo,
				PageSize:                 test.fields.PageSize,
				TLS:                      test.fields.TLS,
				Net:                      test.fields.Net,
				EnableHostVerification:   test.fields.EnableHostVerification,
				DefaultTimestamp:         test.fields.DefaultTimestamp,
				ReconnectInterval:        test.fields.ReconnectInterval,
				MaxWaitSchemaAgreement:   test.fields.MaxWaitSchemaAgreement,
				IgnorePeerAddr:           test.fields.IgnorePeerAddr,
				DisableInitialHostLookup: test.fields.DisableInitialHostLookup,
				DisableNodeStatusEvents:  test.fields.DisableNodeStatusEvents,
				DisableTopologyEvents:    test.fields.DisableTopologyEvents,
				DisableSchemaEvents:      test.fields.DisableSchemaEvents,
				DisableSkipMetadata:      test.fields.DisableSkipMetadata,
				DefaultIdempotence:       test.fields.DefaultIdempotence,
				WriteCoalesceWaitTime:    test.fields.WriteCoalesceWaitTime,
				KVTable:                  test.fields.KVTable,
				VKTable:                  test.fields.VKTable,
				VectorBackupTable:        test.fields.VectorBackupTable,
			}

			got := c.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCassandra_Opts(t *testing.T) {
	type fields struct {
		Hosts                    []string
		CQLVersion               string
		ProtoVersion             int
		Timeout                  string
		ConnectTimeout           string
		Port                     int
		Keyspace                 string
		NumConns                 int
		Consistency              string
		SerialConsistency        string
		Username                 string
		Password                 string
		PoolConfig               *PoolConfig
		RetryPolicy              *RetryPolicy
		ReconnectionPolicy       *ReconnectionPolicy
		HostFilter               *HostFilter
		SocketKeepalive          string
		MaxPreparedStmts         int
		MaxRoutingKeyInfo        int
		PageSize                 int
		TLS                      *TLS
		Net                      *Net
		EnableHostVerification   bool
		DefaultTimestamp         bool
		ReconnectInterval        string
		MaxWaitSchemaAgreement   string
		IgnorePeerAddr           bool
		DisableInitialHostLookup bool
		DisableNodeStatusEvents  bool
		DisableTopologyEvents    bool
		DisableSchemaEvents      bool
		DisableSkipMetadata      bool
		DefaultIdempotence       bool
		WriteCoalesceWaitTime    string
		KVTable                  string
		VKTable                  string
		VectorBackupTable        string
	}
	type want struct {
		wantOpts []cassandra.Option
		err      error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []cassandra.Option, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, gotOpts []cassandra.Option, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(len(gotOpts), len(w.wantOpts)) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOpts, w.wantOpts)
		}
		return nil
	}
	tests := []test{
		func() test {
			cert := testdata.GetTestdataPath("tls/dummyServer.crt")
			key := testdata.GetTestdataPath("tls/dummyServer.key")
			ca := testdata.GetTestdataPath("tls/dummyCa.pem")
			return test{
				name: "return 45 cassandra.Option when no error occurred",
				fields: fields{
					Hosts: []string{
						"cassandra-0.cassandra.default.svc.cluster.local",
						"cassandra-1.cassandra.default.svc.cluster.local",
						"cassandra-2.cassandra.default.svc.cluster.local",
					},
					CQLVersion:        "3.0.0",
					ProtoVersion:      0,
					Timeout:           "600ms",
					ConnectTimeout:    "3s",
					Port:              9042,
					Keyspace:          "vald",
					NumConns:          2,
					Consistency:       "quorum",
					SerialConsistency: "localserial",
					Username:          "root",
					Password:          "password",
					PoolConfig: &PoolConfig{
						DataCenter:               "",
						DCAwareRouting:           false,
						NonLocalReplicasFallback: false,
						ShuffleReplicas:          false,
						TokenAwareHostPolicy:     false,
					},
					RetryPolicy: &RetryPolicy{
						NumRetries:  3,
						MinDuration: "10ms",
						MaxDuration: "1s",
					},
					ReconnectionPolicy: &ReconnectionPolicy{
						MaxRetries:      3,
						InitialInterval: "100ms",
					},
					HostFilter: &HostFilter{
						Enabled:    false,
						DataCenter: "",
						WhiteList:  []string{},
					},
					SocketKeepalive:   "0s",
					MaxPreparedStmts:  1000,
					MaxRoutingKeyInfo: 1000,
					PageSize:          5000,
					TLS: &TLS{
						Enabled:            true,
						Cert:               cert,
						Key:                key,
						CA:                 ca,
						InsecureSkipVerify: false,
					},
					Net: &Net{
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "5m",
							CacheExpiration: "24h",
						},
						Dialer: &Dialer{
							Timeout:          "30s",
							Keepalive:        "10m",
							DualStackEnabled: false,
						},
						TLS: &TLS{
							Enabled:            false,
							Cert:               cert,
							Key:                key,
							CA:                 ca,
							InsecureSkipVerify: false,
						},
						SocketOption: &SocketOption{
							ReusePort:                true,
							ReuseAddr:                true,
							TCPFastOpen:              true,
							TCPCork:                  false,
							TCPDeferAccept:           true,
							IPTransparent:            false,
							IPRecoverDestinationAddr: false,
						},
					},
					EnableHostVerification:   false,
					DefaultTimestamp:         true,
					ReconnectInterval:        "",
					MaxWaitSchemaAgreement:   "",
					IgnorePeerAddr:           false,
					DisableInitialHostLookup: false,
					DisableNodeStatusEvents:  false,
					DisableTopologyEvents:    false,
					DisableSchemaEvents:      false,
					DisableSkipMetadata:      false,
					DefaultIdempotence:       false,
					WriteCoalesceWaitTime:    "200ms",
					KVTable:                  "kv",
					VKTable:                  "vk",
					VectorBackupTable:        "backup_vector",
				},
				want: want{
					wantOpts: make([]cassandra.Option, 45),
					err:      nil,
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil and error when TLS config value in the invalid value",
				fields: fields{
					Hosts: []string{
						"cassandra-0.cassandra.default.svc.cluster.local",
						"cassandra-1.cassandra.default.svc.cluster.local",
						"cassandra-2.cassandra.default.svc.cluster.local",
					},
					CQLVersion:        "3.0.0",
					ProtoVersion:      0,
					Timeout:           "600ms",
					ConnectTimeout:    "3s",
					Port:              9042,
					Keyspace:          "vald",
					NumConns:          2,
					Consistency:       "quorum",
					SerialConsistency: "localserial",
					Username:          "root",
					Password:          "password",
					PoolConfig: &PoolConfig{
						DataCenter:               "",
						DCAwareRouting:           false,
						NonLocalReplicasFallback: false,
						ShuffleReplicas:          false,
						TokenAwareHostPolicy:     false,
					},
					RetryPolicy: &RetryPolicy{
						NumRetries:  3,
						MinDuration: "10ms",
						MaxDuration: "1s",
					},
					ReconnectionPolicy: &ReconnectionPolicy{
						MaxRetries:      3,
						InitialInterval: "100ms",
					},
					HostFilter: &HostFilter{
						Enabled:    false,
						DataCenter: "",
						WhiteList:  []string{},
					},
					SocketKeepalive:   "0s",
					MaxPreparedStmts:  1000,
					MaxRoutingKeyInfo: 1000,
					PageSize:          5000,
					TLS: &TLS{
						Enabled:            true,
						Cert:               "",
						Key:                "",
						InsecureSkipVerify: false,
					},
					Net: &Net{
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "5m",
							CacheExpiration: "24h",
						},
						Dialer: &Dialer{
							Timeout:          "30s",
							Keepalive:        "10m",
							DualStackEnabled: false,
						},
						TLS: &TLS{
							Enabled:            false,
							Cert:               "/path/ro/cert",
							Key:                "/path/to/key",
							CA:                 "/path/to/ca",
							InsecureSkipVerify: false,
						},
						SocketOption: &SocketOption{
							ReusePort:                true,
							ReuseAddr:                true,
							TCPFastOpen:              true,
							TCPCork:                  false,
							TCPDeferAccept:           true,
							IPTransparent:            false,
							IPRecoverDestinationAddr: false,
						},
					},
					EnableHostVerification:   false,
					DefaultTimestamp:         true,
					ReconnectInterval:        "",
					MaxWaitSchemaAgreement:   "",
					IgnorePeerAddr:           false,
					DisableInitialHostLookup: false,
					DisableNodeStatusEvents:  false,
					DisableTopologyEvents:    false,
					DisableSchemaEvents:      false,
					DisableSkipMetadata:      false,
					DefaultIdempotence:       false,
					WriteCoalesceWaitTime:    "200ms",
					KVTable:                  "kv",
					VKTable:                  "vk",
					VectorBackupTable:        "backup_vector",
				},
				want: want{
					wantOpts: nil,
					err:      errors.ErrTLSCertOrKeyNotFound,
				},
			}
		}(),
		func() test {
			cert := testdata.GetTestdataPath("tls/dummyServer.crt")
			key := testdata.GetTestdataPath("tls/dummyServer.key")
			ca := testdata.GetTestdataPath("tls/dummyCa.pem")
			return test{
				name: "return nil and err when net.NewDialer returns error",
				fields: fields{
					Hosts: []string{
						"cassandra-0.cassandra.default.svc.cluster.local",
						"cassandra-1.cassandra.default.svc.cluster.local",
						"cassandra-2.cassandra.default.svc.cluster.local",
					},
					CQLVersion:        "3.0.0",
					ProtoVersion:      0,
					Timeout:           "600ms",
					ConnectTimeout:    "3s",
					Port:              9042,
					Keyspace:          "vald",
					NumConns:          2,
					Consistency:       "quorum",
					SerialConsistency: "localserial",
					Username:          "root",
					Password:          "password",
					PoolConfig: &PoolConfig{
						DataCenter:               "",
						DCAwareRouting:           false,
						NonLocalReplicasFallback: false,
						ShuffleReplicas:          false,
						TokenAwareHostPolicy:     false,
					},
					RetryPolicy: &RetryPolicy{
						NumRetries:  3,
						MinDuration: "10ms",
						MaxDuration: "1s",
					},
					ReconnectionPolicy: &ReconnectionPolicy{
						MaxRetries:      3,
						InitialInterval: "100ms",
					},
					HostFilter: &HostFilter{
						Enabled:    false,
						DataCenter: "",
						WhiteList:  []string{},
					},
					SocketKeepalive:   "0s",
					MaxPreparedStmts:  1000,
					MaxRoutingKeyInfo: 1000,
					PageSize:          5000,
					TLS: &TLS{
						Enabled:            false,
						Cert:               cert,
						Key:                key,
						CA:                 ca,
						InsecureSkipVerify: false,
					},
					Net: &Net{
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "5m",
							CacheExpiration: "1m",
						},
						Dialer: &Dialer{
							Timeout:          "30s",
							Keepalive:        "10m",
							DualStackEnabled: false,
						},
						TLS: &TLS{
							Enabled:            false,
							Cert:               cert,
							Key:                key,
							CA:                 ca,
							InsecureSkipVerify: false,
						},
						SocketOption: &SocketOption{
							ReusePort:                true,
							ReuseAddr:                true,
							TCPFastOpen:              true,
							TCPCork:                  false,
							TCPDeferAccept:           true,
							IPTransparent:            false,
							IPRecoverDestinationAddr: false,
						},
					},
					EnableHostVerification:   false,
					DefaultTimestamp:         true,
					ReconnectInterval:        "",
					MaxWaitSchemaAgreement:   "",
					IgnorePeerAddr:           false,
					DisableInitialHostLookup: false,
					DisableNodeStatusEvents:  false,
					DisableTopologyEvents:    false,
					DisableSchemaEvents:      false,
					DisableSkipMetadata:      false,
					DefaultIdempotence:       false,
					WriteCoalesceWaitTime:    "200ms",
					KVTable:                  "kv",
					VKTable:                  "vk",
					VectorBackupTable:        "backup_vector",
				},
				want: want{
					wantOpts: nil,
					err:      errors.ErrInvalidDNSConfig(5*time.Minute, 1*time.Minute),
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil and err when net.Net.Opts returns error",
				fields: fields{
					Hosts: []string{
						"cassandra-0.cassandra.default.svc.cluster.local",
						"cassandra-1.cassandra.default.svc.cluster.local",
						"cassandra-2.cassandra.default.svc.cluster.local",
					},
					CQLVersion:        "3.0.0",
					ProtoVersion:      0,
					Timeout:           "600ms",
					ConnectTimeout:    "3s",
					Port:              9042,
					Keyspace:          "vald",
					NumConns:          2,
					Consistency:       "quorum",
					SerialConsistency: "localserial",
					Username:          "root",
					Password:          "password",
					PoolConfig: &PoolConfig{
						DataCenter:               "",
						DCAwareRouting:           false,
						NonLocalReplicasFallback: false,
						ShuffleReplicas:          false,
						TokenAwareHostPolicy:     false,
					},
					RetryPolicy: &RetryPolicy{
						NumRetries:  3,
						MinDuration: "10ms",
						MaxDuration: "1s",
					},
					ReconnectionPolicy: &ReconnectionPolicy{
						MaxRetries:      3,
						InitialInterval: "100ms",
					},
					HostFilter: &HostFilter{
						Enabled:    false,
						DataCenter: "",
						WhiteList:  []string{},
					},
					SocketKeepalive:   "0s",
					MaxPreparedStmts:  1000,
					MaxRoutingKeyInfo: 1000,
					PageSize:          5000,
					TLS: &TLS{
						Enabled: false,
					},
					Net: &Net{
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "1m",
							CacheExpiration: "5m",
						},
						Dialer: &Dialer{
							Timeout:          "30s",
							Keepalive:        "10m",
							DualStackEnabled: false,
						},
						TLS: &TLS{
							Enabled: true,
						},
						SocketOption: &SocketOption{
							ReusePort:                true,
							ReuseAddr:                true,
							TCPFastOpen:              true,
							TCPCork:                  false,
							TCPDeferAccept:           true,
							IPTransparent:            false,
							IPRecoverDestinationAddr: false,
						},
					},
					EnableHostVerification:   false,
					DefaultTimestamp:         true,
					ReconnectInterval:        "",
					MaxWaitSchemaAgreement:   "",
					IgnorePeerAddr:           false,
					DisableInitialHostLookup: false,
					DisableNodeStatusEvents:  false,
					DisableTopologyEvents:    false,
					DisableSchemaEvents:      false,
					DisableSkipMetadata:      false,
					DefaultIdempotence:       false,
					WriteCoalesceWaitTime:    "200ms",
					KVTable:                  "kv",
					VKTable:                  "vk",
					VectorBackupTable:        "backup_vector",
				},
				want: want{
					wantOpts: nil,
					err:      errors.ErrTLSCertOrKeyNotFound,
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
			cfg := &Cassandra{
				Hosts:                    test.fields.Hosts,
				CQLVersion:               test.fields.CQLVersion,
				ProtoVersion:             test.fields.ProtoVersion,
				Timeout:                  test.fields.Timeout,
				ConnectTimeout:           test.fields.ConnectTimeout,
				Port:                     test.fields.Port,
				Keyspace:                 test.fields.Keyspace,
				NumConns:                 test.fields.NumConns,
				Consistency:              test.fields.Consistency,
				SerialConsistency:        test.fields.SerialConsistency,
				Username:                 test.fields.Username,
				Password:                 test.fields.Password,
				PoolConfig:               test.fields.PoolConfig,
				RetryPolicy:              test.fields.RetryPolicy,
				ReconnectionPolicy:       test.fields.ReconnectionPolicy,
				HostFilter:               test.fields.HostFilter,
				SocketKeepalive:          test.fields.SocketKeepalive,
				MaxPreparedStmts:         test.fields.MaxPreparedStmts,
				MaxRoutingKeyInfo:        test.fields.MaxRoutingKeyInfo,
				PageSize:                 test.fields.PageSize,
				TLS:                      test.fields.TLS,
				Net:                      test.fields.Net,
				EnableHostVerification:   test.fields.EnableHostVerification,
				DefaultTimestamp:         test.fields.DefaultTimestamp,
				ReconnectInterval:        test.fields.ReconnectInterval,
				MaxWaitSchemaAgreement:   test.fields.MaxWaitSchemaAgreement,
				IgnorePeerAddr:           test.fields.IgnorePeerAddr,
				DisableInitialHostLookup: test.fields.DisableInitialHostLookup,
				DisableNodeStatusEvents:  test.fields.DisableNodeStatusEvents,
				DisableTopologyEvents:    test.fields.DisableTopologyEvents,
				DisableSchemaEvents:      test.fields.DisableSchemaEvents,
				DisableSkipMetadata:      test.fields.DisableSkipMetadata,
				DefaultIdempotence:       test.fields.DefaultIdempotence,
				WriteCoalesceWaitTime:    test.fields.WriteCoalesceWaitTime,
				KVTable:                  test.fields.KVTable,
				VKTable:                  test.fields.VKTable,
				VectorBackupTable:        test.fields.VectorBackupTable,
			}

			gotOpts, err := cfg.Opts()
			if err := checkFunc(test.want, gotOpts, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
