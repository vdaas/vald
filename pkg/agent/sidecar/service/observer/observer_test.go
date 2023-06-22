//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package observer provides storage observer
package observer

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file/watch"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

// NOT IMPLEMENTED BELOW

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantSo StorageObserver
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, StorageObserver, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotSo StorageObserver, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotSo, w.wantSo) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSo, w.wantSo)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           opts:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotSo, err := New(test.args.opts...)
			if err := checkFunc(test.want, gotSo, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
		           ctx:nil,
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			got, err := o.Start(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_PostStop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			err := o.PostStop(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_startTicker(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
		           ctx:nil,
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			got, err := o.startTicker(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_startBackupLoop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
		           ctx:nil,
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			got, err := o.startBackupLoop(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_onWrite(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           name:"",
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           name:"",
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			err := o.onWrite(test.args.ctx, test.args.name)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_onCreate(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           name:"",
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           name:"",
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			err := o.onCreate(test.args.ctx, test.args.name)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_isValidMetadata(t *testing.T) {
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
	}
	type want struct {
		want bool
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got bool, err error) error {
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
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			got, err := o.isValidMetadata()
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_terminate(t *testing.T) {
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			err := o.terminate()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_requestBackup(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           in0:nil,
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           in0:nil,
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			err := o.requestBackup(test.args.in0)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_observer_backup(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		w               watch.Watcher
		dir             string
		eg              errgroup.Group
		checkDuration   time.Duration
		metadataPath    string
		postStopTimeout time.Duration
		watchEnabled    bool
		tickerEnabled   bool
		storage         storage.Storage
		ch              chan struct{}
		hooks           []Hook
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           w:nil,
		           dir:"",
		           eg:nil,
		           checkDuration:nil,
		           metadataPath:"",
		           postStopTimeout:nil,
		           watchEnabled:false,
		           tickerEnabled:false,
		           storage:nil,
		           ch:nil,
		           hooks:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &observer{
				w:               test.fields.w,
				dir:             test.fields.dir,
				eg:              test.fields.eg,
				checkDuration:   test.fields.checkDuration,
				metadataPath:    test.fields.metadataPath,
				postStopTimeout: test.fields.postStopTimeout,
				watchEnabled:    test.fields.watchEnabled,
				tickerEnabled:   test.fields.tickerEnabled,
				storage:         test.fields.storage,
				ch:              test.fields.ch,
				hooks:           test.fields.hooks,
			}

			err := o.backup(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
