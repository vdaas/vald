// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// NOT IMPLEMENTED BELOW
//
// func Test_poolConn_Close(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		delay time.Duration
// 	}
// 	type fields struct {
// 		conn *ClientConn
// 		addr string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           delay:nil,
// 		       },
// 		       fields: fields {
// 		           conn:nil,
// 		           addr:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           delay:nil,
// 		           },
// 		           fields: fields {
// 		           conn:nil,
// 		           addr:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			pc := &poolConn{
// 				conn: test.fields.conn,
// 				addr: test.fields.addr,
// 			}
//
// 			err := pc.Close(test.args.ctx, test.args.delay)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		opts []Option
// 	}
// 	type want struct {
// 		want Conn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.ctx, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_init(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			p.init()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_getSlots(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want *[]atomic.Pointer[poolConn]
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *[]atomic.Pointer[poolConn]) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *[]atomic.Pointer[poolConn]) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.getSlots()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_grow(t *testing.T) {
// 	type args struct {
// 		newSize uint64
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           newSize:0,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           newSize:0,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			p.grow(test.args.newSize)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_load(t *testing.T) {
// 	type args struct {
// 		idx uint64
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want *poolConn
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *poolConn) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *poolConn) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           idx:0,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           idx:0,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.load(test.args.idx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_store(t *testing.T) {
// 	type args struct {
// 		idx uint64
// 		pc  *poolConn
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           idx:0,
// 		           pc:poolConn{},
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           idx:0,
// 		           pc:poolConn{},
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			p.store(test.args.idx, test.args.pc)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_loop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		fn  func(ctx context.Context, idx uint64, pc *poolConn) bool
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			err := p.loop(test.args.ctx, test.args.fn)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_slotCount(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.slotCount()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_flush(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			p.flush()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_refreshConn(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		idx  uint64
// 		pc   *poolConn
// 		addr string
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           pc:poolConn{},
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           pc:poolConn{},
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			err := p.refreshConn(test.args.ctx, test.args.idx, test.args.pc, test.args.addr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Connect(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want Conn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got, err := p.Connect(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_connect(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		ips []string
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		wantC Conn
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotC Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotC, w.wantC) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           ips:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           ips:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			gotC, err := p.connect(test.args.ctx, test.args.ips...)
// 			if err := checkFunc(test.want, gotC, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_singleTargetConnect(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want Conn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got, err := p.singleTargetConnect(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Reconnect(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		force bool
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want Conn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           force:false,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           force:false,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got, err := p.Reconnect(test.args.ctx, test.args.force)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Disconnect(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			err := p.Disconnect(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_dial(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		idx  uint64
// 		addr string
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want *ClientConn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *ClientConn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *ClientConn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got, err := p.dial(test.args.ctx, test.args.idx, test.args.addr)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_getHealthyConn(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		wantIdx uint64
// 		wantPc  *poolConn
// 		wantOk  bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64, *poolConn, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotIdx uint64, gotPc *poolConn, gotOk bool) error {
// 		if !reflect.DeepEqual(gotIdx, w.wantIdx) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIdx, w.wantIdx)
// 		}
// 		if !reflect.DeepEqual(gotPc, w.wantPc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPc, w.wantPc)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			gotIdx, gotPc, gotOk := p.getHealthyConn(test.args.ctx)
// 			if err := checkFunc(test.want, gotIdx, gotPc, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Do(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(conn *ClientConn) error
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			err := p.Do(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Get(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		wantConn *ClientConn
// 		wantOk   bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *ClientConn, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotConn *ClientConn, gotOk bool) error {
// 		if !reflect.DeepEqual(gotConn, w.wantConn) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			gotConn, gotOk := p.Get(test.args.ctx)
// 			if err := checkFunc(test.want, gotConn, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_IsHealthy(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.IsHealthy(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Len(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.Len()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Size(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.Size()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_IsIPConn(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.IsIPConn()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_String(t *testing.T) {
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.String()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_lookupIPAddr(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got, err := p.lookupIPAddr(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_scanGRPCPort(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		wantPort uint16
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint16, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPort uint16, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotPort, w.wantPort) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPort, w.wantPort)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			gotPort, err := p.scanGRPCPort(test.args.ctx)
// 			if err := checkFunc(test.want, gotPort, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestMetrics(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type want struct {
// 		want map[string]int64
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]int64) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]int64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := Metrics(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_isHealthy(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		idx  uint64
// 		conn *ClientConn
// 	}
// 	type fields struct {
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		startPort         uint16
// 		endPort           uint16
// 		host              string
// 		port              uint16
// 		addr              string
// 		isIPAddr          bool
// 		enableDNSLookup   bool
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialOpts          []DialOption
// 		dialTimeout       time.Duration
// 		oldConnCloseDelay time.Duration
// 		bo                backoff.Backoff
// 		errGroup          errgroup.Group
// 		dnsHash           atomic.Pointer[string]
// 		closing           atomic.Bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           conn:nil,
// 		       },
// 		       fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           conn:nil,
// 		           },
// 		           fields: fields {
// 		           connSlots:nil,
// 		           startPort:0,
// 		           endPort:0,
// 		           host:"",
// 		           port:0,
// 		           addr:"",
// 		           isIPAddr:false,
// 		           enableDNSLookup:false,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialOpts:nil,
// 		           dialTimeout:nil,
// 		           oldConnCloseDelay:nil,
// 		           bo:nil,
// 		           errGroup:nil,
// 		           dnsHash:nil,
// 		           closing:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				connSlots:         test.fields.connSlots,
// 				startPort:         test.fields.startPort,
// 				endPort:           test.fields.endPort,
// 				host:              test.fields.host,
// 				port:              test.fields.port,
// 				addr:              test.fields.addr,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialOpts:          test.fields.dialOpts,
// 				dialTimeout:       test.fields.dialTimeout,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				bo:                test.fields.bo,
// 				errGroup:          test.fields.errGroup,
// 				dnsHash:           test.fields.dnsHash,
// 				closing:           test.fields.closing,
// 			}
//
// 			got := p.isHealthy(test.args.ctx, test.args.idx, test.args.conn)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
