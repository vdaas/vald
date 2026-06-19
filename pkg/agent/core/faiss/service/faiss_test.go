// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		cfg  *config.Faiss
// 		opts []Option
// 	}
// 	type want struct {
// 		want Faiss
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Faiss, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Faiss, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           cfg:nil,
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.cfg, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_initFaiss(t *testing.T) {
// 	type args struct {
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.initFaiss(test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_load(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		path string
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.load(test.args.ctx, test.args.path, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_loadKVS(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		path    string
// 		timeout time.Duration
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           timeout:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           timeout:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.loadKVS(test.args.ctx, test.args.path, test.args.timeout)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_mktmp(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.mktmp()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want <-chan error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Train(t *testing.T) {
// 	type args struct {
// 		nb int
// 		xb []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           nb:0,
// 		           xb:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           nb:0,
// 		           xb:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Train(test.args.nb, test.args.xb)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Insert(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Insert(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_InsertWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.InsertWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_insert(t *testing.T) {
// 	type args struct {
// 		uuid       string
// 		xb         []float32
// 		t          int64
// 		validation bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           xb:nil,
// 		           t:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           xb:nil,
// 		           t:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.insert(test.args.uuid, test.args.xb, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Update(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Update(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_UpdateWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.UpdateWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_update(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.update(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_UpdateTimestamp(t *testing.T) {
// 	type args struct {
// 		uuid  string
// 		ts    int64
// 		force bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.UpdateTimestamp(test.args.uuid, test.args.ts, test.args.force)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_readyForUpdate(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.readyForUpdate(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_CreateIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.CreateIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_SaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.SaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_saveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.saveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_moveAndSwitchSavedData(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.moveAndSwitchSavedData(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_CreateAndSaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.CreateAndSaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Search(t *testing.T) {
// 	type args struct {
// 		k      uint32
// 		nprobe uint32
// 		nq     uint32
// 		xq     []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           k:0,
// 		           nprobe:0,
// 		           nq:0,
// 		           xq:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           k:0,
// 		           nprobe:0,
// 		           nq:0,
// 		           xq:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotRes, err := f.Search(test.args.k, test.args.nprobe, test.args.nq, test.args.xq)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Delete(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Delete(test.args.uuid)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_DeleteWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.DeleteWithTime(test.args.uuid, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_delete(t *testing.T) {
// 	type args struct {
// 		uuid       string
// 		t          int64
// 		validation bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           t:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           t:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.delete(test.args.uuid, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Exists(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantOid uint32
// 		wantOk  bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotOid uint32, gotOk bool) error {
// 		if !reflect.DeepEqual(gotOid, w.wantOid) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOid, w.wantOid)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotOid, gotOk := f.Exists(test.args.uuid)
// 			if err := checkFunc(test.want, gotOid, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_GetObject(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantVec       []float32
// 		wantTimestamp int64
// 		err           error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []float32, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotTimestamp int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotVec, gotTimestamp, err := f.GetObject(test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, gotTimestamp, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_IsIndexing(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.IsIndexing()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_IsSaving(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.IsSaving()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_UUIDs(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantUuids []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotUuids []string) error {
// 		if !reflect.DeepEqual(gotUuids, w.wantUuids) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuids, w.wantUuids)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotUuids := f.UUIDs(test.args.ctx)
// 			if err := checkFunc(test.want, gotUuids); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_NumberOfCreateIndexExecution(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.NumberOfCreateIndexExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_NumberOfProactiveGCExecution(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.NumberOfProactiveGCExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_gc(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			f.gc()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Len(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.Len()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_InsertVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.InsertVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_DeleteVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.DeleteVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_GetDimensionSize(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.GetDimensionSize()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_GetTrainSize(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.GetTrainSize()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Close(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Close(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_ListObjectFunc(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		fn  func(uuid string, oid uint32, ts int64) bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			f.ListObjectFunc(test.args.ctx, test.args.fn)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_toSearchResponse(t *testing.T) {
// 	type args struct {
// 		sr []algorithm.SearchResult
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           sr:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           sr:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotRes, err := f.toSearchResponse(test.args.sr)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
