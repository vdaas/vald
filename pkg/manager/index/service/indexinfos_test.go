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

package service

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
)

func Test_newEntryIndexInfos(t *testing.T) {
	type args struct {
		i *payload.Info_Index_Count
	}
	type want struct {
		want *entryIndexInfos
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *entryIndexInfos) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *entryIndexInfos) error {
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
		           i: nil,
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
		           i: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := newEntryIndexInfos(test.args.i)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_indexInfos_Load(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryIndexInfos
		misses int
	}
	type want struct {
		wantValue *payload.Info_Index_Count
		wantOk    bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Index_Count, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotValue *payload.Info_Index_Count, gotOk bool) error {
		if !reflect.DeepEqual(gotValue, w.wantValue) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotValue, w.wantValue)
		}
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           key: "",
		       },
		       fields: fields {
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
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
		           key: "",
		           },
		           fields: fields {
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &indexInfos{
				mu:     test.fields.mu,
				read:   test.fields.read,
				dirty:  test.fields.dirty,
				misses: test.fields.misses,
			}

			gotValue, gotOk := m.Load(test.args.key)
			if err := test.checkFunc(test.want, gotValue, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_entryIndexInfos_load(t *testing.T) {
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
		wantValue *payload.Info_Index_Count
		wantOk    bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Index_Count, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotValue *payload.Info_Index_Count, gotOk bool) error {
		if !reflect.DeepEqual(gotValue, w.wantValue) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotValue, w.wantValue)
		}
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           p: nil,
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
		           p: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &entryIndexInfos{
				p: test.fields.p,
			}

			gotValue, gotOk := e.load()
			if err := test.checkFunc(test.want, gotValue, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_indexInfos_Store(t *testing.T) {
	type args struct {
		key   string
		value *payload.Info_Index_Count
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryIndexInfos
		misses int
	}
	type want struct {
	}
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
		           key: "",
		           value: nil,
		       },
		       fields: fields {
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
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
		           key: "",
		           value: nil,
		           },
		           fields: fields {
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &indexInfos{
				mu:     test.fields.mu,
				read:   test.fields.read,
				dirty:  test.fields.dirty,
				misses: test.fields.misses,
			}

			m.Store(test.args.key, test.args.value)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_entryIndexInfos_tryStore(t *testing.T) {
	type args struct {
		i **payload.Info_Index_Count
	}
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		           i: nil,
		       },
		       fields: fields {
		           p: nil,
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
		           i: nil,
		           },
		           fields: fields {
		           p: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &entryIndexInfos{
				p: test.fields.p,
			}

			got := e.tryStore(test.args.i)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_entryIndexInfos_unexpungeLocked(t *testing.T) {
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
		wantWasExpunged bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotWasExpunged bool) error {
		if !reflect.DeepEqual(gotWasExpunged, w.wantWasExpunged) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWasExpunged, w.wantWasExpunged)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           p: nil,
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
		           p: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &entryIndexInfos{
				p: test.fields.p,
			}

			gotWasExpunged := e.unexpungeLocked()
			if err := test.checkFunc(test.want, gotWasExpunged); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_entryIndexInfos_storeLocked(t *testing.T) {
	type args struct {
		i **payload.Info_Index_Count
	}
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
	}
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
		           i: nil,
		       },
		       fields: fields {
		           p: nil,
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
		           i: nil,
		           },
		           fields: fields {
		           p: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &entryIndexInfos{
				p: test.fields.p,
			}

			e.storeLocked(test.args.i)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_indexInfos_Delete(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryIndexInfos
		misses int
	}
	type want struct {
	}
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
		           key: "",
		       },
		       fields: fields {
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
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
		           key: "",
		           },
		           fields: fields {
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &indexInfos{
				mu:     test.fields.mu,
				read:   test.fields.read,
				dirty:  test.fields.dirty,
				misses: test.fields.misses,
			}

			m.Delete(test.args.key)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_entryIndexInfos_delete(t *testing.T) {
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
		wantHadValue bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotHadValue bool) error {
		if !reflect.DeepEqual(gotHadValue, w.wantHadValue) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHadValue, w.wantHadValue)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           p: nil,
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
		           p: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &entryIndexInfos{
				p: test.fields.p,
			}

			gotHadValue := e.delete()
			if err := test.checkFunc(test.want, gotHadValue); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_indexInfos_Range(t *testing.T) {
	type args struct {
		f func(key string, value *payload.Info_Index_Count) bool
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryIndexInfos
		misses int
	}
	type want struct {
	}
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
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
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
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &indexInfos{
				mu:     test.fields.mu,
				read:   test.fields.read,
				dirty:  test.fields.dirty,
				misses: test.fields.misses,
			}

			m.Range(test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_indexInfos_missLocked(t *testing.T) {
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryIndexInfos
		misses int
	}
	type want struct {
	}
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
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
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
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &indexInfos{
				mu:     test.fields.mu,
				read:   test.fields.read,
				dirty:  test.fields.dirty,
				misses: test.fields.misses,
			}

			m.missLocked()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_indexInfos_dirtyLocked(t *testing.T) {
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryIndexInfos
		misses int
	}
	type want struct {
	}
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
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
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
		           mu: sync.Mutex{},
		           read: nil,
		           dirty: nil,
		           misses: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &indexInfos{
				mu:     test.fields.mu,
				read:   test.fields.read,
				dirty:  test.fields.dirty,
				misses: test.fields.misses,
			}

			m.dirtyLocked()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_entryIndexInfos_tryExpungeLocked(t *testing.T) {
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
		wantIsExpunged bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotIsExpunged bool) error {
		if !reflect.DeepEqual(gotIsExpunged, w.wantIsExpunged) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIsExpunged, w.wantIsExpunged)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           p: nil,
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
		           p: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &entryIndexInfos{
				p: test.fields.p,
			}

			gotIsExpunged := e.tryExpungeLocked()
			if err := test.checkFunc(test.want, gotIsExpunged); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
