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
package grpc

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
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
		UnimplementedValdServer   vald.UnimplementedValdServer
		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
		eg                        errgroup.Group
		gateway                   service.Gateway
		mirror                    service.Mirror
		vAddr                     string
		name                      string
		ip                        string
		streamConcurrency         int
	}
	type want struct {
		wantCe *payload.Object_Location
		err    error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		fields     fields
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
				eg:                        test.fields.eg,
				gateway:                   test.fields.gateway,
				mirror:                    test.fields.mirror,
				vAddr:                     test.fields.vAddr,
				streamConcurrency:         test.fields.streamConcurrency,
				name:                      test.fields.name,
				ip:                        test.fields.ip,
				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
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
		UnimplementedValdServer   vald.UnimplementedValdServer
		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
		eg                        errgroup.Group
		gateway                   service.Gateway
		mirror                    service.Mirror
		vAddr                     string
		name                      string
		ip                        string
		streamConcurrency         int
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		fields     fields
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02", "vald-03",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
							for _, tgt := range targets {
								f(ctx, tgt, cmap[tgt])
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
				eg:                        test.fields.eg,
				gateway:                   test.fields.gateway,
				mirror:                    test.fields.mirror,
				vAddr:                     test.fields.vAddr,
				streamConcurrency:         test.fields.streamConcurrency,
				name:                      test.fields.name,
				ip:                        test.fields.ip,
				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
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
		UnimplementedValdServer   vald.UnimplementedValdServer
		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
		eg                        errgroup.Group
		gateway                   service.Gateway
		mirror                    service.Mirror
		vAddr                     string
		name                      string
		ip                        string
		streamConcurrency         int
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		fields     fields
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ service.MirrorClient, _ ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ service.MirrorClient, _ ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
				eg:                        test.fields.eg,
				gateway:                   test.fields.gateway,
				mirror:                    test.fields.mirror,
				vAddr:                     test.fields.vAddr,
				streamConcurrency:         test.fields.streamConcurrency,
				name:                      test.fields.name,
				ip:                        test.fields.ip,
				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
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
		UnimplementedValdServer   vald.UnimplementedValdServer
		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
		eg                        errgroup.Group
		gateway                   service.Gateway
		mirror                    service.Mirror
		vAddr                     string
		name                      string
		ip                        string
		streamConcurrency         int
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		fields     fields
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid := "test"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
				eg:                        test.fields.eg,
				gateway:                   test.fields.gateway,
				mirror:                    test.fields.mirror,
				vAddr:                     test.fields.vAddr,
				streamConcurrency:         test.fields.streamConcurrency,
				name:                      test.fields.name,
				ip:                        test.fields.ip,
				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
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
		UnimplementedValdServer   vald.UnimplementedValdServer
		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
		eg                        errgroup.Group
		gateway                   service.Gateway
		mirror                    service.Mirror
		vAddr                     string
		name                      string
		ip                        string
		streamConcurrency         int
	}
	type want struct {
		wantLocs *payload.Object_Locations
		err      error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		fields     fields
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
			ctx, cancel := context.WithCancel(t.Context())
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
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
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
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, egctx := errgroup.New(ctx)

			uuid1 := "test01"
			uuid2 := "test02"
			targets := []string{
				"vald-01", "vald-02",
			}
			cmap := map[string]service.MirrorClient{
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
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error) error {
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
				eg:                        test.fields.eg,
				gateway:                   test.fields.gateway,
				mirror:                    test.fields.mirror,
				vAddr:                     test.fields.vAddr,
				streamConcurrency:         test.fields.streamConcurrency,
				name:                      test.fields.name,
				ip:                        test.fields.ip,
				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
			}

			gotLocs, err := s.RemoveByTimestamp(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLocs, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// Test_server_write_broadCastReachedNoTarget pins that a user write whose
// BroadCast returned nil without invoking the callback for any target (every
// target skipped because its pool had no usable conn or its circuit breaker was
// open) is reported as Unavailable instead of OK with an empty result.
func Test_server_write_broadCastReachedNoTarget(t *testing.T) {
	t.Parallel()
	vec := &payload.Object_Vector{Id: "uuid-1", Vector: []float32{1, 2, 3}}
	tests := []struct {
		call func(context.Context, *server) (any, error)
		name string
	}{
		{
			name: "Insert",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Insert(ctx, &payload.Insert_Request{Vector: vec})
			},
		},
		{
			name: "Update",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Update(ctx, &payload.Update_Request{Vector: vec})
			},
		},
		{
			name: "Upsert",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Upsert(ctx, &payload.Upsert_Request{Vector: vec})
			},
		},
		{
			name: "Remove",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.Remove(ctx, &payload.Remove_Request{Id: &payload.Object_ID{Id: vec.GetId()}})
			},
		},
		{
			name: "RemoveByTimestamp",
			call: func(ctx context.Context, s *server) (any, error) {
				return s.RemoveByTimestamp(ctx, &payload.Remove_TimestampRequest{})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var calls int
			s := &server{
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string {
						return "" // request from a user: broadcast to every target
					},
					BroadCastFunc: func(context.Context, func(context.Context, string, service.MirrorClient, ...grpc.CallOption) error) error {
						calls++
						return nil // every target skipped: callback never invoked
					},
				},
			}
			res, err := test.call(t.Context(), s)
			if calls != 1 {
				t.Fatalf("BroadCast calls = %d, want 1", calls)
			}
			if err == nil {
				t.Fatalf("returned OK (%v) for a write that reached no target, want Unavailable", res)
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != codes.Unavailable {
				t.Fatalf("status = %v (ok=%t), want code Unavailable: %v", st, ok, err)
			}
		})
	}
}

// Test_wrapTargetError pins the code and message wrapTargetError produces for
// each class of target error, which every proxy and do* handler relies on.
func Test_wrapTargetError(t *testing.T) {
	t.Parallel()
	const rpc = vald.InsertRPCName
	tests := []struct {
		err     error
		name    string
		wantMsg string
		want    codes.Code
	}{
		{name: "canceled", err: context.Canceled, want: codes.Canceled, wantMsg: rpc + " API canceld"},
		{name: "deadline exceeded", err: context.DeadlineExceeded, want: codes.DeadlineExceeded, wantMsg: rpc + " API deadline exceeded"},
		{name: "target not found", err: errors.ErrTargetNotFound, want: codes.Internal, wantMsg: rpc + " API target not found"},
		{name: "connection not found", err: errors.ErrGRPCClientConnNotFound("*"), want: codes.Internal, wantMsg: rpc + " API connection not found"},
		{name: "gRPC status passes through", err: status.Error(codes.NotFound, "nf"), want: codes.NotFound, wantMsg: "nf"},
		{name: "non-status error", err: errors.New("boom"), want: codes.Internal, wantMsg: "failed to parse " + rpc + " gRPC error response"},
	}
	s := &server{name: "mirror", ip: "192.0.2.1", vAddr: "lb:8081"}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := wrapTargetError(nil, rpc, test.err,
				&errdetails.RequestInfo{RequestId: "id"}, s.vAddrResourceInfo(rpc))
			st, ok := status.FromError(err)
			if !ok || st == nil {
				t.Fatalf("wrapTargetError() = %v, want a gRPC status error", err)
			}
			if st.Code() != test.want || st.Message() != test.wantMsg {
				t.Errorf("status = (%v, %q), want (%v, %q)", st.Code(), st.Message(), test.want, test.wantMsg)
			}
		})
	}
}

func Test_server_vAddrResourceInfo(t *testing.T) {
	t.Parallel()
	s := &server{name: "mirror", ip: "192.0.2.1", vAddr: "lb:8081"}
	got := s.vAddrResourceInfo(vald.SearchRPCName)
	wantType := errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName
	wantName := apiName + ": mirror(192.0.2.1) to lb:8081"
	if got.GetResourceType() != wantType || got.GetResourceName() != wantName {
		t.Errorf("vAddrResourceInfo() = (%q, %q), want (%q, %q)",
			got.GetResourceType(), got.GetResourceName(), wantType, wantName)
	}
}

// Test_statusOf pins that statusOf passes a usable status through and falls
// back to Unknown with the rpc-specific message only when there is none.
func Test_statusOf(t *testing.T) {
	t.Parallel()
	const rpc = vald.UpdateRPCName
	tests := []struct {
		err     error
		name    string
		wantMsg string
		want    codes.Code
	}{
		{name: "gRPC status", err: status.Error(codes.NotFound, "nf"), want: codes.NotFound, wantMsg: "nf"},
		{name: "plain error keeps its message", err: errors.New("boom"), want: codes.Unknown, wantMsg: "boom"},
		{name: "empty status message falls back", err: status.Error(codes.Internal, ""), want: codes.Unknown, wantMsg: "failed to parse " + rpc + " gRPC error response"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			st := statusOf(test.err, rpc)
			if st == nil {
				t.Fatal("statusOf() = nil")
			}
			if st.Code() != test.want || st.Message() != test.wantMsg {
				t.Errorf("statusOf() = (%v, %q), want (%v, %q)", st.Code(), st.Message(), test.want, test.wantMsg)
			}
		})
	}
}

// Test_server_write_concurrentBroadCast runs every BroadCast callback of a user
// write on its own goroutine, as RangeConcurrent does, so that -race catches an
// unguarded write to the request-local result map.
func Test_server_write_concurrentBroadCast(t *testing.T) {
	t.Parallel()
	const targets = 8
	vec := &payload.Object_Vector{Id: "id-1", Vector: []float32{1, 2, 3}}
	tests := []struct {
		call func(context.Context, *server) (*payload.Object_Location, error)
		name string
	}{
		{
			name: "Insert",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Insert(ctx, &payload.Insert_Request{Vector: vec})
			},
		},
		{
			name: "Update",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Update(ctx, &payload.Update_Request{Vector: vec})
			},
		},
		{
			name: "Upsert",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Upsert(ctx, &payload.Upsert_Request{Vector: vec})
			},
		},
		{
			name: "Remove",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Remove(ctx, &payload.Remove_Request{Id: &payload.Object_ID{Id: vec.GetId()}})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s := &server{
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string { return "" },
					BroadCastFunc: func(ctx context.Context, f func(context.Context, string, service.MirrorClient, ...grpc.CallOption) error) error {
						var wg sync.WaitGroup
						for i := range targets {
							wg.Go(func() {
								if err := f(ctx, "target-"+strconv.Itoa(i), probeMirrorClient{}); err != nil {
									t.Errorf("callback for target-%d = %v, want nil", i, err)
								}
							})
						}
						wg.Wait()
						return nil
					},
				},
			}
			type result struct {
				loc *payload.Object_Location
				err error
			}
			done := make(chan result, 1)
			go func() {
				loc, err := test.call(t.Context(), s)
				done <- result{loc: loc, err: err}
			}()
			select {
			case got := <-done:
				if got.err != nil {
					t.Fatalf("err = %v, want nil", got.err)
				}
				if got.loc.GetUuid() != vec.GetId() {
					t.Errorf("uuid = %q, want %q", got.loc.GetUuid(), vec.GetId())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timed out waiting for the handler")
			}
		})
	}
}

// proxiedFrom is what FromForwardedContext reports for a request another mirror
// gateway proxied, which makes the handler forward it to its own LB only.
const proxiedFrom = "mirror-gateway"

// Test_server_proxiedDoError verifies that a proxied Upsert/Remove surfaces the
// error returned by Gateway.Do even when Do fails before invoking the callback
// (e.g. ErrTargetNotFound), instead of returning (nil, nil).
func Test_server_proxiedDoError(t *testing.T) {
	t.Parallel()
	const uuid = "test"
	type test struct {
		call func(ctx context.Context, s *server) (*payload.Object_Location, error)
		name string
	}
	tests := []test{
		{
			name: "Upsert returns error when Do fails without calling f",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Upsert(ctx, &payload.Upsert_Request{
					Vector: &payload.Object_Vector{Id: uuid},
				})
			},
		},
		{
			name: "Remove returns error when Do fails without calling f",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Remove(ctx, &payload.Remove_Request{
					Id: &payload.Object_ID{Id: uuid},
				})
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			s := &server{
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string {
						return proxiedFrom // proxied request
					},
					DoFunc: func(context.Context, string, func(context.Context, string, service.MirrorClient, ...grpc.CallOption) (any, error)) (any, error) {
						return nil, errors.ErrTargetNotFound
					},
				},
			}
			loc, err := test.call(tt.Context(), s)
			if err == nil {
				tt.Fatalf("expected an error, got nil (loc: %v)", loc)
			}
			if loc != nil {
				tt.Errorf("expected nil location, got %v", loc)
			}
			if st, ok := status.FromError(err); !ok || st.Code() != codes.Internal {
				tt.Errorf("expected status code %v, got err: %v", codes.Internal, err)
			}
		})
	}
}

// Test_server_StreamListObject_noReplayOnGatewayFallback models gateway.Do
// falling back to another mirror after the own cluster's listing failed
// mid-stream: the fallback must not open a second listing on the same stream.
func Test_server_StreamListObject_noReplayOnGatewayFallback(t *testing.T) {
	t.Parallel()
	const n = 3
	// bounds the relay loop even if the stream never terminates.
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	failErr := status.New(codes.Unavailable, "agent stream broken").Err()
	own := &listAgentClient{n: n, failErr: failErr}
	peer := &listAgentClient{n: n, failErr: failErr}
	s := &server{
		streamConcurrency: 1,
		gateway: &gatewayMock{
			DoFunc: func(ctx context.Context, target string,
				f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error),
			) (any, error) {
				res, err := f(ctx, target, own)
				if err == nil {
					return res, nil
				}
				return f(ctx, "peer-mirror", peer)
			},
		},
	}
	ss := &listServerStream{ctx: ctx}

	err := s.StreamListObject(&payload.Object_List_Request{}, ss)
	if ctx.Err() != nil {
		t.Fatalf("StreamListObject did not finish before the test deadline: %v", ctx.Err())
	}
	if !status.Is(err, codes.Aborted) {
		t.Errorf("err = %v, want Aborted", err)
	}
	if got := peer.opens.Load(); got != 0 {
		t.Errorf("fallback mirror stream opened %d times, want 0", got)
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if len(ss.ids) != n {
		t.Errorf("forwarded %d messages %v, want %d (no duplicates)", len(ss.ids), ss.ids, n)
	}
}

// Test_server_StreamListObject_abortedOnNonAbortedFallbackError models the
// production RoundRobin fallback, which returns the error of the last conn it
// visited (ErrTargetNotFound for the own LB address) instead of the callback's
// Aborted: the caller must still get Aborted after a partial delivery.
func Test_server_StreamListObject_abortedOnNonAbortedFallbackError(t *testing.T) {
	t.Parallel()
	const n = 3
	// bounds the relay loop even if the stream never terminates.
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	failErr := status.New(codes.Unavailable, "agent stream broken").Err()
	own := &listAgentClient{n: n, failErr: failErr}
	peer := &listAgentClient{n: n, failErr: failErr}
	s := &server{
		streamConcurrency: 1,
		gateway: &gatewayMock{
			DoFunc: func(ctx context.Context, target string,
				f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error),
			) (any, error) {
				if res, err := f(ctx, target, own); err == nil {
					return res, nil
				}
				if res, err := f(ctx, "peer-mirror", peer); err == nil {
					return res, nil
				}
				return nil, errors.ErrTargetNotFound
			},
		},
	}
	ss := &listServerStream{ctx: ctx}

	err := s.StreamListObject(&payload.Object_List_Request{}, ss)
	if ctx.Err() != nil {
		t.Fatalf("StreamListObject did not finish before the test deadline: %v", ctx.Err())
	}
	if !status.Is(err, codes.Aborted) {
		t.Errorf("err = %v, want Aborted", err)
	}
	if got := peer.opens.Load(); got != 0 {
		t.Errorf("fallback mirror stream opened %d times, want 0", got)
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if len(ss.ids) != n {
		t.Errorf("forwarded %d messages %v, want %d (no duplicates)", len(ss.ids), ss.ids, n)
	}
}

// Test_server_Update_repairInsertSkipped verifies that Update does not report
// success when the repair Insert for a NotFound cluster was never attempted
// (DoMulti skipped the target, e.g. no healthy connection).
func Test_server_Update_repairInsertSkipped(t *testing.T) {
	t.Parallel()
	const (
		uuid    = "test"
		agentIP = "127.0.0.1"
		target1 = "vald-01"
		target2 = "vald-02"
	)
	loc := &payload.Object_Location{Uuid: uuid, Ips: []string{agentIP}}
	targets := []string{target1, target2}
	cmap := map[string]service.MirrorClient{
		targets[0]: &clientmock.MirrorClientMock{
			UpdateFunc: func(context.Context, *payload.Update_Request, ...grpc.CallOption) (*payload.Object_Location, error) {
				return loc, nil
			},
		},
		targets[1]: &clientmock.MirrorClientMock{
			UpdateFunc: func(context.Context, *payload.Update_Request, ...grpc.CallOption) (*payload.Object_Location, error) {
				return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
			},
		},
	}
	s := &server{
		gateway: &gatewayMock{
			FromForwardedContextFunc: func(context.Context) string { return "" },
			BroadCastFunc: func(ctx context.Context,
				f func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error,
			) error {
				for _, tgt := range targets {
					if err := f(ctx, tgt, cmap[tgt]); err != nil {
						t.Logf("BroadCast to %s: %v", tgt, err)
					}
				}
				return nil
			},
			// Like OrderedRangeConcurrent for a target without a healthy
			// connection: f is never called and no error is returned.
			DoMultiFunc: func(context.Context, []string,
				func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error,
			) error {
				return nil
			},
		},
	}

	gotLoc, err := s.Update(t.Context(), &payload.Update_Request{
		Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
		Config: &payload.Update_Config{SkipStrictExistCheck: true},
	})
	if !status.Is(err, codes.NotFound) {
		t.Errorf("err = %v, want NotFound", err)
	}
	if gotLoc != nil {
		t.Errorf("loc = %#v, want nil", gotLoc)
	}
}

// fanOutProbe records how many calls to hold overlap. A call stays in flight
// until more than limit calls overlap (then every call is released at once, so
// an unbounded fan-out is seen without waiting) or a short grace period ends
// (so a bounded fan-out still makes progress).
type fanOutProbe struct {
	over   chan struct{}
	mu     sync.Mutex
	limit  int
	cur    int
	peak   int
	calls  int
	closed bool
}

func newFanOutProbe(limit int) *fanOutProbe {
	return &fanOutProbe{over: make(chan struct{}), limit: limit}
}

func (p *fanOutProbe) hold() {
	p.mu.Lock()
	p.cur++
	p.calls++
	p.peak = max(p.peak, p.cur)
	if p.cur > p.limit && !p.closed {
		p.closed = true
		close(p.over)
	}
	p.mu.Unlock()
	timer := time.NewTimer(20 * time.Millisecond)
	select {
	case <-p.over:
	case <-timer.C:
	}
	timer.Stop()
	p.mu.Lock()
	p.cur--
	p.mu.Unlock()
}

func (p *fanOutProbe) stats() (peak, calls int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peak, p.calls
}

// probeMirrorClient answers every write with a location echoing the uuid; the
// uuid "id-0" panics instead when panics is set. Other methods panic through
// the nil embedded interface.
type probeMirrorClient struct {
	service.MirrorClient
	panics bool
}

func (c probeMirrorClient) location(uuid string) (*payload.Object_Location, error) {
	if c.panics && uuid == "id-0" {
		panic("probe panic")
	}
	return &payload.Object_Location{Uuid: uuid}, nil
}

func (c probeMirrorClient) Insert(
	_ context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetVector().GetId())
}

func (c probeMirrorClient) Update(
	_ context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetVector().GetId())
}

func (c probeMirrorClient) Upsert(
	_ context.Context, in *payload.Upsert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetVector().GetId())
}

func (c probeMirrorClient) Remove(
	_ context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.location(in.GetId().GetId())
}

// Test_server_Multi_boundedFanOut pins that the Multi* write handlers run at
// most streamConcurrency items at once, still answer every item in request
// order, and report a panicking item as an error instead of dropping it.
//
// It deliberately does not call t.Parallel: its subtests hold handler
// goroutines open, which the goleak.VerifyNone checks of the parallel
// handler_test.go cases would otherwise report as leaks.
func Test_server_Multi_boundedFanOut(t *testing.T) {
	const (
		n     = 6
		limit = 2
	)
	id := func(i int) string { return "id-" + strconv.Itoa(i) }
	vec := func(i int) *payload.Object_Vector {
		return &payload.Object_Vector{Id: id(i), Vector: []float32{1, 2, 3}}
	}
	tests := []struct {
		call func(context.Context, *server) (*payload.Object_Locations, error)
		name string
	}{
		{
			name: "MultiInsert",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Insert_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Insert_Request{Vector: vec(i)})
				}
				return s.MultiInsert(ctx, reqs)
			},
		},
		{
			name: "MultiUpdate",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Update_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: vec(i)})
				}
				return s.MultiUpdate(ctx, reqs)
			},
		},
		{
			name: "MultiUpsert",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Upsert_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: vec(i)})
				}
				return s.MultiUpsert(ctx, reqs)
			},
		},
		{
			name: "MultiRemove",
			call: func(ctx context.Context, s *server) (*payload.Object_Locations, error) {
				reqs := &payload.Remove_MultiRequest{}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{Id: id(i)}})
				}
				return s.MultiRemove(ctx, reqs)
			},
		},
	}
	for _, test := range tests {
		for _, panics := range []bool{false, true} {
			t.Run(test.name+"/panics="+strconv.FormatBool(panics), func(t *testing.T) {
				t.Parallel()
				p := newFanOutProbe(limit)
				// eg is the shared group the handlers used to fan out on.
				eg, _ := errgroup.New(t.Context())
				s := &server{
					eg: eg,
					gateway: &gatewayMock{
						FromForwardedContextFunc: func(context.Context) string {
							return proxiedFrom // proxied request: forward to the own LB only
						},
						DoFunc: func(ctx context.Context, target string, f func(context.Context, string, service.MirrorClient, ...grpc.CallOption) (any, error)) (any, error) {
							p.hold()
							return f(ctx, target, probeMirrorClient{panics: panics})
						},
					},
					streamConcurrency: limit,
				}
				type result struct {
					res *payload.Object_Locations
					err error
				}
				done := make(chan result, 1)
				go func() {
					res, err := test.call(t.Context(), s)
					done <- result{res: res, err: err}
				}()
				var got result
				select {
				case got = <-done:
				case <-time.After(10 * time.Second):
					t.Fatal("timed out waiting for the handler")
				}
				if peak, _ := p.stats(); peak > limit {
					t.Errorf("peak concurrency = %d, want <= %d", peak, limit)
				}
				if panics {
					// The handlers may re-wrap the recovered panic into a status,
					// so only its presence is pinned here.
					if got.err == nil {
						t.Error("err = nil, want the recovered panic")
					}
					// The panicking item 0 must not make the group skip the queued
					// items: every item ran and every other item has its location.
					if _, calls := p.stats(); calls != n {
						t.Errorf("forwarded %d items, want all %d", calls, n)
					}
					locs := got.res.GetLocations()
					if len(locs) != n {
						t.Fatalf("got %d locations, want %d", len(locs), n)
					}
					for i := 1; i < n; i++ {
						if want := id(i); locs[i].GetUuid() != want {
							t.Errorf("locations[%d].Uuid = %q, want %q", i, locs[i].GetUuid(), want)
						}
					}
					return
				}
				if got.err != nil {
					t.Fatalf("unexpected error: %v", got.err)
				}
				if len(got.res.GetLocations()) != n {
					t.Fatalf("got %d locations, want %d", len(got.res.GetLocations()), n)
				}
				for i, l := range got.res.GetLocations() {
					if want := id(i); l.GetUuid() != want {
						t.Errorf("locations[%d].Uuid = %q, want %q", i, l.GetUuid(), want)
					}
				}
			})
		}
	}
}

// failMirrorClient fails every write with a NotFound error naming the id;
// other methods panic through the nil embedded interface.
type failMirrorClient struct {
	service.MirrorClient
}

func failID(id string) error {
	return status.Error(codes.NotFound, "write failed for "+id)
}

func (failMirrorClient) Insert(
	_ context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetVector().GetId())
}

func (failMirrorClient) Update(
	_ context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetVector().GetId())
}

func (failMirrorClient) Upsert(
	_ context.Context, in *payload.Upsert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetVector().GetId())
}

func (failMirrorClient) Remove(
	_ context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, failID(in.GetId().GetId())
}

// Test_server_multi_joinedErrors pins the error the Multi* write handlers
// return now that per-item failures are collected in a slice and joined once
// after the workers finish: N distinct failures still surface as one status
// carrying the members' code whose message lists each failure exactly once.
func Test_server_multi_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 8
	ids := make([]string, 0, n)
	for i := range n {
		ids = append(ids, "id-"+strconv.Itoa(i))
	}
	vec := func(id string) *payload.Object_Vector {
		return &payload.Object_Vector{Id: id, Vector: []float32{1, 2, 3}}
	}
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "MultiInsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Insert_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Insert_Request{Vector: vec(id)})
				}
				_, err := s.MultiInsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpdate",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Update_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: vec(id)})
				}
				_, err := s.MultiUpdate(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Upsert_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: vec(id)})
				}
				_, err := s.MultiUpsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiRemove",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Remove_MultiRequest{}
				for _, id := range ids {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{Id: id}})
				}
				_, err := s.MultiRemove(ctx, reqs)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, ectx := errgroup.New(ctx)
			s := &server{
				eg: eg,
				gateway: &gatewayMock{
					FromForwardedContextFunc: func(context.Context) string {
						return proxiedFrom // proxied request: forward to the own LB only
					},
					DoFunc: func(ctx context.Context, target string, f func(context.Context, string, service.MirrorClient, ...grpc.CallOption) (any, error)) (any, error) {
						return f(ctx, target, failMirrorClient{})
					},
				},
			}
			done := make(chan error, 1)
			go func() { done <- test.call(ectx, s) }()
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				t.Fatal("timed out waiting for the handler")
			}
			if err == nil {
				t.Fatal("returned nil, want the per-item failures")
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != codes.NotFound {
				t.Fatalf("status = %v (ok=%t), want code NotFound: %v", st, ok, err)
			}
			// The joined members are rendered one per line; a flat join of N
			// distinct failures gives exactly N lines.
			if got := strings.Count(st.Message(), "\n") + 1; got != n {
				t.Errorf("message lines = %d, want %d: %q", got, n, st.Message())
			}
			for _, id := range ids {
				if !strings.Contains(st.Message(), "write failed for "+id) {
					t.Errorf("message %q does not mention %q", st.Message(), id)
				}
			}
		})
	}
}

// listAgentStream yields n vectors and then fails every further Recv with failErr.
type listAgentStream struct {
	vald.Object_StreamListObjectClient
	failErr error
	n       int
	recvd   int
}

func (s *listAgentStream) Recv() (*payload.Object_List_Response, error) {
	if s.recvd >= s.n {
		return nil, s.failErr
	}
	s.recvd++
	return &payload.Object_List_Response{
		Payload: &payload.Object_List_Response_Vector{
			Vector: &payload.Object_Vector{Id: fmt.Sprintf("uuid-%d", s.recvd)},
		},
	}, nil
}

// listAgentClient opens a fresh listAgentStream (restarting from the first
// vector, as a real vald gateway does) on every StreamListObject call.
type listAgentClient struct {
	service.MirrorClient
	failErr error
	n       int
	opens   atomic.Int64
}

func (c *listAgentClient) StreamListObject(
	context.Context, *payload.Object_List_Request, ...grpc.CallOption,
) (vald.Object_StreamListObjectClient, error) {
	c.opens.Add(1)
	return &listAgentStream{n: c.n, failErr: c.failErr}, nil
}

// listServerStream records the ids forwarded to the caller.
type listServerStream struct {
	vald.Object_StreamListObjectServer
	ctx context.Context //nolint:containedctx // a fake grpc.ServerStream must return its stream context
	ids []string
	mu  sync.Mutex
}

func (s *listServerStream) Context() context.Context { return s.ctx }

func (s *listServerStream) Send(res *payload.Object_List_Response) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids = append(s.ids, res.GetVector().GetId())
	return nil
}

// retryLikeExecuteRPC re-invokes f the way the backoff in internal/net/grpc
// gRPCClient.do does: executeRPC treats context errors as final and
// Internal/Unavailable/ResourceExhausted or non-status errors as retryable
// (for a healthy pool).
func retryLikeExecuteRPC(maxAttempts int, f func() error) (err error) {
	for range maxAttempts {
		err = f()
		if err == nil {
			return nil
		}
		if errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
			return err
		}
		st, ok := status.FromError(err)
		if !ok || st == nil {
			continue
		}
		if c := st.Code(); c == codes.Internal || c == codes.Unavailable || c == codes.ResourceExhausted {
			continue
		}
		return err
	}
	return err
}

func Test_server_doStreamListObject_noReplayAfterPartialDelivery(t *testing.T) {
	t.Parallel()
	const maxAttempts = 3
	tests := []struct {
		name      string
		n         int
		wantCode  codes.Code
		wantOpens int64
		wantSent  int
	}{
		{
			name:      "failure after partial delivery is aborted and not replayed",
			n:         3,
			wantCode:  codes.Aborted,
			wantOpens: 1,
			wantSent:  3,
		},
		{
			name:      "failure before the first message keeps its retryable code",
			n:         0,
			wantCode:  codes.Unavailable,
			wantOpens: maxAttempts,
			wantSent:  0,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			// bounds the relay loop even if the stream never terminates.
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			ac := &listAgentClient{
				n:       test.n,
				failErr: status.New(codes.Unavailable, "agent stream broken").Err(),
			}
			s := &server{streamConcurrency: 1}
			ss := &listServerStream{ctx: ctx}
			var forwarded atomic.Bool

			// doStreamListObject is the function gateway.Do re-invokes on retry.
			err := retryLikeExecuteRPC(maxAttempts, func() error {
				return s.doStreamListObject(ctx, ac, &payload.Object_List_Request{}, ss, &forwarded)
			})
			if ctx.Err() != nil {
				tt.Fatalf("doStreamListObject did not finish before the test deadline: %v", ctx.Err())
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != test.wantCode {
				tt.Errorf("error code = %v (err: %v), want %v", st.Code(), err, test.wantCode)
			}
			if got := ac.opens.Load(); got != test.wantOpens {
				tt.Errorf("agent stream opened %d times, want %d", got, test.wantOpens)
			}
			ss.mu.Lock()
			defer ss.mu.Unlock()
			if len(ss.ids) != test.wantSent {
				tt.Errorf("forwarded %d messages %v, want %d (no duplicates)", len(ss.ids), ss.ids, test.wantSent)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Server
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Server, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Server, err error) error {
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Mirror_Targets
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx  context.Context
// 		meta *payload.Object_ID
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 	defaultCheckFunc := func(w want, gotId *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotId, w.wantId) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotId, w.wantId)
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotId, err := s.Exists(test.args.ctx, test.args.meta)
// 			if err := checkFunc(test.want, gotId, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Search(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		stream vald.Search_StreamSearchServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		stream vald.Search_StreamSearchByIDServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchByIDServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// func Test_server_handleInsert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Insert_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.handleInsert(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_handleInsertResult(t *testing.T) {
// 	type args struct {
// 		ctx               context.Context
// 		alreadyExistsTgts []string
// 		req               *payload.Update_Request
// 		result            *sync.Map[string, *errorState]
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           alreadyExistsTgts:nil,
// 		           req:nil,
// 		           result:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           alreadyExistsTgts:nil,
// 		           req:nil,
// 		           result:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.handleInsertResult(test.args.ctx, test.args.alreadyExistsTgts, test.args.req, test.args.result)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_doInsert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Insert_Request
// 		f   func(ctx context.Context) (*payload.Object_Location, error)
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.doInsert(test.args.ctx, test.args.req, test.args.f)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamInsert(t *testing.T) {
// 	type args struct {
// 		stream vald.Insert_StreamInsertServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Insert_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// func Test_server_handleUpdate(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.handleUpdate(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_handleUpdateResult(t *testing.T) {
// 	type args struct {
// 		ctx          context.Context
// 		notFoundTgts []string
// 		req          *payload.Insert_Request
// 		result       *sync.Map[string, *errorState]
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           notFoundTgts:nil,
// 		           req:nil,
// 		           result:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           notFoundTgts:nil,
// 		           req:nil,
// 		           result:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.handleUpdateResult(test.args.ctx, test.args.notFoundTgts, test.args.req, test.args.result)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_doUpdate(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_Request
// 		f   func(ctx context.Context) (*payload.Object_Location, error)
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.doUpdate(test.args.ctx, test.args.req, test.args.f)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpdate(t *testing.T) {
// 	type args struct {
// 		stream vald.Update_StreamUpdateServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// func Test_server_handleUpsert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Upsert_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.handleUpsert(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_doUpsert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Upsert_Request
// 		f   func(ctx context.Context) (*payload.Object_Location, error)
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.doUpsert(test.args.ctx, test.args.req, test.args.f)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpsert(t *testing.T) {
// 	type args struct {
// 		stream vald.Upsert_StreamUpsertServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Upsert_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// func Test_server_handleRemove(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.handleRemove(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_doRemove(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_Request
// 		f   func(ctx context.Context) (*payload.Object_Location, error)
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLoc, err := s.doRemove(test.args.ctx, test.args.req, test.args.f)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamRemove(t *testing.T) {
// 	type args struct {
// 		stream vald.Remove_StreamRemoveServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Remove_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// func Test_server_handleRemoveByTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_TimestampRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
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
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.handleRemoveByTimestamp(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_doRemoveByTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_TimestampRequest
// 		f   func(ctx context.Context) (*payload.Object_Locations, error)
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
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
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			gotLocs, err := s.doRemoveByTimestamp(test.args.ctx, test.args.req, test.args.f)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_GetObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Object_VectorRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
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
// 	type args struct {
// 		stream vald.Object_StreamGetObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamGetObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamListObject(t *testing.T) {
// 	type args struct {
// 		req    *payload.Object_List_Request
// 		stream vald.Object_StreamListObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           req:nil,
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           req:nil,
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamListObject(test.args.req, test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_doStreamListObject(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		client vald.Object_StreamListObjectClient
// 		server vald.Object_StreamListObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           client:nil,
// 		           server:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           client:nil,
// 		           server:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			err := s.doStreamListObject(test.args.ctx, test.args.client, test.args.server)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Flush(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Flush_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
// 	}
// 	type want struct {
// 		want *payload.Info_Index_Count
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Count, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Info_Index_Count, err error) error {
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			got, err := s.Flush(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_isProxied(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		UnimplementedValdServer   vald.UnimplementedValdServer
// 		UnimplementedMirrorServer mirror.UnimplementedMirrorServer
// 		eg                        errgroup.Group
// 		gateway                   service.Gateway
// 		mirror                    service.Mirror
// 		vAddr                     string
// 		name                      string
// 		ip                        string
// 		streamConcurrency         int
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           UnimplementedValdServer:nil,
// 		           UnimplementedMirrorServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           mirror:nil,
// 		           vAddr:"",
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 				UnimplementedValdServer:   test.fields.UnimplementedValdServer,
// 				UnimplementedMirrorServer: test.fields.UnimplementedMirrorServer,
// 				eg:                        test.fields.eg,
// 				gateway:                   test.fields.gateway,
// 				mirror:                    test.fields.mirror,
// 				vAddr:                     test.fields.vAddr,
// 				name:                      test.fields.name,
// 				ip:                        test.fields.ip,
// 				streamConcurrency:         test.fields.streamConcurrency,
// 			}
//
// 			got := s.isProxied(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
