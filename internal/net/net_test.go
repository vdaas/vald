//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package net provides net functionality for grpc
package net

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/goleak"
)

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
		name       string
		args       args
		want       want
		srv        *httptest.Server
		checkFunc  func(want, Conn, error) error
		beforeFunc func(*testing.T, *test)
		afterFunc  func(*testing.T, *test)
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
					w.WriteHeader(200)
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
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
		addr string
	}
	type want struct {
		wantHost string
		wantPort uint16
		isLocal  bool
		isV4     bool
		isV6     bool
		err      error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string, uint16, bool, bool, bool, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotHost string, gotPort uint16, gotIsLocal, gotIsV4, gotIsV6 bool, err error) error {
		if (w.err == nil && err != nil) || (w.err != nil && err == nil) || (err != nil && err.Error() != w.err.Error()) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotHost, w.wantHost) {
			return errors.Errorf("host got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHost, w.wantHost)
		}
		if !reflect.DeepEqual(gotPort, w.wantPort) {
			return errors.Errorf("port got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPort, w.wantPort)
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
				isLocal:  false,
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
			name: "return true if it is local address",
			args: args{
				addr: "localhost:8080",
			},
			want: want{
				wantHost: "localhost",
				wantPort: uint16(8080),
				isV4:     false,
				isV6:     false,
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
				isV4:     false,
				isV6:     false,
				isLocal:  false,
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
				err: &net.AddrError{
					Addr: "dummy",
					Err:  "missing port in address",
				},
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
				isLocal:  false,
				err: &net.AddrError{
					Addr: "192.168.1.1",
					Err:  "missing port in address",
				},
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
				err: &net.AddrError{
					Addr: "2001:db8::1",
					Err:  "too many colons in address",
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

			gotHost, gotPort, gotIsLocal, gotIsV4, gotIsV6, err := Parse(test.args.addr)
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
		wantHost string
		wantPort uint16
		err      error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string, uint16, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		start uint16
		end   uint16
		host  string
	}
	type want struct {
		wantPorts []uint16
		err       error
	}
	type test struct {
		name       string
		args       args
		want       want
		srvs       []*httptest.Server
		checkFunc  func(want, []uint16, error) error
		beforeFunc func(*testing.T, *test)
		afterFunc  func(*test)
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

		// count how want want ports exists in got ports
		cnt := 0
		for _, wp := range w.wantPorts {
			for _, gp := range gotPorts {
				if wp == gp {
					cnt++
					break
				}
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
					w.WriteHeader(200)
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
					w.WriteHeader(200)
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
					w.WriteHeader(200)
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

				for i := 0; i < srvNum; i++ {
					handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(200)
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
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
		name       string
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
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
		name       string
		args       args
		want       want
		checkFunc  func(want, NetworkType) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		name       string
		n          NetworkType
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
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
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
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
