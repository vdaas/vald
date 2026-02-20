//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package grpc

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		name string
// 		opts []Option
// 	}
// 	type want struct {
// 		wantC Client
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Client) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotC Client) error {
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
// 		           name:"",
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
// 		           name:"",
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
// 			gotC := New(test.args.name, test.args.opts...)
// 			if err := checkFunc(test.want, gotC); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_StartConnectionMonitor(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			got, err := g.StartConnectionMonitor(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Range(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			err := g.Range(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_RangeConcurrent(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		concurrency int
// 		f           func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           concurrency:0,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           concurrency:0,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			err := g.RangeConcurrent(test.args.ctx, test.args.concurrency, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_OrderedRange(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		orders []string
// 		f      func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           orders:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           orders:nil,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			err := g.OrderedRange(test.args.ctx, test.args.orders, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_OrderedRangeConcurrent(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		orders      []string
// 		concurrency int
// 		f           func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           orders:nil,
// 		           concurrency:0,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           orders:nil,
// 		           concurrency:0,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			err := g.OrderedRangeConcurrent(test.args.ctx, test.args.orders, test.args.concurrency, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestRoundRobin(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		c   Client
// 		f   func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (R, error)
// 	}
// 	type want struct {
// 		wantData R
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, R, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData R, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           c:nil,
// 		           f:nil,
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
// 		           c:nil,
// 		           f:nil,
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
// 			gotData, err := RoundRobin(test.args.ctx, test.args.c, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_RoundRobin(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		wantData any
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			gotData, err := g.RoundRobin(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Do(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 		f    func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		wantData any
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			gotData, err := g.Do(test.args.ctx, test.args.addr, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_executeRPC(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		p    pool.Conn
// 		addr string
// 		f    func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		wantRes       any
// 		wantRetryable bool
// 		err           error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, bool, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes any, gotRetryable bool, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		if !reflect.DeepEqual(gotRetryable, w.wantRetryable) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRetryable, w.wantRetryable)
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			gotRes, gotRetryable, err := g.executeRPC(test.args.ctx, test.args.p, test.args.addr, test.args.f)
// 			if err := checkFunc(test.want, gotRes, gotRetryable, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_do(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		p    pool.Conn
// 		addr string
// 		f    func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		wantData any
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			gotData, err := g.do(test.args.ctx, test.args.p, test.args.addr, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_GetDialOption(t *testing.T) {
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		want []DialOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []DialOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []DialOption) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			got := g.GetDialOption()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_GetCallOption(t *testing.T) {
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		want []CallOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []CallOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []CallOption) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			got := g.GetCallOption()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_GetBackoff(t *testing.T) {
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		want backoff.Backoff
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, backoff.Backoff) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got backoff.Backoff) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			got := g.GetBackoff()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_SetDisableResolveDNSAddr(t *testing.T) {
// 	type args struct {
// 		addr     string
// 		disabled bool
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           addr:"",
// 		           disabled:false,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           addr:"",
// 		           disabled:false,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			g.SetDisableResolveDNSAddr(test.args.addr, test.args.disabled)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Connect(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		addr  string
// 		dopts []DialOption
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		wantConn pool.Conn
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, pool.Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotConn pool.Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotConn, w.wantConn) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
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
// 		           dopts:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dopts:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			gotConn, err := g.Connect(test.args.ctx, test.args.addr, test.args.dopts...)
// 			if err := checkFunc(test.want, gotConn, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_IsConnected(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			got := g.IsConnected(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Disconnect(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			err := g.Disconnect(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_ConnectedAddrs(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 	}
// 	type want struct {
// 		wantAddrs []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotAddrs []string) error {
// 		if !reflect.DeepEqual(gotAddrs, w.wantAddrs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotAddrs, w.wantAddrs)
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
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			gotAddrs := g.ConnectedAddrs(test.args.in0)
// 			if err := checkFunc(test.want, gotAddrs); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Close(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			err := g.Close(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_rangeConns(t *testing.T) {
// 	type args struct {
// 		action string
// 		fn     func(addr string, p pool.Conn) bool
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
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
// 		           action:"",
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 		           action:"",
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 			}
//
// 			err := g.rangeConns(test.args.action, test.args.fn)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
