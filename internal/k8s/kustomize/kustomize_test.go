//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package kustomize

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func newConfigMap(name string, data map[string]any) unstructured.Unstructured {
	return unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]any{
			"name":      name,
			"namespace": "default",
		},
		"data": data,
	}}
}

func TestMerge(t *testing.T) {
	t.Run("no overlays returns nil", func(t *testing.T) {
		got, err := Merge()
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("single overlay returns a deep copy", func(t *testing.T) {
		base := newConfigMap("cfg", map[string]any{"key": "value"})
		got, err := Merge(base)
		assert.NoError(t, err)
		assert.Equal(t, base.Object, got.Object)

		got.Object["data"].(map[string]any)["key"] = "mutated"
		assert.Equal(t, "value", base.Object["data"].(map[string]any)["key"],
			"mutating the result must not affect the input")
	})

	t.Run("patch overrides base fields and keeps the rest", func(t *testing.T) {
		base := newConfigMap("cfg", map[string]any{"key": "value", "keep": "untouched"})
		patch := newConfigMap("cfg", map[string]any{"key": "patched"})

		got, err := Merge(base, patch)
		assert.NoError(t, err)

		data, ok := got.Object["data"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "patched", data["key"])
		assert.Equal(t, "untouched", data["keep"])
	})

	t.Run("later patches win", func(t *testing.T) {
		base := newConfigMap("cfg", map[string]any{"key": "value"})
		first := newConfigMap("cfg", map[string]any{"key": "first"})
		second := newConfigMap("cfg", map[string]any{"key": "second"})

		got, err := Merge(base, first, second)
		assert.NoError(t, err)
		assert.Equal(t, "second", got.Object["data"].(map[string]any)["key"])
	})

	t.Run("multiple output resources error instead of keeping the last", func(t *testing.T) {
		// A v1 List base is expanded by kustomize into its items, producing
		// more than one output resource, which Merge must reject rather than
		// silently returning only the last one.
		base := unstructured.Unstructured{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "List",
			"items": []any{
				newConfigMap("cfg-a", map[string]any{"key": "a"}).Object,
				newConfigMap("cfg-b", map[string]any{"key": "b"}).Object,
			},
		}}
		patch := newConfigMap("cfg-a", map[string]any{"key": "patched"})

		got, err := Merge(base, patch)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "want exactly 1")
	})
}

// NOT IMPLEMENTED BELOW
