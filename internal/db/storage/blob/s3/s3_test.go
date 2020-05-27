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

package s3

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/session"
	"go.uber.org/goleak"
)

const (
	endpoint        = ""
	region          = ""
	accessKey       = ""
	secretAccessKey = ""
	bucketName      = ""
)

func TestS3Write(t *testing.T) {
	sess, err := session.New(
		session.WithEndpoint(endpoint),
		session.WithRegion(region),
		session.WithAccessKey(accessKey),
		session.WithSecretAccessKey(secretAccessKey),
	).Session()
	if err != nil {
		t.Fatalf("failed to create session: %s", err)
	}

	bucket := New(
		WithSession(sess),
		WithBucket(bucketName),
	)

	ctx := context.Background()

	err = bucket.Open(ctx)
	if err != nil {
		t.Fatalf("bucket open failed: %s", err)
	}

	defer func() {
		err = bucket.Close()
		if err != nil {
			t.Fatalf("bucket close failed: %s", err)
		}
	}()

	w, err := bucket.Writer(ctx, "writer-test.txt")
	if err != nil {
		t.Fatalf("fetch writer failed: %s", err)
	}
	defer func() {
		err = w.Close()
		if err != nil {
			t.Fatalf("writer close failed: %s", err)
		}
	}()

	_, err = w.Write([]byte("Hello from blob world!"))
	if err != nil {
		t.Fatalf("write failed: %s", err)
	}
}

func TestS3Read(t *testing.T) {
	sess, err := session.New(
		session.WithEndpoint(endpoint),
		session.WithRegion(region),
		session.WithAccessKey(accessKey),
		session.WithSecretAccessKey(secretAccessKey),
	).Session()
	if err != nil {
		t.Fatalf("failed to create session: %s", err)
	}

	bucket := New(
		WithSession(sess),
		WithBucket(bucketName),
	)

	ctx := context.Background()

	err = bucket.Open(ctx)
	if err != nil {
		t.Fatalf("bucket open failed: %s", err)
	}

	defer func() {
		err = bucket.Close()
		if err != nil {
			t.Fatalf("bucket close failed: %s", err)
		}
	}()

	r, err := bucket.Reader(ctx, "writer-test.txt")
	if err != nil {
		t.Fatalf("fetch reader failed: %s", err)
	}
	defer func() {
		err = r.Close()
		if err != nil {
			t.Fatalf("reader close failed: %s", err)
		}
	}()

	rbuf := make([]byte, 16)
	_, err = r.Read(rbuf)
	if err != nil {
		t.Fatalf("read failed: %s", err)
	}

	t.Logf("read: %s", string(rbuf))
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want blob.Bucket
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, blob.Bucket) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got blob.Bucket) error {
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

func Test_s3client_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		session         *session.Session
		service         *s3.S3
		bucket          string
		multipartUpload bool
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
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
			s := &s3client{
				session:         test.fields.session,
				service:         test.fields.service,
				bucket:          test.fields.bucket,
				multipartUpload: test.fields.multipartUpload,
			}

			err := s.Open(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_s3client_Close(t *testing.T) {
	type fields struct {
		session         *session.Session
		service         *s3.S3
		bucket          string
		multipartUpload bool
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
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
			s := &s3client{
				session:         test.fields.session,
				service:         test.fields.service,
				bucket:          test.fields.bucket,
				multipartUpload: test.fields.multipartUpload,
			}

			err := s.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_s3client_Reader(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		session         *session.Session
		service         *s3.S3
		bucket          string
		multipartUpload bool
	}
	type want struct {
		want io.ReadCloser
		err  error
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
	defaultCheckFunc := func(w want, got io.ReadCloser, err error) error {
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
		           key: "",
		       },
		       fields: fields {
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
			s := &s3client{
				session:         test.fields.session,
				service:         test.fields.service,
				bucket:          test.fields.bucket,
				multipartUpload: test.fields.multipartUpload,
			}

			got, err := s.Reader(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_s3client_Writer(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		session         *session.Session
		service         *s3.S3
		bucket          string
		multipartUpload bool
	}
	type want struct {
		want io.WriteCloser
		err  error
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
	defaultCheckFunc := func(w want, got io.WriteCloser, err error) error {
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
		           key: "",
		       },
		       fields: fields {
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
		           session: nil,
		           service: nil,
		           bucket: "",
		           multipartUpload: false,
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
			s := &s3client{
				session:         test.fields.session,
				service:         test.fields.service,
				bucket:          test.fields.bucket,
				multipartUpload: test.fields.multipartUpload,
			}

			got, err := s.Writer(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
