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
package cloudstorage

import (
	"context"
	"net/url"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/memblob"
	"gocloud.dev/gcerrors"
)

func Test_client_Open_invalidOptions(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("gs://bucket")
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	tests := []struct {
		client *client
		name   string
	}{
		{name: "nil url (cloud_storage.url is empty)", client: &client{urlOpener: new(gcsblob.URLOpener)}},
		{name: "nil urlOpener", client: &client{url: u}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.client.Open(t.Context())
			var target *errors.ErrInvalidOption
			if !errors.As(err, &target) {
				t.Fatalf("Open() error = %v, want *errors.ErrInvalidOption", err)
			}
			if tc.client.bucket != nil {
				t.Errorf("bucket = %v, want nil", tc.client.bucket)
			}
		})
	}
}

func Test_client_Reader_errorContract(t *testing.T) {
	t.Parallel()

	const key = "backup"

	tests := []struct {
		setup   func(*testing.T) *client
		wantErr func(error) bool
		name    string
		wantEOF bool
	}{
		{
			name: "bucket not opened returns ErrBucketNotOpened",
			setup: func(*testing.T) *client {
				return new(client)
			},
			wantErr: func(err error) bool {
				return errors.Is(err, errors.ErrBucketNotOpened)
			},
		},
		{
			name: "NotFound returns an empty reader and no error",
			setup: func(t *testing.T) *client {
				t.Helper()
				b := memblob.OpenBucket(nil)
				t.Cleanup(func() {
					if err := b.Close(); err != nil {
						t.Errorf("failed to close bucket: %v", err)
					}
				})
				return &client{bucket: b}
			},
			wantEOF: true,
		},
		{
			name: "non-NotFound error returns an error and an untyped nil reader",
			setup: func(t *testing.T) *client {
				t.Helper()
				b := memblob.OpenBucket(nil)
				// NewReader on a closed bucket fails with FailedPrecondition.
				if err := b.Close(); err != nil {
					t.Fatalf("failed to close bucket: %v", err)
				}
				return &client{bucket: b}
			},
			wantErr: func(err error) bool {
				return gcerrors.Code(err) == gcerrors.FailedPrecondition
			},
		},
		{
			name: "existing key returns its content",
			setup: func(t *testing.T) *client {
				t.Helper()
				b := memblob.OpenBucket(nil)
				t.Cleanup(func() {
					if err := b.Close(); err != nil {
						t.Errorf("failed to close bucket: %v", err)
					}
				})
				if err := b.WriteAll(t.Context(), key, []byte("data"), nil); err != nil {
					t.Fatalf("failed to write blob: %v", err)
				}
				return &client{bucket: b, readerOpts: new(blob.ReaderOptions)}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := tc.setup(t)

			rc, err := c.Reader(context.Background(), key)
			if tc.wantErr != nil {
				if err == nil || !tc.wantErr(err) {
					t.Fatalf("unexpected error: %v", err)
				}
				// compare against the untyped nil: a typed nil *blob.Reader
				// stored in the interface must not leak out.
				if rc != nil {
					t.Fatalf("reader must be an untyped nil on error, got %#v", rc)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rc == nil {
				t.Fatal("reader must not be nil on success")
			}
			defer func() {
				if cerr := rc.Close(); cerr != nil {
					t.Errorf("failed to close reader: %v", cerr)
				}
			}()
			got, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("failed to read: %v", err)
			}
			want := "data"
			if tc.wantEOF {
				want = ""
			}
			if string(got) != want {
				t.Errorf("content = %q, want %q", got, want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want iblob.Bucket
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, iblob.Bucket, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got iblob.Bucket, err error) error {
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
// 		})
// 	}
// }
//
// func Test_client_Open(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		url        *url.URL
// 		urlOpener  *gcsblob.URLOpener
// 		bucket     *blob.Bucket
// 		readerOpts *blob.ReaderOptions
// 		writerOpts *blob.WriterOptions
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
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 			c := &client{
// 				url:        test.fields.url,
// 				urlOpener:  test.fields.urlOpener,
// 				bucket:     test.fields.bucket,
// 				readerOpts: test.fields.readerOpts,
// 				writerOpts: test.fields.writerOpts,
// 			}
//
// 			err := c.Open(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_Close(t *testing.T) {
// 	type fields struct {
// 		url        *url.URL
// 		urlOpener  *gcsblob.URLOpener
// 		bucket     *blob.Bucket
// 		readerOpts *blob.ReaderOptions
// 		writerOpts *blob.WriterOptions
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
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 			c := &client{
// 				url:        test.fields.url,
// 				urlOpener:  test.fields.urlOpener,
// 				bucket:     test.fields.bucket,
// 				readerOpts: test.fields.readerOpts,
// 				writerOpts: test.fields.writerOpts,
// 			}
//
// 			err := c.Close()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_Reader(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		key string
// 	}
// 	type fields struct {
// 		url        *url.URL
// 		urlOpener  *gcsblob.URLOpener
// 		bucket     *blob.Bucket
// 		readerOpts *blob.ReaderOptions
// 		writerOpts *blob.WriterOptions
// 	}
// 	type want struct {
// 		want io.ReadCloser
// 		err  error
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
// 	defaultCheckFunc := func(w want, got io.ReadCloser, err error) error {
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
// 		           key:"",
// 		       },
// 		       fields: fields {
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 		           key:"",
// 		           },
// 		           fields: fields {
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 			c := &client{
// 				url:        test.fields.url,
// 				urlOpener:  test.fields.urlOpener,
// 				bucket:     test.fields.bucket,
// 				readerOpts: test.fields.readerOpts,
// 				writerOpts: test.fields.writerOpts,
// 			}
//
// 			got, err := c.Reader(test.args.ctx, test.args.key)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_Writer(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		key string
// 	}
// 	type fields struct {
// 		url        *url.URL
// 		urlOpener  *gcsblob.URLOpener
// 		bucket     *blob.Bucket
// 		readerOpts *blob.ReaderOptions
// 		writerOpts *blob.WriterOptions
// 	}
// 	type want struct {
// 		want io.WriteCloser
// 		err  error
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
// 	defaultCheckFunc := func(w want, got io.WriteCloser, err error) error {
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
// 		           key:"",
// 		       },
// 		       fields: fields {
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 		           key:"",
// 		           },
// 		           fields: fields {
// 		           url:nil,
// 		           urlOpener:nil,
// 		           bucket:nil,
// 		           readerOpts:nil,
// 		           writerOpts:nil,
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
// 			c := &client{
// 				url:        test.fields.url,
// 				urlOpener:  test.fields.urlOpener,
// 				bucket:     test.fields.bucket,
// 				readerOpts: test.fields.readerOpts,
// 				writerOpts: test.fields.writerOpts,
// 			}
//
// 			got, err := c.Writer(test.args.ctx, test.args.key)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
