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

// Package control provides network socket option
package control

import (
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		flag      SocketFlag
		keepAlive int
	}
	type want struct {
		want SocketController
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, SocketController) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got SocketController) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return control when the socket flag and keepalive is 0",
			args: args{
				flag:      0,
				keepAlive: 0,
			},
			want: want{
				want: &control{},
			},
		},
		{
			name: "return control when the socket flag and keepalive is set",
			args: args{
				flag:      ReuseAddr | TCPNoDelay | IPTransparent,
				keepAlive: int(time.Second) * 60,
			},
			want: want{
				want: &control{
					reuseAddr:     true,
					tcpNoDelay:    true,
					ipTransparent: true,
					keepAlive:     int(time.Second) * 60,
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

			got := New(test.args.flag, test.args.keepAlive)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_boolint(t *testing.T) {
	t.Parallel()
	type args struct {
		b bool
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, int) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got int) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return 1 when bool is true",
			args: args{
				b: true,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return 0 when bool is false",
			args: args{
				b: false,
			},
			want: want{
				want: 0,
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

			got := boolint(test.args.b)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_isTCP(t *testing.T) {
	t.Parallel()
	type args struct {
		network string
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
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
			name: "return true when network is tcp",
			args: args{
				network: "tcp",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when network is tcp4",
			args: args{
				network: "tcp4",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when network is tcp6",
			args: args{
				network: "tcp6",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false when network is not tcp",
			args: args{
				network: "udp",
			},
			want: want{
				want: false,
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

			got := isTCP(test.args.network)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_control_GetControl(t *testing.T) {
	t.Parallel()
	type fields struct {
		reusePort                bool
		reuseAddr                bool
		tcpFastOpen              bool
		tcpNoDelay               bool
		tcpCork                  bool
		tcpQuickAck              bool
		tcpDeferAccept           bool
		ipTransparent            bool
		ipRecoverDestinationAddr bool
		keepAlive                int
	}
	type want struct {
		want func(network, addr string, c syscall.RawConn) (err error)
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, func(network, addr string, c syscall.RawConn) (err error)) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got func(network, addr string, c syscall.RawConn) (err error)) error {
		if reflect.ValueOf(w.want).Pointer() != reflect.ValueOf(got).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name:   "return control func success",
			fields: fields{},
			want: want{
				want: new(control).controlFunc,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			ctrl := &control{
				reusePort:                test.fields.reusePort,
				reuseAddr:                test.fields.reuseAddr,
				tcpFastOpen:              test.fields.tcpFastOpen,
				tcpNoDelay:               test.fields.tcpNoDelay,
				tcpCork:                  test.fields.tcpCork,
				tcpQuickAck:              test.fields.tcpQuickAck,
				tcpDeferAccept:           test.fields.tcpDeferAccept,
				ipTransparent:            test.fields.ipTransparent,
				ipRecoverDestinationAddr: test.fields.ipRecoverDestinationAddr,
				keepAlive:                test.fields.keepAlive,
			}

			got := ctrl.GetControl()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_control_controlFunc(t *testing.T) {
	t.Parallel()
	type args struct {
		network string
		address string
		c       syscall.RawConn
	}
	type fields struct {
		reusePort                bool
		reuseAddr                bool
		tcpFastOpen              bool
		tcpNoDelay               bool
		tcpCork                  bool
		tcpQuickAck              bool
		tcpDeferAccept           bool
		ipTransparent            bool
		ipRecoverDestinationAddr bool
		keepAlive                int
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			f, err := file.Open(".", 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			sc, err := f.SyscallConn()
			if err != nil {
				t.Fatal(err)
			}

			return test{
				name: "set socket options success when network is tcp",
				args: args{
					network: "tcp",
					address: "127.0.0.1",
					c:       sc,
				},
				fields: fields{
					reusePort:                true,
					reuseAddr:                true,
					tcpFastOpen:              true,
					tcpNoDelay:               true,
					tcpCork:                  true,
					tcpQuickAck:              true,
					tcpDeferAccept:           true,
					ipTransparent:            true,
					ipRecoverDestinationAddr: true,
					keepAlive:                10,
				},
				want: want{},
				afterFunc: func(a args) {
					f.Close()
				},
			}
		}(),
		func() test {
			f, err := file.Open(".", 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			sc, err := f.SyscallConn()
			if err != nil {
				t.Fatal(err)
			}

			return test{
				name: "set socket options success when network is tcp6",
				args: args{
					network: "tcp6",
					address: "::1",
					c:       sc,
				},
				fields: fields{
					reusePort:                true,
					reuseAddr:                true,
					tcpFastOpen:              true,
					tcpNoDelay:               true,
					tcpCork:                  true,
					tcpQuickAck:              true,
					tcpDeferAccept:           true,
					ipTransparent:            true,
					ipRecoverDestinationAddr: true,
					keepAlive:                10,
				},
				want: want{},
				afterFunc: func(a args) {
					f.Close()
				},
			}
		}(),
		func() test {
			f, err := file.Open(".", 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			sc, err := f.SyscallConn()
			if err != nil {
				t.Fatal(err)
			}

			return test{
				name: "set socket options success when network is file",
				args: args{
					network: "file",
					address: ".",
					c:       sc,
				},
				fields: fields{
					reusePort:                true,
					reuseAddr:                true,
					tcpFastOpen:              true,
					tcpNoDelay:               true,
					tcpCork:                  true,
					tcpQuickAck:              true,
					tcpDeferAccept:           true,
					ipTransparent:            true,
					ipRecoverDestinationAddr: true,
					keepAlive:                10,
				},
				want: want{},
				afterFunc: func(a args) {
					f.Close()
				},
			}
		}(),
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
			ctrl := &control{
				reusePort:                test.fields.reusePort,
				reuseAddr:                test.fields.reuseAddr,
				tcpFastOpen:              test.fields.tcpFastOpen,
				tcpNoDelay:               test.fields.tcpNoDelay,
				tcpCork:                  test.fields.tcpCork,
				tcpQuickAck:              test.fields.tcpQuickAck,
				tcpDeferAccept:           test.fields.tcpDeferAccept,
				ipTransparent:            test.fields.ipTransparent,
				ipRecoverDestinationAddr: test.fields.ipRecoverDestinationAddr,
				keepAlive:                test.fields.keepAlive,
			}

			err := ctrl.controlFunc(test.args.network, test.args.address, test.args.c)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
