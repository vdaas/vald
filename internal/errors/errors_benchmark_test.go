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
package errors

import (
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/bench"
	"github.com/vdaas/vald/internal/test/data/strings"
)

const fastParallelism = 4

var (
	parallelism = 10000

	bigData      = map[string]error{}
	bigDataLen   = 2 << 8
	bigDataCount = 2 << 8
	//nolint:gochecknoglobals // lazy benchmark dataset initialization
	bigDataOnce sync.Once

	smallData = map[string]error{
		"string": New("aaaa"),
		"int":    New("123"),
		"float":  New("99.99"),
		"struct": New("struct{}{}"),
	}
)

func ensureBigData() {
	bigDataOnce.Do(func() {
		for range bigDataCount {
			bigData[strings.Random(bigDataLen)] = New(strings.Random(bigDataLen))
		}
	})
}

func benchmark(b *testing.B, data map[string]error, join func(err1, err2 error) error) {
	b.Helper()
	p := fastParallelism
	if bench.Full() {
		p = parallelism
	}
	b.SetParallelism(p)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var errs error
		for pb.Next() {
			for _, err := range data {
				errs = join(errs, err)
			}
		}
	})
}

func BenchmarkWrapLongData(b *testing.B) {
	bench.SkipUnlessFull(b)
	ensureBigData()
	benchmark(b, bigData,
		func(err1, err2 error) error { return Wrap(err1, err2.Error()) })
}

func BenchmarkWrapShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return Wrap(err1, err2.Error()) })
}

func BenchmarkJoinLongData(b *testing.B) {
	bench.SkipUnlessFull(b)
	ensureBigData()
	benchmark(b, bigData,
		func(err1, err2 error) error { return Join(err1, err2) })
}

func BenchmarkJoinShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return Join(err1, err2) })
}

func BenchmarkStdWrapLongData(b *testing.B) {
	bench.SkipUnlessFull(b)
	ensureBigData()
	benchmark(b, bigData,
		func(err1, err2 error) error { return fmt.Errorf("%s: %w", err2.Error(), err1) })
}

func BenchmarkStdWrapShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return fmt.Errorf("%s: %w", err2.Error(), err1) })
}

func BenchmarkStdJoinLongData(b *testing.B) {
	bench.SkipUnlessFull(b)
	ensureBigData()
	benchmark(b, bigData,
		func(err1, err2 error) error { return Join(err1, err2) })
}

func BenchmarkStdJoinShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return Join(err1, err2) })
}
