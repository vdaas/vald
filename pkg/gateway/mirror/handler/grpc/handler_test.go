package grpc

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

func Test_server_Insert(t *testing.T) {
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			wantLoc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1", "127.0.0.1"},
			}
			return test{
				name: "success insert with new ID",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
					},
				},
				want: want{
					wantCe: wantLoc,
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
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "fail insert with new ID but remove rollback success",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				"vald-lb-gateway-01": &mockClient{
					InsertFunc: func(_ context.Context, _ *payload.Insert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "fail insert with new ID and fail remove rollback",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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

func Test_server_Update(t *testing.T) {
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			wantLoc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1", "127.0.0.1"},
			}
			return test{
				name: "success update with new ID",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: wantLoc,
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
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "fail update with new ID but remove rollback success",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			ovec := &payload.Object_Vector{
				Id:     uuid,
				Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return ovec, nil
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "fail update with new ID but update rollback success",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: ovec,
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "fail update with new ID and fail remove rollback",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			ovec := &payload.Object_Vector{
				Id:     uuid,
				Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
			}
			var cnt uint32
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return ovec, nil
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						if atomic.AddUint32(&cnt, 1) == 1 {
							return loc, nil
						}
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "fail update with new ID and fail update rollback",
				args: args{
					ctx: egctx,
					req: &payload.Update_Request{
						Vector: ovec,
						Config: defaultUpdateConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
			}
			wantLoc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1", "127.0.0.1"},
			}
			return test{
				name: "success upsert with new ID",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: wantLoc,
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
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "fail upsert with new ID but remove rollback success",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			ovec := &payload.Object_Vector{
				Id:     uuid,
				Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return ovec, nil
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "fail upsert with new ID but update rollback success",
				args: args{
					ctx: egctx,
					req: &payload.Upsert_Request{
						Vector: ovec,
						Config: defaultUpsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "fail upsert with new ID and fail remove rollback",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			ovec := &payload.Object_Vector{
				Id:     uuid,
				Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return ovec, nil
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					UpdateFunc: func(_ context.Context, _ *payload.Update_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "fail upsert with new ID and fail update rollback",
				args: args{
					ctx: egctx,
					req: &payload.Upsert_Request{
						Vector: ovec,
						Config: defaultUpsertConfig,
					},
				},
				fields: fields{
					eg: eg,
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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

			gotLoc, err := s.Upsert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Remove(t *testing.T) {
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
			ovec := &payload.Object_Vector{
				Id:     uuid,
				Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return ovec, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
				},
			}
			wantLoc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			return test{
				name: "success remove with existing ID",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
					},
				},
				want: want{
					wantLoc: wantLoc,
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
			ovec := &payload.Object_Vector{
				Id:     uuid,
				Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return ovec, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
			}
			return test{
				name: "fail remove with existing ID but upsert rollback success",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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
			loc := &payload.Object_Location{
				Uuid: uuid,
				Ips:  []string{"127.0.0.1"},
			}
			ovec := &payload.Object_Vector{
				Id:     uuid,
				Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
			}
			cmap := map[string]vald.ClientWithMirror{
				"vald-mirror-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return ovec, nil
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, nil
					},
					UpsertFunc: func(_ context.Context, _ *payload.Upsert_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return nil, status.Error(codes.Internal, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
					},
				},
				"vald-lb-gateway-01": &mockClient{
					GetObjectFunc: func(_ context.Context, _ *payload.Object_VectorRequest, _ ...grpc.CallOption) (*payload.Object_Vector, error) {
						return nil, status.Error(codes.NotFound, errors.ErrObjectIDNotFound(uuid).Error())
					},
					RemoveFunc: func(_ context.Context, _ *payload.Remove_Request, _ ...grpc.CallOption) (*payload.Object_Location, error) {
						return loc, status.Error(codes.Internal, errors.ErrCircuitBreakerOpenState.Error())
					},
				},
			}
			return test{
				name: "fail remove with existing ID and fail upsert rollback",
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
					gateway: &mockGateway{
						FromForwardedContextFunc: func(_ context.Context) string {
							return ""
						},
						BroadCastFunc: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error {
							for tgt, c := range cmap {
								f(ctx, tgt, c)
							}
							return nil
						},
						DoMultiFunc: func(ctx context.Context, targets []string, f func(ctx context.Context, target string, vc vald.ClientWithMirror, _ ...grpc.CallOption) error) error {
							if len(targets) != 1 {
								return errors.New("invalid target")
							}
							if c, ok := cmap[targets[0]]; ok {
								f(ctx, targets[0], c)
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

			gotLoc, err := s.Remove(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want vald.ServerWithMirror
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, vald.Server, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got vald.Server, err error) error {
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Mirror_Targets
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
		want *payload.Mirror_Targets
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Mirror_Targets, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *payload.Mirror_Targets, err error) error {
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
		           ctx:nil,
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			got, err := s.Register(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Advertise(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Mirror_Targets
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
		wantRes *payload.Mirror_Targets
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Mirror_Targets, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Mirror_Targets, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.Advertise(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Exists(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta *payload.Object_ID
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
		wantId *payload.Object_ID
		err    error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_ID, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotId *payload.Object_ID, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotId, w.wantId) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotId, w.wantId)
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
		           meta:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           meta:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotId, err := s.Exists(test.args.ctx, test.args.meta)
			if err := checkFunc(test.want, gotId, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Search(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_Request
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
		wantRes *payload.Search_Response
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.Search(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_SearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_IDRequest
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
		wantRes *payload.Search_Response
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.SearchByID(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamSearch(t *testing.T) {
	type args struct {
		stream vald.Search_StreamSearchServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamSearch(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamSearchByID(t *testing.T) {
	type args struct {
		stream vald.Search_StreamSearchByIDServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamSearchByID(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiSearch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_MultiRequest
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
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiSearch(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiSearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_MultiIDRequest
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
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_LinearSearch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_Request
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
		wantRes *payload.Search_Response
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.LinearSearch(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_LinearSearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_IDRequest
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
		wantRes *payload.Search_Response
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.LinearSearchByID(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamLinearSearch(t *testing.T) {
	type args struct {
		stream vald.Search_StreamLinearSearchServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamLinearSearch(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamLinearSearchByID(t *testing.T) {
	type args struct {
		stream vald.Search_StreamLinearSearchByIDServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamLinearSearchByID(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiLinearSearch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_MultiRequest
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
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiLinearSearch(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiLinearSearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_MultiIDRequest
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
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiLinearSearchByID(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_insert(t *testing.T) {
	type args struct {
		ctx    context.Context
		client vald.InsertClient
		req    *payload.Insert_Request
		opts   []grpc.CallOption
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
		           client:nil,
		           req:nil,
		           opts:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           client:nil,
		           req:nil,
		           opts:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotLoc, err := s.insert(test.args.ctx, test.args.client, test.args.req, test.args.opts...)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamInsert(t *testing.T) {
	type args struct {
		stream vald.Insert_StreamInsertServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamInsert(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiInsert(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Insert_MultiRequest
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
		wantRes *payload.Object_Locations
		err     error
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
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           reqs:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           reqs:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiInsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_update(t *testing.T) {
	type args struct {
		ctx    context.Context
		client vald.UpdateClient
		req    *payload.Update_Request
		opts   []grpc.CallOption
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
		           client:nil,
		           req:nil,
		           opts:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           client:nil,
		           req:nil,
		           opts:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotLoc, err := s.update(test.args.ctx, test.args.client, test.args.req, test.args.opts...)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamUpdate(t *testing.T) {
	type args struct {
		stream vald.Update_StreamUpdateServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamUpdate(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiUpdate(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Update_MultiRequest
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
		wantRes *payload.Object_Locations
		err     error
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
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           reqs:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           reqs:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_upsert(t *testing.T) {
	type args struct {
		ctx    context.Context
		client vald.UpsertClient
		req    *payload.Upsert_Request
		opts   []grpc.CallOption
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
		           client:nil,
		           req:nil,
		           opts:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           client:nil,
		           req:nil,
		           opts:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotLoc, err := s.upsert(test.args.ctx, test.args.client, test.args.req, test.args.opts...)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamUpsert(t *testing.T) {
	type args struct {
		stream vald.Upsert_StreamUpsertServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamUpsert(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiUpsert(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Upsert_MultiRequest
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
		wantRes *payload.Object_Locations
		err     error
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
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           reqs:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           reqs:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiUpsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_remove(t *testing.T) {
	type args struct {
		ctx    context.Context
		client vald.RemoveClient
		req    *payload.Remove_Request
		opts   []grpc.CallOption
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
		want *payload.Object_Location
		err  error
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
	defaultCheckFunc := func(w want, got *payload.Object_Location, err error) error {
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
		           ctx:nil,
		           client:nil,
		           req:nil,
		           opts:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           client:nil,
		           req:nil,
		           opts:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			got, err := s.remove(test.args.ctx, test.args.client, test.args.req, test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamRemove(t *testing.T) {
	type args struct {
		stream vald.Remove_StreamRemoveServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamRemove(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiRemove(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Remove_MultiRequest
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
		wantRes *payload.Object_Locations
		err     error
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
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
		           reqs:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           reqs:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotRes, err := s.MultiRemove(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_GetObject(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Object_VectorRequest
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
		wantVec *payload.Object_Vector
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Vector, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotVec *payload.Object_Vector, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotVec, w.wantVec) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotVec, err := s.GetObject(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotVec, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_getObjects(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Object_VectorRequest
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
		wantVecs *sync.Map
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *sync.Map, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotVecs *sync.Map, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotVecs, w.wantVecs) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVecs, w.wantVecs)
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
		           req:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			gotVecs, err := s.getObjects(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotVecs, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamGetObject(t *testing.T) {
	type args struct {
		stream vald.Object_StreamGetObjectServer
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
		           stream:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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
		           stream:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           mirror:nil,
		           vAddr:"",
		           streamConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServerWithMirror:nil,
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

			err := s.StreamGetObject(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
