//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		want Server
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Server, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Server, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant_error: \"%#v\"", err, w.err)
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_newLocations(t *testing.T) {
	t.Parallel()
	type args struct {
		uuids []string
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantLocs *payload.Object_Locations
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations) error {
		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           uuids: nil,
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
		           uuids: nil,
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

			gotLocs := s.newLocations(test.args.uuids...)
			if err := checkFunc(test.want, gotLocs); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_newLocation(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		want *payload.Object_Location
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *payload.Object_Location) error {
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
		           uuid: "",
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
		           uuid: "",
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

			got := s.newLocation(test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Exists(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		uid *payload.Object_ID
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_ID
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_ID, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_ID, err error) error {
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
		           ctx: nil,
		           uid: nil,
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
		           ctx: nil,
		           uid: nil,
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

			gotRes, err := s.Exists(test.args.ctx, test.args.uid)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Search(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		ctx       context.Context
		insertNum int
		req       *payload.Search_Request
	}
	type fields struct {
		gen func(int, int) [][]float32

		opts []Option

		ngtCfg  *config.NGT
		ngtOpts []service.Option
	}
	type want struct {
		resultSize int
		code       codes.Code
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(fields, args) (Server, error)
		afterFunc  func(args)
	}

	const (
		defaultDimensionSize = 32
	)

	defaultBeforeFunc := func(f fields, a args) (Server, error) {
		eg, ctx := errgroup.New(a.ctx)
		if f.ngtOpts == nil {
			f.ngtOpts = []service.Option{}
		}
		f.ngtOpts = append(f.ngtOpts, service.WithErrGroup(eg), service.WithEnableInMemoryMode(true))
		ngt, err := service.New(f.ngtCfg, f.ngtOpts...)
		if err != nil {
			return nil, err
		}
		if f.opts == nil {
			f.opts = []Option{}
		}
		f.opts = append(f.opts, WithErrGroup(eg), WithNGT(ngt))
		s, err := New(f.opts...)
		if err != nil {
			return nil, err
		}

		reqs := make([]*payload.Insert_Request, a.insertNum)
		for i, v := range f.gen(a.insertNum, f.ngtCfg.Dimension) {
			reqs[i] = &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     strconv.Itoa(i),
					Vector: v,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
				},
			}
		}
		if _, err := s.MultiInsert(ctx, &payload.Insert_MultiRequest{Requests: reqs}); err != nil {
			return nil, err
		}
		if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{PoolSize: 100}); err != nil {
			return nil, err
		}
		return s, nil
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got_code: \"%#v\",\n\t\t\t\twant: \"%#v\"", st.Code(), w.code)
			}
		}
		if gotSize := len(gotRes.GetResults()); gotSize != w.resultSize {
			return errors.Errorf("got size: \"%#v\",\n\t\t\t\twant size: \"%#v\"", gotSize, w.resultSize)
		}
		return nil
	}
	convertVectorUint8ToFloat32 := func(vector []uint8) (ret []float32) {
		ret = make([]float32, len(vector))
		for i, e := range vector {
			ret[i] = float32(e)
		}
		return
	}
	convertVectorsUint8ToFloat32 := func(vectors [][]uint8) (ret [][]float32) {
		ret = make([][]float32, 0, len(vectors))
		for _, v := range vectors {
			ret = append(ret, convertVectorUint8ToFloat32(v))
		}
		return
	}
	ngtConfig := func(dim int, objectType string) *config.NGT {
		return &config.NGT{
			Dimension:        dim,
			DistanceType:     ngt.L2.String(),
			ObjectType:       objectType,
			CreationEdgeSize: 60,
			SearchEdgeSize:   20,
			KVSDB: &config.KVSDB{
				Concurrency: 10,
			},
			VQueue: &config.VQueue{
				InsertBufferPoolSize: 1000,
				DeleteBufferPoolSize: 1000,
			},
		}
	}
	defaultSearch_Config := &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.1,
		Timeout: 1000000000,
	}
	fill := func(f float32) (v []float32) {
		v = make([]float32, defaultDimensionSize)
		for i := range v {
			v[i] = f
		}
		return
	}

	/*
		Search test cases:
		- Equivalence Class Testing
			- case 1.1: success search vector from 1000 vectors (type: uint8)
			- case 1.2: success search vector from 1000 vectors (type: float32)
			- case 2.1: fail search with different dimension vector from 1000 vectors (type: uint8)
			- case 2.2: fail search with different dimension vector from 1000 vectors (type: float32)
		- Boundary Value Testing
			- case 1.1: success search with 0 value (min value) vector from 1000 vectors (type: uint8)
			- case 1.2: success search with +0 value vector from 1000 vectors (type: float32)
			- case 1.3: success search with -0 value vector from 1000 vectors (type: float32)
			- case 2.1: success search with max value vector from 1000 vectors (type: uint8)
			- case 2.2: success search with max value vector from 1000 vectors (type: float32)
			- case 3.1: success search with min value vector from 1000 vectors (type: float32)
			- case 4.1: fail search with NaN value vector from 1000 vectors (type: float32)
			- case 5.1: fail search with Inf value vector from 1000 vectors (type: float32)
			- case 6.1: fail search with -Inf value vector from 1000 vectors (type: float32)
			- case 7.1: fail search with 0 length vector from 1000 vectors (type: uint8)
			- case 7.2: fail search with 0 length vector from 1000 vectors (type: float32)
			- case 8.1: fail search with max dimension vector from 1000 vectors (type: uint8)
			- case 8.2: fail search with max dimension vector from 1000 vectors (type: float32)
			- case 9.1: fail search with nil vector from 1000 vectors (type: uint8)
			- case 9.2: fail search with nil vector from 1000 vectors (type: float32)
		- Decision Table Testing
			- case 1.1: success search with Search_Config.Num=10 from 5 different vectors (type: uint8)
			- case 1.2: success search with Search_Config.Num=10 from 5 different vectors (type: float32)
			- case 2.1: success search with Search_Config.Num=10 from 10 different vectors (type: uint8)
			- case 2.2: success search with Search_Config.Num=10 from 10 different vectors (type: float32)
			- case 3.1: success search with Search_Config.Num=10 from 20 different vectors (type: uint8)
			- case 3.2: success search with Search_Config.Num=10 from 20 different vectors (type: float32)
			- case 4.1: success search with Search_Config.Num=10 from 5 same vectors (type: uint8)
			- case 4.2: success search with Search_Config.Num=10 from 5 same vectors (type: float32)
			- case 5.1: success search with Search_Config.Num=10 from 10 same vectors (type: uint8)
			- case 5.2: success search with Search_Config.Num=10 from 10 same vectors (type: float32)
			- case 6.1: success search with Search_Config.Num=10 from 20 same vectors (type: uint8)
			- case 6.2: success search with Search_Config.Num=10 from 20 same vectors (type: float32)
	*/
	tests := []test{
		// Equivalence Class Testing
		{
			name: "Equivalence Class Testing case 1.1: success search vector (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 1.2: success search vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail search vector with different dimension (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize+1)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Equivalence Class Testing case 2.2: fail search vector with different dimension (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize+1)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},

		// Boundary Value Testing
		{
			name: "Boundary Value Testing case 1.1: success search with 0 value (min value) vector (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(float32(uint8(0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 1.2: success search with +0 value vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(+0.0),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 1.3: success search with -0 value vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(float32(math.Copysign(0, -1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success search with max value vector (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(float32(uint8(math.MaxUint8))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success search with max value vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(math.MaxFloat32),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success search with min value vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(-math.MaxFloat32),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 4.1: fail search with NaN value vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(float32(math.NaN())),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 5.1: fail search with Inf value vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(float32(math.Inf(+1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 6.1: fail search with -Inf value vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: fill(float32(math.Inf(-1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 7.1: fail search with 0 length vector (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: []float32{},
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, "uint8"),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 7.2: fail search with 0 length vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: []float32{},
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 8.1: fail search with max dimension vector (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, math.MaxInt32>>7)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 8.2: fail search with max dimension vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, math.MaxInt32>>7)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.1: fail search with nil vector (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: nil,
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.2: fail search with nil vector (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: nil,
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},

		// Decision Table Testing
		{
			name: "Decision Table Testing case 1.1: success search with Search_Config.Num=10 from 5 different vectors (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 1.2: success search with Search_Config.Num=10 from 5 different vectors (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 2.1: success search with Search_Config.Num=10 from 10 different vectors (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 2.2: success search with Search_Config.Num=10 from 10 different vectors (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 3.1: success search with Search_Config.Num=10 from 20 different vectors (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					return convertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(n, dim))
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 3.2: success search with Search_Config.Num=10 from 20 different vectors (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen:    vector.GaussianDistributedFloat32VectorGenerator,
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 4.1: success search with Search_Config.Num=10 from 5 same vectors (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					v := convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, dim)[0])
					vectors := make([][]float32, n)
					for i := range vectors {
						vectors[i] = v
					}
					return vectors
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 4.2: success search with Search_Config.Num=10 from 5 same vectors (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					v := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
					vectors := make([][]float32, n)
					for i := range vectors {
						vectors[i] = v
					}
					return vectors
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 5.1: success search with Search_Config.Num=10 from 10 same vectors (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					v := convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, dim)[0])
					vectors := make([][]float32, n)
					for i := range vectors {
						vectors[i] = v
					}
					return vectors
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 5.2: success search with Search_Config.Num=10 from 10 same vectors (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					v := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
					vectors := make([][]float32, n)
					for i := range vectors {
						vectors[i] = v
					}
					return vectors
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 6.1: success search with Search_Config.Num=10 from 20 same vectors (type: uint8)",
			args: args{
				ctx:       ctx,
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize)[0]),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					v := convertVectorUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, dim)[0])
					vectors := make([][]float32, n)
					for i := range vectors {
						vectors[i] = v
					}
					return vectors
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 6.2: success search with Search_Config.Num=10 from 20 same vectors (type: float32)",
			args: args{
				ctx:       ctx,
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				gen: func(n, dim int) [][]float32 {
					v := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
					vectors := make([][]float32, n)
					for i := range vectors {
						vectors[i] = v
					}
					return vectors
				},
				ngtCfg: ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			s, err := test.beforeFunc(test.fields, test.args)
			if err != nil {
				tt.Errorf("error = %v", err)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotRes, err := s.Search(ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_SearchByID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		ctx      context.Context
		indexId  string
		searchId string
	}
	type want struct {
		resultSize int
		code       codes.Code
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(args) (Server, error)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got_code: \"%#v\",\n\t\t\t\twant: \"%#v\"", st.Code(), w.code)
			}
		}
		if gotSize := len(gotRes.GetResults()); gotSize != w.resultSize {
			return errors.Errorf("got size: \"%#v\",\n\t\t\t\twant size: \"%#v\"", gotSize, w.resultSize)
		}
		return nil
	}

	const (
		insertNum = 1000
	)

	defaultNgtConfig := &config.NGT{
		Dimension:        128,
		DistanceType:     ngt.L2.String(),
		ObjectType:       ngt.Float.String(),
		CreationEdgeSize: 60,
		SearchEdgeSize:   20,
		KVSDB: &config.KVSDB{
			Concurrency: 10,
		},
		VQueue: &config.VQueue{
			InsertBufferPoolSize: 1000,
			DeleteBufferPoolSize: 1000,
		},
	}
	defaultBeforeFunc := func(a args) (Server, error) {
		eg, ctx := errgroup.New(a.ctx)
		ngt, err := service.New(defaultNgtConfig, service.WithErrGroup(eg), service.WithEnableInMemoryMode(true))
		if err != nil {
			return nil, err
		}

		s, err := New(WithErrGroup(eg), WithNGT(ngt))
		if err != nil {
			return nil, err
		}

		reqs := make([]*payload.Insert_Request, insertNum)
		for i, v := range vector.GaussianDistributedFloat32VectorGenerator(insertNum, defaultNgtConfig.Dimension) {
			reqs[i] = &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     strconv.Itoa(i),
					Vector: v,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
				},
			}
		}
		reqs[0].Vector.Id = a.indexId
		if _, err := s.MultiInsert(ctx, &payload.Insert_MultiRequest{Requests: reqs}); err != nil {
			return nil, err
		}
		if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{PoolSize: 100}); err != nil {
			return nil, err
		}
		return s, nil
	}
	defaultSearch_Config := &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.1,
		Timeout: 1000000000,
	}

	utf8ToSjis := func(s string) string {
		b, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(s), japanese.ShiftJIS.NewEncoder()))
		return string(b)
	}

	utf8ToEucjp := func(s string) string {
		b, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(s), japanese.EUCJP.NewEncoder()))
		return string(b)
	}

	/*
		SearchByID test cases ( focus on ID(string), only test float32 ):
		- Equivalence Class Testing ( 1000 vectors inserted before a search )
			- case 1.1: success search vector
			- case 2.1: fail search with non-existent ID
		- Boundary Value Testing ( 1000 vectors inserted before a search )
			- case 1.1: fail search with ""
			- case 2.1: success search with ^@
			- case 2.2: success search with ^I
			- case 2.3: success search with ^J
			- case 2.4: success search with ^M
			- case 2.5: success search with ^[
			- case 2.6: success search with ^?
			- case 3.1: success search with utf-8 ID from utf-8 index
			- case 3.2: fail search with utf-8 ID from s-jis index
			- case 3.3: fail search with utf-8 ID from euc-jp index
			- case 3.4: fail search with s-jis ID from utf-8 index
			- case 3.5: success search with s-jis ID from s-jis index
			- case 3.6: fail search with s-jis ID from euc-jp index
			- case 3.4: fail search with euc-jp ID from utf-8 index
			- case 3.5: fail search with euc-jp ID from s-jis index
			- case 3.6: success search with euc-jp ID from euc-jp index
			- case 4.1: success search with 😀
		- Decision Table Testing
		    - NONE
	*/
	tests := []test{
		{
			name: "Equivalence Class Testing case 1.1: success search vector",
			args: args{
				ctx:      ctx,
				indexId:  "test",
				searchId: "test",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail search with non-existent ID",
			args: args{
				ctx:      ctx,
				indexId:  "test",
				searchId: "non-existent",
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 1.1: fail search with \"\"",
			args: args{
				ctx:      ctx,
				indexId:  "test",
				searchId: "",
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success search with ^@",
			args: args{
				ctx:      ctx,
				indexId:  string([]byte{0}),
				searchId: string([]byte{0}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: succsess search with ^I",
			args: args{
				ctx:      ctx,
				indexId:  "\t",
				searchId: "\t",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.3: succsess search with ^J",
			args: args{
				ctx:      ctx,
				indexId:  "\n",
				searchId: "\n",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.4: succsess search with ^M",
			args: args{
				ctx:      ctx,
				indexId:  "\r",
				searchId: "\r",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.5: success search with ^[",
			args: args{
				ctx:      ctx,
				indexId:  string([]byte{27}),
				searchId: string([]byte{27}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.6: success search with ^?",
			args: args{
				ctx:      ctx,
				indexId:  string([]byte{127}),
				searchId: string([]byte{127}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success search with utf-8 ID from utf-8 index",
			args: args{
				ctx:      ctx,
				indexId:  "こんにちは",
				searchId: "こんにちは",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.2: fail search with utf-8 ID from s-jis index",
			args: args{
				ctx:      ctx,
				indexId:  utf8ToSjis("こんにちは"),
				searchId: "こんにちは",
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.3: fail search with utf-8 ID from euc-jp index",
			args: args{
				ctx:      ctx,
				indexId:  utf8ToEucjp("こんにちは"),
				searchId: "こんにちは",
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.4: fail search with s-jis ID from utf-8 index",
			args: args{
				ctx:      ctx,
				indexId:  "こんにちは",
				searchId: utf8ToSjis("こんにちは"),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.5: success search with s-jis ID from s-jis index",
			args: args{
				ctx:      ctx,
				indexId:  utf8ToSjis("こんにちは"),
				searchId: utf8ToSjis("こんにちは"),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.6: fail search with s-jis ID from euc-jp index",
			args: args{
				ctx:      ctx,
				indexId:  utf8ToEucjp("こんにちは"),
				searchId: utf8ToSjis("こんにちは"),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.7: fail search with euc-jp ID from utf-8 index",
			args: args{
				ctx:      ctx,
				indexId:  "こんにちは",
				searchId: utf8ToEucjp("こんにちは"),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.8: fail search with euc-jp ID from s-jis index",
			args: args{
				ctx:      ctx,
				indexId:  utf8ToSjis("こんにちは"),
				searchId: utf8ToEucjp("こんにちは"),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.9: success search with euc-jp ID from euc-jp index",
			args: args{
				ctx:      ctx,
				indexId:  utf8ToEucjp("こんにちは"),
				searchId: utf8ToEucjp("こんにちは"),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 4.1: success search with 😀",
			args: args{
				ctx:      ctx,
				indexId:  "😀",
				searchId: "😀",
			},
			want: want{
				resultSize: 10,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			s, err := test.beforeFunc(test.args)
			if err != nil {
				tt.Errorf("error = %v", err)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			req := &payload.Search_IDRequest{
				Id:     test.args.searchId,
				Config: defaultSearch_Config,
			}
			gotRes, err := s.SearchByID(test.args.ctx, req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_toSearchResponse(t *testing.T) {
	t.Parallel()
	type args struct {
		dists []model.Distance
		err   error
	}
	type want struct {
		wantRes *payload.Search_Response
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		           dists: nil,
		           err: nil,
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
		           dists: nil,
		           err: nil,
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

			gotRes, err := toSearchResponse(test.args.dists, test.args.err)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamSearch(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Search_StreamSearchServer
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

			err := s.StreamSearch(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamSearchByID(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Search_StreamSearchByIDServer
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

			err := s.StreamSearchByID(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiSearch(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           reqs: nil,
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
		           ctx: nil,
		           reqs: nil,
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

			gotRes, err := s.MultiSearch(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiSearchByID(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiIDRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           reqs: nil,
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
		           ctx: nil,
		           reqs: nil,
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

			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

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

	// functions to generator vectors for testing
	genF32Vec := func(dist vector.Distribution, dim int) []float32 {
		generator, _ := vector.Float32VectorGenerator(dist)
		return generator(1, dim)[0]
	}
	genIntVec := func(dist vector.Distribution, dim int) []float32 {
		generator, _ := vector.Uint8VectorGenerator(dist)
		ivec := generator(1, dim)[0]

		vec := make([]float32, dim)
		for i := 0; i < dim; i++ {
			vec[i] = float32(ivec[i])
		}
		return vec
	}

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
			ivec := genIntVec(vector.Gaussian, invalidDim)
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec,
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
			ivec := genF32Vec(vector.Gaussian, invalidDim)
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: ivec,
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
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: genIntVec(vector.Gaussian, intVecDim),
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
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: genF32Vec(vector.Gaussian, f32VecDim),
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
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: genIntVec(vector.Uniform, intVecDim),
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
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: genF32Vec(vector.Uniform, f32VecDim),
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
					Vector: []float32{0, 0, 0},
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
					Vector: []float32{0, 0, 0},
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
					Vector: []float32{math.MinInt, math.MinInt, math.MinInt},
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
					Vector: []float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32},
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
					Vector: []float32{math.MaxInt, math.MaxInt, math.MaxInt},
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
					Vector: []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32},
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
					Vector: []float32{nan, nan, nan},
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
			bVec := genIntVec(vector.Gaussian, intVecDim) // used in beforeFunc

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
					s.ngt.Insert(id, bVec)
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
			bVec := genIntVec(vector.Gaussian, intVecDim) // use in beforeFunc

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
					s.ngt.Insert(id, bVec)
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
		wantRes *payload.Object_Locations
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type objectType int
	const (
		Uint8 objectType = iota
		Float
	)

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

	genF32Vec := func(dist vector.Distribution, num int, dim int) [][]float32 {
		generator, _ := vector.Float32VectorGenerator(dist)
		return generator(num, dim)
	}
	genIntVec := func(dist vector.Distribution, num int, dim int) [][]float32 {
		generator, _ := vector.Uint8VectorGenerator(dist)
		ivecs := generator(num, dim)
		result := make([][]float32, num)

		for j, ivec := range ivecs {
			vec := make([]float32, dim)
			for i := 0; i < dim; i++ {
				vec[i] = float32(ivec[i])
			}
			result[j] = vec
		}
		return result
	}

	genMultiInsertReq := func(t objectType, dist vector.Distribution, num int, dim int) *payload.Insert_MultiRequest {
		var vecs [][]float32
		switch t {
		case Float:
			vecs = genF32Vec(dist, num, dim)
		case Uint8:
			vecs = genIntVec(dist, num, dim)
		}

		req := &payload.Insert_MultiRequest{
			Requests: make([]*payload.Insert_Request, num),
		}
		for i, vec := range vecs {
			req.Requests[i] = &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     "uuid-" + strconv.Itoa(i+1),
					Vector: vec,
				},
			}
		}

		return req
	}

	// generate MultiInsert request with the same vector
	genSameVecMultiInsertReq := func(num int, vec []float32) *payload.Insert_MultiRequest {
		req := &payload.Insert_MultiRequest{
			Requests: make([]*payload.Insert_Request, num),
		}
		for i := 0; i < num; i++ {
			req.Requests[i] = &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     "uuid-" + strconv.Itoa(i+1),
					Vector: vec,
				},
			}
		}

		return req
	}

	genObjectLocations := func(num int, name string, ip string) *payload.Object_Locations {
		result := &payload.Object_Locations{
			Locations: make([]*payload.Object_Location, num),
		}

		for i := 0; i < num; i++ {
			result.Locations[i] = &payload.Object_Location{
				Name: name,
				Uuid: "uuid-" + strconv.Itoa(i+1),
				Ips:  []string{ip},
			}
		}
		return result
	}

	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", err, w.err)
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
				- case 4.2: Fail to MultiInsert with all existed vector when SkipStrictExistCheck is false
				- case 4.3: Fail to MultiInsert with 2 existed ID when SkipStrictExistCheck is true
				- case 4.4: Fail to MultiInsert with all existed vector when SkipStrictExistCheck is true
			- existed vector (with 100 insert request in a single MultiInsert request)
				- case 4.1: Success to MultiInsert with 2 existed vector when SkipStrictExistCheck is false
				- case 4.2: Success to MultiInsert with all existed vector when SkipStrictExistCheck is false
				- case 4.3: Success to MultiInsert with 2 existed vector when SkipStrictExistCheck is true
				- case 4.4: Success to MultiInsert with all existed vector when SkipStrictExistCheck is true
			- existed ID & existed vector (with 100 insert request in a single MultiInsert request)
				- case 4.1: Fail to MultiInsert with 2 existed ID & vector when SkipStrictExistCheck is false
				- case 4.2: Fail to MultiInsert with all existed ID & vector when SkipStrictExistCheck is false
				- case 4.3: Fail to MultiInsert with 2 existed ID & vector when SkipStrictExistCheck is true
				- case 4.4: Fail to MultiInsert with all existed ID & vector when SkipStrictExistCheck is true

	*/
	tests := []test{
		{
			name: "Equivalence Class Testing case 1.1: Success to MultiInsert 1 vector (vector type is uint8)",
			args: args{
				ctx:  ctx,
				reqs: genMultiInsertReq(Uint8, vector.Gaussian, 1, intVecDim),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultIntSvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(1, name, ip),
			},
		},
		{
			name: "Equivalence Class Testing case 1.2: Success to MultiInsert 1 vector (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genMultiInsertReq(Float, vector.Gaussian, 1, f32VecDim),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(1, name, ip),
			},
		},
		{
			name: "Equivalence Class Testing case 1.3: Success to MultiInsert 100 vector (vector type is uint8)",
			args: args{
				ctx:  ctx,
				reqs: genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultIntSvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Equivalence Class Testing case 1.4: Success to MultiInsert 100 vector (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 1, intVecDim+1)

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
			req := genMultiInsertReq(Float, vector.Gaussian, 1, f32VecDim+1)

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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim)
			req.Requests[99].Vector.Vector = genIntVec(vector.Gaussian, 1, intVecDim+1)[0]

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
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[99].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[99].Vector.Id,
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			req.Requests[99].Vector.Vector = genF32Vec(vector.Gaussian, 1, f32VecDim+1)[0]

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
						err := errors.ErrIncompatibleDimensionSize(len(req.Requests[99].Vector.Vector), intVecDim)
						err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   req.Requests[99].Vector.Id,
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim)
			for i := 0; i < 100; i += 2 {
				req.Requests[i].Vector.Vector = genIntVec(vector.Gaussian, 1, intVecDim+1)[0]
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			for i := 0; i < 100; i += 2 {
				req.Requests[i].Vector.Vector = genF32Vec(vector.Gaussian, 1, f32VecDim+1)[0]
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim+1)

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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim+1)

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
		{
			name: "Boundary Value Testing case 1.1: Success to MultiInsert with 0 value vector (vector type is uint8)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{0, 0, 0}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultIntSvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 1.2: Success to MultiInsert with 0 value vector (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{0, 0, 0}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 2.1: Success to MultiInsert with min value vector (vector type is uint8)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{math.MinInt, math.MinInt, math.MinInt}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultIntSvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: Success to MultiInsert with min value vector (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 3.1: Success to MultiInsert with max value vector (vector type is uint8)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{math.MaxUint8, math.MaxUint8, math.MaxUint8}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultIntSvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 3.2: Success to MultiInsert with max value vector (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		func() test {
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim)
			for i := 0; i < 100; i += 2 {
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, f32VecDim)
			for i := 0; i < 100; i += 2 {
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim)
			for i := 0; i < 100; i++ {
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, f32VecDim)
			for i := 0; i < 100; i++ {
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim)
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, f32VecDim)
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, intVecDim)
			for i := 0; i < len(req.Requests); i += 2 {
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
			req := genMultiInsertReq(Uint8, vector.Gaussian, 100, f32VecDim)
			for i := 0; i < len(req.Requests); i += 2 {
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
			req := genSameVecMultiInsertReq(100, make([]float32, maxVecDim))

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
			req := genSameVecMultiInsertReq(100, make([]float32, maxVecDim))

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
		{
			name: "Boundary Value Testing case 6.1: Success to MultiInsert with NaN value (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{float32(math.NaN()), float32(math.NaN()), float32(math.NaN())}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 6.2: Success to MultiInsert with +Inf value (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{float32(math.Inf(+1.0)), float32(math.Inf(+1.0)), float32(math.Inf(+1.0))}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 6.3: Success to MultiInsert with -Inf value (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{float32(math.Inf(-1.0)), float32(math.Inf(-1.0)), float32(math.Inf(-1.0))}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		{
			name: "Boundary Value Testing case 6.4: Success to MultiInsert with -0 value (vector type is float32)",
			args: args{
				ctx:  ctx,
				reqs: genSameVecMultiInsertReq(100, []float32{float32(math.Copysign(0, -1.0)), float32(math.Copysign(0, -1.0)), float32(math.Copysign(0, -1.0))}),
			},
			fields: fields{
				name:              name,
				ip:                ip,
				svcCfg:            defaultF32SvcCfg,
				svcOpts:           defaultSvcOpts,
				streamConcurrency: 0,
			},
			want: want{
				wantRes: genObjectLocations(100, name, ip),
			},
		},
		func() test {
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			vid := req.Requests[50].Vector.Id
			req.Requests[50] = nil

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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests); i += 2 {
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests); i += 2 {
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			vid := req.Requests[0].Vector.Id
			for i := 0; i < len(req.Requests); i += 2 {
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
			}

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
			}

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

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
					wantRes: genObjectLocations(100, name, ip),
				},
			}
		}(),
		func() test {
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
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
					wantRes: genObjectLocations(100, name, ip),
				},
			}
		}(),
		func() test {
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

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
					wantRes: genObjectLocations(100, name, ip),
				},
			}
		}(),
		func() test {
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
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
					wantRes: genObjectLocations(100, name, ip),
				},
			}
		}(),
		func() test {
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
			w.Locations[0].Uuid = req.Requests[0].Vector.Id
			//w.Locations[1].Uuid = dupID

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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: false,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			req.Requests[0].Vector.Vector = req.Requests[1].Vector.Vector
			req.Requests[0].Vector.Id = req.Requests[1].Vector.Id

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
			w.Locations[0].Uuid = req.Requests[0].Vector.Id
			// w.Locations[1].Uuid = req.Requests[1].Vector.Id

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
			req := genMultiInsertReq(Float, vector.Gaussian, 100, f32VecDim)
			for i := 1; i < len(req.Requests); i++ {
				req.Requests[i].Vector.Id = req.Requests[0].Vector.Id
				req.Requests[i].Vector.Vector = req.Requests[0].Vector.Vector
			}

			// set insert request config
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}
			for i := 0; i < len(req.Requests); i++ {
				req.Requests[i].Config = cfg
			}

			// set want
			w := genObjectLocations(100, name, ip)
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

			gotRes, err := s.MultiInsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Update(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Update_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
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
		           ctx: nil,
		           req: nil,
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
		           ctx: nil,
		           req: nil,
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

			gotRes, err := s.Update(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamUpdate(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Update_StreamUpdateServer
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

			err := s.StreamUpdate(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiUpdate(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Update_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           reqs: nil,
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
		           ctx: nil,
		           reqs: nil,
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

			gotRes, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Upsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Upsert_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           req: nil,
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
		           ctx: nil,
		           req: nil,
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

			gotLoc, err := s.Upsert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamUpsert(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Upsert_StreamUpsertServer
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

			err := s.StreamUpsert(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiUpsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Upsert_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           reqs: nil,
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
		           ctx: nil,
		           reqs: nil,
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

			gotRes, err := s.MultiUpsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Remove(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Remove_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
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
		           ctx: nil,
		           req: nil,
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
		           ctx: nil,
		           req: nil,
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

			gotRes, err := s.Remove(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamRemove(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Remove_StreamRemoveServer
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

			err := s.StreamRemove(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiRemove(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Remove_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           reqs: nil,
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
		           ctx: nil,
		           reqs: nil,
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

			gotRes, err := s.MultiRemove(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_GetObject(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  *payload.Object_VectorRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Vector
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Vector, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Vector, err error) error {
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
		           ctx: nil,
		           id: nil,
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
		           ctx: nil,
		           id: nil,
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

			gotRes, err := s.GetObject(test.args.ctx, test.args.id)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamGetObject(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Object_StreamGetObjectServer
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

			err := s.StreamGetObject(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_CreateIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		c   *payload.Control_CreateIndexRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
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
		           ctx: nil,
		           c: nil,
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
		           ctx: nil,
		           c: nil,
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

			gotRes, err := s.CreateIndex(test.args.ctx, test.args.c)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_SaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		in1 *payload.Empty
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
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
		           ctx: nil,
		           in1: nil,
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
		           ctx: nil,
		           in1: nil,
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

			gotRes, err := s.SaveIndex(test.args.ctx, test.args.in1)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_CreateAndSaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		c   *payload.Control_CreateIndexRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
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
		           ctx: nil,
		           c: nil,
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
		           ctx: nil,
		           c: nil,
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

			gotRes, err := s.CreateAndSaveIndex(test.args.ctx, test.args.c)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_IndexInfo(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		in1 *payload.Empty
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Info_Index_Count
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Index_Count, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Info_Index_Count, err error) error {
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
		           ctx: nil,
		           in1: nil,
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
		           ctx: nil,
		           in1: nil,
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

			gotRes, err := s.IndexInfo(test.args.ctx, test.args.in1)
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
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.LinearSearch(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
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
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.LinearSearchByID(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
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
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			err := s.StreamLinearSearch(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
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
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			err := s.StreamLinearSearchByID(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiLinearSearch(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiRequest
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.MultiLinearSearch(test.args.ctx, test.args.reqs)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiLinearSearchByID(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiIDRequest
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.MultiLinearSearchByID(test.args.ctx, test.args.reqs)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
