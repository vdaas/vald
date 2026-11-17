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

package service

import (
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// newBenchSnapshotDiscoverer publishes a snapshot of nodes x podsPerNode agent
// pods the way refreshSnapshots does (pods pre-sorted by memory usage).
func newBenchSnapshotDiscoverer(nodes, podsPerNode int) *discoverer {
	const ns, app = "vald", "vald-agent"
	nodeByName := make(map[string]*payload.Info_Node, nodes)
	podsByNode := make(map[string]map[string]map[string][]*payload.Info_Pod, nodes)
	for i := range nodes {
		name := "node-" + strconv.Itoa(i)
		n := &payload.Info_Node{
			Name:   name,
			Memory: &payload.Info_Memory{Usage: float64(i)},
			Pods:   new(payload.Info_Pods),
		}
		nodeByName[name] = n
		pods := make([]*payload.Info_Pod, 0, podsPerNode)
		for j := range podsPerNode {
			pods = append(pods, &payload.Info_Pod{
				AppName:   app,
				Name:      app + "-" + strconv.Itoa(i*podsPerNode+j),
				Namespace: ns,
				Ip:        "10.0." + strconv.Itoa(i) + "." + strconv.Itoa(j),
				Memory:    &payload.Info_Memory{Usage: float64(j)},
				Node:      n,
			})
		}
		podsByNode[name] = map[string]map[string][]*payload.Info_Pod{ns: {app: pods}}
	}
	d := new(discoverer)
	d.nodeByName.Store(&nodeByName)
	d.podsByNode.Store(&podsByNode)
	return d
}

// BenchmarkGetNodes_Poll measures one gateway discovery poll on the discoverer
// side (Nodes response build + marshal) for 300 agents, and reports the CPU it
// costs per second per gateway replica at the LB discoverer poll interval:
// 200ms (old default) and 3s (= discoverer.discovery_duration, the interval at
// which the answered snapshot can change at all).
func BenchmarkGetNodes_Poll(b *testing.B) {
	d := newBenchSnapshotDiscoverer(30, 10)
	req := &payload.Discoverer_Request{Namespace: "vald", Name: "vald-agent"}
	for _, poll := range []time.Duration{200 * time.Millisecond, 3 * time.Second} {
		b.Run("poll="+poll.String(), func(b *testing.B) {
			b.ResetTimer()
			for range b.N {
				nodes, err := d.GetNodes(req)
				if err != nil {
					b.Fatal(err)
				}
				if _, err := nodes.MarshalVT(); err != nil {
					b.Fatal(err)
				}
			}
			perPoll := float64(b.Elapsed().Nanoseconds()) / float64(b.N)
			b.ReportMetric(perPoll*float64(time.Second)/float64(poll)/1e3, "cpu-us/s")
		})
	}
}
