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
	"flag"
	"strings"
	"testing"

	tstr "github.com/vdaas/vald/internal/test/data/strings"
)

var (
	testDataLen   = 2 << 5
	testDataCount = 2 << 8
	testData      = make([]string, 0, testDataCount)
)

func TestMain(m *testing.M) {
	testing.Init()
	flag.Parse()
	if testing.Short() {
		m.Run()
		return
	}
	for i := 0; i < testDataCount; i++ {
		testData = append(testData, tstr.Random(testDataLen))
	}
	m.Run()
	testData = nil
}

func BenchmarkStandardJoin(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := strings.Join(testData, ":")
			_ = s
		}
	})
}

func BenchmarkValdJoin(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := Join(testData, ":")
			_ = s
		}
	})
}
