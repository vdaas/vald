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

// Package service manages the main logic of server.
package service

import (
	"bufio"
	"context"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/kvs"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/vqueue"
)

func TestNew(t *testing.T) {
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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

			gotNn, err := New(test.args.cfg, test.args.opts...)
			if err := checkFunc(test.want, gotNn, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_load(t *testing.T) {
	type args struct {
		ctx  context.Context
		path string
		opts []core.Option
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           path: "",
		           opts: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           path: "",
		           opts: nil,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.load(test.args.ctx, test.args.path, test.args.opts...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_initNGT(t *testing.T) {
	type args struct {
		opts []core.Option
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.initNGT(test.args.opts...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_loadKVS(t *testing.T) {
	type args struct {
		path string
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           path: "",
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           path: "",
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.loadKVS(test.args.path)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.Start(test.args.ctx)
			if err := checkFunc(test.want, got); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got, err := n.Search(test.args.vec, test.args.size, test.args.epsilon, test.args.radius)
			if err := checkFunc(test.want, got, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
	}
	type want struct {
		wantVec []float32
		wantDst []model.Distance
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, []model.Distance, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotVec []float32, gotDst []model.Distance, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotVec, w.wantVec) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			gotVec, gotDst, err := n.SearchByID(test.args.uuid, test.args.size, test.args.epsilon, test.args.radius)
			if err := checkFunc(test.want, gotVec, gotDst, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_LinearSearch(t *testing.T) {
	type args struct {
		vec  []float32
		size uint32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got, err := n.LinearSearch(test.args.vec, test.args.size)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_LinearSearchByID(t *testing.T) {
	type args struct {
		uuid string
		size uint32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
	}
	type want struct {
		wantVec []float32
		wantDst []model.Distance
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, []model.Distance, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotVec []float32, gotDst []model.Distance, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotVec, w.wantVec) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			gotVec, gotDst, err := n.LinearSearchByID(test.args.uuid, test.args.size)
			if err := checkFunc(test.want, gotVec, gotDst, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.Insert(test.args.uuid, test.args.vec)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_InsertWithTime(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.InsertWithTime(test.args.uuid, test.args.vec, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.insert(test.args.uuid, test.args.vec, test.args.t, test.args.validation)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.InsertMultiple(test.args.vecs)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_InsertMultipleWithTime(t *testing.T) {
	type args struct {
		vecs map[string][]float32
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           t: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           t: 0,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.InsertMultipleWithTime(test.args.vecs, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_insertMultiple(t *testing.T) {
	type args struct {
		vecs       map[string][]float32
		now        int64
		validation bool
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           now: 0,
		           validation: false,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           now: 0,
		           validation: false,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.insertMultiple(test.args.vecs, test.args.now, test.args.validation)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.Update(test.args.uuid, test.args.vec)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_UpdateWithTime(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.UpdateWithTime(test.args.uuid, test.args.vec, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_update(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.update(test.args.uuid, test.args.vec, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.UpdateMultiple(test.args.vecs)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_UpdateMultipleWithTime(t *testing.T) {
	type args struct {
		vecs map[string][]float32
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           t: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           t: 0,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.UpdateMultipleWithTime(test.args.vecs, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_updateMultiple(t *testing.T) {
	type args struct {
		vecs map[string][]float32
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           t: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           t: 0,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.updateMultiple(test.args.vecs, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.Delete(test.args.uuid)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_DeleteWithTime(t *testing.T) {
	type args struct {
		uuid string
		t    int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.DeleteWithTime(test.args.uuid, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_delete(t *testing.T) {
	type args struct {
		uuid       string
		t          int64
		validation bool
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           validation: false,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           validation: false,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.delete(test.args.uuid, test.args.t, test.args.validation)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.DeleteMultiple(test.args.uuids...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_DeleteMultipleWithTime(t *testing.T) {
	type args struct {
		uuids []string
		t     int64
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           t: 0,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           t: 0,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.DeleteMultipleWithTime(test.args.uuids, test.args.t)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_deleteMultiple(t *testing.T) {
	type args struct {
		uuids      []string
		now        int64
		validation bool
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           now: 0,
		           validation: false,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           now: 0,
		           validation: false,
		           },
		           fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.deleteMultiple(test.args.uuids, test.args.now, test.args.validation)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.CreateIndex(test.args.ctx, test.args.poolSize)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_removeInvalidIndex(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           ctx: nil,
		       },
		       fields: fields {
		           core: nil,
		           eg: nil,
		           kvs: nil,
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			n.removeInvalidIndex(test.args.ctx)
			if err := checkFunc(test.want); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.SaveIndex(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_saveIndex(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.saveIndex(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.CreateAndSaveIndex(test.args.ctx, test.args.poolSize)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_moveAndSwitchSavedData(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.moveAndSwitchSavedData(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_mktmp(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.mktmp()
			if err := checkFunc(test.want, err); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			gotOid, gotOk := n.Exists(test.args.uuid)
			if err := checkFunc(test.want, gotOid, gotOk); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			gotVec, err := n.GetObject(test.args.uuid)
			if err := checkFunc(test.want, gotVec, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_readyForUpdate(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
	}
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.readyForUpdate(test.args.uuid, test.args.vec)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_IsSaving(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.IsSaving()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_IsIndexing(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.IsIndexing()
			if err := checkFunc(test.want, got); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			gotUuids := n.UUIDs(test.args.ctx)
			if err := checkFunc(test.want, gotUuids); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_NumberOfCreateIndexExecution(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.NumberOfCreateIndexExecution()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_NumberOfProactiveGCExecution(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.NumberOfProactiveGCExecution()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_gc(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			n.gc()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Len(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.Len()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_InsertVQueueBufferLen(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.InsertVQueueBufferLen()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_DeleteVQueueBufferLen(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.DeleteVQueueBufferLen()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_GetDimensionSize(t *testing.T) {
	type fields struct {
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got int) error {
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			got := n.GetDimensionSize()
			if err := checkFunc(test.want, got); err != nil {
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
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
		           fmu: sync.Mutex{},
		           fmap: nil,
		           vq: nil,
		           indexing: nil,
		           saving: nil,
		           cimu: sync.Mutex{},
		           lastNocie: 0,
		           nocie: 0,
		           nogce: 0,
		           inMem: false,
		           dim: 0,
		           alen: 0,
		           lim: nil,
		           dur: nil,
		           sdur: nil,
		           minLit: nil,
		           maxLit: nil,
		           litFactor: nil,
		           enableProactiveGC: false,
		           enableCopyOnWrite: false,
		           path: "",
		           smu: sync.Mutex{},
		           tmpPath: nil,
		           oldPath: "",
		           basePath: "",
		           cowmu: sync.Mutex{},
		           backupGen: 0,
		           poolSize: 0,
		           radius: 0,
		           epsilon: 0,
		           idelay: nil,
		           dcd: false,
		           kvsdbConcurrency: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			n := &ngt{
				core:              test.fields.core,
				eg:                test.fields.eg,
				kvs:               test.fields.kvs,
				fmu:               test.fields.fmu,
				fmap:              test.fields.fmap,
				vq:                test.fields.vq,
				indexing:          test.fields.indexing,
				saving:            test.fields.saving,
				cimu:              test.fields.cimu,
				lastNocie:         test.fields.lastNocie,
				nocie:             test.fields.nocie,
				nogce:             test.fields.nogce,
				inMem:             test.fields.inMem,
				dim:               test.fields.dim,
				alen:              test.fields.alen,
				lim:               test.fields.lim,
				dur:               test.fields.dur,
				sdur:              test.fields.sdur,
				minLit:            test.fields.minLit,
				maxLit:            test.fields.maxLit,
				litFactor:         test.fields.litFactor,
				enableProactiveGC: test.fields.enableProactiveGC,
				enableCopyOnWrite: test.fields.enableCopyOnWrite,
				path:              test.fields.path,
				smu:               test.fields.smu,
				tmpPath:           test.fields.tmpPath,
				oldPath:           test.fields.oldPath,
				basePath:          test.fields.basePath,
				cowmu:             test.fields.cowmu,
				backupGen:         test.fields.backupGen,
				poolSize:          test.fields.poolSize,
				radius:            test.fields.radius,
				epsilon:           test.fields.epsilon,
				idelay:            test.fields.idelay,
				dcd:               test.fields.dcd,
				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
			}

			err := n.Close(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

type index struct {
	uuid string
	vec  []float32
}

func Test_ngt_InsertUpsert(t *testing.T) {
	type args struct {
		idxes []index
	}
	type fields struct {
		svcCfg  *config.NGT
		svcOpts []Option

		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
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
	var (
		// default NGT configuration for test
		kvsdbCfg  = &config.KVSDB{}
		vqueueCfg = &config.VQueue{}
	)
	idxes := readDiffFile("./20221025")
	tests := []test{
		{
			name: "insert & upsert 1",
			args: args{
				idxes: []index{
					{
						uuid: "19731011_puls-0000303:43171:1",
						vec: []float32{
							-0.018963795, 0.13975345, -0.01438206, 0.021553166, -0.1163193, 0.02389161, -0.043322664, 0.051071428, -0.09778391, 0.17714833, 0.070366226, 0.00964489, 0.123230025, -0.0072021615, -0.12053576, -0.25266534, -0.05868739, -0.06946657, -0.13858624, 0.08516157, -0.123261474, 0.14827089, -0.15974368, -0.045519657, -0.15578407, 0.12303146, -0.10828495, 0.061408218, -0.12964113, 0.02870139, 0.2465153, 0.16652209, 0.12929203, -0.24059042, 0.08436171, 0.018702978, -0.005219509, -0.07248656, 0.14141396, -0.13592686, 0.044463452, 0.083398245, -0.12660366, 0.01033342, 0.10309342, 0.1062651, -0.04192098, -0.024938857, -0.17638475, -0.04105866, 0.16464227, 0.09967009, -0.24130723, 0.05822357, 0.065308936, 0.14314993, -0.025086135, -0.014758362, -0.12474182, -0.10733189, -0.16456085, 0.17825097, 0.048106, -0.17680661, 0.02702824, 0.4698901, -0.41418663, 0.4129158, -0.02523511, 0.36365706, -0.30040222, 0.5185276, -0.51615494, 0.5655249, -0.092418216, 0.07052821, 0.40197027, 0.084085345, -0.4798388, -0.5527801, 0.12522608, -0.13325042, -0.4993201, -0.31595704, -0.35174417, 0.27274206, 0.2663609, 0.40937683, -0.45539486, 0.16548769, -0.7304982, 0.2212803, -0.70699847, -0.019003047, 0.63269913, 0.23438695, 0.14075014, 0.034899414, 0.1016486, -0.5403803, 0.049518738, 0.08046312, 0.3919085, -0.1290008, -0.07870986, -0.060331456, 0.1759862, 0.12838659, 0.17038107, -0.09322341, -0.05854857, 0.029937917, -0.0010719631, -0.11727104, 0.1532171, -0.07039561, 0.016434262, -0.14807108, 0.08238894, 0.23709965, 0.15327819, -0.07462129, 0.036771126, 0.07617792, 0.22700097, -0.10173196, 0.15948103, -0.18994464 - 0.018963795, 0.13975345, -0.01438206, 0.021553166, -0.1163193, 0.02389161, -0.043322664, 0.051071428, -0.09778391, 0.17714833, 0.070366226, 0.00964489, 0.123230025, -0.0072021615, -0.12053576, -0.25266534, -0.05868739, -0.06946657, -0.13858624, 0.08516157, -0.123261474, 0.14827089, -0.15974368, -0.045519657, -0.15578407, 0.12303146, -0.10828495, 0.061408218, -0.12964113, 0.02870139, 0.2465153, 0.16652209, 0.12929203, -0.24059042, 0.08436171, 0.018702978, -0.005219509, -0.07248656, 0.14141396, -0.13592686, 0.044463452, 0.083398245, -0.12660366, 0.01033342, 0.10309342, 0.1062651, -0.04192098, -0.024938857, -0.17638475, -0.04105866, 0.16464227, 0.09967009, -0.24130723, 0.05822357, 0.065308936, 0.14314993, -0.025086135, -0.014758362, -0.12474182, -0.10733189, -0.16456085, 0.17825097, 0.048106, -0.17680661, 0.02702824, 0.4698901, -0.41418663, 0.4129158, -0.02523511, 0.36365706, -0.30040222, 0.5185276, -0.51615494, 0.5655249, -0.092418216, 0.07052821, 0.40197027, 0.084085345, -0.4798388, -0.5527801, 0.12522608, -0.13325042, -0.4993201, -0.31595704, -0.35174417, 0.27274206, 0.2663609, 0.40937683, -0.45539486, 0.16548769, -0.7304982, 0.2212803, -0.70699847, -0.019003047, 0.63269913, 0.23438695, 0.14075014, 0.034899414, 0.1016486, -0.5403803, 0.049518738, 0.08046312, 0.3919085, -0.1290008, -0.07870986, -0.060331456, 0.1759862, 0.12838659, 0.17038107, -0.09322341, -0.05854857, 0.029937917, -0.0010719631, -0.11727104, 0.1532171, -0.07039561, 0.016434262, -0.14807108, 0.08238894, 0.23709965, 0.15327819, -0.07462129, 0.036771126, 0.07617792, 0.22700097, -0.10173196, 0.15948103, -0.18994464,
						},
					},
				},
			},
			fields: fields{
				svcCfg: &config.NGT{
					Dimension:    128,
					DistanceType: ngt.Cosine.String(),
					ObjectType:   ngt.Uint8.String(),
					KVSDB:        kvsdbCfg,
					VQueue:       vqueueCfg,
				},
				svcOpts: []Option{
					WithEnableInMemoryMode(true),
				},
			},
		},
		{
			name: "insert & upsert all",
			args: args{
				idxes: idxes,
			},
			fields: fields{
				svcCfg: &config.NGT{
					Dimension:    128,
					DistanceType: ngt.Cosine.String(),
					ObjectType:   ngt.Uint8.String(),
					KVSDB:        kvsdbCfg,
					VQueue:       vqueueCfg,
				},
				svcOpts: []Option{
					WithEnableInMemoryMode(true),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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

			eg, _ := errgroup.New(ctx)
			n, err := New(test.fields.svcCfg, append(test.fields.svcOpts, WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}

			count := 0
			for _, idx := range test.args.idxes {
				count++
				err = n.Insert(idx.uuid, idx.vec)
				if err := checkFunc(test.want, err); err != nil {
					tt.Errorf("error = %v", err)
				}

				if count%100 == 0 {
					err = n.CreateAndSaveIndex(ctx, 100)
					if err != nil {
						tt.Errorf("error creating index: %v", err)
					}
				}

				err = n.Update(idx.uuid, idx.vec)
				if err := checkFunc(test.want, err); err != nil {
					tt.Errorf("error = %v", err)
				}

				if count%100 == 0 {
					err = n.CreateAndSaveIndex(ctx, 100)
					if err != nil {
						tt.Errorf("error creating index: %v", err)
					}
				}
			}
		})
	}
}

func readDiffFile(filePath string) []index {
	// only read all data with insert
	f, err := os.Open(filePath)
	if err != nil {
		glg.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	result := make([]index, 0)
	for scanner.Scan() {
		row := scanner.Text()

		r := strings.Split(row, "\t")
		id := r[0]
		action := r[1]
		var vec []float32

		if action != "delete" {
			vecSlice := r[2:]
			vec = make([]float32, len(vecSlice))
			for i, v := range vecSlice {
				f, err := strconv.ParseFloat(v, 32)
				if err != nil {
					glg.Fatal(err)
				}
				vec[i] = float32(f)
			}

			if len(vec) != 128 {
				glg.Error("invalid vec lengh, ID: %v, vec: %v", id, vec)
			}
		}

		if action == "insert" {
			result = append(result, index{
				uuid: id,
				vec:  vec,
			})
			// _, err = client.Upsert(ctx, &payload.Upsert_Request{
			// 	Vector: &payload.Object_Vector{
			// 		Id:     id,
			// 		Vector: vec,
			// 	},
			// 	Config: &payload.Upsert_Config{
			// 		SkipStrictExistCheck: true,
			// 	},
			// })
			// } else if action == "delete" {
			// 	_, err = client.Remove(ctx, &payload.Remove_Request{
			// 		Id: &payload.Object_ID{
			// 			Id: id,
			// 		},
			// 		Config: &payload.Remove_Config{
			// 			SkipStrictExistCheck: true,
			// 		},
			// 	})
		}
	}
	return result
}
