//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3manager"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	os.Exit(m.Run())
}

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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns writer when option is empty",
			args: args{
				opts: nil,
			},
			want: want{
				want: &writer{
					eg:          errgroup.Get(),
					contentType: "application/octet-stream",
					maxPartSize: 64 * 1024 * 1024,
					s3manager:   s3manager.New(),
				},
			},
		},

		{
			name: "returns writer when option is not empty",
			args: args{
				opts: []Option{
					WithContentType("vdaas"),
				},
			},
			want: want{
				want: &writer{
					eg:          errgroup.Get(),
					contentType: "vdaas",
					maxPartSize: 64 * 1024 * 1024,
					s3manager:   s3manager.New(),
				},
			},
		},

		{
			name: "returns writer and outputs the warn when option is not empty and option apply fails",
			args: args{
				opts: []Option{
					func(w *writer) error {
						return errors.New("err")
					},
				},
			},
			want: want{
				want: &writer{
					eg:          errgroup.Get(),
					contentType: "application/octet-stream",
					maxPartSize: 64 * 1024 * 1024,
					s3manager:   s3manager.New(),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := New(test.args.opts...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_writer_Open(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		eg          errgroup.Group
		s3manager   s3manager.S3Manager
		service     *s3.S3
		bucket      string
		maxPartSize int64
		pw          io.WriteCloser
		wg          *sync.WaitGroup
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
		{
			name: "returns nil when no error occurs",
			args: args{
				ctx: nil,
				key: "vald",
			},
			fields: fields{
				eg: errgroup.Get(),
				s3manager: &MockS3Manager{
					NewUploaderWithClientFunc: func(s3iface.S3API, ...func(*s3manager.Uploader)) s3manager.UploadClient {
						return &MockUploadClient{
							UploadWithContextFunc: func(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
								return &s3manager.UploadOutput{
									Location: "location",
								}, nil
							},
						}
					},
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			w := &writer{
				eg:          test.fields.eg,
				s3manager:   test.fields.s3manager,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				pw:          test.fields.pw,
				wg:          test.fields.wg,
			}

			err := w.Open(test.args.ctx, test.args.key)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_writer_Close(t *testing.T) {
	type fields struct {
		eg          errgroup.Group
		service     *s3.S3
		bucket      string
		maxPartSize int64
		pw          io.WriteCloser
		wg          *sync.WaitGroup
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
			name: "returns nil when no error occurs",
			fields: fields{
				pw: &MockWriteCloser{
					CloseFunc: func() error {
						return nil
					},
				},
				wg: new(sync.WaitGroup),
			},
			want: want{
				err: nil,
			},
		},

		{
			name: "returns error when close error occurs",
			fields: fields{
				pw: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("err")
					},
				},
			},
			want: want{
				err: errors.New("err"),
			},
		},

		{
			name: "returns nil when no error occurs and writer dose not exist",
			fields: fields{
				wg: new(sync.WaitGroup),
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			w := &writer{
				eg:          test.fields.eg,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				pw:          test.fields.pw,
				wg:          test.fields.wg,
			}

			err := w.Close()
			if err := checkFunc(test.want, err); err != nil {
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
		eg          errgroup.Group
		service     *s3.S3
		bucket      string
		maxPartSize int64
		pw          io.WriteCloser
		wg          *sync.WaitGroup
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
			name: "returns nil When the write success",
			args: args{
				p: []byte{},
			},
			fields: fields{
				pw: &MockWriteCloser{
					WriteFunc: func(p []byte) (n int, err error) {
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
			name: "returns error When there is no writer",
			args: args{
				p: []byte{},
			},
			fields: fields{
				pw: nil,
			},
			want: want{
				err: errors.ErrStorageWriterNotOpened,
			},
		},

		{
			name: "returns error When the write fails",
			args: args{
				p: []byte{},
			},
			fields: fields{
				pw: &MockWriteCloser{
					WriteFunc: func(p []byte) (n int, err error) {
						return 0, errors.New("err")
					},
				},
			},
			want: want{
				err: errors.New("err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			w := &writer{
				eg:          test.fields.eg,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				pw:          test.fields.pw,
				wg:          test.fields.wg,
			}

			gotN, err := w.Write(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_writer_upload(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		body io.Reader
	}
	type fields struct {
		eg          errgroup.Group
		s3manager   s3manager.S3Manager
		service     *s3.S3
		bucket      string
		maxPartSize int64
		pw          io.WriteCloser
		wg          *sync.WaitGroup
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fieldsFunc func(*testing.T) fields
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
		{
			name: "returns nil when no error occurs",
			args: args{
				ctx:  context.Background(),
				key:  "vald",
				body: nil,
			},
			fieldsFunc: func(t *testing.T) fields {
				t.Helper()
				return fields{
					s3manager: &MockS3Manager{
						NewUploaderWithClientFunc: func(_ s3iface.S3API, opts ...func(*s3manager.Uploader)) s3manager.UploadClient {
							u := new(s3manager.Uploader)
							for _, opt := range opts {
								opt(u)
							}

							if !reflect.DeepEqual(u.PartSize, int64(100)) {
								t.Errorf("PartSize is invalid. want: %v, but got: %v", 100, u.PartSize)
							}

							return &MockUploadClient{
								UploadWithContextFunc: func(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
									return &s3manager.UploadOutput{
										Location: "location",
									}, nil
								},
							}
						},
					},
					maxPartSize: 100,
				}
			},
			want: want{
				err: nil,
			},
		},

		{
			name: "returns error when upload fails",
			args: args{
				ctx:  context.Background(),
				key:  "vald",
				body: nil,
			},
			fieldsFunc: func(t *testing.T) fields {
				t.Helper()
				return fields{
					s3manager: &MockS3Manager{
						NewUploaderWithClientFunc: func(_ s3iface.S3API, opts ...func(*s3manager.Uploader)) s3manager.UploadClient {
							u := new(s3manager.Uploader)
							for _, opt := range opts {
								opt(u)
							}

							if !reflect.DeepEqual(u.PartSize, int64(100)) {
								t.Errorf("PartSize is invalid. want: %v, but got: %v", 100, u.PartSize)
							}

							return &MockUploadClient{
								UploadWithContextFunc: func(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
									return nil, errors.New("err")
								},
							}
						},
					},
					maxPartSize: 100,
				}
			},
			want: want{
				err: errors.New("err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			fields := test.fieldsFunc(t)

			w := &writer{
				eg:          fields.eg,
				s3manager:   fields.s3manager,
				service:     fields.service,
				bucket:      fields.bucket,
				maxPartSize: fields.maxPartSize,
				pw:          fields.pw,
				wg:          fields.wg,
			}

			err := w.upload(test.args.ctx, test.args.key, test.args.body)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
