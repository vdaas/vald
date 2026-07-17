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

package service

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/metadata"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease"
)

// newBenchRow replicates the per-infra/per-cluster row construction performed
// inside Build so BenchmarkMergeOverlay can exercise mergeOverlay in isolation
// against a realistic, fully-populated row.
func newBenchRow(tb testing.TB, b *vrsBuilder) *valdrelease.ValdRelease {
	tb.Helper()

	infra := b.cr.Spec.Infrastructure[0]

	row := &valdrelease.ValdRelease{}
	row.SetGroupVersionKind(valdrelease.GVK)
	row.SetNamespace(b.cr.Namespace)
	row.SetLabels(metadata.CreateSubResourceLabels(valdrelease.GVK.Kind))

	row.Spec = valdrelease.Values{
		Defaults:   b.buildDefaults(),
		Gateway:    b.buildGateway(),
		Agent:      b.buildAgent(),
		Manager:    b.buildManager(),
		Discoverer: b.buildDiscoverer(),
	}

	agentPool := resolveAgentNodePool(infra)
	row.SetRelationalResources(agentPool.NodeCount, agentPool.MachineResource, valdrelease.ResourceParams{
		AgentPodsPerNode:           b.cfg.AgentPodsPerNode,
		DiscovererDSMaxSurge:       b.cfg.DiscovererDaemonSetMaxSurge,
		DiscovererDSMaxUnavailable: b.cfg.DiscovererDaemonSetMaxUnavailable,
	})
	b.reflectPersistentVolume(row)
	b.applyNodeAffinities(row)

	name, err := b.makeName(b.cr.GetNamespace(), infra.Clusters[0].Name)
	if err != nil {
		tb.Fatal(err)
	}
	row.SetName(name)
	row.SetLabels(mergeLabels(row.GetLabels(), b.buildLabels(0)))
	return row
}

// BenchmarkMergeOverlay measures the per-row overlay merge pipeline (row →
// unstructured conversion + default VRS template + built row + CR overlay
// patch). The conversion is included so the numbers stay comparable with the
// pre-hoisting pipeline, where every cluster paid for it.
func BenchmarkMergeOverlay(b *testing.B) {
	cfg := initConfig(b)
	cr := loadValdOperatorReleaseFromYAML(b)

	builder := newVrsBuilder(cr, cfg, alwaysAvailable())
	if err := builder.parseOverlay(); err != nil {
		b.Fatal(err)
	}
	row := newBenchRow(b, builder)

	b.ReportAllocs()
	for b.Loop() {
		current, err := resource.ToUnstructured(row)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := builder.mergeOverlay(current); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuild measures the full Build loop (validation, overlay parse, row
// construction and the per-cluster merge) exactly as the reconciler drives it.
func BenchmarkBuild(b *testing.B) {
	cfg := initConfig(b)
	cr := loadValdOperatorReleaseFromYAML(b)
	ctx := context.Background()

	b.ReportAllocs()
	for b.Loop() {
		if _, err := newVrsBuilder(cr, cfg, alwaysAvailable()).Build(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
