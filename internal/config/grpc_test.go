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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Addrs: nil,
		           HealthCheckDuration: "",
		           ConnectionPool: ConnectionPool{},
		           Backoff: Backoff{},
		           CallOption: CallOption{},
		           DialOption: DialOption{},
		           TLS: TLS{},
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
		           HealthCheckDuration: "",
		           ConnectionPool: ConnectionPool{},
		           Backoff: Backoff{},
		           CallOption: CallOption{},
		           DialOption: DialOption{},
		           TLS: TLS{},
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *GRPCClientKeepalive) error {
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
		           Time: "",
		           Timeout: "",
		           PermitWithoutStream: false,
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
		           Time: "",
		           Timeout: "",
		           PermitWithoutStream: false,
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           WaitForReady: false,
		           MaxRetryRPCBufferSize: 0,
		           MaxRecvMsgSize: 0,
		           MaxSendMsgSize: 0,
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
		           WaitForReady: false,
		           MaxRetryRPCBufferSize: 0,
		           MaxRecvMsgSize: 0,
		           MaxSendMsgSize: 0,
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
		TCP                         *TCP
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *DialOption) error {
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
		           WriteBufferSize: 0,
		           ReadBufferSize: 0,
		           InitialWindowSize: 0,
		           InitialConnectionWindowSize: 0,
		           MaxMsgSize: 0,
		           BackoffMaxDelay: "",
		           BackoffBaseDelay: "",
		           BackoffJitter: 0,
		           BackoffMultiplier: 0,
		           MinimumConnectionTimeout: "",
		           EnableBackoff: false,
		           Insecure: false,
		           Timeout: "",
		           TCP: TCP{},
		           KeepAlive: GRPCClientKeepalive{},
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
		           WriteBufferSize: 0,
		           ReadBufferSize: 0,
		           InitialWindowSize: 0,
		           InitialConnectionWindowSize: 0,
		           MaxMsgSize: 0,
		           BackoffMaxDelay: "",
		           BackoffBaseDelay: "",
		           BackoffJitter: 0,
		           BackoffMultiplier: 0,
		           MinimumConnectionTimeout: "",
		           EnableBackoff: false,
		           Insecure: false,
		           Timeout: "",
		           TCP: TCP{},
		           KeepAlive: GRPCClientKeepalive{},
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
				TCP:                         test.fields.TCP,
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
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []grpc.Option) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []grpc.Option) error {
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
		           HealthCheckDuration: "",
		           ConnectionPool: ConnectionPool{},
		           Backoff: Backoff{},
		           CallOption: CallOption{},
		           DialOption: DialOption{},
		           TLS: TLS{},
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
		           HealthCheckDuration: "",
		           ConnectionPool: ConnectionPool{},
		           Backoff: Backoff{},
		           CallOption: CallOption{},
		           DialOption: DialOption{},
		           TLS: TLS{},
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
			g := &GRPCClient{
				Addrs:               test.fields.Addrs,
				HealthCheckDuration: test.fields.HealthCheckDuration,
				ConnectionPool:      test.fields.ConnectionPool,
				Backoff:             test.fields.Backoff,
				CallOption:          test.fields.CallOption,
				DialOption:          test.fields.DialOption,
				TLS:                 test.fields.TLS,
			}

			got := g.Opts()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
