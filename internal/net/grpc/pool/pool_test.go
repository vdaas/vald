// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package pool

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

// NOT IMPLEMENTED BELOW

func TestNew(t *testing.T) {
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           opts:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           opts:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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

func Test_pool_init(t *testing.T) {
	type args struct {
		force bool
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           force:false,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           force:false,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			p.init(test.args.force)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_grow(t *testing.T) {
	type args struct {
		size uint64
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           size:0,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           size:0,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			p.grow(test.args.size)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_load(t *testing.T) {
	type args struct {
		idx int
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		wantPc *poolConn
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *poolConn) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotPc *poolConn) error {
		if !reflect.DeepEqual(gotPc, w.wantPc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPc, w.wantPc)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           idx:0,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           idx:0,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotPc := p.load(test.args.idx)
			if err := checkFunc(test.want, gotPc); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_store(t *testing.T) {
	type args struct {
		idx int
		pc  *poolConn
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           idx:0,
		           pc:poolConn{},
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           idx:0,
		           pc:poolConn{},
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			p.store(test.args.idx, test.args.pc)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_loop(t *testing.T) {
	type args struct {
		ctx context.Context
		fn  func(ctx context.Context, idx int, pc *poolConn) bool
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           fn:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           fn:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			err := p.loop(test.args.ctx, test.args.fn)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_len(t *testing.T) {
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got int) error {
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
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got := p.len()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_cap(t *testing.T) {
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got int) error {
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
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got := p.cap()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_flush(t *testing.T) {
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			p.flush()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_refreshConn(t *testing.T) {
	type args struct {
		ctx  context.Context
		idx  int
		pc   *poolConn
		addr string
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           idx:0,
		           pc:poolConn{},
		           addr:"",
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           idx:0,
		           pc:poolConn{},
		           addr:"",
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			err := p.refreshConn(test.args.ctx, test.args.idx, test.args.pc, test.args.addr)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Connect(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
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

func Test_pool_Reconnect(t *testing.T) {
	type args struct {
		ctx   context.Context
		force bool
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           force:false,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           force:false,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
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

func Test_pool_singleTargetConnect(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotC, err := p.singleTargetConnect(test.args.ctx)
			if err := checkFunc(test.want, gotC, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Disconnect(t *testing.T) {
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
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
	type args struct {
		ctx  context.Context
		addr string
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           addr:"",
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           addr:"",
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
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
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		wantHealthy bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotHealthy bool) error {
		if !reflect.DeepEqual(gotHealthy, w.wantHealthy) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHealthy, w.wantHealthy)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotHealthy := p.IsHealthy(test.args.ctx)
			if err := checkFunc(test.want, gotHealthy); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Do(t *testing.T) {
	type args struct {
		ctx context.Context
		f   func(conn *ClientConn) error
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           f:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           f:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			err := p.Do(test.args.ctx, test.args.f)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Get(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got, got1 := p.Get(test.args.ctx)
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_getHealthyConn(t *testing.T) {
	type args struct {
		ctx   context.Context
		cnt   uint64
		retry uint64
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           cnt:0,
		           retry:0,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           cnt:0,
		           retry:0,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			got, got1 := p.getHealthyConn(test.args.ctx, test.args.cnt, test.args.retry)
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Len(t *testing.T) {
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
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
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
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
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
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

func Test_pool_scanGRPCPort(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
				eg:            test.fields.eg,
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

func Test_pool_IsIPConn(t *testing.T) {
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		wantIsIP bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, gotIsIP bool) error {
		if !reflect.DeepEqual(gotIsIP, w.wantIsIP) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIsIP, w.wantIsIP)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotIsIP := p.IsIPConn()
			if err := checkFunc(test.want, gotIsIP); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_String(t *testing.T) {
	type fields struct {
		pool          []atomic.Pointer[poolConn]
		startPort     uint16
		endPort       uint16
		host          string
		port          uint16
		addr          string
		size          atomic.Uint64
		current       atomic.Uint64
		bo            backoff.Backoff
		eg            errgroup.Group
		dopts         []DialOption
		dialTimeout   time.Duration
		roccd         time.Duration
		closing       atomic.Bool
		isIP          bool
		resolveDNS    bool
		reconnectHash atomic.Pointer[string]
	}
	type want struct {
		wantStr string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, gotStr string) error {
		if !reflect.DeepEqual(gotStr, w.wantStr) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStr, w.wantStr)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           pool:nil,
		           startPort:0,
		           endPort:0,
		           host:"",
		           port:0,
		           addr:"",
		           size:nil,
		           current:nil,
		           bo:nil,
		           eg:nil,
		           dopts:nil,
		           dialTimeout:nil,
		           roccd:nil,
		           closing:nil,
		           isIP:false,
		           resolveDNS:false,
		           reconnectHash:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
				eg:            test.fields.eg,
				dopts:         test.fields.dopts,
				dialTimeout:   test.fields.dialTimeout,
				roccd:         test.fields.roccd,
				closing:       test.fields.closing,
				isIP:          test.fields.isIP,
				resolveDNS:    test.fields.resolveDNS,
				reconnectHash: test.fields.reconnectHash,
			}

			gotStr := p.String()
			if err := checkFunc(test.want, gotStr); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_poolConn_Close(t *testing.T) {
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           delay:nil,
		       },
		       fields: fields {
		           conn:nil,
		           addr:"",
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           delay:nil,
		           },
		           fields: fields {
		           conn:nil,
		           addr:"",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			pc := &poolConn{
				conn: test.fields.conn,
				addr: test.fields.addr,
			}

			err := pc.Close(test.args.ctx, test.args.delay)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_isGRPCPort(t *testing.T) {
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           host:"",
		           port:0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           host:"",
		           port:0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           conn:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           conn:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
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
