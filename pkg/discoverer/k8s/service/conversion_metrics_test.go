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

// Package service manages the main logic of server.
package service

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/k8s"
)

// mustPodMetricsList builds a k8s.PodMetricsList from raw JSON. Quantity
// fields (cpu/memory) use the same string encoding the metrics API returns
// ("500m", "128Mi", ...); this sidesteps importing k8s.io/metrics/apis/v1beta1
// directly for the unexported-from-internal/k8s ContainerMetrics type, which
// pkg/ may not import per depguard's k8s-confinement rule.
func mustPodMetricsList(t *testing.T, raw string) *k8s.PodMetricsList {
	t.Helper()
	list := new(k8s.PodMetricsList)
	if err := json.Unmarshal([]byte(raw), list); err != nil {
		t.Fatalf("unmarshal PodMetricsList: %v", err)
	}
	return list
}

func mustNodeMetricsList(t *testing.T, raw string) *k8s.NodeMetricsList {
	t.Helper()
	list := new(k8s.NodeMetricsList)
	if err := json.Unmarshal([]byte(raw), list); err != nil {
		t.Fatalf("unmarshal NodeMetricsList: %v", err)
	}
	return list
}

func mustPodMetrics(t *testing.T, raw string) *k8s.PodMetrics {
	t.Helper()
	pm := new(k8s.PodMetrics)
	if err := json.Unmarshal([]byte(raw), pm); err != nil {
		t.Fatalf("unmarshal PodMetrics: %v", err)
	}
	return pm
}

func TestToPodMetricsMap(t *testing.T) {
	t.Parallel()

	t.Run("empty list returns empty map", func(t *testing.T) {
		t.Parallel()

		got := toPodMetricsMap(mustPodMetricsList(t, `{"items":[]}`))
		if len(got) != 0 {
			t.Errorf("toPodMetricsMap() = %#v, want empty map", got)
		}
	})

	t.Run("averages usage across containers, rounding each container up to whole units first", func(t *testing.T) {
		t.Parallel()

		// resource.Quantity.Value() rounds up to the nearest whole unit away
		// from zero, and that rounding happens per container before the sum
		// is averaged: 500m CPU -> 1, 1500m CPU -> 2, so (1+2)/2 = 1.5, not
		// (0.5+1.5)/2 = 1.0. Memory is already whole bytes, so no rounding
		// changes its sum: (128Mi + 1Ki) / 2 = 67109376.
		list := mustPodMetricsList(t, `{"items":[{
			"metadata":{"name":"agent-0","namespace":"metrics-ns-a"},
			"containers":[
				{"name":"vald-agent","usage":{"cpu":"500m","memory":"128Mi"}},
				{"name":"sidecar","usage":{"cpu":"1500m","memory":"1Ki"}}
			]
		}]}`)

		got := toPodMetricsMap(list)
		want := map[string]PodMetrics{
			"agent-0": {Name: "agent-0", Namespace: "metrics-ns-a", CPU: 1.5, Mem: 67109376},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("toPodMetricsMap() = %#v, want %#v", got, want)
		}
	})

	t.Run("keys multiple pods by name", func(t *testing.T) {
		t.Parallel()

		list := mustPodMetricsList(t, `{"items":[
			{"metadata":{"name":"p1","namespace":"ns1"},"containers":[{"name":"c","usage":{"cpu":"2","memory":"256Mi"}}]},
			{"metadata":{"name":"p2","namespace":"ns2"},"containers":[{"name":"c","usage":{"cpu":"4","memory":"512Mi"}}]}
		]}`)

		got := toPodMetricsMap(list)
		want := map[string]PodMetrics{
			"p1": {Name: "p1", Namespace: "ns1", CPU: 2, Mem: 268435456},
			"p2": {Name: "p2", Namespace: "ns2", CPU: 4, Mem: 536870912},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("toPodMetricsMap() = %#v, want %#v", got, want)
		}
	})

	t.Run("same pod name across different namespaces collides, last item in the list wins", func(t *testing.T) {
		t.Parallel()

		// toPodMetricsMap keys solely by pod name, not (namespace, name); it
		// does not distinguish same-named pods living in different
		// namespaces, so a later list entry silently overwrites an earlier
		// one that shares its name.
		list := mustPodMetricsList(t, `{"items":[
			{"metadata":{"name":"dup","namespace":"ns-a"},"containers":[{"name":"c","usage":{"cpu":"1"}}]},
			{"metadata":{"name":"dup","namespace":"ns-b"},"containers":[{"name":"c","usage":{"cpu":"2"}}]}
		]}`)

		got := toPodMetricsMap(list)
		if len(got) != 1 {
			t.Fatalf("toPodMetricsMap() groups = %d, want 1", len(got))
		}
		want := PodMetrics{Name: "dup", Namespace: "ns-b", CPU: 2, Mem: 0}
		if got["dup"] != want {
			t.Errorf("dup entry = %+v, want %+v", got["dup"], want)
		}
	})

	t.Run("container missing a resource key contributes zero for that resource", func(t *testing.T) {
		t.Parallel()

		// corev1.ResourceList.Cpu()/.Memory() return a zero-value Quantity
		// when the key is absent rather than erroring, so an omitted "cpu"
		// key contributes 0 to the average, not NaN or a skipped container.
		list := mustPodMetricsList(t, `{"items":[{
			"metadata":{"name":"p1","namespace":"metrics-ns-b"},
			"containers":[{"name":"c","usage":{"memory":"128Mi"}}]
		}]}`)

		got := toPodMetricsMap(list)
		want := map[string]PodMetrics{
			"p1": {Name: "p1", Namespace: "metrics-ns-b", CPU: 0, Mem: 134217728},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("toPodMetricsMap() = %#v, want %#v", got, want)
		}
	})

	t.Run("pod with zero containers yields NaN averages", func(t *testing.T) {
		t.Parallel()

		// Both accumulators start and stay at 0.0 when there are no
		// containers to sum, so dividing by len(pod.Containers) == 0 is a
		// 0/0 float division: NaN, not a panic or a zero-valued PodMetrics.
		list := mustPodMetricsList(t, `{"items":[{"metadata":{"name":"p1"},"containers":[]}]}`)

		got := toPodMetricsMap(list)
		pm, ok := got["p1"]
		if !ok {
			t.Fatalf("toPodMetricsMap() missing entry for p1: %#v", got)
		}
		if !math.IsNaN(pm.CPU) || !math.IsNaN(pm.Mem) {
			t.Errorf("p1 = %+v, want CPU and Mem both NaN", pm)
		}
	})
}

func TestToNodeMetricsMap(t *testing.T) {
	t.Parallel()

	t.Run("empty list returns empty map", func(t *testing.T) {
		t.Parallel()

		got := toNodeMetricsMap(mustNodeMetricsList(t, `{"items":[]}`))
		if len(got) != 0 {
			t.Errorf("toNodeMetricsMap() = %#v, want empty map", got)
		}
	})

	t.Run("converts cpu, memory, ephemeral-storage and pods usage without averaging", func(t *testing.T) {
		t.Parallel()

		// Unlike PodMetrics, NodeMetrics has no per-container loop: each
		// field is read straight off the single Usage ResourceList.
		list := mustNodeMetricsList(t, `{"items":[{
			"metadata":{"name":"`+testNodeAName+`"},
			"usage":{"cpu":"2","memory":"4Gi","ephemeral-storage":"10Gi","pods":"110"}
		}]}`)

		got := toNodeMetricsMap(list)
		want := map[string]NodeMetrics{
			testNodeAName: {Name: testNodeAName, CPU: 2, Mem: 4294967296, Storage: 10737418240, Pods: 110},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("toNodeMetricsMap() = %#v, want %#v", got, want)
		}
	})

	t.Run("node missing a resource key defaults to a zero quantity for that resource", func(t *testing.T) {
		t.Parallel()

		list := mustNodeMetricsList(t, `{"items":[{"metadata":{"name":"`+testNodeAName+`"},"usage":{"memory":"4Gi"}}]}`)

		got := toNodeMetricsMap(list)
		want := map[string]NodeMetrics{
			testNodeAName: {Name: testNodeAName, CPU: 0, Mem: 4294967296, Storage: 0, Pods: 0},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("toNodeMetricsMap() = %#v, want %#v", got, want)
		}
	})

	t.Run("duplicate node names collide, last item in the list wins", func(t *testing.T) {
		t.Parallel()

		list := mustNodeMetricsList(t, `{"items":[
			{"metadata":{"name":"dup"},"usage":{"cpu":"1"}},
			{"metadata":{"name":"dup"},"usage":{"cpu":"5"}}
		]}`)

		got := toNodeMetricsMap(list)
		if len(got) != 1 {
			t.Fatalf("toNodeMetricsMap() groups = %d, want 1", len(got))
		}
		want := NodeMetrics{Name: "dup", CPU: 5}
		if got["dup"] != want {
			t.Errorf("dup entry = %+v, want %+v", got["dup"], want)
		}
	})
}

func TestPodMetricsContainersNameIndexer(t *testing.T) {
	t.Parallel()

	t.Run("extracts container names in declaration order", func(t *testing.T) {
		t.Parallel()

		pm := mustPodMetrics(t, `{"metadata":{"name":"p1"},"containers":[
			{"name":"a"},{"name":"b"},{"name":"c"}
		]}`)

		got := podMetricsContainersNameIndexer(pm)
		want := []string{"a", "b", "c"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("podMetricsContainersNameIndexer() = %#v, want %#v", got, want)
		}
	})

	t.Run("pod metrics with no containers yields a non-nil empty slice", func(t *testing.T) {
		t.Parallel()

		pm := mustPodMetrics(t, `{"metadata":{"name":"p1"},"containers":[]}`)

		got := podMetricsContainersNameIndexer(pm)
		if got == nil {
			t.Error("podMetricsContainersNameIndexer() = nil, want non-nil empty slice")
		}
		if len(got) != 0 {
			t.Errorf("podMetricsContainersNameIndexer() = %#v, want empty", got)
		}
	})

	t.Run("non PodMetrics object yields nil", func(t *testing.T) {
		t.Parallel()

		got := podMetricsContainersNameIndexer(&k8s.Pod{})
		if got != nil {
			t.Errorf("podMetricsContainersNameIndexer() = %#v, want nil", got)
		}
	})
}
