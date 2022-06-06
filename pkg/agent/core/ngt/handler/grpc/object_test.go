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
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_Exists(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		ctx      context.Context
		indexId  string
		searchId string
	}
	type want struct {
		code codes.Code
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Object_ID, error) error
		beforeFunc func(args) (Server, error)
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
	defaultBeforeFunc := func(a args) (Server, error) {
		return buildIndex(a.ctx, request.Float, vector.Gaussian, insertNum, defaultInsertConfig, defaultNgtConfig, nil, []string{a.indexId}, nil)
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
				ctx:      ctx,
				indexId:  "test",
				searchId: "test",
			},
			want: want{},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail exists with non-existent ID",
			args: args{
				ctx:      ctx,
				indexId:  "test",
				searchId: "non-existent",
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 1.1: fail exists with \"\"",
			args: args{
				ctx:      ctx,
				indexId:  "test",
				searchId: "",
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success exists with ^@",
			args: args{
				ctx:      ctx,
				indexId:  string([]byte{0}),
				searchId: string([]byte{0}),
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.2: success exists with ^I",
			args: args{
				ctx:      ctx,
				indexId:  "\t",
				searchId: "\t",
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.3: success exists with ^J",
			args: args{
				ctx:      ctx,
				indexId:  "\n",
				searchId: "\n",
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.4: success exists with ^M",
			args: args{
				ctx:      ctx,
				indexId:  "\r",
				searchId: "\r",
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.5: success exists with ^[",
			args: args{
				ctx:      ctx,
				indexId:  string([]byte{27}),
				searchId: string([]byte{27}),
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 2.6: success exists with ^?",
			args: args{
				ctx:      ctx,
				indexId:  string([]byte{127}),
				searchId: string([]byte{127}),
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 3.1: success exists with utf-8 ID from utf-8 index",
			args: args{
				ctx:      ctx,
				indexId:  utf8Str,
				searchId: utf8Str,
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 3.2: fail exists with utf-8 ID from s-jis index",
			args: args{
				ctx:      ctx,
				indexId:  sjisStr,
				searchId: utf8Str,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.3: fail exists with utf-8 ID from euc-jp index",
			args: args{
				ctx:      ctx,
				indexId:  eucjpStr,
				searchId: utf8Str,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.4: fail exists with s-jis ID from utf-8 index",
			args: args{
				ctx:      ctx,
				indexId:  utf8Str,
				searchId: sjisStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.5: success exists with s-jis ID from s-jis index",
			args: args{
				ctx:      ctx,
				indexId:  sjisStr,
				searchId: sjisStr,
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 3.6: fail exists with s-jis ID from euc-jp index",
			args: args{
				ctx:      ctx,
				indexId:  eucjpStr,
				searchId: sjisStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.7: fail exists with euc-jp ID from utf-8 index",
			args: args{
				ctx:      ctx,
				indexId:  utf8Str,
				searchId: eucjpStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.8: fail exists with euc-jp ID from s-jis index",
			args: args{
				ctx:      ctx,
				indexId:  sjisStr,
				searchId: eucjpStr,
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.9: success exists with euc-jp ID from euc-jp index",
			args: args{
				ctx:      ctx,
				indexId:  eucjpStr,
				searchId: eucjpStr,
			},
			want: want{},
		},
		{
			name: "Boundary Value Testing case 4.1: success exists with 😀",
			args: args{
				ctx:      ctx,
				indexId:  "😀",
				searchId: "😀",
			},
			want: want{},
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

			req := &payload.Object_ID{
				Id: test.args.searchId,
			}
			gotRes, err := s.Exists(test.args.ctx, req)
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
