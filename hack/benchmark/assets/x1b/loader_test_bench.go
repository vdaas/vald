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
package x1b

import (
	"testing"
)

const (
	bvecsFile = "../large/sift1b/bigann_base.bvecs"
	fvecsFile = "../large/sift1b/gnd/dis_1000M.fvecs"
	ivecsFile = "../large/sift1b/gnd/idx_1000M.ivecs"
)

func BenchmarkBVecs(b *testing.B) {
	bv, err := NewBVecs(bvecsFile)
	defer func() {
		if err := bv.Close(); err != nil {
			b.Fatal(err)
		}
	}()
	if err != nil {
		b.Fatal(err)
	}

	i := 0
	b.Run("", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < b.N; n++ {
			v, err := bv.Load(i)
			switch err {
			case nil:
				i++
				continue
			case ErrOutOfBounds:
				if err := bv.Close(); err != nil {
					bb.Fatal(err)
				}
				bv, err = NewBVecs(bvecsFile)
				i = 0
			}
			if err != nil {
				bb.Fatal(err)
			}
			bb.Log(v)
		}
	})
}

func BenchmarkFVecs(b *testing.B) {
	fv, err := NewFVecs(fvecsFile)
	defer func() {
		if err := fv.Close(); err != nil {
			b.Fatal(err)
		}
	}()
	if err != nil {
		b.Fatal(err)
	}

	i := 0
	b.Run("", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < b.N; n++ {
			v, err := fv.Load(i)
			switch err {
			case nil:
				i++
				continue
			case ErrOutOfBounds:
				if err := fv.Close(); err != nil {
					bb.Fatal(err)
				}
				fv, err = NewFVecs(fvecsFile)
				i = 0
			}
			if err != nil {
				bb.Fatal(err)
			}
			bb.Log(v)
		}
	})
}

func BenchmarkIVecs(b *testing.B) {
	iv, err := NewIVecs(ivecsFile)
	defer func() {
		if err := iv.Close(); err != nil {
			b.Fatal(err)
		}
	}()

	if err != nil {
		b.Fatal(err)
	}

	i := 0
	b.Run("", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < b.N; n++ {
			v, err := iv.Load(i)
			switch err {
			case nil:
				i++
				continue
			case ErrOutOfBounds:
				if err := iv.Close(); err != nil {
					bb.Fatal(err)
				}
				iv, err = NewIVecs(ivecsFile)
				i = 0
			}
			if err != nil {
				bb.Fatal(err)
			}
			bb.Log(v)
		}
	})
}