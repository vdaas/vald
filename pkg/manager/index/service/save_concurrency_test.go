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
package service

import "testing"

// saveIndexConcurrency is passed to OrderedRangeConcurrent, where 0 means an
// unbounded fan-out of SaveIndex to every agent; it must default to a bound.
func TestNew_saveIndexConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts []Option
		want int
	}{
		{name: "default is serial", want: 1},
		{name: "option is applied", opts: []Option{WithSaveConcurrency(4)}, want: 4},
		{name: "zero keeps default", opts: []Option{WithSaveConcurrency(0)}, want: 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			got, err := New(tc.opts...)
			if err != nil {
				tt.Fatalf("New() error = %v", err)
			}
			idx, ok := got.(*index)
			if !ok {
				tt.Fatalf("New() returned %T, want *index", got)
			}
			if idx.saveIndexConcurrency != tc.want {
				tt.Errorf("saveIndexConcurrency = %d, want %d", idx.saveIndexConcurrency, tc.want)
			}
		})
	}
}
