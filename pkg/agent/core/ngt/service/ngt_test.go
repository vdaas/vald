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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	core "github.com/vdaas/vald/internal/core/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/kvs"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.NGT
	}
	type want struct {
		wantNn NGT
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, NGT, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotNn NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotNn, w.wantNn) {
			return errors.Errorf("got = %v, want %v", gotNn, w.wantNn)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           cfg: nil,
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
		           cfg: nil,
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

			gotNn, err := New(test.args.cfg)
			if err := test.checkFunc(test.want, gotNn, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		want <-chan error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error) error {
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got := n.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Search(t *testing.T) {
	type args struct {
		vec     []float32
		size    uint32
		epsilon float32
		radius  float32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		want []model.Distance
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []model.Distance, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []model.Distance, err error) error {
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
		           vec: nil,
		           size: 0,
		           epsilon: 0,
		           radius: 0,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           vec: nil,
		           size: 0,
		           epsilon: 0,
		           radius: 0,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got, err := n.Search(test.args.vec, test.args.size, test.args.epsilon, test.args.radius)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_SearchByID(t *testing.T) {
	type args struct {
		uuid    string
		size    uint32
		epsilon float32
		radius  float32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		wantDst []model.Distance
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []model.Distance, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotDst []model.Distance, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotDst, w.wantDst) {
			return errors.Errorf("got = %v, want %v", gotDst, w.wantDst)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           uuid: "",
		           size: 0,
		           epsilon: 0,
		           radius: 0,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           size: 0,
		           epsilon: 0,
		           radius: 0,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			gotDst, err := n.SearchByID(test.args.uuid, test.args.size, test.args.epsilon, test.args.radius)
			if err := test.checkFunc(test.want, gotDst, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Insert(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           uuid: "",
		           vec: nil,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           vec: nil,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.Insert(test.args.uuid, test.args.vec)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_insert(t *testing.T) {
	type args struct {
		uuid       string
		vec        []float32
		t          int64
		validation bool
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           uuid: "",
		           vec: nil,
		           t: 0,
		           validation: false,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           vec: nil,
		           t: 0,
		           validation: false,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.insert(test.args.uuid, test.args.vec, test.args.t, test.args.validation)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_InsertMultiple(t *testing.T) {
	type args struct {
		vecs map[string][]float32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           vecs: nil,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           vecs: nil,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.InsertMultiple(test.args.vecs)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Update(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           uuid: "",
		           vec: nil,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           vec: nil,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.Update(test.args.uuid, test.args.vec)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_UpdateMultiple(t *testing.T) {
	type args struct {
		vecs map[string][]float32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           vecs: nil,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           vecs: nil,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.UpdateMultiple(test.args.vecs)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Delete(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           uuid: "",
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.Delete(test.args.uuid)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_delete(t *testing.T) {
	type args struct {
		uuid string
		t    int64
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           uuid: "",
		           t: 0,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           t: 0,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.delete(test.args.uuid, test.args.t)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_DeleteMultiple(t *testing.T) {
	type args struct {
		uuids []string
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           uuids: nil,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuids: nil,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.DeleteMultiple(test.args.uuids...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_GetObject(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		wantVec []float32
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotVec []float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotVec, w.wantVec) {
			return errors.Errorf("got = %v, want %v", gotVec, w.wantVec)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           uuid: "",
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			gotVec, err := n.GetObject(test.args.uuid)
			if err := test.checkFunc(test.want, gotVec, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_CreateIndex(t *testing.T) {
	type args struct {
		ctx      context.Context
		poolSize uint32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           poolSize: 0,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           poolSize: 0,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.CreateIndex(test.args.ctx, test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_SaveIndex(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.SaveIndex(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_CreateAndSaveIndex(t *testing.T) {
	type args struct {
		ctx      context.Context
		poolSize uint32
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           poolSize: 0,
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           poolSize: 0,
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.CreateAndSaveIndex(test.args.ctx, test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Exists(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		wantOid uint32
		wantOk  bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint32, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotOid uint32, gotOk bool) error {
		if !reflect.DeepEqual(gotOid, w.wantOid) {
			return errors.Errorf("got = %v, want %v", gotOid, w.wantOid)
		}
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got = %v, want %v", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           uuid: "",
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			gotOid, gotOk := n.Exists(test.args.uuid)
			if err := test.checkFunc(test.want, gotOid, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_insertCache(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		want  *vcache
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *vcache, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *vcache, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want %v", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           uuid: "",
		       },
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           uuid: "",
		           },
		           fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got, got1 := n.insertCache(test.args.uuid)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_IsIndexing(t *testing.T) {
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got := n.IsIndexing()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_UUIDs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		wantUuids []string
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotUuids []string) error {
		if !reflect.DeepEqual(gotUuids, w.wantUuids) {
			return errors.Errorf("got = %v, want %v", gotUuids, w.wantUuids)
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			gotUuids := n.UUIDs(test.args.ctx)
			if err := test.checkFunc(test.want, gotUuids); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_UncommittedUUIDs(t *testing.T) {
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
	}
	type want struct {
		wantUuids []string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotUuids []string) error {
		if !reflect.DeepEqual(gotUuids, w.wantUuids) {
			return errors.Errorf("got = %v, want %v", gotUuids, w.wantUuids)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			gotUuids := n.UncommittedUUIDs()
			if err := test.checkFunc(test.want, gotUuids); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_NumberOfCreateIndexExecution(t *testing.T) {
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got := n.NumberOfCreateIndexExecution()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Len(t *testing.T) {
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got := n.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_InsertVCacheLen(t *testing.T) {
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got := n.InsertVCacheLen()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_DeleteVCacheLen(t *testing.T) {
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			got := n.DeleteVCacheLen()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Close(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		alen     int
		indexing atomic.Value
		lim      time.Duration
		dur      time.Duration
		sdur     time.Duration
		idelay   time.Duration
		dps      uint32
		ic       uint64
		nocie    uint64
		eg       errgroup.Group
		ivc      *vcaches
		dvc      *vcaches
		path     string
		kvs      kvs.BidiMap
		core     core.NGT
		dcd      bool
		inMem    bool
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
		           alen: 0,
		           indexing: nil,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           idelay: nil,
		           dps: 0,
		           ic: 0,
		           nocie: 0,
		           eg: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           path: "",
		           kvs: nil,
		           core: nil,
		           dcd: false,
		           inMem: false,
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
			n := &ngt{
				alen:     test.fields.alen,
				indexing: test.fields.indexing,
				lim:      test.fields.lim,
				dur:      test.fields.dur,
				sdur:     test.fields.sdur,
				idelay:   test.fields.idelay,
				dps:      test.fields.dps,
				ic:       test.fields.ic,
				nocie:    test.fields.nocie,
				eg:       test.fields.eg,
				ivc:      test.fields.ivc,
				dvc:      test.fields.dvc,
				path:     test.fields.path,
				kvs:      test.fields.kvs,
				core:     test.fields.core,
				dcd:      test.fields.dcd,
				inMem:    test.fields.inMem,
			}

			err := n.Close(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
