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
	"math"
	"testing"
)

func TestErrGRPCClientConnectionClose(t *testing.T) {
	type args struct {
		name string
		err  error
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
			name: "return wrapped ErrGRPCClientConnectionClose error when err is server error and name is 'gateway'",
			args: args{
				err:  New("server error"),
				name: "gateway",
			},
			want: want{
				want: New("gateway's gRPC connection close error: server error"),
			},
		},
		{
			name: "return wrapped ErrGRPCClientConnectionClose error when err is server error and name is empty",
			args: args{
				err:  New("server error"),
				name: "",
			},
			want: want{
				want: New("'s gRPC connection close error: server error"),
			},
		},
		{
			name: "return ErrGRPCClientConnectionClose error when err is nil and name is 'gateway'",
			args: args{
				err:  nil,
				name: "gateway",
			},
			want: want{
				want: New("gateway's gRPC connection close error"),
			},
		},
		{
			name: "return ErrGRPCClientConnectionClose error when err is nil and addr is empty",
			args: args{
				err:  nil,
				name: "",
			},
			want: want{
				want: New("'s gRPC connection close error"),
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

			got := ErrGRPCClientConnectionClose(test.args.name, test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidGRPCPort(t *testing.T) {
	type args struct {
		addr string
		host string
		port uint16
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
			name: "return ErrInvalidGRPCPort error when addr is '127.0.0.1' and host is 'gateway.default.svc.cluster.local' and port is '8080'",
			args: args{
				addr: "127.0.0.1",
				host: "gateway.default.svc.cluster.local",
				port: 8080,
			},
			want: want{
				want: New("invalid gRPC client connection port to addr: 127.0.0.1,\thost: gateway.default.svc.cluster.local\t port: 8080"),
			},
		},
		{
			name: "return ErrInvalidGRPCPort error when addr is empty and host is 'gateway.default.svc.cluster.local' and port is '8080'",
			args: args{
				addr: "",
				host: "gateway.default.svc.cluster.local",
				port: 8080,
			},
			want: want{
				want: New("invalid gRPC client connection port to addr: ,\thost: gateway.default.svc.cluster.local\t port: 8080"),
			},
		},
		{
			name: "return ErrInvalidGRPCPort error when addr is '127.0.0.1' and host is empty and port is '8080'",
			args: args{
				addr: "127.0.0.1",
				host: "",
				port: 8080,
			},
			want: want{
				want: New("invalid gRPC client connection port to addr: 127.0.0.1,\thost: \t port: 8080"),
			},
		},
		{
			name: "return ErrInvalidGRPCPort error when addr is '127.0.0.1' and host is 'gateway.default.svc.cluster.local' and port is '0'",
			args: args{
				addr: "127.0.0.1",
				host: "gateway.default.svc.cluster.local",
				port: 0,
			},
			want: want{
				want: New("invalid gRPC client connection port to addr: 127.0.0.1,\thost: gateway.default.svc.cluster.local\t port: 0"),
			},
		},
		{
			name: "return ErrInvalidGRPCPort error when addr is '127.0.0.1' and host is 'gateway.default.svc.cluster.local' and port is '1'",
			args: args{
				addr: "127.0.0.1",
				host: "gateway.default.svc.cluster.local",
				port: 1,
			},
			want: want{
				want: New("invalid gRPC client connection port to addr: 127.0.0.1,\thost: gateway.default.svc.cluster.local\t port: 1"),
			},
		},
		{
			name: "return ErrInvalidGRPCPort error when addr is '127.0.0.1' and host is 'gateway.default.svc.cluster.local' and port is maximum value of uint16",
			args: args{
				addr: "127.0.0.1",
				host: "gateway.default.svc.cluster.local",
				port: math.MaxUint16,
			},
			want: want{
				want: Errorf("invalid gRPC client connection port to addr: 127.0.0.1,\thost: gateway.default.svc.cluster.local\t port: %d", math.MaxUint16),
			},
		},
		{
			name: "return ErrInvalidGRPCPort error when addr is '127.0.0.1' and host is 'gateway.default.svc.cluster.local' and port is 'MaxUint16-1'",
			args: args{
				addr: "127.0.0.1",
				host: "gateway.default.svc.cluster.local",
				port: math.MaxUint16 - 1,
			},
			want: want{
				want: Errorf("invalid gRPC client connection port to addr: 127.0.0.1,\thost: gateway.default.svc.cluster.local\t port: %d", math.MaxUint16-1),
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

			got := ErrInvalidGRPCPort(test.args.addr, test.args.host, test.args.port)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidGRPCClientConn(t *testing.T) {
	type args struct {
		addr string
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
			name: "return ErrInvalidGRPCClientConn error when addr is '127.0.0.1'",
			args: args{
				addr: "127.0.0.1",
			},
			want: want{
				want: New("invalid gRPC client connection to 127.0.0.1"),
			},
		},
		{
			name: "return ErrInvalidGRPCClientConn error when addr is empty",
			args: args{
				addr: "",
			},
			want: want{
				want: New("invalid gRPC client connection to "),
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

			got := ErrInvalidGRPCClientConn(test.args.addr)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrGRPCLookupIPAddrNotFound(t *testing.T) {
	type args struct {
		host string
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
			name: "return ErrGRPCLookupIPAddrNotFound error when host is 'gateway.vald.svc.cluster.local'",
			args: args{
				host: "gateway.vald.svc.cluster.local",
			},
			want: want{
				want: New("vald internal gRPC client could not find ip addrs for gateway.vald.svc.cluster.local"),
			},
		},
		{
			name: "return ErrGRPCLookupIPAddrNotFound error when host is empty",
			args: args{
				host: "",
			},
			want: want{
				want: New("vald internal gRPC client could not find ip addrs for "),
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

			got := ErrGRPCLookupIPAddrNotFound(test.args.host)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrGRPCClientNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
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
			name: "return ErrGRPCLookupIPAddrNotFound error",
			want: want{
				want: New("vald internal gRPC client not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrGRPCClientNotFound
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrGRPCClientConnNotFound(t *testing.T) {
	type args struct {
		addr string
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
			name: "return ErrGRPCClientConnNotFound error when addr is '127.0.0.1'",
			args: args{
				addr: "127.0.0.1",
			},
			want: want{
				want: New("gRPC client connection not found in 127.0.0.1"),
			},
		},
		{
			name: "return ErrGRPCClientConnNotFound error when addr is empty",
			args: args{
				addr: "",
			},
			want: want{
				want: New("gRPC client connection not found in "),
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

			got := ErrGRPCClientConnNotFound(test.args.addr)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRPCCallFailed(t *testing.T) {
	type args struct {
		addr string
		err  error
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
			name: "return wrapped ErrRPCCallFailed error when err is server error and addr is '127.0.0.1'",
			args: args{
				err:  New("server error"),
				addr: "127.0.0.1",
			},
			want: want{
				want: New("addr: 127.0.0.1: server error"),
			},
		},
		{
			name: "return wrapped ErrRPCCallFailed error when err is server error and addr is empty",
			args: args{
				err:  New("server error"),
				addr: "",
			},
			want: want{
				want: New("addr: : server error"),
			},
		},
		{
			name: "return ErrRPCCallFailed error when err is nil error and addr is '127.0.0.1'",
			args: args{
				err:  nil,
				addr: "127.0.0.1",
			},
			want: want{
				want: New("addr: 127.0.0.1"),
			},
		},
		{
			name: "return ErrRPCCallFailed error when err is nil error and addr is empty",
			args: args{
				err:  nil,
				addr: "",
			},
			want: want{
				want: New("addr: "),
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

			got := ErrRPCCallFailed(test.args.addr, test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrGRPCTargetAddrNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
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
			name: "return ErrGRPCTargetAddrNotFound error",
			want: want{
				want: New("grpc connection target not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrGRPCTargetAddrNotFound
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
