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

package unit

import (
	stderrs "errors"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestParseBytes(t *testing.T) {
	type args struct {
		bs string
	}
	type want struct {
		wantBytes uint64
		err       error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, uint64, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotBytes uint64, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotBytes, w.wantBytes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotBytes, w.wantBytes)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns (1, nil) when bs is `1M`",
			args: args{
				bs: "1M",
			},
			want: want{
				wantBytes: 1048576,
				err:       nil,
			},
		},
		{
			name: "returns (0, nil) when bs is empty",
			want: want{
				wantBytes: 0,
				err:       nil,
			},
		},

		{
			name: "returns (0, nil) when bs is `0`",
			args: args{
				bs: "0",
			},
			want: want{
				wantBytes: 0,
				err:       nil,
			},
		},

		{
			name: "returns (0, error) when bs is `a`",
			args: args{
				bs: "a",
			},
			want: want{
				wantBytes: 0,
				err: func() (err error) {
					err = stderrs.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")
					err = errors.Wrap(err, errors.ErrParseUnitFailed("a").Error())
					return
				}(),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

			gotBytes, err := ParseBytes(test.args.bs)
			if err := checkFunc(test.want, gotBytes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
