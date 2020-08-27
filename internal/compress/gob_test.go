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

// Package compress provides compress functions
package compress

import (
	"io"
	"reflect"
	"testing"
	"unsafe"

	"go.uber.org/goleak"

	"github.com/vdaas/vald/internal/compress/gob"
	"github.com/vdaas/vald/internal/errors"
)

func TestGobCompressVector(t *testing.T) {
	tests := []struct {
		vector []float32
	}{
		{
			vector: []float32{0.1, 0.2, 0.3},
		},
		{
			vector: []float32{0.4, 0.2, 0.3, 0.1},
		},
		{
			vector: []float32{0.1, 0.5, 0.12, 0.13, 1.0},
		},
	}

	for _, tc := range tests {
		gobc, err := NewGob()
		if err != nil {
			t.Fatalf("initialize failed: %s", err)
		}

		compressed, err := gobc.CompressVector(tc.vector)
		if err != nil {
			t.Fatalf("Compress failed: %s", err)
		}

		decompressed, err := gobc.DecompressVector(compressed)
		if err != nil {
			t.Fatalf("Decompress failed: %s", err)
		}
		t.Logf("converted: origin %+v, compressed -> decompressed %+v", tc.vector, decompressed)
		for i := range tc.vector {
			if tc.vector[i] != decompressed[i] {
				t.Fatalf("Invalid convert: origin %+v, compressed -> decompressed %+v", tc.vector, decompressed)
			}
		}
	}
}

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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
					transcoder: func() gob.Transcoder {
						return gob.New()
					}(),
				},
			},
		},

		{
			name: "returns (nil, error) when option apply fails",
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := NewGob(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
								w.Write([]byte("vald"))
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
			name: "returns (nil, error) when Encode fails",
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gobCompressor{}

			got, err := g.CompressVector(test.args.vector)
			if err := test.checkFunc(test.want, got, err); err != nil {
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
		g          *gobCompressor
		want       want
		checkFunc  func(want, []float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
								value := reflect.ValueOf(e)
								vec := []float32{
									1, 2, 3,
								}
								value.SetPointer(unsafe.Pointer(&vec))
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gobCompressor{}

			got, err := g.DecompressVector(test.args.bs)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gobCompressor_Reader(t *testing.T) {
	type args struct {
		src io.ReadCloser
	}
	type want struct {
		want io.ReadCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		g          *gobCompressor
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
		           src: nil,
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
		           src: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gobCompressor{}

			got, err := g.Reader(test.args.src)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_gobCompressor_Writer(t *testing.T) {
	type args struct {
		dst io.WriteCloser
	}
	type want struct {
		want io.WriteCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		g          *gobCompressor
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
		           dst: nil,
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
		           dst: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &gobCompressor{}

			got, err := g.Writer(test.args.dst)
			if err := test.checkFunc(test.want, got, err); err != nil {
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
		decoder *gob.Decoder
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got = %v, want %v", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           p: nil,
		       },
		       fields: fields {
		           src: nil,
		           decoder: nil,
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
		           p: nil,
		           },
		           fields: fields {
		           src: nil,
		           decoder: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			gr := &gobReader{
				src: test.fields.src,
				// decoder: test.fields.decoder,
			}

			gotN, err := gr.Read(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
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
		           src: nil,
		           decoder: nil,
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
		           src: nil,
		           decoder: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			gr := &gobReader{
				src: test.fields.src,
				// decoder: test.fields.decoder,
			}

			err := gr.Close()
			if err := test.checkFunc(test.want, err); err != nil {
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
		encoder *gob.Encoder
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got = %v, want %v", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           p: nil,
		       },
		       fields: fields {
		           dst: nil,
		           encoder: nil,
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
		           p: nil,
		           },
		           fields: fields {
		           dst: nil,
		           encoder: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			gw := &gobWriter{
				dst: test.fields.dst,
				// encoder: test.fields.encoder,
			}

			gotN, err := gw.Write(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_gobWriter_Close(t *testing.T) {
	type fields struct {
		dst     io.WriteCloser
		encoder *gob.Encoder
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
		           dst: nil,
		           encoder: nil,
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
		           dst: nil,
		           encoder: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			gw := &gobWriter{
				dst: test.fields.dst,
				// encoder: test.fields.encoder,
			}

			err := gw.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
