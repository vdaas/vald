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

package blob

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/db/storage/blob/s3"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
	"gocloud.dev/blob"
)

const (
	endpoint        = ""
	region          = ""
	accessKey       = ""
	secretAccessKey = ""
	bucketURL       = ""
)

func TestS3Write(t *testing.T) {
	opener, err := s3.NewSession(
		s3.WithEndpoint(endpoint),
		s3.WithRegion(region),
		s3.WithAccessKey(accessKey),
		s3.WithSecretAccessKey(secretAccessKey),
	).URLOpener()
	if err != nil {
		t.Fatalf("opener initialize failed: %s", err)
	}

	bucket, err := NewBucket(
		WithBucketURLOpener(opener),
		WithBucketURL(bucketURL),
	)
	if err != nil {
		t.Fatalf("bucket initialize failed: %s", err)
	}

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
	opener, err := s3.NewSession(
		s3.WithEndpoint(endpoint),
		s3.WithRegion(region),
		s3.WithAccessKey(accessKey),
		s3.WithSecretAccessKey(secretAccessKey),
	).URLOpener()
	if err != nil {
		t.Fatalf("opener initialize failed: %s", err)
	}

	bucket, err := NewBucket(
		WithBucketURLOpener(opener),
		WithBucketURL(bucketURL),
	)
	if err != nil {
		t.Fatalf("bucket initialize failed: %s", err)
	}

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

func TestNewBucket(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Bucket
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Bucket, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Bucket, err error) error {
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

			got, err := NewBucket(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bucket_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		opener BucketURLOpener
		url    string
		bucket *blob.Bucket
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
			b := &bucket{
				opener: test.fields.opener,
				url:    test.fields.url,
				bucket: test.fields.bucket,
			}

			err := b.Open(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bucket_Close(t *testing.T) {
	type fields struct {
		opener BucketURLOpener
		url    string
		bucket *blob.Bucket
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
			b := &bucket{
				opener: test.fields.opener,
				url:    test.fields.url,
				bucket: test.fields.bucket,
			}

			err := b.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bucket_Reader(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		opener BucketURLOpener
		url    string
		bucket *blob.Bucket
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
			b := &bucket{
				opener: test.fields.opener,
				url:    test.fields.url,
				bucket: test.fields.bucket,
			}

			got, err := b.Reader(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bucket_Writer(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		opener BucketURLOpener
		url    string
		bucket *blob.Bucket
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
		           opener: nil,
		           url: "",
		           bucket: nil,
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
			b := &bucket{
				opener: test.fields.opener,
				url:    test.fields.url,
				bucket: test.fields.bucket,
			}

			got, err := b.Writer(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
