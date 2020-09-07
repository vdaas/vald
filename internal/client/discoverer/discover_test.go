//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package discoverer
package discoverer

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"

	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantD Client
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Client, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotD Client, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotD, w.wantD) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotD, w.wantD)
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotD, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, gotD, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		autoconn     bool
		onDiscover   func(ctx context.Context, c Client, addrs []string) error
		onConnect    func(ctx context.Context, c Client, addr string) error
		onDisconnect func(ctx context.Context, c Client, addr string) error
		client       grpc.Client
		dns          string
		opts         []grpc.Option
		port         int
		addrs        atomic.Value
		dscAddr      string
		dscClient    grpc.Client
		dscDur       time.Duration
		eg           errgroup.Group
		name         string
		namespace    string
		nodeName     string
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				autoconn:     test.fields.autoconn,
				onDiscover:   test.fields.onDiscover,
				onConnect:    test.fields.onConnect,
				onDisconnect: test.fields.onDisconnect,
				client:       test.fields.client,
				dns:          test.fields.dns,
				opts:         test.fields.opts,
				port:         test.fields.port,
				addrs:        test.fields.addrs,
				dscAddr:      test.fields.dscAddr,
				dscClient:    test.fields.dscClient,
				dscDur:       test.fields.dscDur,
				eg:           test.fields.eg,
				name:         test.fields.name,
				namespace:    test.fields.namespace,
				nodeName:     test.fields.nodeName,
			}

			got, err := c.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_GetAddrs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		autoconn     bool
		onDiscover   func(ctx context.Context, c Client, addrs []string) error
		onConnect    func(ctx context.Context, c Client, addr string) error
		onDisconnect func(ctx context.Context, c Client, addr string) error
		client       grpc.Client
		dns          string
		opts         []grpc.Option
		port         int
		addrs        atomic.Value
		dscAddr      string
		dscClient    grpc.Client
		dscDur       time.Duration
		eg           errgroup.Group
		name         string
		namespace    string
		nodeName     string
	}
	type want struct {
		wantAddrs []string
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotAddrs []string) error {
		if !reflect.DeepEqual(gotAddrs, w.wantAddrs) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotAddrs, w.wantAddrs)
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				autoconn:     test.fields.autoconn,
				onDiscover:   test.fields.onDiscover,
				onConnect:    test.fields.onConnect,
				onDisconnect: test.fields.onDisconnect,
				client:       test.fields.client,
				dns:          test.fields.dns,
				opts:         test.fields.opts,
				port:         test.fields.port,
				addrs:        test.fields.addrs,
				dscAddr:      test.fields.dscAddr,
				dscClient:    test.fields.dscClient,
				dscDur:       test.fields.dscDur,
				eg:           test.fields.eg,
				name:         test.fields.name,
				namespace:    test.fields.namespace,
				nodeName:     test.fields.nodeName,
			}

			gotAddrs := c.GetAddrs(test.args.ctx)
			if err := test.checkFunc(test.want, gotAddrs); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_GetClient(t *testing.T) {
	type fields struct {
		autoconn     bool
		onDiscover   func(ctx context.Context, c Client, addrs []string) error
		onConnect    func(ctx context.Context, c Client, addr string) error
		onDisconnect func(ctx context.Context, c Client, addr string) error
		client       grpc.Client
		dns          string
		opts         []grpc.Option
		port         int
		addrs        atomic.Value
		dscAddr      string
		dscClient    grpc.Client
		dscDur       time.Duration
		eg           errgroup.Group
		name         string
		namespace    string
		nodeName     string
	}
	type want struct {
		want grpc.Client
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, grpc.Client) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got grpc.Client) error {
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				autoconn:     test.fields.autoconn,
				onDiscover:   test.fields.onDiscover,
				onConnect:    test.fields.onConnect,
				onDisconnect: test.fields.onDisconnect,
				client:       test.fields.client,
				dns:          test.fields.dns,
				opts:         test.fields.opts,
				port:         test.fields.port,
				addrs:        test.fields.addrs,
				dscAddr:      test.fields.dscAddr,
				dscClient:    test.fields.dscClient,
				dscDur:       test.fields.dscDur,
				eg:           test.fields.eg,
				name:         test.fields.name,
				namespace:    test.fields.namespace,
				nodeName:     test.fields.nodeName,
			}

			got := c.GetClient()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_connect(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}
	type fields struct {
		autoconn     bool
		onDiscover   func(ctx context.Context, c Client, addrs []string) error
		onConnect    func(ctx context.Context, c Client, addr string) error
		onDisconnect func(ctx context.Context, c Client, addr string) error
		client       grpc.Client
		dns          string
		opts         []grpc.Option
		port         int
		addrs        atomic.Value
		dscAddr      string
		dscClient    grpc.Client
		dscDur       time.Duration
		eg           errgroup.Group
		name         string
		namespace    string
		nodeName     string
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				autoconn:     test.fields.autoconn,
				onDiscover:   test.fields.onDiscover,
				onConnect:    test.fields.onConnect,
				onDisconnect: test.fields.onDisconnect,
				client:       test.fields.client,
				dns:          test.fields.dns,
				opts:         test.fields.opts,
				port:         test.fields.port,
				addrs:        test.fields.addrs,
				dscAddr:      test.fields.dscAddr,
				dscClient:    test.fields.dscClient,
				dscDur:       test.fields.dscDur,
				eg:           test.fields.eg,
				name:         test.fields.name,
				namespace:    test.fields.namespace,
				nodeName:     test.fields.nodeName,
			}

			err := c.connect(test.args.ctx, test.args.addr)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_disconnect(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}
	type fields struct {
		autoconn     bool
		onDiscover   func(ctx context.Context, c Client, addrs []string) error
		onConnect    func(ctx context.Context, c Client, addr string) error
		onDisconnect func(ctx context.Context, c Client, addr string) error
		client       grpc.Client
		dns          string
		opts         []grpc.Option
		port         int
		addrs        atomic.Value
		dscAddr      string
		dscClient    grpc.Client
		dscDur       time.Duration
		eg           errgroup.Group
		name         string
		namespace    string
		nodeName     string
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
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
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				autoconn:     test.fields.autoconn,
				onDiscover:   test.fields.onDiscover,
				onConnect:    test.fields.onConnect,
				onDisconnect: test.fields.onDisconnect,
				client:       test.fields.client,
				dns:          test.fields.dns,
				opts:         test.fields.opts,
				port:         test.fields.port,
				addrs:        test.fields.addrs,
				dscAddr:      test.fields.dscAddr,
				dscClient:    test.fields.dscClient,
				dscDur:       test.fields.dscDur,
				eg:           test.fields.eg,
				name:         test.fields.name,
				namespace:    test.fields.namespace,
				nodeName:     test.fields.nodeName,
			}

			err := c.disconnect(test.args.ctx, test.args.addr)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_dnsDiscovery(t *testing.T) {
	type args struct {
		ctx context.Context
		ech chan<- error
	}
	type fields struct {
		autoconn     bool
		onDiscover   func(ctx context.Context, c Client, addrs []string) error
		onConnect    func(ctx context.Context, c Client, addr string) error
		onDisconnect func(ctx context.Context, c Client, addr string) error
		client       grpc.Client
		dns          string
		opts         []grpc.Option
		port         int
		addrs        atomic.Value
		dscAddr      string
		dscClient    grpc.Client
		dscDur       time.Duration
		eg           errgroup.Group
		name         string
		namespace    string
		nodeName     string
	}
	type want struct {
		wantAddrs []string
		err       error
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
	defaultCheckFunc := func(w want, gotAddrs []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotAddrs, w.wantAddrs) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotAddrs, w.wantAddrs)
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
		           ech: nil,
		       },
		       fields: fields {
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
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
		           ech: nil,
		           },
		           fields: fields {
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				autoconn:     test.fields.autoconn,
				onDiscover:   test.fields.onDiscover,
				onConnect:    test.fields.onConnect,
				onDisconnect: test.fields.onDisconnect,
				client:       test.fields.client,
				dns:          test.fields.dns,
				opts:         test.fields.opts,
				port:         test.fields.port,
				addrs:        test.fields.addrs,
				dscAddr:      test.fields.dscAddr,
				dscClient:    test.fields.dscClient,
				dscDur:       test.fields.dscDur,
				eg:           test.fields.eg,
				name:         test.fields.name,
				namespace:    test.fields.namespace,
				nodeName:     test.fields.nodeName,
			}

			gotAddrs, err := c.dnsDiscovery(test.args.ctx, test.args.ech)
			if err := test.checkFunc(test.want, gotAddrs, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_discover(t *testing.T) {
	type args struct {
		ctx context.Context
		ech chan<- error
	}
	type fields struct {
		autoconn     bool
		onDiscover   func(ctx context.Context, c Client, addrs []string) error
		onConnect    func(ctx context.Context, c Client, addr string) error
		onDisconnect func(ctx context.Context, c Client, addr string) error
		client       grpc.Client
		dns          string
		opts         []grpc.Option
		port         int
		addrs        atomic.Value
		dscAddr      string
		dscClient    grpc.Client
		dscDur       time.Duration
		eg           errgroup.Group
		name         string
		namespace    string
		nodeName     string
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
		           ech: nil,
		       },
		       fields: fields {
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
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
		           ech: nil,
		           },
		           fields: fields {
		           autoconn: false,
		           onDiscover: nil,
		           onConnect: nil,
		           onDisconnect: nil,
		           client: nil,
		           dns: "",
		           opts: nil,
		           port: 0,
		           addrs: nil,
		           dscAddr: "",
		           dscClient: nil,
		           dscDur: nil,
		           eg: nil,
		           name: "",
		           namespace: "",
		           nodeName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				autoconn:     test.fields.autoconn,
				onDiscover:   test.fields.onDiscover,
				onConnect:    test.fields.onConnect,
				onDisconnect: test.fields.onDisconnect,
				client:       test.fields.client,
				dns:          test.fields.dns,
				opts:         test.fields.opts,
				port:         test.fields.port,
				addrs:        test.fields.addrs,
				dscAddr:      test.fields.dscAddr,
				dscClient:    test.fields.dscClient,
				dscDur:       test.fields.dscDur,
				eg:           test.fields.eg,
				name:         test.fields.name,
				namespace:    test.fields.namespace,
				nodeName:     test.fields.nodeName,
			}

			err := c.discover(test.args.ctx, test.args.ech)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
