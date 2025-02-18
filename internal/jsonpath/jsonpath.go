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
	"fmt"
	"sort"
	"strconv"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/strings"
)

func flatten(input []any) []any {
	var out []any
	for _, item := range input {
		if inner, ok := item.([]any); ok {
			out = append(out, flatten(inner)...)
		} else {
			out = append(out, item)
		}
	}
	return out
}

func recEval(data any, parts []string) (any, error) {
	if len(parts) == 0 {
		return data, nil
	}

	part := parts[0]
	rest := parts[1:]

	switch typed := data.(type) {
	case map[string]any:
		if part == "length()" {
			return len(typed), nil
		}
		if part == "sum()" {
			var sum float64
			for _, v := range typed {
				if num, ok := v.(float64); ok {
					sum += num
				}
			}
			return sum, nil
		}
		if part == "*" {
			// Return all keys if part is '*'
			result := make([]any, len(typed))
			// Sort keys to ensure consistent order
			keys := make([]string, 0, len(typed))
			for k := range typed {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				v, err := recEval(typed[k], rest)
				result[i] = v
				if err != nil {
					return nil, err
				}
			}
			return flatten(result), nil
		}
		val, exists := typed[part]
		if !exists {
			return nil, fmt.Errorf("key '%s' not found in %v", part, typed)
		}
		return recEval(val, rest)

	case []any:
		if part == "length()" {
			return len(typed), nil
		}
		if part == "sum()" {
			var sum float64
			for _, v := range typed {
				if num, ok := v.(float64); ok {
					sum += num
				}
			}
			return sum, nil
		}
		if part == "*" {
			// Map over all elements in the array
			result := make([]any, len(typed))
			for i, v := range typed {
				res, err := recEval(v, rest)
				if err != nil {
					return nil, err
				}
				result[i] = res
			}
			return flatten(result), nil
		}
		idx, err := strconv.Atoi(part)
		if err == nil {
			if idx < 0 {
				idx += len(typed) // Handle negative indices
			}
			if idx < 0 || idx >= len(typed) {
				return nil, fmt.Errorf("index %d out of range for array of length %d", idx, len(typed))
			}
			val := typed[idx]
			return recEval(val, rest)
		}

		return nil, fmt.Errorf("unexpected array when accessing '%s'", part)

	default:
		return nil, fmt.Errorf("cannot access '%s' on non-object", part)
	}
}

// JSONPathEval supports .field and .length() syntax for JSONPath evaluation.
func JSONPathEval(jsonData []byte, path string) (any, error) {
	var data any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	if !strings.Contains(path, ".") {
		return data, nil
	}

	parts := strings.Split(path, ".")[1:]

	current, err := recEval(data, parts)
	if err != nil {
		return nil, err
	}

	return current, nil
}
