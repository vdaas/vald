package errors

import "testing"

func TestErrgRPCClientConnectionClose(t *testing.T) {
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
			name: "return wrapped ErrgRPCClientConnectionClose error when err is server error and name is 'gateway'",
			args: args{
				err:  New("server error"),
				name: "gateway",
			},
			want: want{
				want: New("gateway's gRPC connection close error: server error"),
			},
		},
		{
			name: "return wrapped ErrgRPCClientConnectionClose error when err is server error and name is empty",
			args: args{
				err:  New("server error"),
				name: "",
			},
			want: want{
				want: New("'s gRPC connection close error: server error"),
			},
		},
		{
			name: "return ErrgRPCClientConnectionClose error when err is nil error and name is 'gateway'",
			args: args{
				err:  nil,
				name: "gateway",
			},
			want: want{
				want: New("gateway's gRPC connection close error"),
			},
		},
		{
			name: "return ErrgRPCClientConnectionClose error when err is nil error and addr is empty",
			args: args{
				err:  nil,
				name: "",
			},
			want: want{
				want: New("'s gRPC connection close error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrgRPCClientConnectionClose(test.args.name, test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrInvalidGRPCClientConn(test.args.addr)
			if err := test.checkFunc(test.want, got); err != nil {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrGRPCLookupIPAddrNotFound(test.args.host)
			if err := test.checkFunc(test.want, got); err != nil {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrGRPCClientNotFound
			if err := test.checkFunc(test.want, got); err != nil {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrGRPCClientConnNotFound(test.args.addr)
			if err := test.checkFunc(test.want, got); err != nil {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrRPCCallFailed(test.args.addr, test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrGRPCTargetAddrNotFound
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
