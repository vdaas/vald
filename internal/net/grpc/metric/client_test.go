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

// Package metric provides metrics functions for grpc
package metric

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.opencensus.io/trace"
)

func TestNewClientHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []ClientOption
	}
	type want struct {
		want *ClientHandler
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *ClientHandler) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *ClientHandler) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return client handler when option is nil",
			args: args{
				opts: nil,
			},
			want: want{
				want: new(ClientHandler),
			},
		},
		{
			name: "return client handler when option is not nil",
			args: args{
				opts: []ClientOption{
					func(h *ClientHandler) {
						h.StartOptions = trace.StartOptions{
							SpanKind: 1,
						}
					},
				},
			},
			want: want{
				want: &ClientHandler{
					StartOptions: trace.StartOptions{
						SpanKind: 1,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := NewClientHandler(test.args.opts...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
