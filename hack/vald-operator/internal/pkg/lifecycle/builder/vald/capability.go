package vald

import (
	"context"
	"fmt"

	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NodePoolCapability represents which node-pool types are available in the
// cluster for a given namespace. It is a pure value: callers resolve it once
// per reconcile and pass it into the Builder, so the Builder itself does not
// need a Kubernetes client to produce its output.
type NodePoolCapability struct {
	HasGeneralPool bool
	HasAgentPool   bool
}

// AlwaysAvailable returns a Capability where both pools are present. This is
// the right default when REQUIRE_NODEPOOL_MATCH is disabled: the controller
// skips node listing entirely and the Builder treats every infra entry as
// schedulable, matching the pre-refactor behavior.
func AlwaysAvailable() NodePoolCapability {
	return NodePoolCapability{HasGeneralPool: true, HasAgentPool: true}
}

// ResolveNodePoolCapability inspects the cluster's Nodes via the given client
// and reports which node pool types are present for the namespace. The label
// prefix must match the one the Builder uses when emitting NodeSelector and
// Tolerations.
func ResolveNodePoolCapability(ctx context.Context, c client.Client, namespace, labelPrefix string) (NodePoolCapability, error) {
	if c == nil {
		return NodePoolCapability{}, fmt.Errorf("k8s client is nil")
	}
	hasGeneral, err := hasNodesForType(ctx, c, namespace, labelPrefix, controllerv1.NodePoolTypeGeneral)
	if err != nil {
		return NodePoolCapability{}, fmt.Errorf("list nodes (general): %w", err)
	}
	hasAgent, err := hasNodesForType(ctx, c, namespace, labelPrefix, controllerv1.NodePoolTypeValdAgent)
	if err != nil {
		return NodePoolCapability{}, fmt.Errorf("list nodes (agent): %w", err)
	}
	return NodePoolCapability{HasGeneralPool: hasGeneral, HasAgentPool: hasAgent}, nil
}

func hasNodesForType(ctx context.Context, c client.Client, namespace, labelPrefix string, nt controllerv1.NodePoolType) (bool, error) {
	nodes := &v1.NodeList{}
	err := c.List(ctx, nodes, client.MatchingLabels{
		labelKey(labelPrefix, nodePoolLabelNamespace): namespace,
		labelKey(labelPrefix, nodePoolLabelType):      string(nt),
	})
	if err != nil {
		return false, err
	}
	return len(nodes.Items) > 0, nil
}
