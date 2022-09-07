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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var serverComparer = []comparator.Option{
	comparator.AllowUnexported(Server{}),
	comparator.IgnoreFields(Server{}, "opts", "quit", "done", "channelzRemoveOnce", "czData", "channelzID"),
	comparator.MutexComparer,
	comparator.CondComparer,
	comparator.WaitGroupComparer,
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNewServer(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []ServerOption
	}
	type want struct {
		want *Server
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *Server) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *Server) error {
		if diff := comparator.Diff(w.want, got, serverComparer...); diff != "" {
			return errors.Errorf(diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns gRPC server when no option is set",
			args: args{
				opts: nil,
			},
			want: want{
				want: grpc.NewServer(),
			},
		},
		{
			name: "returns gRPC server when 1 option is set",
			args: args{
				opts: []ServerOption{MaxSendMsgSize(100)},
			},
			want: want{
				want: grpc.NewServer(),
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

			got := NewServer(test.args.opts...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCreds(t *testing.T) {
	t.Parallel()
	type args struct {
		c credentials.TransportCredentials
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when cred is nil",
			args: args{
				c: nil,
			},
		},
		{
			name: "return server option when cred is not nil",
			args: args{
				c: credentials.NewTLS(nil),
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

			got := Creds(test.args.c)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestKeepaliveParams(t *testing.T) {
	t.Parallel()
	type args struct {
		kp keepalive.ServerParameters
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when keepalive is empty",
			args: args{
				kp: keepalive.ServerParameters{},
			},
		},
		{
			name: "return server option when keepalive is not empty",
			args: args{
				kp: keepalive.ServerParameters{
					MaxConnectionIdle: time.Second,
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

			got := KeepaliveParams(test.args.kp)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMaxRecvMsgSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size int
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := MaxRecvMsgSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMaxSendMsgSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size int
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := MaxSendMsgSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInitialWindowSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size int32
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := InitialWindowSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInitialConnWindowSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size int32
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := InitialConnWindowSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestReadBufferSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size int
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := ReadBufferSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWriteBufferSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size int
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := WriteBufferSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestConnectionTimeout(t *testing.T) {
	t.Parallel()
	type args struct {
		d time.Duration
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when timeout is 0",
			args: args{
				d: 0,
			},
		},
		{
			name: "return server option when timeout is not 0",
			args: args{
				d: time.Second * 30,
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

			got := ConnectionTimeout(test.args.d)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMaxHeaderListSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size uint32
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := MaxHeaderListSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestHeaderTableSize(t *testing.T) {
	t.Parallel()
	type args struct {
		size uint32
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when size is 0",
			args: args{
				size: 0,
			},
		},
		{
			name: "return server option when size is not 0",
			args: args{
				size: 100,
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

			got := HeaderTableSize(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestKeepaliveEnforcementPolicy(t *testing.T) {
	t.Parallel()
	type args struct {
		kep keepalive.EnforcementPolicy
	}
	type want struct {
		want ServerOption
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ServerOption) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ServerOption) error {
		if got == nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return server option when policy is empty",
			args: args{
				kep: keepalive.EnforcementPolicy{},
			},
		},
		{
			name: "return server option when policy is not empty",
			args: args{
				kep: keepalive.EnforcementPolicy{
					MinTime: time.Second,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := KeepaliveEnforcementPolicy(test.args.kep)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
