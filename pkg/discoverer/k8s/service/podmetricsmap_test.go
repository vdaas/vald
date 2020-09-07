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

	"github.com/vdaas/vald/internal/errors"
	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
)

func Test_newEntryPodMetricsMap(t *testing.T) {
	type args struct {
		i mpod.Pod
	}
	type want struct {
		want *entryPodMetricsMap
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *entryPodMetricsMap) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *entryPodMetricsMap) error {
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

			got := newEntryPodMetricsMap(test.args.i)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_podMetricsMap_Load(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryPodMetricsMap
		misses int
	}
	type want struct {
		wantValue mpod.Pod
		wantOk    bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, mpod.Pod, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotValue mpod.Pod, gotOk bool) error {
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
			m := &podMetricsMap{
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

func Test_entryPodMetricsMap_load(t *testing.T) {
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
		wantValue mpod.Pod
		wantOk    bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, mpod.Pod, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotValue mpod.Pod, gotOk bool) error {
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
			e := &entryPodMetricsMap{
				p: test.fields.p,
			}

			gotValue, gotOk := e.load()
			if err := test.checkFunc(test.want, gotValue, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_podMetricsMap_Store(t *testing.T) {
	type args struct {
		key   string
		value mpod.Pod
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryPodMetricsMap
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
			m := &podMetricsMap{
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

func Test_entryPodMetricsMap_tryStore(t *testing.T) {
	type args struct {
		i *mpod.Pod
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
			e := &entryPodMetricsMap{
				p: test.fields.p,
			}

			got := e.tryStore(test.args.i)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_entryPodMetricsMap_unexpungeLocked(t *testing.T) {
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
			e := &entryPodMetricsMap{
				p: test.fields.p,
			}

			gotWasExpunged := e.unexpungeLocked()
			if err := test.checkFunc(test.want, gotWasExpunged); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_entryPodMetricsMap_storeLocked(t *testing.T) {
	type args struct {
		i *mpod.Pod
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
			e := &entryPodMetricsMap{
				p: test.fields.p,
			}

			e.storeLocked(test.args.i)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_podMetricsMap_LoadOrStore(t *testing.T) {
	type args struct {
		key   string
		value mpod.Pod
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryPodMetricsMap
		misses int
	}
	type want struct {
		wantActual mpod.Pod
		wantLoaded bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, mpod.Pod, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotActual mpod.Pod, gotLoaded bool) error {
		if !reflect.DeepEqual(gotActual, w.wantActual) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotActual, w.wantActual)
		}
		if !reflect.DeepEqual(gotLoaded, w.wantLoaded) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoaded, w.wantLoaded)
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
			m := &podMetricsMap{
				mu:     test.fields.mu,
				read:   test.fields.read,
				dirty:  test.fields.dirty,
				misses: test.fields.misses,
			}

			gotActual, gotLoaded := m.LoadOrStore(test.args.key, test.args.value)
			if err := test.checkFunc(test.want, gotActual, gotLoaded); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_entryPodMetricsMap_tryLoadOrStore(t *testing.T) {
	type args struct {
		i mpod.Pod
	}
	type fields struct {
		p unsafe.Pointer
	}
	type want struct {
		wantActual mpod.Pod
		wantLoaded bool
		wantOk     bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, mpod.Pod, bool, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotActual mpod.Pod, gotLoaded bool, gotOk bool) error {
		if !reflect.DeepEqual(gotActual, w.wantActual) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotActual, w.wantActual)
		}
		if !reflect.DeepEqual(gotLoaded, w.wantLoaded) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoaded, w.wantLoaded)
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
			e := &entryPodMetricsMap{
				p: test.fields.p,
			}

			gotActual, gotLoaded, gotOk := e.tryLoadOrStore(test.args.i)
			if err := test.checkFunc(test.want, gotActual, gotLoaded, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_podMetricsMap_Delete(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryPodMetricsMap
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
			m := &podMetricsMap{
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

func Test_entryPodMetricsMap_delete(t *testing.T) {
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
			e := &entryPodMetricsMap{
				p: test.fields.p,
			}

			gotHadValue := e.delete()
			if err := test.checkFunc(test.want, gotHadValue); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_podMetricsMap_Range(t *testing.T) {
	type args struct {
		f func(key string, value mpod.Pod) bool
	}
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryPodMetricsMap
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
			m := &podMetricsMap{
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

func Test_podMetricsMap_missLocked(t *testing.T) {
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryPodMetricsMap
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
			m := &podMetricsMap{
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

func Test_podMetricsMap_dirtyLocked(t *testing.T) {
	type fields struct {
		mu     sync.Mutex
		read   atomic.Value
		dirty  map[string]*entryPodMetricsMap
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
			m := &podMetricsMap{
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

func Test_entryPodMetricsMap_tryExpungeLocked(t *testing.T) {
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
			e := &entryPodMetricsMap{
				p: test.fields.p,
			}

			gotIsExpunged := e.tryExpungeLocked()
			if err := test.checkFunc(test.want, gotIsExpunged); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
