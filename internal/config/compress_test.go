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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_compressAlgorithm_String(t *testing.T) {
	type want struct {
		want string
	}
	type test struct {
		name       string
		ca         compressAlgorithm
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return gob when compressAlgorithm is GOB",
			ca:   GOB,
			want: want{
				want: "gob",
			},
		},
		{
			name: "return gzip when compressAlgorithm is GZIP",
			ca:   GZIP,
			want: want{
				want: "gzip",
			},
		},
		{
			name: "return lz4 when compressAlgorithm is LZ4",
			ca:   LZ4,
			want: want{
				want: "lz4",
			},
		},
		{
			name: "return zstd when compressAlgorithm is ZSTD",
			ca:   ZSTD,
			want: want{
				want: "zstd",
			},
		},
		{
			name: "return unknown when compressAlgorithm is the default value of uint8",
			want: want{
				want: "unknown",
			},
		},
		{
			name: "return unknown when compressAlgorithm is 100",
			ca:   compressAlgorithm(100),
			want: want{
				want: "unknown",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := test.ca.String()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCompressAlgorithm(t *testing.T) {
	type args struct {
		ca string
	}
	type want struct {
		want compressAlgorithm
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, compressAlgorithm) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got compressAlgorithm) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return GOB when ca is gob",
			args: args{
				ca: "gob",
			},
			want: want{
				want: GOB,
			},
		},
		{
			name: "return GOB when ca is gOB",
			args: args{
				ca: "gOB",
			},
			want: want{
				want: GOB,
			},
		},
		{
			name: "return GZIP when ca is gzip",
			args: args{
				ca: "gzip",
			},
			want: want{
				want: GZIP,
			},
		},
		{
			name: "return GZIP when ca is gZIP",
			args: args{
				ca: "gZIP",
			},
			want: want{
				want: GZIP,
			},
		},
		{
			name: "return LZ4 when ca is lz4",
			args: args{
				ca: "lz4",
			},
			want: want{
				want: LZ4,
			},
		},
		{
			name: "return LZ4 when ca is lZ4",
			args: args{
				ca: "lZ4",
			},
			want: want{
				want: LZ4,
			},
		},
		{
			name: "return ZSTD when ca is zstd",
			args: args{
				ca: "zstd",
			},
			want: want{
				want: ZSTD,
			},
		},
		{
			name: "return ZSTD when ca is zSTD",
			args: args{
				ca: "zSTD",
			},
			want: want{
				want: ZSTD,
			},
		},
		{
			name: "return 0 when ca is empty",
			want: want{
				want: 0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := CompressAlgorithm(test.args.ca)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCompressCore_Bind(t *testing.T) {
	type fields struct {
		CompressAlgorithm string
		CompressionLevel  int
	}
	type want struct {
		want *CompressCore
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *CompressCore) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *CompressCore) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return CompressCore when the bind successes",
				fields: fields{
					CompressAlgorithm: "gob",
				},
				want: want{
					want: &CompressCore{
						CompressAlgorithm: "gob",
					},
				},
			}
		}(),
		func() test {
			key := "COMPRESSCORE_BIND_COMPRESS_ALGORITHM"
			wantVal := "gzip"

			return test{
				name: "return CompressCore when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					CompressAlgorithm: "_" + key + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					t.Setenv(key, wantVal)
				},
				want: want{
					want: &CompressCore{
						CompressAlgorithm: wantVal,
					},
				},
			}
		}(),
		func() test {
			key := "COMPRESSCORE_BIND_COMPRESS_ALGORITHM_EMPTY"
			wantVal := ""

			return test{
				name: "return CompressCore when the bind successes but loaded environment variable is empty",
				fields: fields{
					CompressAlgorithm: "_" + key + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					t.Setenv(key, wantVal)
				},
				want: want{
					want: &CompressCore{
						CompressAlgorithm: wantVal,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return CompressCore when the bind successes but loaded environment variable is not found",
				fields: fields{
					CompressAlgorithm: "_COMPRESSCORE_BIND_COMPRESS_ALGORITHM_NOT_EXISTS_",
				},
				want: want{
					want: &CompressCore{
						CompressAlgorithm: "",
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return CompressCore when the bind successes but the field is empty",
				want: want{
					want: &CompressCore{},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &CompressCore{
				CompressAlgorithm: test.fields.CompressAlgorithm,
				CompressionLevel:  test.fields.CompressionLevel,
			}

			got := c.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
