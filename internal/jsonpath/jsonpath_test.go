//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// package jsonpath provides utilities for working with JSONPath in Go.
package jsonpath

import (
	"reflect"
	"testing"
)

func TestJSONPathEval(t *testing.T) {
	jsonStr := []byte(`{
		"test": [1, 2, 3],
		"obj": {"key": "value"},
		"counts": {
			"A": { "stored": 10 },
			"B": { "stored": 5, "uncommited": 5 }
		},
		"array": [
			{"A": "value1"},
			{"B": "value2"}
		]
	}`)

	tests := []struct {
		json     []byte
		name     string
		path     string
		expected any
		wantErr  bool
	}{
		{
			json:     []byte("{}"),
			name:     "empty length",
			path:     "$.length()",
			expected: 0,
			wantErr:  false,
		},
		{
			json:     []byte("{}"),
			name:     "root",
			path:     "$",
			expected: map[string]any{},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "access map length",
			path:     "$.length()",
			expected: 4,
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "access array",
			path:     ".test",
			expected: []any{float64(1), float64(2), float64(3)},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "access array length",
			path:     ".test.length()",
			expected: 3,
			wantErr:  false,
		},
		{
			json:    jsonStr,
			name:    "missing key",
			path:    ".missing",
			wantErr: true,
		},
		{
			json:     jsonStr,
			name:     "length on non-array",
			path:     "$.obj.length()",
			expected: 1,
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "nested key access",
			path:     ".obj.key",
			expected: "value",
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "* for map",
			path:     ".obj.*",
			expected: []any{"value"},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "*.key",
			path:     ".counts.*.stored",
			expected: []any{10.0, 5.0},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "*.*",
			path:     ".counts.*.*",
			expected: []any{10.0, 5.0, 5.0},

			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "*.* for array",
			path:     ".array.*.*",
			expected: []any{"value1", "value2"},

			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONPathEval(tt.json, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr = %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("error details: %v", err)
				}
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got = %#v, expected = %#v", got, tt.expected)
			}
		})
	}
}
