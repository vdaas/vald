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

package compress

import (
	"bytes"
	"reflect"
	"strconv"
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
		want       want
		checkFunc  func(want, Compressor, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
		err  error
		want []byte
	}
	type test struct {
		want       want
		fields     fields
		checkFunc  func(want, []byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
							EncodeFunc: func(e any) error {
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
							EncodeFunc: func(e any) error {
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
		err  error
		want []float32
	}
	type test struct {
		want       want
		checkFunc  func(want, []byte, error, Compressor) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
		err  error
		want []float32
	}
	type test struct {
		want       want
		fields     fields
		checkFunc  func(want, []float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
							DecodeFunc: func(e any) error {
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
							DecodeFunc: func(any) error {
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
		transcoder gob.Transcoder
	}
	type want struct {
		want io.ReadCloser
		err  error
	}
	type test struct {
		want       want
		args       args
		fields     fields
		checkFunc  func(want, io.ReadCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
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
					transcoder: &gob.MockTranscoder{
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
				transcoder: test.fields.transcoder,
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
		want       want
		args       args
		fields     fields
		checkFunc  func(want, io.WriteCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
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
			name: "returns (n, nil) when no error occurs",
			args: args{
				p: make([]byte, 8),
			},
			fields: fields{
				decoder: &gob.MockDecoder{
					DecodeFunc: func(e any) error {
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
					DecodeFunc: func(e any) error {
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

func Test_gobReader_Read_roundTrip(t *testing.T) {
	t.Parallel()

	type test struct {
		wantErr  error
		name     string
		chunks   [][]byte
		truncate int // number of bytes cut from the end of the encoded stream
		readSize int
	}

	seq := func(n int) []byte {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte(i%251 + 1)
		}
		return b
	}
	chunks := [][]byte{seq(1), seq(7), seq(64), {}, seq(300), seq(5)}

	tests := make([]test, 0, 12)
	for _, rs := range []int{1, 3, 7, 8, 64, 100, 300, 301, 4096} {
		tests = append(tests, test{
			name:     "reads all chunks with read buffer size " + strconv.Itoa(rs),
			chunks:   chunks,
			readSize: rs,
		})
	}
	tests = append(tests,
		test{
			name:     "empty chunks only returns EOF",
			chunks:   [][]byte{{}, {}},
			readSize: 16,
		},
		test{
			name:     "no chunk returns EOF",
			readSize: 16,
		},
		test{
			name:     "truncated stream returns io.ErrUnexpectedEOF",
			chunks:   [][]byte{seq(300)},
			truncate: 10,
			readSize: 16,
			wantErr:  io.ErrUnexpectedEOF,
		},
	)

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()

			c, err := NewGob()
			if err != nil {
				tt.Fatal(err)
			}
			encoded := new(bytes.Buffer)
			w, err := c.Writer(&MockWriteCloser{
				WriteFunc: encoded.Write,
				CloseFunc: func() error { return nil },
			})
			if err != nil {
				tt.Fatal(err)
			}
			var want []byte
			for _, ch := range tc.chunks {
				if _, werr := w.Write(ch); werr != nil {
					tt.Fatal(werr)
				}
				want = append(want, ch...)
			}
			stream := encoded.Bytes()
			stream = stream[:len(stream)-tc.truncate]

			newReader := func() io.ReadCloser {
				r, rerr := c.Reader(io.NopCloser(bytes.NewReader(stream)))
				if rerr != nil {
					tt.Fatal(rerr)
				}
				return r
			}

			// manual loop: n must never exceed len(p) and data must land in p.
			r := newReader()
			var got []byte
			p := make([]byte, tc.readSize)
			for {
				n, rerr := r.Read(p)
				if n < 0 || n > len(p) {
					tt.Fatalf("Read returned n=%d for len(p)=%d", n, len(p))
				}
				got = append(got, p[:n]...)
				if rerr != nil {
					if tc.wantErr != nil {
						if !errors.Is(rerr, tc.wantErr) {
							tt.Fatalf("error = %v, want %v", rerr, tc.wantErr)
						}
					} else if !errors.Is(rerr, io.EOF) {
						tt.Fatalf("error = %v, want io.EOF", rerr)
					}
					break
				}
			}
			if tc.wantErr == nil && !bytes.Equal(got, want) {
				tt.Errorf("manual read: got %d bytes, want %d bytes (content differs)", len(got), len(want))
			}

			if tc.wantErr != nil {
				if _, aerr := io.ReadAll(newReader()); !errors.Is(aerr, tc.wantErr) {
					tt.Errorf("io.ReadAll error = %v, want %v", aerr, tc.wantErr)
				}
				return
			}

			all, err := io.ReadAll(newReader())
			if err != nil {
				tt.Fatalf("io.ReadAll: %v", err)
			}
			if !bytes.Equal(all, want) {
				tt.Errorf("io.ReadAll: got %d bytes, want %d bytes (content differs)", len(all), len(want))
			}

			dst := new(bytes.Buffer)
			if _, err := io.Copy(dst, newReader()); err != nil {
				tt.Fatalf("io.Copy: %v", err)
			}
			if !bytes.Equal(dst.Bytes(), want) {
				tt.Errorf("io.Copy: got %d bytes, want %d bytes (content differs)", dst.Len(), len(want))
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
			name: "returns (n, nil) when no error occurs",
			args: args{
				p: []byte{},
			},
			fields: fields{
				encoder: &gob.MockEncoder{
					EncodeFunc: func(e any) error {
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
					EncodeFunc: func(e any) error {
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

// NOT IMPLEMENTED BELOW
