//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.uber.org/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
}

func TestMain(m *testing.M) {
	log.Init()
	os.Exit(m.Run())
}

func TestListen(t *testing.T) {
	type args struct {
		network string
		address string
	}
	type want struct {
		want Listener
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Listener, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Listener, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "listener is created successfully",
			args: args{
				network: "tcp",
				address: ":0",
			},
			checkFunc: func(w want, got Listener, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				return got.Close()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := Listen(test.args.network, test.args.address)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsLocal(t *testing.T) {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := IsLocal(test.args.host)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDial(t *testing.T) {
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
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		func() test {
			srvContent := "test"
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, srvContent)
				w.WriteHeader(200)
			})
			testSrv := httptest.NewServer(handler)

			return test{
				name: "dial return server content",
				args: args{
					network: "tcp",
					addr:    testSrv.URL[len("http://"):],
				},
				checkFunc: func(w want, gotConn Conn, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}

					// read the output from the server and check if it is equals to the count
					fmt.Fprintf(gotConn, "GET / HTTP/1.0\r\n\r\n")
					buf, _ := ioutil.ReadAll(gotConn)
					content := strings.Split(string(buf), "\n")[5] // skip HTTP header
					if content != srvContent {
						return errors.Errorf("invalid content, got: %v, want: %v", content, srvContent)
					}

					return nil
				},
				afterFunc: func(args) {
					testSrv.CloseClientConnections()
					testSrv.Close()
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotConn, err := Dial(test.args.network, test.args.addr)
			if err := test.checkFunc(test.want, gotConn, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		addr string
	}
	type want struct {
		wantHost string
		wantPort uint16
		wantIsIP bool
		err      error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string, uint16, bool, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotHost string, gotPort uint16, gotIsIP bool, err error) error {
		if (w.err == nil && err != nil) || (w.err != nil && err == nil) || (err != nil && err.Error() != w.err.Error()) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotHost, w.wantHost) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHost, w.wantHost)
		}
		if !reflect.DeepEqual(gotPort, w.wantPort) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPort, w.wantPort)
		}
		if !reflect.DeepEqual(gotIsIP, w.wantIsIP) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIsIP, w.wantIsIP)
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
				wantIsIP: true,
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
				wantIsIP: true,
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
				wantIsIP: false,
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
				err: &strconv.NumError{
					Func: "Atoi",
					Num:  "",
					Err:  strconv.ErrSyntax,
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
				wantIsIP: true,
				err: &strconv.NumError{
					Func: "Atoi",
					Num:  "",
					Err:  strconv.ErrSyntax,
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
				wantIsIP: true,
				err: &strconv.NumError{
					Func: "Atoi",
					Num:  "",
					Err:  strconv.ErrSyntax,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotHost, gotPort, gotIsIP, err := Parse(test.args.addr)
			if err := test.checkFunc(test.want, gotHost, gotPort, gotIsIP, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	type args struct {
		addr string
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
			name: "return true if it is IPv6 address",
			args: args{
				addr: "2001:db8::1",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false if it is not IPv6 address",
			args: args{
				addr: "localhost",
			},
			want: want{
				want: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := IsIPv6(test.args.addr)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	type args struct {
		addr string
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
			name: "return true if it is IPv4 address",
			args: args{
				addr: "192.168.1.1",
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false if it is not IPv4 address",
			args: args{
				addr: "localhost",
			},
			want: want{
				want: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := IsIPv4(test.args.addr)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSplitHostPort(t *testing.T) {
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
				err:      &strconv.NumError{"Atoi", "", strconv.ErrSyntax},
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
				err:      &strconv.NumError{"Atoi", "", strconv.ErrSyntax},
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
				err:      &strconv.NumError{"Atoi", "", strconv.ErrSyntax},
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotHost, gotPort, err := SplitHostPort(test.args.hostport)
			if err := test.checkFunc(test.want, gotHost, gotPort, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestScanPorts(t *testing.T) {
	type args struct {
		ctx   context.Context
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
		checkFunc  func(want, []uint16, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		func() test {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			testSrv := httptest.NewServer(handler)

			s := strings.Split(testSrv.URL, ":")
			p, _ := strconv.Atoi(s[len(s)-1])
			srvPort := uint16(p)

			return test{
				name: "return test server port number in given range",
				args: args{
					ctx:   context.Background(),
					host:  "localhost",
					start: srvPort - 5,
					end:   srvPort + 5,
				},
				want: want{
					wantPorts: []uint16{
						srvPort,
					},
				},
				afterFunc: func(args) {
					testSrv.Close()
				},
			}
		}(),
		func() test {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			testSrv := httptest.NewServer(handler)

			s := strings.Split(testSrv.URL, ":")
			p, _ := strconv.Atoi(s[len(s)-1])
			srvPort := uint16(p)

			return test{
				name: "return test server port number when start = end",
				args: args{
					ctx:   context.Background(),
					host:  "localhost",
					start: srvPort,
					end:   srvPort,
				},
				want: want{
					wantPorts: []uint16{
						srvPort,
					},
				},
				afterFunc: func(args) {
					testSrv.Close()
				},
			}
		}(),
		func() test {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			testSrv := httptest.NewServer(handler)

			s := strings.Split(testSrv.URL, ":")
			p, _ := strconv.Atoi(s[len(s)-1])
			srvPort := uint16(p)

			return test{
				name: "return test server port number when start > end",
				args: args{
					ctx:   context.Background(),
					host:  "localhost",
					start: srvPort + 10,
					end:   srvPort - 10,
				},
				want: want{
					wantPorts: []uint16{
						srvPort,
					},
				},
				afterFunc: func(args) {
					testSrv.Close()
				},
			}
		}(),
		func() test {
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

				s := strings.Split(srv.URL, ":")
				p, _ := strconv.Atoi(s[len(s)-1])
				port := uint16(p)
				ports = append(ports, port)

				if port < minPort {
					minPort = port
				}
				if port > maxPort {
					maxPort = port
				}
			}

			return test{
				name: "return multiple test server port number",
				args: args{
					ctx:   context.Background(),
					host:  "localhost",
					start: minPort - 5,
					end:   maxPort + 5,
				},
				want: want{
					wantPorts: ports,
				},
				afterFunc: func(args) {
					for _, s := range srvs {
						s.Close()
					}
				},
			}
		}(),
		{
			name: "return no port available if no port is scanned",
			args: args{
				ctx:   context.Background(),
				host:  "localhost",
				start: 65534,
				end:   65535,
			},
			want: want{
				err: errors.ErrNoPortAvailable("localhost", 65534, 65535),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotPorts, err := ScanPorts(test.args.ctx, test.args.start, test.args.end, test.args.host)
			if err := test.checkFunc(test.want, gotPorts, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if test.afterFunc != nil {
				test.afterFunc(test.args)
			}
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := LoadLocalIP()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
