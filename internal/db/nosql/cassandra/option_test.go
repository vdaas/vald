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

// Package redis provides implementation of Go API for redis interface
package cassandra

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

type queryObserverImpl struct{}

func (queryObserverImpl) ObserveQuery(context.Context, gocql.ObservedQuery) {}

type batchObserverImpl struct{}

func (batchObserverImpl) ObserveBatch(context.Context, gocql.ObservedBatch) {}

type connectObserverImpl struct{}

func (connectObserverImpl) ObserveConnect(gocql.ObservedConnect) {}

type frameHeaderObserverImpl struct{}

func (frameHeaderObserverImpl) ObserveFrameHeader(context.Context, gocql.ObservedFrameHeader) {}

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestWithHosts(t *testing.T) {
	type T = client
	type args struct {
		hosts []string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(*T)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}

		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set host success",
			args: args{
				hosts: []string{"vald.vdaas.org"},
			},
			want: want{
				obj: &T{
					hosts: []string{"vald.vdaas.org"},
				},
			},
		},
		{
			name: "return error if hosts is nil",
			args: args{
				hosts: nil,
			},
			want: want{
				obj: &T{
					hosts: nil,
				},
				err: func() error {
					var nilStr []string
					return errors.NewErrInvalidOption("hosts", nilStr)
				}(),
			},
		},
		{
			name: "set host twice success",
			args: args{
				hosts: []string{"hosts1"},
			},
			beforeFunc: func(obj *T) {
				_ = WithHosts("vald.vdaas.org")(obj)
			},
			want: want{
				obj: &T{
					hosts: []string{"vald.vdaas.org", "hosts1"},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := WithHosts(test.args.hosts...)
			obj := new(T)
			if test.beforeFunc != nil {
				test.beforeFunc(obj)
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDialer(t *testing.T) {
	type T = client
	type args struct {
		der net.Dialer
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		func() test {
			dm := &DialerMock{}
			return test{
				name: "set dialer success",
				args: args{
					der: dm,
				},
				want: want{
					obj: &T{
						dialer:    dm,
						rawDialer: dm,
					},
				},
			}
		}(),
		{
			name: "return error if dialer is nil",
			args: args{
				der: nil,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("dialer", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDialer(test.args.der)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithCQLVersion(t *testing.T) {
	type T = client
	type fields struct {
		cqlVersion string
	}
	type args struct {
		version string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		fields     fields
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set version success",
			args: args{
				version: "1.0",
			},
			want: want{
				obj: &T{
					cqlVersion: "1.0",
				},
			},
		},
		{
			name: "return error when version is empty",
			args: args{
				version: "",
			},
			fields: fields{
				cqlVersion: "1.0",
			},
			want: want{
				obj: &T{
					cqlVersion: "1.0",
				},
				err: errors.NewErrInvalidOption("cqlVersion", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithCQLVersion(test.args.version)
			obj := &T{
				cqlVersion: test.fields.cqlVersion,
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithProtoVersion(t *testing.T) {
	type T = client
	type args struct {
		version int
	}
	type fields struct {
		protoVersion int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set version success",
			args: args{
				version: 1,
			},
			want: want{
				obj: &T{
					protoVersion: 1,
				},
			},
		},
		{
			name: "return error when version < 0",
			args: args{
				version: -1,
			},
			fields: fields{
				protoVersion: 10,
			},
			want: want{
				obj: &T{
					protoVersion: 10,
				},
				err: errors.NewErrInvalidOption("protoVersion", -1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithProtoVersion(test.args.version)
			obj := &T{
				protoVersion: test.fields.protoVersion,
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type T = client
	type args struct {
		dur string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set timeout success",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					timeout: 5 * time.Second,
				},
			},
		},
		{
			name: "set timeout success if the time format is invalid",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					timeout: time.Minute,
				},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("timeout", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithTimeout(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithConnectTimeout(t *testing.T) {
	type T = client
	type args struct {
		dur string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if (w.err == nil && err != nil) || (w.err != nil && err == nil) ||
			(err != nil && w.err.Error() != err.Error()) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set timeout success",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					connectTimeout: 5 * time.Second,
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				dur: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption("connectTimeout", "dummy", errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\"")),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("connectTimeout", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithConnectTimeout(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type T = client
	type args struct {
		port int
	}
	type fields struct {
		port int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set port success",
			args: args{
				port: 8080,
			},
			want: want{
				obj: &T{
					port: 8080,
				},
			},
		},
		{
			name: "set port success (boundary)",
			args: args{
				port: 65535,
			},
			want: want{
				obj: &T{
					port: 65535,
				},
			},
		},
		{
			name: "return error when port <= 0",
			args: args{
				port: -1,
			},
			fields: fields{
				port: 8080,
			},
			want: want{
				err: errors.NewErrInvalidOption("port", -1),
				obj: &T{
					port: 8080,
				},
			},
		},
		{
			name: "return error when port == 0",
			args: args{
				port: 0,
			},
			fields: fields{
				port: 8080,
			},
			want: want{
				err: errors.NewErrInvalidOption("port", 0),
				obj: &T{
					port: 8080,
				},
			},
		},
		{
			name: "return error when port > 65535",
			args: args{
				port: 65536,
			},
			fields: fields{
				port: 8080,
			},
			want: want{
				err: errors.NewErrInvalidOption("port", 65536),
				obj: &T{
					port: 8080,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithPort(test.args.port)
			obj := &T{
				port: test.fields.port,
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithKeyspace(t *testing.T) {
	type T = client
	type args struct {
		keyspace string
	}
	type fields struct {
		keyspace string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set keyspace success",
			args: args{
				keyspace: "keyspace",
			},
			want: want{
				obj: &T{
					keyspace: "keyspace",
				},
			},
		},
		{
			name: "return error when keyspace is empty",
			args: args{
				keyspace: "",
			},
			fields: fields{
				keyspace: "keyspace",
			},
			want: want{
				obj: &T{
					keyspace: "keyspace",
				},
				err: errors.NewErrInvalidOption("keyspace", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithKeyspace(test.args.keyspace)
			obj := &T{
				keyspace: test.fields.keyspace,
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithNumConns(t *testing.T) {
	type T = client
	type args struct {
		numConns int
	}
	type fields struct {
		numConns int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set num conn success",
			args: args{
				numConns: 100,
			},
			want: want{
				obj: &T{
					numConns: 100,
				},
			},
		},
		{
			name: "return error when numConn < 0",
			args: args{
				numConns: -1,
			},
			fields: fields{
				numConns: 100,
			},
			want: want{
				obj: &T{
					numConns: 100,
				},
				err: errors.NewErrInvalidOption("numConns", -1),
			},
		},
		{
			name: "set numConn success when numConn = 0",
			args: args{
				numConns: 0,
			},
			fields: fields{
				numConns: 100,
			},
			want: want{
				obj: &T{
					numConns: 0,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithNumConns(test.args.numConns)
			obj := &T{
				numConns: test.fields.numConns,
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithConsistency(t *testing.T) {
	type T = client
	type args struct {
		consistency string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set consistency level success",
			args: args{
				consistency: "one",
			},
			want: want{
				obj: &T{
					consistency: gocql.One,
				},
			},
		},
		{
			name: "set consistency level success with complex string",
			args: args{
				consistency: "-One_",
			},
			want: want{
				obj: &T{
					consistency: gocql.One,
				},
			},
		},
		{
			name: "return error when consistency is empty",
			args: args{
				consistency: "",
			},
			want: want{
				err: errors.NewErrInvalidOption("consistency", ""),
				obj: &T{},
			},
		},
		{
			name: "return error when consistency is invalid",
			args: args{
				consistency: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption("consistency", "dummy"),
				obj: &T{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithConsistency(test.args.consistency)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithSerialConsistency(t *testing.T) {
	type T = client
	type args struct {
		consistency string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set serial consistency level success",
			args: args{
				consistency: "serial",
			},
			want: want{
				obj: &T{
					serialConsistency: gocql.Serial,
				},
			},
		},
		{
			name: "set serial consistency level success with complex string",
			args: args{
				consistency: "-serial_",
			},
			want: want{
				obj: &T{
					serialConsistency: gocql.Serial,
				},
			},
		},
		{
			name: "return error when consistency is empty",
			args: args{
				consistency: "",
			},
			want: want{
				err: errors.NewErrInvalidOption("serialConsistency", ""),
				obj: &T{},
			},
		},
		{
			name: "return error when consistency is invalid",
			args: args{
				consistency: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption("serialConsistency", "dummy"),
				obj: &T{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithSerialConsistency(test.args.consistency)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithCompressor(t *testing.T) {
	type T = client
	type args struct {
		compressor gocql.Compressor
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set compressor success",
			args: args{
				compressor: &gocql.SnappyCompressor{},
			},
			want: want{
				obj: &T{
					compressor: &gocql.SnappyCompressor{},
				},
			},
		},
		{
			name: "return error when compressor is nil",
			args: args{},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("compressor", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithCompressor(test.args.compressor)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithUsername(t *testing.T) {
	type T = client
	type args struct {
		username string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set username success",
			args: args{
				username: "user",
			},
			want: want{
				obj: &T{
					username: "user",
				},
			},
		},
		{
			name: "return error when username is empty",
			args: args{
				username: "",
			},
			want: want{
				err: errors.NewErrInvalidOption("username", ""),
				obj: &T{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithUsername(test.args.username)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithPassword(t *testing.T) {
	type T = client
	type args struct {
		password string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set password success",
			args: args{
				password: "pass",
			},
			want: want{
				obj: &T{
					password: "pass",
				},
			},
		},
		{
			name: "return error when password is empty",
			args: args{
				password: "",
			},
			want: want{
				err: errors.NewErrInvalidOption("password", ""),
				obj: &T{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithPassword(test.args.password)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithAuthProvider(t *testing.T) {
	type T = client
	type args struct {
		authProvider func(h *gocql.HostInfo) (gocql.Authenticator, error)
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, clientComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		func() test {
			ap := func(h *gocql.HostInfo) (gocql.Authenticator, error) {
				return nil, nil
			}
			return test{
				name: "set auth provider success",
				args: args{
					authProvider: ap,
				},
				want: want{
					obj: &T{
						authProvider: ap,
					},
				},
			}
		}(),
		{
			name: "return error when auth provider is nil",
			args: args{},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("authProvider", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithAuthProvider(test.args.authProvider)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithRetryPolicyNumRetries(t *testing.T) {
	type T = client
	type args struct {
		n int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set numRetries success",
			args: args{
				n: 4,
			},
			want: want{
				obj: &T{
					retryPolicy: retryPolicy{
						numRetries: 4,
					},
				},
			},
		},
		{
			name: "return error when numRetries < 0",
			args: args{
				n: -1,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("retryPolicyNumRetries", -1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithRetryPolicyNumRetries(test.args.n)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithRetryPolicyMinDuration(t *testing.T) {
	type T = client
	type args struct {
		minDuration string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set dur success",
			args: args{
				minDuration: "5s",
			},
			want: want{
				obj: &T{
					retryPolicy: retryPolicy{
						minDuration: 5 * time.Second,
					},
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				minDuration: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption(
					"retryPolicyMinDuration",
					"dummy",
					errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\""),
				),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				minDuration: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("retryPolicyMinDuration", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithRetryPolicyMinDuration(test.args.minDuration)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithRetryPolicyMaxDuration(t *testing.T) {
	type T = client
	type args struct {
		maxDuration string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set dur success",
			args: args{
				maxDuration: "5s",
			},
			want: want{
				obj: &T{
					retryPolicy: retryPolicy{
						maxDuration: 5 * time.Second,
					},
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				maxDuration: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption(
					"retryPolicyMaxDuration",
					"dummy",
					errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\""),
				),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				maxDuration: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("retryPolicyMaxDuration", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithRetryPolicyMaxDuration(test.args.maxDuration)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReconnectionPolicyInitialInterval(t *testing.T) {
	type T = client
	type args struct {
		initialInterval string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set interval success",
			args: args{
				initialInterval: "5s",
			},
			want: want{
				obj: &T{
					reconnectionPolicy: reconnectionPolicy{
						initialInterval: 5 * time.Second,
					},
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				initialInterval: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption(
					"reconnectionPolicyInitialInterval",
					"dummy",
					errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\""),
				),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				initialInterval: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("reconnectionPolicyInitialInterval", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithReconnectionPolicyInitialInterval(test.args.initialInterval)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReconnectionPolicyMaxRetries(t *testing.T) {
	type T = client
	type args struct {
		maxRetries int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set maxRetries success",
			args: args{
				maxRetries: 4,
			},
			want: want{
				obj: &T{
					reconnectionPolicy: reconnectionPolicy{
						maxRetries: 4,
					},
				},
			},
		},
		{
			name: "return error when maxRetries < 0",
			args: args{
				maxRetries: -1,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("maxRetries", -1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithReconnectionPolicyMaxRetries(test.args.maxRetries)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithSocketKeepalive(t *testing.T) {
	type T = client
	type args struct {
		socketKeepalive string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set socket keepalive success",
			args: args{
				socketKeepalive: "5s",
			},
			want: want{
				obj: &T{
					socketKeepalive: 5 * time.Second,
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				socketKeepalive: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption("socketKeepalive", "dummy", errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\"")),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				socketKeepalive: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("socketKeepalive", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithSocketKeepalive(test.args.socketKeepalive)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxPreparedStmts(t *testing.T) {
	type T = client
	type args struct {
		maxPreparedStmts int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set maxPreparedStmts success",
			args: args{
				maxPreparedStmts: 4,
			},
			want: want{
				obj: &T{
					maxPreparedStmts: 4,
				},
			},
		},
		{
			name: "return error when maxPreparedStmts < 0",
			args: args{
				maxPreparedStmts: -1,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("maxPreparedStmts", -1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithMaxPreparedStmts(test.args.maxPreparedStmts)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxRoutingKeyInfo(t *testing.T) {
	type T = client
	type args struct {
		maxRoutingKeyInfo int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set maxRoutingKeyInfo success",
			args: args{
				maxRoutingKeyInfo: 4,
			},
			want: want{
				obj: &T{
					maxRoutingKeyInfo: 4,
				},
			},
		},
		{
			name: "return error when maxRoutingKeyInfo < 0",
			args: args{
				maxRoutingKeyInfo: -1,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("maxRoutingKeyInfo", -1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithMaxRoutingKeyInfo(test.args.maxRoutingKeyInfo)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithPageSize(t *testing.T) {
	type T = client
	type args struct {
		pageSize int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set pageSize success",
			args: args{
				pageSize: 4,
			},
			want: want{
				obj: &T{
					pageSize: 4,
				},
			},
		},
		{
			name: "return error when pageSize < 0",
			args: args{
				pageSize: -1,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("pageSize", -1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithPageSize(test.args.pageSize)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTLS(t *testing.T) {
	type T = client
	type args struct {
		tls *tls.Config
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set compressor success",
			args: args{
				tls: &tls.Config{
					MinVersion: tls.VersionTLS13,
				},
			},
			want: want{
				obj: &T{
					tls: &tls.Config{
						MinVersion: tls.VersionTLS13,
					},
				},
			},
		},
		{
			name: "return error when compressor is nil",
			args: args{},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("tls", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithTLS(test.args.tls)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTLSCertPath(t *testing.T) {
	type T = client
	type args struct {
		certPath string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set cert path success",
			args: args{
				certPath: "cert_path",
			},
			want: want{
				obj: &T{
					tlsCertPath: "cert_path",
				},
			},
		},
		{
			name: "return error if the cert path is empty",
			args: args{
				certPath: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("tlsCertPath", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithTLSCertPath(test.args.certPath)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTLSKeyPath(t *testing.T) {
	type T = client
	type args struct {
		keyPath string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set key path success",
			args: args{
				keyPath: "key_path",
			},
			want: want{
				obj: &T{
					tlsKeyPath: "key_path",
				},
			},
		},
		{
			name: "return error if the key path is empty",
			args: args{
				keyPath: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("tlsKeyPath", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithTLSKeyPath(test.args.keyPath)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTLSCAPath(t *testing.T) {
	type T = client
	type args struct {
		caPath string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set ca path success",
			args: args{
				caPath: "ca_path",
			},
			want: want{
				obj: &T{
					tlsCAPath: "ca_path",
				},
			},
		},
		{
			name: "return error if the ca path is empty",
			args: args{
				caPath: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("tlsCAPath", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithTLSCAPath(test.args.caPath)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithEnableHostVerification(t *testing.T) {
	type T = client
	type args struct {
		enableHostVerification bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set enable host verification success",
			args: args{
				enableHostVerification: true,
			},
			want: want{
				obj: &T{
					enableHostVerification: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithEnableHostVerification(test.args.enableHostVerification)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultTimestamp(t *testing.T) {
	type T = client
	type args struct {
		defaultTimestamp bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default timestamp success",
			args: args{
				defaultTimestamp: true,
			},
			want: want{
				obj: &T{
					defaultTimestamp: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDefaultTimestamp(test.args.defaultTimestamp)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDC(t *testing.T) {
	type T = client
	type args struct {
		name string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set dc success",
			args: args{
				name: "dc",
			},
			want: want{
				obj: &T{
					poolConfig: poolConfig{
						dataCenterName: "dc",
					},
				},
			},
		},
		{
			name: "return error if the dc is empty",
			args: args{
				name: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("DC", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDC(test.args.name)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDCAwareRouting(t *testing.T) {
	type T = client
	type args struct {
		dc_aware_routing bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set DC aware routing success",
			args: args{
				dc_aware_routing: true,
			},
			want: want{
				obj: &T{
					poolConfig: poolConfig{
						enableDCAwareRouting: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDCAwareRouting(test.args.dc_aware_routing)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithNonLocalReplicasFallback(t *testing.T) {
	type T = client
	type args struct {
		non_local_replicas_fallback bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set non local replicas fallback success",
			args: args{
				non_local_replicas_fallback: true,
			},
			want: want{
				obj: &T{
					poolConfig: poolConfig{
						enableNonLocalReplicasFallback: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithNonLocalReplicasFallback(test.args.non_local_replicas_fallback)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithShuffleReplicas(t *testing.T) {
	type T = client
	type args struct {
		shuffleReplicas bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set shuffle replicas success",
			args: args{
				shuffleReplicas: true,
			},
			want: want{
				obj: &T{
					poolConfig: poolConfig{
						enableShuffleReplicas: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithShuffleReplicas(test.args.shuffleReplicas)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTokenAwareHostPolicy(t *testing.T) {
	type T = client
	type args struct {
		tokenAwareHostPolicy bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set token aware host policy success",
			args: args{
				tokenAwareHostPolicy: true,
			},
			want: want{
				obj: &T{
					poolConfig: poolConfig{
						enableTokenAwareHostPolicy: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithTokenAwareHostPolicy(test.args.tokenAwareHostPolicy)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxWaitSchemaAgreement(t *testing.T) {
	type T = client
	type args struct {
		maxWaitSchemaAgreement string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set max wait schema agreement success",
			args: args{
				maxWaitSchemaAgreement: "5s",
			},
			want: want{
				obj: &T{
					maxWaitSchemaAgreement: 5 * time.Second,
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				maxWaitSchemaAgreement: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption(
					"maxWaitSchemaAgreement",
					"dummy",
					errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\""),
				),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				maxWaitSchemaAgreement: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("maxWaitSchemaAgreement", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithMaxWaitSchemaAgreement(test.args.maxWaitSchemaAgreement)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReconnectInterval(t *testing.T) {
	type T = client
	type args struct {
		reconnectInterval string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set reconnect interval success",
			args: args{
				reconnectInterval: "5s",
			},
			want: want{
				obj: &T{
					reconnectInterval: 5 * time.Second,
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				reconnectInterval: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption("reconnectInterval", "dummy", errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\"")),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				reconnectInterval: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("reconnectInterval", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithReconnectInterval(test.args.reconnectInterval)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithIgnorePeerAddr(t *testing.T) {
	type T = client
	type args struct {
		ignorePeerAddr bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set ignore peer addr success",
			args: args{
				ignorePeerAddr: true,
			},
			want: want{
				obj: &T{
					ignorePeerAddr: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithIgnorePeerAddr(test.args.ignorePeerAddr)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDisableInitialHostLookup(t *testing.T) {
	type T = client
	type args struct {
		disableInitialHostLookup bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set disable initial host lookup success",
			args: args{
				disableInitialHostLookup: true,
			},
			want: want{
				obj: &T{
					disableInitialHostLookup: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDisableInitialHostLookup(test.args.disableInitialHostLookup)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDisableNodeStatusEvents(t *testing.T) {
	type T = client
	type args struct {
		disableNodeStatusEvents bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set disable node status events",
			args: args{
				disableNodeStatusEvents: true,
			},
			want: want{
				obj: &T{
					disableNodeStatusEvents: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDisableNodeStatusEvents(test.args.disableNodeStatusEvents)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDisableTopologyEvents(t *testing.T) {
	type T = client
	type args struct {
		disableTopologyEvents bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set disable topology events",
			args: args{
				disableTopologyEvents: true,
			},
			want: want{
				obj: &T{
					disableTopologyEvents: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDisableTopologyEvents(test.args.disableTopologyEvents)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDisableSchemaEvents(t *testing.T) {
	type T = client
	type args struct {
		disableSchemaEvents bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set disable schema events success",
			args: args{
				disableSchemaEvents: true,
			},
			want: want{
				obj: &T{
					disableSchemaEvents: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDisableSchemaEvents(test.args.disableSchemaEvents)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDisableSkipMetadata(t *testing.T) {
	type T = client
	type args struct {
		disableSkipMetadata bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set disable skip metadata success",
			args: args{
				disableSkipMetadata: true,
			},
			want: want{
				obj: &T{
					disableSkipMetadata: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDisableSkipMetadata(test.args.disableSkipMetadata)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithQueryObserver(t *testing.T) {
	type T = client
	type args struct {
		obs QueryObserver
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set query observer success",
			args: args{
				obs: queryObserverImpl{},
			},
			want: want{
				obj: &T{
					queryObserver: queryObserverImpl{},
				},
			},
		},
		{
			name: "set nil query observer fail",
			args: args{
				obs: nil,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("queryObserver", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithQueryObserver(test.args.obs)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBatchObserver(t *testing.T) {
	type T = client
	type args struct {
		obs BatchObserver
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set batch observer success",
			args: args{
				obs: batchObserverImpl{},
			},
			want: want{
				obj: &T{
					batchObserver: batchObserverImpl{},
				},
			},
		},
		{
			name: "set nil batch observer fail",
			args: args{
				obs: nil,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("batchObserver", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithBatchObserver(test.args.obs)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithConnectObserver(t *testing.T) {
	type T = client
	type args struct {
		obs ConnectObserver
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set connect observer success",
			args: args{
				obs: connectObserverImpl{},
			},
			want: want{
				obj: &T{
					connectObserver: connectObserverImpl{},
				},
			},
		},
		{
			name: "set nil connect observer fail",
			args: args{
				obs: nil,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("connectObserver", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithConnectObserver(test.args.obs)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithFrameHeaderObserver(t *testing.T) {
	type T = client
	type args struct {
		obs FrameHeaderObserver
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set frame header observer success",
			args: args{
				obs: frameHeaderObserverImpl{},
			},
			want: want{
				obj: &T{
					frameHeaderObserver: frameHeaderObserverImpl{},
				},
			},
		},
		{
			name: "set nil frame header observer fail",
			args: args{
				obs: nil,
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("frameHeaderObserver", nil),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithFrameHeaderObserver(test.args.obs)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultIdempotence(t *testing.T) {
	type T = client
	type args struct {
		defaultIdempotence bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default idempotent success",
			args: args{
				defaultIdempotence: true,
			},
			want: want{
				obj: &T{
					defaultIdempotence: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDefaultIdempotence(test.args.defaultIdempotence)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithWriteCoalesceWaitTime(t *testing.T) {
	type T = client
	type args struct {
		writeCoalesceWaitTime string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set write coalesce wait time success",
			args: args{
				writeCoalesceWaitTime: "5s",
			},
			want: want{
				obj: &T{
					writeCoalesceWaitTime: 5 * time.Second,
				},
			},
		},
		{
			name: "return error if the time format is invalid",
			args: args{
				writeCoalesceWaitTime: "dummy",
			},
			want: want{
				err: errors.NewErrCriticalOption("writeCoalesceWaitTime", "dummy", errors.New("invalid timeout value: dummy	:timeout parse error out put failed: time: invalid duration \"dummy\"")),
				obj: &T{},
			},
		},
		{
			name: "return error if the time is empty",
			args: args{
				writeCoalesceWaitTime: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("writeCoalesceWaitTime", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithWriteCoalesceWaitTime(test.args.writeCoalesceWaitTime)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithHostFilter(t *testing.T) {
	type T = client
	type args struct {
		flg bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set host filter flag success",
			args: args{
				flg: true,
			},
			want: want{
				obj: &T{
					hostFilter: hostFilter{
						enable: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithHostFilter(test.args.flg)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDCHostFilter(t *testing.T) {
	type T = client
	type args struct {
		dc string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set DC host filter success",
			args: args{
				dc: "DC",
			},
			want: want{
				obj: &T{
					hostFilter: hostFilter{
						dcHost: "DC",
						enable: true,
					},
				},
			},
		},
		{
			name: "return error when  DC host filter empty",
			args: args{
				dc: "",
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("dcHostFilter", ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithDCHostFilter(test.args.dc)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithWhiteListHostFilter(t *testing.T) {
	type T = client
	type args struct {
		list []string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set white list host success",
			args: args{
				list: []string{"list"},
			},
			want: want{
				obj: &T{
					hostFilter: hostFilter{
						whiteList: []string{"list"},
						enable:    true,
					},
				},
			},
		},
		{
			name: "return error when white list host empty",
			args: args{
				list: []string{},
			},
			want: want{
				obj: &T{},
				err: errors.NewErrInvalidOption("whiteListHostFilter", []string{}),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := WithWhiteListHostFilter(test.args.list)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
