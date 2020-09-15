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
	"os"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/db/rdb/mysql/dbr"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

var (
	mysqlComparatorOptions = []comparator.Option{
		comparator.AllowUnexported(mySQLClient{}),
		comparator.Comparer(func(x, y dbr.DBR) bool {
			return reflect.DeepEqual(x, y)
		}),
	}
)

func TestMain(m *testing.M) {
	log.Init()
	os.Exit(m.Run())
}

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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			m := new(mySQLClient)
			for _, opt := range defaultOpts {
				opt(m)
			}
			m.dbr = dbr.New()
			return test{
				name: "return (MySQL, nil) when opts is nil",
				want: want{
					want: m,
				},
			}
		}(),
		func() test {
			n := "vald"
			opts := []Option{
				WithName(n),
			}
			m := new(mySQLClient)
			for _, opt := range append(defaultOpts, opts...) {
				opt(m)
			}
			m.dbr = dbr.New()
			return test{
				name: "return (MySQL, nil) when opts is not empty",
				args: args{
					opts: opts,
				},
				want: want{
					want: m,
				},
			}
		}(),
		func() test {
			m := new(mySQLClient)
			for _, opt := range defaultOpts {
				opt(m)
			}
			m.dbr = dbr.New()
			err := errors.New("error")
			return test{
				name: "return (MySQL, error) when opts fails",
				args: args{
					opts: []Option{
						func(*mySQLClient) error {
							return err
						},
					},
				},
				want: want{
					err: err,
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
		session              dbr.Session
		connected            atomic.Value
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "Open success with tls config when no error occurs",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					db:                   "vdaas",
					host:                 "vald.com",
					port:                 3306,
					user:                 "vdaas",
					pass:                 "vald",
					name:                 "vald-user",
					charset:              "utf8bm4j",
					timezone:             "Local",
					initialPingTimeLimit: 1000 * time.Microsecond,
					initialPingDuration:  10 * time.Microsecond,
					connMaxLifeTime:      1 * time.Microsecond,
					tlsConfig:            new(tls.Config),
					maxOpenConns:         100,
					maxIdleConns:         100,
					session: &dbr.MockSession{
						PingContextFunc: func(ctx context.Context) error {
							return nil
						},
					},
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
					dbr: &dbr.MockDBR{
						OpenFunc: func(driver, dsn string, log EventReceiver) (dbr.Connection, error) {
							conn := &dbr.MockConn{
								NewSessionFunc: func(event EventReceiver) dbr.Session {
									return &dbr.MockSession{
										PingContextFunc: func(ctx context.Context) error {
											return nil
										},
									}
								},
								SetConnMaxLifetimeFunc: func(d time.Duration) {},
								SetMaxIdleConnsFunc:    func(n int) {},
								SetMaxOpenConnsFunc:    func(n int) {},
							}
							return conn, nil
						},
					},
				},
				want: want{},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			dialer, _ := tcp.NewDialer()
			dialerFunc := dialer.GetDialer()
			return test{
				name: "Open success with dialer when no error occurs",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					db:                   "vdaas",
					host:                 "vald.com",
					port:                 3306,
					user:                 "vdaas",
					pass:                 "vald",
					name:                 "vald-user",
					charset:              "utf8bm4j",
					timezone:             "Local",
					initialPingTimeLimit: 1000 * time.Microsecond,
					initialPingDuration:  10 * time.Microsecond,
					connMaxLifeTime:      1 * time.Microsecond,
					dialer:               dialer,
					dialerFunc:           dialerFunc,
					maxOpenConns:         100,
					maxIdleConns:         100,
					session: &dbr.MockSession{
						PingContextFunc: func(ctx context.Context) error {
							return nil
						},
					},
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
					dbr: &dbr.MockDBR{
						OpenFunc: func(driver, dsn string, log EventReceiver) (dbr.Connection, error) {
							conn := &dbr.MockConn{
								NewSessionFunc: func(event EventReceiver) dbr.Session {
									return &dbr.MockSession{
										PingContextFunc: func(ctx context.Context) error {
											return nil
										},
									}
								},
								SetConnMaxLifetimeFunc: func(d time.Duration) {},
								SetMaxIdleConnsFunc:    func(n int) {},
								SetMaxOpenConnsFunc:    func(n int) {},
							}
							return conn, nil
						},
					},
				},
				want: want{},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "returns error when dbr.Open failed",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					db:                   "vdaas",
					host:                 "vald.com",
					port:                 3306,
					user:                 "vdaas",
					pass:                 "vald",
					name:                 "vald-user",
					charset:              "utf8bm4j",
					timezone:             "Local",
					initialPingTimeLimit: 1000 * time.Microsecond,
					initialPingDuration:  10 * time.Microsecond,
					connMaxLifeTime:      1 * time.Microsecond,
					maxOpenConns:         10,
					maxIdleConns:         10,
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
					dbr: &dbr.MockDBR{
						OpenFunc: func(driver, dsn string, log EventReceiver) (dbr.Connection, error) {
							return nil, errors.ErrMySQLConnectionClosed
						},
					},
				},
				want: want{
					err: errors.ErrMySQLConnectionClosed,
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
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
				dbr:                  test.fields.dbr,
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
		initialPingTimeLimit time.Duration
		initialPingDuration  time.Duration
		session              dbr.Session
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "returns nil when no error occurs",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					initialPingTimeLimit: 1 * time.Second,
					initialPingDuration:  1 * time.Microsecond,
					session: &dbr.MockSession{
						PingContextFunc: func(ctx context.Context) error {
							return nil
						},
					},
				},
				want: want{},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			err := errors.New("error")
			return test{
				name: "returns error when session.PingContext returns error",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					initialPingTimeLimit: 100 * time.Microsecond,
					initialPingDuration:  1 * time.Microsecond,
					session: &dbr.MockSession{
						PingContextFunc: func(ctx context.Context) error {
							return err
						},
					},
				},
				want: want{
					err: errors.Wrap(errors.Wrap(errors.ErrMySQLConnectionPingFailed, err.Error()), context.DeadlineExceeded.Error()),
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "returns error when ping failed due to initialPingTimeLimit",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					initialPingTimeLimit: 1 * time.Microsecond,
					initialPingDuration:  10 * time.Microsecond,
					session: &dbr.MockSession{
						PingContextFunc: func(ctx context.Context) error {
							return nil
						},
					},
				},
				want: want{
					err: errors.ErrMySQLConnectionPingFailed,
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
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
			m := &mySQLClient{
				initialPingTimeLimit: test.fields.initialPingTimeLimit,
				initialPingDuration:  test.fields.initialPingDuration,
				session:              test.fields.session,
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
		session   dbr.Session
		connected atomic.Value
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error, *mySQLClient) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error, m *mySQLClient) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if reflect.DeepEqual(m.connected.Load().(bool), false) {
			return errors.Errorf("Close failed")
		}
		return nil
	}
	tests := []test{
		{
			name: "Close success when connection already closed",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				session: &dbr.MockSession{},
				connected: func() atomic.Value {
					v := new(atomic.Value)
					v.Store(false)
					return *v
				}(),
			},
			want: want{},
		},
		{
			name: "Close success when connection is not closed",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				session: &dbr.MockSession{
					CloseFunc: func() error {
						return nil
					},
				},
				connected: func() atomic.Value {
					v := new(atomic.Value)
					v.Store(true)
					return *v
				}(),
			},
			want: want{},
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
			m := &mySQLClient{
				session:   test.fields.session,
				connected: test.fields.connected,
			}

			err := m.Close(test.args.ctx)
			if err := test.checkFunc(test.want, err, m); err != nil {
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return (nil, error) when MySQL connection is closed",
				args: args{
					ctx:  ctx,
					uuid: "",
				},
				fields: fields{
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
				},
				want: want{
					err: errors.ErrMySQLConnectionClosed,
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			err := errors.New("loadcontext error")
			return test{
				name: "return (nil, error) when LoadContext returns error",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					session: &dbr.MockSession{
						SelectFunc: func(column ...string) dbr.SelectStmt {
							m := new(dbr.MockSelect)
							m.FromFunc = func(table interface{}) dbr.SelectStmt {
								return m
							}
							m.WhereFunc = func(query interface{}, value ...interface{}) dbr.SelectStmt {
								return m
							}
							m.LimitFunc = func(n uint64) dbr.SelectStmt {
								return m
							}
							m.LoadContextFunc = func(ctx context.Context, value interface{}) (int, error) {
								return 0, err
							}
							return m
						},
					},
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
					dbr: &dbr.MockDBR{
						EqFunc: func(col string, val interface{}) dbr.Builder {
							return dbr.New().Eq(col, val)
						},
					},
				},
				want: want{
					err: err,
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			uuid := "vdaas-01"
			return test{
				name: "return (nil, error) when meta is not found",
				args: args{
					ctx:  ctx,
					uuid: uuid,
				},
				fields: fields{
					session: &dbr.MockSession{
						SelectFunc: func(column ...string) dbr.SelectStmt {
							s := new(dbr.MockSelect)
							s.FromFunc = func(table interface{}) dbr.SelectStmt {
								return s
							}
							s.WhereFunc = func(query interface{}, value ...interface{}) dbr.SelectStmt {
								return s
							}
							s.LimitFunc = func(n uint64) dbr.SelectStmt {
								return s
							}
							s.LoadContextFunc = func(ctx context.Context, value interface{}) (int, error) {
								var mv *meta
								if reflect.TypeOf(value) == reflect.TypeOf(&mv) {
									return 1, nil
								}
								return 0, errors.New("not found")
							}
							return s
						},
					},
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
					dbr: &dbr.MockDBR{
						EqFunc: func(col string, val interface{}) dbr.Builder {
							return dbr.New().Eq(col, val)
						},
					},
				},
				want: want{
					err: errors.ErrRequiredElementNotFoundByUUID(uuid),
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			uuid := "vdaas-01"
			var m *meta
			m = &meta{
				ID:     1,
				UUID:   uuid,
				Vector: []byte("0.1,0.2"),
			}
			return test{
				name: "return (nil, error) when podIPs are not found",
				args: args{
					ctx:  ctx,
					uuid: uuid,
				},
				fields: fields{
					session: &dbr.MockSession{
						SelectFunc: func(column ...string) dbr.SelectStmt {
							s := new(dbr.MockSelect)
							s.FromFunc = func(table interface{}) dbr.SelectStmt {
								return s
							}
							s.WhereFunc = func(query interface{}, value ...interface{}) dbr.SelectStmt {
								return s
							}
							s.LimitFunc = func(n uint64) dbr.SelectStmt {
								return s
							}
							s.LoadContextFunc = func(ctx context.Context, value interface{}) (int, error) {
								var mv *meta
								var pp []podIP
								if reflect.TypeOf(value) == reflect.TypeOf(&mv) {
									mv = m
									reflect.ValueOf(value).Elem().Set(reflect.ValueOf(mv))
									return 1, nil
								} else if reflect.TypeOf(value) == reflect.TypeOf(&pp) {
									return 0, errors.New("not found")
								}
								return 0, errors.New("not found")
							}
							return s
						},
					},
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
					dbr: &dbr.MockDBR{
						EqFunc: func(col string, val interface{}) dbr.Builder {
							return dbr.New().Eq(col, val)
						},
					},
				},
				want: want{
					err: errors.New("not found"),
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			uuid := "vdaas-01"
			var m *meta
			m = &meta{
				ID:     1,
				UUID:   uuid,
				Vector: []byte("0.1,0.2"),
			}
			var p []podIP
			p = append(p, podIP{
				ID: 1,
				IP: "192.168.1.12",
			})
			return test{
				name: "return (metaVector, nil) when select success",
				args: args{
					ctx:  ctx,
					uuid: uuid,
				},
				fields: fields{
					session: &dbr.MockSession{
						SelectFunc: func(column ...string) dbr.SelectStmt {
							s := new(dbr.MockSelect)
							s.FromFunc = func(table interface{}) dbr.SelectStmt {
								return s
							}
							s.WhereFunc = func(query interface{}, value ...interface{}) dbr.SelectStmt {
								return s
							}
							s.LimitFunc = func(n uint64) dbr.SelectStmt {
								return s
							}
							s.LoadContextFunc = func(ctx context.Context, value interface{}) (int, error) {
								var mv *meta
								var pp []podIP
								if reflect.TypeOf(value) == reflect.TypeOf(&mv) {
									mv = m
									reflect.ValueOf(value).Elem().Set(reflect.ValueOf(mv))
									return 1, nil
								} else if reflect.TypeOf(value) == reflect.TypeOf(&pp) {
									pp = p
									reflect.ValueOf(value).Elem().Set(reflect.ValueOf(pp))
									return 1, nil
								}
								return 0, errors.New("error")
							}
							return s
						},
					},
					connected: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
					dbr: &dbr.MockDBR{
						EqFunc: func(col string, val interface{}) dbr.Builder {
							return dbr.New().Eq(col, val)
						},
					},
				},
				want: want{
					want: &metaVector{
						meta:   *m,
						podIPs: p,
					},
				},
				afterFunc: func(args) {
					go func() {
						cancel()
					}()
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
			m := &mySQLClient{
				session:   test.fields.session,
				connected: test.fields.connected,
				dbr:       test.fields.dbr,
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			m := new(metaVector)
			m.meta.Vector = []byte("0.1,0.2,0.9")
			return test{
				name: "return nil when the len(MetaVector) > 0",
				args: args{
					meta: m,
				},
				want:      want{},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			m := new(metaVector)
			return test{
				name: "return error when the len(MetaVector) is 0",
				args: args{
					meta: m,
				},
				want: want{
					err: errors.ErrRequiredMemberNotFilled("vector"),
				},
				checkFunc: defaultCheckFunc,
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
			}

			err := m.SetMetas(test.args.ctx, test.args.metas...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_mySQLClient_deleteMetaWithTx(t *testing.T) {
	type args struct {
		ctx  context.Context
		tx   dbr.Tx
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
			}

			err := m.deleteMetaWithTx(test.args.ctx, test.args.tx, test.args.uuid)
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
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
		session              dbr.Session
		connected            atomic.Value
		eventReceiver        EventReceiver
		dbr                  dbr.DBR
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           dbr: nil,
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
		           dbr: nil,
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
				dbr:                  test.fields.dbr,
			}

			err := m.RemoveIPs(test.args.ctx, test.args.ips...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
