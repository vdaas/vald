//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package errdetails

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/rpc/errdetails"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/types"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	pproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func Test_decodeDetails(t *testing.T) {
	t.Parallel()
	type args struct {
		objs []any
	}
	tests := []struct {
		name        string
		args        args
		wantDetails []Detail
	}{
		{
			name: "returns nil when objs is nil",
			args: args{
				objs: nil,
			},
			wantDetails: nil,
		},
		{
			name: "returns empty details when objs contains only nil",
			args: args{
				objs: []any{nil, nil},
			},
			wantDetails: []Detail{},
		},
		{
			name: "returns details for *spb.Status",
			args: args{
				objs: []any{
					&spb.Status{
						Code:    1,
						Message: "test",
					},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "google.rpc.Status",
					Message: &spb.Status{
						Code:    1,
						Message: "test",
					},
				},
			},
		},
		{
			name: "returns details for spb.Status",
			args: args{
				objs: []any{
					spb.Status{
						Code:    1,
						Message: "test",
					},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "google.rpc.Status",
					Message: &spb.Status{
						Code:    1,
						Message: "test",
					},
				},
			},
		},
		{
			name: "returns details for *status.Status",
			args: args{
				objs: []any{
					status.New(codes.InvalidArgument, "invalid"),
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "google.rpc.Status",
					Message: &spb.Status{
						Code:    int32(codes.InvalidArgument),
						Message: "invalid",
					},
				},
			},
		},
		{
			name: "returns details for status.Status",
			args: args{
				objs: []any{
					*status.New(codes.InvalidArgument, "invalid"),
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "google.rpc.Status",
					Message: &spb.Status{
						Code:    int32(codes.InvalidArgument),
						Message: "invalid",
					},
				},
			},
		},
		{
			name: "returns details for *status.Status with details",
			args: args{
				objs: func() []any {
					st := status.New(codes.InvalidArgument, "invalid")
					st, err := st.WithDetails(&errdetails.DebugInfo{Detail: "debug"})
					if err != nil {
						t.Fatal(err)
					}
					return []any{st}
				}(),
			},
			wantDetails: []Detail{
				{
					TypeURL: "google.rpc.Status",
					Message: &spb.Status{
						Code:    int32(codes.InvalidArgument),
						Message: "invalid",
					},
				},
				{
					TypeURL: "type.googleapis.com/rpc.v1.DebugInfo",
					Message: &errdetails.DebugInfo{Detail: "debug"},
				},
			},
		},
		{
			name: "returns details for *Detail",
			args: args{
				objs: []any{
					&Detail{
						TypeURL: "custom",
						Message: &errdetails.DebugInfo{},
					},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "custom",
					Message: &errdetails.DebugInfo{},
				},
			},
		},
		{
			name: "returns details for Detail",
			args: args{
				objs: []any{
					Detail{
						TypeURL: "custom",
						Message: &errdetails.DebugInfo{},
					},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "custom",
					Message: &errdetails.DebugInfo{},
				},
			},
		},
		{
			name: "returns details for *info.Detail",
			args: args{
				objs: []any{
					&info.Detail{
						Version: "v1",
					},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "rpc.v1.DebugInfo",
					Message: &errdetails.DebugInfo{
						Detail: `{"vald_version":"v1"}`,
					},
				},
			},
		},
		{
			name: "returns details for info.Detail",
			args: args{
				objs: []any{
					info.Detail{
						Version: "v1",
					},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "rpc.v1.DebugInfo",
					Message: &errdetails.DebugInfo{
						Detail: `{"vald_version":"v1"}`,
					},
				},
			},
		},
		{
			name: "returns details for nested slices",
			args: args{
				objs: []any{
					[]any{
						&spb.Status{Code: 2},
					},
					&spb.Status{Code: 3},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "google.rpc.Status",
					Message: &spb.Status{Code: 2},
				},
				{
					TypeURL: "google.rpc.Status",
					Message: &spb.Status{Code: 3},
				},
			},
		},
		{
			name: "returns details for *types.Any",
			args: args{
				objs: func() []any {
					a, _ := anypb.New(&errdetails.DebugInfo{Detail: "test"})
					return []any{a}
				}(),
			},
			wantDetails: []Detail{
				{
					TypeURL: "type.googleapis.com/rpc.v1.DebugInfo",
					Message: &errdetails.DebugInfo{Detail: "test"},
				},
			},
		},
		{
			name: "returns details for types.Any",
			args: args{
				objs: func() []any {
					a, _ := anypb.New(&errdetails.DebugInfo{Detail: "test"})
					return []any{*a}
				}(),
			},
			wantDetails: []Detail{
				{
					TypeURL: "type.googleapis.com/rpc.v1.DebugInfo",
					Message: &errdetails.DebugInfo{Detail: "test"},
				},
			},
		},
		{
			name: "returns details for *proto.Message",
			args: args{
				objs: func() []any {
					var m proto.Message = &errdetails.DebugInfo{Detail: "test"}
					return []any{&m}
				}(),
			},
			wantDetails: []Detail{
				{
					TypeURL: "rpc.v1.DebugInfo",
					Message: &errdetails.DebugInfo{Detail: "test"},
				},
			},
		},
		{
			name: "returns details for proto.Message (implicit)",
			args: args{
				objs: []any{
					&errdetails.DebugInfo{Detail: "test"},
				},
			},
			wantDetails: []Detail{
				{
					TypeURL: "rpc.v1.DebugInfo",
					Message: &errdetails.DebugInfo{Detail: "test"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetails := decodeDetails(tt.args.objs...)
			if len(gotDetails) != len(tt.wantDetails) {
				t.Errorf("decodeDetails() len = %v, want %v", len(gotDetails), len(tt.wantDetails))
				return
			}
			for i := range gotDetails {
				if gotDetails[i].TypeURL != tt.wantDetails[i].TypeURL {
					t.Errorf("decodeDetails()[%d].TypeURL = %v, want %v", i, gotDetails[i].TypeURL, tt.wantDetails[i].TypeURL)
				}
				if !pproto.Equal(gotDetails[i].Message, tt.wantDetails[i].Message) {
					t.Errorf("decodeDetails()[%d].Message = %v, want %v", i, gotDetails[i].Message, tt.wantDetails[i].Message)
				}
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	t.Parallel()
	type args struct {
		objs []any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns empty string for empty input",
			args: args{
				objs: nil,
			},
			want: "",
		},
		{
			name: "returns <nil> for nil input",
			args: args{
				objs: []any{nil},
			},
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Serialize(tt.args.objs...); got != tt.want {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyToErrorDetail(t *testing.T) {
	t.Parallel()
	type args struct {
		a *types.Any
	}
	tests := []struct {
		name string
		args args
		want proto.Message
	}{
		{
			name: "returns nil for nil input",
			args: args{
				a: nil,
			},
			want: nil,
		},
		{
			name: "converts known type (DebugInfo)",
			args: args{
				a: func() *types.Any {
					a, _ := anypb.New(&errdetails.DebugInfo{Detail: "test"})
					return a
				}(),
			},
			want: &errdetails.DebugInfo{Detail: "test"},
		},
		{
			name: "returns original message for unknown type",
			args: args{
				a: func() *types.Any {
					a, _ := anypb.New(&spb.Status{Code: 1})
					return a
				}(),
			},
			want: &spb.Status{Code: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AnyToErrorDetail(tt.args.a)
			if !pproto.Equal(got, tt.want) {
				t.Errorf("AnyToErrorDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugInfoFromInfoDetail(t *testing.T) {
	t.Parallel()
	type args struct {
		v *info.Detail
	}
	tests := []struct {
		name string
		args args
		want *DebugInfo
	}{
		{
			name: "converts info.Detail to DebugInfo",
			args: args{
				v: &info.Detail{
					Version: "v1",
				},
			},
			want: &DebugInfo{
				Detail: `{"vald_version":"v1"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DebugInfoFromInfoDetail(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebugInfoFromInfoDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestDetail_MarshalJSON(t *testing.T) {
// 	type fields struct {
// 		TypeURL string
// 		Message proto.Message
// 	}
// 	type want struct {
// 		wantBody []byte
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []byte, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotBody []byte, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotBody, w.wantBody) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotBody, w.wantBody)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           TypeURL:"",
// 		           Message:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           TypeURL:"",
// 		           Message:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &Detail{
// 				TypeURL: test.fields.TypeURL,
// 				Message: test.fields.Message,
// 			}
//
// 			gotBody, err := d.MarshalJSON()
// 			if err := checkFunc(test.want, gotBody, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
