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

package s3

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/downloader"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/uploader"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantB blob.Bucket
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, blob.Bucket, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotB blob.Bucket, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotB, w.wantB) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotB, w.wantB)
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotB, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, gotB, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		eg           errgroup.Group
		bucket       string
		s3client     *s3.Client
		client       *http.Client
		logMode      aws.ClientLogMode
		dclient      downloader.Client
		uclient      *uploader.Client
		region       string
		maxPartSize  int64
		maxChunkSize int64
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
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:           test.fields.eg,
				bucket:       test.fields.bucket,
				s3client:     test.fields.s3client,
				client:       test.fields.client,
				logMode:      test.fields.logMode,
				dclient:      test.fields.dclient,
				uclient:      test.fields.uclient,
				region:       test.fields.region,
				maxPartSize:  test.fields.maxPartSize,
				maxChunkSize: test.fields.maxChunkSize,
			}

			err := c.Open(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Close(t *testing.T) {
	type fields struct {
		eg           errgroup.Group
		bucket       string
		s3client     *s3.Client
		client       *http.Client
		logMode      aws.ClientLogMode
		dclient      downloader.Client
		uclient      *uploader.Client
		region       string
		maxPartSize  int64
		maxChunkSize int64
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
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:           test.fields.eg,
				bucket:       test.fields.bucket,
				s3client:     test.fields.s3client,
				client:       test.fields.client,
				logMode:      test.fields.logMode,
				dclient:      test.fields.dclient,
				uclient:      test.fields.uclient,
				region:       test.fields.region,
				maxPartSize:  test.fields.maxPartSize,
				maxChunkSize: test.fields.maxChunkSize,
			}

			err := c.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Reader(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		eg           errgroup.Group
		bucket       string
		s3client     *s3.Client
		client       *http.Client
		logMode      aws.ClientLogMode
		dclient      downloader.Client
		uclient      *uploader.Client
		region       string
		maxPartSize  int64
		maxChunkSize int64
	}
	type want struct {
		wantRc io.ReadCloser
		err    error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, io.ReadCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRc io.ReadCloser, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRc, w.wantRc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRc, w.wantRc)
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
		           key: "",
		       },
		       fields: fields {
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
		           key: "",
		           },
		           fields: fields {
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:           test.fields.eg,
				bucket:       test.fields.bucket,
				s3client:     test.fields.s3client,
				client:       test.fields.client,
				logMode:      test.fields.logMode,
				dclient:      test.fields.dclient,
				uclient:      test.fields.uclient,
				region:       test.fields.region,
				maxPartSize:  test.fields.maxPartSize,
				maxChunkSize: test.fields.maxChunkSize,
			}

			gotRc, err := c.Reader(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, gotRc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Writer(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		eg           errgroup.Group
		bucket       string
		s3client     *s3.Client
		client       *http.Client
		logMode      aws.ClientLogMode
		dclient      downloader.Client
		uclient      *uploader.Client
		region       string
		maxPartSize  int64
		maxChunkSize int64
	}
	type want struct {
		wantWc io.WriteCloser
		err    error
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
	defaultCheckFunc := func(w want, gotWc io.WriteCloser, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotWc, w.wantWc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWc, w.wantWc)
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
		           key: "",
		       },
		       fields: fields {
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
		           key: "",
		           },
		           fields: fields {
		           eg: nil,
		           bucket: "",
		           s3client: nil,
		           client: nil,
		           logMode: nil,
		           dclient: nil,
		           uclient: nil,
		           region: "",
		           maxPartSize: 0,
		           maxChunkSize: 0,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:           test.fields.eg,
				bucket:       test.fields.bucket,
				s3client:     test.fields.s3client,
				client:       test.fields.client,
				logMode:      test.fields.logMode,
				dclient:      test.fields.dclient,
				uclient:      test.fields.uclient,
				region:       test.fields.region,
				maxPartSize:  test.fields.maxPartSize,
				maxChunkSize: test.fields.maxChunkSize,
			}

			gotWc, err := c.Writer(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, gotWc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
