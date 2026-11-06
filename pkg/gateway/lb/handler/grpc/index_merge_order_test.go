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
package grpc

import (
	"slices"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// mergeInfoIndexStatistics must not depend on map iteration order: a minimum
// of 0 is a real value (not "unset"), and IndegreeCount is sized from the
// longest input. Every assignment of the inputs to keys is tried repeatedly so
// that Go's randomized map order visits the agents in many orders.
func Test_mergeInfoIndexStatistics_orderIndependent(t *testing.T) {
	t.Parallel()
	// fresh inputs per call: sumHistograms may reuse an input histogram slice.
	inputs := func() []*payload.Info_Index_Statistics {
		return []*payload.Info_Index_Statistics{
			{Valid: true, MinNumberOfIndegree: 0, MinNumberOfOutdegree: 3, IndegreeCount: []int64{2, 4}},
			{Valid: true, MinNumberOfIndegree: 5, MinNumberOfOutdegree: 1, IndegreeCount: []int64{4, 6, 9}},
			{Valid: true, MinNumberOfIndegree: 2, MinNumberOfOutdegree: 0, IndegreeCount: []int64{6}},
		}
	}
	const (
		wantMinIn  = uint64(0)
		wantMinOut = uint64(0)
	)
	// index 0: (2+4+6)/3, index 1: (4+6)/2, index 2: 9/1
	wantCount := []int64{4, 5, 9}

	perms := [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	keys := []string{"a", "b", "c"}
	for _, perm := range perms {
		for range 50 {
			in := inputs()
			stats := make(map[string]*payload.Info_Index_Statistics, len(perm))
			for i, p := range perm {
				stats[keys[i]] = in[p]
			}
			got := mergeInfoIndexStatistics(stats)
			if got.GetMinNumberOfIndegree() != wantMinIn {
				t.Fatalf("perm %v: MinNumberOfIndegree = %d, want %d", perm, got.GetMinNumberOfIndegree(), wantMinIn)
			}
			if got.GetMinNumberOfOutdegree() != wantMinOut {
				t.Fatalf("perm %v: MinNumberOfOutdegree = %d, want %d", perm, got.GetMinNumberOfOutdegree(), wantMinOut)
			}
			if !slices.Equal(got.GetIndegreeCount(), wantCount) {
				t.Fatalf("perm %v: IndegreeCount = %v, want %v", perm, got.GetIndegreeCount(), wantCount)
			}
		}
	}
}
