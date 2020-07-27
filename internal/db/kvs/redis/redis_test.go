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
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
		goleak.IgnoreTopFunction("github.com/go-redis/redis/v7/internal/pool.(*ConnPool).reaper"),
		goleak.IgnoreTopFunction("github.com/go-redis/redis/v7.(*ClusterClient).reaper"),
	}
)

func TestMain(m *testing.M) {
	log.Init()
	code := m.Run()
	os.Exit(code)
}

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
		if !reflect.DeepEqual(w.wantRc, gotRc) {
			return errors.Errorf("got = %v, want = %v", gotRc, w.wantRc)
		}

		return nil
	}

	tests := []test{
		{
			name: "returns address not found error when options is nil",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				wantRc: nil,
				err:    errors.ErrRedisAddrsNotFound,
			},
		},

		{
			name: "returns ping failed error when options is not nil",
			args: args{
				ctx: context.Background(),
				opts: []Option{
					WithAddrs("127.0.0.0.1"),
					WithInitialPingTimeLimit("1µs"),
					WithInitialPingDuration("10ms"),
				},
			},
			want: want{
				wantRc: nil,
				err:    errors.Wrap(errors.Wrap(nil, errors.ErrRedisConnectionPingFailed.Error()), context.DeadlineExceeded.Error()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

func Test_redisClient_newRedisClient(t *testing.T) {
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
	}
	type want struct {
		want *redisClient
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *redisClient, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *redisClient, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			dialer := func(ctx context.Context, _, _ string) (net.Conn, error) {
				return nil, nil
			}
			connFn := func(c *redis.Conn) error {
				return nil
			}
			cfg := new(tls.Config)

			return test{
				name: "returns Redis implementation when address length is 1",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					addrs:              []string{"127.0.0.1"},
					password:           "pass",
					dialer:             dialer,
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
				},
				want: want{
					want: &redisClient{
						client: redis.NewClient(&redis.Options{
							Addr:               "127.0.0.1",
							Password:           "pass",
							Dialer:             dialer,
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
						}),
					},
					err: nil,
				},
				checkFunc: func(w want, gotRc *redisClient, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got error = %v, want %v", err, w.err)
					}
					if gotRc == nil {
						return errors.New("got is nil")
					}

					var (
						want = w.want.client.(*redis.Client).Options()
						got  = gotRc.client.(*redis.Client).Options()
					)

					opts := []cmp.Option{
						cmpopts.IgnoreUnexported(*want),
						cmpopts.IgnoreUnexported(*got),
						cmp.Comparer(func(want, got *tls.Config) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(ctx context.Context, network, addr string) (net.Conn, error)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(*redis.Conn) error) bool {
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

		func() test {
			dialer := func(ctx context.Context, _, _ string) (net.Conn, error) {
				return nil, nil
			}
			cslots := func() ([]redis.ClusterSlot, error) {
				return nil, nil
			}
			onNewNode := func(*redis.Client) {}
			onConnect := func(c *redis.Conn) error {
				return nil
			}
			cfg := new(tls.Config)

			return test{
				name: "returns Redis implementation when address length is 2",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					addrs:              []string{"127.0.0.1", "127.0.0.2"},
					dialer:             dialer,
					maxRedirects:       1,
					readOnly:           true,
					routeByLatency:     true,
					routeRandomly:      true,
					clusterSlots:       cslots,
					onNewNode:          onNewNode,
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
				},
				want: want{
					want: &redisClient{
						client: redis.NewClusterClient(&redis.ClusterOptions{
							Addrs: []string{
								"127.0.0.1", "127.0.0.2",
							},
							Dialer:             dialer,
							MaxRedirects:       1,
							ReadOnly:           true,
							RouteByLatency:     true,
							RouteRandomly:      true,
							ClusterSlots:       cslots,
							OnNewNode:          onNewNode,
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
						}),
					},
					err: nil,
				},
				checkFunc: func(w want, gotRc *redisClient, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got error = %v, want %v", err, w.err)
					}
					if gotRc == nil {
						return errors.New("got is nil")
					}

					var (
						want = w.want.client.(*redis.ClusterClient).Options()
						got  = gotRc.client.(*redis.ClusterClient).Options()
					)

					opts := []cmp.Option{
						cmpopts.IgnoreUnexported(*want),
						cmpopts.IgnoreUnexported(*got),
						cmp.Comparer(func(want, got func(opt *redis.Options) *redis.Client) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(*redis.Client)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func() ([]redis.ClusterSlot, error)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(ctx context.Context, network, addr string) (net.Conn, error)) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got func(*redis.Conn) error) bool {
							return reflect.ValueOf(want).Pointer() == reflect.ValueOf(got).Pointer()
						}),
						cmp.Comparer(func(want, got *tls.Config) bool {
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

		func() test {
			return test{
				name: "returns address not found error when address length is 0",
				args: args{
					ctx: context.Background(),
				},
				want: want{
					want: nil,
					err:  errors.ErrRedisAddrsNotFound,
				},
			}
		}(),

		func() test {
			return test{
				name: "returns address not found error when address length is 1 but contains empty string",
				fields: fields{
					addrs: []string{""},
				},
				args: args{
					ctx: context.Background(),
				},
				want: want{
					want: nil,
					err:  errors.ErrRedisAddrsNotFound,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
			}

			got, err := rc.newRedisClient(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotR, w.wantR) {
			return errors.Errorf("got = %v, want %v", gotR, w.wantR)
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
				initialPingDuration:  test.fields.initialPingDuration,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				client:               test.fields.client,
			}

			gotR, err := rc.ping(test.args.ctx)
			if err := test.checkFunc(test.want, gotR, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
