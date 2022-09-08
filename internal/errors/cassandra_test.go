//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package errors provides error types and function
package errors

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestErrCassandraInvalidConsistencyType(t *testing.T) {
	type args struct {
		consistency string
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
			name: "returns error when consistency level is 'QUORUM'",
			args: args{
				consistency: "QUORUM",
			},
			want: want{
				want: New("consistetncy type \"QUORUM\" is not defined"),
			},
		},
		{
			name: "returns error when consistency level is empty",
			args: args{
				consistency: "",
			},
			want: want{
				want: New("consistetncy type \"\" is not defined"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraInvalidConsistencyType(test.args.consistency)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewErrCassandraNotFoundIdentity(t *testing.T) {
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
			name: "returns cassandra not found identity error",
			want: want{
				want: &ErrCassandraNotFoundIdentity{
					err: New("cassandra entry not found"),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := NewErrCassandraNotFoundIdentity()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewErrCassandraUnavailableIdentity(t *testing.T) {
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
			name: "returns cassandra unavailable identity error",
			want: want{
				want: &ErrCassandraUnavailableIdentity{
					err: New("cassandra unavailable"),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := NewErrCassandraUnavailableIdentity()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraUnavailable(t *testing.T) {
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
			name: "returns cassandra unavailable identity error",
			want: want{
				want: &ErrCassandraUnavailableIdentity{
					err: New("cassandra unavailable"),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraUnavailable()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraNotFound(t *testing.T) {
	type args struct {
		keys []string
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
			name: "returns cassandra key not found error when keys is not found",
			args: args{
				keys: []string{
					"uuid",
				},
			},
			want: want{
				want: New("cassandra key 'uuid' not found: cassandra entry not found"),
			},
		},
		{
			name: "returns cassandra keys not found error when keys are not found",
			args: args{
				keys: []string{
					"uuid_1",
					"uuid_2",
				},
			},
			want: want{
				want: New("cassandra keys '[uuid_1 uuid_2]' not found: cassandra entry not found"),
			},
		},
		{
			name: "returns nil when keys is nil",
			args: args{
				keys: nil,
			},
			want: want{
				want: nil,
			},
		},
		{
			name: "returns nil when keys is empty",
			args: args{
				keys: []string{},
			},
			want: want{
				want: nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraNotFound(test.args.keys...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraGetOperationFailed(t *testing.T) {
	type args struct {
		key string
		err error
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
			name: "returns wrapped fetch key error when key is 'uuid' and error is database error",
			args: args{
				key: "uuid",
				err: New("database error"),
			},
			want: want{
				want: New("error failed to fetch key (uuid): database error"),
			},
		},
		{
			name: "returns wrapped fetch key error when key is empty and error is database error",
			args: args{
				key: "",
				err: New("database error"),
			},
			want: want{
				want: New("error failed to fetch key (): database error"),
			},
		},
		{
			name: "returns fetch key error when key is 'uuid' and error is nil",
			args: args{
				key: "uuid",
				err: nil,
			},
			want: want{
				want: New("error failed to fetch key (uuid)"),
			},
		},
		{
			name: "returns fetch key error when key is empty and error is nil",
			args: args{
				key: "",
				err: nil,
			},
			want: want{
				want: New("error failed to fetch key ()"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraGetOperationFailed(test.args.key, test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraSetOperationFailed(t *testing.T) {
	type args struct {
		key string
		err error
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
			name: "returns wrapped set key error when key is 'uuid' and error is database error",
			args: args{
				key: "uuid",
				err: New("database error"),
			},
			want: want{
				want: New("error failed to set key (uuid): database error"),
			},
		},
		{
			name: "returns wrapped set key error when key is empty and error is database error",
			args: args{
				key: "",
				err: New("database error"),
			},
			want: want{
				want: New("error failed to set key (): database error"),
			},
		},
		{
			name: "returns set key error when key is 'uuid' and error is nil",
			args: args{
				key: "uuid",
				err: nil,
			},
			want: want{
				want: New("error failed to set key (uuid)"),
			},
		},
		{
			name: "returns set key error when key is empty and error is nil",
			args: args{
				key: "",
				err: nil,
			},
			want: want{
				want: New("error failed to set key ()"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraSetOperationFailed(test.args.key, test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraDeleteOperationFailed(t *testing.T) {
	type args struct {
		key string
		err error
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
			name: "returns wrapped delete key error when key is 'uuid' and error is database error",
			args: args{
				key: "uuid",
				err: New("database error"),
			},
			want: want{
				want: New("error failed to delete key (uuid): database error"),
			},
		},
		{
			name: "returns wrapped delete key error when key is empty and error is database error",
			args: args{
				key: "",
				err: New("database error"),
			},
			want: want{
				want: New("error failed to delete key (): database error"),
			},
		},
		{
			name: "returns delete key error when key is 'uuid' and error is nil",
			args: args{
				key: "uuid",
				err: nil,
			},
			want: want{
				want: New("error failed to delete key (uuid)"),
			},
		},
		{
			name: "returns delete key error when key is empty and error is nil",
			args: args{
				key: "",
				err: nil,
			},
			want: want{
				want: New("error failed to delete key ()"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraDeleteOperationFailed(test.args.key, test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraHostDownDetected(t *testing.T) {
	type args struct {
		nodeInfo string
		err      error
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
			name: "returns wrapped cassandra host down detected error when nodeInfo is '127.0.0.1' and error is database error",
			args: args{
				nodeInfo: "127.0.0.1",
				err:      New("database error"),
			},
			want: want{
				want: New("error cassandra host down detected\t127.0.0.1: database error"),
			},
		},
		{
			name: "returns wrapped cassandra host down detected error when nodeInfo is empty and error is database error",
			args: args{
				nodeInfo: "",
				err:      New("database error"),
			},
			want: want{
				want: New("error cassandra host down detected\t: database error"),
			},
		},
		{
			name: "returns cassandra host down detected error when nodeInfo is '127.0.0.1' and error is nil",
			args: args{
				nodeInfo: "127.0.0.1",
				err:      nil,
			},
			want: want{
				want: New("error cassandra host down detected\t127.0.0.1"),
			},
		},
		{
			name: "returns cassandra host down detected error when nodeInfo is empty and error is nil",
			args: args{
				nodeInfo: "",
				err:      nil,
			},
			want: want{
				want: New("error cassandra host down detected\t"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraHostDownDetected(test.args.err, test.args.nodeInfo)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraFailedToCreateSession(t *testing.T) {
	type args struct {
		err        error
		hosts      []string
		port       int
		cqlVersion string
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
			name: "returns wrapped cassandra failed to create session error when hosts, port and cqlVersion are not nil and error is database error",
			args: args{
				err: New("database error"),
				hosts: []string{
					"vald-cassandra-01.dev.com",
					"vald-cassandra-02.dev.com",
				},
				port:       9042,
				cqlVersion: "3.0.0",
			},
			want: want{
				want: New("error cassandra client failed to create session to hosts: [vald-cassandra-01.dev.com vald-cassandra-02.dev.com]\tport: 9042\tcql_version: 3.0.0 : database error"),
			},
		},
		{
			name: "returns wrapped cassandra failed to create session error when hosts, port and cqlVersion are not nil and error is nil",
			args: args{
				hosts: []string{
					"vald-cassandra-01.dev.com",
					"vald-cassandra-02.dev.com",
				},
				port:       9042,
				cqlVersion: "3.0.0",
			},
			want: want{
				want: New("error cassandra client failed to create session to hosts: [vald-cassandra-01.dev.com vald-cassandra-02.dev.com]\tport: 9042\tcql_version: 3.0.0 "),
			},
		},
		{
			name: "returns wrapped cassandra failed to create session error when hosts and cqlVersion are not nil and error is database error and port is nil",
			args: args{
				err: New("database error"),
				hosts: []string{
					"vald-cassandra-01.dev.com",
					"vald-cassandra-02.dev.com",
				},
				cqlVersion: "3.0.0",
			},
			want: want{
				want: New("error cassandra client failed to create session to hosts: [vald-cassandra-01.dev.com vald-cassandra-02.dev.com]\tport: 0\tcql_version: 3.0.0 : database error"),
			},
		},
		{
			name: "returns wrapped cassandra failed to create session error when hosts is nil, port and cqlVersion not nil and error is database error",
			args: args{
				err:        New("database error"),
				port:       9042,
				cqlVersion: "3.0.0",
			},
			want: want{
				want: New("error cassandra client failed to create session to hosts: []\tport: 9042\tcql_version: 3.0.0 : database error"),
			},
		},
		{
			name: "returns wrapped cassandra failed to create session error when hosts, port are not nil and cqlVersion is nil and error is database error",
			args: args{
				err: New("database error"),
				hosts: []string{
					"vald-cassandra-01.dev.com",
					"vald-cassandra-02.dev.com",
				},
				port: 9042,
			},
			want: want{
				want: New("error cassandra client failed to create session to hosts: [vald-cassandra-01.dev.com vald-cassandra-02.dev.com]\tport: 9042\tcql_version:  : database error"),
			},
		},
		{
			name: "returns wrapped cassandra failed to create session error when all of input are nil or empty",
			args: args{},
			want: want{
				want: New("error cassandra client failed to create session to hosts: []\tport: 0\tcql_version:  "),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCassandraFailedToCreateSession(test.args.err, test.args.hosts, test.args.port, test.args.cqlVersion)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraNotFoundIdentity_Error(t *testing.T) {
	type fields struct {
		err error
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns string when internal error is cassandra not found identity error",
			fields: fields{
				err: New("cassandra not found identity"),
			},
			want: want{
				want: "cassandra not found identity",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			e := &ErrCassandraNotFoundIdentity{
				err: test.fields.err,
			}

			got := e.Error()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraNotFoundIdentity_Unwrap(t *testing.T) {
	type fields struct {
		err error
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !Is(err, w.err) {
			return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns nil when internal error is nil",
				fields: fields{
					err: nil,
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			err := New("cassandra not found identity")
			return test{
				name: "returns internal error when internal error is cassandra not found identity error",
				fields: fields{
					err: err,
				},
				want: want{
					err: err,
				},
				checkFunc: func(w want, err error) error {
					if !reflect.DeepEqual(w.err, err) {
						return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			e := &ErrCassandraNotFoundIdentity{
				err: test.fields.err,
			}

			err := e.Unwrap()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsErrCassandraNotFound(t *testing.T) {
	type args struct {
		err error
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
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns false when error is not cassandra not found identity",
			args: args{
				err: New("database not found"),
			},
			want: want{
				want: false,
			},
		},
		{
			name: "returns true when error is cassandra not found identity",
			args: args{
				err: new(ErrCassandraNotFoundIdentity),
			},
			want: want{
				want: true,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := IsErrCassandraNotFound(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraUnavailableIdentity_Error(t *testing.T) {
	type fields struct {
		err error
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns string when internal error is cassandra unavailable identity error",
			fields: fields{
				err: New("cassandra unavailable identity"),
			},
			want: want{
				want: "cassandra unavailable identity",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			e := &ErrCassandraUnavailableIdentity{
				err: test.fields.err,
			}

			got := e.Error()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraUnavailableIdentity_Unwrap(t *testing.T) {
	type fields struct {
		err error
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !Is(err, w.err) {
			return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns nil when internal error is nil",
				fields: fields{
					err: nil,
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			err := New("cassandra unavailable identity")
			return test{
				name: "returns internal error when internal error is cassandra unavailable identity error",
				fields: fields{
					err: err,
				},
				want: want{
					err: err,
				},
				checkFunc: func(w want, err error) error {
					if !reflect.DeepEqual(w.err, err) {
						return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			e := &ErrCassandraUnavailableIdentity{
				err: test.fields.err,
			}

			err := e.Unwrap()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsErrCassandraUnavailable(t *testing.T) {
	type args struct {
		err error
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
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns false when error is not cassandra unavailable identity",
			args: args{
				err: New("database not found"),
			},
			want: want{
				want: false,
			},
		},
		{
			name: "returns true when error is cassandra unavailable identity",
			args: args{
				err: new(ErrCassandraUnavailableIdentity),
			},
			want: want{
				want: true,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := IsErrCassandraUnavailable(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
