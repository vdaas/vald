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

	"github.com/go-redis/redis/v7"
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
		func() test {
			dialer := func(ctx context.Context, addr, port string) (net.Conn, error) {
				return nil, nil
			}
			connFn := func(*redis.Conn) error {
				return nil
			}
			cfg := new(tls.Config)

			return test{
				name: "returns Redis implementation when address length is 1",
				args: args{
					ctx: context.Background(),
					opts: []Option{
						WithAddrs("127.0.0.1"),
						WithPassword("pass"),
						WithDialer(dialer),
						WithOnConnectFunction(connFn),
						WithDB(1),
						WithRetryLimit(2),
						WithMinimumRetryBackoff("3s"),
						WithMaximumRetryBackoff("4s"),
						WithDialTimeout("5s"),
						WithReadTimeout("6s"),
						WithWriteTimeout("7s"),
						WithPoolSize(8),
						WithMinimumIdleConnection(9),
						WithMaximumConnectionAge("10s"),
						WithPoolTimeout("11s"),
						WithIdleTimeout("12s"),
						WithIdleCheckFrequency("13s"),
						WithTLSConfig(cfg),
						WithPing(false),
					},
				},
				want: want{
					wantRc: redis.NewClient(&redis.Options{
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
					err: nil,
				},
				checkFunc: func(w want, gotRc Redis, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got error = %v, want %v", err, w.err)
					}
					if gotRc == nil {
						return errors.New("got is nil")
					}

					var (
						want = w.wantRc.(*redis.Client).Options()
						got  = gotRc.(*redis.Client).Options()
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
						return errors.Errorf("got = %v, want = %v", got, want)
					}

					return nil
				},
			}
		}(),

		func() test {
			dialer := func(ctx context.Context, addr, port string) (net.Conn, error) {
				return nil, nil
			}
			closterSlots := func() ([]redis.ClusterSlot, error) {
				return nil, nil
			}
			onNewNode := func(*redis.Client) {}
			onConnect := func(*redis.Conn) error {
				return nil
			}
			cfg := new(tls.Config)

			return test{
				name: "returns Redis implementation when address length is 2",
				args: args{
					ctx: context.Background(),
					opts: []Option{
						WithAddrs("127.0.0.1", "192.168.33.10"),
						WithDialer(dialer),
						WithRedirectLimit(1),
						WithReadOnlyFlag(true),
						WithRouteByLatencyFlag(true),
						WithRouteRandomlyFlag(true),
						WithClusterSlots(closterSlots),
						WithOnNewNodeFunction(onNewNode),
						WithOnConnectFunction(onConnect),
						WithPassword("pass"),
						WithRetryLimit(2),
						WithMinimumRetryBackoff("3s"),
						WithMaximumRetryBackoff("4s"),
						WithDialTimeout("5s"),
						WithReadTimeout("6s"),
						WithWriteTimeout("7s"),
						WithPoolSize(8),
						WithMinimumIdleConnection(9),
						WithMaximumConnectionAge("10s"),
						WithPoolTimeout("11s"),
						WithIdleTimeout("12s"),
						WithIdleCheckFrequency("13s"),
						WithTLSConfig(cfg),
						WithPing(false),
					},
				},
				want: want{
					wantRc: redis.NewClusterClient(&redis.ClusterOptions{
						Addrs: []string{
							"127.0.0.1", "192.168.33.10",
						},
						Dialer:             dialer,
						MaxRedirects:       1,
						ReadOnly:           true,
						RouteByLatency:     true,
						RouteRandomly:      true,
						ClusterSlots:       closterSlots,
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
						MinIdleConns:       9,
						MaxConnAge:         10 * time.Second,
						PoolTimeout:        11 * time.Second,
						IdleTimeout:        12 * time.Second,
						IdleCheckFrequency: 13 * time.Second,
						TLSConfig:          cfg,
					}),
					err: nil,
				},
				checkFunc: func(w want, gotRc Redis, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got error = %v, want %v", err, w.err)
					}
					if gotRc == nil {
						return errors.New("got is nil")
					}

					var (
						want = w.wantRc.(*redis.ClusterClient).Options()
						got  = gotRc.(*redis.ClusterClient).Options()
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
				name: "returns redis address not found error when address length is 0",
				args: args{
					ctx:  context.Background(),
					opts: nil,
				},
				want: want{
					wantRc: nil,
					err:    errors.ErrRedisAddrsNotFound,
				},
			}
		}(),

		func() test {
			return test{
				name: "returns redis address not found error when address length is 1 but address contains empty string",
				args: args{
					ctx: context.Background(),
					opts: []Option{
						WithAddrs(""),
					},
				},
				want: want{
					wantRc: nil,
					err:    errors.ErrRedisAddrsNotFound,
				},
			}
		}(),

		func() test {
			err := errors.New("err")
			return test{
				name: "returns ping error when address length is 1 and ping fails",
				args: args{
					ctx: context.Background(),
					opts: []Option{
						WithAddrs("127.0.0.01"),
						WithInitialPingDuration("1ms"),
						WithInitialPingTimeLimit("2ms"),
						WithDialer(func(ctx context.Context, addr string, port string) (net.Conn, error) {
							return nil, err
						}),
					},
				},
				want: want{
					wantRc: nil,
					err: errors.Wrap(errors.Wrap(
						err,
						errors.ErrRedisConnectionPingFailed.Error()), context.DeadlineExceeded.Error()),
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
		initialPingDuration  time.Duration
		initialPingTimeLimit time.Duration
		client               Redis
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
		func() test {
			return test{
				name: "returns nil when the ping success",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					initialPingDuration:  time.Microsecond,
					initialPingTimeLimit: time.Second,
					client: func() Redis {
						return &MockRedis{
							PingFunc: func() *StatusCmd {
								return new(StatusCmd)
							},
						}
					}(),
				},
				want: want{
					err: nil,
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
					initialPingTimeLimit: 5 * time.Millisecond,
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
					err: errors.Wrap(errors.Wrap(err, errors.ErrRedisConnectionPingFailed.Error()), context.DeadlineExceeded.Error()),
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

			err := rc.ping(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
