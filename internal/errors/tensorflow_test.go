//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

package errors

import (
	"math"
	"testing"

	"go.uber.org/goleak"
)

func TestErrFailedToCastTF(t *testing.T) {
	t.Parallel()
	type args struct {
		v interface{}
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrFailedToCastTF error when v is nil",
			args: args{},
			want: want{
				want: New("failed to cast tensorflow result <nil>"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is int value",
			args: args{
				v: func() interface{} {
					var v int
					return v
				}(),
			},
			want: want{
				want: New("failed to cast tensorflow result 0"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MaxInt64 value",
			args: args{
				v: func() interface{} {
					return math.MaxInt64
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", math.MaxInt64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MaxFloat64 value",
			args: args{
				v: func() interface{} {
					return math.MaxFloat64
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", math.MaxFloat64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MinInt64 value",
			args: args{
				v: func() interface{} {
					return math.MinInt64
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", math.MinInt64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is string value",
			args: args{
				v: func() interface{} {
					return "tensorflowObject"
				}(),
			},
			want: want{
				want: New("failed to cast tensorflow result tensorflowObject"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is struct value",
			args: args{
				v: func() interface{} {
					type str struct {
						uuid   int
						vector []float64
					}
					return str{
						uuid:   1,
						vector: make([]float64, 5),
					}
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", struct {
					uuid   int
					vector []float64
				}{
					uuid:   1,
					vector: make([]float64, 5),
				}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is map value",
			args: args{
				v: func() interface{} {
					return map[string]int{"x": 0, "y": int(math.MaxInt64)}
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", map[string]int{
					"x": 0,
					"y": int(math.MaxInt64),
				}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is byte value",
			args: args{
				v: func() interface{} {
					return []byte{0x00, 0x01}
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", []byte{0x00, 0x01}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is bool value",
			args: args{
				v: func() interface{} {
					return true
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", true),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is uintptr value",
			args: args{
				v: func() interface{} {
					var v uintptr
					return v
				}(),
			},
			want: want{
				want: Errorf("failed to cast tensorflow result %+v", 0),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrFailedToCastTF(test.args.v)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInputLength(t *testing.T) {
	t.Parallel()
	type args struct {
		iLength int
		fLength int
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrInputLength error when i s 10 and f is 9",
			args: args{
				iLength: 10,
				fLength: 9,
			},
			want: want{
				want: Errorf("inputs length %d does not match feeds length %d", 10, 9),
			},
		},
		{
			name: "return an ErrInputLength error when i s 10 and f is int(math.MaxInt64)",
			args: args{
				iLength: 10,
				fLength: int(math.MaxInt64),
			},
			want: want{
				want: Errorf("inputs length %d does not match feeds length %d", 10, int(math.MaxInt64)),
			},
		},
		{
			name: "return an ErrInputLength error when i s 10 and f is int(math.MinInt64)",
			args: args{
				iLength: 10,
				fLength: int(math.MinInt64),
			},
			want: want{
				want: Errorf("inputs length %d does not match feeds length %d", 10, int(math.MinInt64)),
			},
		},
		{
			name: "return an ErrInputLength error when i and f is zero value",
			args: args{},
			want: want{
				want: Errorf("inputs length %d does not match feeds length %d", 0, 0),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrInputLength(test.args.iLength, test.args.fLength)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrNilTensorTF(t *testing.T) {
	t.Parallel()
	type args struct {
		v interface{}
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrFailedToCastTF error when v is nil",
			args: args{},
			want: want{
				want: New("nil tensorflow tensor <nil>"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is int value",
			args: args{
				v: func() interface{} {
					var v int
					return v
				}(),
			},
			want: want{
				want: New("nil tensorflow tensor 0"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MaxInt64 value",
			args: args{
				v: func() interface{} {
					return math.MaxInt64
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", math.MaxInt64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MaxFloat64 value",
			args: args{
				v: func() interface{} {
					return math.MaxFloat64
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", math.MaxFloat64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MinInt64 value",
			args: args{
				v: func() interface{} {
					return math.MinInt64
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", math.MinInt64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is string value",
			args: args{
				v: func() interface{} {
					return "tensorflowObject"
				}(),
			},
			want: want{
				want: New("nil tensorflow tensor tensorflowObject"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is struct value",
			args: args{
				v: func() interface{} {
					type str struct {
						uuid   int
						vector []float64
					}
					return str{
						uuid:   1,
						vector: make([]float64, 5),
					}
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", struct {
					uuid   int
					vector []float64
				}{
					uuid:   1,
					vector: make([]float64, 5),
				}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is map value",
			args: args{
				v: func() interface{} {
					return map[string]int{"x": 0, "y": int(math.MaxInt64)}
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", map[string]int{
					"x": 0,
					"y": int(math.MaxInt64),
				}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is byte value",
			args: args{
				v: func() interface{} {
					return []byte{0x00, 0x01}
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", []byte{0x00, 0x01}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is bool value",
			args: args{
				v: func() interface{} {
					return true
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", true),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is uintptr value",
			args: args{
				v: func() interface{} {
					var v uintptr
					return v
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor %+v", 0),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrNilTensorTF(test.args.v)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrNilTensorValueTF(t *testing.T) {
	t.Parallel()
	type args struct {
		v interface{}
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrFailedToCastTF error when v is nil",
			args: args{},
			want: want{
				want: New("nil tensorflow tensor value <nil>"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is int value",
			args: args{
				v: func() interface{} {
					var v int
					return v
				}(),
			},
			want: want{
				want: New("nil tensorflow tensor value 0"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MaxInt64 value",
			args: args{
				v: func() interface{} {
					return math.MaxInt64
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", math.MaxInt64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MaxFloat64 value",
			args: args{
				v: func() interface{} {
					return math.MaxFloat64
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", math.MaxFloat64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is MinInt64 value",
			args: args{
				v: func() interface{} {
					return math.MinInt64
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", math.MinInt64),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is string value",
			args: args{
				v: func() interface{} {
					return "tensorflowObject"
				}(),
			},
			want: want{
				want: New("nil tensorflow tensor value tensorflowObject"),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is struct value",
			args: args{
				v: func() interface{} {
					type str struct {
						uuid   int
						vector []float64
					}
					return str{
						uuid:   1,
						vector: make([]float64, 5),
					}
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", struct {
					uuid   int
					vector []float64
				}{
					uuid:   1,
					vector: make([]float64, 5),
				}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is map value",
			args: args{
				v: func() interface{} {
					return map[string]int{"x": 0, "y": int(math.MaxInt64)}
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", map[string]int{
					"x": 0,
					"y": int(math.MaxInt64),
				}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is byte value",
			args: args{
				v: func() interface{} {
					return []byte{0x00, 0x01}
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", []byte{0x00, 0x01}),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is bool value",
			args: args{
				v: func() interface{} {
					return true
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", true),
			},
		},
		{
			name: "return an ErrFailedToCastTF error when v is uintptr value",
			args: args{
				v: func() interface{} {
					var v uintptr
					return v
				}(),
			},
			want: want{
				want: Errorf("nil tensorflow tensor value %+v", 0),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrNilTensorValueTF(test.args.v)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
