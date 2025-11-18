//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law of a key agreement protocol, the associated documentation
// and conditional statements specifying the rights and obligations of third parties.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package metrics

import (
	"testing"
	"time"
)

func TestExemplar(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		exemplar func() Exemplar
		offers   []struct {
			latency time.Duration
			id      string
		}
		check func(t *testing.T, e Exemplar)
	}

	tests := []testCase{
		{
			name: "offer requests and check snapshot",
			exemplar: func() Exemplar {
				return NewExemplar(3)
			},
			offers: []struct {
				latency time.Duration
				id      string
			}{
				{100 * time.Millisecond, "req-1"},
				{200 * time.Millisecond, "req-2"},
				{50 * time.Millisecond, "req-3"},
				{300 * time.Millisecond, "req-4"},
			},
			check: func(t *testing.T, e Exemplar) {
				snap := e.Snapshot()
				if len(snap) != 3 {
					t.Fatalf("expected snapshot length 3, got %d", len(snap))
				}

				hasReq2 := false
				hasReq4 := false
				hasReq1 := false

				for _, item := range snap {
					switch item.requestID {
					case "req-1":
						hasReq1 = true
					case "req-2":
						hasReq2 = true
					case "req-4":
						hasReq4 = true
					}
				}

				if !hasReq1 || !hasReq2 || !hasReq4 {
					t.Errorf("expected to find req-1, req-2 and req-4 in snapshot")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.exemplar()
			for _, o := range tt.offers {
				e.Offer(o.latency, o.id)
			}
			tt.check(t, e)
		})
	}
}
