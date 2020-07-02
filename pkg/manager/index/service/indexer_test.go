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

// Package service
package service

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/client/discoverer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantIdx Indexer
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Indexer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotIdx Indexer, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotIdx, w.wantIdx) {
			return errors.Errorf("got = %v, want %v", gotIdx, w.wantIdx)
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
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotIdx, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, gotIdx, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_index_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		client                discoverer.Client
		eg                    errgroup.Group
		creationPoolSize      uint32
		indexDuration         time.Duration
		indexDurationLimit    time.Duration
		concurrency           int
		indexInfos            indexInfos
		indexing              atomic.Value
		minUncommitted        uint32
		uuidsCount            uint32
		uncommittedUUIDsCount uint32
		uncommittedUUIDs      atomic.Value
		uuids                 atomic.Value
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
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
		       },
		       fields: fields {
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
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
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			idx := &index{
				client:                test.fields.client,
				eg:                    test.fields.eg,
				creationPoolSize:      test.fields.creationPoolSize,
				indexDuration:         test.fields.indexDuration,
				indexDurationLimit:    test.fields.indexDurationLimit,
				concurrency:           test.fields.concurrency,
				indexInfos:            test.fields.indexInfos,
				indexing:              test.fields.indexing,
				minUncommitted:        test.fields.minUncommitted,
				uuidsCount:            test.fields.uuidsCount,
				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
				uncommittedUUIDs:      test.fields.uncommittedUUIDs,
				uuids:                 test.fields.uuids,
			}

			got, err := idx.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_index_execute(t *testing.T) {
	type args struct {
		ctx                context.Context
		enableLowIndexSkip bool
	}
	type fields struct {
		client                discoverer.Client
		eg                    errgroup.Group
		creationPoolSize      uint32
		indexDuration         time.Duration
		indexDurationLimit    time.Duration
		concurrency           int
		indexInfos            indexInfos
		indexing              atomic.Value
		minUncommitted        uint32
		uuidsCount            uint32
		uncommittedUUIDsCount uint32
		uncommittedUUIDs      atomic.Value
		uuids                 atomic.Value
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
		           enableLowIndexSkip: false,
		       },
		       fields: fields {
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
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
		           enableLowIndexSkip: false,
		           },
		           fields: fields {
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			idx := &index{
				client:                test.fields.client,
				eg:                    test.fields.eg,
				creationPoolSize:      test.fields.creationPoolSize,
				indexDuration:         test.fields.indexDuration,
				indexDurationLimit:    test.fields.indexDurationLimit,
				concurrency:           test.fields.concurrency,
				indexInfos:            test.fields.indexInfos,
				indexing:              test.fields.indexing,
				minUncommitted:        test.fields.minUncommitted,
				uuidsCount:            test.fields.uuidsCount,
				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
				uncommittedUUIDs:      test.fields.uncommittedUUIDs,
				uuids:                 test.fields.uuids,
			}

			err := idx.execute(test.args.ctx, test.args.enableLowIndexSkip)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_index_loadInfos(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		client                discoverer.Client
		eg                    errgroup.Group
		creationPoolSize      uint32
		indexDuration         time.Duration
		indexDurationLimit    time.Duration
		concurrency           int
		indexInfos            indexInfos
		indexing              atomic.Value
		minUncommitted        uint32
		uuidsCount            uint32
		uncommittedUUIDsCount uint32
		uncommittedUUIDs      atomic.Value
		uuids                 atomic.Value
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
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
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
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			idx := &index{
				client:                test.fields.client,
				eg:                    test.fields.eg,
				creationPoolSize:      test.fields.creationPoolSize,
				indexDuration:         test.fields.indexDuration,
				indexDurationLimit:    test.fields.indexDurationLimit,
				concurrency:           test.fields.concurrency,
				indexInfos:            test.fields.indexInfos,
				indexing:              test.fields.indexing,
				minUncommitted:        test.fields.minUncommitted,
				uuidsCount:            test.fields.uuidsCount,
				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
				uncommittedUUIDs:      test.fields.uncommittedUUIDs,
				uuids:                 test.fields.uuids,
			}

			err := idx.loadInfos(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_index_IsIndexing(t *testing.T) {
	type fields struct {
		client                discoverer.Client
		eg                    errgroup.Group
		creationPoolSize      uint32
		indexDuration         time.Duration
		indexDurationLimit    time.Duration
		concurrency           int
		indexInfos            indexInfos
		indexing              atomic.Value
		minUncommitted        uint32
		uuidsCount            uint32
		uncommittedUUIDsCount uint32
		uncommittedUUIDs      atomic.Value
		uuids                 atomic.Value
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		       fields: fields {
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
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
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			idx := &index{
				client:                test.fields.client,
				eg:                    test.fields.eg,
				creationPoolSize:      test.fields.creationPoolSize,
				indexDuration:         test.fields.indexDuration,
				indexDurationLimit:    test.fields.indexDurationLimit,
				concurrency:           test.fields.concurrency,
				indexInfos:            test.fields.indexInfos,
				indexing:              test.fields.indexing,
				minUncommitted:        test.fields.minUncommitted,
				uuidsCount:            test.fields.uuidsCount,
				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
				uncommittedUUIDs:      test.fields.uncommittedUUIDs,
				uuids:                 test.fields.uuids,
			}

			got := idx.IsIndexing()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_index_NumberOfUUIDs(t *testing.T) {
	type fields struct {
		client                discoverer.Client
		eg                    errgroup.Group
		creationPoolSize      uint32
		indexDuration         time.Duration
		indexDurationLimit    time.Duration
		concurrency           int
		indexInfos            indexInfos
		indexing              atomic.Value
		minUncommitted        uint32
		uuidsCount            uint32
		uncommittedUUIDsCount uint32
		uncommittedUUIDs      atomic.Value
		uuids                 atomic.Value
	}
	type want struct {
		want uint32
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint32) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint32) error {
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
		       fields: fields {
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
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
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			idx := &index{
				client:                test.fields.client,
				eg:                    test.fields.eg,
				creationPoolSize:      test.fields.creationPoolSize,
				indexDuration:         test.fields.indexDuration,
				indexDurationLimit:    test.fields.indexDurationLimit,
				concurrency:           test.fields.concurrency,
				indexInfos:            test.fields.indexInfos,
				indexing:              test.fields.indexing,
				minUncommitted:        test.fields.minUncommitted,
				uuidsCount:            test.fields.uuidsCount,
				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
				uncommittedUUIDs:      test.fields.uncommittedUUIDs,
				uuids:                 test.fields.uuids,
			}

			got := idx.NumberOfUUIDs()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_index_NumberOfUncommittedUUIDs(t *testing.T) {
	type fields struct {
		client                discoverer.Client
		eg                    errgroup.Group
		creationPoolSize      uint32
		indexDuration         time.Duration
		indexDurationLimit    time.Duration
		concurrency           int
		indexInfos            indexInfos
		indexing              atomic.Value
		minUncommitted        uint32
		uuidsCount            uint32
		uncommittedUUIDsCount uint32
		uncommittedUUIDs      atomic.Value
		uuids                 atomic.Value
	}
	type want struct {
		want uint32
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint32) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint32) error {
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
		       fields: fields {
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
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
		           client: nil,
		           eg: nil,
		           creationPoolSize: 0,
		           indexDuration: nil,
		           indexDurationLimit: nil,
		           concurrency: 0,
		           indexInfos: indexInfos{},
		           indexing: nil,
		           minUncommitted: 0,
		           uuidsCount: 0,
		           uncommittedUUIDsCount: 0,
		           uncommittedUUIDs: nil,
		           uuids: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			idx := &index{
				client:                test.fields.client,
				eg:                    test.fields.eg,
				creationPoolSize:      test.fields.creationPoolSize,
				indexDuration:         test.fields.indexDuration,
				indexDurationLimit:    test.fields.indexDurationLimit,
				concurrency:           test.fields.concurrency,
				indexInfos:            test.fields.indexInfos,
				indexing:              test.fields.indexing,
				minUncommitted:        test.fields.minUncommitted,
				uuidsCount:            test.fields.uuidsCount,
				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
				uncommittedUUIDs:      test.fields.uncommittedUUIDs,
				uuids:                 test.fields.uuids,
			}

			got := idx.NumberOfUncommittedUUIDs()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
