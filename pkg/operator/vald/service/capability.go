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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
)

const (
	nodePoolLabelNamespace = "namespace"
	nodePoolLabelType      = "type"
	nodePoolLabelRole      = "role"
)

// labelKey renders a NodePool label key with the configured prefix. Exposed
// as a free function (rather than a vrsBuilder method) because tests build
// matching labels on fake Node objects before a builder exists; both call
// sites must use the same logic and prefix value.
func labelKey(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	return prefix + "/" + suffix
}

// nodePoolCapability represents which node-pool types are available in the
// cluster for a given namespace. It is a pure value: callers resolve it once
// per reconcile and pass it into the builder, so the builder itself does not
// need a Kubernetes client to produce its output.
type nodePoolCapability struct {
	HasGeneralPool bool
	HasAgentPool   bool
}

// alwaysAvailable returns a capability where both pools are present. This is
// the right default when node_pool.require_match is disabled: the controller
// skips node listing entirely and the builder treats every infra entry as
// schedulable.
func alwaysAvailable() nodePoolCapability {
	return nodePoolCapability{HasGeneralPool: true, HasAgentPool: true}
}

// resolveNodePoolCapability inspects the cluster's Nodes via the given client
// and reports which node pool types are present for the namespace. The label
// prefix must match the one the builder uses when emitting NodeSelector and
// Tolerations. The Nodes are listed once, filtered by the namespace label,
// and the pool type is classified in Go so the reconcile hot path issues a
// single List call instead of one per pool type.
func resolveNodePoolCapability(
	ctx context.Context, c k8s.Client, namespace, labelPrefix string,
) (nodePoolCapability, error) {
	if c == nil {
		return nodePoolCapability{}, errors.ErrK8sClientIsNil
	}
	nodes, err := resource.ListObjects(ctx, c, &k8s.NodeList{}, k8s.MatchingLabels{
		labelKey(labelPrefix, nodePoolLabelNamespace): namespace,
	})
	if err != nil {
		return nodePoolCapability{}, errors.Wrap(err, "list nodes")
	}

	var capability nodePoolCapability
	typeKey := labelKey(labelPrefix, nodePoolLabelType)
	for i := range nodes.Items {
		switch v1.NodePoolType(nodes.Items[i].Labels[typeKey]) {
		case v1.NodePoolTypeGeneral:
			capability.HasGeneralPool = true
		case v1.NodePoolTypeValdAgent:
			capability.HasAgentPool = true
		}
		if capability.HasGeneralPool && capability.HasAgentPool {
			break
		}
	}
	return capability, nil
}
