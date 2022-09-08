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

	"github.com/vdaas/vald/internal/compress/gob"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewGob(t *testing.T) {
	type args struct {
		opts []GobOption
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
				want: &gobCompressor{
					transcoder: gob.New(),
				},
			},
		},

		{
			name: "returns (nil, error) when option is not nil and option apply fails",
			args: args{
				opts: []GobOption{
					func(c *gobCompressor) error {
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

			got, err := NewGob(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobCompressor_CompressVector(t *testing.T) {
	type args struct {
		vector []float32
	}
	type fields struct {
		transcoder gob.Transcoder
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
				vector: []float32{
					1, 2, 3,
				},
			},
			fields: fields{
				transcoder: &gob.MockTranscoder{
					NewEncoderFunc: func(w io.Writer) gob.Encoder {
						return &gob.MockEncoder{
							EncodeFunc: func(e interface{}) error {
								_, _ = w.Write([]byte("vald"))
								return nil
							},
						}
					},
				},
			},
			want: want{
				want: []byte("vald"),
				err:  nil,
			},
		},

		{
			name: "returns (nil, error) when decode fails",
			args: args{
				vector: []float32{
					1, 2, 3,
				},
			},
			fields: fields{
				transcoder: &gob.MockTranscoder{
					NewEncoderFunc: func(w io.Writer) gob.Encoder {
						return &gob.MockEncoder{
							EncodeFunc: func(e interface{}) error {
								return errors.New("err")
							},
						}
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
			g := &gobCompressor{
				transcoder: test.fields.transcoder,
			}

			got, err := g.CompressVector(test.args.vector)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_E2E_gobCompressor_CompressVector(t *testing.T) {
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

			g, err := NewGob()
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

func Test_gobCompressor_DecompressVector(t *testing.T) {
	type args struct {
		bs []byte
	}
	type fields struct {
		transcoder gob.Transcoder
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
				bs: []byte{},
			},
			fields: fields{
				transcoder: &gob.MockTranscoder{
					NewDecoderFunc: func(io.Reader) gob.Decoder {
						return &gob.MockDecoder{
							DecodeFunc: func(e interface{}) error {
								reflect.ValueOf(e).Elem().Set(reflect.ValueOf(&[]float32{
									1, 2, 3,
								}).Elem())
								return nil
							},
						}
					},
				},
			},
			want: want{
				want: []float32{
					1, 2, 3,
				},
				err: nil,
			},
		},

		{
			name: "returns (nil, error) when decode fails",
			args: args{
				bs: []byte{},
			},
			fields: fields{
				transcoder: &gob.MockTranscoder{
					NewDecoderFunc: func(io.Reader) gob.Decoder {
						return &gob.MockDecoder{
							DecodeFunc: func(interface{}) error {
								return errors.New("err")
							},
						}
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
			g := &gobCompressor{
				transcoder: test.fields.transcoder,
			}

			got, err := g.DecompressVector(test.args.bs)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobCompressor_Reader(t *testing.T) {
	type args struct {
		src io.ReadCloser
	}
	type fields struct {
		transcodr gob.Transcoder
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
			rc := new(MockReadCloser)
			dec := new(gob.MockDecoder)
			return test{
				name: "returns (io.ReadCloser, nil)",
				args: args{
					src: rc,
				},
				fields: fields{
					transcodr: &gob.MockTranscoder{
						NewDecoderFunc: func(r io.Reader) gob.Decoder {
							return dec
						},
					},
				},
				want: want{
					want: &gobReader{
						src:     rc,
						decoder: dec,
					},
					err: nil,
				},
			}
		}(),
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
			g := &gobCompressor{
				transcoder: test.fields.transcodr,
			}

			got, err := g.Reader(test.args.src)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobCompressor_Writer(t *testing.T) {
	type args struct {
		dst io.WriteCloser
	}
	type fields struct {
		transcoder gob.Transcoder
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
			wc := new(MockWriteCloser)
			enc := new(gob.MockEncoder)
			return test{
				name: "returns (io.WriterCloser, nil)",
				args: args{
					dst: wc,
				},
				fields: fields{
					transcoder: &gob.MockTranscoder{
						NewEncoderFunc: func(w io.Writer) gob.Encoder {
							return enc
						},
					},
				},
				want: want{
					want: &gobWriter{
						dst:     wc,
						encoder: enc,
					},
					err: nil,
				},
			}
		}(),
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
			g := &gobCompressor{
				transcoder: test.fields.transcoder,
			}

			got, err := g.Writer(test.args.dst)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobReader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		src     io.ReadCloser
		decoder gob.Decoder
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
			name: "returns (n, nil) when no error occurs",
			args: args{
				p: []byte{},
			},
			fields: fields{
				decoder: &gob.MockDecoder{
					DecodeFunc: func(e interface{}) error {
						reflect.ValueOf(e).Elem().Set(reflect.ValueOf([]byte("vald")))
						return nil
					},
				},
			},
			want: want{
				wantN: 4,
				err:   nil,
			},
		},

		{
			name: "returns (0, error) when decode fails",
			args: args{
				p: []byte{},
			},
			fields: fields{
				decoder: &gob.MockDecoder{
					DecodeFunc: func(e interface{}) error {
						return errors.New("err")
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
			gr := &gobReader{
				decoder: test.fields.decoder,
			}

			gotN, err := gr.Read(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobReader_Close(t *testing.T) {
	type fields struct {
		src     io.ReadCloser
		decoder *gob.Decoder
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
			name: "returns nil",
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
			gr := &gobReader{
				src: test.fields.src,
			}

			err := gr.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		dst     io.WriteCloser
		encoder gob.Encoder
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
			name: "returns (n, nil) when no error occurs",
			args: args{
				p: []byte{},
			},
			fields: fields{
				encoder: &gob.MockEncoder{
					EncodeFunc: func(e interface{}) error {
						reflect.ValueOf(e).Elem().Set(reflect.ValueOf([]byte("vald")))
						return nil
					},
				},
			},
			want: want{
				wantN: 4,
				err:   nil,
			},
		},

		{
			name: "returns (0, error) when encode fails",
			args: args{
				p: []byte{},
			},
			fields: fields{
				encoder: &gob.MockEncoder{
					EncodeFunc: func(e interface{}) error {
						return errors.New("err")
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
			gw := &gobWriter{
				encoder: test.fields.encoder,
			}

			gotN, err := gw.Write(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobWriter_Close(t *testing.T) {
	type fields struct {
		dst     io.WriteCloser
		encoder gob.Encoder
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
			name: "returns nil",
			fields: fields{
				dst: &MockWriteCloser{
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
			gw := &gobWriter{
				dst: test.fields.dst,
			}

			err := gw.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
