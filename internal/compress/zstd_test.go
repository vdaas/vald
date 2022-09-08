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

	"github.com/vdaas/vald/internal/compress/zstd"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

var zstdCompressorComparatorOptions = []comparator.Option{
	comparator.AllowUnexported(zstdCompressor{}),
	comparator.Comparer(func(x, y gobCompressor) bool {
		return reflect.DeepEqual(x, y)
	}),
	comparator.Comparer(func(x, y zstd.EOption) bool {
		if (x == nil && y != nil) || (x != nil && y == nil) {
			return false
		}
		return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
	}),
}

func TestNewZstd(t *testing.T) {
	type args struct {
		opts []ZstdOption
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
		if diff := comparator.Diff(w.want, got, zstdCompressorComparatorOptions...); diff != "" {
			return errors.Errorf("err: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "return zstd when option is nil",
			args: args{
				opts: nil,
			},
			want: want{
				want: &zstdCompressor{
					gobc: func() Compressor {
						gobc, _ := NewGob()
						return gobc
					}(),
					eoptions: []zstd.EOption{
						zstd.WithEncoderLevel(3),
					},
					zstd: zstd.New(),
				},
			},
		},
		{
			name: "set zstd option when option is not nil",
			args: args{
				opts: []ZstdOption{
					WithZstdCompressionLevel(2),
				},
			},
			want: want{
				want: &zstdCompressor{
					gobc: func() Compressor {
						gobc, _ := NewGob()
						return gobc
					}(),
					eoptions: []zstd.EOption{
						zstd.WithEncoderLevel(3),
						zstd.WithEncoderLevel(2),
					},
					zstd: zstd.New(),
				},
			},
		},
		{
			name: "set zstd option when option return error",
			args: args{
				opts: []ZstdOption{
					func(*zstdCompressor) error {
						return errors.New("opts err")
					},
				},
			},
			want: want{
				err: errors.New("opts err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

			got, err := NewZstd(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_zstdCompressor_CompressVector(t *testing.T) {
	type args struct {
		vector []float32
	}
	type fields struct {
		gobc     Compressor
		eoptions []zstd.EOption
		zstd     zstd.Zstd
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
			name: "return compressed vector",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return ([]byte("vdaas/vald")), nil
					},
				},
				zstd: &zstd.MockZstd{
					NewWriterFunc: func(w io.Writer, opts ...zstd.EOption) (zstd.Encoder, error) {
						return &zstd.MockEncoder{
							WriteFunc: func([]byte) (int, error) {
								return 0, nil
							},
							CloseFunc: func() error {
								return nil
							},
							ReadFromFunc: func(r io.Reader) (int64, error) {
								return io.Copy(w, r)
							},
						}, nil
					},
				},
			},
			want: want{
				want: []byte("vdaas/vald"),
			},
		},
		{
			name: "return error when gobc failed to compress vector",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return nil, errors.New("gobc err")
					},
				},
			},
			want: want{
				err: errors.New("gobc err"),
			},
		},
		{
			name: "return error when writer cannot be init",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return ([]byte("vdaas/vald")), nil
					},
				},
				zstd: &zstd.MockZstd{
					NewWriterFunc: func(w io.Writer, opts ...zstd.EOption) (zstd.Encoder, error) {
						return nil, errors.New("new writer err")
					},
				},
			},
			want: want{
				err: errors.New("new writer err"),
			},
		},
		{
			name: "return error when writer cannot read from source",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return ([]byte("vdaas/vald")), nil
					},
				},
				zstd: &zstd.MockZstd{
					NewWriterFunc: func(w io.Writer, opts ...zstd.EOption) (zstd.Encoder, error) {
						return &zstd.MockEncoder{
							WriteFunc: func([]byte) (int, error) {
								return 0, nil
							},
							ReadFromFunc: func(r io.Reader) (int64, error) {
								return 0, errors.New("readFrom err")
							},
						}, nil
					},
				},
			},
			want: want{
				err: errors.New("readFrom err"),
			},
		},
		{
			name: "return error when writer cannot be closed",
			args: args{
				vector: []float32{0.1, 0.2, 0.3},
			},
			fields: fields{
				gobc: &MockCompressor{
					CompressVectorFunc: func(vector []float32) (bytes []byte, err error) {
						return ([]byte("vdaas/vald")), nil
					},
				},
				zstd: &zstd.MockZstd{
					NewWriterFunc: func(w io.Writer, opts ...zstd.EOption) (zstd.Encoder, error) {
						return &zstd.MockEncoder{
							WriteFunc: func([]byte) (int, error) {
								return 0, nil
							},
							CloseFunc: func() error {
								return errors.New("close err")
							},
							ReadFromFunc: func(r io.Reader) (int64, error) {
								return io.Copy(w, r)
							},
						}, nil
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
			defer goleak.VerifyNone(tt)
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
			z := &zstdCompressor{
				gobc:     test.fields.gobc,
				eoptions: test.fields.eoptions,
				zstd:     test.fields.zstd,
			}

			got, err := z.CompressVector(test.args.vector)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_E2E_zstdCompressor_CompressVector(t *testing.T) {
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		decompressed, err := l.DecompressVector(got)
		if err != nil {
			return errors.Errorf("decompress error: %v", err)
		}
		if !reflect.DeepEqual(decompressed, w.want) {
			return errors.Errorf("got = %v, want %v", decompressed, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns same vector after decompress the compressed data",
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
			defer goleak.VerifyNone(tt)
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

			g, err := NewZstd()
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

func Test_zstdCompressor_DecompressVector(t *testing.T) {
	type args struct {
		bs []byte
	}
	type fields struct {
		gobc     Compressor
		eoptions []zstd.EOption
		zstd     zstd.Zstd
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
			name: "return decompress data",
			args: args{
				bs: []byte("vdaas/vald"),
			},
			fields: fields{
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return []float32{0.1, 0.2, 0.3}, nil
					},
				},
				zstd: &zstd.MockZstd{
					NewReaderFunc: func(r io.Reader, opts ...zstd.DOption) (zstd.Decoder, error) {
						return &zstd.MockDecoder{
							ReadFunc: func([]byte) (int, error) {
								return 0, nil
							},
							CloseFunc: func() {
							},
							WriteToFunc: func(w io.Writer) (int64, error) {
								return io.Copy(w, r)
							},
						}, nil
					},
				},
			},
			want: want{
				want: []float32{0.1, 0.2, 0.3},
			},
		},
		{
			name: "return error when failed to init reader",
			args: args{
				bs: []byte("vdaas/vald"),
			},
			fields: fields{
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return []float32{0.1, 0.2, 0.3}, nil
					},
				},
				zstd: &zstd.MockZstd{
					NewReaderFunc: func(r io.Reader, opts ...zstd.DOption) (zstd.Decoder, error) {
						return nil, errors.New("new reader err")
					},
				},
			},
			want: want{
				err: errors.New("new reader err"),
			},
		},
		{
			name: "return error when error write to buffer",
			args: args{
				bs: []byte("vdaas/vald"),
			},
			fields: fields{
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return []float32{0.1, 0.2, 0.3}, nil
					},
				},
				zstd: &zstd.MockZstd{
					NewReaderFunc: func(r io.Reader, opts ...zstd.DOption) (zstd.Decoder, error) {
						return &zstd.MockDecoder{
							ReadFunc: func([]byte) (int, error) {
								return 0, nil
							},
							CloseFunc: func() {
							},
							WriteToFunc: func(w io.Writer) (int64, error) {
								return 0, errors.New("write to err")
							},
						}, nil
					},
				},
			},
			want: want{
				err: errors.New("write to err"),
			},
		},
		{
			name: "return error when error to decompress vecotr",
			args: args{
				bs: []byte("vdaas/vald"),
			},
			fields: fields{
				gobc: &MockCompressor{
					DecompressVectorFunc: func(bytes []byte) (vector []float32, err error) {
						return nil, errors.New("decom vec err")
					},
				},
				zstd: &zstd.MockZstd{
					NewReaderFunc: func(r io.Reader, opts ...zstd.DOption) (zstd.Decoder, error) {
						return &zstd.MockDecoder{
							ReadFunc: func([]byte) (int, error) {
								return 0, nil
							},
							CloseFunc: func() {
							},
							WriteToFunc: func(w io.Writer) (int64, error) {
								return 0, nil
							},
						}, nil
					},
				},
			},
			want: want{
				err: errors.New("decom vec err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			z := &zstdCompressor{
				gobc:     test.fields.gobc,
				eoptions: test.fields.eoptions,
				zstd:     test.fields.zstd,
			}

			got, err := z.DecompressVector(test.args.bs)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_zstdCompressor_Reader(t *testing.T) {
	type args struct {
		src io.ReadCloser
	}
	type fields struct {
		gobc     Compressor
		eoptions []zstd.EOption
		zstd     zstd.Zstd
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
			d := &zstd.MockDecoder{
				ReadFunc: func([]byte) (int, error) {
					return 0, nil
				},
				CloseFunc: func() {
				},
				WriteToFunc: func(w io.Writer) (int64, error) {
					return 0, nil
				},
			}
			return test{
				name: "return read closer",
				args: args{
					src: nil,
				},
				fields: fields{
					zstd: &zstd.MockZstd{
						NewReaderFunc: func(r io.Reader, opts ...zstd.DOption) (zstd.Decoder, error) {
							return d, nil
						},
					},
				},
				want: want{
					want: &zstdReader{
						src: nil,
						r:   d,
					},
				},
			}
		}(),
		{
			name: "return closer error when failed to init reader",
			args: args{
				src: nil,
			},
			fields: fields{
				zstd: &zstd.MockZstd{
					NewReaderFunc: func(r io.Reader, opts ...zstd.DOption) (zstd.Decoder, error) {
						return nil, errors.New("new reader err")
					},
				},
			},
			want: want{
				err: errors.New("new reader err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			z := &zstdCompressor{
				gobc:     test.fields.gobc,
				eoptions: test.fields.eoptions,
				zstd:     test.fields.zstd,
			}

			got, err := z.Reader(test.args.src)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_zstdCompressor_Writer(t *testing.T) {
	type args struct {
		dst io.WriteCloser
	}
	type fields struct {
		gobc     Compressor
		eoptions []zstd.EOption
		zstd     zstd.Zstd
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
			e := &zstd.MockEncoder{
				WriteFunc: func([]byte) (int, error) {
					return 0, nil
				},
				CloseFunc: func() error {
					return nil
				},
				ReadFromFunc: func(r io.Reader) (int64, error) {
					return 0, nil
				},
			}
			return test{
				name: "return writer",
				args: args{
					dst: nil,
				},
				fields: fields{
					zstd: &zstd.MockZstd{
						NewWriterFunc: func(w io.Writer, opts ...zstd.EOption) (zstd.Encoder, error) {
							return e, nil
						},
					},
				},
				want: want{
					want: &zstdWriter{
						dst: nil,
						w:   e,
					},
				},
			}
		}(),
		{
			name: "return error when failed to init writer",
			args: args{
				dst: nil,
			},
			fields: fields{
				zstd: &zstd.MockZstd{
					NewWriterFunc: func(w io.Writer, opts ...zstd.EOption) (zstd.Encoder, error) {
						return nil, errors.New("new writer err")
					},
				},
			},
			want: want{
				err: errors.New("new writer err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			z := &zstdCompressor{
				gobc:     test.fields.gobc,
				eoptions: test.fields.eoptions,
				zstd:     test.fields.zstd,
			}

			got, err := z.Writer(test.args.dst)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_zstdReader_Read(t *testing.T) {
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
			name: "returns n when read success",
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
		{
			name: "returns error when read failed",
			args: args{
				p: []byte{},
			},
			fields: fields{
				r: &MockReadCloser{
					ReadFunc: func(p []byte) (int, error) {
						return 0, errors.New("read err")
					},
				},
			},
			want: want{
				err: errors.New("read err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			z := &zstdReader{
				src: test.fields.src,
				r:   test.fields.r,
			}

			gotN, err := z.Read(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_zstdReader_Close(t *testing.T) {
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
			name: "returns nil when close success",
			fields: fields{
				src: &MockReadCloser{
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
			name: "returns error when failed close",
			fields: fields{
				src: &MockReadCloser{
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
			defer goleak.VerifyNone(tt)
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
			z := &zstdReader{
				src: test.fields.src,
				r:   test.fields.r,
			}

			err := z.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_zstdWriter_Write(t *testing.T) {
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
			name: "returns n when write success",
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
		{
			name: "returns error when write failed",
			args: args{
				p: []byte{},
			},
			fields: fields{
				w: &MockWriteCloser{
					WriteFunc: func(p []byte) (int, error) {
						return 0, errors.New("write err")
					},
				},
			},
			want: want{
				err: errors.New("write err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			z := &zstdWriter{
				dst: test.fields.dst,
				w:   test.fields.w,
			}

			gotN, err := z.Write(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_zstdWriter_Close(t *testing.T) {
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
			name: "returns error when w close failed",
			fields: fields{
				dst: &MockWriteCloser{
					CloseFunc: func() error {
						return nil
					},
				},
				w: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("w close err")
					},
				},
			},
			want: want{
				err: errors.New("w close err"),
			},
		},
		{
			name: "returns error when dst close failed",
			fields: fields{
				dst: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("dst close err")
					},
				},
				w: &MockWriteCloser{
					CloseFunc: func() error {
						return nil
					},
				},
			},
			want: want{
				err: errors.New("dst close err"),
			},
		},
		{
			name: "returns error when dst and close failed",
			fields: fields{
				dst: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("dst close err")
					},
				},
				w: &MockWriteCloser{
					CloseFunc: func() error {
						return errors.New("w close err")
					},
				},
			},
			want: want{
				err: errors.Wrap(errors.New("dst close err"), "w close err"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			z := &zstdWriter{
				dst: test.fields.dst,
				w:   test.fields.w,
			}

			err := z.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
