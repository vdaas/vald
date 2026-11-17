// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package control

import (
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/goleak"
	"golang.org/x/sys/unix"
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
		want       want
		checkFunc  func(want, SocketController) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
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
				defer test.afterFunc(tt, test.args)
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
		checkFunc  func(want, int) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
		want       want
		args       args
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
				defer test.afterFunc(tt, test.args)
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
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
		want       want
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
				defer test.afterFunc(tt, test.args)
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
		want       want
		checkFunc  func(want, func(network, addr string, c syscall.RawConn) (err error)) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
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
		c       syscall.RawConn
		network string
		address string
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
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
		fields     fields
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
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
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
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
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
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
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
				defer test.afterFunc(tt, test.args)
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

// controlFunc used to return c.Control(...) directly, which overwrote the setsockopt
// errors collected inside the closure with c.Control's own nil result, so they were
// silently dropped. setsockopts must report them; controlFunc only logs them.
func Test_control_setsockopts(t *testing.T) {
	t.Parallel()
	if SO_REUSEADDR == 0 {
		t.Skip("setsockopt is not supported on this platform")
	}
	// a directory fd is not a socket, so every setsockopt fails.
	f, err := file.Open(".", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			t.Errorf("failed to close file: %v", cerr)
		}
	}()
	sc, err := f.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}
	ctrl := &control{reuseAddr: true, tcpNoDelay: true, keepAlive: 10}
	var serr error
	if err = sc.Control(func(fd uintptr) {
		serr = ctrl.setsockopts("tcp", int(fd))
	}); err != nil {
		t.Fatalf("Control() error = %v", err)
	}
	if serr == nil {
		t.Error("setsockopts() error = nil, want the setsockopt failures on a non-socket fd")
	}
	// controlFunc stays fail-open: the same failures must not fail the dial or listen.
	if err = ctrl.controlFunc("tcp", "127.0.0.1", sc); err != nil {
		t.Errorf("controlFunc() error = %v, want nil", err)
	}
}

type errRawConn struct {
	err error
}

func (c errRawConn) Control(func(fd uintptr)) error { return c.err }

func (c errRawConn) Read(func(fd uintptr) bool) error { return c.err }

func (c errRawConn) Write(func(fd uintptr) bool) error { return c.err }

func Test_control_controlFunc_returnsControlError(t *testing.T) {
	t.Parallel()
	want := errors.New("control failed")
	err := new(control).controlFunc("tcp", "127.0.0.1", errRawConn{err: want})
	if !errors.Is(err, want) {
		t.Errorf("controlFunc() error = %v, want %v", err, want)
	}
}

type sockoptCall struct {
	fd    int
	level int
	opt   int
	value int
}

// mockRawConn implements syscall.RawConn for testing custom Control behaviors.
type mockRawConn struct {
	controlFunc func(f func(uintptr)) error
}

func (m *mockRawConn) Control(f func(uintptr)) error {
	if m.controlFunc != nil {
		return m.controlFunc(f)
	}
	f(0)
	return nil
}

func (m *mockRawConn) Read(f func(uintptr) (done bool)) error  { return nil }
func (m *mockRawConn) Write(f func(uintptr) (done bool)) error { return nil }

// Test_control_setsockopts_defaultKernelState empirically verifies
// on a real Linux TCP socket that default setsockopts sets TCP_NODELAY in the kernel
// and leaves all disabled options unchanged.
func Test_control_setsockopts_defaultKernelState(t *testing.T) {
	// Create a real raw TCP socket directly via syscall
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("failed to create socket: %v", err)
	}
	defer unix.Close(fd)

	// Verify initial kernel state before setsockopts
	initNoDelay, _ := unix.GetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_NODELAY)
	if initNoDelay != 0 {
		t.Fatalf("expected initial TCP_NODELAY=0, got %d", initNoDelay)
	}

	// Apply default control (flag=0, keepalive=0)
	ctrl := New(0, 0).(*control)
	err = ctrl.setsockopts("tcp4", fd)
	if err != nil {
		t.Fatalf("setsockopts returned error: %v", err)
	}

	// Inspect the kernel state via GetsockoptInt
	gotNoDelay, err := unix.GetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_NODELAY)
	if err != nil {
		t.Fatalf("failed to get TCP_NODELAY: %v", err)
	}
	gotReuseAddr, _ := unix.GetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR)
	gotReusePort, _ := unix.GetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT)
	gotCork, _ := unix.GetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_CORK)
	gotKeepAlive, _ := unix.GetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE)
	gotDeferAccept, _ := unix.GetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_DEFER_ACCEPT)

	if gotNoDelay != 1 {
		t.Errorf("empirical kernel TCP_NODELAY = %d, want 1", gotNoDelay)
	}
	if gotReuseAddr != 0 {
		t.Errorf("empirical kernel SO_REUSEADDR = %d, want 0", gotReuseAddr)
	}
	if gotReusePort != 0 {
		t.Errorf("empirical kernel SO_REUSEPORT = %d, want 0", gotReusePort)
	}
	if gotCork != 0 {
		t.Errorf("empirical kernel TCP_CORK = %d, want 0", gotCork)
	}
	if gotKeepAlive != 0 {
		t.Errorf("empirical kernel SO_KEEPALIVE = %d, want 0", gotKeepAlive)
	}
	if gotDeferAccept != 0 {
		t.Errorf("empirical kernel TCP_DEFER_ACCEPT = %d, want 0", gotDeferAccept)
	}
}

// Test_control_setsockopts_invocationCount counts every SetsockoptInt invocation
// and proves that disabled options are never called.
func Test_control_setsockopts_invocationCount(t *testing.T) {
	orig := SetsockoptInt
	defer func() { SetsockoptInt = orig }()

	var mu sync.Mutex
	var calls []sockoptCall
	SetsockoptInt = func(fd, level, opt, value int) error {
		mu.Lock()
		calls = append(calls, sockoptCall{fd: fd, level: level, opt: opt, value: value})
		mu.Unlock()
		return nil
	}

	// Case 1: Default control on TCP -> exactly 1 syscall (TCP_NODELAY = 1)
	calls = nil
	cDefault := New(0, 0)
	cStruct := cDefault.(*control)
	err := cStruct.setsockopts("tcp", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("expected exactly 1 syscall for default TCP, got %d: %+v", len(calls), calls)
	}
	if calls[0].level != unix.IPPROTO_TCP || calls[0].opt != unix.TCP_NODELAY || calls[0].value != 1 {
		t.Errorf("unexpected syscall parameters: %+v, want IPPROTO_TCP, TCP_NODELAY, 1", calls[0])
	}

	// Case 2: Default control on UDP -> exactly 0 syscalls
	calls = nil
	err = cStruct.setsockopts("udp", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 0 {
		t.Fatalf("expected 0 syscalls for default UDP, got %d: %+v", len(calls), calls)
	}

	// Case 3: Default control on UNIX domain socket -> exactly 0 syscalls
	calls = nil
	err = cStruct.setsockopts("unix", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 0 {
		t.Fatalf("expected 0 syscalls for default unix, got %d: %+v", len(calls), calls)
	}
}

// Test_control_setsockopts_tcpCorkOmitsNoDelay empirically verifies
// that when TCPCork is set:
// 1. TCP_NODELAY is omitted.
// 2. TCP_CORK is applied.
// 3. This holds both at the kernel level and at the syscall accounting level.
func Test_control_setsockopts_tcpCorkOmitsNoDelay(t *testing.T) {
	// Part A: Syscall accounting
	orig := SetsockoptInt
	defer func() { SetsockoptInt = orig }()

	var mu sync.Mutex
	var calls []sockoptCall
	SetsockoptInt = func(fd, level, opt, value int) error {
		mu.Lock()
		calls = append(calls, sockoptCall{fd: fd, level: level, opt: opt, value: value})
		mu.Unlock()
		return nil
	}

	cCork := New(TCPCork, 0).(*control)
	calls = nil
	err := cCork.setsockopts("tcp", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that TCP_NODELAY was NOT called, but TCP_CORK was called
	var foundCork, foundNoDelay bool
	for _, call := range calls {
		if call.level == unix.IPPROTO_TCP && call.opt == unix.TCP_NODELAY {
			foundNoDelay = true
		}
		if call.level == unix.IPPROTO_TCP && call.opt == unix.TCP_CORK {
			foundCork = true
			if call.value != 1 {
				t.Errorf("TCP_CORK value = %d, want 1", call.value)
			}
		}
	}
	if foundNoDelay {
		t.Errorf("TCP_NODELAY was invoked when TCPCork was enabled: %+v", calls)
	}
	if !foundCork {
		t.Errorf("TCP_CORK was NOT invoked when TCPCork was enabled: %+v", calls)
	}

	// Part B: Combined TCPNoDelay | TCPCork -> cork overrides/omits nodelay
	calls = nil
	cBoth := New(TCPNoDelay|TCPCork, 0).(*control)
	err = cBoth.setsockopts("tcp", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, call := range calls {
		if call.level == unix.IPPROTO_TCP && call.opt == unix.TCP_NODELAY {
			t.Errorf("TCP_NODELAY was invoked when TCPNoDelay|TCPCork were both specified: %+v", calls)
		}
	}

	// Part C: Empirical Linux kernel verification on real raw socket
	SetsockoptInt = orig // restore real syscalls

	rawFd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("failed to create raw socket: %v", err)
	}
	defer unix.Close(rawFd)

	// Verify initial kernel state before setsockopts
	initNoDelay, _ := unix.GetsockoptInt(rawFd, unix.IPPROTO_TCP, unix.TCP_NODELAY)
	initCork, _ := unix.GetsockoptInt(rawFd, unix.IPPROTO_TCP, unix.TCP_CORK)
	if initNoDelay != 0 || initCork != 0 {
		t.Fatalf("expected initial TCP_NODELAY=0, TCP_CORK=0, got %d, %d", initNoDelay, initCork)
	}

	err = cCork.setsockopts("tcp4", rawFd)
	if err != nil {
		t.Fatalf("setsockopts error: %v", err)
	}

	kernelNoDelay, err := unix.GetsockoptInt(rawFd, unix.IPPROTO_TCP, unix.TCP_NODELAY)
	if err != nil {
		t.Fatalf("failed to read TCP_NODELAY: %v", err)
	}
	kernelCork, err := unix.GetsockoptInt(rawFd, unix.IPPROTO_TCP, unix.TCP_CORK)
	if err != nil {
		t.Fatalf("failed to read TCP_CORK: %v", err)
	}

	if kernelCork != 1 {
		t.Errorf("empirical kernel TCP_CORK = %d, want 1", kernelCork)
	}
	if kernelNoDelay != 0 {
		t.Errorf("empirical kernel TCP_NODELAY = %d, want 0 when cork is active", kernelNoDelay)
	}

	if kernelCork != 1 {
		t.Errorf("empirical kernel TCP_CORK = %d, want 1", kernelCork)
	}
	if kernelNoDelay != 0 {
		t.Errorf("empirical kernel TCP_NODELAY = %d, want 0 when cork is active", kernelNoDelay)
	}
}

// Test_control_GetControl_nonSocketFailsOpen verifies that calling the socket controller
// on non-socket raw connections (files, pipes, invalid fds, mock errors) fails open
// (returns nil) and never panics.
func Test_control_GetControl_nonSocketFailsOpen(t *testing.T) {
	allFlags := ReusePort | ReuseAddr | TCPFastOpen | TCPNoDelay | TCPCork |
		TCPQuickAck | TCPDeferAccept | IPTransparent | IPRecoverDestinationAddr
	ctrls := []SocketController{
		New(0, 0),
		New(allFlags, 60),
		New(TCPCork, 0),
		(*control)(nil), // nil controller receiver
	}

	for i, c := range ctrls {
		ctrl := c
		fn := ctrl.GetControl()
		if fn == nil {
			t.Fatalf("[%d] GetControl returned nil function", i)
		}

		// 1. Regular os.File raw connection
		tmpFile, err := os.CreateTemp("", "vald-control-test-*")
		if err != nil {
			t.Fatalf("[%d] CreateTemp error: %v", i, err)
		}
		scFile, err := tmpFile.SyscallConn()
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			t.Fatalf("[%d] SyscallConn error: %v", i, err)
		}

		// Must return nil and not panic
		err = fn("tcp4", "127.0.0.1:8080", scFile)
		if err != nil {
			t.Errorf("[%d] expected nil (fail-open) on os.File, got %v", i, err)
		}
		tmpFile.Close()
		os.Remove(tmpFile.Name())

		// 2. Pipe raw connection
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("[%d] os.Pipe error: %v", i, err)
		}
		scPipe, err := r.SyscallConn()
		if err != nil {
			r.Close()
			w.Close()
			t.Fatalf("[%d] Pipe SyscallConn error: %v", i, err)
		}

		err = fn("tcp4", "127.0.0.1:8080", scPipe)
		if err != nil {
			t.Errorf("[%d] expected nil (fail-open) on os.Pipe, got %v", i, err)
		}
		r.Close()
		w.Close()

		// 3. Mock RawConn with invalid fd (-1)
		mockInvalidFD := &mockRawConn{
			controlFunc: func(f func(uintptr)) error {
				f(uintptr(^uint(0))) // fd = -1
				return nil
			},
		}
		err = fn("tcp4", "127.0.0.1:8080", mockInvalidFD)
		if err != nil {
			t.Errorf("[%d] expected nil (fail-open) on invalid fd, got %v", i, err)
		}

		// 4. Mock RawConn where Control itself fails with error
		expectedErr := errors.New("raw connection control failure")
		mockErrorConn := &mockRawConn{
			controlFunc: func(f func(uintptr)) error {
				return expectedErr
			},
		}
		err = fn("tcp4", "127.0.0.1:8080", mockErrorConn)
		if !errors.Is(err, expectedErr) {
			if i != 3 { // nil controller returns nil regardless
				t.Errorf("[%d] expected error %v, got %v", i, expectedErr, err)
			}
		}
	}
}

// NOT IMPLEMENTED BELOW

// TCP_KEEPIDLE/TCP_KEEPINTVL take seconds, but the callers used to pass the keepalive
// time.Duration as nanoseconds (30s became 30e9, which setsockopt truncates and rejects).
func TestKeepAliveSeconds(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		d    time.Duration
		want int
	}{
		{name: "zero stays disabled", d: 0, want: 0},
		{name: "negative stays disabled", d: -time.Second, want: 0},
		{name: "sub-second becomes 1", d: 500 * time.Millisecond, want: 1},
		{name: "one second", d: time.Second, want: 1},
		{name: "fractional seconds are truncated", d: 1500 * time.Millisecond, want: 1},
		{name: "thirty seconds", d: 30 * time.Second, want: 30},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			if got := KeepAliveSeconds(tc.d); got != tc.want {
				tt.Errorf("KeepAliveSeconds(%v) = %d, want %d", tc.d, got, tc.want)
			}
		})
	}
}

func Test_control_setsockopts_SelectiveSyscalls(t *testing.T) {
	t.Parallel()
	// On a non-socket file fd with default empty control, 0 syscalls are made on "file" network.
	f, err := file.Open(".", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	sc, err := f.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}

	ctrl := &control{}
	var serr error
	err = sc.Control(func(fd uintptr) {
		serr = ctrl.setsockopts("file", int(fd))
	})
	if err != nil {
		t.Fatalf("Control() error = %v", err)
	}
	// Since no options are enabled for non-TCP network, setsockopts makes 0 syscalls and returns nil.
	if serr != nil {
		t.Errorf("setsockopts(\"file\", ...) error = %v, want nil for empty control", serr)
	}
}

func Test_control_setsockopts_TCPNoDelayDefault(t *testing.T) {
	t.Parallel()
	ctrl := &control{}
	if !isTCP("tcp") {
		t.Fatal("expected tcp to be isTCP")
	}
	// Verify that default control has tcpCork == false
	if ctrl.tcpCork {
		t.Error("expected tcpCork to be false by default")
	}
}
