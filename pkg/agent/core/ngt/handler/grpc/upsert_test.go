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
package grpc

import (
	"context"
	"math"
	"reflect"
	"testing"

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
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_Upsert(t *testing.T) {
	t.Parallel()

	// optIdx is used for additional vector before upsert test.
	type optIdx struct {
		id  string
		vec []float32
	}

	type args struct {
		optIdx optIdx
		req    *payload.Upsert_Request
	}
	type want struct {
		code     codes.Code
		wantUUID string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(context.Context, optIdx) (Server, error)
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code(), w.code)
			}
		} else {
			if uuid := gotRes.GetUuid(); w.wantUUID != uuid {
				return errors.Errorf("got uuid: \"%#v\",\n\t\t\t\twant uuid: \"%#v\"", uuid, w.wantUUID)
			}
		}
		return nil
	}

	const (
		defaultInsertNum = 1000
		dimension        = 128
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

	defaultUpsertConfig := &payload.Upsert_Config{
		SkipStrictExistCheck: true,
	}
	defaultInsertConfig := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}
	defaultBeforeFunc := func(objectType string, insertNum int) func(context.Context, optIdx) (Server, error) {
		cfg := &config.NGT{
			Dimension:        dimension,
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

		return func(ctx context.Context, opt optIdx) (Server, error) {
			var overwriteID []string
			if opt.id != "" {
				overwriteID = []string{
					opt.id,
				}
			}
			var overwriteVec [][]float32
			if opt.vec != nil {
				overwriteVec = [][]float32{
					opt.vec,
				}
			}
			return buildIndex(ctx, request.Float, vector.Gaussian, insertNum, defaultInsertConfig, cfg, nil, overwriteID, overwriteVec)
		}
	}

	/*
		Upsert test cases (only test float32 unless otherwise specified):
		- Equivalence Class Testing ( 1000 vectors inserted before an upsert )
			- case 1.1: success upsert with new ID
			- case 1.2: success upsert with existent ID
			- case 2.1: fail upsert with one different dimension vector (type: uint8)
			- case 2.2: fail upsert with one different dimension vector (type: float32)
		- Boundary Value Testing ( 1000 vectors inserted before an upsert if not specified )
			- case 1.1: fail upsert with "" as ID
			- case 1.2: fail upsert to empty index with "" as ID
			- case 2.1: success upsert with ^@ as duplicated ID
			- case 2.2: success upsert with ^@ as new ID
			- case 2.3: success upsert to empty index with ^@ as new ID
			- case 2.4: success upsert with ^I as duplicated ID
			- case 2.5: success upsert with ^I as new ID
			- case 2.6: success upsert to empty index with ^I as new ID
			- case 2.7: success upsert with ^J as duplicated ID
			- case 2.8: success upsert with ^J as new ID
			- case 2.9: success upsert to empty index with ^J as new ID
			- case 2.10: success upsert with ^M as duplicated ID
			- case 2.11: success upsert with ^M as new ID
			- case 2.12: success upsert to empty index with ^M as new ID
			- case 2.13: success upsert with ^[ as duplicated ID
			- case 2.14: success upsert with ^[ as new ID
			- case 2.15: success upsert to empty index with ^[ as new ID
			- case 2.16: success upsert with ^? as duplicated ID
			- case 2.17: success upsert with ^? as new ID
			- case 2.18: success upsert to empty index with ^? as new ID
			- case 3.1: success upsert with utf-8 ID from utf-8 index
			- case 3.2: success upsert with utf-8 ID from s-jis index
			- case 3.3: success upsert with utf-8 ID from euc-jp index
			- case 3.4: success upsert with s-jis ID from utf-8 index
			- case 3.5: success upsert with s-jis ID from s-jis index
			- case 3.6: success upsert with s-jis ID from euc-jp index
			- case 3.7: success upsert with euc-jp ID from utf-8 index
			- case 3.8: success upsert with euc-jp ID from s-jis index
			- case 3.9: success upsert with euc-jp ID from euc-jp index
			- case 4.1: success upsert with 😀 as duplicated ID
			- case 4.2: success upsert with 😀 as new ID
			- case 4.3: success upsert to empty index with 😀 as new ID
			- case 5.1: success upsert with one 0 value vector with duplicated ID (type: uint8)
			- case 5.2: success upsert with one 0 value vector with new ID (type: uint8)
			- case 5.3: success upsert to empty index with one 0 value vector with new ID (type: uint8)
			- case 5.4: success upsert with one +0 value vector with duplicated ID (type: float32)
			- case 5.5: success upsert with one +0 value vector with new ID (type: float32)
			- case 5.6: success upsert to empty index with one +0 value vector with new ID (type: float32)
			- case 5.7: success upsert with one -0 value vector with duplicated ID (type: float32)
			- case 5.8: success upsert with one -0 value vector with new ID (type: float32)
			- case 5.9: success upsert to empty index with one -0 value vector with new ID (type: float32)
			- case 6.1: success upsert with one min value vector with duplicated ID (type: uint8)
			- case 6.2: success upsert with one min value vector with new ID (type: uint8)
			- case 6.3: success upsert to empty index with one min value vector with new ID (type: uint8)
			- case 6.4: success upsert with one min value vector with duplicated ID (type: float32)
			- case 6.5: success upsert with one min value vector with new ID (type: float32)
			- case 6.6: success upsert to empty index with one min value vector with new ID (type: float32)
			- case 7.1: success upsert with one max value vector with duplicated ID (type: uint8)
			- case 7.2: success upsert with one max value vector with new ID (type: uint8)
			- case 7.3: success upsert to empty index with one max value vector with new ID (type: uint8)
			- case 7.4: success upsert with one max value vector with duplicated ID (type: float32)
			- case 7.5: success upsert with one max value vector (type: float32)
			- case 7.6: success upsert to empty index with one max value vector (type: float32)
			- case 8.1: success upsert with one NaN value vector with duplicated ID (type: float32) // NOTE: To fix it, it is necessary to check all of vector value
			- case 8.2: success upsert with one NaN value vector with new ID (type: float32) // NOTE: To fix it, it is necessary to check all of vector value
			- case 8.3: success upsert to empty index with one NaN value vector with new ID (type: float32) // NOTE: To fix it, it is necessary to check all of vector value
			- case 9.1: success upsert with one +inf value vector with duplicated ID (type: float32)
			- case 9.2: success upsert with one +inf value vector with new ID (type: float32)
			- case 9.3: success upsert to empty index with one +inf value vector with new ID (type: float32)
			- case 9.4: success upsert with one -inf value vector with duplicated ID (type: float32)
			- case 9.5: success upsert with one -inf value vector with new ID (type: float32)
			- case 9.6: success upsert to empty index with one -inf value vector with new ID (type: float32)
			- case 10.1: fail upsert with one nil vector
			- case 10.2: fail upsert to empty index with one nil vector
			- case 11.1: fail upsert with one empty vector
			- case 11.2: fail upsert to empty index with one empty vector
		- Decision Table Testing
			- case 1.1: fail upsert with one duplicated vector, duplicated ID and SkipStrictExistCheck is true
			- case 1.2: success upsert with one different vector, duplicated ID and SkipStrictExistsCheck is true
			- case 1.3: success upsert with one duplicated vector, different ID and SkipStrictExistCheck is true
			- case 1.4: success upsert with one different vector, different ID and SkipStrictExistsCheck is true
			- case 2.1: fail upsert with one duplicated vector, duplicated ID and SkipStrictExistCheck is false
			- case 2.2: success upsert with one different vector, duplicated ID and SkipStrictExistsCheck is false
			- case 2.3: success upsert with one duplicated vector, different ID and SkipStrictExistCheck is false
			- case 2.4: success upsert with one different vector, different ID and SkipStrictExistsCheck is false
	*/
	tests := []test{
		{
			name: "Equivalent Class Testing case 1.1: success upsert with new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Equivalent Class Testing case 1.2: success upsert with existent ID",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Equivalent Class Testing case 2.1: fail upsert with one different dimension vector (type: uint8)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, dimension+1))[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), defaultInsertNum),
		},
		{
			name: "Equivalent Class Testing case 2.2: fail upsert with one different dimension vector (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension+1)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 1.1: fail upsert with \"\" as ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 1.2: fail upsert to empty index with \"\" as ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success upsert with ^@ as duplicated ID",
			args: args{
				optIdx: optIdx{
					id: string([]byte{0}),
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{0}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{0}),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success upsert with ^@ as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{0}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{0}),
			},
		},
		{
			name: "Boundary Value Testing case 2.3: success upsert to empty index with ^@ as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{0}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: string([]byte{0}),
			},
		},
		{
			name: "Boundary Value Testing case 2.4: success upsert with ^I as duplicated ID",
			args: args{
				optIdx: optIdx{
					id: "\t",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\t",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "\t",
			},
		},
		{
			name: "Boundary Value Testing case 2.5: success upsert with ^I as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\t",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "\t",
			},
		},
		{
			name: "Boundary Value Testing case 2.6: success upsert to empty index with ^I as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\t",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "\t",
			},
		},
		{
			name: "Boundary Value Testing case 2.7: success upsert with ^J as duplicated ID",
			args: args{
				optIdx: optIdx{
					id: "\n",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\n",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "\n",
			},
		},
		{
			name: "Boundary Value Testing case 2.8: success upsert with ^J as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\n",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "\n",
			},
		},
		{
			name: "Boundary Value Testing case 2.9: success upsert to empty index with ^J as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\n",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "\n",
			},
		},
		{
			name: "Boundary Value Testing case 2.10: success upsert with ^M as duplicated ID",
			args: args{
				optIdx: optIdx{
					id: "\r",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\r",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "\r",
			},
		},
		{
			name: "Boundary Value Testing case 2.11: success upsert with ^M as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\r",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "\r",
			},
		},
		{
			name: "Boundary Value Testing case 2.12: success upsert to empty index with ^M as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "\r",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "\r",
			},
		},
		{
			name: "Boundary Value Testing case 2.13: success upsert with ^[ as duplicated ID",
			args: args{
				optIdx: optIdx{
					id: string([]byte{27}),
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{27}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{27}),
			},
		},
		{
			name: "Boundary Value Testing case 2.14: success upsert with ^[ as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{27}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{27}),
			},
		},
		{
			name: "Boundary Value Testing case 2.15: success upsert to empty index with ^[ as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{27}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: string([]byte{27}),
			},
		},
		{
			name: "Boundary Value Testing case 2.16: success upsert with ^? as duplicated ID",
			args: args{
				optIdx: optIdx{
					id: string([]byte{127}),
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{127}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{127}),
			},
		},
		{
			name: "Boundary Value Testing case 2.17: success upsert with ^? as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{127}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{127}),
			},
		},
		{
			name: "Boundary Value Testing case 2.18: success upsert to empty index with ^? as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{127}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: string([]byte{127}),
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success upsert with utf-8 ID from utf-8 index",
			args: args{
				optIdx: optIdx{
					id: utf8Str,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     utf8Str,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: utf8Str,
			},
		},
		{
			name: "Boundary Value Testing case 3.2: success upsert with utf-8 ID from s-jis index",
			args: args{
				optIdx: optIdx{
					id: sjisStr,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     utf8Str,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: utf8Str,
			},
		},
		{
			name: "Boundary Value Testing case 3.3: success upsert with utf-8 ID from euc-jp index",
			args: args{
				optIdx: optIdx{
					id: eucjpStr,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     utf8Str,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: utf8Str,
			},
		},
		{
			name: "Boundary Value Testing case 3.4: success upsert with s-jis ID from utf-8 index",
			args: args{
				optIdx: optIdx{
					id: utf8Str,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     sjisStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: sjisStr,
			},
		},
		{
			name: "Boundary Value Testing case 3.5: success upsert with s-jis ID from s-jis index",
			args: args{
				optIdx: optIdx{
					id: sjisStr,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     sjisStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: sjisStr,
			},
		},
		{
			name: "Boundary Value Testing case 3.6: success upsert with s-jis ID from euc-jp index",
			args: args{
				optIdx: optIdx{
					id: eucjpStr,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     sjisStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: sjisStr,
			},
		},
		{
			name: "Boundary Value Testing case 3.7: success upsert with euc-jp ID from utf-8 index",
			args: args{
				optIdx: optIdx{
					id: utf8Str,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     eucjpStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: eucjpStr,
			},
		},
		{
			name: "Boundary Value Testing case 3.8: success upsert with euc-jp ID from s-jis index",
			args: args{
				optIdx: optIdx{
					id: sjisStr,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     eucjpStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: eucjpStr,
			},
		},
		{
			name: "Boundary Value Testing case 3.9: success upsert with euc-jp ID from euc-jp index",
			args: args{
				optIdx: optIdx{
					id: eucjpStr,
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     eucjpStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: eucjpStr,
			},
		},
		{
			name: "Boundary Value Testing case 4.1: success upsert with 😀 as duplicated ID",
			args: args{
				optIdx: optIdx{
					id: "😀",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "😀",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "😀",
			},
		},
		{
			name: "Boundary Value Testing case 4.2: success upsert with 😀 as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "😀",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "😀",
			},
		},
		{
			name: "Boundary Value Testing case 4.3: success upsert to empty index with 😀 as new ID",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "😀",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "😀",
			},
		},
		{
			name: "Boundary Value Testing case 5.1: success upsert with one 0 value vector with duplicated ID (type: uint8)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 1000),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.2: success upsert with one 0 value vector with new ID (type: uint8)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 1000),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.3: success upsert to empty index with one 0 value vector with new ID (type: uint8)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.4: success upsert with one +0 value vector with duplicated ID (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, 0),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.5: success upsert with one +0 value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, 0),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.6: success upsert to empty index with one +0 value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, 0),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.7: success upsert with one -0 value vector with duplicated ID (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Copysign(0, -1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.8: success upsert with one -0 value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Copysign(0, -1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.9: success upsert to empty index with one -0 value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Copysign(0, -1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.1: success upsert with one min value vector with duplicated ID (type: uint8)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 1000),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.2: success upsert with one min value vector with new ID (type: uint8)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 1000),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.3: success upsert to empty index with one min value vector with new ID (type: uint8)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.4: success upsert with one min value vector with duplicated ID (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, -math.MaxFloat32),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.5: success upsert with one min value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, -math.MaxFloat32),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.6: success upsert with one min value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, -math.MaxFloat32),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.1: success upsert with one max value vector with  ID (type: uint8)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.MaxUint8)),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 1000),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.2: success upsert with one max value vector with new ID (type: uint8)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.MaxUint8)),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 1000),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.3: success upsert to empty index with one max value vector with new ID (type: uint8)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.MaxUint8)),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Uint8.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.4: success upsert with one max value vector with duplicated ID (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, math.MaxFloat32),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.5: success upsert with one max value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, math.MaxFloat32),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.6: success upsert to empty index with one max value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, math.MaxFloat32),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 8.1: success upsert with one NaN value vector with duplicated ID (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.NaN())),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 8.2: success upsert with one NaN value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.NaN())),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 8.3: success upsert to empty index with one NaN value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.NaN())),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 9.1: success upsert with one +inf value vector with duplicated ID (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 9.2: success upsert with one +inf value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 9.3: success upsert to empty index with one +inf value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 9.4: success upsert with one -inf value vector with duplicated ID (type: float32)",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(-1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 9.5: success upsert with one -inf value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(-1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 9.6: success upsert to empty index with one -inf value vector with new ID (type: float32)",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(-1.0))),
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 10.1: fail upsert with one nil vector",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: nil,
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 10.2: fail upsert to empty with one nil vector",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: nil,
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 11.1: fail upsert with one empty vector",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: []float32{},
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 11.2: fail upsert to empty index with one empty vector",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: []float32{},
					},
					Config: defaultUpsertConfig,
				},
			},
			beforeFunc: defaultBeforeFunc(ngt.Float.String(), 0),
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Decision Table Testing case 1.1: fail upsert with one duplicated vector, duplicated ID and SkipStrictExistCheck is true",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					optIdx: optIdx{
						id:  "test",
						vec: vector,
					},
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     "test",
							Vector: vector,
						},
						Config: defaultUpsertConfig,
					},
				}
			}(),
			want: want{
				code: codes.AlreadyExists,
			},
		},
		{
			name: "Decision Table Testing case 1.2: success upsert with one different vector, duplicated ID and SkipStrictExistCheck is true",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Decision Table Testing case 1.3: success upsert with one duplicated vector, different ID and SkipStrictExistCheck is true",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					optIdx: optIdx{
						id:  "test",
						vec: vector,
					},
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     "uuid-2", // the first uuid is overwritten, so use the second one
							Vector: vector,
						},
						Config: defaultUpsertConfig,
					},
				}
			}(),
			want: want{
				wantUUID: "uuid-2",
			},
		},
		{
			name: "Decision Table Testing case 1.4: success upsert with one different vector, different ID and SkipStrictExistCheck is true",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpsertConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Decision Table Testing case 2.1: fail upsert with one duplicated vector, duplicated ID and SkipStrictExistCheck is false",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					optIdx: optIdx{
						id:  "test",
						vec: vector,
					},
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     "test",
							Vector: vector,
						},
						Config: &payload.Upsert_Config{
							SkipStrictExistCheck: false,
						},
					},
				}
			}(),
			want: want{
				code: codes.AlreadyExists,
			},
		},
		{
			name: "Decision Table Testing case 2.2: success upsert with one duplicated vector, duplicated ID and SkipStrictExistCheck is false",
			args: args{
				optIdx: optIdx{
					id: "test",
				},
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: &payload.Upsert_Config{
						SkipStrictExistCheck: false,
					},
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Decision Table Testing case 2.3: success upsert with one duplicated vector, different ID and SkipStrictExistCheck is false",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					optIdx: optIdx{
						id:  "test",
						vec: vector,
					},
					req: &payload.Upsert_Request{
						Vector: &payload.Object_Vector{
							Id:     "uuid-2", // the first uuid is overwritten, so use the second one
							Vector: vector,
						},
						Config: &payload.Upsert_Config{
							SkipStrictExistCheck: false,
						},
					},
				}
			}(),
			want: want{
				wantUUID: "uuid-2",
			},
		},
		{
			name: "Decision Table Testing case 2.4: success upsert with one different vector, different ID and SkipStrictExistCheck is false",
			args: args{
				req: &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: &payload.Upsert_Config{
						SkipStrictExistCheck: false,
					},
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			defer cancel()
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc(ngt.Float.String(), defaultInsertNum)
			}
			s, err := test.beforeFunc(ctx, test.args.optIdx)
			if err != nil {
				tt.Errorf("error = %v", err)
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			gotLoc, err := s.Upsert(ctx, test.args.req)
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
