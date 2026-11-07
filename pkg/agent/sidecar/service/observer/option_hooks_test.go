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
package observer

import "testing"

// TestWithHooks_accumulates pins that repeated WithHooks calls append instead
// of replacing the hooks registered by an earlier call.
func TestWithHooks_accumulates(t *testing.T) {
	t.Parallel()

	h1, h2, h3 := new(recordHook), new(recordHook), new(recordHook)
	tests := []struct {
		name  string
		calls [][]Hook
		want  []Hook
	}{
		{name: "single call", calls: [][]Hook{{h1}}, want: []Hook{h1}},
		{name: "second call keeps the first hooks", calls: [][]Hook{{h1}, {h2, h3}}, want: []Hook{h1, h2, h3}},
		{name: "nil call is a no-op", calls: [][]Hook{{h1}, nil}, want: []Hook{h1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o := new(observer)
			for _, hooks := range tt.calls {
				if err := WithHooks(hooks...)(o); err != nil {
					t.Fatalf("WithHooks returned error: %v", err)
				}
			}
			if len(o.hooks) != len(tt.want) {
				t.Fatalf("hooks = %v, want %v", o.hooks, tt.want)
			}
			for i := range tt.want {
				if o.hooks[i] != tt.want[i] {
					t.Errorf("hooks[%d] = %p, want %p", i, o.hooks[i], tt.want[i])
				}
			}
		})
	}
}
