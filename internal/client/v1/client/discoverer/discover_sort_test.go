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

package discoverer

import (
	"slices"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// A nil Memory must sort as usage 0. Treating it as "equal to everything" is
// not a strict weak ordering, and [5, nil, 1] stayed unsorted with it.
func Test_compareMemory_nilMemory(t *testing.T) {
	t.Parallel()
	mem := func(u float64) *payload.Info_Memory { return &payload.Info_Memory{Usage: u} }
	const (
		nilName = "nil"
		low     = "low"
		high    = "high"
	)
	want := []string{nilName, low, high}

	t.Run("nodes", func(t *testing.T) {
		t.Parallel()
		nodes := []*payload.Info_Node{
			{Name: high, Memory: mem(5)},
			{Name: nilName},
			{Name: low, Memory: mem(1)},
		}
		slices.SortFunc(nodes, compareNodeMemory)
		got := make([]string, 0, len(nodes))
		for _, n := range nodes {
			got = append(got, n.GetName())
		}
		if !slices.Equal(got, want) {
			t.Errorf("node order = %v, want %v", got, want)
		}
	})
	t.Run("pods", func(t *testing.T) {
		t.Parallel()
		pods := []*payload.Info_Pod{
			{Name: high, Memory: mem(5)},
			{Name: nilName},
			{Name: low, Memory: mem(1)},
		}
		slices.SortFunc(pods, comparePodMemory)
		got := make([]string, 0, len(pods))
		for _, p := range pods {
			got = append(got, p.GetName())
		}
		if !slices.Equal(got, want) {
			t.Errorf("pod order = %v, want %v", got, want)
		}
	})
}
