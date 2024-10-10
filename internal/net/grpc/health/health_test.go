//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package health provides generic functionality for grpc health checks.
package health

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/test/goleak"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRegister(t *testing.T) {
	t.Parallel()
	type args struct {
		srv *grpc.Server
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		func() test {
			srv := grpc.NewServer()
			return test{
				name: "success to register the health check server",
				args: args{
					srv: srv,
				},
				checkFunc: func(w want) error {
					if _, ok := srv.GetServiceInfo()[healthpb.Health_ServiceDesc.ServiceName]; !ok {
						return errors.New("health check server not registered")
					}

					return nil
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

			Register(test.args.srv)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
