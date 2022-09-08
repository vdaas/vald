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

	"github.com/vdaas/vald/internal/compress/lz4"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

func TestNewLZ4(t *testing.T) {
	type args struct {
		opts []LZ4Option
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
				want: &lz4Compressor{
					gobc: func() (gob Compressor) {
						gob, _ = NewGob()
						return
					}(),
					compressionLevel: 0,
					lz4:              lz4.New(),
				},
			},
		},
		{
			name: "returns (Compressor, nil) when option is not empty",
			args: args{
				opts: []LZ4Option{
					WithLZ4CompressionLevel(-1),
				},
			},
			want: want{
				want: &lz4Compressor{
					gobc: func() (gob Compressor) {
						gob, _ = NewGob()
						return
					}(),
					compressionLevel: -1,
					lz4:              lz4.New(),
				},
			},
		},
		{
			name: "returns (nil, error) when option apply fails",
			args: args{
				opts: []LZ4Option{
					func(*lz4Compressor) error {
						return errors.New("opts err")
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("opts err"),
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

			got, err := NewLZ4(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Compressor_CompressVector(t *testing.T) {
	type args struct {
		vector []float32
	}
	type fields struct {
		gobc             Compressor
		compressionLevel int
		lz4              lz4.LZ4
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
		checkFunc  func(want, []byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []byte, err error) error {
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
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return []byte("vdaas/vald"), nil
					},
				},
				compressionLevel: 0,
				lz4: &lz4.MockLZ4{
					NewWriterLevelFunc: func(w io.Writer, level int) lz4.Writer {
						return &lz4.MockWriter{
							WriteFunc: w.Write,
							FlushFunc: func() error {
								return nil
							},
							CloseFunc: func() error {
								return nil
							},
						}
					},
				},
			},
			want: want{
				want: []byte("vdaas/vald"),
			},
		},
		{
			name: "returns (nil, error) when compress vector fails",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return nil, errors.New("compress err")
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("compress err"),
			},
		},
		{
			name: "returns (nil, error) when Write fails",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return []byte("vdaas/vald"), nil
					},
				},
				compressionLevel: 0,
				lz4: &lz4.MockLZ4{
					NewWriterLevelFunc: func(w io.Writer, level int) lz4.Writer {
						return &lz4.MockWriter{
							WriteFunc: func(p []byte) (int, error) {
								return 0, errors.New("write err")
							},
							CloseFunc: func() error {
								return nil
							},
						}
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("write err"),
			},
		},
		{
			name: "returns (nil, error) when Flush fails",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return []byte("vdaas/vald"), nil
					},
				},
				compressionLevel: 0,
				lz4: &lz4.MockLZ4{
					NewWriterLevelFunc: func(w io.Writer, level int) lz4.Writer {
						return &lz4.MockWriter{
							WriteFunc: w.Write,
							FlushFunc: func() error {
								return errors.New("flush err")
							},
							CloseFunc: func() error {
								return nil
							},
						}
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("flush err"),
			},
		},
		{
			name: "returns (nil, error) when Close fails",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return []byte("vdaas/vald"), nil
					},
				},
				compressionLevel: 0,
				lz4: &lz4.MockLZ4{
					NewWriterLevelFunc: func(w io.Writer, level int) lz4.Writer {
						return &lz4.MockWriter{
							WriteFunc: w.Write,
							FlushFunc: func() error {
								return nil
							},
							CloseFunc: func() error {
								return errors.New("close err")
							},
						}
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("close err"),
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
			l := &lz4Compressor{
				gobc:             test.fields.gobc,
				compressionLevel: test.fields.compressionLevel,
				lz4:              test.fields.lz4,
			}

			got, err := l.CompressVector(test.args.vector)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_E2E_lz4Compressor_CompressVector(t *testing.T) {
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
	defaultCheckFunc := func(w want, got []byte, err error, l Compressor) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		decompressed, err := l.DecompressVector(got)
		if err != nil {
			return errors.Errorf("decompress error: %v", err)
		}
		if !reflect.DeepEqual(decompressed, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", decompressed, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns ([]byte, nil) when no error occurs",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			want: want{
				want: []float32{0.1, 0.2, 0.3},
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

			l, err := NewLZ4()
			if err != nil {
				t.Fatal(err)
			}

			got, err := l.CompressVector(test.args.vector)
			if err := checkFunc(test.want, got, err, l); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Compressor_DecompressVector(t *testing.T) {
	type args struct {
		bs []byte
	}
	type fields struct {
		gobc             Compressor
		compressionLevel int
		lz4              lz4.LZ4
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
		checkFunc  func(want, []float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float32, err error) error {
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
			name: "returns ([]float32, nil) when no error occurs",
			args: args{
				bs: []byte("vdaas/vald"),
			},
			fields: fields{
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return []float32{0.1, 0.2, 0.3}, nil
					},
				},
				compressionLevel: 0,
				lz4: &lz4.MockLZ4{
					NewReaderFunc: func(r io.Reader) lz4.Reader {
						return &lz4.MockReader{
							ReadFunc: r.Read,
						}
					},
				},
			},
			want: want{
				want: []float32{0.1, 0.2, 0.3},
			},
		},
		{
			name: "returns (nil, error) when Copy fails",
			args: args{
				bs: []byte("vdaas/vald"),
			},
			fields: fields{
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return []float32{0.1, 0.2, 0.3}, nil
					},
				},
				compressionLevel: 0,
				lz4: &lz4.MockLZ4{
					NewReaderFunc: func(r io.Reader) lz4.Reader {
						return &lz4.MockReader{
							ReadFunc: func(p []byte) (int, error) {
								return 0, errors.New("copy err")
							},
						}
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("copy err"),
			},
		},
		{
			name: "returns (nil, error) when decompresse fails",
			args: args{
				bs: []byte("vdaas/vald"),
			},
			fields: fields{
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return nil, errors.New("decompresse err")
					},
				},
				compressionLevel: 0,
				lz4: &lz4.MockLZ4{
					NewReaderFunc: func(r io.Reader) lz4.Reader {
						return &lz4.MockReader{
							ReadFunc: r.Read,
						}
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.New("decompresse err"),
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
			l := &lz4Compressor{
				gobc:             test.fields.gobc,
				compressionLevel: test.fields.compressionLevel,
				lz4:              test.fields.lz4,
			}

			got, err := l.DecompressVector(test.args.bs)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Compressor_Reader(t *testing.T) {
	type args struct {
		src io.ReadCloser
	}
	type fields struct {
		gobc             Compressor
		compressionLevel int
		lz4              lz4.LZ4
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
				src = new(lz4Reader)
				r   = new(lz4Reader)
			)
			return test{
				name: "returns (io.ReadCloser, nil) when no error occurs",
				args: args{
					src: src,
				},
				fields: fields{
					lz4: &lz4.MockLZ4{
						NewReaderFunc: func(io.Reader) lz4.Reader {
							return r
						},
					},
				},
				want: want{
					want: &lz4Reader{
						src: src,
						r:   r,
					},
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
			l := &lz4Compressor{
				gobc:             test.fields.gobc,
				compressionLevel: test.fields.compressionLevel,
				lz4:              test.fields.lz4,
			}

			got, err := l.Reader(test.args.src)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Compressor_Writer(t *testing.T) {
	type args struct {
		dst io.WriteCloser
	}
	type fields struct {
		gobc             Compressor
		compressionLevel int
		lz4              lz4.LZ4
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
				dst = new(lz4.MockWriter)
				w   = new(lz4.MockWriter)
			)
			return test{
				name: "returns (io.WriteCloser, nil) when no erro occurs",
				args: args{
					dst: dst,
				},
				fields: fields{
					lz4: &lz4.MockLZ4{
						NewWriterFunc: func(io.Writer) lz4.Writer {
							return w
						},
					},
				},
				want: want{
					want: &lz4Writer{
						dst: dst,
						w:   w,
					},
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
			l := &lz4Compressor{
				gobc:             test.fields.gobc,
				compressionLevel: test.fields.compressionLevel,
				lz4:              test.fields.lz4,
			}

			got, err := l.Writer(test.args.dst)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Reader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		src io.ReadCloser
		r   io.Reader
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
			name: "returns (n, nil) when read success",
			args: args{
				p: []byte{},
			},
			fields: fields{
				r: &MockReadCloser{
					ReadFunc: func(p []byte) (int, error) {
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
			l := &lz4Reader{
				src: test.fields.src,
				r:   test.fields.r,
			}

			gotN, err := l.Read(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Reader_Close(t *testing.T) {
	type fields struct {
		src io.ReadCloser
		r   io.Reader
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
			name: "returns nil when readClose success",
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
			l := &lz4Reader{
				src: test.fields.src,
				r:   test.fields.r,
			}

			err := l.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Writer_Write(t *testing.T) {
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
			name: "returns (n, nil) when write success",
			args: args{
				p: []byte{},
			},
			fields: fields{
				w: &MockWriteCloser{
					WriteFunc: func(p []byte) (int, error) {
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
			l := &lz4Writer{
				dst: test.fields.dst,
				w:   test.fields.w,
			}

			gotN, err := l.Write(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_lz4Writer_Close(t *testing.T) {
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
			name: "returns nil when writeClose success",
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
						return nil
					},
				},
				w: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("close err")
					},
				},
			},
			want: want{
				err: errors.New("close err"),
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
			l := &lz4Writer{
				dst: test.fields.dst,
				w:   test.fields.w,
			}

			err := l.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
