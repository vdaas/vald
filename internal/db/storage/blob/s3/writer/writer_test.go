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

package writer

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Writer
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Writer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Writer) error {
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

			got := New(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			err := w.Open(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_Close(t *testing.T) {
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			err := w.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got = %v, want %v", gotN, w.wantN)
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			gotN, err := w.Write(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_upload(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got = %v, want %v", gotN, w.wantN)
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			gotN, err := w.upload(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_multipartUpload(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got = %v, want %v", gotN, w.wantN)
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			gotN, err := w.multipartUpload(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_uploadPart(t *testing.T) {
	type args struct {
		p []byte
		n int64
	}
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
	}
	type want struct {
		want *s3.CompletedPart
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *s3.CompletedPart, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *s3.CompletedPart, err error) error {
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
		           p: nil,
		           n: 0,
		       },
		       fields: fields {
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           n: 0,
		           },
		           fields: fields {
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			got, err := w.uploadPart(test.args.p, test.args.n)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_abortMultipartUpload(t *testing.T) {
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			err := w.abortMultipartUpload()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_writer_completeMultipartUpload(t *testing.T) {
	type fields struct {
		service        *s3.S3
		bucket         string
		key            string
		maxPartSize    int64
		multipart      bool
		ctx            context.Context
		resp           *s3.CreateMultipartUploadOutput
		completedParts []*s3.CompletedPart
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
		           service: nil,
		           bucket: "",
		           key: "",
		           maxPartSize: 0,
		           multipart: false,
		           ctx: nil,
		           resp: nil,
		           completedParts: nil,
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
			w := &writer{
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				key:            test.fields.key,
				maxPartSize:    test.fields.maxPartSize,
				multipart:      test.fields.multipart,
				ctx:            test.fields.ctx,
				resp:           test.fields.resp,
				completedParts: test.fields.completedParts,
			}

			err := w.completeMultipartUpload()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
