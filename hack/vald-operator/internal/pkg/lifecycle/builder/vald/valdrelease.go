package vald

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/metadata"

	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	nodePoolLabelNamespace = "namespace"
	nodePoolLabelType      = "type"
	nodePoolLabelRole      = "role"
)

type VrsBuilder struct {
	GVK        schema.GroupVersionKind
	List       *unstructured.UnstructuredList
	CR         *controllerv1.ValdOperatorRelease
	Config     *config.Config
	Capability NodePoolCapability
	Rules      DomainRules
}

func NewVrsBuilder(cr *controllerv1.ValdOperatorRelease, cfg *config.Config, capability NodePoolCapability, rules DomainRules) *VrsBuilder {
	return &VrsBuilder{
		GVK:        valdrelease.GVK,
		List:       &unstructured.UnstructuredList{},
		CR:         cr,
		Config:     cfg,
		Capability: capability,
		Rules:      rules,
	}
}

// Build produces the desired ValdRelease objects. It is a pure function of
// (CR, Config, Capability): no Kubernetes API calls are made here. The caller
// is responsible for resolving the NodePoolCapability before constructing the
// Builder (see ResolveNodePoolCapability).
func (b *VrsBuilder) Build(_ context.Context) (client.ObjectList, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	resourceParams := valdrelease.ResourceParams{
		AgentPodsPerNode:           b.Config.AgentPodsPerNode,
		DiscovererDSMaxSurge:       b.Config.DiscovererDaemonSetMaxSurge,
		DiscovererDSMaxUnavailable: b.Config.DiscovererDaemonSetMaxUnavailable,
	}

	for i, infra := range b.CR.Spec.Infrastructure {
		if !infra.Active {
			continue
		}

		if b.Config.RequireNodePoolMatch && !b.Capability.HasGeneralPool {
			continue
		}

		// At first, set the metadata fields.
		b.List.SetGroupVersionKind(b.GVK)
		row := &valdrelease.ValdRelease{}
		row.SetGroupVersionKind(b.GVK)
		row.SetNamespace(b.CR.Namespace)
		row.SetLabels(metadata.CreateSubResourceLabels(valdrelease.GVK.Kind))

		// Next, set the basic fields by input (b.CR)
		row.Spec = valdrelease.Spec{
			Defaults:   b.buildDefaults(),
			Gateway:    b.buildGateway(),
			Agent:      b.buildAgent(),
			Manager:    b.buildManager(),
			Discoverer: b.buildDiscoverer(),
		}

		// Resolve the agent's NodePool spec — the general-pool fallback rule
		// lives in the domain layer.
		agentPool := b.Rules.ResolveAgentNodePool(infra)
		row.SetRelationalResources(agentPool.NodeCount, agentPool.MachineResource, resourceParams)

		// Reflect optional settings (PV, node affinities) after resources are confirmed.
		b.reflectPersistentVolume(row)
		b.applyNodeAffinities(row)

		// Finally, reflect the settings for each cluster and apply the default values and Overlay patches.
		for _, cluster := range infra.Clusters {
			name, err := b.makeName(b.CR.GetNamespace(), cluster.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to make name for cluster %s: %w", cluster.Name, err)
			}
			row.SetName(name)
			row.SetLabels(b.buildLabels(i))
			u, err := b.mergeOverlay(row)
			if err != nil {
				return nil, fmt.Errorf("failed to merge overlay: %w", err)
			}

			b.List.Items = append(b.List.Items, *u)
		}
	}

	b.List.SetGroupVersionKind(valdrelease.GVK)
	return b.List, nil
}

func (b *VrsBuilder) validate() error {
	if b.CR == nil {
		return fmt.Errorf("CR is nil")
	}

	spec := b.CR.Spec
	if len(spec.Infrastructure) == 0 {
		return fmt.Errorf("CR.Spec.Infrastructure is empty")
	}
	for _, infra := range spec.Infrastructure {
		if len(infra.Clusters) == 0 {
			return fmt.Errorf("CR.Spec.Infrastructure.Clusters is empty for role %s", infra.Role)
		}
		for _, cluster := range infra.Clusters {
			if cluster.Name == "" {
				return fmt.Errorf("CR.Spec.Infrastructure.Clusters.Name is empty for role %s", infra.Role)
			}
		}
	}
	return nil
}
