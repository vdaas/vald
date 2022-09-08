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

package cassandra

import (
	"reflect"
	"testing"

	"github.com/gocql/gocql"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewConvictionPolicy(t *testing.T) {
	type want struct {
		want gocql.ConvictionPolicy
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, gocql.ConvictionPolicy) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got gocql.ConvictionPolicy) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return conviction policy success",
			want: want{
				want: new(convictionPolicy),
			},
		},
	}

	for _, tc := range tests {
		test := tc
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

			got := NewConvictionPolicy()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_convictionPolicy_AddFailure(t *testing.T) {
	type args struct {
		err  error
		host *gocql.HostInfo
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		c          *convictionPolicy
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "AddFailure successfully log the error",
			args: args{
				err: errors.New("dummy"),
				host: func() *gocql.HostInfo {
					h := &gocql.HostInfo{}
					h.SetConnectAddress(net.IPv4(127, 0, 0, 1))
					return h
				}(),
			},
			want: want{
				want: true,
			},
		},
	}

	for _, tc := range tests {
		test := tc
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
			c := &convictionPolicy{}

			got := c.AddFailure(test.args.err, test.args.host)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_convictionPolicy_Reset(t *testing.T) {
	type args struct {
		host *gocql.HostInfo
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		c          *convictionPolicy
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		{
			name: "reset execute success",
			args: args{
				host: func() *gocql.HostInfo {
					h := &gocql.HostInfo{}
					h.SetConnectAddress(net.IPv4(127, 0, 0, 1))
					return h
				}(),
			},
		},
	}

	for _, tc := range tests {
		test := tc
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
			c := &convictionPolicy{}

			c.Reset(test.args.host)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
