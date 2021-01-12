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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/kvs"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		cfg  *config.NGT
		opts []Option
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotNn, w.wantNn) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNn, w.wantNn)
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
		           cfg: nil,
		           opts: nil,
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

			gotNn, err := New(test.args.cfg, test.args.opts...)
			if err := test.checkFunc(test.want, gotNn, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_initNGT(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []core.Option
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           opts: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.initNGT(test.args.opts...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_loadKVS(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.loadKVS()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Start(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Search(t *testing.T) {
	t.Parallel()
	type args struct {
		vec     []float32
		size    uint32
		epsilon float32
		radius  float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           vec: nil,
		           size: 0,
		           epsilon: 0,
		           radius: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got, err := n.Search(test.args.vec, test.args.size, test.args.epsilon, test.args.radius)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_SearchByID(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid    string
		size    uint32
		epsilon float32
		radius  float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotDst, w.wantDst) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDst, w.wantDst)
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			gotDst, err := n.SearchByID(test.args.uuid, test.args.size, test.args.epsilon, test.args.radius)
			if err := test.checkFunc(test.want, gotDst, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Insert(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
		vec  []float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           uuid: "",
		           vec: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.Insert(test.args.uuid, test.args.vec)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_insert(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid       string
		vec        []float32
		t          int64
		validation bool
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           uuid: "",
		           vec: nil,
		           t: 0,
		           validation: false,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.insert(test.args.uuid, test.args.vec, test.args.t, test.args.validation)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_InsertMultiple(t *testing.T) {
	t.Parallel()
	type args struct {
		vecs map[string][]float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           vecs: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.InsertMultiple(test.args.vecs)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Update(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
		vec  []float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           uuid: "",
		           vec: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.Update(test.args.uuid, test.args.vec)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_UpdateMultiple(t *testing.T) {
	t.Parallel()
	type args struct {
		vecs map[string][]float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           vecs: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.UpdateMultiple(test.args.vecs)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Delete(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           uuid: "",
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.Delete(test.args.uuid)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_delete(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           uuid: "",
		           t: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.delete(test.args.uuid, test.args.t)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_DeleteMultiple(t *testing.T) {
	t.Parallel()
	type args struct {
		uuids []string
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           uuids: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.DeleteMultiple(test.args.uuids...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_GetObject(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotVec, w.wantVec) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			gotVec, err := n.GetObject(test.args.uuid)
			if err := test.checkFunc(test.want, gotVec, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_CreateIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx      context.Context
		poolSize uint32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           poolSize: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.CreateIndex(test.args.ctx, test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_SaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.SaveIndex(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_saveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.saveIndex(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_CreateAndSaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx      context.Context
		poolSize uint32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           poolSize: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.CreateAndSaveIndex(test.args.ctx, test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Exists(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOid, w.wantOid)
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
		           uuid: "",
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			gotOid, gotOk := n.Exists(test.args.uuid)
			if err := test.checkFunc(test.want, gotOid, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_readyForUpdate(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
		vec  []float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
	}
	type want struct {
		wantReady bool
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
	defaultCheckFunc := func(w want, gotReady bool) error {
		if !reflect.DeepEqual(gotReady, w.wantReady) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotReady, w.wantReady)
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			gotReady := n.readyForUpdate(test.args.uuid, test.args.vec)
			if err := test.checkFunc(test.want, gotReady); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_insertCache(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		       args: args {
		           uuid: "",
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got, got1 := n.insertCache(test.args.uuid)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_IsSaving(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.IsSaving()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_IsIndexing(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.IsIndexing()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_UUIDs(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuids, w.wantUuids)
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			gotUuids := n.UUIDs(test.args.ctx)
			if err := test.checkFunc(test.want, gotUuids); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_UncommittedUUIDs(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuids, w.wantUuids)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			gotUuids := n.UncommittedUUIDs()
			if err := test.checkFunc(test.want, gotUuids); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_NumberOfCreateIndexExecution(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.NumberOfCreateIndexExecution()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_NumberOfProactiveGCExecution(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.NumberOfProactiveGCExecution()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_gc(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			n.gc()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Len(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_InsertVCacheLen(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.InsertVCacheLen()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_DeleteVCacheLen(t *testing.T) {
	t.Parallel()
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			got := n.DeleteVCacheLen()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Close(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		ivc               *vcaches
		dvc               *vcaches
		indexing          atomic.Value
		saving            atomic.Value
		lastNoice         uint64
		ic                uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		path              string
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           ivc: vcaches{},
		           dvc: vcaches{},
		           indexing: nil,
		           saving: nil,
		           lastNoice: 0,
		           ic: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           path: "",
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				ivc:               test.fields.ivc,
				dvc:               test.fields.dvc,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				lastNoice:         test.fields.lastNoice,
				ic:                test.fields.ic,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				path:              test.fields.path,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
			}

			err := n.Close(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
