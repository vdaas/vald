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

// Package credentials provides generic functionality for grpc credentials setting
package credentials

import (
	"crypto/tls"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"google.golang.org/grpc/credentials"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNewTLS(t *testing.T) {
	t.Parallel()
	type args struct {
		c *tls.Config
	}
	type want struct {
		want credentials.TransportCredentials
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, credentials.TransportCredentials) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got credentials.TransportCredentials) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return credential when config is nil",
			args: args{
				c: nil,
			},
			want: want{
				want: credentials.NewTLS(nil),
			},
		},
		{
			name: "return credential when config is not nil",
			args: args{
				c: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			},
			want: want{
				want: credentials.NewTLS(&tls.Config{
					MinVersion: tls.VersionTLS12,
				}),
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

			got := NewTLS(test.args.c)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
