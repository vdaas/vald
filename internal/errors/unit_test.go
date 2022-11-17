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
package errors

import (
	"testing"
)

func TestErrParseUnitFailed(t *testing.T) {
	type args struct {
		str string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\", \n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrParseUnitFailed error when str is not empty.",
			args: args{
				str: "parse target string",
			},
			want: want{
				want: New("failed to parse: 'parse target string'"),
			},
		},
		{
			name: "return an ErrParseUnitFailed error when str is empty.",
			args: args{},
			want: want{
				want: New("failed to parse: ''"),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := ErrParseUnitFailed(test.args.str)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
