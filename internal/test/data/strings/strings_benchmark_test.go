// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package strings

import (
	"fmt"
	"testing"
)

var (
	parallelism = 1000
	short       = 10
	long        = 10000
)

func benchmark(b *testing.B, f func(l int) string) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for i := short; i <= long; i *= 10 {
		cnt := i
		b.Run(fmt.Sprint(cnt), func(bb *testing.B) {
			bb.Helper()
			bb.SetParallelism(parallelism)
			bb.ReportAllocs()
			bb.ResetTimer()
			bb.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = f(cnt)
				}
			})
		})
	}
}

func BenchmarkRandom(b *testing.B) {
	benchmark(b, Random)
}
