// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm/qbg"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/qbg/service"
)

func Test_server_Remove(t *testing.T) {
	t.Parallel()

	type args struct {
		indexID  string
		removeID string
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
		beforeFunc func(*testing.T, context.Context, args) (Server, error)
		afterFunc  func(args)
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
			if !reflect.DeepEqual(gotRes.Uuid, w.wantUUID) {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantUUID)
			}
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

	defaultNgtConfig := &config.QBG{
		Dimension:        128,
		DistanceType:     qbg.L2.String(),
		ObjectType:       qbg.Float.String(),
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
	defaultBeforeFunc := func(t *testing.T, ctx context.Context, a args) (Server, error) {
		t.Helper()
		eg, ctx := errgroup.New(ctx)
		qbg, err := newIndexedQBGService(ctx, eg, request.Float, vector.Gaussian, insertNum, defaultInsertConfig, defaultNgtConfig, nil, []string{a.indexID}, nil)
		if err != nil {
			return nil, err
		}
		s, err := New(WithErrGroup(eg), WithQBG(qbg))
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	/*
		Remove test cases ( focus on ID(string), only test float32 ):
		- Equivalence Class Testing ( 1000 vectors inserted before a search )
			- case 1.1: success remove vector
			- case 2.1: fail remove with non-existent ID
		- Boundary Value Testing ( 1000 vectors inserted before a search )
			- case 1.1: fail remove with ""
			- case 2.1: success remove with ^@
			- case 2.2: success remove with ^I
			- case 2.3: success remove with ^J
			- case 2.4: success remove with ^M
			- case 2.5: success remove with ^[
			- case 2.6: success remove with ^?
			- case 3.1: success remove with utf-8 ID from utf-8 index
			- case 3.2: fail remove with utf-8 ID from s-jis index
			- case 3.3: fail remove with utf-8 ID from euc-jp index
			- case 3.4: fail remove with s-jis ID from utf-8 index
			- case 3.5: success remove with s-jis ID from s-jis index
			- case 3.6: fail remove with s-jis ID from euc-jp index
			- case 3.4: fail remove with euc-jp ID from utf-8 index
			- case 3.5: fail remove with euc-jp ID from s-jis index
			- case 3.6: success remove with euc-jp ID from euc-jp index
			- case 4.1: success remove with 😀
		- Decision Table Testing
		    - NONE
	*/
	tests := []test{
		{
			name: "Equivalence Class Testing case 1.1: success exists vector",
			args: args{
				indexID:  "test",
				removeID: "test",
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail exists with non-existent ID",
			args: args{
				indexID:  "test",
				removeID: "non-existent",
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 1.1: fail exists with \"\"",
			args: args{
				indexID:  "test",
				removeID: "",
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success exists with ^@",
			args: args{
				indexID:  string([]byte{0}),
				removeID: string([]byte{0}),
			},
			want: want{
				wantUUID: string([]byte{0}),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success exists with ^I",
			args: args{
				indexID:  "\t",
				removeID: "\t",
			},
			want: want{
				wantUUID: "\t",
			},
		},
		{
			name: "Boundary Value Testing case 2.3: success exists with ^J",
			args: args{
				indexID:  "\n",
				removeID: "\n",
			},
			want: want{
				wantUUID: "\n",
			},
		},
		{
			name: "Boundary Value Testing case 2.4: success exists with ^M",
			args: args{
				indexID:  "\r",
				removeID: "\r",
			},
			want: want{
				wantUUID: "\r",
			},
		},
		{
			name: "Boundary Value Testing case 2.5: success exists with ^[",
			args: args{
				indexID:  string([]byte{27}),
				removeID: string([]byte{27}),
			},
			want: want{
				wantUUID: string([]byte{27}),
			},
		},
		{
			name: "Boundary Value Testing case 2.6: success exists with ^?",
			args: args{
				indexID:  string([]byte{127}),
				removeID: string([]byte{127}),
			},
			want: want{
				wantUUID: string([]byte{127}),
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success exists with utf-8 ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				removeID: utf8Str,
			},
			want: want{
				wantUUID: utf8Str,
			},
		},
		{
			name: "Boundary Value Testing case 3.2: fail exists with utf-8 ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				removeID: utf8Str,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.3: fail exists with utf-8 ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				removeID: utf8Str,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.4: fail exists with s-jis ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				removeID: sjisStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.5: success exists with s-jis ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				removeID: sjisStr,
			},
			want: want{
				wantUUID: sjisStr,
			},
		},
		{
			name: "Boundary Value Testing case 3.6: fail exists with s-jis ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				removeID: sjisStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.7: fail exists with euc-jp ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				removeID: eucjpStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.8: fail exists with euc-jp ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				removeID: eucjpStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.9: success exists with euc-jp ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				removeID: eucjpStr,
			},
			want: want{
				wantUUID: eucjpStr,
			},
		},
		{
			name: "Boundary Value Testing case 4.1: success exists with 😀",
			args: args{
				indexID:  "😀",
				removeID: "😀",
			},
			want: want{
				wantUUID: "😀",
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
			s, err := test.beforeFunc(tt, ctx, test.args)
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

			req := &payload.Remove_Request{
				Id: &payload.Object_ID{
					Id: test.args.removeID,
				},
			}
			gotRes, err := s.Remove(ctx, req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_RemoveByTimestamp(t *testing.T) {
	type args struct {
		req *payload.Remove_TimestampRequest
	}
	type want struct {
		code    codes.Code
		wantLen int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(context.Context, args) (Server, func(context.Context) error, error)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code().String(), w.code.String())
			}
		}
		if len(gotLocs.GetLocations()) != w.wantLen {
			return errors.Errorf("got Len: \"%#v\",\n\t\t\t\twant Len: \"%#v\"", len(gotLocs.GetLocations()), w.wantLen)
		}
		return nil
	}

	defaultNgtConfig := &config.QBG{
		Dimension:        128,
		DistanceType:     qbg.L2.String(),
		ObjectType:       qbg.Float.String(),
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

	defaultInsertNum := 100
	defaultTimestamp := int64(1000)

	createInsertReq := func(num int) (*payload.Insert_MultiRequest, error) {
		return request.GenMultiInsertReq(
			request.Float,
			vector.Gaussian,
			num,
			defaultNgtConfig.Dimension,
			&payload.Insert_Config{
				SkipStrictExistCheck: true,
				Timestamp:            defaultTimestamp,
			},
		)
	}

	defaultBeforeFunc := func(ctx context.Context, _ args) (Server, func(ctx context.Context) error, error) {
		eg, ctx := errgroup.New(ctx)
		qbg, err := service.New(defaultNgtConfig,
			service.WithErrGroup(eg),
			service.WithEnableInMemoryMode(true),
		)
		if err != nil {
			return nil, nil, err
		}

		s, err := New(
			WithErrGroup(eg),
			WithQBG(qbg),
		)
		if err != nil {
			return nil, nil, err
		}

		req, err := createInsertReq(defaultInsertNum)
		if err != nil {
			return nil, nil, err
		}
		for _, req := range req.GetRequests() {
			_, err := s.Insert(ctx, req)
			if err != nil {
				return nil, nil, err
			}
		}
		return s, func(ctx context.Context) error {
			return qbg.CreateIndex(ctx, 1000)
		}, err
	}
	tests := []test{
		{
			name: "succeeds if all vector data is deleted when the operator is Eq",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp,
							Operator:  payload.Remove_Timestamp_Eq,
						},
					},
				},
			},
			want: want{
				code:    codes.OK,
				wantLen: defaultInsertNum,
			},
		},
		{
			name: "succeeds if all vector data is deleted when the operator is Ge",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp,
							Operator:  payload.Remove_Timestamp_Ge,
						},
					},
				},
			},
			want: want{
				code:    codes.OK,
				wantLen: defaultInsertNum,
			},
		},
		{
			name: "succeeds if one vector data is deleted when the operator are Gt and Lt",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp,
							Operator:  payload.Remove_Timestamp_Gt,
						},
						{
							Timestamp: defaultTimestamp + 2,
							Operator:  payload.Remove_Timestamp_Lt,
						},
					},
				},
			},
			// Insert two additional vectors.
			beforeFunc: func(ctx context.Context, a args) (Server, func(context.Context) error, error) {
				s, fn, err := defaultBeforeFunc(ctx, a)
				if err != nil {
					return nil, nil, err
				}

				req, err := createInsertReq(2)
				if err != nil {
					return nil, nil, err
				}
				for i := range req.GetRequests() {
					req.Requests[i].Vector.Id += "-" + strconv.Itoa(i+1)
					req.Requests[i].Config.Timestamp += int64(i + 1)
					_, err := s.Insert(ctx, req.Requests[i])
					if err != nil {
						return nil, nil, err
					}
				}
				return s, fn, nil
			},
			want: want{
				code:    codes.OK,
				wantLen: 1,
			},
		},
		{
			name: "succeeds if all vector data is deleted when the operator is Le",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp,
							Operator:  payload.Remove_Timestamp_Le,
						},
					},
				},
			},
			want: want{
				code:    codes.OK,
				wantLen: defaultInsertNum,
			},
		},
		{
			name: "succeeds if all vector data is deleted when the operator is Gt and Lt",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp / 2,
							Operator:  payload.Remove_Timestamp_Gt,
						},
						{
							Timestamp: defaultTimestamp * 2,
							Operator:  payload.Remove_Timestamp_Lt,
						},
					},
				},
			},
			want: want{
				code:    codes.OK,
				wantLen: defaultInsertNum,
			},
		},
		{
			name: "fails if the target vector is not found when the operator is Eq",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp * 2,
							Operator:  payload.Remove_Timestamp_Eq,
						},
					},
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "fails if target vector is not found when the operator is Gt",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp * 2,
							Operator:  payload.Remove_Timestamp_Gt,
						},
					},
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "fails if all vector data is deleted when the operator is Lt",
			args: args{
				req: &payload.Remove_TimestampRequest{
					Timestamps: []*payload.Remove_Timestamp{
						{
							Timestamp: defaultTimestamp / 2,
							Operator:  payload.Remove_Timestamp_Lt,
						},
					},
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() { cancel() })

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			s, fn, err := test.beforeFunc(ctx, test.args)
			if err != nil {
				t.Errorf("error = %v", err)
				return
			}
			if err := fn(ctx); err != nil {
				t.Errorf("error = %v", err)
				return
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotLocs, err := s.RemoveByTimestamp(ctx, test.args.req)
			if err := checkFunc(test.want, gotLocs, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_timestampOpsFunc(t *testing.T) {
	type args struct {
		timestamp int64
		ts        []*payload.Remove_Timestamp
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if got != w.want {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return true when the timestamp is equal",
			args: args{
				timestamp: 1000,
				ts: []*payload.Remove_Timestamp{
					{
						Timestamp: 1000,
						Operator:  payload.Remove_Timestamp_Eq,
					},
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when timestamps is within range",
			args: args{
				timestamp: 1001,
				ts: []*payload.Remove_Timestamp{
					{
						Timestamp: 1000,
						Operator:  payload.Remove_Timestamp_Gt,
					},
					{
						Timestamp: 2000,
						Operator:  payload.Remove_Timestamp_Lt,
					},
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when timestamps is not within range",
			args: args{
				timestamp: 900,
				ts: []*payload.Remove_Timestamp{
					{
						Timestamp: 1000,
						Operator:  payload.Remove_Timestamp_Gt,
					},
					{
						Timestamp: 2000,
						Operator:  payload.Remove_Timestamp_Lt,
					},
				},
			},
			want: want{
				want: false,
			},
		},
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

			got := timestampOpsFunc(test.args.ts)(test.args.timestamp)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_timestampOpFunc(t *testing.T) {
	type args struct {
		timestamp int64
		ts        *payload.Remove_Timestamp
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if got != w.want {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return true when the timestamp is equal",
			args: args{
				timestamp: 1000,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Eq,
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the timestamp is not equal",
			args: args{
				timestamp: 1100,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Ne,
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the timestamp greater or equal",
			args: args{
				timestamp: 1000,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Ge,
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the timestamp is greater",
			args: args{
				timestamp: 1100,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Gt,
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the timestamp is less or equal",
			args: args{
				timestamp: 1000,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Le,
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the timestamp is less",
			args: args{
				timestamp: 900,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Lt,
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false when the operator is invalid",
			args: args{
				timestamp: 1000,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Operator(100),
				},
			},
			want: want{
				want: false,
			},
		},

		{
			name: "return false when the timestamp does not match the Eq operator",
			args: args{
				timestamp: 1100,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Eq,
				},
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when the timestamp does not match the Ne operator",
			args: args{
				timestamp: 1000,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Ne,
				},
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when the timestamp does not match the Ge operator",
			args: args{
				timestamp: 900,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Ge,
				},
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when the timestamp does not match the Gt operator",
			args: args{
				timestamp: 900,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Gt,
				},
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when the timestamp does not match the Le operator",
			args: args{
				timestamp: 1100,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Le,
				},
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when the timestamp does not match the Lt operator",
			args: args{
				timestamp: 1100,
				ts: &payload.Remove_Timestamp{
					Timestamp: 1000,
					Operator:  payload.Remove_Timestamp_Lt,
				},
			},
			want: want{
				want: false,
			},
		},
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

			got := timestampOpFunc(test.args.ts)(test.args.timestamp)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_StreamRemove(t *testing.T) {
// 	type args struct {
// 		stream vald.Remove_StreamRemoveServer
// 	}
// 	type fields struct {
// 		name                     string
// 		ip                       string
// 		qbg                      service.QBG
// 		eg                       errgroup.Group
// 		streamConcurrency        int
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
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
// 		           name:"",
// 		           ip:"",
// 		           qbg:nil,
// 		           eg:nil,
// 		           streamConcurrency:0,
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
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
// 		           name:"",
// 		           ip:"",
// 		           qbg:nil,
// 		           eg:nil,
// 		           streamConcurrency:0,
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
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
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				qbg:                      test.fields.qbg,
// 				eg:                       test.fields.eg,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 			}
//
// 			err := s.StreamRemove(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		name                     string
// 		ip                       string
// 		qbg                      service.QBG
// 		eg                       errgroup.Group
// 		streamConcurrency        int
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
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
// 		           name:"",
// 		           ip:"",
// 		           qbg:nil,
// 		           eg:nil,
// 		           streamConcurrency:0,
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
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
// 		           name:"",
// 		           ip:"",
// 		           qbg:nil,
// 		           eg:nil,
// 		           streamConcurrency:0,
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
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
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				qbg:                      test.fields.qbg,
// 				eg:                       test.fields.eg,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 			}
//
// 			gotRes, err := s.MultiRemove(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
