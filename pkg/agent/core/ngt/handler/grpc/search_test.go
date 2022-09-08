// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_Search(t *testing.T) {
	t.Parallel()

	type args struct {
		insertNum int
		req       *payload.Search_Request
	}
	type fields struct {
		objectType   request.ObjectType
		distribution vector.Distribution
		overwriteVec [][]float32

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
		beforeFunc func(context.Context, fields, args) (Server, error)
		afterFunc  func(args)
	}

	const (
		defaultDimensionSize = 32
	)

	defaultInsertConfig := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}
	defaultBeforeFunc := func(ctx context.Context, f fields, a args) (Server, error) {
		return buildIndex(ctx, f.objectType, f.distribution, a.insertNum, defaultInsertConfig, f.ngtCfg, f.ngtOpts, nil, f.overwriteVec)
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
		Timeout: int64(time.Second),
	}
	genSameVecs := func(ot request.ObjectType, n int, dim int) [][]float32 {
		var vecs [][]float32
		var err error
		switch ot {
		case request.Float:
			vecs, err = vector.GenF32Vec(vector.Gaussian, 1, dim)
		case request.Uint8:
			vecs, err = vector.GenUint8Vec(vector.Gaussian, 1, dim)
		}
		if err != nil {
			t.Error(err)
		}

		res := make([][]float32, n)
		for i := 0; i < n; i++ {
			res[i] = vecs[0]
		}

		return res
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
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 1.2: success search vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail search vector with different dimension (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize+1))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Equivalence Class Testing case 2.2: fail search vector with different dimension (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize+1)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
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
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(uint8(0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 1.2: success search with +0 value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, +0.0),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 1.3: success search with -0 value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.Copysign(0, -1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success search with max value vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.MaxUint8)),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success search with max value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, math.MaxFloat32),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success search with min value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, -math.MaxFloat32),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 4.1: fail search with NaN value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.NaN())),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 5.1: fail search with Inf value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.Inf(+1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 6.1: fail search with -Inf value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.Inf(-1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 7.1: fail search with 0 length vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: []float32{},
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, "uint8"),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 7.2: fail search with 0 length vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: []float32{},
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 8.1: fail search with max dimension vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, math.MaxInt32>>7))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 8.2: fail search with max dimension vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, math.MaxInt32>>7)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.1: fail search with nil vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: nil,
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.2: fail search with nil vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: nil,
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
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
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 1.2: success search with Search_Config.Num=10 from 5 different vectors (type: float32)",
			args: args{
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 2.1: success search with Search_Config.Num=10 from 10 different vectors (type: uint8)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 2.2: success search with Search_Config.Num=10 from 10 different vectors (type: float32)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 3.1: success search with Search_Config.Num=10 from 20 different vectors (type: uint8)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 3.2: success search with Search_Config.Num=10 from 20 different vectors (type: float32)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 4.1: success search with Search_Config.Num=10 from 5 same vectors (type: uint8)",
			args: args{
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Uint8, 5, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 4.2: success search with Search_Config.Num=10 from 5 same vectors (type: float32)",
			args: args{
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Float, 5, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 5.1: success search with Search_Config.Num=10 from 10 same vectors (type: uint8)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Uint8, 10, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 5.2: success search with Search_Config.Num=10 from 10 same vectors (type: float32)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Float, 10, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 6.1: success search with Search_Config.Num=10 from 20 same vectors (type: uint8)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Uint8, 20, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 6.2: success search with Search_Config.Num=10 from 20 same vectors (type: float32)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Float, 20, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			s, err := test.beforeFunc(ctx, test.fields, test.args)
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
		indexID  string
		searchID string
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

	utf8Str := "こんにちは"
	eucjpStr, err := conv.Utf8ToEucjp(utf8Str)
	if err != nil {
		t.Error(err)
	}
	sjisStr, err := conv.Utf8ToSjis(utf8Str)
	if err != nil {
		t.Error(err)
	}

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
	defaultInsertConfig := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}
	defaultBeforeFunc := func(a args) (Server, error) {
		return buildIndex(a.ctx, request.Float, vector.Gaussian, insertNum, defaultInsertConfig, defaultNgtConfig, nil, []string{a.indexID}, nil)
	}
	defaultSearch_Config := &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.1,
		Timeout: 1000000000,
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
				indexID:  "test",
				searchID: "test",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail search with non-existent ID",
			args: args{
				ctx:      ctx,
				indexID:  "test",
				searchID: "non-existent",
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
				indexID:  "test",
				searchID: "",
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success search with ^@",
			args: args{
				ctx:      ctx,
				indexID:  string([]byte{0}),
				searchID: string([]byte{0}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success search with ^I",
			args: args{
				ctx:      ctx,
				indexID:  "\t",
				searchID: "\t",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.3: success search with ^J",
			args: args{
				ctx:      ctx,
				indexID:  "\n",
				searchID: "\n",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.4: success search with ^M",
			args: args{
				ctx:      ctx,
				indexID:  "\r",
				searchID: "\r",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.5: success search with ^[",
			args: args{
				ctx:      ctx,
				indexID:  string([]byte{27}),
				searchID: string([]byte{27}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.6: success search with ^?",
			args: args{
				ctx:      ctx,
				indexID:  string([]byte{127}),
				searchID: string([]byte{127}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success search with utf-8 ID from utf-8 index",
			args: args{
				ctx:      ctx,
				indexID:  utf8Str,
				searchID: utf8Str,
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.2: fail search with utf-8 ID from s-jis index",
			args: args{
				ctx:      ctx,
				indexID:  sjisStr,
				searchID: utf8Str,
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
				indexID:  eucjpStr,
				searchID: utf8Str,
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
				indexID:  utf8Str,
				searchID: sjisStr,
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
				indexID:  sjisStr,
				searchID: sjisStr,
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.6: fail search with s-jis ID from euc-jp index",
			args: args{
				ctx:      ctx,
				indexID:  eucjpStr,
				searchID: sjisStr,
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
				indexID:  utf8Str,
				searchID: eucjpStr,
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
				indexID:  sjisStr,
				searchID: eucjpStr,
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
				indexID:  eucjpStr,
				searchID: eucjpStr,
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 4.1: success search with 😀",
			args: args{
				ctx:      ctx,
				indexID:  "😀",
				searchID: "😀",
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
				Id:     test.args.searchID,
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
