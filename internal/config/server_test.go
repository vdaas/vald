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

// Package config providers configuration type and load configuration logic
package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/servers/server"
	"go.uber.org/goleak"
)

func TestServers_Bind(t *testing.T) {
	type fields struct {
		Servers              []*Server
		HealthCheckServers   []*Server
		MetricsServers       []*Server
		StartUpStrategy      []string
		ShutdownStrategy     []string
		FullShutdownDuration string
		TLS                  *TLS
	}
	type want struct {
		want *Servers
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Servers) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Servers) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			server := &Server{
				Name:          "vald-agent-ngt",
				Network:       "tcp",
				Host:          "0.0.0.0",
				Port:          uint16(8081),
				SocketPath:    "",
				Mode:          "REST",
				ProbeWaitTime: "3s",
				HTTP: &HTTP{
					ShutdownDuration:  "5s",
					HandlerTimeout:    "5s",
					IdleTimeout:       "1s",
					ReadHeaderTimeout: "1s",
					ReadTimeout:       "1s",
					WriteTimeout:      "1s",
				},
				GRPC: &GRPC{
					BidirectionalStreamConcurrency: 20,
					MaxReceiveMessageSize:          5,
					MaxSendMessageSize:             5,
					InitialWindowSize:              1,
					InitialConnWindowSize:          1,
					Keepalive: &GRPCKeepalive{
						MaxConnIdle:     "3",
						MaxConnAge:      "30s",
						MaxConnAgeGrace: "45s",
						Time:            "60s",
						Timeout:         "90s",
					},
					WriteBufferSize:   3,
					ReadBufferSize:    3,
					ConnectionTimeout: "3s",
					MaxHeaderListSize: 5,
					HeaderTableSize:   1,
					Interceptors: []string{
						"RecoverInterceptor",
					},
					EnableReflection: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            false,
					IPRecoverDestinationAddr: false,
				},
				Restart: false,
			}
			startUpStrategy := []string{
				"vald-agent-ngt",
			}
			shutdownStrategy := []string{
				"vald-agent-ngt",
			}
			fullShutdownDuration := "600s"
			tls := &TLS{
				Enabled: false,
			}
			return test{
				name: "return Servers when all parameters are set",
				fields: fields{
					Servers: []*Server{
						server,
					},
					HealthCheckServers: []*Server{
						server,
					},
					MetricsServers: []*Server{
						server,
					},
					StartUpStrategy:      startUpStrategy,
					ShutdownStrategy:     shutdownStrategy,
					FullShutdownDuration: fullShutdownDuration,
					TLS:                  tls,
				},
				want: want{
					want: &Servers{
						Servers: []*Server{
							server,
						},
						HealthCheckServers: []*Server{
							server,
						},
						MetricsServers: []*Server{
							server,
						},
						StartUpStrategy:      startUpStrategy,
						ShutdownStrategy:     shutdownStrategy,
						FullShutdownDuration: fullShutdownDuration,
						TLS:                  tls,
					},
				},
			}
		}(),
		func() test {
			server := &Server{
				Name:          "vald-agent-ngt",
				Network:       "tcp",
				Host:          "0.0.0.0",
				Port:          uint16(8081),
				SocketPath:    "",
				Mode:          "REST",
				ProbeWaitTime: "3s",
				HTTP: &HTTP{
					ShutdownDuration:  "5s",
					HandlerTimeout:    "5s",
					IdleTimeout:       "1s",
					ReadHeaderTimeout: "1s",
					ReadTimeout:       "1s",
					WriteTimeout:      "1s",
				},
				GRPC: &GRPC{
					BidirectionalStreamConcurrency: 20,
					MaxReceiveMessageSize:          5,
					MaxSendMessageSize:             5,
					InitialWindowSize:              1,
					InitialConnWindowSize:          1,
					Keepalive: &GRPCKeepalive{
						MaxConnIdle:     "3",
						MaxConnAge:      "30s",
						MaxConnAgeGrace: "45s",
						Time:            "60s",
						Timeout:         "90s",
					},
					WriteBufferSize:   3,
					ReadBufferSize:    3,
					ConnectionTimeout: "3s",
					MaxHeaderListSize: 5,
					HeaderTableSize:   1,
					Interceptors: []string{
						"RecoverInterceptor",
					},
					EnableReflection: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            false,
					IPRecoverDestinationAddr: false,
				},
				Restart: false,
			}
			startUpStrategy := []string{
				"vald-agent-ngt",
			}
			shutdownStrategy := []string{
				"vald-agent-ngt",
			}
			fullShutdownDuration := "600s"
			tls := &TLS{
				Enabled: false,
			}
			return test{
				name: "return Servers when TLS is not set",
				fields: fields{
					Servers: []*Server{
						server,
					},
					HealthCheckServers: []*Server{
						server,
					},
					MetricsServers: []*Server{
						server,
					},
					StartUpStrategy:      startUpStrategy,
					ShutdownStrategy:     shutdownStrategy,
					FullShutdownDuration: fullShutdownDuration,
				},
				want: want{
					want: &Servers{
						Servers: []*Server{
							server,
						},
						HealthCheckServers: []*Server{
							server,
						},
						MetricsServers: []*Server{
							server,
						},
						StartUpStrategy:      startUpStrategy,
						ShutdownStrategy:     shutdownStrategy,
						FullShutdownDuration: fullShutdownDuration,
						TLS:                  tls,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Servers when all parameters are not set",
				fields: fields{},
				want: want{
					want: &Servers{
						StartUpStrategy:  make([]string, 0),
						ShutdownStrategy: make([]string, 0),
						TLS: &TLS{
							Enabled: false,
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &Servers{
				Servers:              test.fields.Servers,
				HealthCheckServers:   test.fields.HealthCheckServers,
				MetricsServers:       test.fields.MetricsServers,
				StartUpStrategy:      test.fields.StartUpStrategy,
				ShutdownStrategy:     test.fields.ShutdownStrategy,
				FullShutdownDuration: test.fields.FullShutdownDuration,
				TLS:                  test.fields.TLS,
			}

			got := s.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestServers_GetGRPCStreamConcurrency(t *testing.T) {
	type fields struct {
		Servers              []*Server
		HealthCheckServers   []*Server
		MetricsServers       []*Server
		StartUpStrategy      []string
		ShutdownStrategy     []string
		FullShutdownDuration string
		TLS                  *TLS
	}
	type want struct {
		wantC int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotC int) error {
		if !reflect.DeepEqual(gotC, w.wantC) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
		}
		return nil
	}
	tests := []test{
		func() test {
			servers := []*Server{
				{
					Name:          "vald-agent-ngt",
					Network:       "tcp",
					Host:          "0.0.0.0",
					Port:          uint16(8081),
					SocketPath:    "",
					Mode:          "GRPC",
					ProbeWaitTime: "3s",
					HTTP: &HTTP{
						ShutdownDuration:  "5s",
						HandlerTimeout:    "5s",
						IdleTimeout:       "1s",
						ReadHeaderTimeout: "1s",
						ReadTimeout:       "1s",
						WriteTimeout:      "1s",
					},
					GRPC: &GRPC{
						BidirectionalStreamConcurrency: 20,
						MaxReceiveMessageSize:          5,
						MaxSendMessageSize:             5,
						InitialWindowSize:              1,
						InitialConnWindowSize:          1,
						Keepalive: &GRPCKeepalive{
							MaxConnIdle:     "3",
							MaxConnAge:      "30s",
							MaxConnAgeGrace: "45s",
							Time:            "60s",
							Timeout:         "90s",
						},
						WriteBufferSize:   3,
						ReadBufferSize:    3,
						ConnectionTimeout: "3s",
						MaxHeaderListSize: 5,
						HeaderTableSize:   1,
						Interceptors: []string{
							"RecoverInterceptor",
						},
						EnableReflection: true,
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
					Restart: false,
				},
			}
			startUpStrategy := []string{
				"vald-agent-ngt",
			}
			shutdownStrategy := []string{
				"vald-agent-ngt",
			}
			fullShutdownDuration := "600s"
			return test{
				name: "return 20 when servers not nil and whose GRPC BidirectionalStreamConcurrency is 20",
				fields: fields{
					Servers:              servers,
					HealthCheckServers:   servers,
					MetricsServers:       servers,
					StartUpStrategy:      startUpStrategy,
					ShutdownStrategy:     shutdownStrategy,
					FullShutdownDuration: fullShutdownDuration,
				},
				want: want{
					wantC: 20,
				},
			}
		}(),
		func() test {
			servers := []*Server{
				{
					Name:          "vald-agent-ngt",
					Network:       "tcp",
					Host:          "0.0.0.0",
					Port:          uint16(8081),
					SocketPath:    "",
					Mode:          "GRPC",
					ProbeWaitTime: "3s",
					HTTP: &HTTP{
						ShutdownDuration:  "5s",
						HandlerTimeout:    "5s",
						IdleTimeout:       "1s",
						ReadHeaderTimeout: "1s",
						ReadTimeout:       "1s",
						WriteTimeout:      "1s",
					},
					SocketOption: &SocketOption{
						ReusePort:                true,
						ReuseAddr:                true,
						TCPFastOpen:              true,
						TCPNoDelay:               true,
						TCPCork:                  false,
						TCPQuickAck:              true,
						TCPDeferAccept:           true,
						IPTransparent:            false,
						IPRecoverDestinationAddr: false,
					},
					Restart: false,
				},
			}
			startUpStrategy := []string{
				"vald-agent-ngt",
			}
			shutdownStrategy := []string{
				"vald-agent-ngt",
			}
			fullShutdownDuration := "600s"
			return test{
				name: "return 0 when servers not nil and GRPC is nil",
				fields: fields{
					Servers:              servers,
					HealthCheckServers:   servers,
					MetricsServers:       servers,
					StartUpStrategy:      startUpStrategy,
					ShutdownStrategy:     shutdownStrategy,
					FullShutdownDuration: fullShutdownDuration,
				},
				want: want{
					wantC: 0,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &Servers{
				Servers:              test.fields.Servers,
				HealthCheckServers:   test.fields.HealthCheckServers,
				MetricsServers:       test.fields.MetricsServers,
				StartUpStrategy:      test.fields.StartUpStrategy,
				ShutdownStrategy:     test.fields.ShutdownStrategy,
				FullShutdownDuration: test.fields.FullShutdownDuration,
				TLS:                  test.fields.TLS,
			}

			gotC := s.GetGRPCStreamConcurrency()
			if err := test.checkFunc(test.want, gotC); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestHTTP_Bind(t *testing.T) {
	type fields struct {
		ShutdownDuration  string
		HandlerTimeout    string
		IdleTimeout       string
		ReadHeaderTimeout string
		ReadTimeout       string
		WriteTimeout      string
	}
	type want struct {
		want *HTTP
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *HTTP) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *HTTP) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			shutdownDuration := "5s"
			handlerTimeout := "5s"
			idleTimeout := "1s"
			readHeaderTimeout := "1s"
			readTimeout := "1s"
			writeTimeout := "1s"
			return test{
				name: "return HTTP when all parameters are set",
				fields: fields{
					ShutdownDuration:  shutdownDuration,
					HandlerTimeout:    handlerTimeout,
					IdleTimeout:       idleTimeout,
					ReadHeaderTimeout: readHeaderTimeout,
					ReadTimeout:       readTimeout,
					WriteTimeout:      writeTimeout,
				},
				want: want{
					want: &HTTP{
						ShutdownDuration:  shutdownDuration,
						HandlerTimeout:    handlerTimeout,
						IdleTimeout:       idleTimeout,
						ReadHeaderTimeout: readHeaderTimeout,
						ReadTimeout:       readTimeout,
						WriteTimeout:      writeTimeout,
					},
				},
			}
		}(),
		func() test {
			p := map[string]string{
				"SHUTDOWN_DURATION":   "5s",
				"HANDLER_TIMEOUT":     "5s",
				"IDLE_TIMEOUT":        "1s",
				"READ_HEADER_TIMEOUT": "1s",
				"READ_TIMEOUT":        "1s",
				"WRITE_TIMEOUT":       "1s",
			}
			return test{
				name: "return HTTP when all parameters are set as environment value",
				fields: fields{
					ShutdownDuration:  "_SHUTDOWN_DURATION_",
					HandlerTimeout:    "_HANDLER_TIMEOUT_",
					IdleTimeout:       "_IDLE_TIMEOUT_",
					ReadHeaderTimeout: "_READ_HEADER_TIMEOUT_",
					ReadTimeout:       "_READ_TIMEOUT_",
					WriteTimeout:      "_WRITE_TIMEOUT_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range p {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &HTTP{
						ShutdownDuration:  "5s",
						HandlerTimeout:    "5s",
						IdleTimeout:       "1s",
						ReadHeaderTimeout: "1s",
						ReadTimeout:       "1s",
						WriteTimeout:      "1s",
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return HTTP when all parameters are not set",
				fields: fields{},
				want: want{
					want: &HTTP{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			h := &HTTP{
				ShutdownDuration:  test.fields.ShutdownDuration,
				HandlerTimeout:    test.fields.HandlerTimeout,
				IdleTimeout:       test.fields.IdleTimeout,
				ReadHeaderTimeout: test.fields.ReadHeaderTimeout,
				ReadTimeout:       test.fields.ReadTimeout,
				WriteTimeout:      test.fields.WriteTimeout,
			}

			got := h.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGRPC_Bind(t *testing.T) {
	type fields struct {
		BidirectionalStreamConcurrency int
		MaxReceiveMessageSize          int
		MaxSendMessageSize             int
		InitialWindowSize              int
		InitialConnWindowSize          int
		Keepalive                      *GRPCKeepalive
		WriteBufferSize                int
		ReadBufferSize                 int
		ConnectionTimeout              string
		MaxHeaderListSize              int
		HeaderTableSize                int
		Interceptors                   []string
		EnableReflection               bool
	}
	type want struct {
		want *GRPC
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *GRPC) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *GRPC) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			bidirectionalStreamConcurrency := 20
			maxReceiveMessageSize := 5
			maxSendMessageSize := 5
			initialWindowSize := 1
			initialConnWindowSize := 1
			keepalive := &GRPCKeepalive{
				MaxConnIdle:     "3",
				MaxConnAge:      "30s",
				MaxConnAgeGrace: "45s",
				Time:            "60s",
				Timeout:         "90s",
			}
			writeBufferSize := 3
			readBufferSize := 3
			connectionTimeout := "3s"
			maxHeaderListSize := 5
			headerTableSize := 1
			interceptors := []string{
				"RecoverInterceptor",
			}
			enableReflection := true
			return test{
				name: "return GRPC when all parameters are set",
				fields: fields{
					BidirectionalStreamConcurrency: bidirectionalStreamConcurrency,
					MaxReceiveMessageSize:          maxReceiveMessageSize,
					MaxSendMessageSize:             maxSendMessageSize,
					InitialWindowSize:              initialWindowSize,
					InitialConnWindowSize:          initialConnWindowSize,
					Keepalive:                      keepalive,
					WriteBufferSize:                writeBufferSize,
					ReadBufferSize:                 readBufferSize,
					ConnectionTimeout:              connectionTimeout,
					MaxHeaderListSize:              maxHeaderListSize,
					HeaderTableSize:                headerTableSize,
					Interceptors:                   interceptors,
					EnableReflection:               enableReflection,
				},
				want: want{
					want: &GRPC{
						BidirectionalStreamConcurrency: bidirectionalStreamConcurrency,
						MaxReceiveMessageSize:          maxReceiveMessageSize,
						MaxSendMessageSize:             maxSendMessageSize,
						InitialWindowSize:              initialWindowSize,
						InitialConnWindowSize:          initialConnWindowSize,
						Keepalive:                      keepalive,
						WriteBufferSize:                writeBufferSize,
						ReadBufferSize:                 readBufferSize,
						ConnectionTimeout:              connectionTimeout,
						MaxHeaderListSize:              maxHeaderListSize,
						HeaderTableSize:                headerTableSize,
						Interceptors:                   interceptors,
						EnableReflection:               enableReflection,
					},
				},
			}
		}(),
		func() test {
			p := map[string]string{
				"CONNECTION_TIMEOUT": "3s",
				"INTERCEPTORS":       "RecoverInterceptor",
			}
			bidirectionalStreamConcurrency := 20
			maxReceiveMessageSize := 5
			maxSendMessageSize := 5
			initialWindowSize := 1
			initialConnWindowSize := 1
			keepalive := &GRPCKeepalive{
				MaxConnIdle:     "3",
				MaxConnAge:      "30s",
				MaxConnAgeGrace: "45s",
				Time:            "60s",
				Timeout:         "90s",
			}
			writeBufferSize := 3
			readBufferSize := 3
			maxHeaderListSize := 5
			headerTableSize := 1
			enableReflection := true
			return test{
				name: "return GRPC when some parameters are set as environment value",
				fields: fields{
					BidirectionalStreamConcurrency: bidirectionalStreamConcurrency,
					MaxReceiveMessageSize:          maxReceiveMessageSize,
					MaxSendMessageSize:             maxSendMessageSize,
					InitialWindowSize:              initialWindowSize,
					InitialConnWindowSize:          initialConnWindowSize,
					Keepalive:                      keepalive,
					WriteBufferSize:                writeBufferSize,
					ReadBufferSize:                 readBufferSize,
					ConnectionTimeout:              "_CONNECTION_TIMEOUT_",
					MaxHeaderListSize:              maxHeaderListSize,
					HeaderTableSize:                headerTableSize,
					Interceptors: []string{
						"_INTERCEPTORS_",
					},
					EnableReflection: enableReflection,
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range p {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &GRPC{
						BidirectionalStreamConcurrency: bidirectionalStreamConcurrency,
						MaxReceiveMessageSize:          maxReceiveMessageSize,
						MaxSendMessageSize:             maxSendMessageSize,
						InitialWindowSize:              initialWindowSize,
						InitialConnWindowSize:          initialConnWindowSize,
						Keepalive:                      keepalive,
						WriteBufferSize:                writeBufferSize,
						ReadBufferSize:                 readBufferSize,
						ConnectionTimeout:              "3s",
						MaxHeaderListSize:              maxHeaderListSize,
						HeaderTableSize:                headerTableSize,
						Interceptors: []string{
							"RecoverInterceptor",
						},
						EnableReflection: enableReflection,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return GRPC when all parameters are not set",
				fields: fields{},
				want: want{
					want: &GRPC{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &GRPC{
				BidirectionalStreamConcurrency: test.fields.BidirectionalStreamConcurrency,
				MaxReceiveMessageSize:          test.fields.MaxReceiveMessageSize,
				MaxSendMessageSize:             test.fields.MaxSendMessageSize,
				InitialWindowSize:              test.fields.InitialWindowSize,
				InitialConnWindowSize:          test.fields.InitialConnWindowSize,
				Keepalive:                      test.fields.Keepalive,
				WriteBufferSize:                test.fields.WriteBufferSize,
				ReadBufferSize:                 test.fields.ReadBufferSize,
				ConnectionTimeout:              test.fields.ConnectionTimeout,
				MaxHeaderListSize:              test.fields.MaxHeaderListSize,
				HeaderTableSize:                test.fields.HeaderTableSize,
				Interceptors:                   test.fields.Interceptors,
				EnableReflection:               test.fields.EnableReflection,
			}

			got := g.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGRPCKeepalive_Bind(t *testing.T) {
	type fields struct {
		MaxConnIdle     string
		MaxConnAge      string
		MaxConnAgeGrace string
		Time            string
		Timeout         string
	}
	type want struct {
		want *GRPCKeepalive
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *GRPCKeepalive) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *GRPCKeepalive) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			maxConnIdle := "3"
			maxConnAge := "30s"
			maxConnAgeGrace := "45s"
			time := "60s"
			timeout := "90s"
			return test{
				name: "return GPRCKeepalive when the parameters are set",
				fields: fields{
					MaxConnIdle:     maxConnIdle,
					MaxConnAge:      maxConnAge,
					MaxConnAgeGrace: maxConnAgeGrace,
					Time:            time,
					Timeout:         timeout,
				},
				want: want{
					want: &GRPCKeepalive{
						MaxConnIdle:     maxConnIdle,
						MaxConnAge:      maxConnAge,
						MaxConnAgeGrace: maxConnAgeGrace,
						Time:            time,
						Timeout:         timeout,
					},
				},
			}
		}(),
		func() test {
			p := map[string]string{
				"MAX_CONN_IDLE":      "3",
				"MAX_CONN_AGE":       "30s",
				"MAX_CONN_AGE_GRACE": "45s",
				"TIME":               "60s",
				"TIMEOUT":            "90s",
			}
			return test{
				name: "return GPRCKeepalive when the parameters are set as environment value",
				fields: fields{
					MaxConnIdle:     "_MAX_CONN_IDLE_",
					MaxConnAge:      "_MAX_CONN_AGE_",
					MaxConnAgeGrace: "_MAX_CONN_AGE_GRACE_",
					Time:            "_TIME_",
					Timeout:         "_TIMEOUT_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range p {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &GRPCKeepalive{
						MaxConnIdle:     "3",
						MaxConnAge:      "30s",
						MaxConnAgeGrace: "45s",
						Time:            "60s",
						Timeout:         "90s",
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return GPRCKeepalive when the parameters are not set",
				fields: fields{},
				want: want{
					want: &GRPCKeepalive{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			k := &GRPCKeepalive{
				MaxConnIdle:     test.fields.MaxConnIdle,
				MaxConnAge:      test.fields.MaxConnAge,
				MaxConnAgeGrace: test.fields.MaxConnAgeGrace,
				Time:            test.fields.Time,
				Timeout:         test.fields.Timeout,
			}

			got := k.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestServer_Bind(t *testing.T) {
	type fields struct {
		Name          string
		Network       string
		Host          string
		Port          uint16
		SocketPath    string
		Mode          string
		ProbeWaitTime string
		HTTP          *HTTP
		GRPC          *GRPC
		SocketOption  *SocketOption
		Restart       bool
	}
	type want struct {
		want *Server
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Server) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Server) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			name := "vald-agent-ngt"
			network := "tcp"
			host := "0.0.0.0"
			port := uint16(8081)
			socketPath := "/var/run/docker.sock"
			mode := "REST"
			probeWaitTime := "3s"
			http := &HTTP{
				ShutdownDuration:  "5s",
				HandlerTimeout:    "5s",
				IdleTimeout:       "1s",
				ReadHeaderTimeout: "1s",
				ReadTimeout:       "1s",
				WriteTimeout:      "1s",
			}
			grpc := &GRPC{
				BidirectionalStreamConcurrency: 20,
				MaxReceiveMessageSize:          5,
				MaxSendMessageSize:             5,
				InitialWindowSize:              1,
				InitialConnWindowSize:          1,
				Keepalive: &GRPCKeepalive{
					MaxConnIdle:     "3",
					MaxConnAge:      "30s",
					MaxConnAgeGrace: "45s",
					Time:            "60s",
					Timeout:         "90s",
				},
				WriteBufferSize:   3,
				ReadBufferSize:    3,
				ConnectionTimeout: "3s",
				MaxHeaderListSize: 5,
				HeaderTableSize:   1,
				Interceptors: []string{
					"RecoverInterceptor",
				},
				EnableReflection: true,
			}
			socketOption := &SocketOption{
				ReusePort:                true,
				ReuseAddr:                true,
				TCPFastOpen:              true,
				TCPNoDelay:               true,
				TCPCork:                  false,
				TCPQuickAck:              true,
				TCPDeferAccept:           true,
				IPTransparent:            false,
				IPRecoverDestinationAddr: false,
			}
			return test{
				name: "return Server when all parameters are set",
				fields: fields{
					Name:          name,
					Network:       network,
					Host:          host,
					Port:          port,
					SocketPath:    socketPath,
					Mode:          mode,
					ProbeWaitTime: probeWaitTime,
					HTTP:          http,
					GRPC:          grpc,
					SocketOption:  socketOption,
					Restart:       false,
				},
				want: want{
					want: &Server{
						Name:          name,
						Network:       network,
						Host:          host,
						Port:          port,
						SocketPath:    socketPath,
						Mode:          mode,
						ProbeWaitTime: probeWaitTime,
						HTTP:          http,
						GRPC:          grpc,
						SocketOption:  socketOption,
						Restart:       false,
					},
				},
			}
		}(),
		func() test {
			p := map[string]string{
				"NAME":            "vald-agent-ngt",
				"NETWORK":         "tcp",
				"HOST":            "0.0.0.0",
				"SOCKET_PATH":     "/var/run/docker.sock",
				"MODE":            "REST",
				"PROBE_WAIT_TIME": "3s",
			}
			port := uint16(8081)
			http := &HTTP{
				ShutdownDuration:  "5s",
				HandlerTimeout:    "5s",
				IdleTimeout:       "1s",
				ReadHeaderTimeout: "1s",
				ReadTimeout:       "1s",
				WriteTimeout:      "1s",
			}
			grpc := &GRPC{
				BidirectionalStreamConcurrency: 20,
				MaxReceiveMessageSize:          5,
				MaxSendMessageSize:             5,
				InitialWindowSize:              1,
				InitialConnWindowSize:          1,
				Keepalive: &GRPCKeepalive{
					MaxConnIdle:     "3",
					MaxConnAge:      "30s",
					MaxConnAgeGrace: "45s",
					Time:            "60s",
					Timeout:         "90s",
				},
				WriteBufferSize:   3,
				ReadBufferSize:    3,
				ConnectionTimeout: "3s",
				MaxHeaderListSize: 5,
				HeaderTableSize:   1,
				Interceptors: []string{
					"RecoverInterceptor",
				},
				EnableReflection: true,
			}
			socketOption := &SocketOption{
				ReusePort:                true,
				ReuseAddr:                true,
				TCPFastOpen:              true,
				TCPNoDelay:               true,
				TCPCork:                  false,
				TCPQuickAck:              true,
				TCPDeferAccept:           true,
				IPTransparent:            false,
				IPRecoverDestinationAddr: false,
			}
			return test{
				name: "return Server when all parameters are set",
				fields: fields{
					Name:          "_NAME_",
					Network:       "_NETWORK_",
					Host:          "_HOST_",
					Port:          port,
					SocketPath:    "_SOCKET_PATH_",
					Mode:          "_MODE_",
					ProbeWaitTime: "_PROBE_WAIT_TIME_",
					HTTP:          http,
					GRPC:          grpc,
					SocketOption:  socketOption,
					Restart:       false,
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range p {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &Server{
						Name:          "vald-agent-ngt",
						Network:       "tcp",
						Host:          "0.0.0.0",
						Port:          port,
						SocketPath:    "/var/run/docker.sock",
						Mode:          "REST",
						ProbeWaitTime: "3s",
						HTTP:          http,
						GRPC:          grpc,
						SocketOption:  socketOption,
						Restart:       false,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Server when all parameters are not set",
				fields: fields{},
				want: want{
					want: &Server{
						SocketOption: &SocketOption{},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &Server{
				Name:          test.fields.Name,
				Network:       test.fields.Network,
				Host:          test.fields.Host,
				Port:          test.fields.Port,
				SocketPath:    test.fields.SocketPath,
				Mode:          test.fields.Mode,
				ProbeWaitTime: test.fields.ProbeWaitTime,
				HTTP:          test.fields.HTTP,
				GRPC:          test.fields.GRPC,
				SocketOption:  test.fields.SocketOption,
				Restart:       test.fields.Restart,
			}

			got := s.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestServer_Opts(t *testing.T) {
	type fields struct {
		Name          string
		Network       string
		Host          string
		Port          uint16
		SocketPath    string
		Mode          string
		ProbeWaitTime string
		HTTP          *HTTP
		GRPC          *GRPC
		SocketOption  *SocketOption
		Restart       bool
	}
	type want struct {
		want []server.Option
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []server.Option) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []server.Option) error {
		if !reflect.DeepEqual(len(got), len(w.want)) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return 13 server.Options when NETWORK is tcp, MODE is REST",
			fields: fields{
				Name:          "vald-agent-ngt",
				Network:       "tcp",
				Host:          "0.0.0.0",
				Port:          uint16(8081),
				SocketPath:    "",
				Mode:          "REST",
				ProbeWaitTime: "3s",
				HTTP: &HTTP{
					ShutdownDuration:  "5s",
					HandlerTimeout:    "5s",
					IdleTimeout:       "1s",
					ReadHeaderTimeout: "1s",
					ReadTimeout:       "1s",
					WriteTimeout:      "1s",
				},
				GRPC: &GRPC{
					BidirectionalStreamConcurrency: 20,
					MaxReceiveMessageSize:          5,
					MaxSendMessageSize:             5,
					InitialWindowSize:              1,
					InitialConnWindowSize:          1,
					Keepalive: &GRPCKeepalive{
						MaxConnIdle:     "3",
						MaxConnAge:      "30s",
						MaxConnAgeGrace: "45s",
						Time:            "60s",
						Timeout:         "90s",
					},
					WriteBufferSize:   3,
					ReadBufferSize:    3,
					ConnectionTimeout: "3s",
					MaxHeaderListSize: 5,
					HeaderTableSize:   1,
					Interceptors: []string{
						"RecoverInterceptor",
					},
					EnableReflection: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            false,
					IPRecoverDestinationAddr: false,
				},
				Restart: false,
			},
			want: want{
				want: make([]server.Option, 13),
			},
		},
		{
			name: "return 13 server.Options when NETWORK is empty, MODE is REST",
			fields: fields{
				Name:          "vald-agent-ngt",
				Host:          "0.0.0.0",
				Port:          uint16(8081),
				SocketPath:    "",
				Mode:          "REST",
				ProbeWaitTime: "3s",
				HTTP: &HTTP{
					ShutdownDuration:  "5s",
					HandlerTimeout:    "5s",
					IdleTimeout:       "1s",
					ReadHeaderTimeout: "1s",
					ReadTimeout:       "1s",
					WriteTimeout:      "1s",
				},
				GRPC: &GRPC{
					BidirectionalStreamConcurrency: 20,
					MaxReceiveMessageSize:          5,
					MaxSendMessageSize:             5,
					InitialWindowSize:              1,
					InitialConnWindowSize:          1,
					Keepalive: &GRPCKeepalive{
						MaxConnIdle:     "3",
						MaxConnAge:      "30s",
						MaxConnAgeGrace: "45s",
						Time:            "60s",
						Timeout:         "90s",
					},
					WriteBufferSize:   3,
					ReadBufferSize:    3,
					ConnectionTimeout: "3s",
					MaxHeaderListSize: 5,
					HeaderTableSize:   1,
					Interceptors: []string{
						"RecoverInterceptor",
					},
					EnableReflection: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            false,
					IPRecoverDestinationAddr: false,
				},
				Restart: false,
			},
			want: want{
				want: make([]server.Option, 13),
			},
		},
		{
			name: "return 13 server.Options when NETWORK is empty, MODE is GRPC",
			fields: fields{
				Name:          "vald-agent-ngt",
				Host:          "0.0.0.0",
				Port:          uint16(8081),
				SocketPath:    "",
				Mode:          "GRPC",
				ProbeWaitTime: "3s",
				HTTP: &HTTP{
					ShutdownDuration:  "5s",
					HandlerTimeout:    "5s",
					IdleTimeout:       "1s",
					ReadHeaderTimeout: "1s",
					ReadTimeout:       "1s",
					WriteTimeout:      "1s",
				},
				GRPC: &GRPC{
					BidirectionalStreamConcurrency: 20,
					MaxReceiveMessageSize:          5,
					MaxSendMessageSize:             5,
					InitialWindowSize:              1,
					InitialConnWindowSize:          1,
					Keepalive: &GRPCKeepalive{
						MaxConnIdle:     "3",
						MaxConnAge:      "30s",
						MaxConnAgeGrace: "45s",
						Time:            "60s",
						Timeout:         "90s",
					},
					WriteBufferSize:   3,
					ReadBufferSize:    3,
					ConnectionTimeout: "3s",
					MaxHeaderListSize: 5,
					HeaderTableSize:   1,
					Interceptors: []string{
						"RecoverInterceptor",
					},
					EnableReflection: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            false,
					IPRecoverDestinationAddr: false,
				},
				Restart: false,
			},
			want: want{
				want: make([]server.Option, 25),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &Server{
				Name:          test.fields.Name,
				Network:       test.fields.Network,
				Host:          test.fields.Host,
				Port:          test.fields.Port,
				SocketPath:    test.fields.SocketPath,
				Mode:          test.fields.Mode,
				ProbeWaitTime: test.fields.ProbeWaitTime,
				HTTP:          test.fields.HTTP,
				GRPC:          test.fields.GRPC,
				SocketOption:  test.fields.SocketOption,
				Restart:       test.fields.Restart,
			}

			got := s.Opts()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
