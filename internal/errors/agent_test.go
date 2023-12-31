//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package errors
package errors

import (
	"testing"
)

func TestErrObjectNotFound(t *testing.T) {
	type args struct {
		err  error
		uuid string
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
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrObjectNotFound error when err is agent error and uuid is 550e8400-e29b-41d4",
			args: args{
				err:  New("agent error"),
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("uuid 550e8400-e29b-41d4's object not found: agent error"),
			},
		},
		{
			name: "return a wrapped ErrObjectNotFound error when err is agent error and uuid is empty",
			args: args{
				err:  New("agent error"),
				uuid: "",
			},
			want: want{
				want: New("uuid 's object not found: agent error"),
			},
		},
		{
			name: "return an ErrObjectNotFound error when err is nil and uuid is 550e8400-e29b-41d4",
			args: args{
				err:  nil,
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("uuid 550e8400-e29b-41d4's object not found"),
			},
		},
		{
			name: "return an ErrObjectNotFound error when err is nil and uuid is empty",
			args: args{
				err:  nil,
				uuid: "",
			},
			want: want{
				want: New("uuid 's object not found"),
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
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := ErrObjectNotFound(test.args.err, test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
