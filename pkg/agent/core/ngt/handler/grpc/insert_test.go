package grpc

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_Insert(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *payload.Insert_Request
	}
	type fields struct {
		name              string
		ip                string
		streamConcurrency int
		svcCfg            *config.NGT
		svcOpts           []service.Option
	}
	type want struct {
		wantRes *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*server)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err.Error(), w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}

	// common variables for test
	const (
		name      = "vald-agent-ngt-1" // agent name
		id        = "uuid-1"           // insert request id
		intVecDim = 3                  // int vector dimension
		f32VecDim = 3                  // float32 vector dimension
	)
	var (
		ip     = net.LoadLocalIP()        // agent ip address
		intVec = []float32{1, 2, 3}       // int vector of the insert request
		f32Vec = []float32{1.5, 2.3, 3.6} // float32 vector of the insert request

		// default NGT configuration for test
		kvsdbCfg  = &config.KVSDB{}
		vqueueCfg = &config.VQueue{}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*
		- Equivalence Class Testing
			- uint8, float32
				- case 1.1: Insert vector success (vector type is uint8)
				- case 1.2: Insert vector success (vector type is float32)
				- case 2.1: Insert vector with different dimension (vector type is uint8)
				- case 2.2: Insert vector with different dimension (vector type is float32)
				- case 3.1: Insert gaussian distributed vector success (vector type is uint8)
				- case 3.2: Insert gaussian distributed vector success (vector type is float32)
				- case 4.1: Insert uniform distributed vector success (vector type is uint8)
				- case 4.2: Insert uniform distributed vector success (vector type is float32)

		- Boundary Value Testing
			- uint8, float32
				- case 1.1: Insert vector with 0 value success (vector type is uint8)
				- case 1.1: Insert vector with 0 value success (vector type is float32)
				- case 2.1: Insert vector with min value success (vector type is uint8)
				- case 2.2: Insert vector with min value success (vector type is float32)
				- case 3.1: Insert vector with max value success (vector type is uint8)
				- case 3.2: Insert vector with max value success (vector type is float32)
				- case 4.1: Insert with empty UUID fail (vector type is uint8)
				- case 4.2: Insert with empty UUID fail (vector type is float32)

			- float32
				- case 5: Insert vector with NaN value fail (vector type is float32)

			- case 6: Insert nil insert request fail
				* IncompatibleDimensionSize error will be returned.
			- case 7: Insert nil vector fail
				* IncompatibleDimensionSize error will be returned.
			- case 8: Insert empty insert vector fail
				* IncompatibleDimensionSize error will be returned.

		- Decision Table Testing
			- duplicated ID, duplicated vector, duplicated ID & vector
				- case 1.1: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID)
					* AlreadyExists error will be returned.
				- case 1.2: Insert duplicated request success when SkipStrictExistCheck is false (duplicated vector)
				- case 1.3: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID & vector)
				- case 2.1: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID)
					* SkipStrictExistCheck flag is not used in agent handler, so the result is same as case 1.
				- case 2.2: Insert duplicated request success when SkipStrictExistCheck is true (duplicated vector)
				- case 2.3: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID & vector)
	*/
	tests := []test{
		// Equivalence Class Testing
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
			}

			return test{
				name: "Equivalence Class Testing case 1.1: Insert vector success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: f32Vec,
				},
			}

			return test{
				name: "Equivalence Class Testing case 1.2: Insert vector success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			invalidDim := intVecDim + 1
			ivec, err := vector.GenUint8Vec(vector.Gaussian, 1, invalidDim)
			if err != nil {
				t.Error(err)
			}
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec[0],
				},
			}

			return test{
				name: "Equivalence Class Testing case 2.1: Insert vector with different dimension (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(ivec), 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			invalidDim := f32VecDim + 1
			ivec, err := vector.GenF32Vec(vector.Gaussian, 1, invalidDim)
			if err != nil {
				t.Error(err)
			}
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec[0],
				},
			}

			return test{
				name: "Equivalence Class Testing case 2.2: Insert vector with different dimension (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(ivec), 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			ivec, err := vector.GenUint8Vec(vector.Gaussian, 1, intVecDim)
			if err != nil {
				t.Error(err)
			}

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec[0],
				},
			}

			return test{
				name: "Equivalence Class Testing case 3.1: Insert gaussian distributed vector success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			ivec, err := vector.GenF32Vec(vector.Gaussian, 1, f32VecDim)
			if err != nil {
				t.Error(err)
			}

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec[0],
				},
			}

			return test{
				name: "Equivalence Class Testing case 3.2: Insert gaussian distributed vector success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			ivec, err := vector.GenUint8Vec(vector.Uniform, 1, intVecDim)
			if err != nil {
				t.Error(err)
			}

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec[0],
				},
			}

			return test{
				name: "Equivalence Class Testing case 4.1: Insert uniform distributed vector success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			ivec, err := vector.GenF32Vec(vector.Uniform, 1, f32VecDim)
			if err != nil {
				t.Error(err)
			}

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec[0],
				},
			}

			return test{
				name: "Equivalence Class Testing case 4.2: Insert uniform distributed vector success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),

		// Boundary Value Testing
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vector.GenSameValueVec(intVecDim, 0),
				},
			}

			return test{
				name: "Boundary Value Testing case 1.1: Insert vector with 0 value success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vector.GenSameValueVec(f32VecDim, 0),
				},
			}

			return test{
				name: "Boundary Value Testing case 1.2: Insert vector with 0 value success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vector.GenSameValueVec(intVecDim, math.MinInt),
				},
			}

			return test{
				name: "Boundary Value Testing case 2.1: Insert vector with min value success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vector.GenSameValueVec(f32VecDim, -math.MaxFloat32),
				},
			}

			return test{
				name: "Boundary Value Testing case 2.2: Insert vector with min value success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vector.GenSameValueVec(intVecDim, math.MaxInt),
				},
			}

			return test{
				name: "Boundary Value Testing case 3.1: Insert vector with max value success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vector.GenSameValueVec(f32VecDim, math.MaxFloat32),
				},
			}

			return test{
				name: "Boundary Value Testing case 3.2: Insert vector with max value success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     "",
					Vector: intVec,
				},
			}

			return test{
				name: "Boundary Value Testing case 4.1: Insert with empty UUID fail (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrUUIDNotFound(0)
						err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API empty uuid \"%s\" was given", req.GetVector().GetId()), err,
							&errdetails.RequestInfo{
								RequestId:   req.GetVector().GetId(),
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "uuid",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     "",
					Vector: f32Vec,
				},
			}

			return test{
				name: "Boundary Value Testing case 4.2: Insert with empty UUID fail (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrUUIDNotFound(0)
						err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API empty uuid \"%s\" was given", req.GetVector().GetId()), err,
							&errdetails.RequestInfo{
								RequestId:   req.GetVector().GetId(),
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "uuid",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     "",
					Vector: f32Vec,
				},
			}

			return test{
				name: "Boundary Value Testing case 4.2: Insert with empty UUID fail (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrUUIDNotFound(0)
						err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API empty uuid \"%s\" was given", req.GetVector().GetId()), err,
							&errdetails.RequestInfo{
								RequestId:   req.GetVector().GetId(),
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "uuid",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			nan := float32(math.NaN())
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vector.GenSameValueVec(f32VecDim, nan),
				},
			}

			return test{
				name: "Boundary Value Testing case 5: Insert vector with NaN value fail (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "Boundary Value Testing case 6: Insert nil insert request fail",
				args: args{
					ctx: ctx,
					req: nil,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					// IncompatibleDimensionSize error will be returned
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   "",
								ServingData: errdetails.Serialize(nil),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: nil,
				},
			}

			return test{
				name: "Boundary Value Testing case 7: Insert nil vector fail",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					// IncompatibleDimensionSize error will be returned
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   id,
								ServingData: errdetails.Serialize(nil),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: []float32{},
				},
			}

			return test{
				name: "Boundary Value Testing case 8: Insert empty insert vector fail",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					// IncompatibleDimensionSize error will be returned
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   id,
								ServingData: errdetails.Serialize(nil),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),

		// Decision Table Testing
		func() test {
			bVecs, err := vector.GenUint8Vec(vector.Gaussian, 1, intVecDim) // used in beforeFunc
			if err != nil {
				t.Error(err)
			}

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: false,
				},
			}

			return test{
				name: "Decision Table Testing case 1.1: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, bVecs[0])
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			bId := "uuid-2" // use in beforeFunc

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: false,
				},
			}

			return test{
				name: "Decision Table Testing case 1.2: Insert duplicated request success when SkipStrictExistCheck is false (duplicated vector)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(bId, intVec)
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: false,
				},
			}

			return test{
				name: "Decision Table Testing case 1.3: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID & vector)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, intVec)
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			bVec, err := vector.GenUint8Vec(vector.Gaussian, 1, intVecDim) // use in beforeFunc
			if err != nil {
				t.Error(err)
			}

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
				},
			}

			return test{
				name: "Decision Table Testing case 2.1: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, bVec[0])
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			bId := "uuid-2" // use in beforeFunc

			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
				},
			}

			return test{
				name: "Decision Table Testing case 2.2: Insert duplicated request success when SkipStrictExistCheck is true (duplicated vector)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(bId, intVec)
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
				},
			}

			return test{
				name: "Decision Table Testing case 2.3: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID & vector)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, intVec)
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(test.fields.svcCfg, append(test.fields.svcOpts, service.WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}

			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               ngt,
				eg:                eg,
				streamConcurrency: test.fields.streamConcurrency,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(s)
			}

			gotRes, err := s.Insert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamInsert(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Insert_StreamInsertServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamInsert(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiInsert(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		reqs *payload.Insert_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		streamConcurrency int
		svcCfg            *config.NGT
		svcOpts           []service.Option
	}
	type want struct {
		wantRes    *payload.Object_Locations
		err        error
		containErr []error // check the function output error contain one of the error or not
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(*testing.T, *server)
		afterFunc  func(args)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// common variables for test
	const (
		name      = "vald-agent-ngt-1" // agent name
		id        = "uuid-1"           // insert request id
		intVecDim = 3                  // int vector dimension
		f32VecDim = 3                  // float32 vector dimension
		maxVecDim = 1 << 18            // reference value for testing, this value is temporary
	)
	var (
		ip = net.LoadLocalIP() // agent ip address

		// default NGT configuration for test
		defaultIntSvcCfg = &config.NGT{
			Dimension:    intVecDim,
			DistanceType: ngt.Angle.String(),
			ObjectType:   ngt.Uint8.String(),
			KVSDB:        &config.KVSDB{},
			VQueue:       &config.VQueue{},
		}
		defaultF32SvcCfg = &config.NGT{
			Dimension:    f32VecDim,
			DistanceType: ngt.Angle.String(),
			ObjectType:   ngt.Float.String(),
			KVSDB:        &config.KVSDB{},
			VQueue:       &config.VQueue{},
		}
		defaultSvcOpts = []service.Option{
			service.WithEnableInMemoryMode(true),
		}
	)

	genAlreadyExistsErr := func(uuid string, req *payload.Insert_MultiRequest, name, ip string) error {
		return status.WrapWithAlreadyExists(fmt.Sprintf("MultiInsert API uuids [%v] already exists", uuid),
			errors.ErrUUIDAlreadyExists(uuid),
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiInsert",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
			})
	}

	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
		if w.containErr == nil {
			if !errors.Is(err, w.err) {
				return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", err, w.err)
			}
		} else {
			exist := false
			for _, e := range w.containErr {
				if errors.Is(err, e) {
					exist = true
					break
				}
			}
			if !exist {
				return errors.Errorf("got_error: \"%v\",\n\t\t\t\tshould contain one of the error: \"%v\"", err, w.containErr)
			}
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}

	/*
		- Equivalence Class Testing
			- uint8, float32
				- case 1.1: Success to MultiInsert 1 vector (vector type is uint8)
				- case 1.2: Success to MultiInsert 1 vector (vector type is float32)
				- case 1.3: Success to MultiInsert 100 vector (vector type is uint8)
				- case 1.4: Success to MultiInsert 100 vector (vector type is float32)
				- case 1.5: Success to MultiInsert 0 vector (vector type is uint8)
				- case 1.6: Success to MultiInsert 0 vector (vector type is float32)
				- case 2.1: Fail to MultiInsert 1 vector with different dimension (vector type is uint8)
				- case 2.2: Fail to MultiInsert 1 vector with different dimension (vector type is float32)
				- case 3.1: Fail to MultiInsert 100 vector with 1 vector with different dimension (vector type is uint8)
				- case 3.2: Fail to MultiInsert 100 vector with 1 vector with different dimension (vector type is float32)
				- case 3.3: Fail to MultiInsert 100 vector with 50 vector with different dimension (vector type is uint8)
				- case 3.4: Fail to MultiInsert 100 vector with 50 vector with different dimension (vector type is float32)
				- case 3.5: Fail to MultiInsert 100 vector with all vector with different dimension (vector type is uint8)
				- case 3.6: Fail to MultiInsert 100 vector with all vector with different dimension (vector type is float32)

		- Boundary Value Testing
			- uint8, float32 (with 100 insert request in a single MultiInsert request)
				- case 1.1: Success to MultiInsert with 0 value vector (vector type is uint8)
				- case 1.2: Success to MultiInsert with 0 value vector (vector type is float32)
				- case 2.1: Success to MultiInsert with min value vector (vector type is uint8)
				- case 2.2: Success to MultiInsert with min value vector (vector type is float32)
				- case 3.1: Success to MultiInsert with max value vector (vector type is uint8)
				- case 3.2: Success to MultiInsert with max value vector (vector type is float32)
				- case 4.1: Fail to MultiInsert with 1 request with empty UUID (vector type is uint8)
				- case 4.2: Fail to MultiInsert with 1 request with empty UUID (vector type is float32)
				- case 4.3: Fail to MultiInsert with 50 request with empty UUID (vector type is uint8)
				- case 4.4: Fail to MultiInsert with 50 request with empty UUID (vector type is float32)
				- case 4.5: Fail to MultiInsert with all request with empty UUID (vector type is uint8)
				- case 4.6: Fail to MultiInsert with all request with empty UUID (vector type is float32)
				- case 5.1: Fail to MultiInsert with 1 vector with maximum dimension (vector type is uint8)
				- case 5.2: Fail to MultiInsert with 1 vector with maximum dimension (vector type is float32)
				- case 5.3: Fail to MultiInsert with 50 vector with maximum dimension (vector type is uint8)
				- case 5.4: Fail to MultiInsert with 50 vector with maximum dimension (vector type is float32)
				- case 5.5: Fail to MultiInsert with all vector with maximum dimension (vector type is uint8)
				- case 5.6: Fail to MultiInsert with all vector with maximum dimension (vector type is float32)

			- float32 (with 100 insert request in a single MultiInsert request)
				- case 6.1: Success to MultiInsert with NaN value (vector type is float32)
				- case 6.2: Success to MultiInsert with +Inf value (vector type is float32)
				- case 6.3: Success to MultiInsert with -Inf value (vector type is float32)
				- case 6.4: Success to MultiInsert with -0 value (vector type is float32)

			- others  (with 100 insert request in a single MultiInsert request)
				- case 7.1: Fail to MultiInsert with 1 vector with nil insert request
				- case 7.2: Fail to MultiInsert with 50 vector with nil insert request
				- case 7.3: Fail to MultiInsert with all vector with nil insert request
				- case 8.1: Fail to MultiInsert with 1 vector with nil vector
				- case 8.2: Fail to MultiInsert with 50 vector with nil vector
				- case 8.3: Fail to MultiInsert with all vector with nil vector
				- case 9.1: Fail to MultiInsert with 1 vector with empty insert vector
				- case 9.2: Fail to MultiInsert with 50 vector with empty insert vector
				- case 9.3: Fail to MultiInsert with all vector with empty insert vector

		- Decision Table Testing
			- duplicated ID (with 100 insert request in a single MultiInsert request)
				- case 1.1: Success to MultiInsert with 2 duplicated ID when SkipStrictExistCheck is false
				- case 1.2: Success to MultiInsert with all duplicated ID when SkipStrictExistCheck is false
				- case 1.3: Success to MultiInsert with 2 duplicated ID when SkipStrictExistCheck is true
				- case 1.4: Success to MultiInsert with all duplicated ID when SkipStrictExistCheck is true
			- duplicated vector (with 100 insert request in a single MultiInsert request)
				- case 2.1: Success to MultiInsert with 2 duplicated vector when SkipStrictExistCheck is false
				- case 2.2: Success to MultiInsert with all duplicated vector when SkipStrictExistCheck is false
				- case 2.3: Success to MultiInsert with 2 duplicated vector when SkipStrictExistCheck is true
				- case 2.4: Success to MultiInsert with all duplicated vector when SkipStrictExistCheck is true
			- duplicated ID & duplicated vector (with 100 insert request in a single MultiInsert request)
				- case 3.1: Success to MultiInsert with 2 duplicated ID & vector when SkipStrictExistCheck is false
				- case 3.2: Success to MultiInsert with all duplicated ID & vector when SkipStrictExistCheck is false
				- case 3.3: Success to MultiInsert with 2 duplicated ID & vector when SkipStrictExistCheck is true
				- case 3.4: Success to MultiInsert with all duplicated ID & vector when SkipStrictExistCheck is true

			// existed in NGT test cases
			- existed ID (with 100 insert request in a single MultiInsert request)
				- case 4.1: Fail to MultiInsert with 2 existed ID when SkipStrictExistCheck is false
				- case 4.2: Fail to MultiInsert with all existed ID when SkipStrictExistCheck is false
				- case 4.3: Fail to MultiInsert with 2 existed ID when SkipStrictExistCheck is true
				- case 4.4: Fail to MultiInsert with all existed ID when SkipStrictExistCheck is true
			- existed vector (with 100 insert request in a single MultiInsert request)
				- case 5.1: Success to MultiInsert with 2 existed vector when SkipStrictExistCheck is false
				- case 5.2: Success to MultiInsert with all existed vector when SkipStrictExistCheck is false
				- case 5.3: Success to MultiInsert with 2 existed vector when SkipStrictExistCheck is true
				- case 5.4: Success to MultiInsert with all existed vector when SkipStrictExistCheck is true
			- existed ID & existed vector (with 100 insert request in a single MultiInsert request)
				- case 6.1: Fail to MultiInsert with 2 existed ID & vector when SkipStrictExistCheck is false
				- case 6.2: Fail to MultiInsert with all existed ID & vector when SkipStrictExistCheck is false
				- case 6.3: Fail to MultiInsert with 2 existed ID & vector when SkipStrictExistCheck is true
				- case 6.4: Fail to MultiInsert with all existed ID & vector when SkipStrictExistCheck is true

	*/
	tests := []test{
		func() test {
			insertNum := 1
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 1.1: Success to MultiInsert 1 vector (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 1
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 1.2: Success to MultiInsert 1 vector (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 1.3: Success to MultiInsert 100 vector (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 1.4: Success to MultiInsert 100 vector (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		{
			name: "Equivalence Class Testing case 1.5: Success to MultiInsert 0 vector (vector type is uint8)",
			args: args{
				ctx: ctx,
				reqs: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{},
				},
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultIntSvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: nil,
			},
		},
		{
			name: "Equivalence Class Testing case 1.6: Success to MultiInsert 0 vector (vector type is float32)",
			args: args{
				ctx: ctx,
				reqs: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{},
				},
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: nil,
			},
		},
		func() test {
			insertNum := 1
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim+1, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 2.1: Fail to MultiInsert 1 vector with different dimension (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultIntSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 1
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim+1, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 2.2: Fail to MultiInsert 1 vector with different dimension (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}

			invalidVecs, err := vector.GenUint8Vec(vector.Gaussian, 1, intVecDim+1)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = invalidVecs[0]

			return test{
				name: "Equivalence Class Testing case 3.1: Fail to MultiInsert 100 vector with 1 vector with different dimension (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultIntSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			invalidVecs, err := vector.GenF32Vec(vector.Gaussian, 1, f32VecDim+1)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = invalidVecs[0]

			return test{
				name: "Equivalence Class Testing case 3.2: Fail to MultiInsert 100 vector with 1 vector with different dimension (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),

		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}

			invalidCnt := len(req.Requests) / 2
			invalidVec, err := vector.GenUint8Vec(vector.Gaussian, invalidCnt, intVecDim+1)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < invalidCnt; i++ {
				req.Requests[i].Vector.Vector = invalidVec[i]
			}

			return test{
				name: "Equivalence Class Testing case 3.3: Fail to MultiInsert 100 vector with 50 vector with different dimension (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultIntSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			invalidCnt := len(req.Requests) / 2
			invalidVec, err := vector.GenF32Vec(vector.Gaussian, invalidCnt, f32VecDim+1)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < invalidCnt; i++ {
				req.Requests[i].Vector.Vector = invalidVec[i]
			}

			return test{
				name: "Equivalence Class Testing case 3.4: Fail to MultiInsert 100 vector with 50 vector with different dimension (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim+1, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 3.5: Fail to MultiInsert 100 vector with all vector with different dimension (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim+1, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 3.6: Fail to MultiInsert 100 vector with all vector with different dimension (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 1.1: Success to MultiInsert with 0 value vector (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(intVecDim, 0), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 1.2: Success to MultiInsert with 0 value vector (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(f32VecDim, 0), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 2.1: Success to MultiInsert with min value vector (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(intVecDim, math.MinInt), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 2.2: Success to MultiInsert with min value vector (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(f32VecDim, -math.MaxFloat32), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 3.1: Success to MultiInsert with max value vector (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(intVecDim, math.MaxUint8), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 3.2: Success to MultiInsert with max value vector (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(f32VecDim, math.MaxFloat32), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Id = ""

			uuids := make([]string, 0, len(req.Requests))
			for _, r := range req.Requests {
				uuids = append(uuids, r.Vector.Id)
			}

			return test{
				name: "Boundary Value Testing case 4.1: Fail to MultiInsert with 1 request with empty UUID (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: status.WrapWithInvalidArgument(fmt.Sprintf("MultiInsert API invalid uuids \"%v\" detected", uuids), errors.ErrUUIDNotFound(0),
						&errdetails.RequestInfo{
							RequestId:   strings.Join(uuids, ", "),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequestFieldViolation{
								{
									Field:       "uuid",
									Description: errors.ErrUUIDNotFound(0).Error(),
								},
							},
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.MultiInsert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Id = ""

			uuids := make([]string, 0, len(req.Requests))
			for _, r := range req.Requests {
				uuids = append(uuids, r.Vector.Id)
			}

			return test{
				name: "Boundary Value Testing case 4.2: Fail to MultiInsert with 1 request with empty UUID (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: status.WrapWithInvalidArgument(fmt.Sprintf("MultiInsert API invalid uuids \"%v\" detected", uuids), errors.ErrUUIDNotFound(0),
						&errdetails.RequestInfo{
							RequestId:   strings.Join(uuids, ", "),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequestFieldViolation{
								{
									Field:       "uuid",
									Description: errors.ErrUUIDNotFound(0).Error(),
								},
							},
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.MultiInsert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < len(req.Requests)/2; i++ {
				req.Requests[i].Vector.Id = ""
			}

			uuids := make([]string, 0, len(req.Requests))
			for _, r := range req.Requests {
				uuids = append(uuids, r.Vector.Id)
			}

			return test{
				name: "Boundary Value Testing case 4.3: Fail to MultiInsert with 50 request with empty UUID (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: status.WrapWithInvalidArgument(fmt.Sprintf("MultiInsert API invalid uuids \"%v\" detected", uuids), errors.ErrUUIDNotFound(0),
						&errdetails.RequestInfo{
							RequestId:   strings.Join(uuids, ", "),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequestFieldViolation{
								{
									Field:       "uuid",
									Description: errors.ErrUUIDNotFound(0).Error(),
								},
							},
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.MultiInsert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < len(req.Requests)/2; i++ {
				req.Requests[i].Vector.Id = ""
			}

			uuids := make([]string, 0, len(req.Requests))
			for _, r := range req.Requests {
				uuids = append(uuids, r.Vector.Id)
			}

			return test{
				name: "Boundary Value Testing case 4.4: Fail to MultiInsert with 50 request with empty UUID (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: status.WrapWithInvalidArgument(fmt.Sprintf("MultiInsert API invalid uuids \"%v\" detected", uuids), errors.ErrUUIDNotFound(0),
						&errdetails.RequestInfo{
							RequestId:   strings.Join(uuids, ", "),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequestFieldViolation{
								{
									Field:       "uuid",
									Description: errors.ErrUUIDNotFound(0).Error(),
								},
							},
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.MultiInsert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = ""
			}

			uuids := make([]string, 0, len(req.Requests))
			for _, r := range req.Requests {
				uuids = append(uuids, r.Vector.Id)
			}

			return test{
				name: "Boundary Value Testing case 4.5: Fail to MultiInsert with all request with empty UUID (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: status.WrapWithInvalidArgument(fmt.Sprintf("MultiInsert API invalid uuids \"%v\" detected", uuids), errors.ErrUUIDNotFound(0),
						&errdetails.RequestInfo{
							RequestId:   strings.Join(uuids, ", "),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequestFieldViolation{
								{
									Field:       "uuid",
									Description: errors.ErrUUIDNotFound(0).Error(),
								},
							},
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.MultiInsert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = ""
			}

			uuids := make([]string, 0, len(req.Requests))
			for _, r := range req.Requests {
				uuids = append(uuids, r.Vector.Id)
			}

			return test{
				name: "Boundary Value Testing case 4.6: Fail to MultiInsert with all request with empty UUID (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: status.WrapWithInvalidArgument(fmt.Sprintf("MultiInsert API invalid uuids \"%v\" detected", uuids), errors.ErrUUIDNotFound(0),
						&errdetails.RequestInfo{
							RequestId:   strings.Join(uuids, ", "),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequestFieldViolation{
								{
									Field:       "uuid",
									Description: errors.ErrUUIDNotFound(0).Error(),
								},
							},
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.MultiInsert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = make([]float32, maxVecDim)

			return test{
				name: "Boundary Value Testing case 5.1: Fail to MultiInsert with 1 vector with maximum dimension (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = make([]float32, maxVecDim)

			return test{
				name: "Boundary Value Testing case 5.1: Fail to MultiInsert with 1 vector with maximum dimension (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, intVecDim, nil)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < len(req.Requests)/2; i++ {
				req.Requests[i].Vector.Vector = make([]float32, maxVecDim)
			}

			return test{
				name: "Boundary Value Testing case 5.3: Fail to MultiInsert with 50 vector with maximum dimension (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			for i := 0; i < len(req.Requests)/2; i++ {
				req.Requests[i].Vector.Vector = make([]float32, maxVecDim)
			}

			return test{
				name: "Boundary Value Testing case 5.4: Fail to MultiInsert with 50 vector with maximum dimension (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req := request.GenSameVecMultiInsertReq(insertNum, make([]float32, maxVecDim), nil)

			return test{
				name: "Boundary Value Testing case 5.5: Fail to MultiInsert with all vector with maximum dimension (vector type is uint8)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultIntSvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req := request.GenSameVecMultiInsertReq(insertNum, make([]float32, maxVecDim), nil)

			return test{
				name: "Boundary Value Testing case 5.6: Fail to MultiInsert with all vector with maximum dimension (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[0].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[0].Vector.Id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 6.1: Success to MultiInsert with NaN value (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(f32VecDim, float32(math.NaN())), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 6.2: Success to MultiInsert with +Inf value (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(f32VecDim, float32(math.Inf(+1.0))), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 6.3: Success to MultiInsert with -Inf value (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(f32VecDim, float32(math.Inf(-1.0))), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			return test{
				name: "Boundary Value Testing case 6.4: Success to MultiInsert with -0 value (vector type is float32)",
				args: args{
					ctx:  ctx,
					reqs: request.GenSameVecMultiInsertReq(insertNum, vector.GenSameValueVec(f32VecDim, float32(math.Copysign(0, -1.0))), nil),
				},
				fields: fields{
					name:              name,
					ip:                ip,
					svcCfg:            defaultF32SvcCfg,
					svcOpts:           defaultSvcOpts,
					streamConcurrency: 0,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			req.Requests[0] = nil

			return test{
				name: "Boundary Value Testing case 7.1: Fail to MultiInsert with 1 vector with nil insert request",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests)/2; i++ {
				req.Requests[i] = nil
			}

			return test{
				name: "Boundary Value Testing case 7.2: Fail to MultiInsert with 50 vector with nil insert request",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i] = nil
			}

			return test{
				name: "Boundary Value Testing case 7.3: Fail to MultiInsert with all vector with nil insert request",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			req.Requests[0].Vector.Vector = nil

			return test{
				name: "Boundary Value Testing case 8.1: Fail to MultiInsert with 1 vector with nil vector",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests)/2; i++ {
				req.Requests[i].Vector.Vector = nil
			}

			return test{
				name: "Boundary Value Testing case 8.2: Fail to MultiInsert with 50 vector with nil vector",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Vector = nil
			}

			return test{
				name: "Boundary Value Testing case 8.3: Fail to MultiInsert with all vector with nil vector",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			req.Requests[0].Vector.Vector = []float32{}

			return test{
				name: "Boundary Value Testing case 9.1: Fail to MultiInsert with 1 vector with empty insert vector",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests)/2; i++ {
				req.Requests[i].Vector.Vector = []float32{}
			}

			return test{
				name: "Boundary Value Testing case 9.2: Fail to MultiInsert with 50 vector with empty insert vector",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Vector = []float32{}
			}

			return test{
				name: "Boundary Value Testing case 9.3: Fail to MultiInsert with all vector with empty insert vector",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   vid,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.MultiInsert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			w.Locations[0].Uuid = req.Requests[0].Vector.Id

			return test{
				name: "Decision Table Testing case 1.1: Success to MultiInsert with 2 duplicated ID when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
			}

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			for _, l := range w.Locations {
				l.Uuid = req.Requests[0].Vector.Id
			}

			return test{
				name: "Decision Table Testing case 1.2: Success to MultiInsert with all duplicated ID when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			w.Locations[0].Uuid = req.Requests[0].Vector.Id
			// w.Locations[1].Uuid = dupID

			return test{
				name: "Decision Table Testing case 1.3: Success to MultiInsert with 2 duplicated ID when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
			}

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			for _, l := range w.Locations {
				l.Uuid = req.Requests[0].Vector.Id
			}

			return test{
				name: "Decision Table Testing case 1.4: Success to MultiInsert with all duplicated ID when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector

			return test{
				name: "Decision Table Testing case 2.1: Success to MultiInsert with 2 duplicated vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			return test{
				name: "Decision Table Testing case 2.2: Success to MultiInsert with all duplicated vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector

			return test{
				name: "Decision Table Testing case 2.3: Success to MultiInsert with 2 duplicated vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			return test{
				name: "Decision Table Testing case 2.4: Success to MultiInsert with all duplicated vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			w.Locations[0].Uuid = req.Requests[0].Vector.Id

			return test{
				name: "Decision Table Testing case 3.1: Success to MultiInsert with 2 duplicated ID & vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			for _, l := range w.Locations {
				l.Uuid = req.Requests[0].Vector.Id
			}

			return test{
				name: "Decision Table Testing case 3.2: Success to MultiInsert with all duplicated ID & vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			w.Locations[0].Uuid = req.Requests[0].Vector.Id

			return test{
				name: "Decision Table Testing case 3.3: Success to MultiInsert with 2 duplicated ID & vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, cfg)
			if err != nil {
				t.Error(err)
			}
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			// set want
			w := request.GenObjectLocations(insertNum, name, ip)
			for _, l := range w.Locations {
				l.Uuid = req.Requests[0].Vector.Id
			}

			return test{
				name: "Decision Table Testing case 3.4: Success to MultiInsert with all duplicated ID & vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				want: want{
					wantRes: w,
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 4.1: Fail to MultiInsert with 2 existed ID when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					vecs, err := vector.GenF32Vec(vector.Gaussian, 2, f32VecDim)
					if err != nil {
						t.Error(err)
					}
					for i := 0; i < 2; i++ {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     req.Requests[i].Vector.Id,
								Vector: vecs[i],
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: false,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}

					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 2,
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: []error{
						genAlreadyExistsErr(req.Requests[0].Vector.Id, req, name, ip),
						genAlreadyExistsErr(req.Requests[1].Vector.Id, req, name, ip),
					},
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			wantErrs := make([]error, 100)
			for i := 0; i < len(req.Requests); i++ {
				wantErrs[i] = genAlreadyExistsErr(req.Requests[i].Vector.Id, req, name, ip)
			}

			return test{
				name: "Decision Table Testing case 4.2: Fail to MultiInsert with all existed ID when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					vecs, err := vector.GenF32Vec(vector.Gaussian, insertNum, f32VecDim)
					if err != nil {
						t.Error(err)
					}
					for i, r := range req.Requests {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     r.Vector.Id,
								Vector: vecs[i],
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: false,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: wantErrs,
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 4.3: Fail to MultiInsert with 2 existed ID when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					vecs, err := vector.GenF32Vec(vector.Gaussian, 2, f32VecDim)
					if err != nil {
						t.Error(err)
					}
					for i := 0; i < 2; i++ {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     req.Requests[i].Vector.Id,
								Vector: vecs[i],
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: true,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}

					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 2,
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: []error{
						genAlreadyExistsErr(req.Requests[0].Vector.Id, req, name, ip),
						genAlreadyExistsErr(req.Requests[1].Vector.Id, req, name, ip),
					},
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			wantErrs := make([]error, 100)
			for i := 0; i < len(req.Requests); i++ {
				wantErrs[i] = genAlreadyExistsErr(req.Requests[i].Vector.Id, req, name, ip)
			}

			return test{
				name: "Decision Table Testing case 4.4: Fail to MultiInsert with all existed ID when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					vecs, err := vector.GenF32Vec(vector.Gaussian, insertNum, f32VecDim)
					if err != nil {
						t.Error(err)
					}
					for i, r := range req.Requests {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     r.Vector.Id,
								Vector: vecs[i],
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: true,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: wantErrs,
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 5.1: Success to MultiInsert with 2 existed vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					// insert same request with different ID
					for i := 0; i < 2; i++ {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     fmt.Sprintf("nonexistid%d", i),
								Vector: req.Requests[i].Vector.Vector,
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: false,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantRes: request.GenObjectLocations(100, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 5.2: Success to MultiInsert with all existed vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					// insert same request with different ID
					for i := range req.Requests {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     fmt.Sprintf("nonexistid%d", i),
								Vector: req.Requests[i].Vector.Vector,
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: false,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 5.3: Success to MultiInsert with 2 existed vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					// insert same request with different ID
					for i := 0; i < 2; i++ {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     fmt.Sprintf("nonexistid%d", i),
								Vector: req.Requests[i].Vector.Vector,
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: true,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 5.4: Success to MultiInsert with all existed vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					// insert same request with different ID
					for i := range req.Requests {
						ir := &payload.Insert_Request{
							Vector: &payload.Object_Vector{
								Id:     fmt.Sprintf("nonexistid%d", i),
								Vector: req.Requests[i].Vector.Vector,
							},
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: true,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					wantRes: request.GenObjectLocations(insertNum, name, ip),
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 6.1: Fail to MultiInsert with 2 existed ID & vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					for i := 0; i < 2; i++ {
						ir := &payload.Insert_Request{
							Vector: req.Requests[i].Vector,
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: false,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}

					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 2,
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: []error{
						genAlreadyExistsErr(req.Requests[0].Vector.Id, req, name, ip),
						genAlreadyExistsErr(req.Requests[1].Vector.Id, req, name, ip),
					},
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			wantErrs := make([]error, 100)
			for i := 0; i < len(req.Requests); i++ {
				wantErrs[i] = genAlreadyExistsErr(req.Requests[i].Vector.Id, req, name, ip)
			}

			return test{
				name: "Decision Table Testingcase 6.2: Fail to MultiInsert with all existed ID & vector when SkipStrictExistCheck is false",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					for _, r := range req.Requests {
						ir := &payload.Insert_Request{
							Vector: r.Vector,
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: false,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: wantErrs,
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 6.3: Fail to MultiInsert with 2 existed ID & vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					for i := 0; i < 2; i++ {
						ir := &payload.Insert_Request{
							Vector: req.Requests[i].Vector,
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: true,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}

					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 2,
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: []error{
						genAlreadyExistsErr(req.Requests[0].Vector.Id, req, name, ip),
						genAlreadyExistsErr(req.Requests[1].Vector.Id, req, name, ip),
					},
				},
			}
		}(),
		func() test {
			insertNum := 100
			req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, f32VecDim, nil)
			if err != nil {
				t.Error(err)
			}

			wantErrs := make([]error, 100)
			for i := 0; i < len(req.Requests); i++ {
				wantErrs[i] = genAlreadyExistsErr(req.Requests[i].Vector.Id, req, name, ip)
			}

			return test{
				name: "Decision Table Testing case 6.4: Fail to MultiInsert with all existed ID & vector when SkipStrictExistCheck is true",
				args: args{
					ctx:  ctx,
					reqs: req,
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  defaultF32SvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, s *server) {
					for _, r := range req.Requests {
						ir := &payload.Insert_Request{
							Vector: r.Vector,
							Config: &payload.Insert_Config{
								SkipStrictExistCheck: true,
							},
						}
						if _, err := s.Insert(ctx, ir); err != nil {
							t.Fatal(err)
						}
					}
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(len(req.Requests)),
					}); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					containErr: wantErrs,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(test.fields.svcCfg, append(test.fields.svcOpts, service.WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}

			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               ngt,
				eg:                eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, s)
			}

			gotRes, err := s.MultiInsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
