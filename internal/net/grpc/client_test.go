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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/singleflight"
	"go.uber.org/goleak"
	gbackoff "google.golang.org/grpc/backoff"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		wantC Client
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Client) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotC Client) error {
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotC := New(test.args.opts...)
			if err := test.checkFunc(test.want, gotC); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_StartConnectionMonitor(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			got, err := g.StartConnectionMonitor(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_Range(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		f   func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
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
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			err := g.Range(test.args.ctx, test.args.f)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_RangeConcurrent(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx         context.Context
		concurrency int
		f           func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
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
		           concurrency: 0,
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           concurrency: 0,
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			err := g.RangeConcurrent(test.args.ctx, test.args.concurrency, test.args.f)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_OrderedRange(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		orders []string
		f      func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
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
		           orders: nil,
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           orders: nil,
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			err := g.OrderedRange(test.args.ctx, test.args.orders, test.args.f)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_OrderedRangeConcurrent(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx         context.Context
		orders      []string
		concurrency int
		f           func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
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
		           orders: nil,
		           concurrency: 0,
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           orders: nil,
		           concurrency: 0,
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			err := g.OrderedRangeConcurrent(test.args.ctx, test.args.orders, test.args.concurrency, test.args.f)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_RoundRobin(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		f   func(ctx context.Context, conn *ClientConn, copts ...CallOption) (interface{}, error)
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		wantData interface{}
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotData interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotData, w.wantData) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			gotData, err := g.RoundRobin(test.args.ctx, test.args.f)
			if err := test.checkFunc(test.want, gotData, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_Do(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		addr string
		f    func(ctx context.Context, conn *ClientConn, copts ...CallOption) (interface{}, error)
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		wantData interface{}
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotData interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotData, w.wantData) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			gotData, err := g.Do(test.args.ctx, test.args.addr, test.args.f)
			if err := test.checkFunc(test.want, gotData, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_do(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx           context.Context
		p             pool.Conn
		addr          string
		enableBackoff bool
		f             func(ctx context.Context, conn *ClientConn, copts ...CallOption) (interface{}, error)
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		wantData interface{}
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotData interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotData, w.wantData) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
		           p: nil,
		           addr: "",
		           enableBackoff: false,
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           p: nil,
		           addr: "",
		           enableBackoff: false,
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			gotData, err := g.do(test.args.ctx, test.args.p, test.args.addr, test.args.enableBackoff, test.args.f)
			if err := test.checkFunc(test.want, gotData, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_GetDialOption(t *testing.T) {
	t.Parallel()
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		want []DialOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []DialOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []DialOption) error {
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			got := g.GetDialOption()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_GetCallOption(t *testing.T) {
	t.Parallel()
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		want []CallOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []CallOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []CallOption) error {
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			got := g.GetCallOption()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_Connect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		addr  string
		dopts []DialOption
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		wantConn pool.Conn
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, pool.Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotConn pool.Conn, err error) error {
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
		           dopts: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           dopts: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			gotConn, err := g.Connect(test.args.ctx, test.args.addr, test.args.dopts...)
			if err := test.checkFunc(test.want, gotConn, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_Disconnect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		addr string
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
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
		           addr: "",
		       },
		       fields: fields {
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			err := g.Disconnect(test.args.ctx, test.args.addr)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_ConnectedAddrs(t *testing.T) {
	t.Parallel()
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
	}
	type want struct {
		want []string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []string) error {
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			got := g.ConnectedAddrs()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gRPCClient_Close(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		addrs               map[string]struct{}
		atomicAddrs         AtomicAddrs
		poolSize            uint64
		clientCount         uint64
		conns               grpcConns
		hcDur               time.Duration
		prDur               time.Duration
		enablePoolRebalance bool
		resolveDNS          bool
		dopts               []DialOption
		copts               []CallOption
		roccd               string
		eg                  errgroup.Group
		bo                  backoff.Backoff
		gbo                 gbackoff.Config
		mcd                 time.Duration
		group               singleflight.Group
		crl                 sync.Map
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
		           addrs: nil,
		           atomicAddrs: nil,
		           poolSize: 0,
		           clientCount: 0,
		           conns: grpcConns{},
		           hcDur: nil,
		           prDur: nil,
		           enablePoolRebalance: false,
		           resolveDNS: false,
		           dopts: nil,
		           copts: nil,
		           roccd: "",
		           eg: nil,
		           bo: nil,
		           gbo: nil,
		           mcd: nil,
		           group: nil,
		           crl: sync.Map{},
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gRPCClient{
				addrs:               test.fields.addrs,
				atomicAddrs:         test.fields.atomicAddrs,
				poolSize:            test.fields.poolSize,
				clientCount:         test.fields.clientCount,
				conns:               test.fields.conns,
				hcDur:               test.fields.hcDur,
				prDur:               test.fields.prDur,
				enablePoolRebalance: test.fields.enablePoolRebalance,
				resolveDNS:          test.fields.resolveDNS,
				dopts:               test.fields.dopts,
				copts:               test.fields.copts,
				roccd:               test.fields.roccd,
				eg:                  test.fields.eg,
				bo:                  test.fields.bo,
				gbo:                 test.fields.gbo,
				mcd:                 test.fields.mcd,
				group:               test.fields.group,
				crl:                 test.fields.crl,
			}

			err := g.Close(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
