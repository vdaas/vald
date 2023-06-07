// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/test/data/strings"
)

var (
	parallelism = 10000

	bigData      = map[string]error{}
	bigDataLen   = 2 << 8
	bigDataCount = 2 << 8

	smallData = map[string]error{
		"string": New("aaaa"),
		"int":    New("123"),
		"float":  New("99.99"),
		"struct": New("struct{}{}"),
	}
)

func TestMain(m *testing.M) {
	testing.Init()
	flag.Parse()
	if testing.Short() {
		m.Run()
		return
	}
	for i := 0; i < bigDataCount; i++ {
		bigData[strings.Random(bigDataLen)] = New(strings.Random(bigDataLen))
	}
	m.Run()
	bigData = nil
}

func benchmark(b *testing.B, data map[string]error,
	join func(err1, err2 error) error,
) {
	b.Helper()
	b.SetParallelism(parallelism)
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
	benchmark(b, bigData,
		func(err1, err2 error) error { return Wrap(err1, err2.Error()) })
}

func BenchmarkWrapShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return Wrap(err1, err2.Error()) })
}

func BenchmarkJoinLongData(b *testing.B) {
	benchmark(b, bigData,
		func(err1, err2 error) error { return Join(err1, err2) })
}

func BenchmarkJoinShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return Join(err1, err2) })
}

func BenchmarkStdWrapLongData(b *testing.B) {
	benchmark(b, bigData,
		func(err1, err2 error) error { return fmt.Errorf("%s: %w", err2.Error(), err1) })
}

func BenchmarkStdWrapShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return fmt.Errorf("%s: %w", err2.Error(), err1) })
}

func BenchmarkStdJoinLongData(b *testing.B) {
	benchmark(b, bigData,
		func(err1, err2 error) error { return errors.Join(err1, err2) })
}

func BenchmarkStdJoinShortData(b *testing.B) {
	benchmark(b, smallData,
		func(err1, err2 error) error { return errors.Join(err1, err2) })
}
