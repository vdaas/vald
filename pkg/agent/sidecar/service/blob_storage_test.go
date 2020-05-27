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
	"io"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/compress"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNewBlobStorage(t *testing.T) {
	type args struct {
		opts []BlobStorageOption
	}
	type want struct {
		want BlobStorage
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, BlobStorage, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got BlobStorage, err error) error {
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

			got, err := NewBlobStorage(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bs_initCompressor(t *testing.T) {
	type fields struct {
		storageType       string
		bucketName        string
		filename          string
		suffix            string
		endpoint          string
		region            string
		accessKey         string
		secretAccessKey   string
		token             string
		multipartUpload   bool
		compressAlgorithm string
		compressionLevel  int
		bucket            blob.Bucket
		compressor        compress.Compressor
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
			b := &bs{
				storageType:       test.fields.storageType,
				bucketName:        test.fields.bucketName,
				filename:          test.fields.filename,
				suffix:            test.fields.suffix,
				endpoint:          test.fields.endpoint,
				region:            test.fields.region,
				accessKey:         test.fields.accessKey,
				secretAccessKey:   test.fields.secretAccessKey,
				token:             test.fields.token,
				multipartUpload:   test.fields.multipartUpload,
				compressAlgorithm: test.fields.compressAlgorithm,
				compressionLevel:  test.fields.compressionLevel,
				bucket:            test.fields.bucket,
				compressor:        test.fields.compressor,
			}

			err := b.initCompressor()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bs_initBucket(t *testing.T) {
	type fields struct {
		storageType       string
		bucketName        string
		filename          string
		suffix            string
		endpoint          string
		region            string
		accessKey         string
		secretAccessKey   string
		token             string
		multipartUpload   bool
		compressAlgorithm string
		compressionLevel  int
		bucket            blob.Bucket
		compressor        compress.Compressor
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
			b := &bs{
				storageType:       test.fields.storageType,
				bucketName:        test.fields.bucketName,
				filename:          test.fields.filename,
				suffix:            test.fields.suffix,
				endpoint:          test.fields.endpoint,
				region:            test.fields.region,
				accessKey:         test.fields.accessKey,
				secretAccessKey:   test.fields.secretAccessKey,
				token:             test.fields.token,
				multipartUpload:   test.fields.multipartUpload,
				compressAlgorithm: test.fields.compressAlgorithm,
				compressionLevel:  test.fields.compressionLevel,
				bucket:            test.fields.bucket,
				compressor:        test.fields.compressor,
			}

			err := b.initBucket()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bs_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		storageType       string
		bucketName        string
		filename          string
		suffix            string
		endpoint          string
		region            string
		accessKey         string
		secretAccessKey   string
		token             string
		multipartUpload   bool
		compressAlgorithm string
		compressionLevel  int
		bucket            blob.Bucket
		compressor        compress.Compressor
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
			b := &bs{
				storageType:       test.fields.storageType,
				bucketName:        test.fields.bucketName,
				filename:          test.fields.filename,
				suffix:            test.fields.suffix,
				endpoint:          test.fields.endpoint,
				region:            test.fields.region,
				accessKey:         test.fields.accessKey,
				secretAccessKey:   test.fields.secretAccessKey,
				token:             test.fields.token,
				multipartUpload:   test.fields.multipartUpload,
				compressAlgorithm: test.fields.compressAlgorithm,
				compressionLevel:  test.fields.compressionLevel,
				bucket:            test.fields.bucket,
				compressor:        test.fields.compressor,
			}

			got, err := b.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bs_Reader(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		storageType       string
		bucketName        string
		filename          string
		suffix            string
		endpoint          string
		region            string
		accessKey         string
		secretAccessKey   string
		token             string
		multipartUpload   bool
		compressAlgorithm string
		compressionLevel  int
		bucket            blob.Bucket
		compressor        compress.Compressor
	}
	type want struct {
		wantR io.Reader
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, io.Reader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotR io.Reader, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotR, w.wantR) {
			return errors.Errorf("got = %v, want %v", gotR, w.wantR)
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
			b := &bs{
				storageType:       test.fields.storageType,
				bucketName:        test.fields.bucketName,
				filename:          test.fields.filename,
				suffix:            test.fields.suffix,
				endpoint:          test.fields.endpoint,
				region:            test.fields.region,
				accessKey:         test.fields.accessKey,
				secretAccessKey:   test.fields.secretAccessKey,
				token:             test.fields.token,
				multipartUpload:   test.fields.multipartUpload,
				compressAlgorithm: test.fields.compressAlgorithm,
				compressionLevel:  test.fields.compressionLevel,
				bucket:            test.fields.bucket,
				compressor:        test.fields.compressor,
			}

			gotR, err := b.Reader(test.args.ctx)
			if err := test.checkFunc(test.want, gotR, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bs_Writer(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		storageType       string
		bucketName        string
		filename          string
		suffix            string
		endpoint          string
		region            string
		accessKey         string
		secretAccessKey   string
		token             string
		multipartUpload   bool
		compressAlgorithm string
		compressionLevel  int
		bucket            blob.Bucket
		compressor        compress.Compressor
	}
	type want struct {
		wantW io.WriteCloser
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, io.WriteCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotW io.WriteCloser, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotW, w.wantW) {
			return errors.Errorf("got = %v, want %v", gotW, w.wantW)
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
		           storageType: "",
		           bucketName: "",
		           filename: "",
		           suffix: "",
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           multipartUpload: false,
		           compressAlgorithm: "",
		           compressionLevel: 0,
		           bucket: nil,
		           compressor: nil,
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
			b := &bs{
				storageType:       test.fields.storageType,
				bucketName:        test.fields.bucketName,
				filename:          test.fields.filename,
				suffix:            test.fields.suffix,
				endpoint:          test.fields.endpoint,
				region:            test.fields.region,
				accessKey:         test.fields.accessKey,
				secretAccessKey:   test.fields.secretAccessKey,
				token:             test.fields.token,
				multipartUpload:   test.fields.multipartUpload,
				compressAlgorithm: test.fields.compressAlgorithm,
				compressionLevel:  test.fields.compressionLevel,
				bucket:            test.fields.bucket,
				compressor:        test.fields.compressor,
			}

			gotW, err := b.Writer(test.args.ctx)
			if err := test.checkFunc(test.want, gotW, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
