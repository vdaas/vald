// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
// 		wantC Conn
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
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
// 			gotC, err := New(test.args.ctx, test.args.opts...)
// 			if err := checkFunc(test.want, gotC, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_init(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 	}
// 	type want struct {
// 		want []atomic.Pointer[poolConn]
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []atomic.Pointer[poolConn]) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []atomic.Pointer[poolConn]) error {
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 	}
// 	type want struct {
// 		wantRidx uint64
// 		wantPc   *poolConn
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64, *poolConn) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRidx uint64, gotPc *poolConn) error {
// 		if !reflect.DeepEqual(gotRidx, w.wantRidx) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRidx, w.wantRidx)
// 		}
// 		if !reflect.DeepEqual(gotPc, w.wantPc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPc, w.wantPc)
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 			}
//
// 			gotRidx, gotPc := p.load(test.args.idx)
// 			if err := checkFunc(test.want, gotRidx, gotPc); err != nil {
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 	}
// 	type want struct {
// 		wantPc *poolConn
// 		wantOk bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *poolConn, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPc *poolConn, gotOk bool) error {
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 			}
//
// 			gotPc, gotOk := p.getHealthyConn(test.args.ctx)
// 			if err := checkFunc(test.want, gotPc, gotOk); err != nil {
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
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
// 		in0 context.Context
// 	}
// 	type want struct {
// 		want map[string]uint64
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]uint64) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]uint64) error {
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
// 		           in0:nil,
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
// 		           in0:nil,
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
// 			got := Metrics(test.args.in0)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_isHealthy(t *testing.T) {
// 	type args struct {
// 		idx  uint64
// 		conn *ClientConn
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 	}
// 	type want struct {
// 		wantState   connectivity.State
// 		wantHealthy bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, connectivity.State, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotState connectivity.State, gotHealthy bool) error {
// 		if !reflect.DeepEqual(gotState, w.wantState) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotState, w.wantState)
// 		}
// 		if !reflect.DeepEqual(gotHealthy, w.wantHealthy) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHealthy, w.wantHealthy)
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
// 		           conn:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 		           conn:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
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
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 			}
//
// 			gotState, gotHealthy := p.isHealthy(test.args.idx, test.args.conn)
// 			if err := checkFunc(test.want, gotState, gotHealthy); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
