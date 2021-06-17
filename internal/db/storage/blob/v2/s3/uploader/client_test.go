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

// Package uploader controlls upload operations for AWS S3 objects
package uploader

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/pool"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"go.uber.org/goleak"
)

func Test_rangeParam(t *testing.T) {
	type args struct {
		offset int64
		size   int64
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           offset: 0,
		           size: 0,
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
		           offset: 0,
		           size: 0,
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

			got := rangeParam(test.args.offset, test.args.size)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Client
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Client, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Client, err error) error {
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

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Open(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *s3.PutObjectInput
	}
	type fields struct {
		partSize       int64
		concurrency    int
		maxUploadParts int32
		client         manager.UploadAPIClient
		clientOptions  []func(*s3.Options)
		bufferProvider manager.ReadSeekerWriteToProvider
		copier         io.Copier
		eg             errgroup.Group
		bo             backoff.Backoff
	}
	type want struct {
		wantRc io.WriteCloser
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
	defaultCheckFunc := func(w want, gotRc io.WriteCloser, err error) error {
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
		           input: nil,
		       },
		       fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
		           input: nil,
		           },
		           fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
				partSize:       test.fields.partSize,
				concurrency:    test.fields.concurrency,
				maxUploadParts: test.fields.maxUploadParts,
				client:         test.fields.client,
				clientOptions:  test.fields.clientOptions,
				bufferProvider: test.fields.bufferProvider,
				copier:         test.fields.copier,
				eg:             test.fields.eg,
				bo:             test.fields.bo,
			}

			gotRc, err := c.Open(test.args.ctx, test.args.input)
			if err := test.checkFunc(test.want, gotRc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Upload(t *testing.T) {
	type args struct {
		ctx   context.Context
		r     io.Reader
		input *s3.PutObjectInput
	}
	type fields struct {
		partSize       int64
		concurrency    int
		maxUploadParts int32
		client         manager.UploadAPIClient
		clientOptions  []func(*s3.Options)
		bufferProvider manager.ReadSeekerWriteToProvider
		copier         io.Copier
		eg             errgroup.Group
		bo             backoff.Backoff
	}
	type want struct {
		wantWritten int64
		err         error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int64, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotWritten int64, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotWritten, w.wantWritten) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWritten, w.wantWritten)
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
		           r: nil,
		           input: nil,
		       },
		       fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
		           r: nil,
		           input: nil,
		           },
		           fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
				partSize:       test.fields.partSize,
				concurrency:    test.fields.concurrency,
				maxUploadParts: test.fields.maxUploadParts,
				client:         test.fields.client,
				clientOptions:  test.fields.clientOptions,
				bufferProvider: test.fields.bufferProvider,
				copier:         test.fields.copier,
				eg:             test.fields.eg,
				bo:             test.fields.bo,
			}

			gotWritten, err := c.Upload(test.args.ctx, test.args.r, test.args.input)
			if err := test.checkFunc(test.want, gotWritten, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_downloadChunk(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *s3.GetObjectInput
		w   *rangeWriter
	}
	type fields struct {
		partSize       int64
		concurrency    int
		maxUploadParts int32
		client         manager.UploadAPIClient
		clientOptions  []func(*s3.Options)
		bufferProvider manager.ReadSeekerWriteToProvider
		copier         io.Copier
		eg             errgroup.Group
		bo             backoff.Backoff
	}
	type want struct {
		wantN     int64
		wantTotal int64
		err       error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int64, int64, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int64, gotTotal int64, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
		}
		if !reflect.DeepEqual(gotTotal, w.wantTotal) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTotal, w.wantTotal)
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
		           in: nil,
		           w: rangeWriter{},
		       },
		       fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
		           in: nil,
		           w: rangeWriter{},
		           },
		           fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
				partSize:       test.fields.partSize,
				concurrency:    test.fields.concurrency,
				maxUploadParts: test.fields.maxUploadParts,
				client:         test.fields.client,
				clientOptions:  test.fields.clientOptions,
				bufferProvider: test.fields.bufferProvider,
				copier:         test.fields.copier,
				eg:             test.fields.eg,
				bo:             test.fields.bo,
			}

			gotN, gotTotal, err := c.downloadChunk(test.args.ctx, test.args.in, test.args.w)
			if err := test.checkFunc(test.want, gotN, gotTotal, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_tryDownloadChunk(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *s3.GetObjectInput
	}
	type fields struct {
		partSize       int64
		concurrency    int
		maxUploadParts int32
		client         manager.UploadAPIClient
		clientOptions  []func(*s3.Options)
		bufferProvider manager.ReadSeekerWriteToProvider
		copier         io.Copier
		eg             errgroup.Group
		bo             backoff.Backoff
	}
	type want struct {
		wantN     int64
		wantTotal int64
		wantW     string
		err       error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int64, int64, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int64, gotTotal int64, gotW string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
		}
		if !reflect.DeepEqual(gotTotal, w.wantTotal) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTotal, w.wantTotal)
		}
		if !reflect.DeepEqual(gotW, w.wantW) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotW, w.wantW)
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
		           in: nil,
		       },
		       fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
		           in: nil,
		           },
		           fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
				partSize:       test.fields.partSize,
				concurrency:    test.fields.concurrency,
				maxUploadParts: test.fields.maxUploadParts,
				client:         test.fields.client,
				clientOptions:  test.fields.clientOptions,
				bufferProvider: test.fields.bufferProvider,
				copier:         test.fields.copier,
				eg:             test.fields.eg,
				bo:             test.fields.bo,
			}
			w := &bytes.Buffer{}

			gotN, gotTotal, err := c.tryDownloadChunk(test.args.ctx, test.args.in, w)
			if err := test.checkFunc(test.want, gotN, gotTotal, err, w.String()); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_rangeWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		w       io.WriterAt
		offset  int64
		size    int64
		current int64
		rng     string
	}
	type want struct {
		wantN int
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           p: nil,
		       },
		       fields: fields {
		           w: nil,
		           offset: 0,
		           size: 0,
		           current: 0,
		           rng: "",
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
		           p: nil,
		           },
		           fields: fields {
		           w: nil,
		           offset: 0,
		           size: 0,
		           current: 0,
		           rng: "",
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
			r := &rangeWriter{
				w:       test.fields.w,
				offset:  test.fields.offset,
				size:    test.fields.size,
				current: test.fields.current,
				rng:     test.fields.rng,
			}

			gotN, err := r.Write(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_nextReader(t *testing.T) {
	type fields struct {
		partSize       int64
		concurrency    int
		maxUploadParts int32
		client         manager.UploadAPIClient
		clientOptions  []func(*s3.Options)
		bufferProvider manager.ReadSeekerWriteToProvider
		partPool       pool.ByteSlicePool
		copier         io.Copier
		eg             errgroup.Group
		bo             backoff.Backoff
	}
	type want struct {
		want  io.ReadSeeker
		want1 int
		want2 func()
		err   error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, io.ReadSeeker, int, func(), error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got io.ReadSeeker, got1 int, got2 func(), err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		if !reflect.DeepEqual(got2, w.want2) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got2, w.want2)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           partPool: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
		           partSize: 0,
		           concurrency: 0,
		           maxUploadParts: 0,
		           client: nil,
		           clientOptions: nil,
		           bufferProvider: nil,
		           partPool: nil,
		           copier: nil,
		           eg: nil,
		           bo: nil,
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
				partSize:       test.fields.partSize,
				concurrency:    test.fields.concurrency,
				maxUploadParts: test.fields.maxUploadParts,
				client:         test.fields.client,
				clientOptions:  test.fields.clientOptions,
				bufferProvider: test.fields.bufferProvider,
				partPool:       test.fields.partPool,
				copier:         test.fields.copier,
				eg:             test.fields.eg,
				bo:             test.fields.bo,
			}

			got, got1, got2, err := c.nextReader()
			if err := test.checkFunc(test.want, got, got1, got2, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
