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

package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestToUnstructured(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "cfg", Namespace: "ns"},
		Data:       map[string]string{"key": "value"},
	}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	us, err := ToUnstructured(cm)
	assert.NoError(t, err)
	assert.Equal(t, "ConfigMap", us.GetKind())
	assert.Equal(t, "cfg", us.GetName())
	assert.Equal(t, "ns", us.GetNamespace())

	data, ok := us.Object["data"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", data["key"])
}

func TestObjectsOf(t *testing.T) {
	items := []corev1.ConfigMap{
		{ObjectMeta: metav1.ObjectMeta{Name: "a", Namespace: "ns"}},
		{ObjectMeta: metav1.ObjectMeta{Name: "b", Namespace: "ns"}},
	}

	out := ObjectsOf(items)
	assert.Len(t, out, 2)
	assert.Equal(t, "a", out[0].GetName())
	assert.Equal(t, "b", out[1].GetName())

	// The returned objects must alias the input slice elements, not copies.
	out[0].SetName("mutated")
	assert.Equal(t, "mutated", items[0].GetName())
}

// NOT IMPLEMENTED BELOW
