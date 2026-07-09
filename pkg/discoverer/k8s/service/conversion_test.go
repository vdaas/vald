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
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/k8s"
)

func newRunningPod(name, app string) k8s.Pod {
	return k8s.Pod{
		ObjectMeta: k8s.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels:    map[string]string{"app": app},
		},
		Spec: k8s.PodSpec{
			NodeName:   "node-a",
			Containers: []k8s.Container{{Name: "vald-agent"}},
		},
		Status: k8s.PodStatus{
			Phase: k8s.PodRunning,
			PodIP: "10.0.0.1",
		},
	}
}

func TestPodsByAppName(t *testing.T) {
	t.Parallel()

	t.Run("groups pods by app label", func(t *testing.T) {
		t.Parallel()

		list := &k8s.PodList{Items: []k8s.Pod{
			newRunningPod("agent-0", "vald-agent"),
			newRunningPod("agent-1", "vald-agent"),
			newRunningPod("gateway-0", "vald-lb-gateway"),
		}}

		got := podsByAppName(list, "")
		if len(got) != 2 {
			t.Fatalf("podsByAppName() groups = %d, want 2", len(got))
		}
		if len(got["vald-agent"]) != 2 {
			t.Errorf("vald-agent pods = %d, want 2", len(got["vald-agent"]))
		}
		if len(got["vald-lb-gateway"]) != 1 {
			t.Errorf("vald-lb-gateway pods = %d, want 1", len(got["vald-lb-gateway"]))
		}
		if got["vald-agent"][0].IP != "10.0.0.1" || got["vald-agent"][0].NodeName != "node-a" {
			t.Errorf("unexpected pod conversion: %+v", got["vald-agent"][0])
		}
	})

	t.Run("skips non-running and namespace-mismatched pods", func(t *testing.T) {
		t.Parallel()

		pending := newRunningPod("pending-0", "vald-agent")
		pending.Status.Phase = k8s.PodPending
		other := newRunningPod("other-ns-0", "vald-agent")
		other.Namespace = "other"

		list := &k8s.PodList{Items: []k8s.Pod{
			pending,
			other,
			newRunningPod("agent-0", "vald-agent"),
		}}

		got := podsByAppName(list, "default")
		if len(got["vald-agent"]) != 1 {
			t.Errorf("vald-agent pods = %d, want 1", len(got["vald-agent"]))
		}
	})

	t.Run("does not over-reserve capacity per group", func(t *testing.T) {
		t.Parallel()

		// One large group plus one single-pod group: reserving
		// len(list.Items) for every group previously inflated the small
		// group's backing array to the full list size (~14x waste measured).
		items := make([]k8s.Pod, 0, 101)
		for i := range 100 {
			items = append(items, newRunningPod("agent-"+strconv.Itoa(i), "vald-agent"))
		}
		items = append(items, newRunningPod("gateway-0", "vald-lb-gateway"))
		list := &k8s.PodList{Items: items}

		got := podsByAppName(list, "")
		if n := len(got["vald-lb-gateway"]); n != 1 {
			t.Fatalf("vald-lb-gateway pods = %d, want 1", n)
		}
		if c := cap(got["vald-lb-gateway"]); c >= len(list.Items) {
			t.Errorf("single-pod group capacity = %d, want natural growth (< %d)", c, len(list.Items))
		}
	})
}
