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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/servers/server"
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Servers: nil,
		           HealthCheckServers: nil,
		           MetricsServers: nil,
		           StartUpStrategy: nil,
		           ShutdownStrategy: nil,
		           FullShutdownDuration: "",
		           TLS: TLS{},
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           Servers: nil,
		           HealthCheckServers: nil,
		           MetricsServers: nil,
		           StartUpStrategy: nil,
		           ShutdownStrategy: nil,
		           FullShutdownDuration: "",
		           TLS: TLS{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Servers: nil,
		           HealthCheckServers: nil,
		           MetricsServers: nil,
		           StartUpStrategy: nil,
		           ShutdownStrategy: nil,
		           FullShutdownDuration: "",
		           TLS: TLS{},
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           Servers: nil,
		           HealthCheckServers: nil,
		           MetricsServers: nil,
		           StartUpStrategy: nil,
		           ShutdownStrategy: nil,
		           FullShutdownDuration: "",
		           TLS: TLS{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *HTTP) error {
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
		       fields: fields {
		           ShutdownDuration: "",
		           HandlerTimeout: "",
		           IdleTimeout: "",
		           ReadHeaderTimeout: "",
		           ReadTimeout: "",
		           WriteTimeout: "",
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           ShutdownDuration: "",
		           HandlerTimeout: "",
		           IdleTimeout: "",
		           ReadHeaderTimeout: "",
		           ReadTimeout: "",
		           WriteTimeout: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
	}
	type want struct {
		want *GRPC
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *GRPC) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *GRPC) error {
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
		       fields: fields {
		           BidirectionalStreamConcurrency: 0,
		           MaxReceiveMessageSize: 0,
		           MaxSendMessageSize: 0,
		           InitialWindowSize: 0,
		           InitialConnWindowSize: 0,
		           Keepalive: GRPCKeepalive{},
		           WriteBufferSize: 0,
		           ReadBufferSize: 0,
		           ConnectionTimeout: "",
		           MaxHeaderListSize: 0,
		           HeaderTableSize: 0,
		           Interceptors: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           BidirectionalStreamConcurrency: 0,
		           MaxReceiveMessageSize: 0,
		           MaxSendMessageSize: 0,
		           InitialWindowSize: 0,
		           InitialConnWindowSize: 0,
		           Keepalive: GRPCKeepalive{},
		           WriteBufferSize: 0,
		           ReadBufferSize: 0,
		           ConnectionTimeout: "",
		           MaxHeaderListSize: 0,
		           HeaderTableSize: 0,
		           Interceptors: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *GRPCKeepalive) error {
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
		       fields: fields {
		           MaxConnIdle: "",
		           MaxConnAge: "",
		           MaxConnAgeGrace: "",
		           Time: "",
		           Timeout: "",
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           MaxConnIdle: "",
		           MaxConnAge: "",
		           MaxConnAgeGrace: "",
		           Time: "",
		           Timeout: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
		Host          string
		Port          uint
		Mode          string
		ProbeWaitTime string
		HTTP          *HTTP
		GRPC          *GRPC
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Server) error {
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
		       fields: fields {
		           Name: "",
		           Host: "",
		           Port: 0,
		           Mode: "",
		           ProbeWaitTime: "",
		           HTTP: HTTP{},
		           GRPC: GRPC{},
		           Restart: false,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           Name: "",
		           Host: "",
		           Port: 0,
		           Mode: "",
		           ProbeWaitTime: "",
		           HTTP: HTTP{},
		           GRPC: GRPC{},
		           Restart: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			s := &Server{
				Name:          test.fields.Name,
				Host:          test.fields.Host,
				Port:          test.fields.Port,
				Mode:          test.fields.Mode,
				ProbeWaitTime: test.fields.ProbeWaitTime,
				HTTP:          test.fields.HTTP,
				GRPC:          test.fields.GRPC,
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
		Host          string
		Port          uint
		Mode          string
		ProbeWaitTime string
		HTTP          *HTTP
		GRPC          *GRPC
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
		       fields: fields {
		           Name: "",
		           Host: "",
		           Port: 0,
		           Mode: "",
		           ProbeWaitTime: "",
		           HTTP: HTTP{},
		           GRPC: GRPC{},
		           Restart: false,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           Name: "",
		           Host: "",
		           Port: 0,
		           Mode: "",
		           ProbeWaitTime: "",
		           HTTP: HTTP{},
		           GRPC: GRPC{},
		           Restart: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			s := &Server{
				Name:          test.fields.Name,
				Host:          test.fields.Host,
				Port:          test.fields.Port,
				Mode:          test.fields.Mode,
				ProbeWaitTime: test.fields.ProbeWaitTime,
				HTTP:          test.fields.HTTP,
				GRPC:          test.fields.GRPC,
				Restart:       test.fields.Restart,
			}

			got := s.Opts()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
