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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestCassandra_Bind(t *testing.T) {
	t.Parallel()
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
		TCP                      *TCP
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Cassandra) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Hosts: nil,
		           CQLVersion: "",
		           ProtoVersion: 0,
		           Timeout: "",
		           ConnectTimeout: "",
		           Port: 0,
		           Keyspace: "",
		           NumConns: 0,
		           Consistency: "",
		           SerialConsistency: "",
		           Username: "",
		           Password: "",
		           PoolConfig: PoolConfig{},
		           RetryPolicy: RetryPolicy{},
		           ReconnectionPolicy: ReconnectionPolicy{},
		           HostFilter: HostFilter{},
		           SocketKeepalive: "",
		           MaxPreparedStmts: 0,
		           MaxRoutingKeyInfo: 0,
		           PageSize: 0,
		           TLS: TLS{},
		           TCP: TCP{},
		           EnableHostVerification: false,
		           DefaultTimestamp: false,
		           ReconnectInterval: "",
		           MaxWaitSchemaAgreement: "",
		           IgnorePeerAddr: false,
		           DisableInitialHostLookup: false,
		           DisableNodeStatusEvents: false,
		           DisableTopologyEvents: false,
		           DisableSchemaEvents: false,
		           DisableSkipMetadata: false,
		           DefaultIdempotence: false,
		           WriteCoalesceWaitTime: "",
		           KVTable: "",
		           VKTable: "",
		           VectorBackupTable: "",
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
		           Hosts: nil,
		           CQLVersion: "",
		           ProtoVersion: 0,
		           Timeout: "",
		           ConnectTimeout: "",
		           Port: 0,
		           Keyspace: "",
		           NumConns: 0,
		           Consistency: "",
		           SerialConsistency: "",
		           Username: "",
		           Password: "",
		           PoolConfig: PoolConfig{},
		           RetryPolicy: RetryPolicy{},
		           ReconnectionPolicy: ReconnectionPolicy{},
		           HostFilter: HostFilter{},
		           SocketKeepalive: "",
		           MaxPreparedStmts: 0,
		           MaxRoutingKeyInfo: 0,
		           PageSize: 0,
		           TLS: TLS{},
		           TCP: TCP{},
		           EnableHostVerification: false,
		           DefaultTimestamp: false,
		           ReconnectInterval: "",
		           MaxWaitSchemaAgreement: "",
		           IgnorePeerAddr: false,
		           DisableInitialHostLookup: false,
		           DisableNodeStatusEvents: false,
		           DisableTopologyEvents: false,
		           DisableSchemaEvents: false,
		           DisableSkipMetadata: false,
		           DefaultIdempotence: false,
		           WriteCoalesceWaitTime: "",
		           KVTable: "",
		           VKTable: "",
		           VectorBackupTable: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
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
				TCP:                      test.fields.TCP,
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
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCassandra_Opts(t *testing.T) {
	t.Parallel()
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
		TCP                      *TCP
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotOpts []cassandra.Option, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotOpts, w.wantOpts) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOpts, w.wantOpts)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Hosts: nil,
		           CQLVersion: "",
		           ProtoVersion: 0,
		           Timeout: "",
		           ConnectTimeout: "",
		           Port: 0,
		           Keyspace: "",
		           NumConns: 0,
		           Consistency: "",
		           SerialConsistency: "",
		           Username: "",
		           Password: "",
		           PoolConfig: PoolConfig{},
		           RetryPolicy: RetryPolicy{},
		           ReconnectionPolicy: ReconnectionPolicy{},
		           HostFilter: HostFilter{},
		           SocketKeepalive: "",
		           MaxPreparedStmts: 0,
		           MaxRoutingKeyInfo: 0,
		           PageSize: 0,
		           TLS: TLS{},
		           TCP: TCP{},
		           EnableHostVerification: false,
		           DefaultTimestamp: false,
		           ReconnectInterval: "",
		           MaxWaitSchemaAgreement: "",
		           IgnorePeerAddr: false,
		           DisableInitialHostLookup: false,
		           DisableNodeStatusEvents: false,
		           DisableTopologyEvents: false,
		           DisableSchemaEvents: false,
		           DisableSkipMetadata: false,
		           DefaultIdempotence: false,
		           WriteCoalesceWaitTime: "",
		           KVTable: "",
		           VKTable: "",
		           VectorBackupTable: "",
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
		           Hosts: nil,
		           CQLVersion: "",
		           ProtoVersion: 0,
		           Timeout: "",
		           ConnectTimeout: "",
		           Port: 0,
		           Keyspace: "",
		           NumConns: 0,
		           Consistency: "",
		           SerialConsistency: "",
		           Username: "",
		           Password: "",
		           PoolConfig: PoolConfig{},
		           RetryPolicy: RetryPolicy{},
		           ReconnectionPolicy: ReconnectionPolicy{},
		           HostFilter: HostFilter{},
		           SocketKeepalive: "",
		           MaxPreparedStmts: 0,
		           MaxRoutingKeyInfo: 0,
		           PageSize: 0,
		           TLS: TLS{},
		           TCP: TCP{},
		           EnableHostVerification: false,
		           DefaultTimestamp: false,
		           ReconnectInterval: "",
		           MaxWaitSchemaAgreement: "",
		           IgnorePeerAddr: false,
		           DisableInitialHostLookup: false,
		           DisableNodeStatusEvents: false,
		           DisableTopologyEvents: false,
		           DisableSchemaEvents: false,
		           DisableSkipMetadata: false,
		           DefaultIdempotence: false,
		           WriteCoalesceWaitTime: "",
		           KVTable: "",
		           VKTable: "",
		           VectorBackupTable: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
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
				TCP:                      test.fields.TCP,
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
			if err := test.checkFunc(test.want, gotOpts, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
