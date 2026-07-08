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
package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestIndex_SetResources(t *testing.T) {
	idx := &Index{}
	idx.SetResources()

	assert.NotNil(t, idx.Resources)
	assert.NotEmpty(t, idx.Resources.Requests[v1.ResourceCPU])
	assert.NotEmpty(t, idx.Resources.Requests[v1.ResourceMemory])
	assert.NotEmpty(t, idx.Resources.Limits[v1.ResourceCPU])
	assert.NotEmpty(t, idx.Resources.Limits[v1.ResourceMemory])
}

func TestIndex_SetTopologySpreadConstraints(t *testing.T) {
	idx := &Index{}
	idx.SetTopologySpreadConstraints()

	assert.Len(t, idx.TopologySpreadConstraints, 1)
	assert.Equal(t, "manager-index", idx.TopologySpreadConstraints[0].LabelSelector.MatchLabels["app.kubernetes.io/component"])
}

// NOT IMPLEMENTED BELOW
