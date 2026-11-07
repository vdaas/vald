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

// Package service manages the main logic of server.
package service

import (
	"testing"

	"github.com/vdaas/vald/internal/k8s"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
)

// TestDiscovererListControllers_clusterScoped verifies that the Node and
// NodeMetrics lists are not restricted to the pod namespace. Both resources are
// cluster-scoped (namespace ""), so an InNamespace filter listed nothing and
// left d.nodes / d.nodeMetrics empty.
func TestDiscovererListControllers_clusterScoped(t *testing.T) {
	t.Parallel()
	const nodeName = "node-a"
	tests := []struct {
		obj    k8s.Object
		stored func(d *discoverer) bool
		name   string
	}{
		{
			name: "node discoverer",
			obj:  &k8s.Node{Name: nodeName},
			stored: func(d *discoverer) bool {
				_, ok := d.nodes.Load(nodeName)
				return ok
			},
		},
		{
			name: "node metrics discoverer",
			obj:  &k8s.NodeMetrics{Name: nodeName},
			stored: func(d *discoverer) bool {
				_, ok := d.nodeMetrics.Load(nodeName)
				return ok
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			scheme := k8s.NewScheme()
			if err := k8s.AddCoreToScheme(scheme); err != nil {
				t.Fatal(err)
			}
			if err := k8s.AddMetricsToScheme(scheme); err != nil {
				t.Fatal(err)
			}
			mgr := mock.NewListManagerMock(
				mock.NewFakeClientBuilder().WithScheme(scheme).WithObjects(test.obj).Build(),
			)

			d := &discoverer{namespace: "vald"}
			var found bool
			for _, rc := range d.listControllers(nil) {
				if rc.GetName() != test.name {
					continue
				}
				found = true
				r := rc.NewReconciler(t.Context(), mgr)
				if _, err := r.Reconcile(t.Context(), k8s.Request{}); err != nil {
					t.Fatalf("Reconcile() error = %v", err)
				}
			}
			if !found {
				t.Fatalf("no list controller named %q", test.name)
			}
			if !test.stored(d) {
				t.Errorf("%s did not store the cluster-scoped %s while the pod namespace is %q",
					test.name, nodeName, d.namespace)
			}
		})
	}
}
