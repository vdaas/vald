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

package cassandra

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/vdaas/vald/internal/errors"

	"go.uber.org/goleak"
)

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
		       args: args {
		           opts: nil,
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
		           args: args {
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           hosts: nil,
		           cqlVersion: "",
		           protoVersion: 0,
		           timeout: nil,
		           connectTimeout: nil,
		           port: 0,
		           keyspace: "",
		           numConns: 0,
		           consistency: nil,
		           compressor: nil,
		           username: "",
		           password: "",
		           authProvider: nil,
		           retryPolicy: struct{numRetries int; minDuration time.Duration; maxDuration time.Duration}{},
		           reconnectionPolicy: struct{initialInterval time.Duration; maxRetries int}{},
		           poolConfig: struct{dataCenterName string; enableDCAwareRouting bool; enableShuffleReplicas bool; enableNonLocalReplicasFallback bool; enableTokenAwareHostPolicy bool}{},
		           hostFilter: struct{enable bool; dcHost string; whiteList []string}{},
		           socketKeepalive: nil,
		           maxPreparedStmts: 0,
		           maxRoutingKeyInfo: 0,
		           pageSize: 0,
		           serialConsistency: nil,
		           tls: nil,
		           tlsCertPath: "",
		           tlsKeyPath: "",
		           tlsCAPath: "",
		           enableHostVerification: false,
		           defaultTimestamp: false,
		           reconnectInterval: nil,
		           maxWaitSchemaAgreement: nil,
		           ignorePeerAddr: false,
		           disableInitialHostLookup: false,
		           disableNodeStatusEvents: false,
		           disableTopologyEvents: false,
		           disableSchemaEvents: false,
		           disableSkipMetadata: false,
		           defaultIdempotence: false,
		           dialer: nil,
		           writeCoalesceWaitTime: nil,
		           cluster: nil,
		           session: nil,
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
		           args: args {
		           ctx: nil,
		           },
		           fields: fields {
		           hosts: nil,
		           cqlVersion: "",
		           protoVersion: 0,
		           timeout: nil,
		           connectTimeout: nil,
		           port: 0,
		           keyspace: "",
		           numConns: 0,
		           consistency: nil,
		           compressor: nil,
		           username: "",
		           password: "",
		           authProvider: nil,
		           retryPolicy: struct{numRetries int; minDuration time.Duration; maxDuration time.Duration}{},
		           reconnectionPolicy: struct{initialInterval time.Duration; maxRetries int}{},
		           poolConfig: struct{dataCenterName string; enableDCAwareRouting bool; enableShuffleReplicas bool; enableNonLocalReplicasFallback bool; enableTokenAwareHostPolicy bool}{},
		           hostFilter: struct{enable bool; dcHost string; whiteList []string}{},
		           socketKeepalive: nil,
		           maxPreparedStmts: 0,
		           maxRoutingKeyInfo: 0,
		           pageSize: 0,
		           serialConsistency: nil,
		           tls: nil,
		           tlsCertPath: "",
		           tlsKeyPath: "",
		           tlsCAPath: "",
		           enableHostVerification: false,
		           defaultTimestamp: false,
		           reconnectInterval: nil,
		           maxWaitSchemaAgreement: nil,
		           ignorePeerAddr: false,
		           disableInitialHostLookup: false,
		           disableNodeStatusEvents: false,
		           disableTopologyEvents: false,
		           disableSchemaEvents: false,
		           disableSkipMetadata: false,
		           defaultIdempotence: false,
		           dialer: nil,
		           writeCoalesceWaitTime: nil,
		           cluster: nil,
		           session: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
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
			if err := test.checkFunc(test.want, err); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           hosts: nil,
		           cqlVersion: "",
		           protoVersion: 0,
		           timeout: nil,
		           connectTimeout: nil,
		           port: 0,
		           keyspace: "",
		           numConns: 0,
		           consistency: nil,
		           compressor: nil,
		           username: "",
		           password: "",
		           authProvider: nil,
		           retryPolicy: struct{numRetries int; minDuration time.Duration; maxDuration time.Duration}{},
		           reconnectionPolicy: struct{initialInterval time.Duration; maxRetries int}{},
		           poolConfig: struct{dataCenterName string; enableDCAwareRouting bool; enableShuffleReplicas bool; enableNonLocalReplicasFallback bool; enableTokenAwareHostPolicy bool}{},
		           hostFilter: struct{enable bool; dcHost string; whiteList []string}{},
		           socketKeepalive: nil,
		           maxPreparedStmts: 0,
		           maxRoutingKeyInfo: 0,
		           pageSize: 0,
		           serialConsistency: nil,
		           tls: nil,
		           tlsCertPath: "",
		           tlsKeyPath: "",
		           tlsCAPath: "",
		           enableHostVerification: false,
		           defaultTimestamp: false,
		           reconnectInterval: nil,
		           maxWaitSchemaAgreement: nil,
		           ignorePeerAddr: false,
		           disableInitialHostLookup: false,
		           disableNodeStatusEvents: false,
		           disableTopologyEvents: false,
		           disableSchemaEvents: false,
		           disableSkipMetadata: false,
		           defaultIdempotence: false,
		           dialer: nil,
		           writeCoalesceWaitTime: nil,
		           cluster: nil,
		           session: nil,
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
		           args: args {
		           ctx: nil,
		           },
		           fields: fields {
		           hosts: nil,
		           cqlVersion: "",
		           protoVersion: 0,
		           timeout: nil,
		           connectTimeout: nil,
		           port: 0,
		           keyspace: "",
		           numConns: 0,
		           consistency: nil,
		           compressor: nil,
		           username: "",
		           password: "",
		           authProvider: nil,
		           retryPolicy: struct{numRetries int; minDuration time.Duration; maxDuration time.Duration}{},
		           reconnectionPolicy: struct{initialInterval time.Duration; maxRetries int}{},
		           poolConfig: struct{dataCenterName string; enableDCAwareRouting bool; enableShuffleReplicas bool; enableNonLocalReplicasFallback bool; enableTokenAwareHostPolicy bool}{},
		           hostFilter: struct{enable bool; dcHost string; whiteList []string}{},
		           socketKeepalive: nil,
		           maxPreparedStmts: 0,
		           maxRoutingKeyInfo: 0,
		           pageSize: 0,
		           serialConsistency: nil,
		           tls: nil,
		           tlsCertPath: "",
		           tlsKeyPath: "",
		           tlsCAPath: "",
		           enableHostVerification: false,
		           defaultTimestamp: false,
		           reconnectInterval: nil,
		           maxWaitSchemaAgreement: nil,
		           ignorePeerAddr: false,
		           disableInitialHostLookup: false,
		           disableNodeStatusEvents: false,
		           disableTopologyEvents: false,
		           disableSchemaEvents: false,
		           disableSkipMetadata: false,
		           defaultIdempotence: false,
		           dialer: nil,
		           writeCoalesceWaitTime: nil,
		           cluster: nil,
		           session: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
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
			if err := test.checkFunc(test.want, err); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           stmt: "",
		           names: nil,
		       },
		       fields: fields {
		           hosts: nil,
		           cqlVersion: "",
		           protoVersion: 0,
		           timeout: nil,
		           connectTimeout: nil,
		           port: 0,
		           keyspace: "",
		           numConns: 0,
		           consistency: nil,
		           compressor: nil,
		           username: "",
		           password: "",
		           authProvider: nil,
		           retryPolicy: struct{numRetries int; minDuration time.Duration; maxDuration time.Duration}{},
		           reconnectionPolicy: struct{initialInterval time.Duration; maxRetries int}{},
		           poolConfig: struct{dataCenterName string; enableDCAwareRouting bool; enableShuffleReplicas bool; enableNonLocalReplicasFallback bool; enableTokenAwareHostPolicy bool}{},
		           hostFilter: struct{enable bool; dcHost string; whiteList []string}{},
		           socketKeepalive: nil,
		           maxPreparedStmts: 0,
		           maxRoutingKeyInfo: 0,
		           pageSize: 0,
		           serialConsistency: nil,
		           tls: nil,
		           tlsCertPath: "",
		           tlsKeyPath: "",
		           tlsCAPath: "",
		           enableHostVerification: false,
		           defaultTimestamp: false,
		           reconnectInterval: nil,
		           maxWaitSchemaAgreement: nil,
		           ignorePeerAddr: false,
		           disableInitialHostLookup: false,
		           disableNodeStatusEvents: false,
		           disableTopologyEvents: false,
		           disableSchemaEvents: false,
		           disableSkipMetadata: false,
		           defaultIdempotence: false,
		           dialer: nil,
		           writeCoalesceWaitTime: nil,
		           cluster: nil,
		           session: nil,
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
		           args: args {
		           stmt: "",
		           names: nil,
		           },
		           fields: fields {
		           hosts: nil,
		           cqlVersion: "",
		           protoVersion: 0,
		           timeout: nil,
		           connectTimeout: nil,
		           port: 0,
		           keyspace: "",
		           numConns: 0,
		           consistency: nil,
		           compressor: nil,
		           username: "",
		           password: "",
		           authProvider: nil,
		           retryPolicy: struct{numRetries int; minDuration time.Duration; maxDuration time.Duration}{},
		           reconnectionPolicy: struct{initialInterval time.Duration; maxRetries int}{},
		           poolConfig: struct{dataCenterName string; enableDCAwareRouting bool; enableShuffleReplicas bool; enableNonLocalReplicasFallback bool; enableTokenAwareHostPolicy bool}{},
		           hostFilter: struct{enable bool; dcHost string; whiteList []string}{},
		           socketKeepalive: nil,
		           maxPreparedStmts: 0,
		           maxRoutingKeyInfo: 0,
		           pageSize: 0,
		           serialConsistency: nil,
		           tls: nil,
		           tlsCertPath: "",
		           tlsKeyPath: "",
		           tlsCAPath: "",
		           enableHostVerification: false,
		           defaultTimestamp: false,
		           reconnectInterval: nil,
		           maxWaitSchemaAgreement: nil,
		           ignorePeerAddr: false,
		           disableInitialHostLookup: false,
		           disableNodeStatusEvents: false,
		           disableTopologyEvents: false,
		           disableSchemaEvents: false,
		           disableSkipMetadata: false,
		           defaultIdempotence: false,
		           dialer: nil,
		           writeCoalesceWaitTime: nil,
		           cluster: nil,
		           session: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
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
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           table: "",
		           columns: nil,
		           cmps: nil,
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
		           args: args {
		           table: "",
		           columns: nil,
		           cmps: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotStmt, gotNames := Select(test.args.table, test.args.columns, test.args.cmps...)
			if err := test.checkFunc(test.want, gotStmt, gotNames); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           table: "",
		           cmps: nil,
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
		           args: args {
		           table: "",
		           cmps: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Delete(test.args.table, test.args.cmps...)
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           table: "",
		           columns: nil,
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
		           args: args {
		           table: "",
		           columns: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Insert(test.args.table, test.args.columns...)
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           table: "",
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
		           args: args {
		           table: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Update(test.args.table)
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Batch()
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           column: "",
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
		           args: args {
		           column: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Eq(test.args.column)
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           column: "",
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
		           args: args {
		           column: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := In(test.args.column)
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           column: "",
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
		           args: args {
		           column: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Contains(test.args.column)
			if err := test.checkFunc(test.want, got); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           err: nil,
		           keys: nil,
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
		           args: args {
		           err: nil,
		           keys: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := WrapErrorWithKeys(test.args.err, test.args.keys...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
