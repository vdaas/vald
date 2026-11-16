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
package grpc

import (
	"fmt"
	"math"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/faiss/service"
)

// TestMain initializes the global logger and info provider once, before any parallel test runs,
// and checks the package for leaked goroutines.
func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	if err := info.Init(""); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	goleak.VerifyTestMain(m)
}

// stubFaiss returns a fixed error from Search; every other method panics through the nil embedded interface.
type stubFaiss struct {
	service.Faiss
	err error
	dim int
}

func (f *stubFaiss) GetDimensionSize() int { return f.dim }

func (f *stubFaiss) Search(_, _, _ uint32, _ []float32) (*payload.Search_Response, error) {
	return nil, f.err
}

// stubFaissWrite records whether the service was reached; every other method panics through the nil embedded interface.
type stubFaissWrite struct {
	stubFaiss
	called bool
}

func (f *stubFaissWrite) InsertWithTime(string, []float32, int64) error {
	f.called = true
	return nil
}

func (f *stubFaissWrite) UpdateWithTime(string, []float32, int64) error {
	f.called = true
	return nil
}

// nanOrInfValues are the non-finite vector components that must be rejected up front.
func nanOrInfValues() map[string]float32 {
	return map[string]float32{
		"NaN":  float32(math.NaN()),
		"+Inf": float32(math.Inf(1)),
		"-Inf": float32(math.Inf(-1)),
	}
}

func wantInvalidArgument(t *testing.T, err error) {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("got non-status error: %v", err)
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("code = %v, want %v", st.Code(), codes.InvalidArgument)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Server
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Server, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Server, err error) error {
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
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_newLocations(t *testing.T) {
// 	type args struct {
// 		uuids []string
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations) error {
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuids:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           uuids:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			gotLocs := s.newLocations(test.args.uuids...)
// 			if err := checkFunc(test.want, gotLocs); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_newLocation(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		want *payload.Object_Location
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Location) error {
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
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			got := s.newLocation(test.args.uuid)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
