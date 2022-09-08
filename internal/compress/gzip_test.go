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

// Package compress provides compress functions
package compress

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/compress/gzip"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

func TestNewGzip(t *testing.T) {
	type args struct {
		opts []GzipOption
	}
	type want struct {
		want Compressor
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Compressor, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Compressor, err error) error {
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
			name: "returns (Compressor, nil) when option is empty",
			args: args{
				opts: nil,
			},
			want: want{
				want: &gzipCompressor{
					gobc: func() (gob Compressor) {
						gob, _ = NewGob()
						return
					}(),
					compressionLevel: gzip.DefaultCompression,
					gzip:             gzip.New(),
				},
			},
		},

		{
			name: "returns (nil, error) when option apply fails",
			args: args{
				opts: []GzipOption{
					func(*gzipCompressor) error {
						return errors.New("err")
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

			got, err := NewGzip(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipCompressor_CompressVector(t *testing.T) {
	type args struct {
		vector []float32
	}
	type fields struct {
		gobc             Compressor
		compressionLevel int
		gzip             gzip.Gzip
	}
	type want struct {
		want []byte
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []byte, error, *gzipCompressor) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []byte, err error, _ *gzipCompressor) error {
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
			name: "returns ([]byte, nil) when no error occurs",
			args: args{
				vector: []float32{1, 2},
			},
			fields: fields{
				compressionLevel: gzip.DefaultCompression,
				gzip: &gzip.MockGzip{
					NewWriterLevelFunc: func(w io.Writer, level int) (gzip.Writer, error) {
						return &gzip.MockWriter{
							WriteFunc: w.Write,
							CloseFunc: func() error {
								return nil
							},
						}, nil
					},
				},
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return []byte("vdaas/vald"), nil
					},
				},
			},
			want: want{
				want: []byte("vdaas/vald"),
				err:  nil,
			},
		},

		{
			name: "returns (nil, error) when initialize writer level fails",
			args: args{
				vector: []float32{1, 2},
			},
			fields: fields{
				compressionLevel: gzip.DefaultCompression,
				gzip: &gzip.MockGzip{
					NewWriterLevelFunc: func(w io.Writer, level int) (gzip.Writer, error) {
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
			name: "returns (nil, error) when compress vector fails",
			args: args{
				vector: []float32{1, 2},
			},
			fields: fields{
				compressionLevel: gzip.DefaultCompression,
				gzip: &gzip.MockGzip{
					NewWriterLevelFunc: func(w io.Writer, level int) (gzip.Writer, error) {
						return new(gzip.MockWriter), nil
					},
				},
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
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
			name: "returns (nil, error) when write fails",
			args: args{
				vector: []float32{1, 2},
			},
			fields: fields{
				compressionLevel: gzip.DefaultCompression,
				gzip: &gzip.MockGzip{
					NewWriterLevelFunc: func(w io.Writer, level int) (gzip.Writer, error) {
						return &gzip.MockWriter{
							WriteFunc: func(p []byte) (n int, err error) {
								return 0, errors.New("err")
							},
						}, nil
					},
				},
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return []byte{}, nil
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("err"),
			},
		},

		{
			name: "returns (nil, error) when close fails",
			args: args{
				vector: []float32{1, 2},
			},
			fields: fields{
				compressionLevel: gzip.DefaultCompression,
				gzip: &gzip.MockGzip{
					NewWriterLevelFunc: func(w io.Writer, level int) (gzip.Writer, error) {
						return &gzip.MockWriter{
							WriteFunc: func(p []byte) (n int, err error) {
								return 10, nil
							},
							CloseFunc: func() error {
								return errors.New("err")
							},
						}, nil
					},
				},
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return []byte{}, nil
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

			g := &gzipCompressor{
				gobc:             test.fields.gobc,
				compressionLevel: test.fields.compressionLevel,
				gzip:             test.fields.gzip,
			}

			got, err := g.CompressVector(test.args.vector)
			if err := checkFunc(test.want, got, err, g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_E2E_gzipCompressor_CompressVector(t *testing.T) {
	type args struct {
		vector []float32
	}
	type want struct {
		want []float32
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []byte, error, Compressor) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []byte, err error, g Compressor) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		decompressed, err := g.DecompressVector(got)
		if err != nil {
			return errors.Errorf("decompress error: %v", err)
		}
		if !reflect.DeepEqual(decompressed, w.want) {
			return errors.Errorf("decompressed got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "compression success",
			args: args{
				vector: []float32{
					0.1, 0.2, 0.3, 0.4,
				},
			},
			want: want{
				want: []float32{
					0.1, 0.2, 0.3, 0.4,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			g, err := NewGzip()
			if err != nil {
				t.Fatal(err)
			}

			got, err := g.CompressVector(test.args.vector)
			if err := checkFunc(test.want, got, err, g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipCompressor_DecompressVector(t *testing.T) {
	type args struct {
		bs []byte
	}
	type fields struct {
		gobc Compressor
		gzip gzip.Gzip
	}
	type want struct {
		want []float32
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, error, *gzipCompressor) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float32, err error, _ *gzipCompressor) error {
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
			name: "return ([]float32, nil) when no error occurs internally",
			fields: fields{
				gzip: &gzip.MockGzip{
					NewReaderFunc: func(r io.Reader) (gzip.Reader, error) {
						return &gzip.MockReader{
							ReadFunc: func(p []byte) (n int, err error) {
								return 10, io.EOF
							},
						}, nil
					},
				},
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return []float32{1, 2, 3}, nil
					},
				},
			},
			want: want{
				want: []float32{1, 2, 3},
				err:  nil,
			},
		},

		{
			name: "return (nil, error) when initialize reader fails",
			fields: fields{
				gzip: &gzip.MockGzip{
					NewReaderFunc: func(r io.Reader) (gzip.Reader, error) {
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
			name: "return (nil, error) when copy fails",
			fields: fields{
				gzip: &gzip.MockGzip{
					NewReaderFunc: func(r io.Reader) (gzip.Reader, error) {
						return &gzip.MockReader{
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

		{
			name: "return (nil, error) when decompress vector fails",
			fields: fields{
				gzip: &gzip.MockGzip{
					NewReaderFunc: func(r io.Reader) (gzip.Reader, error) {
						return &gzip.MockReader{
							ReadFunc: func(p []byte) (n int, err error) {
								return 10, io.EOF
							},
						}, nil
					},
				},
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return nil, errors.New("err")
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
			g := &gzipCompressor{
				gobc: test.fields.gobc,
				gzip: test.fields.gzip,
			}

			got, err := g.DecompressVector(test.args.bs)
			if err := checkFunc(test.want, got, err, g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipCompressor_Reader(t *testing.T) {
	type args struct {
		src io.ReadCloser
	}
	type fields struct {
		gzip gzip.Gzip
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			var (
				src = new(gzip.MockReader)
				r   = new(gzip.MockReader)
			)
			return test{
				name: "returns (io.ReadCloser, nil) when no error occurs internally",
				args: args{
					src: src,
				},
				fields: fields{
					gzip: &gzip.MockGzip{
						NewReaderFunc: func(io.Reader) (gzip.Reader, error) {
							return r, nil
						},
					},
				},
				want: want{
					want: &gzipReader{
						src: src,
						r:   r,
					},
				},
			}
		}(),

		func() test {
			src := new(gzip.MockReader)
			return test{
				name: "returns (io.ReadCloser, nil) when no error occurs internally",
				args: args{
					src: src,
				},
				fields: fields{
					gzip: &gzip.MockGzip{
						NewReaderFunc: func(io.Reader) (gzip.Reader, error) {
							return nil, errors.New("err")
						},
					},
				},
				want: want{
					want: nil,
					err:  errors.New("err"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			g := &gzipCompressor{
				gzip: test.fields.gzip,
			}

			got, err := g.Reader(test.args.src)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipCompressor_Writer(t *testing.T) {
	type args struct {
		dst io.WriteCloser
	}
	type fields struct {
		compressionLevel int
		gzip             gzip.Gzip
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			var (
				dst = new(gzip.MockWriter)
				w   = new(gzip.MockWriter)
			)
			return test{
				name: "returns (io.WriteCloser, nil) when no error occurs internally",
				args: args{
					dst: dst,
				},
				fields: fields{
					gzip: &gzip.MockGzip{
						NewWriterLevelFunc: func(io.Writer, int) (gzip.Writer, error) {
							return w, nil
						},
					},
				},
				want: want{
					want: &gzipWriter{
						dst: dst,
						w:   w,
					},
					err: nil,
				},
			}
		}(),

		func() test {
			return test{
				name: "returns (io.WriteCloser, nil) when no error occurs internally",
				args: args{
					dst: new(gzip.MockWriter),
				},
				fields: fields{
					gzip: &gzip.MockGzip{
						NewWriterLevelFunc: func(io.Writer, int) (gzip.Writer, error) {
							return nil, errors.New("err")
						},
					},
				},
				want: want{
					want: nil,
					err:  errors.New("err"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			g := &gzipCompressor{
				compressionLevel: test.fields.compressionLevel,
				gzip:             test.fields.gzip,
			}

			got, err := g.Writer(test.args.dst)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipReader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		src io.ReadCloser
		r   io.ReadCloser
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
			name: "returns nil when read success",
			args: args{
				p: []byte{},
			},
			fields: fields{
				r: &MockReadCloser{
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
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			g := &gzipReader{
				src: test.fields.src,
				r:   test.fields.r,
			}

			gotN, err := g.Read(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipReader_Close(t *testing.T) {
	type fields struct {
		src io.ReadCloser
		r   io.ReadCloser
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
				src: &MockReadCloser{
					CloseFunc: func() error {
						return nil
					},
				},
				r: &MockReadCloser{
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
			name: "returns error when close fails",
			fields: fields{
				src: &MockReadCloser{
					CloseFunc: func() error {
						return errors.New("serr")
					},
				},
				r: &MockReadCloser{
					CloseFunc: func() error {
						return errors.New("rerr")
					},
				},
			},
			want: want{
				err: errors.Wrap(errors.New("serr"), errors.New("rerr").Error()),
			},
		},

		{
			name: "returns error when close of r fails",
			fields: fields{
				src: &MockReadCloser{
					CloseFunc: func() error {
						return nil
					},
				},
				r: &MockReadCloser{
					CloseFunc: func() error {
						return errors.New("rerr")
					},
				},
			},
			want: want{
				err: errors.Wrap(nil, errors.New("rerr").Error()),
			},
		},

		{
			name: "returns error when close of src fails",
			fields: fields{
				src: &MockReadCloser{
					CloseFunc: func() error {
						return errors.New("serr")
					},
				},
				r: &MockReadCloser{
					CloseFunc: func() error {
						return nil
					},
				},
			},
			want: want{
				err: errors.New("serr"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			g := &gzipReader{
				src: test.fields.src,
				r:   test.fields.r,
			}

			err := g.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		dst io.WriteCloser
		w   io.WriteCloser
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
			name: "returns nil when write success",
			args: args{
				p: []byte{},
			},
			fields: fields{
				w: &MockWriteCloser{
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
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			g := &gzipWriter{
				dst: test.fields.dst,
				w:   test.fields.w,
			}

			gotN, err := g.Write(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gzipWriter_Close(t *testing.T) {
	type fields struct {
		dst io.WriteCloser
		w   io.WriteCloser
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
				dst: &MockWriteCloser{
					CloseFunc: func() error {
						return nil
					},
				},
				w: &MockWriteCloser{
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
			name: "returns error when close fails",
			fields: fields{
				dst: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("derr")
					},
				},
				w: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("werr")
					},
				},
			},
			want: want{
				err: errors.Wrap(errors.New("derr"), errors.New("werr").Error()),
			},
		},

		{
			name: "returns error when close of w fails",
			fields: fields{
				dst: &MockWriteCloser{
					CloseFunc: func() error {
						return nil
					},
				},
				w: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("werr")
					},
				},
			},
			want: want{
				err: errors.Wrap(nil, errors.New("werr").Error()),
			},
		},

		{
			name: "returns error when close of dst fails",
			fields: fields{
				dst: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("derr")
					},
				},
				w: &MockWriteCloser{
					CloseFunc: func() error {
						return nil
					},
				},
			},
			want: want{
				err: errors.New("derr"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			g := &gzipWriter{
				dst: test.fields.dst,
				w:   test.fields.w,
			}

			err := g.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
