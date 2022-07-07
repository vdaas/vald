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
package request

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestGenMultiRemoveReq(t *testing.T) {
	type args struct {
		num int
		cfg *payload.Remove_Config
	}
	type want struct {
		want *payload.Remove_MultiRequest
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Remove_MultiRequest) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *payload.Remove_MultiRequest) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultRemoveCfg := &payload.Remove_Config{
		SkipStrictExistCheck: true,
	}
	tests := []test{
		{
			name: "success to generate 1 remove request",
			args: args{
				num: 1,
				cfg: defaultRemoveCfg,
			},
			want: want{
				want: &payload.Remove_MultiRequest{
					Requests: []*payload.Remove_Request{
						{
							Id: &payload.Object_ID{
								Id: "uuid-1",
							},
							Config: defaultRemoveCfg,
						},
					},
				},
			},
		},
		{
			name: "success to generate 5 remove request",
			args: args{
				num: 5,
				cfg: defaultRemoveCfg,
			},
			want: want{
				want: &payload.Remove_MultiRequest{
					Requests: []*payload.Remove_Request{
						{
							Id: &payload.Object_ID{
								Id: "uuid-1",
							},
							Config: defaultRemoveCfg,
						},
						{
							Id: &payload.Object_ID{
								Id: "uuid-2",
							},
							Config: defaultRemoveCfg,
						},
						{
							Id: &payload.Object_ID{
								Id: "uuid-3",
							},
							Config: defaultRemoveCfg,
						},
						{
							Id: &payload.Object_ID{
								Id: "uuid-4",
							},
							Config: defaultRemoveCfg,
						},
						{
							Id: &payload.Object_ID{
								Id: "uuid-5",
							},
							Config: defaultRemoveCfg,
						},
					},
				},
			},
		},
		{
			name: "success to generate 1 remove request with cfg is nil",
			args: args{
				num: 1,
				cfg: nil,
			},
			want: want{
				want: &payload.Remove_MultiRequest{
					Requests: []*payload.Remove_Request{
						{
							Id: &payload.Object_ID{
								Id: "uuid-1",
							},
							Config: nil,
						},
					},
				},
			},
		},
		{
			name: "success to generate 0 remove request",
			args: args{
				num: 0,
				cfg: nil,
			},
			want: want{
				want: &payload.Remove_MultiRequest{
					Requests: []*payload.Remove_Request{},
				},
			},
		},
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

			got := GenMultiRemoveReq(test.args.num, test.args.cfg)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
