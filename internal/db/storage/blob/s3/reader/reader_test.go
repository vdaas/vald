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

package reader

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/vdaas/vald/internal/backoff"
	ctxio "github.com/vdaas/vald/internal/db/storage/blob/s3/reader/io"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
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
		want Reader
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Reader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Reader, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns reader when option is empty",
			args: args{
				opts: nil,
			},
			want: want{
				want: &reader{
					eg:             errgroup.Get(),
					ctxio:          ctxio.New(),
					maxChunkSize:   512 * 1024 * 1024,
					backoffEnabled: false,
				},
			},
		},

		{
			name: "returns reader when option is not empty",
			args: args{
				opts: []Option{
					WithBackoff(true),
				},
			},
			want: want{
				want: &reader{
					eg:             errgroup.Get(),
					ctxio:          ctxio.New(),
					maxChunkSize:   512 * 1024 * 1024,
					backoffEnabled: true,
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reader_Open(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		eg             errgroup.Group
		backoffEnabled bool
		service        s3iface.S3API
		bucket         string
		pr             io.ReadCloser
		wg             *sync.WaitGroup
		ctxio          ctxio.IO
		backoffOpts    []backoff.Option
		maxChunkSize   int64
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
		afterFunc  func(args, *testing.T)
		hookFunc   func(*reader)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			return test{
				name: "returns nil when context is canceled",
				args: args{
					ctx: cctx,
					key: "vald",
				},
				fields: fields{
					eg: eg,
				},
				want: want{
					err: nil,
				},
				beforeFunc: func(args) {
					cancel()
				},
				afterFunc: func(_ args, t *testing.T) {
					t.Helper()
					if err := eg.Wait(); err != nil {
						t.Errorf("want: %v, but got: %v", nil, err)
					}
				},
			}
		}(),

		func() test {
			ctx := context.Background()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			return test{
				name: "returns nil when backoff is enabled and s3 service returns an error",
				args: args{
					ctx: cctx,
					key: "vald",
				},
				fields: fields{
					eg: eg,
					service: &MockS3API{
						GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
							return nil, wantErr
						},
					},
					backoffEnabled: true,
					backoffOpts: []backoff.Option{
						backoff.WithRetryCount(1),
					},
				},
				want: want{
					err: nil,
				},
				afterFunc: func(_ args, t *testing.T) {
					t.Helper()

					if err := eg.Wait(); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := context.Background()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			return test{
				name: "returns nil when backoff is disabled and s3 service returns an error",
				args: args{
					ctx: cctx,
				},
				fields: fields{
					eg: eg,
					service: &MockS3API{
						GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
							return nil, wantErr
						},
					},
				},
				want: want{
					err: nil,
				},
				afterFunc: func(_ args, t *testing.T) {
					t.Helper()

					if err := eg.Wait(); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := context.Background()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			return test{
				name: "returns nil when backoff is disabled and the reader creation fails",
				args: args{
					ctx: cctx,
				},
				fields: fields{
					eg: eg,
					service: &MockS3API{
						GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
							return new(s3.GetObjectOutput), nil
						},
					},
					ctxio: &MockIO{
						NewReaderWithContextFunc: func(ctx context.Context, r io.Reader) (io.Reader, error) {
							return nil, wantErr
						},
						NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
							return &MockReadCloser{
								CloseFunc: func() error {
									return nil
								},
								ReadFunc: func(p []byte) (n int, err error) {
									return 1, io.EOF
								},
							}, nil
						},
					},
				},
				want: want{
					err: nil,
				},
				afterFunc: func(_ args, t *testing.T) {
					t.Helper()

					if err := eg.Wait(); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := context.Background()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			return test{
				name: "returns nil when backoff is disabled and the reader copy fails",
				args: args{
					ctx: cctx,
				},
				fields: fields{
					eg: eg,
					service: &MockS3API{
						GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
							return new(s3.GetObjectOutput), nil
						},
					},
					ctxio: &MockIO{
						NewReaderWithContextFunc: func(ctx context.Context, r io.Reader) (io.Reader, error) {
							return &MockReadCloser{
								ReadFunc: func(p []byte) (n int, err error) {
									return 0, wantErr
								},
							}, nil
						},
						NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
							return &MockReadCloser{
								CloseFunc: func() error {
									return nil
								},
								ReadFunc: func(p []byte) (n int, err error) {
									return 0, io.EOF
								},
							}, nil
						},
					},
				},
				want: want{
					err: nil,
				},
				afterFunc: func(_ args, t *testing.T) {
					t.Helper()

					if err := eg.Wait(); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := context.Background()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			var roopCnt uint64
			return test{
				name: "returns nil when backoff is disable and multiple reads success",
				args: args{
					ctx: cctx,
				},
				fields: fields{
					eg:           eg,
					maxChunkSize: 10,
					service: &MockS3API{
						GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
							return new(s3.GetObjectOutput), nil
						},
					},
					ctxio: &MockIO{
						NewReaderWithContextFunc: func(ctx context.Context, r io.Reader) (io.Reader, error) {
							return &MockReadCloser{
								ReadFunc: func(p []byte) (n int, err error) {
									if atomic.CompareAndSwapUint64(&roopCnt, 0, 1) {
										return 0, io.EOF
									}
									atomic.AddUint64(&roopCnt, 1)
									return 10, io.EOF
								},
							}, nil
						},
						NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
							return &MockReadCloser{
								CloseFunc: func() error {
									return nil
								},
								ReadFunc: func(p []byte) (n int, err error) {
									return 10, io.EOF
								},
							}, nil
						},
					},
				},
				want: want{
					err: nil,
				},
				hookFunc: func(r *reader) {
					go func() {
						bytes := [][]byte{
							make([]byte, 10),
							make([]byte, 0),
						}
						for {
							select {
							case <-cctx.Done():
								return
							default:
								if atomic.LoadUint64(&roopCnt) == 0 {
									if _, err := r.Read(bytes[0]); errors.Is(err, io.EOF) {
										return
									}
								} else {
									if _, err := r.Read(bytes[1]); errors.Is(err, io.EOF) {
										return
									}
								}
							}
						}
					}()
				},
				afterFunc: func(_ args, t *testing.T) {
					t.Helper()

					if err := eg.Wait(); err != nil {
						t.Errorf("want: %v, but got: %v", nil, err)
					}

					cancel()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, t)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			r := &reader{
				eg:             test.fields.eg,
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				pr:             test.fields.pr,
				wg:             test.fields.wg,
				ctxio:          test.fields.ctxio,
				backoffEnabled: test.fields.backoffEnabled,
				backoffOpts:    test.fields.backoffOpts,
				maxChunkSize:   test.fields.maxChunkSize,
			}

			err := r.Open(test.args.ctx, test.args.key)
			if test.hookFunc != nil {
				test.hookFunc(r)
			}
			if err := checkFunc(test.want, err); err != nil {
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
			name: "returns nil when the close is nil",
			fields: fields{
				wg: new(sync.WaitGroup),
			},
			want: want{
				err: nil,
			},
		},

		{
			name: "returns nil when the close fails",
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
			r := &reader{
				eg:      test.fields.eg,
				service: test.fields.service,
				bucket:  test.fields.bucket,
				pr:      test.fields.pr,
				wg:      test.fields.wg,
			}

			err := r.Close()
			if err := checkFunc(test.want, err); err != nil {
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
			name: "returns error when read is nil",
			args: args{
				p: []byte{},
			},
			want: want{
				wantN: 0,
				err:   errors.ErrStorageReaderNotOpened,
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
			r := &reader{
				eg:      test.fields.eg,
				service: test.fields.service,
				bucket:  test.fields.bucket,
				pr:      test.fields.pr,
				wg:      test.fields.wg,
			}

			gotN, err := r.Read(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reader_getObjectWithBackoff(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		offset int64
		length int64
	}
	type fields struct {
		eg             errgroup.Group
		service        s3iface.S3API
		bucket         string
		pr             io.ReadCloser
		wg             *sync.WaitGroup
		backoffEnabled bool
		backoffOpts    []backoff.Option
		maxChunkSize   int64
		ctxio          ctxio.IO
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
		{
			name: "returns (Reader, nil) when no error occurs",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 1,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return new(s3.GetObjectOutput), nil
					},
				},
				ctxio: &MockIO{
					NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
						return &MockReadCloser{
							CloseFunc: func() error {
								return nil
							},
							ReadFunc: func(p []byte) (n int, err error) {
								return 1, io.EOF
							},
						}, nil
					},
				},
			},
			want: want{
				want: func() *bytes.Buffer {
					buf := new(bytes.Buffer)
					var b byte
					buf.WriteByte(b)
					return buf
				}(),
				err: nil,
			},
		},

		{
			name: "returns error when s3 service returns error and backoff fails",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 1,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return nil, errors.New("err")
					},
				},
				backoffEnabled: false,
				backoffOpts: []backoff.Option{
					backoff.WithRetryCount(1),
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
			r := &reader{
				eg:             test.fields.eg,
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				pr:             test.fields.pr,
				wg:             test.fields.wg,
				backoffEnabled: test.fields.backoffEnabled,
				backoffOpts:    test.fields.backoffOpts,
				maxChunkSize:   test.fields.maxChunkSize,
				ctxio:          test.fields.ctxio,
			}

			got, err := r.getObjectWithBackoff(test.args.ctx, test.args.key, test.args.offset, test.args.length)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reader_getObject(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		offset int64
		length int64
	}
	type fields struct {
		eg             errgroup.Group
		service        s3iface.S3API
		bucket         string
		pr             io.ReadCloser
		wg             *sync.WaitGroup
		backoffEnabled bool
		backoffOpts    []backoff.Option
		maxChunkSize   int64
		ctxio          ctxio.IO
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
		{
			name: "returns (Reader, nil) when no error occurs",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return new(s3.GetObjectOutput), nil
					},
				},
				ctxio: &MockIO{
					NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
						return &MockReadCloser{
							CloseFunc: func() error {
								return nil
							},
							ReadFunc: func(p []byte) (n int, err error) {
								return 1, io.EOF
							},
						}, nil
					},
				},
			},
			want: want{
				want: func() *bytes.Buffer {
					buf := new(bytes.Buffer)
					var b byte
					buf.WriteByte(b)
					return buf
				}(),
				err: nil,
			},
		},

		{
			name: "returns (Reader, nil) when the reader close error occurs and output warning",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return new(s3.GetObjectOutput), nil
					},
				},
				ctxio: &MockIO{
					NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
						return &MockReadCloser{
							CloseFunc: func() error {
								return errors.New("err")
							},
							ReadFunc: func(p []byte) (n int, err error) {
								return 1, io.EOF
							},
						}, nil
					},
				},
			},
			want: want{
				want: func() *bytes.Buffer {
					buf := new(bytes.Buffer)
					var b byte
					buf.WriteByte(b)
					return buf
				}(),
				err: nil,
			},
		},

		{
			name: "returns ErrBlobNoSuchBucket when s3 service returns error and error code is ErrBlobNoSuchBucket",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return nil, awserr.New(s3.ErrCodeNoSuchBucket, "", nil)
					},
				},
				bucket: "vald",
			},
			want: want{
				want: nil,
				err:  errors.NewErrBlobNoSuchBucket(awserr.New(s3.ErrCodeNoSuchBucket, "", nil), "vald"),
			},
		},

		{
			name: "returns nil when s3 service returns error and error code is ErrCodeNoSuchKey",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return nil, awserr.New(s3.ErrCodeNoSuchKey, "", nil)
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.NewErrBlobNoSuchKey(awserr.New(s3.ErrCodeNoSuchKey, "", nil), "vald"),
			},
		},

		{
			name: "returns nil when s3 service returns error and error code is ErrCodeNoSuchKey",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return nil, awserr.New("InvalidRange", "", nil)
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.NewErrBlobInvalidChunkRange(awserr.New("InvalidRange", "", nil), "bytes=2-11"),
			},
		},

		{
			name: "returns s3 error when s3 service returns error and error code is `Invalid`",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return nil, awserr.New("Invalid", "", nil)
					},
				},
			},
			want: want{
				want: nil,
				err:  awserr.New("Invalid", "", nil),
			},
		},

		{
			name: "returns error when reader creation fails",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return new(s3.GetObjectOutput), nil
					},
				},
				ctxio: &MockIO{
					NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
						return nil, errors.New("err")
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("err"),
			},
		},

		{
			name: "returns error when failed to copy to buffer",
			args: args{
				ctx:    context.Background(),
				key:    "vald",
				offset: 2,
				length: 10,
			},
			fields: fields{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return new(s3.GetObjectOutput), nil
					},
				},
				ctxio: &MockIO{
					NewReadCloserWithContextFunc: func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
						return &MockReadCloser{
							CloseFunc: func() error {
								return nil
							},
							ReadFunc: func(p []byte) (n int, err error) {
								return 0, errors.New("err")
							},
						}, nil
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("err"),
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
			r := &reader{
				eg:             test.fields.eg,
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				pr:             test.fields.pr,
				wg:             test.fields.wg,
				backoffEnabled: test.fields.backoffEnabled,
				backoffOpts:    test.fields.backoffOpts,
				maxChunkSize:   test.fields.maxChunkSize,
				ctxio:          test.fields.ctxio,
			}

			got, err := r.getObject(test.args.ctx, test.args.key, test.args.offset, test.args.length)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
