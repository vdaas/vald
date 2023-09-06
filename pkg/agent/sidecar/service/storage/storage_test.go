//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package storage provides blob storage service
package storage

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Storage
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Storage, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Storage, err error) error {
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_bs_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                        errgroup.Group
// 		storageType               string
// 		bucketName                string
// 		filename                  string
// 		suffix                    string
// 		s3Opts                    []s3.Option
// 		s3SessionOpts             []session.Option
// 		cloudStorageOpts          []cloudstorage.Option
// 		cloudStorageURLOpenerOpts []urlopener.Option
// 		compressAlgorithm         string
// 		compressionLevel          int
// 		bucket                    blob.Bucket
// 		compressor                compress.Compressor
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 			b := &bs{
// 				eg:                        test.fields.eg,
// 				storageType:               test.fields.storageType,
// 				bucketName:                test.fields.bucketName,
// 				filename:                  test.fields.filename,
// 				suffix:                    test.fields.suffix,
// 				s3Opts:                    test.fields.s3Opts,
// 				s3SessionOpts:             test.fields.s3SessionOpts,
// 				cloudStorageOpts:          test.fields.cloudStorageOpts,
// 				cloudStorageURLOpenerOpts: test.fields.cloudStorageURLOpenerOpts,
// 				compressAlgorithm:         test.fields.compressAlgorithm,
// 				compressionLevel:          test.fields.compressionLevel,
// 				bucket:                    test.fields.bucket,
// 				compressor:                test.fields.compressor,
// 			}
//
// 			got, err := b.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_bs_Stop(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		eg                        errgroup.Group
// 		storageType               string
// 		bucketName                string
// 		filename                  string
// 		suffix                    string
// 		s3Opts                    []s3.Option
// 		s3SessionOpts             []session.Option
// 		cloudStorageOpts          []cloudstorage.Option
// 		cloudStorageURLOpenerOpts []urlopener.Option
// 		compressAlgorithm         string
// 		compressionLevel          int
// 		bucket                    blob.Bucket
// 		compressor                compress.Compressor
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
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 			b := &bs{
// 				eg:                        test.fields.eg,
// 				storageType:               test.fields.storageType,
// 				bucketName:                test.fields.bucketName,
// 				filename:                  test.fields.filename,
// 				suffix:                    test.fields.suffix,
// 				s3Opts:                    test.fields.s3Opts,
// 				s3SessionOpts:             test.fields.s3SessionOpts,
// 				cloudStorageOpts:          test.fields.cloudStorageOpts,
// 				cloudStorageURLOpenerOpts: test.fields.cloudStorageURLOpenerOpts,
// 				compressAlgorithm:         test.fields.compressAlgorithm,
// 				compressionLevel:          test.fields.compressionLevel,
// 				bucket:                    test.fields.bucket,
// 				compressor:                test.fields.compressor,
// 			}
//
// 			err := b.Stop(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_bs_Reader(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                        errgroup.Group
// 		storageType               string
// 		bucketName                string
// 		filename                  string
// 		suffix                    string
// 		s3Opts                    []s3.Option
// 		s3SessionOpts             []session.Option
// 		cloudStorageOpts          []cloudstorage.Option
// 		cloudStorageURLOpenerOpts []urlopener.Option
// 		compressAlgorithm         string
// 		compressionLevel          int
// 		bucket                    blob.Bucket
// 		compressor                compress.Compressor
// 	}
// 	type want struct {
// 		wantR io.ReadCloser
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, io.ReadCloser, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotR io.ReadCloser, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotR, w.wantR) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotR, w.wantR)
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
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 			b := &bs{
// 				eg:                        test.fields.eg,
// 				storageType:               test.fields.storageType,
// 				bucketName:                test.fields.bucketName,
// 				filename:                  test.fields.filename,
// 				suffix:                    test.fields.suffix,
// 				s3Opts:                    test.fields.s3Opts,
// 				s3SessionOpts:             test.fields.s3SessionOpts,
// 				cloudStorageOpts:          test.fields.cloudStorageOpts,
// 				cloudStorageURLOpenerOpts: test.fields.cloudStorageURLOpenerOpts,
// 				compressAlgorithm:         test.fields.compressAlgorithm,
// 				compressionLevel:          test.fields.compressionLevel,
// 				bucket:                    test.fields.bucket,
// 				compressor:                test.fields.compressor,
// 			}
//
// 			gotR, err := b.Reader(test.args.ctx)
// 			if err := checkFunc(test.want, gotR, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_bs_Writer(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                        errgroup.Group
// 		storageType               string
// 		bucketName                string
// 		filename                  string
// 		suffix                    string
// 		s3Opts                    []s3.Option
// 		s3SessionOpts             []session.Option
// 		cloudStorageOpts          []cloudstorage.Option
// 		cloudStorageURLOpenerOpts []urlopener.Option
// 		compressAlgorithm         string
// 		compressionLevel          int
// 		bucket                    blob.Bucket
// 		compressor                compress.Compressor
// 	}
// 	type want struct {
// 		wantW io.WriteCloser
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, io.WriteCloser, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotW io.WriteCloser, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotW, w.wantW) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotW, w.wantW)
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
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 			b := &bs{
// 				eg:                        test.fields.eg,
// 				storageType:               test.fields.storageType,
// 				bucketName:                test.fields.bucketName,
// 				filename:                  test.fields.filename,
// 				suffix:                    test.fields.suffix,
// 				s3Opts:                    test.fields.s3Opts,
// 				s3SessionOpts:             test.fields.s3SessionOpts,
// 				cloudStorageOpts:          test.fields.cloudStorageOpts,
// 				cloudStorageURLOpenerOpts: test.fields.cloudStorageURLOpenerOpts,
// 				compressAlgorithm:         test.fields.compressAlgorithm,
// 				compressionLevel:          test.fields.compressionLevel,
// 				bucket:                    test.fields.bucket,
// 				compressor:                test.fields.compressor,
// 			}
//
// 			gotW, err := b.Writer(test.args.ctx)
// 			if err := checkFunc(test.want, gotW, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_bs_StorageInfo(t *testing.T) {
// 	type fields struct {
// 		eg                        errgroup.Group
// 		storageType               string
// 		bucketName                string
// 		filename                  string
// 		suffix                    string
// 		s3Opts                    []s3.Option
// 		s3SessionOpts             []session.Option
// 		cloudStorageOpts          []cloudstorage.Option
// 		cloudStorageURLOpenerOpts []urlopener.Option
// 		compressAlgorithm         string
// 		compressionLevel          int
// 		bucket                    blob.Bucket
// 		compressor                compress.Compressor
// 	}
// 	type want struct {
// 		want *StorageInfo
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *StorageInfo) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *StorageInfo) error {
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
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 		           eg:nil,
// 		           storageType:"",
// 		           bucketName:"",
// 		           filename:"",
// 		           suffix:"",
// 		           s3Opts:nil,
// 		           s3SessionOpts:nil,
// 		           cloudStorageOpts:nil,
// 		           cloudStorageURLOpenerOpts:nil,
// 		           compressAlgorithm:"",
// 		           compressionLevel:0,
// 		           bucket:nil,
// 		           compressor:nil,
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
// 			b := &bs{
// 				eg:                        test.fields.eg,
// 				storageType:               test.fields.storageType,
// 				bucketName:                test.fields.bucketName,
// 				filename:                  test.fields.filename,
// 				suffix:                    test.fields.suffix,
// 				s3Opts:                    test.fields.s3Opts,
// 				s3SessionOpts:             test.fields.s3SessionOpts,
// 				cloudStorageOpts:          test.fields.cloudStorageOpts,
// 				cloudStorageURLOpenerOpts: test.fields.cloudStorageURLOpenerOpts,
// 				compressAlgorithm:         test.fields.compressAlgorithm,
// 				compressionLevel:          test.fields.compressionLevel,
// 				bucket:                    test.fields.bucket,
// 				compressor:                test.fields.compressor,
// 			}
//
// 			got := b.StorageInfo()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
