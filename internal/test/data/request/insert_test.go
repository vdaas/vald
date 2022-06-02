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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
)

var defaultMultiInsertReqComparators = []cmp.Option{
	comparator.IgnoreUnexported(payload.Insert_Request{}),
	comparator.IgnoreUnexported(payload.Insert_MultiRequest{}),
	comparator.IgnoreUnexported(payload.Object_Vector{}),
	comparator.IgnoreUnexported(payload.Insert_Config{}),
}

func TestGenMultiInsertReq(t *testing.T) {
	type args struct {
		t    ObjectType
		dist vector.Distribution
		num  int
		dim  int
		cfg  *payload.Insert_Config
	}
	type want struct {
		want *payload.Insert_MultiRequest
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Insert_MultiRequest, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	dim := 10
	comparators := append(defaultMultiInsertReqComparators, comparator.IgnoreFields(payload.Object_Vector{}, "Vector"))

	defaultCheckFunc := func(w want, got *payload.Insert_MultiRequest, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(got, w.want, comparators...); diff != "" {
			return errors.Errorf("diff: %v", diff)
		}
		if len(got.Requests) != 0 && len(got.Requests[0].Vector.Vector) != dim {
			return errors.New("vector length not match")
		}
		return nil
	}
	tests := []test{
		{
			name: "success to generate 1 float request",
			args: args{
				t:    Float,
				dist: vector.Gaussian,
				num:  1,
				dim:  dim,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-1",
							},
						},
					},
				},
			},
		},
		{
			name: "success to generate 1 uint8 request",
			args: args{
				t:    Uint8,
				dist: vector.Gaussian,
				num:  1,
				dim:  dim,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-1",
							},
						},
					},
				},
			},
		},
		{
			name: "success to generate 5 float request",
			args: args{
				t:    Float,
				dist: vector.Gaussian,
				num:  5,
				dim:  dim,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-1",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-2",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-3",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-4",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-5",
							},
						},
					},
				},
			},
		},
		{
			name: "success to generate 5 uint8 request",
			args: args{
				t:    Uint8,
				dist: vector.Gaussian,
				num:  5,
				dim:  dim,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-1",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-2",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-3",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-4",
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id: "uuid-5",
							},
						},
					},
				},
			},
		},
		{
			name: "success to generate 0 float request",
			args: args{
				t:    Float,
				dist: vector.Gaussian,
				num:  0,
				dim:  dim,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{},
				},
			},
		},
		{
			name: "success to generate 0 uint8 request",
			args: args{
				t:    Uint8,
				dist: vector.Gaussian,
				num:  0,
				dim:  dim,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{},
				},
			},
		},
		// max num and max dim test is ignored due to test timeout
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

			got, err := GenMultiInsertReq(test.args.t, test.args.dist, test.args.num, test.args.dim, test.args.cfg)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGenSameVecMultiInsertReq(t *testing.T) {
	type args struct {
		num int
		vec []float32
		cfg *payload.Insert_Config
	}
	type want struct {
		want *payload.Insert_MultiRequest
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Insert_MultiRequest) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *payload.Insert_MultiRequest) error {
		if diff := comparator.Diff(got, w.want, defaultMultiInsertReqComparators...); diff != "" {
			return errors.Errorf("diff: %v", diff)
		}
		return nil
	}
	tests := []test{
		func() test {
			vecs, err := vector.GenF32Vec(vector.Gaussian, 1, 10)
			if err != nil {
				t.Error(err)
			}
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			return test{
				name: "success to generate 1 same vector request",
				args: args{
					num: 1,
					vec: vecs[0],
					cfg: cfg,
				},
				want: want{
					want: &payload.Insert_MultiRequest{
						Requests: []*payload.Insert_Request{
							{
								Vector: &payload.Object_Vector{
									Id:     "uuid-1",
									Vector: vecs[0],
								},
								Config: cfg,
							},
						},
					},
				},
			}
		}(),
		func() test {
			vecs, err := vector.GenF32Vec(vector.Gaussian, 1, 10)
			if err != nil {
				t.Error(err)
			}
			cfg := &payload.Insert_Config{
				SkipStrictExistCheck: true,
			}

			return test{
				name: "success to generate 5 same vector request",
				args: args{
					num: 5,
					vec: vecs[0],
					cfg: cfg,
				},
				want: want{
					want: &payload.Insert_MultiRequest{
						Requests: []*payload.Insert_Request{
							{
								Vector: &payload.Object_Vector{
									Id:     "uuid-1",
									Vector: vecs[0],
								},
								Config: cfg,
							},
							{
								Vector: &payload.Object_Vector{
									Id:     "uuid-2",
									Vector: vecs[0],
								},
								Config: cfg,
							},
							{
								Vector: &payload.Object_Vector{
									Id:     "uuid-3",
									Vector: vecs[0],
								},
								Config: cfg,
							},
							{
								Vector: &payload.Object_Vector{
									Id:     "uuid-4",
									Vector: vecs[0],
								},
								Config: cfg,
							},
							{
								Vector: &payload.Object_Vector{
									Id:     "uuid-5",
									Vector: vecs[0],
								},
								Config: cfg,
							},
						},
					},
				},
			}
		}(),
		func() test {
			vecs, err := vector.GenF32Vec(vector.Gaussian, 1, 10)
			if err != nil {
				t.Error(err)
			}
			var cfg *payload.Insert_Config

			return test{
				name: "success to generate 1 same vector request when cfg is nil",
				args: args{
					num: 1,
					vec: vecs[0],
					cfg: cfg,
				},
				want: want{
					want: &payload.Insert_MultiRequest{
						Requests: []*payload.Insert_Request{
							{
								Vector: &payload.Object_Vector{
									Id:     "uuid-1",
									Vector: vecs[0],
								},
								Config: cfg,
							},
						},
					},
				},
			}
		}(),
		{
			name: "success to generate 0 same vector request",
			args: args{
				num: 0,
				vec: []float32{1, 2, 3},
				cfg: nil,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{},
				},
			},
		},
		{
			name: "success to generate empty vector request",
			args: args{
				num: 1,
				vec: []float32{},
				cfg: nil,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{
						{
							Vector: &payload.Object_Vector{
								Id:     "uuid-1",
								Vector: []float32{},
							},
						},
					},
				},
			},
		},
		{
			name: "success to generate multiple empty vector request",
			args: args{
				num: 5,
				vec: []float32{},
				cfg: nil,
			},
			want: want{
				want: &payload.Insert_MultiRequest{
					Requests: []*payload.Insert_Request{
						{
							Vector: &payload.Object_Vector{
								Id:     "uuid-1",
								Vector: []float32{},
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id:     "uuid-2",
								Vector: []float32{},
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id:     "uuid-3",
								Vector: []float32{},
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id:     "uuid-4",
								Vector: []float32{},
							},
						},
						{
							Vector: &payload.Object_Vector{
								Id:     "uuid-5",
								Vector: []float32{},
							},
						},
					},
				},
			},
		},
		// max num test is ignored due to test timeout
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

			got := GenSameVecMultiInsertReq(test.args.num, test.args.vec, test.args.cfg)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
