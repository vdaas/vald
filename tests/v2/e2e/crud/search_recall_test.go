//go:build e2e

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

// package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"math"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
)

// TestRecordRecall covers the recordRecall guard conditions (nil plan, nil
// Metrics, Metrics.Enabled == false, nil Collector) as well as the success
// path, mirroring the plan.Metrics.Enabled / plan.Collector != nil guard
// already used throughout grpc_test.go for Record().
func TestRecordRecall(t *testing.T) {
	t.Parallel()

	t.Run("nil plan does not panic and is a no-op", func(t *testing.T) {
		t.Parallel()
		recordRecall(nil, 0.5)
	})

	t.Run("nil Metrics does not panic and is a no-op", func(t *testing.T) {
		t.Parallel()
		c, err := metrics.NewCollector()
		if err != nil {
			t.Fatalf("metrics.NewCollector() error = %v", err)
		}
		plan := &config.Execution{Collector: c}
		recordRecall(plan, 0.5)
		if snap := c.GlobalSnapshot(); snap.Recalls != nil && snap.Recalls.Total != 0 {
			t.Errorf("recordRecall must not record when Metrics is nil, got Recalls=%+v", snap.Recalls)
		}
	})

	t.Run("Metrics disabled is a no-op", func(t *testing.T) {
		t.Parallel()
		c, err := metrics.NewCollector()
		if err != nil {
			t.Fatalf("metrics.NewCollector() error = %v", err)
		}
		plan := &config.Execution{Metrics: &config.Metrics{Enabled: false}, Collector: c}
		recordRecall(plan, 0.5)
		if snap := c.GlobalSnapshot(); snap.Recalls != nil && snap.Recalls.Total != 0 {
			t.Errorf("recordRecall must not record when Metrics.Enabled is false, got Recalls=%+v", snap.Recalls)
		}
	})

	t.Run("nil Collector does not panic and is a no-op", func(t *testing.T) {
		t.Parallel()
		plan := &config.Execution{Metrics: &config.Metrics{Enabled: true}}
		recordRecall(plan, 0.5)
	})

	t.Run("records into plan.Collector when enabled", func(t *testing.T) {
		t.Parallel()
		c, err := metrics.NewCollector()
		if err != nil {
			t.Fatalf("metrics.NewCollector() error = %v", err)
		}
		plan := &config.Execution{Metrics: &config.Metrics{Enabled: true}, Collector: c}
		recordRecall(plan, 0.75)
		snap := c.GlobalSnapshot()
		if snap.Recalls == nil || snap.Recalls.Total != 1 {
			t.Fatalf("recordRecall did not record into plan.Collector: Recalls=%+v", snap.Recalls)
		}
		if snap.Recalls.Mean != 0.75 {
			t.Errorf("Recalls.Mean = %v, want 0.75", snap.Recalls.Mean)
		}
	})
}

// TestCheckUnarySearchResponse_RecordsRecall exercises the wiring between
// calculateRecall and plan.Collector via checkUnarySearchResponse, without
// requiring a live gRPC connection: the callback is invoked directly with a
// hand-built *payload.Search_Response, matching how single() in grpc_test.go
// invokes it after a real gRPC call.
func TestCheckUnarySearchResponse_RecordsRecall(t *testing.T) {
	t.Parallel()

	c, err := metrics.NewCollector()
	if err != nil {
		t.Fatalf("metrics.NewCollector() error = %v", err)
	}
	plan := &config.Execution{Metrics: &config.Metrics{Enabled: true}, Collector: c}

	// 2 of the 4 returned IDs ("0","1") are within the expected neighbor set
	// {0,1,2,3}; the other 2 ("98","99") are not, so recall == 0.5.
	neighbors := iter.NewCycle([][]int{{0, 1, 2, 3}}, 1, 0, nil)
	res := &payload.Search_Response{
		RequestId: "0-Unknown",
		Results: []*payload.Object_Distance{
			{Id: "0"},
			{Id: "1"},
			{Id: "99"},
			{Id: "98"},
		},
	}

	cb := checkUnarySearchResponse(neighbors, plan)
	if !cb(t, 0, res, nil) {
		t.Fatal("checkUnarySearchResponse callback returned false")
	}

	wantRecall := recall(t, []string{"0", "1", "99", "98"}, []int{0, 1, 2, 3})
	snap := c.GlobalSnapshot()
	if snap.Recalls == nil || snap.Recalls.Total != 1 {
		t.Fatalf("checkUnarySearchResponse did not record recall into plan.Collector: Recalls=%+v", snap.Recalls)
	}
	if math.Abs(snap.Recalls.Mean-wantRecall) > 1e-9 {
		t.Errorf("Recalls.Mean = %v, want %v", snap.Recalls.Mean, wantRecall)
	}
}

// TestCheckUnarySearchResponse_NilPlan_DoesNotPanic guards against a
// regression where checkUnarySearchResponse's recall computation happens
// regardless of whether metrics collection is configured for the plan (e.g.
// plan == nil, as processSearch's callers may pass when metrics are
// disabled entirely).
func TestCheckUnarySearchResponse_NilPlan_DoesNotPanic(t *testing.T) {
	t.Parallel()

	neighbors := iter.NewCycle([][]int{{0, 1}}, 1, 0, nil)
	res := &payload.Search_Response{
		RequestId: "0-Unknown",
		Results:   []*payload.Object_Distance{{Id: "0"}},
	}

	cb := checkUnarySearchResponse(neighbors, nil)
	if !cb(t, 0, res, nil) {
		t.Fatal("checkUnarySearchResponse callback returned false")
	}
}
