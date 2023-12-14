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
package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	clientmock "github.com/vdaas/vald/internal/test/mock/client"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

func Test_server_Insert(t *testing.T) { // skipcq: GO-R1005
	t.Parallel()
	const dimension = 128
	defaultInsertConfig := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}
	type args struct {
		ctx context.Context
		req *payload.Insert_Request
	}
	type fields struct {
		eg                                errgroup.Group
		gateway                           service.Gateway
		mirror                            service.Mirror
		vAddr                             string
		streamConcurrency                 int
		name                              string
		ip                                string
		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
	}
	type want struct {
		wantCe *payload.Object_Location
		err    error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotCe *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			gotSt, gotOk := status.FromError(err)
			wantSt, wantOk := status.FromError(w.err)
			if gotOk != wantOk || gotSt.Code() != wantSt.Code() {
				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
			}
		}
		if !reflect.DeepEqual(gotCe, w.wantCe) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotCe, w.wantCe)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			return test{
				name: "Success: insert with new ID",
				args: args{
					ctx: egctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultInsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantCe: &payload.Object_Location{
						Uuid: uuid,
						Ips:  []string{"127.0.0.1", "127.0.0.1"},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return &payload.Object_Location{
							Uuid: uuid,
							Ips:  []string{"127.0.0.1"},
						}, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			return test{
				name: "Success: when the last status codes are (OK, OK) after updating the target that returned AlreadyExists",
				args: args{
					ctx: egctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultInsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, target := range targets {
								if c, ok := cmap[target]; !ok {
									return errors.ErrTargetNotFound
								} else {
									f(ctx, target, c)
								}
							}
							return nil
						},
					},
				},
				want: want{
					wantCe: &payload.Object_Location{
						Uuid: uuid,
						Ips:  []string{"127.0.0.1", "127.0.0.1"},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return &payload.Object_Location{
							Uuid: uuid,
							Ips:  []string{"127.0.0.1"},
						}, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Success: when the last status codes are (OK, AlreadyExists) after updating the target that returned AlreadyExists",
				args: args{
					ctx: egctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultInsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, target := range targets {
								if c, ok := cmap[target]; !ok {
									return errors.New("target not found")
								} else {
									f(ctx, target, c)
								}
							}
							return nil
						},
					},
				},
				want: want{
					wantCe: &payload.Object_Location{
						Uuid: uuid,
						Ips:  []string{"127.0.0.1"},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (AlreadyExists, AlreadyExists)",
				args: args{
					ctx: egctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultInsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.AlreadyExists, vald.InsertRPCName+" API target same vector already exists"),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (OK, Internal)",
				args: args{
					ctx: egctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultInsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (Internal, Internal)",
				args: args{
					ctx: egctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultInsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.Join(
						status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error()),
						status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
					).Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return &payload.Object_Location{
							Uuid: uuid,
							Ips:  []string{"127.0.0.1"},
						}, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "Fail: when the last status codes are (OK, Internal) after updating the target that returned AlreadyExists",
				args: args{
					ctx: egctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultInsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, target := range targets {
								if c, ok := cmap[target]; !ok {
									return errors.New("target not found")
								} else {
									f(ctx, target, c)
								}
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
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
			s := &server{
				eg:                                test.fields.eg,
				gateway:                           test.fields.gateway,
				mirror:                            test.fields.mirror,
				vAddr:                             test.fields.vAddr,
				streamConcurrency:                 test.fields.streamConcurrency,
				name:                              test.fields.name,
				ip:                                test.fields.ip,
				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
			}

			gotCe, err := s.Insert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotCe, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Update(t *testing.T) { // skipcq: GO-R1005
	t.Parallel()
	const dimension = 128
	defaultUpdateConfig := &payload.Update_Config{
		SkipStrictExistCheck: true,
	}
	type args struct {
		ctx context.Context
		req *payload.Update_Request
	}
	type fields struct {
		eg                                errgroup.Group
		gateway                           service.Gateway
		mirror                            service.Mirror
		vAddr                             string
		streamConcurrency                 int
		name                              string
		ip                                string
		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			gotSt, gotOk := status.FromError(err)
			wantSt, wantOk := status.FromError(w.err)
			if gotOk != wantOk || gotSt.Code() != wantSt.Code() {
				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
			}
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			return test{
				name: "Success: update with new ID",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: &payload.Object_Location{
						Uuid: uuid,
						Ips:  []string{"127.0.0.1", "127.0.0.1"},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Success: when the status codes are (AlreadyExists, OK)",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: &payload.Object_Location{
						Uuid: uuid,
						Ips:  []string{"127.0.0.1"},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[2]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			return test{
				name: "Success: when the last status codes are (OK, OK, OK) after inserting the target that returned NotFound",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, target := range targets {
								if c, ok := cmap[target]; !ok {
									return errors.ErrTargetNotFound
								} else {
									f(ctx, target, c)
								}
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: &payload.Object_Location{
						Uuid: uuid,
						Ips: []string{
							"127.0.0.1", "127.0.0.1", "127.0.0.1",
						},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[2]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Success: when the last status codes are (OK, OK, AlreadyExists) after inserting the target that returned NotFound",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, target := range targets {
								if c, ok := cmap[target]; !ok {
									return errors.ErrTargetNotFound
								} else {
									f(ctx, target, c)
								}
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: &payload.Object_Location{
						Uuid: uuid,
						Ips: []string{
							"127.0.0.1", "127.0.0.1",
						},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (NotFound, NotFound)",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.NotFound, vald.UpdateRPCName+" API id "+uuid+" not found"),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (Internal, OK)",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (Internal, AlreadyExists)",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.Join(
						status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
						status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error()),
					).Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
				targets[2]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Fail: when the last status codes are (AlreadyExists, AlreadyExists, AlreadyExists) after inserting the target that returned NotFound",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, target := range targets {
								if c, ok := cmap[target]; !ok {
									return errors.ErrTargetNotFound
								} else {
									f(ctx, target, c)
								}
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.AlreadyExists, vald.InsertRPCName+" for "+vald.UpdateRPCName+" API target same vector already exists"),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				targets[2]: &clientmock.MirrorClientMock{
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			return test{
				name: "Fail: when the last status codes are (OK, OK, Internal) after inserting the target that returned NotFound",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, target := range targets {
								if c, ok := cmap[target]; !ok {
									return errors.New("target not found")
								} else {
									f(ctx, target, c)
								}
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
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
			s := &server{
				eg:                                test.fields.eg,
				gateway:                           test.fields.gateway,
				mirror:                            test.fields.mirror,
				vAddr:                             test.fields.vAddr,
				streamConcurrency:                 test.fields.streamConcurrency,
				name:                              test.fields.name,
				ip:                                test.fields.ip,
				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
			}

			gotLoc, err := s.Update(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Upsert(t *testing.T) {
	t.Parallel()
	const dimension = 128
	defaultUpsertConfig := &payload.Upsert_Config{
		SkipStrictExistCheck: true,
	}
	type args struct {
		ctx context.Context
		req *payload.Upsert_Request
	}
	type fields struct {
		eg                                errgroup.Group
		gateway                           service.Gateway
		mirror                            service.Mirror
		vAddr                             string
		streamConcurrency                 int
		name                              string
		ip                                string
		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			gotSt, gotOk := status.FromError(err)
			wantSt, wantOk := status.FromError(w.err)
			if gotOk != wantOk || gotSt.Code() != wantSt.Code() {
				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
			}
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			return test{
				name: "Success: upsert with new ID",
				args: args{
					ctx: egctx,
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: &payload.Object_Location{
						Uuid: uuid,
						Ips:  []string{"127.0.0.1", "127.0.0.1"},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpsertFunc: func(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Success: when the status codes are (AlreadyExists, OK)",
				args: args{
					ctx: egctx,
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: loc,
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpsertFunc: func(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.AlreadyExists, errors.ErrMetaDataAlreadyExists(uuid).Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (AlreadyExists, AlreadyExists)",
				args: args{
					ctx: egctx,
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.AlreadyExists, vald.UpsertRPCName+" API target same vector already exists"),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (Internal, OK)",
				args: args{
					ctx: egctx,
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "Fail: upsert when the status codes are (Internal, Internal)",
				args: args{
					ctx: egctx,
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     uuid,
							Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
						},
						Config: defaultUpsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.Join(
						status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
						status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error()),
					).Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
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
			s := &server{
				eg:                                test.fields.eg,
				gateway:                           test.fields.gateway,
				mirror:                            test.fields.mirror,
				vAddr:                             test.fields.vAddr,
				streamConcurrency:                 test.fields.streamConcurrency,
				name:                              test.fields.name,
				ip:                                test.fields.ip,
				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
			}

			gotLoc, err := s.Upsert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Remove(t *testing.T) {
	t.Parallel()
	const dimension = 128
	defaultRemoveConfig := &payload.Remove_Config{
		SkipStrictExistCheck: true,
	}
	type args struct {
		ctx context.Context
		req *payload.Remove_Request
	}
	type fields struct {
		eg                                errgroup.Group
		gateway                           service.Gateway
		mirror                            service.Mirror
		vAddr                             string
		streamConcurrency                 int
		name                              string
		ip                                string
		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			gotSt, gotOk := status.FromError(err)
			wantSt, wantOk := status.FromError(w.err)
			if gotOk != wantOk || gotSt.Code() != wantSt.Code() {
				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
			}
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			return test{
				name: "Success: remove with existing ID",
				args: args{
					ctx: egctx,
					req: &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: uuid,
						},
						Config: defaultRemoveConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: &payload.Object_Location{
						Uuid: uuid,
						Ips:  []string{"127.0.0.1", "127.0.0.1"},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
				},
			}
			return test{
				name: "Success: when the status codes are (NotFound, OK)",
				args: args{
					ctx: egctx,
					req: &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: uuid,
						},
						Config: defaultRemoveConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: loc,
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (Internal, OK)",
				args: args{
					ctx: egctx,
					req: &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: uuid,
						},
						Config: defaultRemoveConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (Internal, Internal)",
				args: args{
					ctx: egctx,
					req: &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: uuid,
						},
						Config: defaultRemoveConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.Join(
						status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
						status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error()),
					).Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrIndexNotFound.Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrIndexNotFound.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (NotFound, NotFound)",
				args: args{
					ctx: egctx,
					req: &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: uuid,
						},
						Config: defaultRemoveConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.NotFound, vald.RemoveRPCName+" API id "+uuid+" not found"),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
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
			s := &server{
				eg:                                test.fields.eg,
				gateway:                           test.fields.gateway,
				mirror:                            test.fields.mirror,
				vAddr:                             test.fields.vAddr,
				streamConcurrency:                 test.fields.streamConcurrency,
				name:                              test.fields.name,
				ip:                                test.fields.ip,
				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
			}

			gotLoc, err := s.Remove(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_RemoveByTimestamp(t *testing.T) {
	t.Parallel()
	defaultRemoveByTimestampReq := &payload.Remove_TimestampRequest{
		Timestamps: []*payload.Remove_Timestamp{},
	}
	type args struct {
		ctx context.Context
		req *payload.Remove_TimestampRequest
	}
	type fields struct {
		eg                                errgroup.Group
		gateway                           service.Gateway
		mirror                            service.Mirror
		vAddr                             string
		streamConcurrency                 int
		name                              string
		ip                                string
		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
	}
	type want struct {
		wantLocs *payload.Object_Locations
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
		if !errors.Is(err, w.err) {
			gotSt, gotOk := status.FromError(err)
			wantSt, wantOk := status.FromError(w.err)
			if gotOk != wantOk || gotSt.Code() != wantSt.Code() {
				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
			}
		}
		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			loc := &payload.Object_Location{
				Uuid: "test",
				Ips: []string{
					"127.0.0.1",
				},
			}
			loc2 := &payload.Object_Location{
				Uuid: "test02",
				Ips: []string{
					"127.0.0.1",
				},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return &payload.Object_Locations{
							Locations: []*payload.Object_Location{
								loc,
							},
						}, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return &payload.Object_Locations{
							Locations: []*payload.Object_Location{
								loc2,
							},
						}, nil
					},
				},
			}
			return test{
				name: "Success: removeByTimestamp",
				args: args{
					ctx: egctx,
					req: defaultRemoveByTimestampReq,
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLocs: &payload.Object_Locations{
						Locations: []*payload.Object_Location{
							loc, loc2,
						},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			loc := &payload.Object_Location{
				Uuid: "test",
				Ips: []string{
					"127.0.0.1",
				},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return &payload.Object_Locations{
							Locations: []*payload.Object_Location{
								loc,
							},
						}, nil
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound("test02").Error())
					},
				},
			}
			return test{
				name: "Success: when the status codes are (NotFound, OK)",
				args: args{
					ctx: egctx,
					req: defaultRemoveByTimestampReq,
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					wantLocs: &payload.Object_Locations{
						Locations: []*payload.Object_Location{
							loc,
						},
					},
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (Internal, Internal)",
				args: args{
					ctx: egctx,
					req: defaultRemoveByTimestampReq,
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal, errors.Join(
						status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error()),
						status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error()),
					).Error()),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, egctx := errgroup.New(ctx)

			uuid1 := "test01"
			uuid2 := "test02"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]vald.ClientWithMirror{
				targets[0]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid1).Error())
					},
				},
				targets[1]: &clientmock.MirrorClientMock{
					RemoveByTimestampFunc: func(_ context.Context, _ *payload.Remove_TimestampRequest, _ ...grpc.CallOption) (*payload.Object_Locations, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid2).Error())
					},
				},
			}
			return test{
				name: "Fail: when the status codes are (NotFound, NotFound)",
				args: args{
					ctx: egctx,
					req: defaultRemoveByTimestampReq,
				},
				fields: fields{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
					},
				},
				want: want{
					err: status.Error(codes.NotFound, vald.RemoveByTimestampRPCName+" API target not found"),
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
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
			s := &server{
				eg:                                test.fields.eg,
				gateway:                           test.fields.gateway,
				mirror:                            test.fields.mirror,
				vAddr:                             test.fields.vAddr,
				streamConcurrency:                 test.fields.streamConcurrency,
				name:                              test.fields.name,
				ip:                                test.fields.ip,
				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
			}

			gotLocs, err := s.RemoveByTimestamp(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLocs, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW

// func TestNew(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want vald.ServerWithMirror
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, vald.Server, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got vald.Server, err error) error {
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Register(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Mirror_Targets
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		want *payload.Mirror_Targets
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Mirror_Targets, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Mirror_Targets, err error) error {
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			got, err := s.Register(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Exists(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx  context.Context
// 		meta *payload.Object_ID
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantId *payload.Object_ID
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_ID, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotID *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotID, w.wantId) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotID, w.wantId)
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
// 		           meta:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           meta:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotID, err := s.Exists(test.args.ctx, test.args.meta)
// 			if err := checkFunc(test.want, gotID, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Search(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.Search(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_SearchByID(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.SearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearch(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Search_StreamSearchServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearchByID(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Search_StreamSearchByIDServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearch(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiSearch(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearchByID(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_LinearSearch(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.LinearSearch(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_LinearSearchByID(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.LinearSearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearch(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamLinearSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearchByID(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchByIDServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamLinearSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearch(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiLinearSearch(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearchByID(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiLinearSearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamInsert(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Insert_StreamInsertServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamInsert(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiInsert(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Insert_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Locations
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiInsert(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpdate(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Update_StreamUpdateServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamUpdate(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpdate(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Locations
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpsert(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Upsert_StreamUpsertServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamUpsert(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpsert(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Upsert_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Locations
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiUpsert(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamRemove(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Remove_StreamRemoveServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamRemove(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiRemove(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Remove_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Locations
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotRes, err := s.MultiRemove(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_GetObject(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Object_VectorRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
// 	}
// 	type want struct {
// 		wantVec *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			gotVec, err := s.GetObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamGetObject(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		stream vald.Object_StreamGetObjectServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		gateway                           service.Gateway
// 		mirror                            service.Mirror
// 		vAddr                             string
// 		streamConcurrency                 int
// 		name                              string
// 		ip                                string
// 		UnimplementedValdServerWithMirror vald.UnimplementedValdServerWithMirror
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           streamConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServerWithMirror:nil,
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
// 			s := &server{
// 				eg:                                test.fields.eg,
// 				gateway:                           test.fields.gateway,
// 				mirror:                            test.fields.mirror,
// 				vAddr:                             test.fields.vAddr,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				UnimplementedValdServerWithMirror: test.fields.UnimplementedValdServerWithMirror,
// 			}
//
// 			err := s.StreamGetObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
