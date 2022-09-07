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

// Package pool provides grpc connection pool client
package pool

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		opts []Option
	}
	type want struct {
		wantC Conn
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotC Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotC, w.wantC) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			gotC, err := New(test.args.ctx, test.args.opts...)
			if err := checkFunc(test.want, gotC, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Connect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		wantC Conn
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotC Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotC, w.wantC) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotC, err := p.Connect(test.args.ctx)
			if err := checkFunc(test.want, gotC, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_load(t *testing.T) {
	t.Parallel()
	type args struct {
		idx int
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		wantPc *poolConn
		wantOk bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *poolConn, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotPc *poolConn, gotOk bool) error {
		if !reflect.DeepEqual(gotPc, w.wantPc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPc, w.wantPc)
		}
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           idx: 0,
		       },
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           idx: 0,
		           },
		           fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotPc, gotOk := p.load(test.args.idx)
			if err := checkFunc(test.want, gotPc, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_connect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		wantC Conn
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotC Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotC, w.wantC) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotC, err := p.connect(test.args.ctx)
			if err := checkFunc(test.want, gotC, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Disconnect(t *testing.T) {
	t.Parallel()
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
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
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			err := p.Disconnect()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_dial(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		addr string
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		wantConn *ClientConn
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *ClientConn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotConn *ClientConn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotConn, w.wantConn) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
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
		           addr: "",
		       },
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           addr: "",
		           },
		           fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotConn, err := p.dial(test.args.ctx, test.args.addr)
			if err := checkFunc(test.want, gotConn, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_IsHealthy(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		           ctx: nil,
		       },
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got := p.IsHealthy(test.args.ctx)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Do(t *testing.T) {
	t.Parallel()
	type args struct {
		f func(conn *ClientConn) error
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
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
		           f: nil,
		       },
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           f: nil,
		           },
		           fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			err := p.Do(test.args.f)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Get(t *testing.T) {
	t.Parallel()
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		want  *ClientConn
		want1 bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *ClientConn, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *ClientConn, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got, got1 := p.Get()
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_get(t *testing.T) {
	t.Parallel()
	type args struct {
		retry uint64
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		want  *ClientConn
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *ClientConn, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *ClientConn, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           retry: 0,
		       },
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           retry: 0,
		           },
		           fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got, got1 := p.get(test.args.retry)
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Len(t *testing.T) {
	t.Parallel()
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got := p.Len()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Size(t *testing.T) {
	t.Parallel()
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got := p.Size()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_lookupIPAddr(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		wantIps []string
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotIps []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotIps, w.wantIps) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIps, w.wantIps)
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotIps, err := p.lookupIPAddr(test.args.ctx)
			if err := checkFunc(test.want, gotIps, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Reconnect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		force bool
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
	}
	type want struct {
		wantC Conn
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotC Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotC, w.wantC) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
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
		           force: false,
		       },
		       fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           force: false,
		           },
		           fields: fields {
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotC, err := p.Reconnect(test.args.ctx, test.args.force)
			if err := checkFunc(test.want, gotC, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_scanGRPCPort(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Value
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          uint64
		current       uint64
		bo            backoff.Backoff
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Value
		isIP          bool
		resolveDNS    bool
		reconnectHash string
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
		           pool: nil,
		           startPort: 0,
		           endPort: 0,
		           host: "",
		           port: 0,
		           addr: "",
		           size: 0,
		           current: 0,
		           bo: nil,
		           dopts: nil,
		           dialTimeout: nil,
		           roccd: nil,
		           closing: nil,
		           isIP: false,
		           resolveDNS: false,
		           reconnectHash: "",
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
			p := &pool{
				pool:          test.fields.pool,
				startPort:     test.fields.startPort,
				endPort:       test.fields.endPort,
				host:          test.fields.host,
				port:          test.fields.port,
				addr:          test.fields.addr,
				size:          test.fields.size,
				current:       test.fields.current,
				bo:            test.fields.bo,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			err := p.scanGRPCPort(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_isGRPCPort(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		host string
		port uint16
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		           ctx: nil,
		           host: "",
		           port: 0,
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
		           host: "",
		           port: 0,
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

			got := isGRPCPort(test.args.ctx, test.args.host, test.args.port)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_isHealthy(t *testing.T) {
	t.Parallel()
	type args struct {
		conn *ClientConn
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		           conn: nil,
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
		           conn: nil,
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

			got := isHealthy(test.args.conn)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_poolConn_Close(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		delay time.Duration
	}
	type fields struct {
		conn *ClientConn
		addr string
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
		           delay: nil,
		       },
		       fields: fields {
		           conn: nil,
		           addr: "",
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
		           delay: nil,
		           },
		           fields: fields {
		           conn: nil,
		           addr: "",
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
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			pc := &poolConn{
				conn: test.fields.conn,
				addr: test.fields.addr,
			}

			err := pc.Close(test.args.ctx, test.args.delay)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
