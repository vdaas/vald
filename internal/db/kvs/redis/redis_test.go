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

package redis

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []Option
	}
	type want struct {
		wantRc Redis
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Redis, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRc Redis, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRc, w.wantRc) {
			return errors.Errorf("got = %v, want %v", gotRc, w.wantRc)
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
		           ctx: nil,
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

			gotRc, err := New(test.args.ctx, test.args.opts...)
			if err := test.checkFunc(test.want, gotRc, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_redisClient_ping(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		addrs                []string
		clusterSlots         func() ([]redis.ClusterSlot, error)
		db                   int
		dialTimeout          time.Duration
		dialer               func(ctx context.Context, network, addr string) (net.Conn, error)
		idleCheckFrequency   time.Duration
		idleTimeout          time.Duration
		initialPingDuration  time.Duration
		initialPingTimeLimit time.Duration
		keyPref              string
		maxConnAge           time.Duration
		maxRedirects         int
		maxRetries           int
		maxRetryBackoff      time.Duration
		minIdleConns         int
		minRetryBackoff      time.Duration
		onConnect            func(*redis.Conn) error
		onNewNode            func(*redis.Client)
		password             string
		poolSize             int
		poolTimeout          time.Duration
		readOnly             bool
		readTimeout          time.Duration
		routeByLatency       bool
		routeRandomly        bool
		tlsConfig            *tls.Config
		writeTimeout         time.Duration
		client               Redis
		pingEnabled          bool
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           addrs: nil,
		           clusterSlots: nil,
		           db: 0,
		           dialTimeout: nil,
		           dialer: nil,
		           idleCheckFrequency: nil,
		           idleTimeout: nil,
		           initialPingDuration: nil,
		           initialPingTimeLimit: nil,
		           keyPref: "",
		           maxConnAge: nil,
		           maxRedirects: 0,
		           maxRetries: 0,
		           maxRetryBackoff: nil,
		           minIdleConns: 0,
		           minRetryBackoff: nil,
		           onConnect: nil,
		           onNewNode: nil,
		           password: "",
		           poolSize: 0,
		           poolTimeout: nil,
		           readOnly: false,
		           readTimeout: nil,
		           routeByLatency: false,
		           routeRandomly: false,
		           tlsConfig: nil,
		           writeTimeout: nil,
		           client: nil,
		           pingEnabled: false,
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
		           addrs: nil,
		           clusterSlots: nil,
		           db: 0,
		           dialTimeout: nil,
		           dialer: nil,
		           idleCheckFrequency: nil,
		           idleTimeout: nil,
		           initialPingDuration: nil,
		           initialPingTimeLimit: nil,
		           keyPref: "",
		           maxConnAge: nil,
		           maxRedirects: 0,
		           maxRetries: 0,
		           maxRetryBackoff: nil,
		           minIdleConns: 0,
		           minRetryBackoff: nil,
		           onConnect: nil,
		           onNewNode: nil,
		           password: "",
		           poolSize: 0,
		           poolTimeout: nil,
		           readOnly: false,
		           readTimeout: nil,
		           routeByLatency: false,
		           routeRandomly: false,
		           tlsConfig: nil,
		           writeTimeout: nil,
		           client: nil,
		           pingEnabled: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			rc := &redisClient{
				addrs:                test.fields.addrs,
				clusterSlots:         test.fields.clusterSlots,
				db:                   test.fields.db,
				dialTimeout:          test.fields.dialTimeout,
				dialer:               test.fields.dialer,
				idleCheckFrequency:   test.fields.idleCheckFrequency,
				idleTimeout:          test.fields.idleTimeout,
				initialPingDuration:  test.fields.initialPingDuration,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				keyPref:              test.fields.keyPref,
				maxConnAge:           test.fields.maxConnAge,
				maxRedirects:         test.fields.maxRedirects,
				maxRetries:           test.fields.maxRetries,
				maxRetryBackoff:      test.fields.maxRetryBackoff,
				minIdleConns:         test.fields.minIdleConns,
				minRetryBackoff:      test.fields.minRetryBackoff,
				onConnect:            test.fields.onConnect,
				onNewNode:            test.fields.onNewNode,
				password:             test.fields.password,
				poolSize:             test.fields.poolSize,
				poolTimeout:          test.fields.poolTimeout,
				readOnly:             test.fields.readOnly,
				readTimeout:          test.fields.readTimeout,
				routeByLatency:       test.fields.routeByLatency,
				routeRandomly:        test.fields.routeRandomly,
				tlsConfig:            test.fields.tlsConfig,
				writeTimeout:         test.fields.writeTimeout,
				client:               test.fields.client,
				pingEnabled:          test.fields.pingEnabled,
			}

			err := rc.ping(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
