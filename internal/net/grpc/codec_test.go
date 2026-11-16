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

package grpc

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	_ "github.com/vdaas/vald/internal/net/grpc/logger"
	"github.com/vdaas/vald/internal/sync"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/mem"
)

func TestCodec_Marshal(t *testing.T) {
	t.Parallel()
	type args struct {
		v any
	}
	type want struct {
		err  error
		want []byte
	}
	type test struct {
		c          Codec
		args       args
		checkFunc  func(want, mem.BufferSlice, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got mem.BufferSlice, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\t\t\twant: \"%#v\"", err, w.err)
		}
		gotBytes := got.Materialize()
		if len(gotBytes) == 0 && len(w.want) == 0 {
			return nil
		}
		if !reflect.DeepEqual(gotBytes, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\t\twant: \"%#v\"", gotBytes, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return marshal result when val is vtproto message",
			args: args{
				v: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				},
			},
			checkFunc: func(_ want, b mem.BufferSlice, e error) error {
				if e != nil {
					return e
				}
				if b.Len() == 0 {
					return errors.New("return byte is empty")
				}
				return nil
			},
		},
		{
			name: "return marshal result when val is empty proto message",
			args: args{
				v: &payload.Empty{},
			},
			want: want{
				want: []byte{},
				err:  nil,
			},
		},
		{
			name: "return error when val is not proto message",
			args: args{
				v: []int{1},
			},
			want: want{
				want: nil,
				err:  errors.New("proto: failed to marshal, message is []int, want proto.Message"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			c := Codec{
				fallback: encoding.GetCodecV2(Name),
			}

			got, err := c.Marshal(test.args.v)
			if err := checkFunc(test.want, got, err); err != nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestCodec_Unmarshal(t *testing.T) {
	t.Parallel()
	type args struct {
		v    any
		data []byte
	}
	type want struct {
		err error
	}
	type test struct {
		c          Codec
		want       want
		checkFunc  func(test, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(t test, err error) error {
		if !errors.Is(err, t.want.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\t\twant: \"%#v\"", err, t.want.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "unmarshal success to parse data into v",
			args: args{
				data: func() []byte {
					b, _ := Codec{}.Marshal(&payload.Object_Vector{
						Id:     "1",
						Vector: []float32{1.0, 2.1},
					})
					return b.Materialize()
				}(),
				v: &payload.Object_Vector{},
			},
			checkFunc: func(t test, _ error) error {
				if !reflect.DeepEqual(t.args.v, &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				}) {
					return errors.New("unmarshal result is not correct")
				}
				return nil
			},
		},
		{
			name: "return error when data is invalid",
			args: args{
				data: []byte{0, 1, 2},
				v:    &payload.Object_Vector{},
			},
			want: want{
				err: errors.New("proto: Object_Vector: illegal tag 0 (wire type 0)"),
			},
		},
		{
			name: "return error when v is invalid",
			args: args{
				data: []byte{0, 1, 2},
				v:    Codec{},
			},
			want: want{
				err: errors.New("failed to unmarshal, message is grpc.Codec, want proto.Message"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			c := Codec{
				fallback: encoding.GetCodecV2(Name),
			}

			err := c.Unmarshal(mem.BufferSlice{mem.SliceBuffer(test.args.data)}, test.args.v)
			if err := checkFunc(test, err); err != nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestCodec_Name(t *testing.T) {
	t.Parallel()
	type want struct {
		want string
	}
	type test struct {
		c          Codec
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return codec name",
			want: want{
				want: "proto",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			c := Codec{}

			got := c.Name()
			if err := checkFunc(test.want, got); err != nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

// testSizedVTMessage implements vtprotoMessage with exact arbitrary byte sizes.
type testSizedVTMessage struct {
	data []byte
}

func (m *testSizedVTMessage) SizeVT() int {
	if m == nil {
		return 0
	}
	return len(m.data)
}

func (m *testSizedVTMessage) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	if m == nil || len(m.data) == 0 {
		return 0, nil
	}
	copy(dAtA[:len(m.data)], m.data)
	return len(m.data), nil
}

func (m *testSizedVTMessage) UnmarshalVT(dAtA []byte) error {
	m.data = make([]byte, len(dAtA))
	copy(m.data, dAtA)
	return nil
}

func makeSizedVTMessage(size int) *testSizedVTMessage {
	return &testSizedVTMessage{
		data: bytes.Repeat([]byte{0x5a}, size),
	}
}

// TestCodec_Marshal_boundaryPayloads verifies Codec behavior across all required
// boundary payload sizes: 0B, 1B, 1023B, 1024B, 1025B, 2047B, 2048B, 64KB, and multi-MB.
func TestCodec_Marshal_boundaryPayloads(t *testing.T) {
	pool := mem.NewTieredBufferPool(
		2<<10,   // 2KB
		4<<10,   // 4KB
		8<<10,   // 8KB
		16<<10,  // 16KB
		32<<10,  // 32KB
		64<<10,  // 64KB
		128<<10, // 128KB
		256<<10, // 256KB
		512<<10, // 512KB
		1<<20,   // 1MB
		2<<20,   // 2MB
		4<<20,   // 4MB
	)

	c := Codec{
		buffer:   pool,
		fallback: encoding.GetCodecV2(Name),
	}

	tests := []struct {
		name       string
		targetSize int
		msg        vtprotoMessage
	}{
		{
			name:       "0-byte empty payload (payload.Empty)",
			targetSize: 0,
			msg:        &payload.Empty{},
		},
		{
			name:       "0-byte empty VTMessage",
			targetSize: 0,
			msg:        makeSizedVTMessage(0),
		},
		{
			name:       "1-byte payload",
			targetSize: 1,
			msg:        makeSizedVTMessage(1),
		},
		{
			name:       "1023-byte payload (1B below 1KB boundary)",
			targetSize: 1023,
			msg:        makeSizedVTMessage(1023),
		},
		{
			name:       "1024-byte payload (exact 1KB threshold)",
			targetSize: 1024,
			msg:        makeSizedVTMessage(1024),
		},
		{
			name:       "1025-byte payload (1B above 1KB threshold)",
			targetSize: 1025,
			msg:        makeSizedVTMessage(1025),
		},
		{
			name:       "2047-byte payload (1B below 2KB tier)",
			targetSize: 2047,
			msg:        makeSizedVTMessage(2047),
		},
		{
			name:       "2048-byte payload (exact 2KB tier boundary)",
			targetSize: 2048,
			msg:        makeSizedVTMessage(2048),
		},
		{
			name:       "65536-byte payload (64KB boundary)",
			targetSize: 65536,
			msg:        makeSizedVTMessage(65536),
		},
		{
			name:       "1MB payload (1048576 bytes)",
			targetSize: 1024 * 1024,
			msg:        makeSizedVTMessage(1024 * 1024),
		},
		{
			name:       "2MB payload (2097152 bytes)",
			targetSize: 2 * 1024 * 1024,
			msg:        makeSizedVTMessage(2 * 1024 * 1024),
		},
		{
			name:       "4MB payload (4194304 bytes)",
			targetSize: 4 * 1024 * 1024,
			msg:        makeSizedVTMessage(4 * 1024 * 1024),
		},
		{
			name:       "Real protobuf Object_Blob with 1536d vector (6144B payload)",
			targetSize: 6154, // 6144B object + tag/varint header + Id
			msg: &payload.Object_Blob{
				Id:     "vec-1",
				Object: bytes.Repeat([]byte{0xaa}, 6144),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Check size computation
			gotSize := tt.msg.SizeVT()
			if tt.targetSize > 0 && gotSize != tt.targetSize {
				t.Fatalf("SizeVT() = %d, want %d", gotSize, tt.targetSize)
			}

			// 2. Marshal
			bufSlice, err := c.Marshal(tt.msg)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			defer bufSlice.Free()

			if bufSlice.Len() != gotSize {
				t.Fatalf("bufSlice.Len() = %d, want %d", bufSlice.Len(), gotSize)
			}

			// 3. Unmarshal back into target
			switch orig := tt.msg.(type) {
			case *payload.Empty:
				var unmarshaled payload.Empty
				if err := c.Unmarshal(bufSlice, &unmarshaled); err != nil {
					t.Fatalf("Unmarshal failed: %v", err)
				}
			case *testSizedVTMessage:
				var unmarshaled testSizedVTMessage
				if err := c.Unmarshal(bufSlice, &unmarshaled); err != nil {
					t.Fatalf("Unmarshal failed: %v", err)
				}
				if unmarshaled.SizeVT() != gotSize {
					t.Fatalf("unmarshaled SizeVT() = %d, want %d", unmarshaled.SizeVT(), gotSize)
				}
				if !bytes.Equal(unmarshaled.data, orig.data) {
					t.Fatal("unmarshaled bytes mismatch with original")
				}
			case *payload.Object_Blob:
				var unmarshaled payload.Object_Blob
				if err := c.Unmarshal(bufSlice, &unmarshaled); err != nil {
					t.Fatalf("Unmarshal failed: %v", err)
				}
				if unmarshaled.Id != orig.Id {
					t.Fatalf("unmarshaled Id = %q, want %q", unmarshaled.Id, orig.Id)
				}
				if !bytes.Equal(unmarshaled.Object, orig.Object) {
					t.Fatal("unmarshaled Object mismatch with original")
				}
			}
		})
	}
}

// TestCodec_nilBuffer verifies that Codec functions safely without
// panics when c.buffer is nil across all boundary payload sizes.
func TestCodec_nilBuffer(t *testing.T) {
	cNil := Codec{
		buffer:   nil,
		fallback: encoding.GetCodecV2(Name),
	}

	sizes := []int{0, 1, 1023, 1024, 1025, 2047, 2048, 65536, 1024 * 1024}

	for _, sz := range sizes {
		t.Run(fmt.Sprintf("size_%d", sz), func(t *testing.T) {
			msg := makeSizedVTMessage(sz)

			// Marshal without pool
			bufSlice, err := cNil.Marshal(msg)
			if err != nil {
				t.Fatalf("Marshal with nil buffer failed for size %d: %v", sz, err)
			}
			defer bufSlice.Free()

			if bufSlice.Len() != sz {
				t.Fatalf("bufSlice.Len() = %d, want %d", bufSlice.Len(), sz)
			}

			// Unmarshal without pool
			var unmarshaled testSizedVTMessage
			if err := cNil.Unmarshal(bufSlice, &unmarshaled); err != nil {
				t.Fatalf("Unmarshal with nil buffer failed for size %d: %v", sz, err)
			}

			if unmarshaled.SizeVT() != sz {
				t.Fatalf("unmarshaled SizeVT() = %d, want %d", unmarshaled.SizeVT(), sz)
			}
		})
	}
}

// TestCodec_bufferPoolMemoryStability stress-tests buffer recycling
// across 500,000 iterations for edge-case sizes (1B, 1023B, 1024B, 1025B) to empirically
// verify that mem.NewBuffer does NOT leak buffers and RSS/HeapAlloc remains stable.
func TestCodec_bufferPoolMemoryStability(t *testing.T) {
	pool := mem.NewTieredBufferPool(
		2<<10,  // 2KB
		4<<10,  // 4KB
		8<<10,  // 8KB
		16<<10, // 16KB
		64<<10, // 64KB
	)

	c := Codec{
		buffer:   pool,
		fallback: encoding.GetCodecV2(Name),
	}

	// Boundary sizes around 1KB threshold
	boundarySizes := []int{1, 1023, 1024, 1025, 2048}

	for _, sz := range boundarySizes {
		t.Run(fmt.Sprintf("boundary_size_%d_stress", sz), func(t *testing.T) {
			msg := makeSizedVTMessage(sz)

			// Warm up pool
			for range 100 {
				b, err := c.Marshal(msg)
				if err != nil {
					t.Fatal(err)
				}
				b.Free()
			}

			runtime.GC()
			var memBefore runtime.MemStats
			runtime.ReadMemStats(&memBefore)

			iterations := 500_000
			if testing.Short() {
				iterations = 1_000
			}
			for range iterations {
				bufSlice, err := c.Marshal(msg)
				if err != nil {
					t.Fatalf("Marshal failed at iteration: %v", err)
				}

				var target testSizedVTMessage
				if err := c.Unmarshal(bufSlice, &target); err != nil {
					bufSlice.Free()
					t.Fatalf("Unmarshal failed: %v", err)
				}

				// Crucial: Free the buffer slice back to the pool
				bufSlice.Free()
			}

			runtime.GC()
			var memAfter runtime.MemStats
			runtime.ReadMemStats(&memAfter)

			// If 500,000 allocations of 2KB leaked without being returned to the pool,
			// HeapAlloc would grow by at least 1 GB (500,000 * 2048 bytes = 1,024,000,000 bytes).
			// With proper buffer pooling, memory in use should remain bounded under 32MB.
			heapDelta := int64(memAfter.HeapAlloc) - int64(memBefore.HeapAlloc)
			t.Logf("Size %d: HeapAlloc before=%d KB, after=%d KB, delta=%d KB across %d iterations",
				sz, memBefore.HeapAlloc/1024, memAfter.HeapAlloc/1024, heapDelta/1024, iterations)

			const maxAllowedHeapGrowth = 32 * 1024 * 1024 // 32MB
			if heapDelta > maxAllowedHeapGrowth {
				t.Fatalf("Memory leak detected! HeapAlloc grew by %d MB (limit: %d MB) for payload size %d",
					heapDelta/(1024*1024), maxAllowedHeapGrowth/(1024*1024), sz)
			}
		})
	}
}

// TestCodec_concurrent verifies thread-safety and absence of data
// races during concurrent buffer pool usage across 16 parallel goroutines.
func TestCodec_concurrent(t *testing.T) {
	pool := mem.NewTieredBufferPool(
		2<<10,  // 2KB
		4<<10,  // 4KB
		8<<10,  // 8KB
		64<<10, // 64KB
	)

	c := Codec{
		buffer:   pool,
		fallback: encoding.GetCodecV2(Name),
	}

	const goroutines = 16
	const iterationsPerGoroutine = 10_000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := range goroutines {
		go func(workerID int) {
			defer wg.Done()

			sizes := []int{0, 1, 1023, 1024, 1025, 2048, 4096}
			for i := range iterationsPerGoroutine {
				sz := sizes[(workerID+i)%len(sizes)]
				msg := makeSizedVTMessage(sz)

				bufSlice, err := c.Marshal(msg)
				if err != nil {
					t.Errorf("worker %d: Marshal failed: %v", workerID, err)
					return
				}

				if bufSlice.Len() != sz {
					t.Errorf("worker %d: got len %d, want %d", workerID, bufSlice.Len(), sz)
					bufSlice.Free()
					return
				}

				var unmarshaled testSizedVTMessage
				if err := c.Unmarshal(bufSlice, &unmarshaled); err != nil {
					t.Errorf("worker %d: Unmarshal failed: %v", workerID, err)
					bufSlice.Free()
					return
				}

				bufSlice.Free()
			}
		}(g)
	}

	wg.Wait()
}

// NOT IMPLEMENTED BELOW
//
// func TestInitCodec(t *testing.T) {
// 	type args struct {
// 		size int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
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
// 		           size:0,
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
// 			InitCodec(test.args.size)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
