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

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/strings"
)

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

	current := data
	for _, part := range parts {
		switch typed := current.(type) {
		case map[string]any:
			if part == "length()" {
				return len(typed), nil
			}
			val, exists := typed[part]
			if !exists {
				return nil, fmt.Errorf("key '%s' not found", part)
			}
			current = val

		case []any:
			if part == "length()" {
				return len(typed), nil
			}
			return nil, fmt.Errorf("unexpected array when accessing '%s'", part)

		default:
			return nil, fmt.Errorf("cannot access '%s' on non-object", part)
		}
	}

	return current, nil
}
