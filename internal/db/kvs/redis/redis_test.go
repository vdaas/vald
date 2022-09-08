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

package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
	goleak.IgnoreTopFunction("github.com/go-redis/redis/v8/internal/pool.(*ConnPool).reaper"),
	goleak.IgnoreTopFunction("github.com/go-redis/redis/v8.(*ClusterClient).reaper"),
}

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	code := m.Run()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantRc *redisClient
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Connector, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRc Connector, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(w.wantRc, gotRc) {
			return errors.Errorf("got = %v, want = %v", gotRc, w.wantRc)
		}

		return nil
	}

	tests := []test{
		{
			name: "returns Connector instance",
			args: args{},
			want: want{
				wantRc: &redisClient{
					initialPingDuration:  30 * time.Millisecond,
					initialPingTimeLimit: 5 * time.Minute,
					network:              net.TCP.String(),
				},
				err: nil,
			},
		},
		func() test {
			dummyErr := errors.New("error")

			opt := dummyWithFunc(dummyErr)
			return test{
				name: "returns error when applying options failed",
				args: args{
					opts: []Option{
						opt,
					},
				},
				want: want{
					wantRc: nil,
					err:    errors.ErrOptionFailed(dummyErr, reflect.ValueOf(opt)),
				},
				checkFunc: func(w want, gotRc Connector, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}

					return nil
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

			gotRc, err := New(test.args.opts...)
			if err := checkFunc(test.want, gotRc, err); err != nil {
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
		initialPingDuration  time.Duration
		initialPingTimeLimit time.Duration
		client               Redis
	}
	type want struct {
		wantR Redis
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Redis, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotR Redis, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotR, w.wantR) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotR, w.wantR)
		}
		return nil
	}
	tests := []test{
		func() test {
			r := &MockRedis{
				PingFunc: func() *StatusCmd {
					return new(StatusCmd)
				},
			}

			return test{
				name: "returns nil when the ping success",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					initialPingDuration:  time.Millisecond,
					initialPingTimeLimit: time.Second,
					client:               r,
				},
				want: want{
					wantR: r,
					err:   nil,
				},
			}
		}(),

		func() test {
			err := errors.New("err")

			return test{
				name: "returns ping failed error when the ping fails and reached the ping time limit",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					initialPingDuration:  time.Millisecond,
					initialPingTimeLimit: 3 * time.Millisecond,
					client: func() Redis {
						return &MockRedis{
							PingFunc: func() (cmd *StatusCmd) {
								cmd = new(StatusCmd)
								cmd.SetErr(err)
								return
							},
						}
					}(),
				},
				want: want{
					wantR: nil,
					err:   errors.Wrap(errors.Wrap(err, errors.ErrRedisConnectionPingFailed.Error()), context.DeadlineExceeded.Error()),
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
			rc := &redisClient{
				initialPingDuration:  test.fields.initialPingDuration,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				client:               test.fields.client,
			}

			gotR, err := rc.ping(test.args.ctx)
			if err := checkFunc(test.want, gotR, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisClient_setClient(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		addrs                []string
		clusterSlots         func(context.Context) ([]redis.ClusterSlot, error)
		db                   int
		dialTimeout          time.Duration
		dialer               net.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
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
		onConnect            func(ctx context.Context, conn *redis.Conn) error
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
		hooks                []Hook
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
			name: "returns error when addrs not specified",
			args: args{context.Background()},
			want: want{
				err: errors.ErrRedisAddrsNotFound,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "returns error when addrs is empty",
			args: args{context.Background()},
			fields: fields{
				addrs: []string{},
			},
			want: want{
				err: errors.ErrRedisAddrsNotFound,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "returns nil when addrs is single addr",
			args: args{context.Background()},
			fields: fields{
				addrs: []string{"127.0.0.1:6379"},
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "returns nil when addrs is single addr and it is empty string",
			args: args{context.Background()},
			fields: fields{
				addrs: []string{""},
			},
			want: want{
				err: errors.ErrRedisAddrsNotFound,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "returns nil when addrs is multiple addrs",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				addrs: []string{
					"127.0.0.1:6379",
					"127.0.0.2:6379",
				},
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "returns error when addrs is multiple addrs and it has empty string",
			fields: fields{
				addrs: []string{
					"",
					"127.0.0.2:6379",
				},
			},
			want: want{
				err: errors.ErrRedisAddrsNotFound,
			},
			checkFunc: defaultCheckFunc,
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
			rc := &redisClient{
				addrs:                test.fields.addrs,
				clusterSlots:         test.fields.clusterSlots,
				db:                   test.fields.db,
				dialTimeout:          test.fields.dialTimeout,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
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
				hooks:                test.fields.hooks,
			}

			err := rc.setClient(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisClient_newClient(t *testing.T) {
	type fields struct {
		addrs                []string
		clusterSlots         func(context.Context) ([]redis.ClusterSlot, error)
		network              string
		db                   int
		dialTimeout          time.Duration
		dialer               net.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
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
		onConnect            func(ctx context.Context, conn *redis.Conn) error
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
		hooks                []Hook
	}
	type want struct {
		want *redis.Client
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *redis.Client, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *redis.Client, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns error when first element of addrs is empty string",
			fields: fields{
				addrs: []string{""},
			},
			want: want{
				want: nil,
				err:  errors.ErrRedisAddrsNotFound,
			},
			checkFunc: defaultCheckFunc,
		},
		func() test {
			dialer := func(ctx context.Context, _, _ string) (net.Conn, error) {
				return nil, nil
			}
			connFn := func(ctx context.Context, c *redis.Conn) error {
				return nil
			}
			cfg := new(tls.Config)
			hook := &dummyHook{}

			rc := redis.NewClient(&redis.Options{
				Addr:               "127.0.0.1:6379",
				Password:           "pass",
				Dialer:             dialer,
				Network:            net.TCP.String(),
				OnConnect:          connFn,
				DB:                 1,
				MaxRetries:         2,
				MinRetryBackoff:    3 * time.Second,
				MaxRetryBackoff:    4 * time.Second,
				DialTimeout:        5 * time.Second,
				ReadTimeout:        6 * time.Second,
				WriteTimeout:       7 * time.Second,
				PoolSize:           8,
				MinIdleConns:       9,
				MaxConnAge:         10 * time.Second,
				PoolTimeout:        11 * time.Second,
				IdleTimeout:        12 * time.Second,
				IdleCheckFrequency: 13 * time.Second,
				TLSConfig:          cfg,
			})
			rc.AddHook(hook)

			return test{
				name: "returns redis.Client successfully",
				fields: fields{
					addrs:              []string{"127.0.0.1:6379"},
					network:            net.TCP.String(),
					password:           "pass",
					dialerFunc:         dialer,
					onConnect:          connFn,
					db:                 1,
					maxRetries:         2,
					minRetryBackoff:    3 * time.Second,
					maxRetryBackoff:    4 * time.Second,
					dialTimeout:        5 * time.Second,
					readTimeout:        6 * time.Second,
					writeTimeout:       7 * time.Second,
					poolSize:           8,
					minIdleConns:       9,
					maxConnAge:         10 * time.Second,
					poolTimeout:        11 * time.Second,
					idleTimeout:        12 * time.Second,
					idleCheckFrequency: 13 * time.Second,
					tlsConfig:          cfg,
					hooks:              []redis.Hook{hook},
				},
				want: want{
					want: rc,
				},
				checkFunc: func(w want, gotc *redis.Client, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					if gotc == nil {
						return errors.New("got is nil")
					}

					var (
						want = w.want.Options()
						got  = gotc.Options()
					)

					opts := []cmp.Option{
						cmpopts.IgnoreUnexported(*want),
						cmpopts.IgnoreUnexported(*got),
						cmpopts.IgnoreFields(redis.Options{}, "OnConnect"),
						cmp.Comparer(func(want, got *tls.Config) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(ctx context.Context, network, addr string) (net.Conn, error)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(*redis.Conn) error) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got []redis.Hook) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
					}
					if diff := cmp.Diff(want, got, opts...); diff != "" {
						return errors.Errorf("client options diff: %s", diff)
					}

					return nil
				},
			}
		}(),
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
			rc := &redisClient{
				addrs:                test.fields.addrs,
				clusterSlots:         test.fields.clusterSlots,
				db:                   test.fields.db,
				dialTimeout:          test.fields.dialTimeout,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
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
				hooks:                test.fields.hooks,
			}

			got, err := rc.newClient(context.Background())
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisClient_newClusterClient(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		addrs                []string
		clusterSlots         func(context.Context) ([]redis.ClusterSlot, error)
		db                   int
		dialTimeout          time.Duration
		dialer               net.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
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
		onConnect            func(ctx context.Context, conn *redis.Conn) error
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
		hooks                []Hook
	}
	type want struct {
		want *redis.ClusterClient
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *redis.ClusterClient, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *redis.ClusterClient, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns error when first element of addrs is empty string",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				addrs: []string{""},
			},
			want: want{
				want: nil,
				err:  errors.ErrRedisAddrsNotFound,
			},
			checkFunc: defaultCheckFunc,
		},
		func() test {
			dialer := func(ctx context.Context, _, _ string) (net.Conn, error) {
				return nil, nil
			}
			cslots := func(ctx context.Context) ([]redis.ClusterSlot, error) {
				return nil, nil
			}
			onConnect := func(ctx context.Context, c *redis.Conn) error {
				return nil
			}
			cfg := new(tls.Config)
			hook := &dummyHook{}

			rc := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:              []string{"127.0.0.1:6379"},
				Dialer:             dialer,
				MaxRedirects:       1,
				ReadOnly:           true,
				RouteByLatency:     true,
				RouteRandomly:      true,
				ClusterSlots:       cslots,
				OnConnect:          onConnect,
				Password:           "pass",
				MaxRetries:         2,
				MinRetryBackoff:    3 * time.Second,
				MaxRetryBackoff:    4 * time.Second,
				DialTimeout:        5 * time.Second,
				ReadTimeout:        6 * time.Second,
				WriteTimeout:       7 * time.Second,
				PoolSize:           8,
				MaxConnAge:         9 * time.Second,
				IdleTimeout:        10 * time.Second,
				IdleCheckFrequency: 11 * time.Second,
				TLSConfig:          cfg,
			})
			rc.AddHook(hook)

			return test{
				name: "returns redis.Client successfully",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					addrs:              []string{"127.0.0.1:6379"},
					dialerFunc:         dialer,
					maxRedirects:       1,
					readOnly:           true,
					routeByLatency:     true,
					routeRandomly:      true,
					clusterSlots:       cslots,
					onConnect:          onConnect,
					password:           "pass",
					maxRetries:         2,
					minRetryBackoff:    3 * time.Second,
					maxRetryBackoff:    4 * time.Second,
					dialTimeout:        5 * time.Second,
					readTimeout:        6 * time.Second,
					writeTimeout:       7 * time.Second,
					poolSize:           8,
					maxConnAge:         9 * time.Second,
					idleTimeout:        10 * time.Second,
					idleCheckFrequency: 11 * time.Second,
					tlsConfig:          cfg,
					hooks:              []redis.Hook{hook},
				},
				want: want{
					want: rc,
				},
				checkFunc: func(w want, gotc *redis.ClusterClient, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					if gotc == nil {
						return errors.New("got is nil")
					}

					var (
						want = w.want.Options()
						got  = gotc.Options()
					)

					opts := []cmp.Option{
						cmpopts.IgnoreUnexported(*want),
						cmpopts.IgnoreUnexported(*got),
						cmp.Comparer(func(want, got func(opt *redis.Options) *redis.Client) bool {
							// TODO fix this code later
							return true
						}),
						cmp.Comparer(func(want, got func(*redis.Client)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(context.Context) ([]redis.ClusterSlot, error)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func() ([]redis.ClusterSlot, error)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(ctx context.Context, network, addr string) (net.Conn, error)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(context.Context, *redis.Conn) error) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(*redis.Conn) error) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got *tls.Config) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got []redis.Hook) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
					}
					if diff := cmp.Diff(want, got, opts...); diff != "" {
						fmt.Println(diff)
						return errors.Errorf("got = %v, want = %v", got, want)
					}

					return nil
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
			rc := &redisClient{
				addrs:                test.fields.addrs,
				clusterSlots:         test.fields.clusterSlots,
				db:                   test.fields.db,
				dialTimeout:          test.fields.dialTimeout,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
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
				hooks:                test.fields.hooks,
			}

			got, err := rc.newClusterClient(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisClient_Connect(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		addrs                []string
		clusterSlots         func(ctx context.Context) ([]redis.ClusterSlot, error)
		db                   int
		dialTimeout          time.Duration
		dialer               net.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
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
		onConnect            func(ctx context.Context, conn *redis.Conn) error
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
		hooks                []Hook
	}
	type want struct {
		want Redis
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Redis, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Redis, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			dialer, err := net.NewDialer()
			if err != nil {
				panic(err)
			}

			return test{
				name: "returns error when addrs not specified",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					addrs:  []string{""},
					dialer: dialer,
				},
				want: want{
					err: errors.ErrRedisAddrsNotFound,
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			dialer, err := net.NewDialer()
			if err != nil {
				panic(err)
			}

			return test{
				name: "returns error when an invalid addrs specified",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					addrs:                []string{"127.0.0.1:6379"},
					initialPingTimeLimit: time.Microsecond,
					initialPingDuration:  10 * time.Millisecond,
					dialer:               dialer,
				},
				want: want{
					err: errors.Wrap(errors.Wrap(nil, errors.ErrRedisConnectionPingFailed.Error()), context.DeadlineExceeded.Error()),
				},
				checkFunc: defaultCheckFunc,
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
			rc := &redisClient{
				addrs:                test.fields.addrs,
				clusterSlots:         test.fields.clusterSlots,
				db:                   test.fields.db,
				dialTimeout:          test.fields.dialTimeout,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
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
				hooks:                test.fields.hooks,
			}

			got, err := rc.Connect(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
