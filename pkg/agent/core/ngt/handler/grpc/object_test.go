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
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_Exists(t *testing.T) {
	t.Parallel()

	type args struct {
		indexID  string
		searchID string
	}
	type want struct {
		code codes.Code
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Object_ID, error) error
		beforeFunc func(context.Context, args) (Server, error)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_ID, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code(), w.code)
			}
		}
		return nil
	}

	const (
		insertNum = 1000
		dim       = 128
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
		Dimension:        dim,
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
	defaultBeforeFunc := func(ctx context.Context, a args) (Server, error) {
		return buildIndex(ctx, request.Float, vector.Gaussian, insertNum, defaultInsertConfig, defaultNgtConfig, nil, []string{a.indexID}, nil)
	}

	/*
		Exists test cases (focus on ID(string), only test float32):
		- Equivalence Class Testing ( 1000 vectors inserted before a search )
			- case 1.1: success exists vector
			- case 2.1: fail exists with non-existent ID
		- Boundary Value Testing ( 1000 vectors inserted before a search )
			- case 1.1: fail exists with ""
			- case 2.1: success exists with ^@
			- case 2.2: success exists with ^I
			- case 2.3: success exists with ^J
			- case 2.4: success exists with ^M
			- case 2.5: success exists with ^[
			- case 2.6: success exists with ^?
			- case 3.1: success exists with utf-8 ID from utf-8 index
			- case 3.2: fail exists with utf-8 ID from s-jis index
			- case 3.3: fail exists with utf-8 ID from euc-jp index
			- case 3.4: fail exists with s-jis ID from utf-8 index
			- case 3.5: success exists with s-jis ID from s-jis index
			- case 3.6: fail exists with s-jis ID from euc-jp index
			- case 3.4: fail exists with euc-jp ID from utf-8 index
			- case 3.5: fail exists with euc-jp ID from s-jis index
			- case 3.6: success exists with euc-jp ID from euc-jp index
			- case 4.1: success exists with 😀
		- Decision Table Testing
		    - NONE
	*/
	tests := []test{
		{
			name: "Equivalence Class Testing case 1.1: success exists vector",
			args: args{
				indexID:  "test",
				searchID: "test",
			},
			want: want{},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail exists with non-existent ID",
			args: args{
				indexID:  "test",
				searchID: "non-existent",
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 1.1: fail exists with \"\"",
			args: args{
				indexID:  "test",
				searchID: "",
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success exists with ^@",
			args: args{
				indexID:  string([]byte{0}),
				searchID: string([]byte{0}),
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.2: success exists with ^I",
			args: args{
				indexID:  "\t",
				searchID: "\t",
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.3: success exists with ^J",
			args: args{
				indexID:  "\n",
				searchID: "\n",
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.4: success exists with ^M",
			args: args{
				indexID:  "\r",
				searchID: "\r",
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.5: success exists with ^[",
			args: args{
				indexID:  string([]byte{27}),
				searchID: string([]byte{27}),
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.6: success exists with ^?",
			args: args{
				indexID:  string([]byte{127}),
				searchID: string([]byte{127}),
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 3.1: success exists with utf-8 ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				searchID: utf8Str,
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 3.2: fail exists with utf-8 ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				searchID: utf8Str,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.3: fail exists with utf-8 ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				searchID: utf8Str,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.4: fail exists with s-jis ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				searchID: sjisStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.5: success exists with s-jis ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				searchID: sjisStr,
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 3.6: fail exists with s-jis ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				searchID: sjisStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.7: fail exists with euc-jp ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				searchID: eucjpStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.8: fail exists with euc-jp ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				searchID: eucjpStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.9: success exists with euc-jp ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				searchID: eucjpStr,
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 4.1: success exists with 😀",
			args: args{
				indexID:  "😀",
				searchID: "😀",
			},
			want: want{},
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
			s, err := test.beforeFunc(ctx, test.args)
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

			req := &payload.Object_ID{
				Id: test.args.searchID,
			}
			gotRes, err := s.Exists(ctx, req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_GetObject(t *testing.T) {
	t.Parallel()
	type args struct {
		id *payload.Object_VectorRequest
	}
	type fields struct {
		srvOpts []Option
		svcCfg  *config.NGT
		svcOpts []service.Option
	}
	type want struct {
		wantRes *payload.Object_Vector
		errCode codes.Code
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Vector, error) error
		beforeFunc func(*testing.T, context.Context, args, Server)
		afterFunc  func(args)
	}

	// common variables for test
	const (
		name      = "vald-agent-ngt-1" // agent name
		dim       = 3                  //  vector dimension
		id        = "uuid-1"           // id for getObject request
		insertCnt = 1000               // default insert count
	)
	var (
		ip = net.LoadLocalIP() // agent ip address

		// default NGT configuration for test
		kvsdbCfg  = &config.KVSDB{}
		vqueueCfg = &config.VQueue{}

		defaultSvcCfg = &config.NGT{
			Dimension:    dim,
			DistanceType: ngt.Angle.String(),
			ObjectType:   ngt.Float.String(),
			KVSDB:        kvsdbCfg,
			VQueue:       vqueueCfg,
		}
		defaultSvcOpts = []service.Option{
			service.WithEnableInMemoryMode(true),
		}

		defaultInsertConfig = &payload.Insert_Config{}
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

	insertAndCreateIndex := func(t *testing.T, ctx context.Context, s Server, req *payload.Insert_MultiRequest) {
		if _, err := s.MultiInsert(ctx, req); err != nil {
			t.Fatal(err)
		}
		if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
			PoolSize: uint32(len(req.Requests)),
		}); err != nil {
			t.Fatal(err)
		}
	}

	defaultCheckFunc := func(w want, gotRes *payload.Object_Vector, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.errCode {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code(), w.errCode)
			}
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	/*
		- Equivalence Class Testing (1000 vectors inserted)
			- case 1.1: success to get object (type: uint8)
			- case 2.1: success to get object (type: float32)
		- Boundary Value Testing (1000 float32 vectors inserted)
			- case 1.1: fail to get object with ""
			- case 2.1: success to get object with ^@
			- case 2.2: success to get object with ^I
			- case 2.3: success to get object with ^J
			- case 2.4: success to get object with ^M
			- case 2.5: success to get object with ^[
			- case 2.6: success to get object with ^?
			- case 2.7: success to get object with utf-8 ID from utf-8 index
			- case 3.1: fail to get object with utf-8 ID from s-jis index
			- case 3.2: fail to get object with utf-8 ID from euc-jp index
			- case 3.3: fail to get object with s-jis ID from utf-8 index
			- case 3.4: success to get object with s-jis ID from s-jis index
			- case 4.1: fail to get object with s-jis ID from euc-jp index
			- case 4.2: fail to get object with euc-jp ID from utf-8 index
			- case 4.3: fail to get object with euc-jp ID from s-jis index
			- case 4.4: success to get object with euc-jp ID from euc-jp index
			- case 5.1: success to get object with 😀
		- Decision Table Testing
		    - NONE
	*/
	tests := []test{
		func() test {
			ir, err := request.GenMultiInsertReq(request.Uint8, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector

			return test{
				name: "Equivalence Class Testing case 1.1: success to get object (type: uint8)",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg: &config.NGT{
						Dimension:    dim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector

			return test{
				name: "Equivalence Class Testing case 2.1: success to get object (type: float32)",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}

			return test{
				name: `Boundary Value Testing case 1.1: fail to get object with ""`,
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: "",
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					errCode: codes.InvalidArgument,
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = string([]byte{0})

			return test{
				name: "Boundary Value Testing case 2.1: success to get object with ^@",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = "\t"

			return test{
				name: "Boundary Value Testing case 2.2: success to get object with ^I",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = "\n"

			return test{
				name: "Boundary Value Testing case 2.3: success to get object with ^J",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = "\r"

			return test{
				name: "Boundary Value Testing case 2.4: success to get object with ^M",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = string([]byte{27})

			return test{
				name: "Boundary Value Testing case 2.5: success to get object with ^[",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = string([]byte{127})

			return test{
				name: "Boundary Value Testing case 2.6: success to get object with ^?",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = utf8Str

			return test{
				name: "Boundary Value Testing case 2.7: success to get object with utf-8 ID from utf-8 index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = sjisStr

			return test{
				name: "Boundary Value Testing case 3.1: fail to get object with utf-8 ID from s-jis index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: utf8Str,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					errCode: codes.NotFound,
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = eucjpStr

			return test{
				name: "Boundary Value Testing case 3.2: fail to get object with utf-8 ID from euc-jp index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: utf8Str,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					errCode: codes.NotFound,
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = utf8Str

			return test{
				name: "Boundary Value Testing case 3.3: fail to get object with s-jis ID from utf-8 index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: sjisStr,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					errCode: codes.NotFound,
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = sjisStr

			return test{
				name: "Boundary Value Testing case 3.4: success to get object with s-jis ID from s-jis index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = eucjpStr

			return test{
				name: "Boundary Value Testing case 4.1: fail to get object with s-jis ID from euc-jp index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: sjisStr,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					errCode: codes.NotFound,
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = utf8Str

			return test{
				name: "Boundary Value Testing case 4.2: fail to get object with euc-jp ID from utf-8 index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: eucjpStr,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					errCode: codes.NotFound,
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = sjisStr

			return test{
				name: "Boundary Value Testing case 4.3: fail to get object with euc-jp ID from s-jis index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: eucjpStr,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					errCode: codes.NotFound,
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = eucjpStr

			return test{
				name: "Boundary Value Testing case 4.4: success to get object with euc-jp ID from euc-jp index",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
		func() test {
			ir, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
			if err != nil {
				t.Fatal(err)
			}
			reqVec := ir.Requests[0].Vector
			reqVec.Id = "😀"

			return test{
				name: "Boundary Value Testing case 5.1: success to get object with 😀",
				args: args{
					id: &payload.Object_VectorRequest{
						Id: &payload.Object_ID{
							Id: reqVec.Id,
						},
					},
				},
				fields: fields{
					srvOpts: []Option{
						WithName(name),
						WithIP(ip),
					},
					svcCfg:  defaultSvcCfg,
					svcOpts: defaultSvcOpts,
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					insertAndCreateIndex(t, ctx, s, ir)
				},
				want: want{
					wantRes: &payload.Object_Vector{
						Id:     reqVec.Id,
						Vector: reqVec.Vector,
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

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

			s, err := New(append(test.fields.srvOpts, WithNGT(ngt), WithErrGroup(eg))...)
			if err != nil {
				t.Errorf("failed to init service, err: %v", err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, ctx, test.args, s)
			}

			gotRes, err := s.GetObject(ctx, test.args.id)
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
