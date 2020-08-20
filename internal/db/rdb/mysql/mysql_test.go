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

package mysql

import (
	"context"
	"crypto/tls"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/tcp"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want MySQL
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, MySQL, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got MySQL, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
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
		           args: args {
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.Open(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_Ping(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.Ping(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_Close(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.Close(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_GetMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		want MetaVector
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, MetaVector, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got MetaVector, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           uuid: "",
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           uuid: "",
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			got, err := m.GetMeta(test.args.ctx, test.args.uuid)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_GetIPs(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		want []string
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           uuid: "",
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           uuid: "",
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			got, err := m.GetIPs(test.args.ctx, test.args.uuid)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_validateMeta(t *testing.T) {
	type args struct {
		meta MetaVector
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           meta: nil,
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
		           args: args {
		           meta: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := validateMeta(test.args.meta)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_SetMeta(t *testing.T) {
	type args struct {
		ctx context.Context
		mv  MetaVector
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           mv: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           mv: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.SetMeta(test.args.ctx, test.args.mv)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_SetMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []MetaVector
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           metas: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           metas: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.SetMetas(test.args.ctx, test.args.metas...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_deleteMetaWithTx(t *testing.T) {
	type args struct {
		ctx  context.Context
		tx   *dbr.Tx
		uuid string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           tx: nil,
		           uuid: "",
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
		           args: args {
		           ctx: nil,
		           tx: nil,
		           uuid: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := deleteMetaWithTx(test.args.ctx, test.args.tx, test.args.uuid)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_DeleteMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           uuid: "",
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           uuid: "",
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.DeleteMeta(test.args.ctx, test.args.uuid)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_DeleteMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           uuids: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           uuids: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.DeleteMetas(test.args.ctx, test.args.uuids...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_SetIPs(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
		ips  []string
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           uuid: "",
		           ips: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           uuid: "",
		           ips: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.SetIPs(test.args.ctx, test.args.uuid, test.args.ips...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_RemoveIPs(t *testing.T) {
	type args struct {
		ctx context.Context
		ips []string
	}
	type fields struct {
		db                   string
		host                 string
		port                 int
		user                 string
		pass                 string
		name                 string
		charset              string
		timezone             string
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		connMaxLifeTime      time.Duration
		dialer               tcp.Dialer
		dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
		tlsConfig            *tls.Config
		maxOpenConns         int
		maxIdleConns         int
		session              *dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           ips: nil,
		       },
		       fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
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
		           args: args {
		           ctx: nil,
		           ips: nil,
		           },
		           fields: fields {
		           db: "",
		           host: "",
		           port: 0,
		           user: "",
		           pass: "",
		           name: "",
		           charset: "",
		           timezone: "",
		           initialPingTimeLimit: nil,
		           initialPingDuration: nil,
		           connMaxLifeTime: nil,
		           dialer: nil,
		           dialerFunc: nil,
		           tlsConfig: nil,
		           maxOpenConns: 0,
		           maxIdleConns: 0,
		           session: nil,
		           connected: nil,
		           eventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &mySQLClient{
				db:                   test.fields.db,
				host:                 test.fields.host,
				port:                 test.fields.port,
				user:                 test.fields.user,
				pass:                 test.fields.pass,
				name:                 test.fields.name,
				charset:              test.fields.charset,
				timezone:             test.fields.timezone,
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				connMaxLifeTime:      test.fields.connMaxLifeTime,
				dialer:               test.fields.dialer,
				dialerFunc:           test.fields.dialerFunc,
				tlsConfig:            test.fields.tlsConfig,
				maxOpenConns:         test.fields.maxOpenConns,
				maxIdleConns:         test.fields.maxIdleConns,
				session:              test.fields.session,
				connected:            test.fields.connected,
				eventReceiver:        test.fields.eventReceiver,
			}

			err := m.RemoveIPs(test.args.ctx, test.args.ips...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
