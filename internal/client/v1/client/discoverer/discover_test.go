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

// Package discoverer
package discoverer

import (
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/mock"
)

func Test_client_GetReadClient(t *testing.T) {
	type fields struct {
		client              grpc.Client
		readClient          grpc.Client
		readReplicaReplicas uint64
		roundRobin          atomic.Uint64
	}
	type test struct {
		want   grpc.Client
		name   string
		fields fields
	}

	mockClient := mock.ClientInternal{}
	mockClient.On("GetAddrs").Return([]string{"read write client"})
	mockReadClient := mock.ClientInternal{}
	mockReadClient.On("GetAddrs").Return([]string{"read replica client"})

	tests := []test{
		{
			name: "returns primary client when there is no read replica",
			fields: fields{
				client:              &mockClient,
				readClient:          nil,
				readReplicaReplicas: 1,
			},
			want: &mockClient,
		},
		func() test {
			var counter atomic.Uint64
			counter.Store(0)
			return test{
				name: "returns read client when there is read replica and the counter increments to anything other than 0",
				fields: fields{
					client:              &mockClient,
					readClient:          &mockReadClient,
					readReplicaReplicas: 1,
					//nolint:govet,copylocks
					//skipcq: VET-V0008
					roundRobin: counter,
				},
				want: &mockReadClient,
			}
		}(),
		func() test {
			var counter atomic.Uint64
			counter.Store(1)
			return test{
				name: "returns primary client when there is read replica and the counter increments to 0",
				fields: fields{
					client:              &mockClient,
					readClient:          &mockReadClient,
					readReplicaReplicas: 1,
					//nolint:govet,copylocks
					//skipcq: VET-V0008
					roundRobin: counter,
				},
				want: &mockClient,
			}
		}(),
		func() test {
			var counter atomic.Uint64
			counter.Store(3)
			return test{
				name: "returns primary client when there is read replica and the counter increments to 0(replicas: 3)",
				fields: fields{
					client:              &mockClient,
					readClient:          &mockReadClient,
					readReplicaReplicas: 3,
					//nolint:govet,copylocks
					//skipcq: VET-V0008
					roundRobin: counter,
				},
				want: &mockClient,
			}
		}(),
	}
	//nolint:govet,copylocks
	//skipcq: VET-V0008
	for _, tc := range tests {
		//nolint:govet,copylocks
		//skipcq: VET-V0008
		test := tc
		t.Run(test.name, func(t *testing.T) {
			c := &client{
				client:              test.fields.client,
				readClient:          test.fields.readClient,
				readReplicaReplicas: test.fields.readReplicaReplicas,
				//nolint:govet,copylocks
				//skipcq: VET-V0008
				roundRobin: test.fields.roundRobin,
			}
			got := c.GetReadClient()
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetReadClient() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_client_GetReadClient_concurrent(t *testing.T) {
	mockClient := mock.ClientInternal{}
	mockClient.On("GetAddrs").Return([]string{"read write client"})
	mockReadClient := mock.ClientInternal{}
	mockReadClient.On("GetAddrs").Return([]string{"read replica client"})

	c := &client{
		client:              &mockClient,
		readClient:          &mockReadClient,
		readReplicaReplicas: 100,
		roundRobin:          atomic.Uint64{},
	}

	eg, _ := errgroup.New(t.Context())
	for range 150 {
		eg.Go(func() error {
			c.GetReadClient()
			return nil
		})
	}

	err := eg.Wait()
	require.NoError(t, err)

	require.EqualValues(t, uint64(49), c.roundRobin.Load(), "atomic operation did not happen in the concurrent calls")
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		wantD Client
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Client, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotD Client, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotD, w.wantD) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotD, w.wantD)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
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
// 			gotD, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, gotD, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			got, err := c.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_GetAddrs(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			gotAddrs := c.GetAddrs(test.args.ctx)
// 			if err := checkFunc(test.want, gotAddrs); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_GetClient(t *testing.T) {
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
// 	}
// 	type want struct {
// 		want grpc.Client
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, grpc.Client) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got grpc.Client) error {
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			got := c.GetClient()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_connect(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			err := c.connect(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_disconnect(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			err := c.disconnect(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_dnsDiscovery(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
// 	}
// 	type want struct {
// 		wantAddrs []string
// 		err       error
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
// 	defaultCheckFunc := func(w want, gotAddrs []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			gotAddrs, err := c.dnsDiscovery(test.args.ctx)
// 			if err := checkFunc(test.want, gotAddrs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_discover(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			err := c.discover(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_updateDiscoveryInfo(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
// 	}
// 	type want struct {
// 		wantConnected []string
// 		err           error
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
// 	defaultCheckFunc := func(w want, gotConnected []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotConnected, w.wantConnected) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConnected, w.wantConnected)
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			gotConnected, err := c.updateDiscoveryInfo(test.args.ctx)
// 			if err := checkFunc(test.want, gotConnected, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_discoverNodes(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
// 	}
// 	type want struct {
// 		wantNodes *payload.Info_Nodes
// 		err       error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Nodes, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotNodes *payload.Info_Nodes, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotNodes, w.wantNodes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNodes, w.wantNodes)
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			gotNodes, err := c.discoverNodes(test.args.ctx)
// 			if err := checkFunc(test.want, gotNodes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_discoverAddrs(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		nodes *payload.Info_Nodes
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
// 	}
// 	type want struct {
// 		wantAddrs []string
// 		err       error
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
// 	defaultCheckFunc := func(w want, gotAddrs []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           ctx:nil,
// 		           nodes:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           nodes:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			gotAddrs, err := c.discoverAddrs(test.args.ctx, test.args.nodes)
// 			if err := checkFunc(test.want, gotAddrs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_disconnectOldAddrs(t *testing.T) {
// 	type args struct {
// 		ctx            context.Context
// 		oldAddrs       []string
// 		connectedAddrs []string
// 	}
// 	type fields struct {
// 		client              grpc.Client
// 		readClient          grpc.Client
// 		eg                  errgroup.Group
// 		dscClient           grpc.Client
// 		addrs               atomic.Pointer[[]string]
// 		onDisconnect        func(ctx context.Context, c Client, addr string) error
// 		onDiscover          func(ctx context.Context, c Client, addrs []string) error
// 		onConnect           func(ctx context.Context, c Client, addr string) error
// 		nodeName            string
// 		name                string
// 		namespace           string
// 		dns                 string
// 		opts                []grpc.Option
// 		dscDur              time.Duration
// 		port                int
// 		readReplicaReplicas uint64
// 		roundRobin          atomic.Uint64
// 		autoconn            bool
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
// 		           oldAddrs:nil,
// 		           connectedAddrs:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 		           oldAddrs:nil,
// 		           connectedAddrs:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           readClient:nil,
// 		           eg:nil,
// 		           dscClient:nil,
// 		           addrs:nil,
// 		           onDisconnect:nil,
// 		           onDiscover:nil,
// 		           onConnect:nil,
// 		           nodeName:"",
// 		           name:"",
// 		           namespace:"",
// 		           dns:"",
// 		           opts:nil,
// 		           dscDur:nil,
// 		           port:0,
// 		           readReplicaReplicas:0,
// 		           roundRobin:nil,
// 		           autoconn:false,
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
// 			c := &client{
// 				client:              test.fields.client,
// 				readClient:          test.fields.readClient,
// 				eg:                  test.fields.eg,
// 				dscClient:           test.fields.dscClient,
// 				addrs:               test.fields.addrs,
// 				onDisconnect:        test.fields.onDisconnect,
// 				onDiscover:          test.fields.onDiscover,
// 				onConnect:           test.fields.onConnect,
// 				nodeName:            test.fields.nodeName,
// 				name:                test.fields.name,
// 				namespace:           test.fields.namespace,
// 				dns:                 test.fields.dns,
// 				opts:                test.fields.opts,
// 				dscDur:              test.fields.dscDur,
// 				port:                test.fields.port,
// 				readReplicaReplicas: test.fields.readReplicaReplicas,
// 				roundRobin:          test.fields.roundRobin,
// 				autoconn:            test.fields.autoconn,
// 			}
//
// 			err := c.disconnectOldAddrs(test.args.ctx, test.args.oldAddrs, test.args.connectedAddrs)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
