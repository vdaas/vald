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

package reader

import (
	"context"
	"io"
	"reflect"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Reader
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Reader) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Reader) error {
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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

			got := New(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reader_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		eg      errgroup.Group
		service *s3.S3
		bucket  string
		key     string
		pr      io.ReadCloser
		wg      *sync.WaitGroup
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
		           service: nil,
		           bucket: "",
		           key: "",
		           pr: nil,
		           wg: sync.WaitGroup{},
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
		           service: nil,
		           bucket: "",
		           key: "",
		           pr: nil,
		           wg: sync.WaitGroup{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			r := &reader{
				eg:      test.fields.eg,
				service: test.fields.service,
				bucket:  test.fields.bucket,
				key:     test.fields.key,
				pr:      test.fields.pr,
				wg:      test.fields.wg,
			}

			err := r.Open(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reader_Close(t *testing.T) {
	type fields struct {
		eg      errgroup.Group
		service *s3.S3
		bucket  string
		key     string
		pr      io.ReadCloser
		wg      *sync.WaitGroup
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
		{
			name: "returns nil when close success",
			fields: fields{
				pr: &MockReadCloser{
					CloseFunc: func() error {
						return nil
					},
				},
			},
			want: want{
				err: nil,
			},
		},

		{
			name: "returns nil when close is nil",
			fields: fields{
				wg: new(sync.WaitGroup),
			},
			want: want{
				err: nil,
			},
		},

		{
			name: "returns nil when close fails",
			fields: fields{
				pr: &MockReadCloser{
					CloseFunc: func() error {
						return errors.New("err")
					},
				},
			},
			want: want{
				err: errors.New("err"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reader{
				eg:      test.fields.eg,
				service: test.fields.service,
				bucket:  test.fields.bucket,
				key:     test.fields.key,
				pr:      test.fields.pr,
				wg:      test.fields.wg,
			}

			err := r.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		eg      errgroup.Group
		service *s3.S3
		bucket  string
		key     string
		pr      io.ReadCloser
		wg      *sync.WaitGroup
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
		{
			name: "returns (10, nil) when read success",
			args: args{
				p: []byte{},
			},
			fields: fields{
				pr: &MockReadCloser{
					ReadFunc: func(p []byte) (n int, err error) {
						return 10, nil
					},
				},
			},
			want: want{
				wantN: 10,
				err:   nil,
			},
		},

		{
			name: "returns error when read fails",
			args: args{
				p: []byte{},
			},
			fields: fields{
				pr: &MockReadCloser{
					ReadFunc: func(p []byte) (n int, err error) {
						return 0, errors.New("err")
					},
				},
			},
			want: want{
				wantN: 0,
				err:   errors.New("err"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reader{
				eg:      test.fields.eg,
				service: test.fields.service,
				bucket:  test.fields.bucket,
				key:     test.fields.key,
				pr:      test.fields.pr,
				wg:      test.fields.wg,
			}

			gotN, err := r.Read(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reader_getObjectWithBackoff(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset int64
		length int64
	}
	type fields struct {
		eg             errgroup.Group
		service        *s3.S3
		bucket         string
		key            string
		pr             io.ReadCloser
		wg             *sync.WaitGroup
		backoffEnabled bool
		backoffOpts    []backoff.Option
		maxChunkSize   int64
	}
	type want struct {
		want io.Reader
		err  error
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
	defaultCheckFunc := func(w want, got io.Reader, err error) error {
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
		           ctx: nil,
		           offset: 0,
		           length: 0,
		       },
		       fields: fields {
		           eg: nil,
		           service: nil,
		           bucket: "",
		           key: "",
		           pr: nil,
		           wg: nil,
		           backoffEnabled: false,
		           backoffOpts: nil,
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
		           offset: 0,
		           length: 0,
		           },
		           fields: fields {
		           eg: nil,
		           service: nil,
		           bucket: "",
		           key: "",
		           pr: nil,
		           wg: nil,
		           backoffEnabled: false,
		           backoffOpts: nil,
		           maxChunkSize: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			r := &reader{
				eg:             test.fields.eg,
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				pr:             test.fields.pr,
				wg:             test.fields.wg,
				backoffEnabled: test.fields.backoffEnabled,
				backoffOpts:    test.fields.backoffOpts,
				maxChunkSize:   test.fields.maxChunkSize,
			}

			got, err := r.getObjectWithBackoff(test.args.ctx, test.args.offset, test.args.length)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reader_getObject(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset int64
		length int64
	}
	type fields struct {
		eg             errgroup.Group
		service        *s3.S3
		bucket         string
		key            string
		pr             io.ReadCloser
		wg             *sync.WaitGroup
		backoffEnabled bool
		backoffOpts    []backoff.Option
		maxChunkSize   int64
	}
	type want struct {
		want io.Reader
		err  error
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
	defaultCheckFunc := func(w want, got io.Reader, err error) error {
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
		           ctx: nil,
		           offset: 0,
		           length: 0,
		       },
		       fields: fields {
		           eg: nil,
		           service: nil,
		           bucket: "",
		           key: "",
		           pr: nil,
		           wg: nil,
		           backoffEnabled: false,
		           backoffOpts: nil,
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
		           offset: 0,
		           length: 0,
		           },
		           fields: fields {
		           eg: nil,
		           service: nil,
		           bucket: "",
		           key: "",
		           pr: nil,
		           wg: nil,
		           backoffEnabled: false,
		           backoffOpts: nil,
		           maxChunkSize: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			r := &reader{
				eg:             test.fields.eg,
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				pr:             test.fields.pr,
				wg:             test.fields.wg,
				backoffEnabled: test.fields.backoffEnabled,
				backoffOpts:    test.fields.backoffOpts,
				maxChunkSize:   test.fields.maxChunkSize,
			}

			got, err := r.getObject(test.args.ctx, test.args.offset, test.args.length)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
