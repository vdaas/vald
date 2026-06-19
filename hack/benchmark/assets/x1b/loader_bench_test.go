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
package x1b

import "testing"

const (
	bvecsFile = "../large/sift1b/bigann_base.bvecs"
	fvecsFile = "../large/sift1b/gnd/dis_1000M.fvecs"
	ivecsFile = "../large/sift1b/gnd/idx_1000M.ivecs"
)

func BenchmarkBVecs(b *testing.B) {
	bv, err := NewUint8Vectors(bvecsFile)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := bv.Close(); err != nil {
			b.Fatal(err)
		}
	}()

	b.Run(bvecsFile, func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()

		size := bv.Size()
		i := 0
		for bb.Loop() {
			_, err := bv.Load(i % size)
			if err != nil {
				bb.Fatal(err)
			}
			i++
		}
	})
}

func BenchmarkFVecs(b *testing.B) {
	fv, err := NewFloatVectors(fvecsFile)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := fv.Close(); err != nil {
			b.Fatal(err)
		}
	}()

	b.Run(fvecsFile, func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()

		size := fv.Size()
		i := 0
		for bb.Loop() {
			_, err := fv.Load(i % size)
			if err != nil {
				bb.Fatal(err)
			}
			i++
		}
	})
}

func BenchmarkIVecs(b *testing.B) {
	iv, err := NewInt32Vectors(ivecsFile)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := iv.Close(); err != nil {
			b.Fatal(err)
		}
	}()

	b.Run(ivecsFile, func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()

		size := iv.Size()
		i := 0
		for bb.Loop() {
			_, err := iv.Load(i % size)
			if err != nil {
				bb.Fatal(err)
			}
			i++
		}
	})
}
