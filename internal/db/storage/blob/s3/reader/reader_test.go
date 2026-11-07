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

package reader

import (
	"bytes"
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/vdaas/vald/internal/backoff"
	ctxio "github.com/vdaas/vald/internal/db/storage/blob/s3/reader/io"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
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
		want       want
		checkFunc  func(want, Reader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
		service        s3iface.S3API
		ctxio          ctxio.IO
		bucket         string
		backoffOpts    []backoff.Option
		maxChunkSize   int64
		backoffEnabled bool
	}
	type want struct {
		err error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		hookFunc   func(io.ReadCloser)
		args       args
		name       string
		fields     fields
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := t.Context()
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
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					cancel()
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					if err := eg.Wait(); err != nil {
						t.Errorf("want: %v, but got: %v", nil, err)
					}
				},
			}
		}(),

		func() test {
			ctx := t.Context()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			var rc io.ReadCloser
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
				hookFunc: func(r io.ReadCloser) {
					rc = r
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()

					// the read error is delivered to the consumer through the stream, not to the errgroup.
					if err := eg.Wait(); err != nil {
						t.Errorf("want: %v, but got: %v", nil, err)
					}
					if _, err := readAllWithin(t, rc); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := t.Context()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			var rc io.ReadCloser
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
				hookFunc: func(r io.ReadCloser) {
					rc = r
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()

					// the read error is delivered to the consumer through the stream, not to the errgroup.
					if err := eg.Wait(); err != nil {
						t.Errorf("want: %v, but got: %v", nil, err)
					}
					if _, err := readAllWithin(t, rc); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := t.Context()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			var rc io.ReadCloser
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
				hookFunc: func(r io.ReadCloser) {
					rc = r
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()

					// the read error is delivered to the consumer through the stream, not to the errgroup.
					if err := eg.Wait(); err != nil {
						t.Errorf("want: %v, but got: %v", nil, err)
					}
					if _, err := readAllWithin(t, rc); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := t.Context()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			wantErr := errors.New("err")
			var rc io.ReadCloser
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
				hookFunc: func(r io.ReadCloser) {
					rc = r
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()

					// the read error is delivered to the consumer through the stream, not to the errgroup.
					if err := eg.Wait(); err != nil {
						t.Errorf("want: %v, but got: %v", nil, err)
					}
					if _, err := readAllWithin(t, rc); !errors.Is(err, wantErr) {
						t.Errorf("want: %v, but got: %v", wantErr, err)
					}

					cancel()
				},
			}
		}(),

		func() test {
			ctx := t.Context()
			cctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.New(ctx)

			var roopCnt atomic.Uint64
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
									if roopCnt.CompareAndSwap(0, 1) {
										return 0, io.EOF
									}
									roopCnt.Add(1)
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
				hookFunc: func(r io.ReadCloser) {
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
								if roopCnt.Load() == 0 {
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
				afterFunc: func(t *testing.T, _ args) {
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			r := &reader{
				eg:             test.fields.eg,
				service:        test.fields.service,
				bucket:         test.fields.bucket,
				ctxio:          test.fields.ctxio,
				backoffEnabled: test.fields.backoffEnabled,
				backoffOpts:    test.fields.backoffOpts,
				maxChunkSize:   test.fields.maxChunkSize,
			}

			rc, err := r.Open(test.args.ctx, test.args.key)
			if test.hookFunc != nil {
				test.hookFunc(rc)
			}
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_stream_Close(t *testing.T) {
	type fields struct {
		pr io.ReadCloser
		wg *sync.WaitGroup
	}
	type want struct {
		err error
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
		name       string
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
			r := &stream{
				pr: test.fields.pr,
				wg: test.fields.wg,
			}

			err := r.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_stream_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		pr io.ReadCloser
		wg *sync.WaitGroup
	}
	type want struct {
		err   error
		wantN int
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
			r := &stream{
				pr: test.fields.pr,
				wg: test.fields.wg,
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
		ctxio          ctxio.IO
		bucket         string
		backoffOpts    []backoff.Option
		maxChunkSize   int64
		backoffEnabled bool
	}
	type want struct {
		want io.Reader
		err  error
	}
	type test struct {
		want       want
		checkFunc  func(want, io.Reader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		fields     fields
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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
				backoffEnabled: test.fields.backoffEnabled,
				backoffOpts:    test.fields.backoffOpts,
				maxChunkSize:   test.fields.maxChunkSize,
				ctxio:          test.fields.ctxio,
			}

			got, err := r.getObjectWithBackoff(test.args.ctx, nil, test.args.key, test.args.offset, test.args.length)
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
		ctxio          ctxio.IO
		bucket         string
		backoffOpts    []backoff.Option
		maxChunkSize   int64
		backoffEnabled bool
	}
	type want struct {
		want io.Reader
		err  error
	}
	type test struct {
		want       want
		checkFunc  func(want, io.Reader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		fields     fields
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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
				ctx:    t.Context(),
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

// NOT IMPLEMENTED BELOW

// Test_reader_getObjectWithBackoff_concurrentNonRetryableErrors calls getObjectWithBackoff concurrently
// with errors that must not be retried. The errors.As targets used to be package-level variables,
// which is a data race under -race.
func Test_reader_getObjectWithBackoff_concurrentNonRetryableErrors(t *testing.T) {
	t.Parallel()

	const callers = 32
	tests := []struct {
		isWant func(error) bool
		name   string
		code   string
	}{
		{
			name:   "no such bucket is not retried",
			code:   s3.ErrCodeNoSuchBucket,
			isWant: errors.IsErrBlobNoSuchBucket,
		},
		{
			name:   "no such key is not retried",
			code:   s3.ErrCodeNoSuchKey,
			isWant: errors.IsErrBlobNoSuchKey,
		},
		{
			name: "invalid chunk range is not retried",
			code: "InvalidRange",
			isWant: func(err error) bool {
				var target *errors.ErrBlobInvalidChunkRange
				return errors.As(err, &target)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var calls atomic.Int64
			r := &reader{
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						calls.Add(1)
						return nil, awserr.New(test.code, "err", nil)
					},
				},
				backoffEnabled: true,
			}
			bo := backoff.New(backoff.WithRetryCount(1))

			start := make(chan struct{})
			errs := make(chan error, callers)
			for range callers {
				go func() {
					<-start
					_, err := r.getObjectWithBackoff(context.Background(), bo, "vald", 0, 10)
					errs <- err
				}()
			}
			close(start)

			timeout := time.After(10 * time.Second)
			for range callers {
				select {
				case err := <-errs:
					if !test.isWant(err) {
						t.Errorf("got error: %v", err)
					}
				case <-timeout:
					t.Fatal("timed out waiting for getObjectWithBackoff callers")
				}
			}
			if got := calls.Load(); got != callers {
				t.Errorf("service called %d times, want %d (no retry)", got, callers)
			}
		})
	}
}

// testTimeout bounds every blocking operation in the regression tests below.
const testTimeout = 5 * time.Second

// readAllWithin reads r until EOF or error in its own goroutine and fails the test if it blocks for testTimeout.
func readAllWithin(t *testing.T, r io.Reader) ([]byte, error) {
	t.Helper()
	type result struct {
		err error
		b   []byte
	}
	ch := make(chan result, 1)
	go func() {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r)
		ch <- result{err: err, b: buf.Bytes()}
	}()
	select {
	case res := <-ch:
		return res.b, res.err
	case <-time.After(testTimeout):
		t.Fatal("timed out reading")
		return nil, nil
	}
}

func newBody(s string) *s3.GetObjectOutput {
	return &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader([]byte(s)))}
}

func Test_reader_Read_returnsGoroutineError(t *testing.T) {
	t.Parallel()
	errGet := errors.New("get object failed")
	tests := []struct {
		name    string
		get     func(call int32) (*s3.GetObjectOutput, error)
		wantErr error
		want    string
	}{
		{
			name: "returns the error instead of io.EOF when the first chunk fails",
			get: func(int32) (*s3.GetObjectOutput, error) {
				return nil, errGet
			},
			wantErr: errGet,
		},
		{
			name: "returns the data read so far and the error when a later chunk fails",
			get: func(call int32) (*s3.GetObjectOutput, error) {
				if call == 1 {
					return newBody("abcd"), nil
				}
				return nil, errGet
			},
			wantErr: errGet,
			want:    "abcd",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			var calls atomic.Int32
			r := &reader{
				eg: errgroup.Get(),
				service: &MockS3API{
					GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
						return tc.get(calls.Add(1))
					},
				},
				ctxio:        ctxio.New(),
				maxChunkSize: 4,
			}
			rc, err := r.Open(tt.Context(), "key")
			if err != nil {
				tt.Fatal(err)
			}
			got, err := readAllWithin(tt, rc)
			if !errors.Is(err, tc.wantErr) {
				tt.Errorf("Read error = %v, want %v", err, tc.wantErr)
			}
			if string(got) != tc.want {
				tt.Errorf("read %q, want %q", got, tc.want)
			}
			if err := rc.Close(); err != nil {
				tt.Errorf("Close() = %v, want nil", err)
			}
		})
	}
}

func Test_reader_Close_waitsForGoroutine(t *testing.T) {
	t.Parallel()
	release := make(chan struct{})
	var once sync.Once
	releaseGet := func() { once.Do(func() { close(release) }) }
	defer releaseGet()
	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()
	r := &reader{
		eg: errgroup.Get(),
		service: &MockS3API{
			GetObjectWithContextFunc: func(ctx aws.Context, _ *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
				select {
				case <-release:
					return newBody("abcd"), nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			},
		},
		ctxio:        ctxio.New(),
		maxChunkSize: 4,
	}
	rc, err := r.Open(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- rc.Close() }()

	// The goroutine is still inside GetObject, so Close must not have returned yet.
	select {
	case err := <-done:
		t.Fatalf("Close returned before the goroutine finished: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	releaseGet()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Close() = %v, want nil", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("Close did not return after the goroutine finished")
	}
}

func Test_reader_getObject_closesBodyOnWrapError(t *testing.T) {
	t.Parallel()
	var closed atomic.Bool
	errWrap := errors.New("wrap failed")
	r := &reader{
		service: &MockS3API{
			GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{Body: &MockReadCloser{
					CloseFunc: func() error {
						closed.Store(true)
						return nil
					},
				}}, nil
			},
		},
		ctxio: &MockIO{
			NewReadCloserWithContextFunc: func(context.Context, io.ReadCloser) (io.ReadCloser, error) {
				return nil, errWrap
			},
		},
	}
	if _, err := r.getObject(t.Context(), "key", 0, 4); !errors.Is(err, errWrap) {
		t.Errorf("getObject() error = %v, want %v", err, errWrap)
	}
	if !closed.Load() {
		t.Error("resp.Body was not closed")
	}
}

// Test_reader_Open_streamsAreIndependent simulates a restore retry that does not wait for the previous
// attempt: the consumer of the old stream is still blocked in Read when a new Open happens, and the old
// stream is closed only afterwards. Open returns a handle per stream, so the old Read and Close must
// only use the pipe, wait group and backoff of the old stream. The new stream must deliver its full
// content, and closing it must succeed. The new stream's retry loop is held until the end, so closing
// the backoff of the new stream from the old Close would block the old Close.
func Test_reader_Open_streamsAreIndependent(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 2*testTimeout)
	defer cancel()

	var (
		releaseOld, releaseNew = make(chan struct{}), make(chan struct{})
		oldInGet, newInRetry   = make(chan struct{}), make(chan struct{})
		newCalls               atomic.Int32
	)
	r := &reader{
		eg: errgroup.Get(),
		service: &MockS3API{
			GetObjectWithContextFunc: func(ctx aws.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
				if *in.Key == "old" {
					close(oldInGet)
					select {
					case <-releaseOld:
						return newBody("abcd"), nil
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				}
				// the first call fails with a retryable error, so the second one runs inside bo.Do
				switch newCalls.Add(1) {
				case 1:
					return nil, awserr.New("Invalid", "", nil)
				case 2:
				default:
					// the next chunk request after the full chunk: empty, which ends the stream
					return newBody(""), nil
				}
				close(newInRetry)
				select {
				case <-releaseNew:
					return newBody("abcd"), nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			},
		},
		ctxio:          ctxio.New(),
		maxChunkSize:   4,
		backoffEnabled: true,
		backoffOpts:    []backoff.Option{backoff.WithInitialDuration("1ms"), backoff.WithRetryCount(3)},
	}
	wait := func(ch <-chan struct{}, what string) {
		t.Helper()
		select {
		case <-ch:
		case <-time.After(testTimeout):
			t.Fatalf("timed out waiting for %s", what)
		}
	}
	waitErr := func(ch <-chan error, what string) error {
		t.Helper()
		select {
		case err := <-ch:
			return err
		case <-time.After(testTimeout):
			t.Fatalf("timed out waiting for %s", what)
			return nil
		}
	}

	oldRC, err := r.Open(ctx, "old")
	if err != nil {
		t.Fatal(err)
	}
	wait(oldInGet, "the old stream to call GetObject")

	// The old consumer reads the old stream, which has no data yet.
	oldReadDone := make(chan error, 1)
	go func() {
		_, rerr := oldRC.Read(make([]byte, 1))
		oldReadDone <- rerr
	}()

	// The retry opens the new stream before the old one is closed.
	newRC, err := r.Open(ctx, "new")
	if err != nil {
		t.Fatal(err)
	}
	wait(newInRetry, "the new stream to enter its retry loop")

	oldCloseDone := make(chan error, 1)
	go func() { oldCloseDone <- oldRC.Close() }()

	// Closing the old stream closes only the old pipe, which ends the blocked old Read.
	if rerr := waitErr(oldReadDone, "the old Read to end"); !errors.Is(rerr, io.ErrClosedPipe) {
		t.Fatalf("Read of the old stream = %v, want %v", rerr, io.ErrClosedPipe)
	}

	// Let the old goroutine finish. The old Close must return although the new backoff is still busy.
	close(releaseOld)
	if cerr := waitErr(oldCloseDone, "Close of the old stream (it must not wait for the new stream)"); cerr != nil {
		t.Errorf("Close() of the old stream = %v, want nil", cerr)
	}

	// A later Read of the old stream must not read from the new stream. The new stream has no data yet,
	// so reading from it would block here.
	oldReadAgain := make(chan error, 1)
	go func() {
		_, rerr := oldRC.Read(make([]byte, 1))
		oldReadAgain <- rerr
	}()
	if rerr := waitErr(oldReadAgain, "the second Read of the old stream"); !errors.Is(rerr, io.ErrClosedPipe) {
		t.Errorf("second Read of the old stream = %v, want %v", rerr, io.ErrClosedPipe)
	}

	// The new stream must be intact.
	close(releaseNew)
	got, err := readAllWithin(t, newRC)
	if err != nil || string(got) != "abcd" {
		t.Errorf("read of the new stream = %q, %v, want %q, nil", got, err, "abcd")
	}
	if cerr := newRC.Close(); cerr != nil {
		t.Errorf("Close() of the new stream = %v, want nil", cerr)
	}
}

// Test_reader_Close_beforeEOFIsNotGroupError pins that a consumer closing the
// stream before reading everything (the restorer stops early when the target
// already exists) is a normal stop for the goroutine. Returning the resulting
// io.ErrClosedPipe recorded it on the errgroup, which is the process-wide group
// in the sidecar, so the daemon shutdown failed.
func Test_reader_Close_beforeEOFIsNotGroupError(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()
	eg, _ := errgroup.New(ctx)
	r := &reader{
		eg: eg,
		service: &MockS3API{
			GetObjectWithContextFunc: func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error) {
				// Nothing reads the stream, so the copy into the pipe blocks on
				// this body until Close closes the reading end.
				return newBody(string(make([]byte, 64<<10))), nil
			},
		},
		ctxio:        ctxio.New(),
		maxChunkSize: 1 << 20,
	}
	rc, err := r.Open(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	closed := make(chan error, 1)
	go func() { closed <- rc.Close() }()
	select {
	case err := <-closed:
		if err != nil {
			t.Errorf("Close() = %v, want nil", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("Close did not return")
	}
	waited := make(chan error, 1)
	go func() { waited <- eg.Wait() }()
	select {
	case err := <-waited:
		if err != nil {
			t.Errorf("eg.Wait() = %v, want nil: closing the stream early is not a failure", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("eg.Wait did not return")
	}
}
