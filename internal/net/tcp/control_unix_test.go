// +build linux,!windows,!wasm,!js,!darwin

//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package tcp provides tcp option
package tcp

import (
	"net"
	"syscall"
	"testing"

	"github.com/vdaas/vald/internal/errors"

	"go.uber.org/goleak"
)

func TestControl(t *testing.T) {
	type args struct {
		network string
		address string
		c       syscall.RawConn
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		argsFunc   func(*testing.T) args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns nil when no error occurs internally",
			argsFunc: func(t *testing.T) args {
				t.Helper()

				addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
				if err != nil {
					t.Fatal(err)
				}
				ls, err := net.ListenTCP("tcp", addr)
				if err != nil {
					t.Fatal(err)
				}
				c, err := ls.SyscallConn()
				if err != nil {
					t.Fatal(err)
				}

				return args{
					network: "tcp",
					address: "127.0.0.1:1234",
					c:       c,
				}
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)

			args := test.argsFunc(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := Control(args.network, args.address, args.c)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
