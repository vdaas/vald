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

	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/errors"
	testdata "github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestRedis_Bind(t *testing.T) {
	type fields struct {
		Addrs                []string
		DB                   int
		DialTimeout          string
		IdleCheckFrequency   string
		IdleTimeout          string
		InitialPingDuration  string
		InitialPingTimeLimit string
		KVPrefix             string
		KeyPref              string
		MaxConnAge           string
		MaxRedirects         int
		MaxRetries           int
		MaxRetryBackoff      string
		MinIdleConns         int
		MinRetryBackoff      string
		Network              string
		Password             string
		PoolSize             int
		PoolTimeout          string
		PrefixDelimiter      string
		ReadOnly             bool
		ReadTimeout          string
		RouteByLatency       bool
		RouteRandomly        bool
		SentinelPassword     string
		SentinelMasterName   string
		Net                  *Net
		TLS                  *TLS
		Username             string
		VKPrefix             string
		WriteTimeout         string
	}
	type want struct {
		want *Redis
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Redis) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Redis) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			addrs := []string{"redis.default.svc.cluster.local:6379"}
			db := 0
			dialTimeout := "5s"
			idleCheckFrequency := "1m"
			IdleTimeout := "5m"
			initialPingDuration := "5s"
			initialPingTimelimit := "30s"
			keyPref := "vald"
			maxConnAge := "0s"
			maxRedirects := 3
			maxRetries := 1
			maxRetryBackoff := "512s"
			minIdleConns := 0
			minRetryBackoff := "8ms"
			network := "tcp"
			password := "password"
			poolSize := 10
			poolTimeout := "4s"
			prefixDelimiter := "+"
			readOnly := true
			readTimeout := "3s"
			routeByLatency := false
			routeRandomly := true
			sentinelPassword := ""
			sentinelMasterName := ""
			kvPrefix := ""
			vkPrefix := ""
			username := "vald"
			writeTimeout := "3s"
			tls := &TLS{
				Enabled: false,
			}
			net := &Net{
				DNS: &DNS{
					CacheEnabled:    true,
					RefreshDuration: "1h",
					CacheExpiration: "24h",
				},
				Dialer: &Dialer{
					Timeout:          "5s",
					Keepalive:        "5m",
					DualStackEnabled: false,
				},
				TLS: tls,
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            false,
					IPRecoverDestinationAddr: false,
				},
			}
			return test{
				name: "return Redis when parameters are set",
				fields: fields{
					Addrs:                addrs,
					DB:                   db,
					DialTimeout:          dialTimeout,
					IdleCheckFrequency:   idleCheckFrequency,
					IdleTimeout:          IdleTimeout,
					InitialPingDuration:  initialPingDuration,
					InitialPingTimeLimit: initialPingTimelimit,
					KVPrefix:             kvPrefix,
					KeyPref:              keyPref,
					MaxConnAge:           maxConnAge,
					MaxRedirects:         maxRedirects,
					MaxRetries:           maxRetries,
					MaxRetryBackoff:      maxRetryBackoff,
					MinIdleConns:         minIdleConns,
					MinRetryBackoff:      minRetryBackoff,
					Network:              network,
					Password:             password,
					PoolSize:             poolSize,
					PoolTimeout:          poolTimeout,
					PrefixDelimiter:      prefixDelimiter,
					ReadOnly:             readOnly,
					ReadTimeout:          readTimeout,
					RouteByLatency:       routeByLatency,
					RouteRandomly:        routeRandomly,
					SentinelPassword:     sentinelPassword,
					SentinelMasterName:   sentinelMasterName,
					Net:                  net,
					TLS:                  tls,
					Username:             username,
					VKPrefix:             vkPrefix,
					WriteTimeout:         writeTimeout,
				},
				want: want{
					want: &Redis{
						Addrs:                addrs,
						DB:                   db,
						DialTimeout:          dialTimeout,
						IdleCheckFrequency:   idleCheckFrequency,
						IdleTimeout:          IdleTimeout,
						InitialPingDuration:  initialPingDuration,
						InitialPingTimeLimit: initialPingTimelimit,
						KVPrefix:             kvPrefix,
						KeyPref:              keyPref,
						MaxConnAge:           maxConnAge,
						MaxRedirects:         maxRedirects,
						MaxRetries:           maxRetries,
						MaxRetryBackoff:      maxRetryBackoff,
						MinIdleConns:         minIdleConns,
						MinRetryBackoff:      minRetryBackoff,
						Network:              network,
						Password:             password,
						PoolSize:             poolSize,
						PoolTimeout:          poolTimeout,
						PrefixDelimiter:      prefixDelimiter,
						ReadOnly:             readOnly,
						ReadTimeout:          readTimeout,
						RouteByLatency:       routeByLatency,
						RouteRandomly:        routeRandomly,
						SentinelPassword:     sentinelPassword,
						SentinelMasterName:   sentinelMasterName,
						Net:                  net,
						TLS:                  tls,
						Username:             username,
						VKPrefix:             vkPrefix,
						WriteTimeout:         writeTimeout,
					},
				},
			}
		}(),
		func() test {
			envPrefix := "REDIS_BIND_"
			p := map[string]string{
				envPrefix + "ADDRS":                "redis.default.svc.cluster.local:6379",
				envPrefix + "DIAL_TIMEOUT":         "5s",
				envPrefix + "IDLE_CHECK_FREQUENCY": "1m",
				envPrefix + "IDLE_TIMEOUT":         "5m",
				envPrefix + "KEY_PREF":             "vald",
				envPrefix + "MAX_CONN_AGE":         "0s",
				envPrefix + "MAX_RETRY_BACKOFF":    "512s",
				envPrefix + "MIN_RETRY_BACKOFF":    "8ms",
				envPrefix + "NETWORK":              "tcp",
				envPrefix + "PASSWORD":             "password",
				envPrefix + "POOL_TIMEOUT":         "4s",
				envPrefix + "PREFIX_DELIMITER":     "_",
				envPrefix + "READ_TIMEOUT":         "3s",
				envPrefix + "SENTINEL_PASSWORD":    "",
				envPrefix + "SENTINEL_MASTER_NAME": "",
				envPrefix + "KV_PREFIX":            "",
				envPrefix + "VK_PREFIX":            "",
				envPrefix + "USERNAME":             "vald",
				envPrefix + "WRITE_TIMEOUT":        "3s",
			}
			db := 0
			maxRedirects := 3
			maxRetries := 1
			minIdleConns := 0
			poolSize := 10
			readOnly := true
			routeByLatency := false
			routeRandomly := true
			tls := &TLS{
				Enabled: false,
			}
			net := &Net{
				DNS: &DNS{
					CacheEnabled:    true,
					RefreshDuration: "1h",
					CacheExpiration: "24h",
				},
				Dialer: &Dialer{
					Timeout:          "5s",
					Keepalive:        "5m",
					DualStackEnabled: false,
				},
				TLS: tls,
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            false,
					IPRecoverDestinationAddr: false,
				},
			}
			return test{
				name: "return Redis when parameters are set as environment value",
				fields: fields{
					Addrs:                []string{"_" + envPrefix + "ADDRS_"},
					DB:                   db,
					DialTimeout:          "_" + envPrefix + "DIAL_TIMEOUT_",
					IdleCheckFrequency:   "_" + envPrefix + "IDLE_CHECK_FREQUENCY_",
					IdleTimeout:          "_" + envPrefix + "IDLE_TIMEOUT_",
					InitialPingDuration:  "",
					InitialPingTimeLimit: "",
					KVPrefix:             "_" + envPrefix + "KV_PREFIX_",
					KeyPref:              "_" + envPrefix + "KEY_PREF_",
					MaxConnAge:           "_" + envPrefix + "MAX_CONN_AGE_",
					MaxRedirects:         maxRedirects,
					MaxRetries:           maxRetries,
					MaxRetryBackoff:      "_" + envPrefix + "MAX_RETRY_BACKOFF_",
					MinIdleConns:         minIdleConns,
					MinRetryBackoff:      "_" + envPrefix + "MIN_RETRY_BACKOFF_",
					Network:              "_" + envPrefix + "NETWORK_",
					Password:             "_" + envPrefix + "PASSWORD_",
					PoolSize:             poolSize,
					PoolTimeout:          "_" + envPrefix + "POOL_TIMEOUT_",
					PrefixDelimiter:      "_" + envPrefix + "PREFIX_DELIMITER_",
					ReadOnly:             readOnly,
					ReadTimeout:          "_" + envPrefix + "READ_TIMEOUT_",
					RouteByLatency:       routeByLatency,
					RouteRandomly:        routeRandomly,
					SentinelPassword:     "_" + envPrefix + "SENTINEL_PASSWORD_",
					SentinelMasterName:   "_" + envPrefix + "SENTINEL_MASTER_NAME_",
					Net:                  net,
					TLS:                  tls,
					Username:             "_" + envPrefix + "USERNAME_",
					VKPrefix:             "_" + envPrefix + "VK_PREFIX_",
					WriteTimeout:         "_" + envPrefix + "WRITE_TIMEOUT_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &Redis{
						Addrs: []string{
							"redis.default.svc.cluster.local:6379",
						},
						DB:                   db,
						DialTimeout:          "5s",
						IdleCheckFrequency:   "1m",
						IdleTimeout:          "5m",
						InitialPingDuration:  "",
						InitialPingTimeLimit: "",
						KVPrefix:             "",
						KeyPref:              "vald",
						MaxConnAge:           "0s",
						MaxRedirects:         maxRedirects,
						MaxRetries:           maxRetries,
						MaxRetryBackoff:      "512s",
						MinIdleConns:         minIdleConns,
						MinRetryBackoff:      "8ms",
						Network:              "tcp",
						Password:             "password",
						PoolSize:             poolSize,
						PoolTimeout:          "4s",
						PrefixDelimiter:      "_",
						ReadOnly:             readOnly,
						ReadTimeout:          "3s",
						RouteByLatency:       routeByLatency,
						RouteRandomly:        routeRandomly,
						SentinelPassword:     "",
						SentinelMasterName:   "",
						Net:                  net,
						TLS:                  tls,
						Username:             "vald",
						VKPrefix:             "",
						WriteTimeout:         "3s",
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Redis when parameters are not set",
				fields: fields{},
				want: want{
					want: &Redis{
						Net: &Net{},
						TLS: &TLS{},
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
			r := &Redis{
				Addrs:                test.fields.Addrs,
				DB:                   test.fields.DB,
				DialTimeout:          test.fields.DialTimeout,
				IdleCheckFrequency:   test.fields.IdleCheckFrequency,
				IdleTimeout:          test.fields.IdleTimeout,
				InitialPingDuration:  test.fields.InitialPingDuration,
				InitialPingTimeLimit: test.fields.InitialPingTimeLimit,
				KVPrefix:             test.fields.KVPrefix,
				KeyPref:              test.fields.KeyPref,
				MaxConnAge:           test.fields.MaxConnAge,
				MaxRedirects:         test.fields.MaxRedirects,
				MaxRetries:           test.fields.MaxRetries,
				MaxRetryBackoff:      test.fields.MaxRetryBackoff,
				MinIdleConns:         test.fields.MinIdleConns,
				MinRetryBackoff:      test.fields.MinRetryBackoff,
				Network:              test.fields.Network,
				Password:             test.fields.Password,
				PoolSize:             test.fields.PoolSize,
				PoolTimeout:          test.fields.PoolTimeout,
				PrefixDelimiter:      test.fields.PrefixDelimiter,
				ReadOnly:             test.fields.ReadOnly,
				ReadTimeout:          test.fields.ReadTimeout,
				RouteByLatency:       test.fields.RouteByLatency,
				RouteRandomly:        test.fields.RouteRandomly,
				SentinelPassword:     test.fields.SentinelPassword,
				SentinelMasterName:   test.fields.SentinelMasterName,
				Net:                  test.fields.Net,
				TLS:                  test.fields.TLS,
				Username:             test.fields.Username,
				VKPrefix:             test.fields.VKPrefix,
				WriteTimeout:         test.fields.WriteTimeout,
			}

			got := r.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRedis_Opts(t *testing.T) {
	type fields struct {
		Addrs                []string
		DB                   int
		DialTimeout          string
		IdleCheckFrequency   string
		IdleTimeout          string
		InitialPingDuration  string
		InitialPingTimeLimit string
		KVPrefix             string
		KeyPref              string
		MaxConnAge           string
		MaxRedirects         int
		MaxRetries           int
		MaxRetryBackoff      string
		MinIdleConns         int
		MinRetryBackoff      string
		Network              string
		Password             string
		PoolSize             int
		PoolTimeout          string
		PrefixDelimiter      string
		ReadOnly             bool
		ReadTimeout          string
		RouteByLatency       bool
		RouteRandomly        bool
		SentinelPassword     string
		SentinelMasterName   string
		Net                  *Net
		TLS                  *TLS
		Username             string
		VKPrefix             string
		WriteTimeout         string
	}
	type want struct {
		wantOpts []redis.Option
		err      error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []redis.Option, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotOpts []redis.Option, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(len(gotOpts), len(w.wantOpts)) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOpts, w.wantOpts)
		}
		return nil
	}
	tests := []test{
		{
			name: "return 26 []redis.Options and nil error when all parameters are set",
			fields: fields{
				Addrs: []string{
					"redis-01.default.svc.cluster.local:6379",
					"redis-02.default.svc.cluster.local:6379",
				},
				DB:                   0,
				DialTimeout:          "5s",
				IdleCheckFrequency:   "1m",
				IdleTimeout:          "5m",
				InitialPingDuration:  "",
				InitialPingTimeLimit: "",
				KVPrefix:             "",
				KeyPref:              "vald",
				MaxConnAge:           "0s",
				MaxRedirects:         3,
				MaxRetries:           1,
				MaxRetryBackoff:      "512s",
				MinIdleConns:         0,
				MinRetryBackoff:      "8ms",
				Network:              "tcp",
				Password:             "password",
				PoolSize:             10,
				PoolTimeout:          "4s",
				PrefixDelimiter:      "_",
				ReadOnly:             true,
				ReadTimeout:          "3s",
				RouteByLatency:       false,
				RouteRandomly:        true,
				SentinelPassword:     "",
				SentinelMasterName:   "",
				Net: &Net{
					DNS: &DNS{
						CacheEnabled:    true,
						RefreshDuration: "1h",
						CacheExpiration: "24h",
					},
					Dialer: &Dialer{
						Timeout:          "5s",
						Keepalive:        "5m",
						DualStackEnabled: false,
					},
					TLS: &TLS{
						Enabled: false,
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
				},
				TLS: &TLS{
					Enabled: true,
					Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
					Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
					CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
				},
				Username:     "vald",
				VKPrefix:     "",
				WriteTimeout: "3s",
			},
			want: want{
				wantOpts: make([]redis.Option, 27),
			},
		},
		{
			name: "return 26 []redis.Options and nil error when addrs is nil",
			fields: fields{
				DB:                   0,
				DialTimeout:          "5s",
				IdleCheckFrequency:   "1m",
				IdleTimeout:          "5m",
				InitialPingDuration:  "",
				InitialPingTimeLimit: "",
				KVPrefix:             "",
				KeyPref:              "vald",
				MaxConnAge:           "0s",
				MaxRedirects:         3,
				MaxRetries:           1,
				MaxRetryBackoff:      "512s",
				MinIdleConns:         0,
				MinRetryBackoff:      "8ms",
				Network:              "tcp",
				Password:             "password",
				PoolSize:             10,
				PoolTimeout:          "4s",
				PrefixDelimiter:      "_",
				ReadOnly:             true,
				ReadTimeout:          "3s",
				RouteByLatency:       false,
				RouteRandomly:        true,
				SentinelPassword:     "",
				SentinelMasterName:   "",
				Net: &Net{
					DNS: &DNS{
						CacheEnabled:    true,
						RefreshDuration: "1h",
						CacheExpiration: "24h",
					},
					Dialer: &Dialer{
						Timeout:          "5s",
						Keepalive:        "5m",
						DualStackEnabled: false,
					},
					TLS: &TLS{
						Enabled: false,
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
				},
				TLS: &TLS{
					Enabled: false,
				},
				Username:     "vald",
				VKPrefix:     "",
				WriteTimeout: "3s",
			},
			want: want{
				wantOpts: make([]redis.Option, 26),
			},
		},
		{
			name: "return 26 []redis.Options and nil error when Network is empty",
			fields: fields{
				Addrs: []string{
					"redis.default.svc.cluster.local:6379",
				},
				DB:                   0,
				DialTimeout:          "5s",
				IdleCheckFrequency:   "1m",
				IdleTimeout:          "5m",
				InitialPingDuration:  "",
				InitialPingTimeLimit: "",
				KVPrefix:             "",
				KeyPref:              "vald",
				MaxConnAge:           "0s",
				MaxRedirects:         3,
				MaxRetries:           1,
				MaxRetryBackoff:      "512s",
				MinIdleConns:         0,
				MinRetryBackoff:      "8ms",
				Password:             "password",
				PoolSize:             10,
				PoolTimeout:          "4s",
				PrefixDelimiter:      "_",
				ReadOnly:             true,
				ReadTimeout:          "3s",
				RouteByLatency:       false,
				RouteRandomly:        true,
				SentinelPassword:     "",
				SentinelMasterName:   "",
				Net: &Net{
					DNS: &DNS{
						CacheEnabled:    true,
						RefreshDuration: "1h",
						CacheExpiration: "24h",
					},
					Dialer: &Dialer{
						Timeout:          "5s",
						Keepalive:        "5m",
						DualStackEnabled: false,
					},
					TLS: &TLS{
						Enabled: false,
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
				},
				TLS: &TLS{
					Enabled: false,
				},
				Username:     "vald",
				VKPrefix:     "",
				WriteTimeout: "3s",
			},
			want: want{
				wantOpts: make([]redis.Option, 26),
			},
		},
		{
			name: "return nil []redis.Options and error when Net.TLS has invalid parameter",
			fields: fields{
				Addrs: []string{
					"redis.default.svc.cluster.local:6379",
				},
				DB:                   0,
				DialTimeout:          "5s",
				IdleCheckFrequency:   "1m",
				IdleTimeout:          "5m",
				InitialPingDuration:  "",
				InitialPingTimeLimit: "",
				KVPrefix:             "",
				KeyPref:              "vald",
				MaxConnAge:           "0s",
				MaxRedirects:         3,
				MaxRetries:           1,
				MaxRetryBackoff:      "512s",
				MinIdleConns:         0,
				MinRetryBackoff:      "8ms",
				Network:              "tcp",
				Password:             "password",
				PoolSize:             10,
				PoolTimeout:          "4s",
				PrefixDelimiter:      "_",
				ReadOnly:             true,
				ReadTimeout:          "3s",
				RouteByLatency:       false,
				RouteRandomly:        true,
				SentinelPassword:     "",
				SentinelMasterName:   "",
				Net: &Net{
					DNS: &DNS{
						CacheEnabled:    true,
						RefreshDuration: "1h",
						CacheExpiration: "24h",
					},
					Dialer: &Dialer{
						Timeout:          "5s",
						Keepalive:        "5m",
						DualStackEnabled: false,
					},
					TLS: &TLS{
						Enabled: true,
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
				},
				TLS: &TLS{
					Enabled: false,
				},
				Username:     "vald",
				VKPrefix:     "",
				WriteTimeout: "3s",
			},
			want: want{
				wantOpts: nil,
				err:      errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			name: "return nil []redis.Options and error when TLS has invalid parameter",
			fields: fields{
				Addrs: []string{
					"redis.default.svc.cluster.local:6379",
				},
				DB:                   0,
				DialTimeout:          "5s",
				IdleCheckFrequency:   "1m",
				IdleTimeout:          "5m",
				InitialPingDuration:  "",
				InitialPingTimeLimit: "",
				KVPrefix:             "",
				KeyPref:              "vald",
				MaxConnAge:           "0s",
				MaxRedirects:         3,
				MaxRetries:           1,
				MaxRetryBackoff:      "512s",
				MinIdleConns:         0,
				MinRetryBackoff:      "8ms",
				Network:              "tcp",
				Password:             "password",
				PoolSize:             10,
				PoolTimeout:          "4s",
				PrefixDelimiter:      "_",
				ReadOnly:             true,
				ReadTimeout:          "3s",
				RouteByLatency:       false,
				RouteRandomly:        true,
				SentinelPassword:     "",
				SentinelMasterName:   "",
				Net: &Net{
					DNS: &DNS{
						CacheEnabled:    true,
						RefreshDuration: "1h",
						CacheExpiration: "24h",
					},
					Dialer: &Dialer{
						Timeout:          "5s",
						Keepalive:        "5m",
						DualStackEnabled: false,
					},
					TLS: &TLS{
						Enabled: false,
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
				},
				TLS: &TLS{
					Enabled: true,
				},
				Username:     "vald",
				VKPrefix:     "",
				WriteTimeout: "3s",
			},
			want: want{
				wantOpts: nil,
				err:      errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			name: "return nil []redis.Options and error when Dialer has invalid parameter",
			fields: fields{
				Addrs: []string{
					"redis.default.svc.cluster.local:6379",
				},
				DB:                   0,
				DialTimeout:          "5s",
				IdleCheckFrequency:   "1m",
				IdleTimeout:          "5m",
				InitialPingDuration:  "",
				InitialPingTimeLimit: "",
				KVPrefix:             "",
				KeyPref:              "vald",
				MaxConnAge:           "0s",
				MaxRedirects:         3,
				MaxRetries:           1,
				MaxRetryBackoff:      "512s",
				MinIdleConns:         0,
				MinRetryBackoff:      "8ms",
				Network:              "tcp",
				Password:             "password",
				PoolSize:             10,
				PoolTimeout:          "4s",
				PrefixDelimiter:      "_",
				ReadOnly:             true,
				ReadTimeout:          "3s",
				RouteByLatency:       false,
				RouteRandomly:        true,
				SentinelPassword:     "",
				SentinelMasterName:   "",
				Net: &Net{
					DNS: &DNS{
						CacheEnabled:    true,
						RefreshDuration: "12h",
						CacheExpiration: "1h",
					},
					Dialer: &Dialer{
						Timeout:          "5s",
						Keepalive:        "5m",
						DualStackEnabled: false,
					},
					TLS: &TLS{
						Enabled: false,
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
				},
				TLS: &TLS{
					Enabled: false,
				},
				Username:     "vald",
				VKPrefix:     "",
				WriteTimeout: "3s",
			},
			want: want{
				wantOpts: nil,
				err:      errors.ErrInvalidDNSConfig(12*time.Hour, time.Hour),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			r := &Redis{
				Addrs:                test.fields.Addrs,
				DB:                   test.fields.DB,
				DialTimeout:          test.fields.DialTimeout,
				IdleCheckFrequency:   test.fields.IdleCheckFrequency,
				IdleTimeout:          test.fields.IdleTimeout,
				InitialPingDuration:  test.fields.InitialPingDuration,
				InitialPingTimeLimit: test.fields.InitialPingTimeLimit,
				KVPrefix:             test.fields.KVPrefix,
				KeyPref:              test.fields.KeyPref,
				MaxConnAge:           test.fields.MaxConnAge,
				MaxRedirects:         test.fields.MaxRedirects,
				MaxRetries:           test.fields.MaxRetries,
				MaxRetryBackoff:      test.fields.MaxRetryBackoff,
				MinIdleConns:         test.fields.MinIdleConns,
				MinRetryBackoff:      test.fields.MinRetryBackoff,
				Network:              test.fields.Network,
				Password:             test.fields.Password,
				PoolSize:             test.fields.PoolSize,
				PoolTimeout:          test.fields.PoolTimeout,
				PrefixDelimiter:      test.fields.PrefixDelimiter,
				ReadOnly:             test.fields.ReadOnly,
				ReadTimeout:          test.fields.ReadTimeout,
				RouteByLatency:       test.fields.RouteByLatency,
				RouteRandomly:        test.fields.RouteRandomly,
				SentinelPassword:     test.fields.SentinelPassword,
				SentinelMasterName:   test.fields.SentinelMasterName,
				Net:                  test.fields.Net,
				TLS:                  test.fields.TLS,
				Username:             test.fields.Username,
				VKPrefix:             test.fields.VKPrefix,
				WriteTimeout:         test.fields.WriteTimeout,
			}

			gotOpts, err := r.Opts()
			if err := checkFunc(test.want, gotOpts, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
