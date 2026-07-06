package vald

import (
	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease"
	v1 "k8s.io/api/core/v1"
)

func (b *VrsBuilder) buildNodeSelector(nt controllerv1.NodePoolType) map[string]string {
	nt = b.effectiveNodePoolType(nt)
	return map[string]string{
		b.labelKey(nodePoolLabelNamespace): b.CR.GetNamespace(),
		b.labelKey(nodePoolLabelType):      string(nt),
	}
}

func (b *VrsBuilder) buildToleration(nt controllerv1.NodePoolType) *[]v1.Toleration {
	nt = b.effectiveNodePoolType(nt)
	return &[]v1.Toleration{
		{
			Key:      b.labelKey(nodePoolLabelNamespace),
			Operator: v1.TolerationOpEqual,
			Value:    b.CR.GetNamespace(),
			Effect:   v1.TaintEffectNoSchedule,
		},
		{
			Key:      b.labelKey(nodePoolLabelType),
			Operator: v1.TolerationOpEqual,
			Value:    string(nt),
		},
	}
}

// effectiveNodePoolType returns the pool type to use after applying the
// agent → general fallback. The decision is a pure function of the
// pre-resolved Capability — no I/O, no context.
func (b *VrsBuilder) effectiveNodePoolType(nt controllerv1.NodePoolType) controllerv1.NodePoolType {
	if nt == controllerv1.NodePoolTypeValdAgent && !b.Capability.HasAgentPool {
		return controllerv1.NodePoolTypeGeneral
	}
	return nt
}

// applyNodeAffinities sets NodeSelector and Tolerations on every component
// of the ValdRelease. Agent goes to the (possibly-fallback'd) ValdAgent pool;
// every other component (gateway lb, discoverer, manager index + creator +
// saver) goes to the General pool. This replaces a previous reflect-based
// implementation that decided pool placement from the runtime struct name —
// renaming any of those types silently lost the affinity, which the explicit
// listing here makes impossible.
func (b *VrsBuilder) applyNodeAffinities(v *valdrelease.ValdRelease) {
	agentNS := b.buildNodeSelector(controllerv1.NodePoolTypeValdAgent)
	agentTol := b.buildToleration(controllerv1.NodePoolTypeValdAgent)
	generalNS := b.buildNodeSelector(controllerv1.NodePoolTypeGeneral)
	generalTol := b.buildToleration(controllerv1.NodePoolTypeGeneral)

	v.Spec.Agent.NodeSelector = agentNS
	v.Spec.Agent.Tolerations = agentTol

	v.Spec.Gateway.Lb.NodeSelector = generalNS
	v.Spec.Gateway.Lb.Tolerations = generalTol

	v.Spec.Discoverer.NodeSelector = generalNS
	v.Spec.Discoverer.Tolerations = generalTol

	v.Spec.Manager.Index.NodeSelector = generalNS
	v.Spec.Manager.Index.Tolerations = generalTol

	if v.Spec.Manager.Index.Saver != nil {
		v.Spec.Manager.Index.Saver.NodeSelector = generalNS
		v.Spec.Manager.Index.Saver.Tolerations = generalTol
	}
	if v.Spec.Manager.Index.Creator != nil {
		v.Spec.Manager.Index.Creator.NodeSelector = generalNS
		v.Spec.Manager.Index.Creator.Tolerations = generalTol
	}
}
