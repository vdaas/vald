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

	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestRedis_Bind(t *testing.T) {
	t.Parallel()
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
		TCP                  *TCP
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Redis) error {
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
		           Addrs: nil,
		           DB: 0,
		           DialTimeout: "",
		           IdleCheckFrequency: "",
		           IdleTimeout: "",
		           InitialPingDuration: "",
		           InitialPingTimeLimit: "",
		           KVPrefix: "",
		           KeyPref: "",
		           MaxConnAge: "",
		           MaxRedirects: 0,
		           MaxRetries: 0,
		           MaxRetryBackoff: "",
		           MinIdleConns: 0,
		           MinRetryBackoff: "",
		           Network: "",
		           Password: "",
		           PoolSize: 0,
		           PoolTimeout: "",
		           PrefixDelimiter: "",
		           ReadOnly: false,
		           ReadTimeout: "",
		           RouteByLatency: false,
		           RouteRandomly: false,
		           TCP: TCP{},
		           TLS: TLS{},
		           Username: "",
		           VKPrefix: "",
		           WriteTimeout: "",
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
		           Addrs: nil,
		           DB: 0,
		           DialTimeout: "",
		           IdleCheckFrequency: "",
		           IdleTimeout: "",
		           InitialPingDuration: "",
		           InitialPingTimeLimit: "",
		           KVPrefix: "",
		           KeyPref: "",
		           MaxConnAge: "",
		           MaxRedirects: 0,
		           MaxRetries: 0,
		           MaxRetryBackoff: "",
		           MinIdleConns: 0,
		           MinRetryBackoff: "",
		           Network: "",
		           Password: "",
		           PoolSize: 0,
		           PoolTimeout: "",
		           PrefixDelimiter: "",
		           ReadOnly: false,
		           ReadTimeout: "",
		           RouteByLatency: false,
		           RouteRandomly: false,
		           TCP: TCP{},
		           TLS: TLS{},
		           Username: "",
		           VKPrefix: "",
		           WriteTimeout: "",
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
				TCP:                  test.fields.TCP,
				TLS:                  test.fields.TLS,
				Username:             test.fields.Username,
				VKPrefix:             test.fields.VKPrefix,
				WriteTimeout:         test.fields.WriteTimeout,
			}

			got := r.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRedis_Opts(t *testing.T) {
	t.Parallel()
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
		TCP                  *TCP
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
		           Addrs: nil,
		           DB: 0,
		           DialTimeout: "",
		           IdleCheckFrequency: "",
		           IdleTimeout: "",
		           InitialPingDuration: "",
		           InitialPingTimeLimit: "",
		           KVPrefix: "",
		           KeyPref: "",
		           MaxConnAge: "",
		           MaxRedirects: 0,
		           MaxRetries: 0,
		           MaxRetryBackoff: "",
		           MinIdleConns: 0,
		           MinRetryBackoff: "",
		           Network: "",
		           Password: "",
		           PoolSize: 0,
		           PoolTimeout: "",
		           PrefixDelimiter: "",
		           ReadOnly: false,
		           ReadTimeout: "",
		           RouteByLatency: false,
		           RouteRandomly: false,
		           TCP: TCP{},
		           TLS: TLS{},
		           Username: "",
		           VKPrefix: "",
		           WriteTimeout: "",
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
		           Addrs: nil,
		           DB: 0,
		           DialTimeout: "",
		           IdleCheckFrequency: "",
		           IdleTimeout: "",
		           InitialPingDuration: "",
		           InitialPingTimeLimit: "",
		           KVPrefix: "",
		           KeyPref: "",
		           MaxConnAge: "",
		           MaxRedirects: 0,
		           MaxRetries: 0,
		           MaxRetryBackoff: "",
		           MinIdleConns: 0,
		           MinRetryBackoff: "",
		           Network: "",
		           Password: "",
		           PoolSize: 0,
		           PoolTimeout: "",
		           PrefixDelimiter: "",
		           ReadOnly: false,
		           ReadTimeout: "",
		           RouteByLatency: false,
		           RouteRandomly: false,
		           TCP: TCP{},
		           TLS: TLS{},
		           Username: "",
		           VKPrefix: "",
		           WriteTimeout: "",
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
				TCP:                  test.fields.TCP,
				TLS:                  test.fields.TLS,
				Username:             test.fields.Username,
				VKPrefix:             test.fields.VKPrefix,
				WriteTimeout:         test.fields.WriteTimeout,
			}

			gotOpts, err := r.Opts()
			if err := test.checkFunc(test.want, gotOpts, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
