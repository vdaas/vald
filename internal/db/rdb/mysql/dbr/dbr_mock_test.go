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

package dbr

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestMockDBR_Open(t *testing.T) {
	t.Parallel()
	type args struct {
		driver string
		dsn    string
		log    EventReceiver
	}
	type fields struct {
		OpenFunc func(driver, dsn string, log EventReceiver) (Connection, error)
		EqFunc   func(col string, val interface{}) Builder
	}
	type want struct {
		want Connection
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Connection, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Connection, err error) error {
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
		           driver: "",
		           dsn: "",
		           log: nil,
		       },
		       fields: fields {
		           OpenFunc: nil,
		           EqFunc: nil,
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
		           driver: "",
		           dsn: "",
		           log: nil,
		           },
		           fields: fields {
		           OpenFunc: nil,
		           EqFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			d := &MockDBR{
				OpenFunc: test.fields.OpenFunc,
				EqFunc:   test.fields.EqFunc,
			}

			got, err := d.Open(test.args.driver, test.args.dsn, test.args.log)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockDBR_Eq(t *testing.T) {
	t.Parallel()
	type args struct {
		col string
		val interface{}
	}
	type fields struct {
		OpenFunc func(driver, dsn string, log EventReceiver) (Connection, error)
		EqFunc   func(col string, val interface{}) Builder
	}
	type want struct {
		want Builder
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Builder) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Builder) error {
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
		           col: "",
		           val: nil,
		       },
		       fields: fields {
		           OpenFunc: nil,
		           EqFunc: nil,
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
		           col: "",
		           val: nil,
		           },
		           fields: fields {
		           OpenFunc: nil,
		           EqFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			d := &MockDBR{
				OpenFunc: test.fields.OpenFunc,
				EqFunc:   test.fields.EqFunc,
			}

			got := d.Eq(test.args.col, test.args.val)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSession_Select(t *testing.T) {
	t.Parallel()
	type args struct {
		column []string
	}
	type fields struct {
		SelectFunc      func(column ...string) SelectStmt
		BeginFunc       func() (Tx, error)
		CloseFunc       func() error
		PingContextFunc func(ctx context.Context) error
	}
	type want struct {
		want SelectStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, SelectStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got SelectStmt) error {
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
		           column: nil,
		       },
		       fields: fields {
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
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
		           column: nil,
		           },
		           fields: fields {
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockSession{
				SelectFunc:      test.fields.SelectFunc,
				BeginFunc:       test.fields.BeginFunc,
				CloseFunc:       test.fields.CloseFunc,
				PingContextFunc: test.fields.PingContextFunc,
			}

			got := s.Select(test.args.column...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSession_Begin(t *testing.T) {
	t.Parallel()
	type fields struct {
		SelectFunc      func(column ...string) SelectStmt
		BeginFunc       func() (Tx, error)
		CloseFunc       func() error
		PingContextFunc func(ctx context.Context) error
	}
	type want struct {
		want Tx
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, Tx, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Tx, err error) error {
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
		       fields: fields {
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
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
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
		           },
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
			s := &MockSession{
				SelectFunc:      test.fields.SelectFunc,
				BeginFunc:       test.fields.BeginFunc,
				CloseFunc:       test.fields.CloseFunc,
				PingContextFunc: test.fields.PingContextFunc,
			}

			got, err := s.Begin()
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSession_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		SelectFunc      func(column ...string) SelectStmt
		BeginFunc       func() (Tx, error)
		CloseFunc       func() error
		PingContextFunc func(ctx context.Context) error
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
		       fields: fields {
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
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
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
		           },
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
			s := &MockSession{
				SelectFunc:      test.fields.SelectFunc,
				BeginFunc:       test.fields.BeginFunc,
				CloseFunc:       test.fields.CloseFunc,
				PingContextFunc: test.fields.PingContextFunc,
			}

			err := s.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSession_PingContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		SelectFunc      func(column ...string) SelectStmt
		BeginFunc       func() (Tx, error)
		CloseFunc       func() error
		PingContextFunc func(ctx context.Context) error
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
		       },
		       fields: fields {
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
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
		           SelectFunc: nil,
		           BeginFunc: nil,
		           CloseFunc: nil,
		           PingContextFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockSession{
				SelectFunc:      test.fields.SelectFunc,
				BeginFunc:       test.fields.BeginFunc,
				CloseFunc:       test.fields.CloseFunc,
				PingContextFunc: test.fields.PingContextFunc,
			}

			err := s.PingContext(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockTx_Commit(t *testing.T) {
	t.Parallel()
	type fields struct {
		CommitFunc                  func() error
		RollbackFunc                func() error
		RollbackUnlessCommittedFunc func()
		InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
		InsertIntoFunc              func(table string) InsertStmt
		SelectFunc                  func(column ...string) SelectStmt
		DeleteFromFunc              func(table string) DeleteStmt
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
		       fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
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
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
		           },
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
			t := &MockTx{
				CommitFunc:                  test.fields.CommitFunc,
				RollbackFunc:                test.fields.RollbackFunc,
				RollbackUnlessCommittedFunc: test.fields.RollbackUnlessCommittedFunc,
				InsertBySqlFunc:             test.fields.InsertBySqlFunc,
				InsertIntoFunc:              test.fields.InsertIntoFunc,
				SelectFunc:                  test.fields.SelectFunc,
				DeleteFromFunc:              test.fields.DeleteFromFunc,
			}

			err := t.Commit()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockTx_Rollback(t *testing.T) {
	t.Parallel()
	type fields struct {
		CommitFunc                  func() error
		RollbackFunc                func() error
		RollbackUnlessCommittedFunc func()
		InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
		InsertIntoFunc              func(table string) InsertStmt
		SelectFunc                  func(column ...string) SelectStmt
		DeleteFromFunc              func(table string) DeleteStmt
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
		       fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
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
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
		           },
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
			t := &MockTx{
				CommitFunc:                  test.fields.CommitFunc,
				RollbackFunc:                test.fields.RollbackFunc,
				RollbackUnlessCommittedFunc: test.fields.RollbackUnlessCommittedFunc,
				InsertBySqlFunc:             test.fields.InsertBySqlFunc,
				InsertIntoFunc:              test.fields.InsertIntoFunc,
				SelectFunc:                  test.fields.SelectFunc,
				DeleteFromFunc:              test.fields.DeleteFromFunc,
			}

			err := t.Rollback()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockTx_RollbackUnlessCommitted(t *testing.T) {
	t.Parallel()
	type fields struct {
		CommitFunc                  func() error
		RollbackFunc                func() error
		RollbackUnlessCommittedFunc func()
		InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
		InsertIntoFunc              func(table string) InsertStmt
		SelectFunc                  func(column ...string) SelectStmt
		DeleteFromFunc              func(table string) DeleteStmt
	}
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
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
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
		           },
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
			t := &MockTx{
				CommitFunc:                  test.fields.CommitFunc,
				RollbackFunc:                test.fields.RollbackFunc,
				RollbackUnlessCommittedFunc: test.fields.RollbackUnlessCommittedFunc,
				InsertBySqlFunc:             test.fields.InsertBySqlFunc,
				InsertIntoFunc:              test.fields.InsertIntoFunc,
				SelectFunc:                  test.fields.SelectFunc,
				DeleteFromFunc:              test.fields.DeleteFromFunc,
			}

			t.RollbackUnlessCommitted()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockTx_InsertBySql(t *testing.T) {
	t.Parallel()
	type args struct {
		query string
		value []interface{}
	}
	type fields struct {
		CommitFunc                  func() error
		RollbackFunc                func() error
		RollbackUnlessCommittedFunc func()
		InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
		InsertIntoFunc              func(table string) InsertStmt
		SelectFunc                  func(column ...string) SelectStmt
		DeleteFromFunc              func(table string) DeleteStmt
	}
	type want struct {
		want InsertStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, InsertStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got InsertStmt) error {
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
		           query: "",
		           value: nil,
		       },
		       fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
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
		           query: "",
		           value: nil,
		           },
		           fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			t := &MockTx{
				CommitFunc:                  test.fields.CommitFunc,
				RollbackFunc:                test.fields.RollbackFunc,
				RollbackUnlessCommittedFunc: test.fields.RollbackUnlessCommittedFunc,
				InsertBySqlFunc:             test.fields.InsertBySqlFunc,
				InsertIntoFunc:              test.fields.InsertIntoFunc,
				SelectFunc:                  test.fields.SelectFunc,
				DeleteFromFunc:              test.fields.DeleteFromFunc,
			}

			got := t.InsertBySql(test.args.query, test.args.value...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockTx_InsertInto(t *testing.T) {
	t.Parallel()
	type args struct {
		table string
	}
	type fields struct {
		CommitFunc                  func() error
		RollbackFunc                func() error
		RollbackUnlessCommittedFunc func()
		InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
		InsertIntoFunc              func(table string) InsertStmt
		SelectFunc                  func(column ...string) SelectStmt
		DeleteFromFunc              func(table string) DeleteStmt
	}
	type want struct {
		want InsertStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, InsertStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got InsertStmt) error {
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
		           table: "",
		       },
		       fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
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
		           table: "",
		           },
		           fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			t := &MockTx{
				CommitFunc:                  test.fields.CommitFunc,
				RollbackFunc:                test.fields.RollbackFunc,
				RollbackUnlessCommittedFunc: test.fields.RollbackUnlessCommittedFunc,
				InsertBySqlFunc:             test.fields.InsertBySqlFunc,
				InsertIntoFunc:              test.fields.InsertIntoFunc,
				SelectFunc:                  test.fields.SelectFunc,
				DeleteFromFunc:              test.fields.DeleteFromFunc,
			}

			got := t.InsertInto(test.args.table)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockTx_Select(t *testing.T) {
	t.Parallel()
	type args struct {
		column []string
	}
	type fields struct {
		CommitFunc                  func() error
		RollbackFunc                func() error
		RollbackUnlessCommittedFunc func()
		InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
		InsertIntoFunc              func(table string) InsertStmt
		SelectFunc                  func(column ...string) SelectStmt
		DeleteFromFunc              func(table string) DeleteStmt
	}
	type want struct {
		want SelectStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, SelectStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got SelectStmt) error {
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
		           column: nil,
		       },
		       fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
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
		           column: nil,
		           },
		           fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			t := &MockTx{
				CommitFunc:                  test.fields.CommitFunc,
				RollbackFunc:                test.fields.RollbackFunc,
				RollbackUnlessCommittedFunc: test.fields.RollbackUnlessCommittedFunc,
				InsertBySqlFunc:             test.fields.InsertBySqlFunc,
				InsertIntoFunc:              test.fields.InsertIntoFunc,
				SelectFunc:                  test.fields.SelectFunc,
				DeleteFromFunc:              test.fields.DeleteFromFunc,
			}

			got := t.Select(test.args.column...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockTx_DeleteFrom(t *testing.T) {
	t.Parallel()
	type args struct {
		table string
	}
	type fields struct {
		CommitFunc                  func() error
		RollbackFunc                func() error
		RollbackUnlessCommittedFunc func()
		InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
		InsertIntoFunc              func(table string) InsertStmt
		SelectFunc                  func(column ...string) SelectStmt
		DeleteFromFunc              func(table string) DeleteStmt
	}
	type want struct {
		want DeleteStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, DeleteStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got DeleteStmt) error {
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
		           table: "",
		       },
		       fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
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
		           table: "",
		           },
		           fields: fields {
		           CommitFunc: nil,
		           RollbackFunc: nil,
		           RollbackUnlessCommittedFunc: nil,
		           InsertBySqlFunc: nil,
		           InsertIntoFunc: nil,
		           SelectFunc: nil,
		           DeleteFromFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			t := &MockTx{
				CommitFunc:                  test.fields.CommitFunc,
				RollbackFunc:                test.fields.RollbackFunc,
				RollbackUnlessCommittedFunc: test.fields.RollbackUnlessCommittedFunc,
				InsertBySqlFunc:             test.fields.InsertBySqlFunc,
				InsertIntoFunc:              test.fields.InsertIntoFunc,
				SelectFunc:                  test.fields.SelectFunc,
				DeleteFromFunc:              test.fields.DeleteFromFunc,
			}

			got := t.DeleteFrom(test.args.table)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockConn_NewSession(t *testing.T) {
	t.Parallel()
	type args struct {
		event EventReceiver
	}
	type fields struct {
		NewSessionFunc         func(event EventReceiver) Session
		SetConnMaxLifetimeFunc func(d time.Duration)
		SetMaxIdleConnsFunc    func(n int)
		SetMaxOpenConnsFunc    func(n int)
	}
	type want struct {
		want Session
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Session) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Session) error {
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
		           event: nil,
		       },
		       fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
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
		           event: nil,
		           },
		           fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &MockConn{
				NewSessionFunc:         test.fields.NewSessionFunc,
				SetConnMaxLifetimeFunc: test.fields.SetConnMaxLifetimeFunc,
				SetMaxIdleConnsFunc:    test.fields.SetMaxIdleConnsFunc,
				SetMaxOpenConnsFunc:    test.fields.SetMaxOpenConnsFunc,
			}

			got := c.NewSession(test.args.event)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockConn_SetConnMaxLifetime(t *testing.T) {
	t.Parallel()
	type args struct {
		d time.Duration
	}
	type fields struct {
		NewSessionFunc         func(event EventReceiver) Session
		SetConnMaxLifetimeFunc func(d time.Duration)
		SetMaxIdleConnsFunc    func(n int)
		SetMaxOpenConnsFunc    func(n int)
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           d: nil,
		       },
		       fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
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
		           d: nil,
		           },
		           fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &MockConn{
				NewSessionFunc:         test.fields.NewSessionFunc,
				SetConnMaxLifetimeFunc: test.fields.SetConnMaxLifetimeFunc,
				SetMaxIdleConnsFunc:    test.fields.SetMaxIdleConnsFunc,
				SetMaxOpenConnsFunc:    test.fields.SetMaxOpenConnsFunc,
			}

			c.SetConnMaxLifetime(test.args.d)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockConn_SetMaxIdleConns(t *testing.T) {
	t.Parallel()
	type args struct {
		n int
	}
	type fields struct {
		NewSessionFunc         func(event EventReceiver) Session
		SetConnMaxLifetimeFunc func(d time.Duration)
		SetMaxIdleConnsFunc    func(n int)
		SetMaxOpenConnsFunc    func(n int)
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           n: 0,
		       },
		       fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
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
		           n: 0,
		           },
		           fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &MockConn{
				NewSessionFunc:         test.fields.NewSessionFunc,
				SetConnMaxLifetimeFunc: test.fields.SetConnMaxLifetimeFunc,
				SetMaxIdleConnsFunc:    test.fields.SetMaxIdleConnsFunc,
				SetMaxOpenConnsFunc:    test.fields.SetMaxOpenConnsFunc,
			}

			c.SetMaxIdleConns(test.args.n)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockConn_SetMaxOpenConns(t *testing.T) {
	t.Parallel()
	type args struct {
		n int
	}
	type fields struct {
		NewSessionFunc         func(event EventReceiver) Session
		SetConnMaxLifetimeFunc func(d time.Duration)
		SetMaxIdleConnsFunc    func(n int)
		SetMaxOpenConnsFunc    func(n int)
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           n: 0,
		       },
		       fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
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
		           n: 0,
		           },
		           fields: fields {
		           NewSessionFunc: nil,
		           SetConnMaxLifetimeFunc: nil,
		           SetMaxIdleConnsFunc: nil,
		           SetMaxOpenConnsFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &MockConn{
				NewSessionFunc:         test.fields.NewSessionFunc,
				SetConnMaxLifetimeFunc: test.fields.SetConnMaxLifetimeFunc,
				SetMaxIdleConnsFunc:    test.fields.SetMaxIdleConnsFunc,
				SetMaxOpenConnsFunc:    test.fields.SetMaxOpenConnsFunc,
			}

			c.SetMaxOpenConns(test.args.n)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSelect_From(t *testing.T) {
	t.Parallel()
	type args struct {
		table interface{}
	}
	type fields struct {
		FromFunc        func(table interface{}) SelectStmt
		WhereFunc       func(query interface{}, value ...interface{}) SelectStmt
		LimitFunc       func(n uint64) SelectStmt
		LoadContextFunc func(ctx context.Context, value interface{}) (int, error)
	}
	type want struct {
		want SelectStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, SelectStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got SelectStmt) error {
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
		           table: nil,
		       },
		       fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
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
		           table: nil,
		           },
		           fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockSelect{
				FromFunc:        test.fields.FromFunc,
				WhereFunc:       test.fields.WhereFunc,
				LimitFunc:       test.fields.LimitFunc,
				LoadContextFunc: test.fields.LoadContextFunc,
			}

			got := s.From(test.args.table)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSelect_Where(t *testing.T) {
	t.Parallel()
	type args struct {
		query interface{}
		value []interface{}
	}
	type fields struct {
		FromFunc        func(table interface{}) SelectStmt
		WhereFunc       func(query interface{}, value ...interface{}) SelectStmt
		LimitFunc       func(n uint64) SelectStmt
		LoadContextFunc func(ctx context.Context, value interface{}) (int, error)
	}
	type want struct {
		want SelectStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, SelectStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got SelectStmt) error {
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
		           query: nil,
		           value: nil,
		       },
		       fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
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
		           query: nil,
		           value: nil,
		           },
		           fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockSelect{
				FromFunc:        test.fields.FromFunc,
				WhereFunc:       test.fields.WhereFunc,
				LimitFunc:       test.fields.LimitFunc,
				LoadContextFunc: test.fields.LoadContextFunc,
			}

			got := s.Where(test.args.query, test.args.value...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSelect_Limit(t *testing.T) {
	t.Parallel()
	type args struct {
		n uint64
	}
	type fields struct {
		FromFunc        func(table interface{}) SelectStmt
		WhereFunc       func(query interface{}, value ...interface{}) SelectStmt
		LimitFunc       func(n uint64) SelectStmt
		LoadContextFunc func(ctx context.Context, value interface{}) (int, error)
	}
	type want struct {
		want SelectStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, SelectStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got SelectStmt) error {
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
		           n: 0,
		       },
		       fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
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
		           n: 0,
		           },
		           fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockSelect{
				FromFunc:        test.fields.FromFunc,
				WhereFunc:       test.fields.WhereFunc,
				LimitFunc:       test.fields.LimitFunc,
				LoadContextFunc: test.fields.LoadContextFunc,
			}

			got := s.Limit(test.args.n)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockSelect_LoadContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		value interface{}
	}
	type fields struct {
		FromFunc        func(table interface{}) SelectStmt
		WhereFunc       func(query interface{}, value ...interface{}) SelectStmt
		LimitFunc       func(n uint64) SelectStmt
		LoadContextFunc func(ctx context.Context, value interface{}) (int, error)
	}
	type want struct {
		want int
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got int, err error) error {
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
		           value: nil,
		       },
		       fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
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
		           value: nil,
		           },
		           fields: fields {
		           FromFunc: nil,
		           WhereFunc: nil,
		           LimitFunc: nil,
		           LoadContextFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockSelect{
				FromFunc:        test.fields.FromFunc,
				WhereFunc:       test.fields.WhereFunc,
				LimitFunc:       test.fields.LimitFunc,
				LoadContextFunc: test.fields.LoadContextFunc,
			}

			got, err := s.LoadContext(test.args.ctx, test.args.value)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockInsert_Columns(t *testing.T) {
	t.Parallel()
	type args struct {
		column []string
	}
	type fields struct {
		ColumnsFunc     func(column ...string) InsertStmt
		ExecContextFunc func(ctx context.Context) (sql.Result, error)
		RecordFunc      func(structValue interface{}) InsertStmt
	}
	type want struct {
		want InsertStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, InsertStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got InsertStmt) error {
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
		           column: nil,
		       },
		       fields: fields {
		           ColumnsFunc: nil,
		           ExecContextFunc: nil,
		           RecordFunc: nil,
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
		           column: nil,
		           },
		           fields: fields {
		           ColumnsFunc: nil,
		           ExecContextFunc: nil,
		           RecordFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockInsert{
				ColumnsFunc:     test.fields.ColumnsFunc,
				ExecContextFunc: test.fields.ExecContextFunc,
				RecordFunc:      test.fields.RecordFunc,
			}

			got := s.Columns(test.args.column...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockInsert_ExecContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		ColumnsFunc     func(column ...string) InsertStmt
		ExecContextFunc func(ctx context.Context) (sql.Result, error)
		RecordFunc      func(structValue interface{}) InsertStmt
	}
	type want struct {
		want sql.Result
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, sql.Result, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got sql.Result, err error) error {
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
		       },
		       fields: fields {
		           ColumnsFunc: nil,
		           ExecContextFunc: nil,
		           RecordFunc: nil,
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
		           ColumnsFunc: nil,
		           ExecContextFunc: nil,
		           RecordFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockInsert{
				ColumnsFunc:     test.fields.ColumnsFunc,
				ExecContextFunc: test.fields.ExecContextFunc,
				RecordFunc:      test.fields.RecordFunc,
			}

			got, err := s.ExecContext(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockInsert_Record(t *testing.T) {
	t.Parallel()
	type args struct {
		structValue interface{}
	}
	type fields struct {
		ColumnsFunc     func(column ...string) InsertStmt
		ExecContextFunc func(ctx context.Context) (sql.Result, error)
		RecordFunc      func(structValue interface{}) InsertStmt
	}
	type want struct {
		want InsertStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, InsertStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got InsertStmt) error {
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
		           structValue: nil,
		       },
		       fields: fields {
		           ColumnsFunc: nil,
		           ExecContextFunc: nil,
		           RecordFunc: nil,
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
		           structValue: nil,
		           },
		           fields: fields {
		           ColumnsFunc: nil,
		           ExecContextFunc: nil,
		           RecordFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockInsert{
				ColumnsFunc:     test.fields.ColumnsFunc,
				ExecContextFunc: test.fields.ExecContextFunc,
				RecordFunc:      test.fields.RecordFunc,
			}

			got := s.Record(test.args.structValue)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockDelete_ExecContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		ExecContextFunc func(ctx context.Context) (sql.Result, error)
		WhereFunc       func(query interface{}, value ...interface{}) DeleteStmt
	}
	type want struct {
		want sql.Result
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, sql.Result, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got sql.Result, err error) error {
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
		       },
		       fields: fields {
		           ExecContextFunc: nil,
		           WhereFunc: nil,
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
		           ExecContextFunc: nil,
		           WhereFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockDelete{
				ExecContextFunc: test.fields.ExecContextFunc,
				WhereFunc:       test.fields.WhereFunc,
			}

			got, err := s.ExecContext(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockDelete_Where(t *testing.T) {
	t.Parallel()
	type args struct {
		query interface{}
		value []interface{}
	}
	type fields struct {
		ExecContextFunc func(ctx context.Context) (sql.Result, error)
		WhereFunc       func(query interface{}, value ...interface{}) DeleteStmt
	}
	type want struct {
		want DeleteStmt
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, DeleteStmt) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got DeleteStmt) error {
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
		           query: nil,
		           value: nil,
		       },
		       fields: fields {
		           ExecContextFunc: nil,
		           WhereFunc: nil,
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
		           query: nil,
		           value: nil,
		           },
		           fields: fields {
		           ExecContextFunc: nil,
		           WhereFunc: nil,
		           },
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &MockDelete{
				ExecContextFunc: test.fields.ExecContextFunc,
				WhereFunc:       test.fields.WhereFunc,
			}

			got := s.Where(test.args.query, test.args.value...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
