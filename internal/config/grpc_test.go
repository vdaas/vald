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
	"io/fs"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	testdata "github.com/vdaas/vald/internal/test"
	"go.uber.org/goleak"
)

func Test_newGRPCClientConfig(t *testing.T) {
	t.Parallel()
	type want struct {
		want *GRPCClient
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *GRPCClient) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *GRPCClient) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return GRPCClient when called newGRPCClientConfig()",
			want: want{
				want: &GRPCClient{
					DialOption: &DialOption{
						Insecure: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := newGRPCClientConfig()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGRPCClient_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Addrs               []string
		HealthCheckDuration string
		ConnectionPool      *ConnectionPool
		Backoff             *Backoff
		CallOption          *CallOption
		DialOption          *DialOption
		TLS                 *TLS
	}
	type want struct {
		want *GRPCClient
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *GRPCClient) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *GRPCClient) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			addrs := []string{
				"10.40.3.342",
				"10.40.98.17",
				"10.40.84.215",
			}
			healthcheck := "30s"
			return test{
				name: "return GRPCClient when only addrs and healthcheck duration are set",
				fields: fields{
					Addrs:               addrs,
					HealthCheckDuration: healthcheck,
				},
				want: want{
					want: &GRPCClient{
						Addrs:               addrs,
						HealthCheckDuration: healthcheck,
						ConnectionPool:      &ConnectionPool{},
						DialOption: &DialOption{
							Insecure: true,
						},
						TLS: &TLS{
							Enabled: false,
						},
					},
				},
			}
		}(),
		func() test {
			addrs := []string{
				"10.40.3.342",
				"10.40.98.17",
				"10.40.84.215",
			}
			healthcheck := "30s"
			connectionPool := &ConnectionPool{
				ResolveDNS:           true,
				EnableRebalance:      true,
				RebalanceDuration:    "5m",
				Size:                 100,
				OldConnCloseDuration: "3m",
			}
			backoffOpts := &Backoff{
				InitialDuration:  "5m",
				BackoffTimeLimit: "10m",
				MaximumDuration:  "15m",
				JitterLimit:      "3m",
				BackoffFactor:    3,
				RetryCount:       100,
				EnableErrorLog:   true,
			}
			callOpts := &CallOption{
				WaitForReady:          true,
				MaxRetryRPCBufferSize: 100,
				MaxRecvMsgSize:        1000,
				MaxSendMsgSize:        1000,
			}
			dialOpts := &DialOption{
				WriteBufferSize:             10000,
				ReadBufferSize:              10000,
				InitialWindowSize:           100,
				InitialConnectionWindowSize: 100,
				MaxMsgSize:                  1000,
				BackoffMaxDelay:             "3m",
				BackoffBaseDelay:            "1m",
				BackoffJitter:               100,
				BackoffMultiplier:           10,
				MinimumConnectionTimeout:    "5m",
				EnableBackoff:               true,
				Insecure:                    true,
				Timeout:                     "5m",
				Net:                         &Net{},
				KeepAlive: &GRPCClientKeepalive{
					Time:                "100s",
					Timeout:             "300s",
					PermitWithoutStream: true,
				},
			}
			tls := &TLS{
				Enabled: true,
				Cert:    "cert",
				Key:     "key",
				CA:      "ca",
			}
			return test{
				name: "return GRPCClient when all parameters are set",
				fields: fields{
					Addrs:               addrs,
					HealthCheckDuration: healthcheck,
					ConnectionPool:      connectionPool,
					Backoff:             backoffOpts,
					CallOption:          callOpts,
					DialOption:          dialOpts,
					TLS:                 tls,
				},
				want: want{
					want: &GRPCClient{
						Addrs:               addrs,
						HealthCheckDuration: healthcheck,
						ConnectionPool:      connectionPool,
						Backoff:             backoffOpts,
						CallOption:          callOpts,
						DialOption:          dialOpts,
						TLS:                 tls,
					},
				},
			}
		}(),
		func() test {
			addrs := []string{
				"10.40.3.342",
				"10.40.98.17",
				"10.40.84.215",
			}
			key := "HEALTH_CHECK_DURATION"
			value := "30s"
			return test{
				name: "return GRPCClient when only healthcheck duration is set as environment value",
				fields: fields{
					Addrs:               addrs,
					HealthCheckDuration: "_" + key + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					if err := os.Setenv(key, value); err != nil {
						t.Fatal(err)
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					if err := os.Unsetenv(key); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: &GRPCClient{
						Addrs:               addrs,
						HealthCheckDuration: value,
						ConnectionPool:      &ConnectionPool{},
						DialOption: &DialOption{
							Insecure: true,
						},
						TLS: &TLS{
							Enabled: false,
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return GRPCClient when all parameters are not set",
				fields: fields{},
				want: want{
					want: &GRPCClient{
						ConnectionPool: &ConnectionPool{},
						DialOption: &DialOption{
							Insecure: true,
						},
						TLS: &TLS{
							Enabled: false,
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &GRPCClient{
				Addrs:               test.fields.Addrs,
				HealthCheckDuration: test.fields.HealthCheckDuration,
				ConnectionPool:      test.fields.ConnectionPool,
				Backoff:             test.fields.Backoff,
				CallOption:          test.fields.CallOption,
				DialOption:          test.fields.DialOption,
				TLS:                 test.fields.TLS,
			}

			got := g.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGRPCClientKeepalive_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Time                string
		Timeout             string
		PermitWithoutStream bool
	}
	type want struct {
		want *GRPCClientKeepalive
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *GRPCClientKeepalive) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *GRPCClientKeepalive) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			time := "100s"
			timeout := "300s"
			permitWithoutStream := true
			return test{
				name: "return GRPCClientKeepalive when parameters are set",
				fields: fields{
					Time:                time,
					Timeout:             timeout,
					PermitWithoutStream: permitWithoutStream,
				},
				want: want{
					want: &GRPCClientKeepalive{
						Time:                time,
						Timeout:             timeout,
						PermitWithoutStream: permitWithoutStream,
					},
				},
			}
		}(),
		func() test {
			p := map[string]string{
				"TIME":    "100s",
				"TIMEOUT": "300s",
			}
			return test{
				name: "return GRPCClientKeepalive when parameters are set as environment value",
				fields: fields{
					Time:    "_TIME_",
					Timeout: "_TIMEOUT_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for key, value := range p {
						if err := os.Setenv(key, value); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for key := range p {
						if err := os.Unsetenv(key); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &GRPCClientKeepalive{
						Time:    "100s",
						Timeout: "300s",
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return GRPCClientKeepalive when all parameters are not set",
				fields: fields{},
				want: want{
					want: &GRPCClientKeepalive{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &GRPCClientKeepalive{
				Time:                test.fields.Time,
				Timeout:             test.fields.Timeout,
				PermitWithoutStream: test.fields.PermitWithoutStream,
			}

			got := g.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCallOption_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		WaitForReady          bool
		MaxRetryRPCBufferSize int
		MaxRecvMsgSize        int
		MaxSendMsgSize        int
	}
	type want struct {
		want *CallOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *CallOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *CallOption) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			waitForReady := true
			maxRetryRPCBufferSize := 100
			maxRecvMsgSize := 1000
			maxSendMsgSize := 1000
			return test{
				name: "return CallOption when all parameters are set",
				fields: fields{
					WaitForReady:          waitForReady,
					MaxRetryRPCBufferSize: maxRetryRPCBufferSize,
					MaxRecvMsgSize:        maxRecvMsgSize,
					MaxSendMsgSize:        maxSendMsgSize,
				},
				want: want{
					want: &CallOption{
						WaitForReady:          waitForReady,
						MaxRetryRPCBufferSize: maxRetryRPCBufferSize,
						MaxRecvMsgSize:        maxRecvMsgSize,
						MaxSendMsgSize:        maxSendMsgSize,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return CallOption when all parameters are not set",
				fields: fields{},
				want: want{
					want: &CallOption{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &CallOption{
				WaitForReady:          test.fields.WaitForReady,
				MaxRetryRPCBufferSize: test.fields.MaxRetryRPCBufferSize,
				MaxRecvMsgSize:        test.fields.MaxRecvMsgSize,
				MaxSendMsgSize:        test.fields.MaxSendMsgSize,
			}

			got := c.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDialOption_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		WriteBufferSize             int
		ReadBufferSize              int
		InitialWindowSize           int
		InitialConnectionWindowSize int
		MaxMsgSize                  int
		BackoffMaxDelay             string
		BackoffBaseDelay            string
		BackoffJitter               float64
		BackoffMultiplier           float64
		MinimumConnectionTimeout    string
		EnableBackoff               bool
		Insecure                    bool
		Timeout                     string
		Net                         *Net
		KeepAlive                   *GRPCClientKeepalive
	}
	type want struct {
		want *DialOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *DialOption) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *DialOption) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			writeBufferSize := 10000
			readBufferSize := 10000
			initialWindowSize := 100
			initialConnectionWindowSize := 100
			maxMsgSize := 1000
			backoffMaxDelay := "3m"
			backoffBaseDelay := "1m"
			backoffJitter := float64(100)
			backoffMultiplier := float64(10)
			minimumConnectionTimeout := "5m"
			enableBackoff := true
			insecure := true
			timeout := "5m"
			net := &Net{}
			keepAlive := &GRPCClientKeepalive{
				Time:                "100s",
				Timeout:             "300s",
				PermitWithoutStream: true,
			}
			return test{
				name: "return DialOption when all parameters are set",
				fields: fields{
					WriteBufferSize:             writeBufferSize,
					ReadBufferSize:              readBufferSize,
					InitialWindowSize:           initialWindowSize,
					InitialConnectionWindowSize: initialConnectionWindowSize,
					MaxMsgSize:                  maxMsgSize,
					BackoffMaxDelay:             backoffMaxDelay,
					BackoffBaseDelay:            backoffBaseDelay,
					BackoffJitter:               backoffJitter,
					BackoffMultiplier:           backoffMultiplier,
					MinimumConnectionTimeout:    minimumConnectionTimeout,
					EnableBackoff:               enableBackoff,
					Insecure:                    insecure,
					Timeout:                     timeout,
					Net:                         net,
					KeepAlive:                   keepAlive,
				},
				want: want{
					want: &DialOption{
						WriteBufferSize:             writeBufferSize,
						ReadBufferSize:              readBufferSize,
						InitialWindowSize:           initialWindowSize,
						InitialConnectionWindowSize: initialConnectionWindowSize,
						MaxMsgSize:                  maxMsgSize,
						BackoffMaxDelay:             backoffMaxDelay,
						BackoffBaseDelay:            backoffBaseDelay,
						BackoffJitter:               backoffJitter,
						BackoffMultiplier:           backoffMultiplier,
						MinimumConnectionTimeout:    minimumConnectionTimeout,
						EnableBackoff:               enableBackoff,
						Insecure:                    insecure,
						Timeout:                     timeout,
						Net:                         net,
						KeepAlive:                   keepAlive,
					},
				},
			}
		}(),
		func() test {
			p := map[string]string{
				"BACKOFF_MAX_DELAY": "3m",
				"TIMEOUT":           "3m",
			}
			return test{
				name: "return DialOption when parameters are set as environment value",
				fields: fields{
					BackoffMaxDelay: "_BACKOFF_MAX_DELAY_",
					Timeout:         "_TIMEOUT_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for key, value := range p {
						if err := os.Setenv(key, value); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for key := range p {
						if err := os.Unsetenv(key); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &DialOption{
						BackoffMaxDelay: "3m",
						Timeout:         "3m",
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return DialOption when all parameters are not set",
				fields: fields{},
				want: want{
					want: &DialOption{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			d := &DialOption{
				WriteBufferSize:             test.fields.WriteBufferSize,
				ReadBufferSize:              test.fields.ReadBufferSize,
				InitialWindowSize:           test.fields.InitialWindowSize,
				InitialConnectionWindowSize: test.fields.InitialConnectionWindowSize,
				MaxMsgSize:                  test.fields.MaxMsgSize,
				BackoffMaxDelay:             test.fields.BackoffMaxDelay,
				BackoffBaseDelay:            test.fields.BackoffBaseDelay,
				BackoffJitter:               test.fields.BackoffJitter,
				BackoffMultiplier:           test.fields.BackoffMultiplier,
				MinimumConnectionTimeout:    test.fields.MinimumConnectionTimeout,
				EnableBackoff:               test.fields.EnableBackoff,
				Insecure:                    test.fields.Insecure,
				Timeout:                     test.fields.Timeout,
				Net:                         test.fields.Net,
				KeepAlive:                   test.fields.KeepAlive,
			}

			got := d.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGRPCClient_Opts(t *testing.T) {
	t.Parallel()
	type fields struct {
		Addrs               []string
		HealthCheckDuration string
		ConnectionPool      *ConnectionPool
		Backoff             *Backoff
		CallOption          *CallOption
		DialOption          *DialOption
		TLS                 *TLS
	}
	type want struct {
		want []grpc.Option
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []grpc.Option, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotOpts []grpc.Option, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(len(gotOpts), len(w.want)) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOpts, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return 25 grpc.Option and nil error when all parameters are set",
			fields: fields{
				Addrs: []string{
					"10.40.3.342",
					"10.40.98.17",
					"10.40.84.215",
				},
				HealthCheckDuration: "30s",
				ConnectionPool: &ConnectionPool{
					ResolveDNS:           true,
					EnableRebalance:      true,
					RebalanceDuration:    "5m",
					Size:                 100,
					OldConnCloseDuration: "3m",
				},
				Backoff: &Backoff{
					InitialDuration:  "5m",
					BackoffTimeLimit: "10m",
					MaximumDuration:  "15m",
					JitterLimit:      "3m",
					BackoffFactor:    3,
					RetryCount:       100,
					EnableErrorLog:   true,
				},
				CallOption: &CallOption{
					WaitForReady:          true,
					MaxRetryRPCBufferSize: 100,
					MaxRecvMsgSize:        1000,
					MaxSendMsgSize:        1000,
				},
				DialOption: &DialOption{
					WriteBufferSize:             10000,
					ReadBufferSize:              10000,
					InitialWindowSize:           100,
					InitialConnectionWindowSize: 100,
					MaxMsgSize:                  1000,
					BackoffMaxDelay:             "3m",
					BackoffBaseDelay:            "1m",
					BackoffJitter:               100,
					BackoffMultiplier:           10,
					MinimumConnectionTimeout:    "5m",
					EnableBackoff:               true,
					Insecure:                    false,
					Timeout:                     "5m",
					Net: &Net{
						Dialer: &Dialer{
							Timeout: "10m",
						},
						TLS: &TLS{
							Enabled: true,
							Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
							Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
							CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
						},
					},
					KeepAlive: &GRPCClientKeepalive{
						Time:                "100s",
						Timeout:             "300s",
						PermitWithoutStream: true,
					},
				},
				TLS: &TLS{
					Enabled: true,
					Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
					Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
					CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
				},
			},
			want: want{
				want: make([]grpc.Option, 25),
			},
		},
		{
			name: "return nil grpc.Option and an error when dns error is ocurred",
			fields: fields{
				Addrs: []string{
					"10.40.3.342",
					"10.40.98.17",
					"10.40.84.215",
				},
				HealthCheckDuration: "30s",
				ConnectionPool: &ConnectionPool{
					ResolveDNS:           true,
					EnableRebalance:      true,
					RebalanceDuration:    "5m",
					Size:                 100,
					OldConnCloseDuration: "3m",
				},
				Backoff: &Backoff{
					InitialDuration:  "5m",
					BackoffTimeLimit: "10m",
					MaximumDuration:  "15m",
					JitterLimit:      "3m",
					BackoffFactor:    3,
					RetryCount:       100,
					EnableErrorLog:   true,
				},
				CallOption: &CallOption{
					WaitForReady:          true,
					MaxRetryRPCBufferSize: 100,
					MaxRecvMsgSize:        1000,
					MaxSendMsgSize:        1000,
				},
				DialOption: &DialOption{
					WriteBufferSize:             10000,
					ReadBufferSize:              10000,
					InitialWindowSize:           100,
					InitialConnectionWindowSize: 100,
					MaxMsgSize:                  1000,
					BackoffMaxDelay:             "1m",
					BackoffBaseDelay:            "3m",
					BackoffJitter:               100,
					BackoffMultiplier:           10,
					MinimumConnectionTimeout:    "5m",
					EnableBackoff:               true,
					Insecure:                    false,
					Timeout:                     "5m",
					Net: &Net{
						Dialer: &Dialer{
							Timeout: "10m",
						},
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "3m",
							CacheExpiration: "1m",
						},
						TLS: &TLS{
							Enabled: true,
							Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
							Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
							CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
						},
					},
					KeepAlive: &GRPCClientKeepalive{
						Time:                "100s",
						Timeout:             "300s",
						PermitWithoutStream: true,
					},
				},
				TLS: &TLS{
					Enabled: true,
					Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
					Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
					CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
				},
			},
			want: want{
				want: make([]grpc.Option, 0),
				err:  errors.ErrInvalidDNSConfig(3*time.Minute, time.Minute),
			},
		},
		{
			name: "return nil grpc.Option and an error when tls error is ocurred",
			fields: fields{
				Addrs: []string{
					"10.40.3.342",
					"10.40.98.17",
					"10.40.84.215",
				},
				HealthCheckDuration: "30s",
				ConnectionPool: &ConnectionPool{
					ResolveDNS:           true,
					EnableRebalance:      true,
					RebalanceDuration:    "5m",
					Size:                 100,
					OldConnCloseDuration: "3m",
				},
				Backoff: &Backoff{
					InitialDuration:  "5m",
					BackoffTimeLimit: "10m",
					MaximumDuration:  "15m",
					JitterLimit:      "3m",
					BackoffFactor:    3,
					RetryCount:       100,
					EnableErrorLog:   true,
				},
				CallOption: &CallOption{
					WaitForReady:          true,
					MaxRetryRPCBufferSize: 100,
					MaxRecvMsgSize:        1000,
					MaxSendMsgSize:        1000,
				},
				DialOption: &DialOption{
					WriteBufferSize:             10000,
					ReadBufferSize:              10000,
					InitialWindowSize:           100,
					InitialConnectionWindowSize: 100,
					MaxMsgSize:                  1000,
					BackoffMaxDelay:             "1m",
					BackoffBaseDelay:            "3m",
					BackoffJitter:               100,
					BackoffMultiplier:           10,
					MinimumConnectionTimeout:    "5m",
					EnableBackoff:               true,
					Insecure:                    false,
					Timeout:                     "5m",
					Net: &Net{
						Dialer: &Dialer{
							Timeout: "10m",
						},
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "1m",
							CacheExpiration: "3m",
						},
						TLS: &TLS{
							Enabled: true,
							Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
							Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
							CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
						},
					},
					KeepAlive: &GRPCClientKeepalive{
						Time:                "100s",
						Timeout:             "300s",
						PermitWithoutStream: true,
					},
				},
				TLS: &TLS{
					Enabled: true,
					Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
					Key:     "tls/dummy/Server.key",
					CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
				},
			},
			want: want{
				want: make([]grpc.Option, 0),
				err: &fs.PathError{
					Op:   "open",
					Path: "tls/dummy/Server.key",
					Err:  syscall.Errno(0x2),
				},
			},
		},
		{
			name: "return nil grpc.Option and an error when net.TLS.Opts error is ocurred",
			fields: fields{
				Addrs: []string{
					"10.40.3.342",
					"10.40.98.17",
					"10.40.84.215",
				},
				HealthCheckDuration: "30s",
				ConnectionPool: &ConnectionPool{
					ResolveDNS:           true,
					EnableRebalance:      true,
					RebalanceDuration:    "5m",
					Size:                 100,
					OldConnCloseDuration: "3m",
				},
				Backoff: &Backoff{
					InitialDuration:  "5m",
					BackoffTimeLimit: "10m",
					MaximumDuration:  "15m",
					JitterLimit:      "3m",
					BackoffFactor:    3,
					RetryCount:       100,
					EnableErrorLog:   true,
				},
				CallOption: &CallOption{
					WaitForReady:          true,
					MaxRetryRPCBufferSize: 100,
					MaxRecvMsgSize:        1000,
					MaxSendMsgSize:        1000,
				},
				DialOption: &DialOption{
					WriteBufferSize:             10000,
					ReadBufferSize:              10000,
					InitialWindowSize:           100,
					InitialConnectionWindowSize: 100,
					MaxMsgSize:                  1000,
					BackoffMaxDelay:             "1m",
					BackoffBaseDelay:            "3m",
					BackoffJitter:               100,
					BackoffMultiplier:           10,
					MinimumConnectionTimeout:    "5m",
					EnableBackoff:               true,
					Insecure:                    false,
					Timeout:                     "5m",
					Net: &Net{
						Dialer: &Dialer{
							Timeout: "10m",
						},
						DNS: &DNS{
							CacheEnabled:    true,
							RefreshDuration: "1m",
							CacheExpiration: "3m",
						},
						TLS: &TLS{
							Enabled: true,
						},
					},
					KeepAlive: &GRPCClientKeepalive{
						Time:                "100s",
						Timeout:             "300s",
						PermitWithoutStream: true,
					},
				},
				TLS: &TLS{
					Enabled: true,
					Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
					Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
					CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
				},
			},
			want: want{
				want: make([]grpc.Option, 0),
				err:  errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			name:   "return 1 grpc.Option when all parameters are set",
			fields: fields{},
			want: want{
				want: make([]grpc.Option, 1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &GRPCClient{
				Addrs:               test.fields.Addrs,
				HealthCheckDuration: test.fields.HealthCheckDuration,
				ConnectionPool:      test.fields.ConnectionPool,
				Backoff:             test.fields.Backoff,
				CallOption:          test.fields.CallOption,
				DialOption:          test.fields.DialOption,
				TLS:                 test.fields.TLS,
			}

			got, err := g.Opts()
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
