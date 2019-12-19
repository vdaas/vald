//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

import "testing"

func TestLZ4CompressVector(t *testing.T) {
	tests := []struct {
		vector []float64
	}{
		{
			vector: []float64{0.1, 0.2, 0.3},
		},
		{
			vector: []float64{0.4, 0.2, 0.3, 0.1},
		},
		{
			vector: []float64{0.1, 0.5, 0.12, 0.13, 1.0},
		},
	}

	for _, tc := range tests {
		lz4c := NewLZ4()
		compressed, err := lz4c.CompressVector(tc.vector)
		if err != nil {
			t.Fatalf("Compress failed: %s", err)
		}

		decompressed, err := lz4c.DecompressVector(compressed)
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
