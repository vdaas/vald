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
package internal

import "testing"

func TestRecall(t *testing.T) {
	x := []string{
		"foo", "bar", "hoge", "huga",
	}
	y := []string{
		"foo", "baz", "hoge", "huga", "moge",
	}

	tests := []struct {
		k      int
		recall float64
	}{
		{1, 1.0},
		{2, 0.5},
		{3, 2.0 / 3.0},
		{4, 3.0 / 4.0},
	}

	for _, tt := range tests {
		if Recall(x, y, tt.k) != tt.recall {
			t.Errorf("Recall@%d is %f", tt.k, tt.recall)
		}
	}
}
