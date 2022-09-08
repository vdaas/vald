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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func Test_newAddr(t *testing.T) {
	t.Parallel()
	type args struct {
		addrList map[string]struct{}
	}
	type want struct {
		want AtomicAddrs
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, AtomicAddrs) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got AtomicAddrs) error {
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
		           addrList: nil,
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
		           addrList: nil,
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

			got := newAddr(test.args.addrList)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_atomicAddrs_GetAll(t *testing.T) {
	t.Parallel()
	type fields struct {
		addrs      atomic.Value
		dupCheck   map[string]bool
		mu         sync.RWMutex
		addrSeeker uint64
	}
	type want struct {
		want  []string
		want1 bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []string, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []string, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
			a := &atomicAddrs{
				addrs:      test.fields.addrs,
				dupCheck:   test.fields.dupCheck,
				mu:         test.fields.mu,
				addrSeeker: test.fields.addrSeeker,
			}

			got, got1 := a.GetAll()
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_atomicAddrs_Range(t *testing.T) {
	t.Parallel()
	type args struct {
		f func(addr string) bool
	}
	type fields struct {
		addrs      atomic.Value
		dupCheck   map[string]bool
		mu         sync.RWMutex
		addrSeeker uint64
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
		           f: nil,
		       },
		       fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
		           f: nil,
		           },
		           fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
			a := &atomicAddrs{
				addrs:      test.fields.addrs,
				dupCheck:   test.fields.dupCheck,
				mu:         test.fields.mu,
				addrSeeker: test.fields.addrSeeker,
			}

			a.Range(test.args.f)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_atomicAddrs_Add(t *testing.T) {
	t.Parallel()
	type args struct {
		addr string
	}
	type fields struct {
		addrs      atomic.Value
		dupCheck   map[string]bool
		mu         sync.RWMutex
		addrSeeker uint64
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
		           addr: "",
		       },
		       fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
		           addr: "",
		           },
		           fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
			a := &atomicAddrs{
				addrs:      test.fields.addrs,
				dupCheck:   test.fields.dupCheck,
				mu:         test.fields.mu,
				addrSeeker: test.fields.addrSeeker,
			}

			a.Add(test.args.addr)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_atomicAddrs_Delete(t *testing.T) {
	t.Parallel()
	type args struct {
		addr string
	}
	type fields struct {
		addrs      atomic.Value
		dupCheck   map[string]bool
		mu         sync.RWMutex
		addrSeeker uint64
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
		           addr: "",
		       },
		       fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
		           addr: "",
		           },
		           fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
			a := &atomicAddrs{
				addrs:      test.fields.addrs,
				dupCheck:   test.fields.dupCheck,
				mu:         test.fields.mu,
				addrSeeker: test.fields.addrSeeker,
			}

			a.Delete(test.args.addr)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_atomicAddrs_Next(t *testing.T) {
	t.Parallel()
	type fields struct {
		addrs      atomic.Value
		dupCheck   map[string]bool
		mu         sync.RWMutex
		addrSeeker uint64
	}
	type want struct {
		want  string
		want1 bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
			a := &atomicAddrs{
				addrs:      test.fields.addrs,
				dupCheck:   test.fields.dupCheck,
				mu:         test.fields.mu,
				addrSeeker: test.fields.addrSeeker,
			}

			got, got1 := a.Next()
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_atomicAddrs_Len(t *testing.T) {
	t.Parallel()
	type fields struct {
		addrs      atomic.Value
		dupCheck   map[string]bool
		mu         sync.RWMutex
		addrSeeker uint64
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
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
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
		           addrs: nil,
		           dupCheck: nil,
		           mu: sync.RWMutex{},
		           addrSeeker: 0,
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
			a := &atomicAddrs{
				addrs:      test.fields.addrs,
				dupCheck:   test.fields.dupCheck,
				mu:         test.fields.mu,
				addrSeeker: test.fields.addrSeeker,
			}

			got := a.Len()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
