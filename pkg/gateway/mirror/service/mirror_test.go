// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package service

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/test/goleak"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

func Test_mirr_Connect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		targets []*payload.Mirror_Target
	}
	type fields struct {
		gatewayAddr  string
		selfMirrAddr string
		gateway      Gateway
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
		func() test {
			gatewayAddr := "192.168.1.2:8081"
			selfMirrorAddr := "192.168.1.3:8081"
			return test{
				name: "Succeeded to connect to other mirror gateways",
				args: args{
					ctx: context.Background(),
					targets: []*payload.Mirror_Target{
						{
							Host: "192.168.2.2",
							Port: 8081,
						},
						{
							Host: "192.168.3.2",
							Port: 8081,
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return false
								},
								ConnectFunc: func(_ context.Context, _ string, _ ...grpc.DialOption) (conn pool.Conn, err error) {
									return conn, err
								},
							}
						},
					},
				},
			}
		}(),
		func() test {
			gatewayAddr := "192.168.1.2:8081"
			selfMirrorAddr := "192.168.1.3:8081"
			return test{
				name: "Failed to connect to other mirror gateways due to an invalid address",
				args: args{
					ctx: context.Background(),
					targets: []*payload.Mirror_Target{
						{
							Host: "192.168.2.2",
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return false
								},
								ConnectFunc: func(_ context.Context, _ string, _ ...grpc.DialOption) (pool.Conn, error) {
									return nil, errors.New("missing port in address")
								},
							}
						},
					},
				},
				want: want{
					err: errors.New("missing port in address"),
				},
			}
		}(),
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
			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			err = m.Connect(test.args.ctx, test.args.targets...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mirr_Disconnect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		targets []*payload.Mirror_Target
	}
	type fields struct {
		gatewayAddr  string
		selfMirrAddr string
		gateway      Gateway
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
		func() test {
			gatewayAddr := "192.168.1.2:8081"
			selfMirrorAddr := "192.168.1.3:8081"
			return test{
				name: "Succeeded to disconnect to other mirror gateways",
				args: args{
					ctx: context.Background(),
					targets: []*payload.Mirror_Target{
						{
							Host: "192.168.2.2",
							Port: 8081,
						},
						{
							Host: "192.168.3.2",
							Port: 8081,
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return true
								},
								DisconnectFunc: func(_ context.Context, _ string) error {
									return nil
								},
							}
						},
					},
				},
			}
		}(),
		func() test {
			gatewayAddr := "192.168.1.2:8081"
			selfMirrorAddr := "192.168.1.3:8081"
			return test{
				name: "Failed to connect to other mirror gateways due to an invalid address",
				args: args{
					ctx: context.Background(),
					targets: []*payload.Mirror_Target{
						{
							Host: "192.168.2.2",
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return true
								},
								DisconnectFunc: func(_ context.Context, _ string) error {
									return errors.New("missing port in address")
								},
							}
						},
					},
				},
				want: want{
					err: errors.New("missing port in address"),
				},
			}
		}(),
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
			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			err = m.Disconnect(test.args.ctx, test.args.targets...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mirr_MirrorTargets(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		gatewayAddr  string
		selfMirrAddr string
		gateway      Gateway
	}
	type want struct {
		want []*payload.Mirror_Target
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []*payload.Mirror_Target, error) error
		beforeFunc func(*testing.T, Mirror)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got []*payload.Mirror_Target, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if len(got) != len(w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		sort.Slice(got, func(i, j int) bool {
			return got[i].Host > got[j].Host
		})
		sort.Slice(w.want, func(i, j int) bool {
			return w.want[i].Host > w.want[j].Host
		})
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			gatewayAddr := "192.168.1.2:8081"
			selfMirrorAddr := "192.168.1.3:8081"
			connectTargets := []*payload.Mirror_Target{
				{
					Host: "192.168.1.2", // gateway addresses
					Port: 8081,
				},
				{
					Host: "192.168.2.2", // other mirror address
					Port: 8081,
				},
				{
					Host: "192.168.3.2", // other mirror address
					Port: 8081,
				},
			}
			connected := make(map[string]bool)
			return test{
				name: "returns only the addresses of the mirror gateways",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					gatewayAddr:  gatewayAddr,
					selfMirrAddr: selfMirrorAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (conn pool.Conn, err error) {
									connected[addr] = true
									return conn, err
								},
								IsConnectedFunc: func(_ context.Context, addr string) bool {
									return connected[addr]
								},
							}
						},
					},
				},
				beforeFunc: func(t *testing.T, m Mirror) {
					t.Helper()
					if err := m.Connect(context.Background(), connectTargets...); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: []*payload.Mirror_Target{
						{
							Host: "192.168.2.2", // other mirror address
							Port: 8081,
						},
						{
							Host: "192.168.3.2", // other mirror address
							Port: 8081,
						},
						{
							Host: "192.168.1.3", // self mirror address
							Port: 8081,
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, m)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := m.MirrorTargets(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mirr_connectedOtherMirrorAddrs(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		gatewayAddr  string
		selfMirrAddr string
		gateway      Gateway
	}
	type want struct {
		want []string
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func(*testing.T, Mirror)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got []string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			gatewayAddr := "192.168.1.2:8081"
			selfMirrorAddr := "192.168.1.3:8081"
			connectTargets := []*payload.Mirror_Target{
				{
					Host: "192.168.1.2", // gateway addresses
					Port: 8081,
				},
				{
					Host: "192.168.2.2", // other mirror address
					Port: 8081,
				},
			}
			connected := make(map[string]bool)
			return test{
				name: "returns only the address of the other mirror gateway",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (conn pool.Conn, err error) {
									connected[addr] = true
									return conn, err
								},
								IsConnectedFunc: func(_ context.Context, addr string) bool {
									return connected[addr]
								},
							}
						},
					},
				},
				beforeFunc: func(t *testing.T, m Mirror) {
					t.Helper()
					if err := m.Connect(context.Background(), connectTargets...); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: []string{
						"192.168.2.2:8081",
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, m)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			if mirr, ok := m.(*mirr); ok {
				got := mirr.connectedOtherMirrorAddrs(test.args.ctx)
				if err := checkFunc(test.want, got); err != nil {
					tt.Errorf("error = %v", err)
				}
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewMirror(t *testing.T) {
// 	type args struct {
// 		opts []MirrorOption
// 	}
// 	type want struct {
// 		want Mirror
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Mirror, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Mirror, err error) error {
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
// 			got, err := NewMirror(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_mirr_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		addrl         sync.Map[string, any]
// 		selfMirrTgts  []*payload.Mirror_Target
// 		selfMirrAddrl sync.Map[string, any]
// 		gwAddrl       sync.Map[string, any]
// 		eg            errgroup.Group
// 		registerDur   time.Duration
// 		gateway       Gateway
// 	}
// 	type want struct {
// 		want <-chan error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error) error {
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
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 			m := &mirr{
// 				addrl:         test.fields.addrl,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				selfMirrAddrl: test.fields.selfMirrAddrl,
// 				gwAddrl:       test.fields.gwAddrl,
// 				eg:            test.fields.eg,
// 				registerDur:   test.fields.registerDur,
// 				gateway:       test.fields.gateway,
// 			}
//
// 			got := m.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_mirr_Disconnect(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		targets []*payload.Mirror_Target
// 	}
// 	type fields struct {
// 		addrl         sync.Map[string, any]
// 		selfMirrTgts  []*payload.Mirror_Target
// 		selfMirrAddrl sync.Map[string, any]
// 		gwAddrl       sync.Map[string, any]
// 		eg            errgroup.Group
// 		registerDur   time.Duration
// 		gateway       Gateway
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
// 		           targets:nil,
// 		       },
// 		       fields: fields {
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 		           targets:nil,
// 		           },
// 		           fields: fields {
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 			m := &mirr{
// 				addrl:         test.fields.addrl,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				selfMirrAddrl: test.fields.selfMirrAddrl,
// 				gwAddrl:       test.fields.gwAddrl,
// 				eg:            test.fields.eg,
// 				registerDur:   test.fields.registerDur,
// 				gateway:       test.fields.gateway,
// 			}
//
// 			err := m.Disconnect(test.args.ctx, test.args.targets...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_mirr_IsConnected(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		addrl         sync.Map[string, any]
// 		selfMirrTgts  []*payload.Mirror_Target
// 		selfMirrAddrl sync.Map[string, any]
// 		gwAddrl       sync.Map[string, any]
// 		eg            errgroup.Group
// 		registerDur   time.Duration
// 		gateway       Gateway
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
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 			m := &mirr{
// 				addrl:         test.fields.addrl,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				selfMirrAddrl: test.fields.selfMirrAddrl,
// 				gwAddrl:       test.fields.gwAddrl,
// 				eg:            test.fields.eg,
// 				registerDur:   test.fields.registerDur,
// 				gateway:       test.fields.gateway,
// 			}
//
// 			got := m.IsConnected(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_mirr_RangeAllMirrorAddr(t *testing.T) {
// 	type args struct {
// 		f func(addr string, _ any) bool
// 	}
// 	type fields struct {
// 		addrl         sync.Map[string, any]
// 		selfMirrTgts  []*payload.Mirror_Target
// 		selfMirrAddrl sync.Map[string, any]
// 		gwAddrl       sync.Map[string, any]
// 		eg            errgroup.Group
// 		registerDur   time.Duration
// 		gateway       Gateway
// 	}
// 	type want struct {
// 	}
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           addrl:nil,
// 		           selfMirrTgts:nil,
// 		           selfMirrAddrl:nil,
// 		           gwAddrl:nil,
// 		           eg:nil,
// 		           registerDur:nil,
// 		           gateway:nil,
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
// 			m := &mirr{
// 				addrl:         test.fields.addrl,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				selfMirrAddrl: test.fields.selfMirrAddrl,
// 				gwAddrl:       test.fields.gwAddrl,
// 				eg:            test.fields.eg,
// 				registerDur:   test.fields.registerDur,
// 				gateway:       test.fields.gateway,
// 			}
//
// 			m.RangeAllMirrorAddr(test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
