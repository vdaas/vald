// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
//
// 		})
// 	}
// }
//
// func Test_faiss_initFaiss(t *testing.T) {
// 	type args struct {
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.initFaiss(test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.load(test.args.ctx, test.args.path, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.loadKVS(test.args.ctx, test.args.path, test.args.timeout)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_mktmp(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.mktmp()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.Train(test.args.nb, test.args.xb)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.Insert(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.InsertWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.insert(test.args.uuid, test.args.xb, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.Update(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.UpdateWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.update(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.readyForUpdate(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_CreateIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.CreateIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_SaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.SaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_saveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.saveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_moveAndSwitchSavedData(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.moveAndSwitchSavedData(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_CreateAndSaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.CreateAndSaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_Search(t *testing.T) {
// 	type args struct {
// 		k  uint32
// 		nq uint32
// 		xq []float32
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
// 	}
// 	type want struct {
// 		want []model.Distance
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []model.Distance, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []model.Distance, err error) error {
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
// 		           k:0,
// 		           nq:0,
// 		           xq:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           nq:0,
// 		           xq:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got, err := f.Search(test.args.k, test.args.nq, test.args.xq)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_Delete(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.Delete(test.args.uuid)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.DeleteWithTime(test.args.uuid, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.delete(test.args.uuid, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_Exists(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
// 	}
// 	type want struct {
// 		want  uint32
// 		want1 bool
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
// 	defaultCheckFunc := func(w want, got uint32, got1 bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		if !reflect.DeepEqual(got1, w.want1) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got, got1 := f.Exists(test.args.uuid)
// 			if err := checkFunc(test.want, got, got1); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_IsIndexing(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.IsIndexing()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_IsSaving(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.IsSaving()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_NumberOfCreateIndexExecution(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.NumberOfCreateIndexExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_NumberOfProactiveGCExecution(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.NumberOfProactiveGCExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_gc(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
// 	}
// 	type want struct {
// 	}
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
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
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.Len()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_InsertVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.InsertVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_DeleteVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.DeleteVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_GetDimensionSize(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.GetDimensionSize()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_GetTrainSize(t *testing.T) {
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			got := f.GetTrainSize()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_faiss_Close(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core              core.Faiss
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		fmap              map[string]int64
// 		vq                vqueue.Queue
// 		addVecs           []float32
// 		addIds            []int64
// 		isTrained         bool
// 		trainSize         int
// 		icnt              uint64
// 		indexing          atomic.Value
// 		saving            atomic.Value
// 		lastNocie         uint64
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		inMem             bool
// 		dim               int
// 		nlist             int
// 		m                 int
// 		alen              int
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		enableProactiveGC bool
// 		enableCopyOnWrite bool
// 		path              string
// 		tmpPath           atomic.Value
// 		oldPath           string
// 		basePath          string
// 		dcd               bool
// 		idelay            time.Duration
// 		kvsdbConcurrency  int
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           addVecs:nil,
// 		           addIds:nil,
// 		           isTrained:false,
// 		           trainSize:0,
// 		           icnt:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           inMem:false,
// 		           dim:0,
// 		           nlist:0,
// 		           m:0,
// 		           alen:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           dcd:false,
// 		           idelay:nil,
// 		           kvsdbConcurrency:0,
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
// 				core:              test.fields.core,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				fmap:              test.fields.fmap,
// 				vq:                test.fields.vq,
// 				addVecs:           test.fields.addVecs,
// 				addIds:            test.fields.addIds,
// 				isTrained:         test.fields.isTrained,
// 				trainSize:         test.fields.trainSize,
// 				icnt:              test.fields.icnt,
// 				indexing:          test.fields.indexing,
// 				saving:            test.fields.saving,
// 				lastNocie:         test.fields.lastNocie,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				inMem:             test.fields.inMem,
// 				dim:               test.fields.dim,
// 				nlist:             test.fields.nlist,
// 				m:                 test.fields.m,
// 				alen:              test.fields.alen,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				path:              test.fields.path,
// 				tmpPath:           test.fields.tmpPath,
// 				oldPath:           test.fields.oldPath,
// 				basePath:          test.fields.basePath,
// 				dcd:               test.fields.dcd,
// 				idelay:            test.fields.idelay,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 			}
//
// 			err := f.Close(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
