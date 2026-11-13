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

package net

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
)

// bracketedIPv6 is an IPv6 literal in the bracketed host form, without a port.
const bracketedIPv6 = "[2001:db8::1]"

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	goleak.VerifyTestMain(m)
}

func TestIsLocal(t *testing.T) {
	t.Parallel()
	type args struct {
		host string
	}
	type want struct {
		want bool
	}
	type test struct {
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
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
			name: "return true if it is host is `localhost`",
			args: args{
				host: "localhost",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true if it is host is local IPv4 address",
			args: args{
				host: "127.0.0.1",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true if it is host is local IPv6 address",
			args: args{
				host: "::1",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false if the host is not an address",
			args: args{
				host: "dummy",
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false if the host is empty",
			args: args{
				host: "",
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

			got := IsLocal(test.args.host)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDialContext(t *testing.T) {
	t.Parallel()
	type args struct {
		network string
		addr    string
	}
	type want struct {
		wantConn Conn
		err      error
	}
	type test struct {
		want       want
		srv        *httptest.Server
		checkFunc  func(want, Conn, error) error
		beforeFunc func(*testing.T, *test)
		afterFunc  func(*testing.T, *test)
		args       args
		name       string
	}
	defaultCheckFunc := func(w want, gotConn Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotConn, w.wantConn) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
		}
		return nil
	}
	tests := []test{
		{
			name: "dial return server content",
			args: args{
				network: TCP.String(),
			},
			beforeFunc: func(t *testing.T, test *test) {
				t.Helper()
				srvContent := "test"

				test.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, srvContent)
				}))
				test.args.addr = test.srv.URL[len("http://"):]
			},
			checkFunc: func(w want, gotConn Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				srvContent := "test"

				// read the output from the server and check if it is equals to the count
				fmt.Fprintf(gotConn, "GET / HTTP/1.0\r\n\r\n")
				buf, err := io.ReadAll(gotConn)
				if err != nil || buf == nil {
					return errors.Errorf("error or buffer is nil,\terror: %v, buf: %v", err, buf)
				}
				content := strings.Split(conv.Btoa(buf), "\n")[5] // skip HTTP header
				if content != srvContent {
					return errors.Errorf("invalid content, got: %v, want: %v", content, srvContent)
				}

				return nil
			},
			afterFunc: func(t *testing.T, test *test) {
				t.Helper()
				test.srv.Client().CloseIdleConnections()
				test.srv.CloseClientConnections()
				test.srv.Close()
			},
		},
	}

	for i := range tests {
		test := &tests[i]
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx := tt.Context()
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotConn, err := DialContext(ctx, test.args.network, test.args.addr)
			if err := checkFunc(test.want, gotConn, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if test.afterFunc != nil {
				test.afterFunc(tt, test)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		addr string
	}
	type want struct {
		err      error
		wantHost string
		wantPort uint16
		isLocal  bool
		isV4     bool
		isV6     bool
	}
	type test struct {
		checkFunc  func(want, string, uint16, bool, bool, bool, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		args       args
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, gotHost string, gotPort uint16, gotIsLocal, gotIsV4, gotIsV6 bool, err error) error {
		if (w.err == nil && err != nil) || (w.err != nil && err == nil) || (err != nil && err.Error() != w.err.Error()) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotHost, w.wantHost) {
			return errors.Errorf("host got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHost, w.wantHost)
		}
		if !reflect.DeepEqual(gotPort, w.wantPort) {
			return errors.Errorf("port got: \"%d\",\n\t\t\t\twant: \"%d\"", gotPort, w.wantPort)
		}
		if !reflect.DeepEqual(gotIsLocal, w.isLocal) {
			return errors.Errorf("isLocal got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIsLocal, w.isLocal)
		}
		if !reflect.DeepEqual(gotIsV4, w.isV4) {
			return errors.Errorf("isV4 got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIsV4, w.isV4)
		}
		if !reflect.DeepEqual(gotIsV6, w.isV6) {
			return errors.Errorf("isV6 got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIsV6, w.isV6)
		}
		return nil
	}
	tests := []test{
		{
			name: "parse success with IPv4 address",
			args: args{
				addr: "192.168.1.1:8080",
			},
			want: want{
				wantHost: "192.168.1.1",
				wantPort: uint16(8080),
				isV4:     true,
				isV6:     false,
				isLocal:  true,
			},
		},
		{
			name: "parse success with IPv6 address",
			args: args{
				addr: "[2001:db8::1]:8080",
			},
			want: want{
				wantHost: "2001:db8::1",
				wantPort: uint16(8080),
				isV4:     false,
				isV6:     true,
				isLocal:  false,
			},
		},
		{
			name: "parse success with bracketed IPv6 address without port",
			args: args{
				addr: bracketedIPv6,
			},
			want: want{
				wantHost: "2001:db8::1",
				wantPort: uint16(80),
				isV4:     false,
				isV6:     true,
				isLocal:  false,
			},
		},
		{
			name: "return true if it is local address",
			args: args{
				addr: "localhost:8080",
			},
			want: want{
				wantHost: "localhost",
				wantPort: uint16(8080),
				isV4:     false,
				isV6:     true,
				isLocal:  true,
			},
		},
		{
			name: "parse success with hostname",
			args: args{
				addr: "google.com:8080",
			},
			want: want{
				wantHost: "google.com",
				wantPort: uint16(8080),
				isLocal:  false,
			},
			checkFunc: func(w want, gotHost string, gotPort uint16, gotIsLocal bool, gotIsV4 bool, gotIsV6 bool, err error) error {
				if gotHost != w.wantHost {
					return errors.Errorf("gotHost: %q, want: %q", gotHost, w.wantHost)
				}
				if gotPort != w.wantPort {
					return errors.Errorf("gotPort: %d, want: %d", gotPort, w.wantPort)
				}
				if gotIsLocal != w.isLocal {
					return errors.Errorf("gotIsLocal: %v, want: %v", gotIsLocal, w.isLocal)
				}
				return nil
			},
		},
		{
			name: "return default port when parse failed if it is not an address",
			args: args{
				addr: "dummy",
			},
			want: want{
				wantHost: "dummy",
				wantPort: uint16(80),
			},
			checkFunc: func(w want, s string, u uint16, b1, b2, b3 bool, err error) error {
				if strings.Contains(err.Error(), "failed to lookup ip for dummy") {
					return nil
				}
				return err
			},
		},
		{
			name: "return default port when parse failed if port number missing in IPv4 address",
			args: args{
				addr: "192.168.1.1",
			},
			want: want{
				wantHost: "192.168.1.1",
				wantPort: uint16(80),
				isV4:     true,
				isV6:     false,
				isLocal:  true,
			},
		},
		{
			name: "return default port when parse failed if port number missing in IPv6 address",
			args: args{
				addr: "2001:db8::1",
			},
			want: want{
				wantHost: "2001:db8::1",
				wantPort: uint16(80),
				isV4:     false,
				isV6:     true,
				isLocal:  false,
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

			gotHost, gotPort, gotIsLocal, gotIsV4, gotIsV6, err := Parse(tt.Context(), test.args.addr)
			if err := checkFunc(test.want, gotHost, gotPort, gotIsLocal, gotIsV4, gotIsV6, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSplitHostPort(t *testing.T) {
	t.Parallel()
	type args struct {
		hostport string
	}
	type want struct {
		err      error
		wantHost string
		wantPort uint16
	}
	type test struct {
		checkFunc  func(want, string, uint16, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		want       want
	}
	defaultCheckFunc := func(w want, gotHost string, gotPort uint16, err error) error {
		if (w.err == nil && err != nil) || (w.err != nil && err == nil) || (err != nil && err.Error() != w.err.Error()) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotHost, w.wantHost) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHost, w.wantHost)
		}
		if !reflect.DeepEqual(gotPort, w.wantPort) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPort, w.wantPort)
		}
		return nil
	}
	tests := []test{
		{
			name: "parse success with IPv4 address",
			args: args{
				hostport: "192.168.1.1:8080",
			},
			want: want{
				wantHost: "192.168.1.1",
				wantPort: uint16(8080),
			},
		},
		{
			name: "parse success with IPv6 address",
			args: args{
				hostport: "[2001:db8::1]:8080",
			},
			want: want{
				wantHost: "2001:db8::1",
				wantPort: uint16(8080),
			},
		},
		{
			name: "parse success with hostname",
			args: args{
				hostport: "google.com:8080",
			},
			want: want{
				wantHost: "google.com",
				wantPort: uint16(8080),
			},
		},
		{
			name: "return default port when parse failed if it is not an address",
			args: args{
				hostport: "dummy",
			},
			want: want{
				wantHost: "dummy",
				wantPort: uint16(80),
				err: &net.AddrError{
					Addr: "dummy",
					Err:  "missing port in address",
				},
			},
		},
		{
			name: "return default port when parse failed if port number missing in IPv4 address",
			args: args{
				hostport: "192.168.1.1",
			},
			want: want{
				wantHost: "192.168.1.1",
				wantPort: uint16(80),
				err: &net.AddrError{
					Addr: "192.168.1.1",
					Err:  "missing port in address",
				},
			},
		},
		{
			name: "return default port when parse failed if port number missing in IPv6 address",
			args: args{
				hostport: "2001:db8::1",
			},
			want: want{
				wantHost: "2001:db8::1",
				wantPort: uint16(80),
				err: &net.AddrError{
					Addr: "2001:db8::1",
					Err:  "too many colons in address",
				},
			},
		},
		{
			name: "return the unbracketed host and default port when port number missing in bracketed IPv6 address",
			args: args{
				hostport: bracketedIPv6,
			},
			want: want{
				wantHost: "2001:db8::1",
				wantPort: uint16(80),
				err: &net.AddrError{
					Addr: bracketedIPv6,
					Err:  "missing port in address",
				},
			},
		},
		{
			name: "parse success with default IPv4 address",
			args: args{
				hostport: ":8080",
			},
			want: want{
				wantHost: "127.0.0.1",
				wantPort: uint16(8080),
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

			gotHost, gotPort, err := SplitHostPort(test.args.hostport)
			if err := checkFunc(test.want, gotHost, gotPort, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestScanPorts(t *testing.T) {
	t.Parallel()
	type args struct {
		host  string
		start uint16
		end   uint16
	}
	type want struct {
		err       error
		wantPorts []uint16
	}
	type test struct {
		want       want
		checkFunc  func(want, []uint16, error) error
		beforeFunc func(*testing.T, *test)
		afterFunc  func(*test)
		args       args
		name       string
		srvs       []*httptest.Server
	}
	defaultAfterFunc := func(t *test) {
		for _, s := range t.srvs {
			s.Client().CloseIdleConnections()
			s.CloseClientConnections()
			s.Close()
		}
	}
	defaultCheckFunc := func(w want, gotPorts []uint16, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		// count how want ports exists in got ports
		cnt := 0
		for _, wp := range w.wantPorts {
			if slices.Contains(gotPorts, wp) {
				cnt++
			}
		}

		if cnt != len(w.wantPorts) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPorts, w.wantPorts)
		}
		return nil
	}
	tests := []test{
		{
			name: "return test server port number in given range",
			args: args{
				host: "localhost",
			},
			beforeFunc: func(t *testing.T, test *test) {
				t.Helper()
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))

				_, p, err := SplitHostPort(srv.URL[len("http://"):])
				if err != nil {
					t.Error(err)
				}

				test.srvs = []*httptest.Server{
					srv,
				}

				// set args to server port range
				test.args.start = p - 5
				test.args.end = p + 5

				// set want to server port
				test.want.wantPorts = []uint16{
					p,
				}
			},
		},
		{
			name: "return test server port number when start = end",
			args: args{
				host: "localhost",
			},
			beforeFunc: func(t *testing.T, test *test) {
				t.Helper()
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))

				_, p, err := SplitHostPort(srv.URL[len("http://"):])
				if err != nil {
					t.Error(err)
				}

				test.srvs = []*httptest.Server{
					srv,
				}

				// set args to server port range
				test.args.start = p
				test.args.end = p

				// set want to server port
				test.want.wantPorts = []uint16{
					p,
				}
			},
		},
		{
			name: "return test server port number when start > end",
			args: args{
				host: "localhost",
			},
			beforeFunc: func(t *testing.T, test *test) {
				t.Helper()
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))

				_, p, err := SplitHostPort(srv.URL[len("http://"):])
				if err != nil {
					t.Error(err)
				}

				test.srvs = []*httptest.Server{
					srv,
				}

				// set args to server port range
				test.args.start = p + 10
				test.args.end = p - 10

				// set want to server port
				test.want.wantPorts = []uint16{
					p,
				}
			},
		},
		{
			name: "return multiple test server port number",
			args: args{
				host: "localhost",
			},
			beforeFunc: func(t *testing.T, test *test) {
				t.Helper()
				srvNum := 20

				srvs := make([]*httptest.Server, 0, srvNum)
				ports := make([]uint16, 0, srvNum)
				minPort := uint16(math.MaxUint16)
				maxPort := uint16(0)

				for range srvNum {
					handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
					})
					srv := httptest.NewServer(handler)
					srvs = append(srvs, srv)

					_, p, err := SplitHostPort(srv.URL[len("http://"):])
					if err != nil {
						t.Error(err)
					}

					ports = append(ports, p)

					if p < minPort {
						minPort = p
					}
					if p > maxPort {
						maxPort = p
					}
				}

				test.srvs = srvs

				// set args to server port range
				test.args.start = minPort - 10
				test.args.end = maxPort + 10

				// set want to server port
				test.want.wantPorts = ports
			},
		},
		{
			name: "return no port available if no port is scanned",
			args: args{
				host:  "localhost",
				start: 65534,
				end:   65535,
			},
			want: want{
				err: errors.ErrNoPortAvailable("localhost", 65534, 65535),
			},
		},
	}

	for i := range tests {
		test := &tests[i]
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx := tt.Context()
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			afterFunc := test.afterFunc
			if test.afterFunc == nil {
				afterFunc = defaultAfterFunc
			}
			defer afterFunc(test)

			gotPorts, err := ScanPorts(ctx, test.args.start, test.args.end, test.args.host)
			if err := checkFunc(test.want, gotPorts, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			afterFunc(test)
		})
	}
}

func TestLoadLocalIP(t *testing.T) {
	t.Parallel()
	type want struct {
		want string
	}
	type test struct {
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "get local ip",
			want: want{
				want: "127.0.0.1",
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

			got := LoadLocalIP()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNetworkTypeFromString(t *testing.T) {
	t.Parallel()
	type args struct {
		str string
	}
	type want struct {
		want NetworkType
	}
	type test struct {
		checkFunc  func(want, NetworkType) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		want       want
	}
	defaultCheckFunc := func(w want, got NetworkType) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return UNIX when the string is UNIX",
			args: args{
				str: "UNIX",
			},
			want: want{
				want: UNIX,
			},
		},
		{
			name: "return UNIXGRAM when the string is UNIXGRAM",
			args: args{
				str: "UNIXGRAM",
			},
			want: want{
				want: UNIXGRAM,
			},
		},
		{
			name: "return UNIXPACKET when the string is UNIXPACKET",
			args: args{
				str: "UNIXPACKET",
			},
			want: want{
				want: UNIXPACKET,
			},
		},
		{
			name: "return ICMP when the string is ICMP",
			args: args{
				str: "ICMP",
			},
			want: want{
				want: ICMP,
			},
		},
		{
			name: "return ICMP6 when the string is ipv6-icmp",
			args: args{
				str: "ipv6-icmp",
			},
			want: want{
				want: ICMP6,
			},
		},
		{
			name: "return IGMP when the string is IGMP",
			args: args{
				str: "IGMP",
			},
			want: want{
				want: IGMP,
			},
		},
		{
			name: "return TCP when the string is TCP",
			args: args{
				str: "TCP",
			},
			want: want{
				want: TCP,
			},
		},
		{
			name: "return TCP4 when the string is TCP4",
			args: args{
				str: "TCP4",
			},
			want: want{
				want: TCP4,
			},
		},
		{
			name: "return TCP6 when the string is TCP6",
			args: args{
				str: "TCP6",
			},
			want: want{
				want: TCP6,
			},
		},
		{
			name: "return UDP when the string is UDP",
			args: args{
				str: "UDP",
			},
			want: want{
				want: UDP,
			},
		},
		{
			name: "return UDP4 when the string is UDP4",
			args: args{
				str: "UDP4",
			},
			want: want{
				want: UDP4,
			},
		},
		{
			name: "return UDP6 when the string is UDP6",
			args: args{
				str: "UDP6",
			},
			want: want{
				want: UDP6,
			},
		},
		{
			name: "return UNKNOWN when the string is invalid string",
			args: args{
				str: "invalid type",
			},
			want: want{
				want: Unknown,
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

			got := NetworkTypeFromString(test.args.str)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNetworkType_String(t *testing.T) {
	t.Parallel()
	type want struct {
		want string
	}
	type test struct {
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
		name       string
		want       want
		n          NetworkType
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return unix when the type is UNIX",
			n:    UNIX,
			want: want{
				want: "unix",
			},
		},
		{
			name: "return unixgram when the type is UNIXGRAM",
			n:    UNIXGRAM,
			want: want{
				want: "unixgram",
			},
		},
		{
			name: "return unixpacket when the type is UNIXPACKET",
			n:    UNIXPACKET,
			want: want{
				want: "unixpacket",
			},
		},
		{
			name: "return tcp when the type is TCP",
			n:    TCP,
			want: want{
				want: "tcp",
			},
		},
		{
			name: "return tcp4 when the type is TCP4",
			n:    TCP4,
			want: want{
				want: "tcp4",
			},
		},
		{
			name: "return tcp6 when the type is TCP6",
			n:    TCP6,
			want: want{
				want: "tcp6",
			},
		},
		{
			name: "return udp when the type is UDP",
			n:    UDP,
			want: want{
				want: "udp",
			},
		},
		{
			name: "return udp4 when the type is UDP4",
			n:    UDP4,
			want: want{
				want: "udp4",
			},
		},
		{
			name: "return udp6 when the type is UDP6",
			n:    UDP6,
			want: want{
				want: "udp6",
			},
		},
		{
			name: "return icmp when the type is ICMP",
			n:    ICMP,
			want: want{
				want: "icmp",
			},
		},
		{
			name: "return igmp when the type is IGMP",
			n:    IGMP,
			want: want{
				want: "igmp",
			},
		},
		{
			name: "return ipv6-icmp when the type is ICMP6",
			n:    ICMP6,
			want: want{
				want: "ipv6-icmp",
			},
		},
		{
			name: "return unknown when the type is Unknown",
			n:    Unknown,
			want: want{
				want: "unknown",
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

			got := test.n.String()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestJoinHostPort(t *testing.T) {
	t.Parallel()
	type args struct {
		host string
		port uint16
	}
	type want struct {
		want string
	}
	type test struct {
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		want       want
		args       args
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ipv4 host port",
			args: args{
				host: "127.0.0.1",
				port: 8080,
			},
			want: want{
				want: "127.0.0.1:8080",
			},
		},
		{
			name: "return ipv6 host port",
			args: args{
				host: "2001:db8::1",
				port: 8081,
			},
			want: want{
				want: "[2001:db8::1]:8081",
			},
		},
		{
			name: "return hostname port",
			args: args{
				host: "www.example.com",
				port: 80,
			},
			want: want{
				want: "www.example.com:80",
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

			got := JoinHostPort(test.args.host, test.args.port)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDistinctAddrs(t *testing.T) {
	t.Parallel()
	// Named once and referenced by identifier below (rather than repeating the
	// literals) to avoid tripping goconst across the many test cases that need
	// to reuse the same addresses to exercise duplicate removal.
	const (
		addr1 = "10.0.0.1:8081"
		addr2 = "10.0.0.2:8081"
		addr3 = "10.0.0.3:8081"
	)
	type args struct {
		addrs []string
	}
	type want struct {
		want []string
	}
	type test struct {
		checkFunc  func(want, []string) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		want       want
	}
	defaultCheckFunc := func(w want, got []string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return empty slice when addrs is nil",
			args: args{
				addrs: nil,
			},
			want: want{
				want: []string{},
			},
		},
		{
			name: "return empty slice when addrs is empty",
			args: args{
				addrs: []string{},
			},
			want: want{
				want: []string{},
			},
		},
		{
			name: "return the same addrs when there are no duplicates",
			args: args{
				addrs: []string{addr1, addr2, addr3},
			},
			want: want{
				want: []string{addr1, addr2, addr3},
			},
		},
		{
			name: "remove duplicated addrs while preserving order of first occurrence",
			args: args{
				addrs: []string{addr2, addr1, addr2, addr3, addr1},
			},
			want: want{
				want: []string{addr2, addr1, addr3},
			},
		},
		{
			name: "return single element slice when all addrs are the same",
			args: args{
				addrs: []string{addr1, addr1, addr1},
			},
			want: want{
				want: []string{addr1},
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

			got := DistinctAddrs(test.args.addrs)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// TestIsIPLiteral_cornerCases checks IsIPLiteral
// against corner cases: leading/trailing whitespace, IPv4-in-IPv6,
// zone IDs (e.g. fe80::1%eth0), invalid formats, ports, and bracketed without ports.
func TestIsIPLiteral_cornerCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		addr string
		want bool
	}{
		// 1. Whitespace Corner Cases
		{name: "leading space IPv4", addr: " 192.168.1.1", want: false},
		{name: "trailing space IPv4", addr: "192.168.1.1 ", want: false},
		{name: "surrounding space IPv4", addr: " 192.168.1.1 ", want: false},
		{name: "tab padded IPv4", addr: "\t192.168.1.1\t", want: false},
		{name: "newline padded IPv4", addr: "\n192.168.1.1\n", want: false},
		{name: "leading space IPv6", addr: " ::1", want: false},
		{name: "trailing space IPv6", addr: "::1 ", want: false},
		{name: "leading space with port", addr: " 192.168.1.1:8080", want: false},
		// Note on trailing space with port: net.SplitHostPort splits before parsing host,
		// so host "192.168.1.1" parses as IP literal while port string contains the space.
		{name: "trailing space on port", addr: "192.168.1.1:8080 ", want: true},
		{name: "bracketed IPv6 with leading space", addr: " [::1]:8080", want: false},
		{name: "bracketed IPv6 with trailing space on port", addr: "[::1]:8080 ", want: true},

		// 2. IPv4-in-IPv6 / IPv4-mapped IPv6
		{name: "IPv4-mapped IPv6 standard", addr: "::ffff:192.0.2.1", want: true},
		{name: "IPv4-mapped IPv6 loopback", addr: "::ffff:127.0.0.1", want: true},
		{name: "IPv4-mapped IPv6 bracketed", addr: "[::ffff:192.0.2.1]", want: true},
		{name: "IPv4-mapped IPv6 bracketed with port", addr: "[::ffff:192.0.2.1]:8080", want: true},
		{name: "IPv4-mapped IPv6 uppercase", addr: "::FFFF:192.0.2.1", want: true},
		{name: "IPv4-compatible IPv6 (deprecated)", addr: "::192.0.2.1", want: true},

		// 3. Zone IDs (RFC 4007 / RFC 6874)
		{name: "IPv6 link-local with zone ID eth0", addr: "fe80::1%eth0", want: true},
		{name: "IPv6 link-local bracketed with zone ID eth0", addr: "[fe80::1%eth0]", want: true},
		{name: "IPv6 link-local bracketed with zone ID and port", addr: "[fe80::1%eth0]:8080", want: true},
		{name: "IPv6 link-local with numeric zone ID", addr: "fe80::1%1", want: true},
		{name: "IPv6 link-local bracketed with numeric zone ID", addr: "[fe80::1%1]", want: true},
		{name: "IPv6 link-local bracketed with numeric zone and port", addr: "[fe80::1%1]:443", want: true},
		{name: "IPv6 zone ID without brackets but with port", addr: "fe80::1%eth0:8080", want: true},

		// 4. Bracketed Literals without Ports
		{name: "bracketed IPv6 loopback", addr: "[::1]", want: true},
		{name: "bracketed IPv6 unspecified", addr: "[::]", want: true},
		{name: "bracketed full IPv6", addr: "[2001:db8::1]", want: true},
		{name: "bracketed IPv4", addr: "[192.168.1.1]", want: true},
		{name: "bracketed empty", addr: "[]", want: false},
		{name: "bracketed whitespace", addr: "[ ]", want: false},
		{name: "bracketed hostname", addr: "[localhost]", want: false},
		{name: "bracketed domain", addr: "[example.com]", want: false},
		{name: "unclosed opening bracket", addr: "[::1", want: false},
		{name: "unopened closing bracket", addr: "::1]", want: false},
		{name: "double bracketed", addr: "[[::1]]", want: false},

		// 5. Invalid Formats & Edge Cases
		{name: "IPv4 octet overflow", addr: "256.256.256.256", want: false},
		{name: "IPv4 5 octets", addr: "1.2.3.4.5", want: false},
		{name: "IPv4 3 octets", addr: "1.2.3", want: false},
		{name: "IPv4 trailing dot", addr: "192.168.1.1.", want: false},
		{name: "IPv6 triple colon", addr: ":::1", want: false},
		{name: "IPv6 invalid hex", addr: "2001:xyz::1", want: false},
		{name: "empty string", addr: "", want: false},
		{name: "single colon", addr: ":", want: false},
		{name: "double colon alone (unspecified IPv6)", addr: "::", want: true},
		{name: "just port", addr: ":8080", want: false},

		// 6. Hostnames with and without Ports (Must be false)
		{name: "hostname localhost", addr: "localhost", want: false},
		{name: "hostname localhost with port", addr: "localhost:8080", want: false},
		{name: "fully qualified domain", addr: "vald-agent.default.svc.cluster.local", want: false},
		{name: "fully qualified domain with port", addr: "vald-agent.default.svc.cluster.local:8081", want: false},
		{name: "hostname with numbers", addr: "host-1.local", want: false},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			got := IsIPLiteral(test.addr)
			if got != test.want {
				tt.Errorf("IsIPLiteral(%q) = %v, want %v", test.addr, got, test.want)
			}
		})
	}
}

// TestParse_robustness stress-tests Parse with
// missing port addresses, empty hostnames, malformed strings, and concurrent invocations.
func TestParse_robustness(t *testing.T) {
	t.Parallel()

	inputs := []struct {
		addr        string
		checkHost   string
		checkPort   uint16
		expectError bool
	}{
		// Missing port addresses
		{addr: "192.168.1.1", expectError: false, checkHost: "192.168.1.1", checkPort: defaultPort},
		{addr: "127.0.0.1", expectError: false, checkHost: "127.0.0.1", checkPort: defaultPort},
		{addr: "[::1]", expectError: false, checkHost: "::1", checkPort: defaultPort},
		{addr: "::1", expectError: false, checkHost: "::1", checkPort: defaultPort},
		{addr: "[2001:db8::1]", expectError: false, checkHost: "2001:db8::1", checkPort: defaultPort},

		// Empty hostnames
		{addr: ":8080", expectError: false, checkHost: "127.0.0.1", checkPort: 8080},
		{addr: ":", expectError: false, checkHost: "127.0.0.1", checkPort: defaultPort},
		{addr: ":0", expectError: false, checkHost: "127.0.0.1", checkPort: defaultPort},
		{addr: "", expectError: false, checkHost: "localhost", checkPort: defaultPort},

		// Malformed inputs (graceful error, no panic)
		{addr: "   ", expectError: true},
		{addr: "   :8080", expectError: true},
		{addr: "invalid_unresolvable_hostname_999999.local", expectError: true},
		{addr: "256.256.256.256:8080", expectError: true},
		{addr: "[::1:8080", expectError: true},
	}

	for _, tc := range inputs {
		test := tc
		t.Run("input_"+test.addr, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			host, port, _, _, _, err := Parse(ctx, test.addr)
			if test.expectError {
				if err == nil {
					tt.Errorf("Parse(%q) expected error, got nil (host=%q, port=%d)", test.addr, host, port)
				}
			} else {
				if err != nil {
					tt.Errorf("Parse(%q) unexpected error: %v", test.addr, err)
				}
				if test.checkHost != "" && host != test.checkHost {
					tt.Errorf("Parse(%q) host = %q, want %q", test.addr, host, test.checkHost)
				}
				if test.checkPort != 0 && port != test.checkPort {
					tt.Errorf("Parse(%q) port = %d, want %d", test.addr, port, test.checkPort)
				}
			}
		})
	}
}

// TestParse_concurrent verifies thread-safety, zero panics,
// and zero memory growth under concurrent execution of Parse.
func TestParse_concurrent(t *testing.T) {
	t.Parallel()

	const concurrency = 20
	const iterations = 500

	inputs := []string{
		"127.0.0.1",
		"127.0.0.1:8080",
		":8080",
		"[::1]",
		"[::1]:9090",
		":",
		"localhost:8080",
	}

	eg, _ := errgroup.New(context.Background())
	for i := range concurrency {
		eg.Go(func() error {
			for j := range iterations {
				addr := inputs[(i+j)%len(inputs)]
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				host, port, _, _, _, err := Parse(ctx, addr)
				cancel()
				if err != nil {
					return fmt.Errorf("unexpected Parse error for %q: %w", addr, err)
				}
				if host == "" || port == 0 {
					return fmt.Errorf("invalid Parse result for %q: host=%q port=%d", addr, host, port)
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("concurrent Parse stress failed: %v", err)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestIsUDP(t *testing.T) {
// 	type args struct {
// 		network string
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           network:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           network:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := IsUDP(test.args.network)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestIsTCP(t *testing.T) {
// 	type args struct {
// 		network string
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           network:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           network:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := IsTCP(test.args.network)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_lookupNetIP(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		rsv  *net.Resolver
// 		host string
// 	}
// 	type want struct {
// 		want []netip.Addr
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		want
// 		checkFunc  func(want, []netip.Addr, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []netip.Addr, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           rsv:net.Resolver{},
// 		           host:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           rsv:net.Resolver{},
// 		           host:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := lookupNetIP(test.args.ctx, test.args.rsv, test.args.host)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// TestScanPorts_contextDone verifies that ScanPorts terminates and reports the context error
// (not partial ports, not ErrNoPortAvailable) when the context is already done.
// The uint16 loop variable used to wrap to 0 after 65535, so start == 0 && end == 65535 never
// finished. The contexts are done up front so that no real large scan is performed.
func TestScanPorts_contextDone(t *testing.T) {
	t.Parallel()
	type test struct {
		ctxFunc func(*testing.T) (context.Context, context.CancelFunc)
		want    error
		name    string
		start   uint16
		end     uint16
		openSrv bool
	}
	cancelled := func(t *testing.T) (context.Context, context.CancelFunc) {
		t.Helper()
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		return ctx, cancel
	}
	expired := func(t *testing.T) (context.Context, context.CancelFunc) {
		t.Helper()
		return context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
	}
	tests := []test{
		{
			name:    "cancelled context on the full range returns context.Canceled",
			ctxFunc: cancelled,
			start:   0,
			end:     math.MaxUint16,
			want:    context.Canceled,
		},
		{
			name:    "expired deadline on the full range returns context.DeadlineExceeded",
			ctxFunc: expired,
			start:   0,
			end:     math.MaxUint16,
			want:    context.DeadlineExceeded,
		},
		{
			name:    "pre-cancelled context over a range with an open port returns context.Canceled and no ports",
			ctxFunc: cancelled,
			openSrv: true,
			want:    context.Canceled,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			start, end := test.start, test.end
			if test.openSrv {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				defer srv.Close()
				_, p, err := SplitHostPort(srv.URL[len("http://"):])
				if err != nil {
					tt.Fatal(err)
				}
				start, end = p, p
			}
			ctx, cancel := test.ctxFunc(tt)
			defer cancel()

			type result struct {
				err   error
				ports []uint16
			}
			done := make(chan result, 1)
			go func() {
				ports, err := ScanPorts(ctx, start, end, "localhost")
				done <- result{ports: ports, err: err}
			}()

			select {
			case r := <-done:
				if !errors.Is(r.err, test.want) {
					tt.Errorf("got error: %v, want: %v", r.err, test.want)
				}
				if len(r.ports) != 0 {
					tt.Errorf("unexpected ports returned for a done context: %v", r.ports)
				}
			case <-time.After(30 * time.Second):
				tt.Fatal("timed out waiting for ScanPorts to terminate")
			}
		})
	}
}

func TestIsIPLiteral(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		addr string
		want bool
	}{
		{name: "IPv4 literal", addr: "192.168.1.1", want: true},
		{name: "IPv4 loopback", addr: "127.0.0.1", want: true},
		{name: "IPv4 with port", addr: "192.168.1.1:8080", want: true},
		{name: "IPv6 loopback", addr: "::1", want: true},
		{name: "IPv6 literal", addr: "2001:db8::1", want: true},
		{name: "IPv6 bracketed", addr: "[::1]", want: true},
		{name: "IPv6 bracketed with port", addr: "[::1]:8080", want: true},
		{name: "IPv6 full bracketed with port", addr: "[2001:db8::1]:443", want: true},
		{name: "IPv4-mapped IPv6", addr: "::ffff:192.0.2.1", want: true},
		{name: "hostname localhost", addr: "localhost", want: false},
		{name: "hostname with port", addr: "localhost:8080", want: false},
		{name: "hostname domain", addr: "example.com", want: false},
		{name: "hostname domain with port", addr: "example.com:443", want: false},
		{name: "Kubernetes headless service", addr: "vald-agent-headless.default.svc.cluster.local:8081", want: false},
		{name: "empty string", addr: "", want: false},
		{name: "invalid IP", addr: "999.999.999.999", want: false},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			got := IsIPLiteral(test.addr)
			if got != test.want {
				tt.Errorf("IsIPLiteral(%q) = %v, want %v", test.addr, got, test.want)
			}
		})
	}
}

func Test_isMissingPortErr(t *testing.T) {
	t.Parallel()
	_, _, err := net.SplitHostPort("localhost")
	if !isMissingPortErr(err) {
		t.Errorf("isMissingPortErr(%v) = false, want true for missing port error", err)
	}
	if isMissingPortErr(nil) {
		t.Error("isMissingPortErr(nil) = true, want false")
	}
	if isMissingPortErr(errors.New("other error")) {
		t.Error("isMissingPortErr(other error) = true, want false")
	}
}
