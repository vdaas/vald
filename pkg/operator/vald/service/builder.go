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
	"encoding/json"
	"maps"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/kustomize"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/metadata"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/agent"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/defaults"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/discoverer"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/gateway"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/manager"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/operator/vald/config"
)

// maxK8sNameLength is the maximum number of characters allowed in a Kubernetes
// resource name (RFC 1123 / Kubernetes validation).
const maxK8sNameLength = 63

type vrsBuilder struct {
	list *valdrelease.ValdReleaseList
	cr   *v1.ValdOperatorRelease
	cfg  *config.Config
	// overlay is the CR-supplied overlay patch decoded once per Build; per
	// cluster it is DeepCopy'd instead of re-unmarshalling Overlay.Raw, which
	// is a reconcile hot path (raw size × clusters × requeue frequency).
	overlay    *k8s.Unstructured
	gvk        k8s.GroupVersionKind
	capability nodePoolCapability
}

func newVrsBuilder(
	cr *v1.ValdOperatorRelease, cfg *config.Config, capability nodePoolCapability,
) *vrsBuilder {
	return &vrsBuilder{
		gvk:        valdrelease.GVK,
		list:       &valdrelease.ValdReleaseList{},
		cr:         cr,
		cfg:        cfg,
		capability: capability,
	}
}

// Build produces the desired ValdRelease objects. It is a pure function of
// (CR, Config, Capability): no Kubernetes API calls are made here. The caller
// is responsible for resolving the nodePoolCapability before constructing the
// builder (see resolveNodePoolCapability).
func (b *vrsBuilder) Build(_ context.Context) ([]k8s.Object, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	// Decode the overlay patch once up front; makeOverlayPatch clones it for
	// every cluster instead of re-parsing the raw JSON on each iteration.
	if err := b.parseOverlay(); err != nil {
		return nil, errors.Wrap(err, "failed to parse overlay")
	}

	resourceParams := valdrelease.ResourceParams{
		AgentPodsPerNode:           b.cfg.AgentPodsPerNode,
		DiscovererDSMaxSurge:       b.cfg.DiscovererDaemonSetMaxSurge,
		DiscovererDSMaxUnavailable: b.cfg.DiscovererDaemonSetMaxUnavailable,
	}

	for i, infra := range b.cr.Spec.Infrastructure {
		if !infra.Active {
			continue
		}

		if b.cfg.RequireNodePoolMatch && !b.capability.HasGeneralPool {
			continue
		}

		row := &valdrelease.ValdRelease{}
		row.SetGroupVersionKind(b.gvk)
		row.SetNamespace(b.cr.Namespace)
		row.SetLabels(metadata.CreateSubResourceLabels(valdrelease.GVK.Kind))

		row.Spec = valdrelease.Spec{
			Defaults:   b.buildDefaults(),
			Gateway:    b.buildGateway(),
			Agent:      b.buildAgent(),
			Manager:    b.buildManager(),
			Discoverer: b.buildDiscoverer(),
		}

		// Resolve the agent's NodePool spec, applying the general-pool
		// fallback rule.
		agentPool := resolveAgentNodePool(infra)
		row.SetRelationalResources(agentPool.NodeCount, agentPool.MachineResource, resourceParams)

		// Reflect optional settings (PV, node affinities) after resources are confirmed.
		b.reflectPersistentVolume(row)
		b.applyNodeAffinities(row)

		for _, cluster := range infra.Clusters {
			name, err := b.makeName(b.cr.GetNamespace(), cluster.Name)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to make name for cluster %s", cluster.Name)
			}
			row.SetName(name)
			// Merge the per-infra labels on top of the existing ones so the
			// managed-by / managed-resource labels set above survive.
			row.SetLabels(mergeLabels(row.GetLabels(), b.buildLabels(i)))
			u, err := b.mergeOverlay(row)
			if err != nil {
				return nil, errors.Wrap(err, "failed to merge overlay")
			}

			b.list.Items = append(b.list.Items, *u)
		}
	}

	b.list.SetGroupVersionKind(valdrelease.GVK)
	return resource.ObjectsOf(b.list.Items), nil
}

// fetchExistingVrs lists the ValdRelease objects currently present in the
// namespace, restoring the GVK the client may have stripped, so the syncer
// can prune the ones the current build no longer produces.
func fetchExistingVrs(ctx context.Context, c k8s.Client, namespace string) ([]k8s.Object, error) {
	exists := &valdrelease.ValdReleaseList{}
	exists.SetGroupVersionKind(valdrelease.GVK)

	if _, err := resource.ListObjects(ctx, c, exists, k8s.InNamespace(namespace)); err != nil {
		return nil, errors.Wrap(err, "failed to list existing ValdRelease objects")
	}

	out := resource.ObjectsOf(exists.Items)
	for _, obj := range out {
		// Recover GVK if client stripped it
		if obj.GetObjectKind().GroupVersionKind().Empty() {
			obj.GetObjectKind().SetGroupVersionKind(valdrelease.GVK)
		}
	}
	return out, nil
}

func (b *vrsBuilder) validate() error {
	if b.cr == nil {
		return errors.ErrCRIsNil
	}

	spec := b.cr.Spec
	if len(spec.Infrastructure) == 0 {
		return errors.ErrCRSpecInfrastructureIsEmpty
	}
	for _, infra := range spec.Infrastructure {
		if len(infra.Clusters) == 0 {
			return errors.ErrCRSpecInfrastructureClustersIsEmptyForRole(string(infra.Role))
		}
		for _, cluster := range infra.Clusters {
			if cluster.Name == "" {
				return errors.ErrCRSpecInfrastructureClustersNameIsEmptyForRole(string(infra.Role))
			}
		}
	}
	return nil
}

func (b *vrsBuilder) buildDefaults() defaults.Defaults {
	vald := b.cr.Spec.VectorEngine.Vald
	return defaults.Defaults{
		Logging: b.buildLogging(vald.Defaults.LogLevel),
	}
}

func (b *vrsBuilder) buildAgent() agent.Agent {
	input := b.cr.Spec.VectorEngine.Vald.Agent

	mu := k8s.IntOrStringFrom(b.cfg.AgentMaxUnavailable)
	ms := k8s.IntOrStringFrom(b.cfg.AgentMaxSurge)

	a := agent.Agent{
		Logging: b.buildLogging(input.LogLevel),
		Kind:    common.KindTypeStatefulSet,
		RollingUpdate: &discoverer.RollingUpdateValdRelease{
			MaxUnavailable: &mu,
			MaxSurge:       &ms,
		},
		NGT: *b.buildAgentNgt(),
	}
	return a
}

func (b *vrsBuilder) buildAgentNgt() *agent.NGT {
	input := b.cr.Spec.VectorEngine.Vald.Agent.Ngt
	return &agent.NGT{
		Dimension:          input.Dimension,
		DistanceType:       input.DistanceType,
		ObjectType:         input.ObjectType,
		SearchEdgeSize:     input.SearchEdgeSize,
		CreationEdgeSize:   input.CreationEdgeSize,
		EnableInMemoryMode: b.cfg.AgentEnableInMemoryMode,
	}
}

func (b *vrsBuilder) buildGateway() gateway.Gateway {
	return gateway.Gateway{
		Lb: *b.buildLb(),
	}
}

func (b *vrsBuilder) buildLb() *gateway.Lb {
	inputGw := b.cr.Spec.VectorEngine.Vald.Gateway

	return &gateway.Lb{
		Logging: b.buildLogging(inputGw.LogLevel),
		Hpa: &gateway.Hpa{
			TargetCPUUtilizationPercentage: b.cfg.GatewayHpaTargetCPUUtilization,
		},
		GatewayConfig: gateway.GatewayConfig{
			IndexReplica: inputGw.IndexReplica,
		},
		ServiceType: b.getGatewayServiceType(inputGw.ServiceType),
		Ingress:     b.buildIngress(inputGw.Ingress),
	}
}

// buildIngress renders the gateway ingress. cfg.EnableIngress is the
// operator-level gate: when it is false, the ingress stays disabled even if
// the CR spec enables it. cfg.GatewayIngressAnnotations seed the ingress
// annotations; the CR overlay is merged on top of the built row afterwards,
// so CR-side annotations win per key.
func (b *vrsBuilder) buildIngress(in *v1.GatewayIngress) gateway.Ingress {
	base := gateway.Ingress{
		DefaultBackend: gateway.DefaultBackend{Enabled: false},
		PathType:       b.getGatewayIngressPathType(),
		ServicePort:    b.cfg.GatewayIngressServicePort,
	}
	if len(b.cfg.GatewayIngressAnnotations) > 0 {
		base.Annotations = mergeLabels(nil, b.cfg.GatewayIngressAnnotations)
	}
	if !b.cfg.EnableIngress || in == nil || !in.Enabled {
		return base
	}
	base.Enabled = true
	base.Host = in.Host
	return base
}

func (b *vrsBuilder) getGatewayIngressPathType() k8s.PathType {
	switch b.cfg.GatewayIngressPathType {
	case string(k8s.PathTypeExact):
		return k8s.PathTypeExact
	case string(k8s.PathTypeImplementationSpecific):
		return k8s.PathTypeImplementationSpecific
	default:
		return k8s.PathTypePrefix
	}
}

// getGatewayServiceType resolves the gateway service type with the priority
// CR spec > operator config > NodePort default.
func (b *vrsBuilder) getGatewayServiceType(st string) k8s.ServiceType {
	if st == "" {
		st = b.cfg.GatewayServiceType
	}
	switch st {
	case string(k8s.ServiceTypeClusterIP):
		return k8s.ServiceTypeClusterIP
	case string(k8s.ServiceTypeLoadBalancer):
		return k8s.ServiceTypeLoadBalancer
	default:
		return k8s.ServiceTypeNodePort
	}
}

func (b *vrsBuilder) buildDiscoverer() discoverer.Discoverer {
	input := b.cr.Spec.VectorEngine.Vald.Discoverer

	d := discoverer.Discoverer{
		Logging: b.buildLogging(input.LogLevel),
		ClusterRole: discoverer.ClusterRole{
			Name: b.cr.Namespace,
		},
		ClusterRoleBinding: discoverer.ClusterRoleBinding{
			Name: b.cr.Namespace,
		},
		Kind: common.KindType(input.Kind),
	}
	return d
}

func (b *vrsBuilder) buildManager() manager.Manager {
	indexer := b.cr.Spec.VectorEngine.Vald.Indexer

	m := manager.Manager{
		Index: manager.Index{
			Logging: b.buildLogging(indexer.LogLevel),
			Enabled: indexer.Manager,
		},
	}

	if m.Index.Enabled {
		m.Index.Indexer = manager.Indexer{
			AutoIndexDurationLimit:     b.cfg.IndexerAutoIndexDurationLimit,
			AutoSaveIndexDurationLimit: b.cfg.IndexerAutoSaveIndexDurationLimit,
			AutoIndexCheckDuration:     indexer.IndexDuration,
			AutoSaveIndexWaitDuration:  indexer.SaveDuration,
			Concurrency:                &indexer.Concurrency,
		}
		return m
	}

	m.Index.Creator = &manager.Creator{
		Enabled:     true,
		Schedule:    indexer.IndexSchedule,
		Suspend:     indexer.IndexSuspend,
		Concurrency: &indexer.Concurrency,
	}
	m.Index.Saver = &manager.Saver{
		Enabled:  true,
		Schedule: indexer.SaveSchedule,
		Suspend:  indexer.SaveSuspend,
	}
	return m
}

func (b *vrsBuilder) makeName(str ...string) (string, error) {
	name := strings.Join(str, "-")
	if len(name) > maxK8sNameLength {
		return "", errors.ErrNameIsTooLong(name)
	}
	return name, nil
}

func (b *vrsBuilder) buildLabels(iKey int) map[string]string {
	return map[string]string{
		b.labelKey(nodePoolLabelType): string(b.cr.Spec.Infrastructure[iKey].Type),
		b.labelKey(nodePoolLabelRole): string(b.cr.Spec.Infrastructure[iKey].Role),
	}
}

// mergeLabels returns a new map with overlay entries applied on top of base.
// Neither input is mutated, so callers can keep reusing the base map.
func mergeLabels(base, overlay map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(overlay))
	maps.Copy(merged, base)
	maps.Copy(merged, overlay)
	return merged
}

func (b *vrsBuilder) buildLogging(ll string) *defaults.Logging {
	vald := b.cr.Spec.VectorEngine.Vald
	l := ll
	if l == "" {
		if vald.Defaults.LogLevel != "" {
			l = vald.Defaults.LogLevel
		} else {
			l = b.cfg.VrsLogLevel
		}
	}
	return &defaults.Logging{
		Level:  l,
		Format: b.cfg.VrsLogFormat,
		Logger: b.cfg.VrsLogger,
	}
}

func (b *vrsBuilder) labelKey(suffix string) string {
	return labelKey(b.cfg.NodePoolLabelPrefix, suffix)
}

// mergeOverlay layers the default VRS template, the built row and the
// CR-supplied overlay patch on top of each other and decodes the merged
// result back into a ValdRelease.
func (b *vrsBuilder) mergeOverlay(row k8s.Object) (*valdrelease.ValdRelease, error) {
	current, err := resource.ToUnstructured(row)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert row to unstructured")
	}

	// Copy the parsed default VRS before mutating; the cached version must remain intact.
	baseVrs := b.cfg.DefaultVrs.Us.DeepCopy()
	baseVrs.SetName(current.GetName())
	baseVrs.SetNamespace(current.GetNamespace())

	patches := []k8s.Unstructured{*current}

	if patch := b.makeOverlayPatch(row); patch != nil {
		patches = append(patches, *patch)
	}

	allPatches := append([]k8s.Unstructured{*baseVrs}, patches...)
	merged, err := kustomize.Merge(allPatches...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to patch unstructured")
	}

	mergedRaw, err := json.Marshal(merged)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal merged unstructured")
	}
	mergedVRS := &valdrelease.ValdRelease{}
	if err := json.Unmarshal(mergedRaw, mergedVRS); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal into valdrelease")
	}
	return mergedVRS, nil
}

// parseOverlay decodes the CR-supplied overlay patch once per Build. An empty
// overlay leaves b.overlay nil, which disables the patch step entirely.
func (b *vrsBuilder) parseOverlay() error {
	patchRaw := b.cr.Spec.VectorEngine.Vald.Overlay.Raw
	if len(patchRaw) == 0 {
		b.overlay = nil
		return nil
	}

	var obj map[string]any
	if err := json.Unmarshal(patchRaw, &obj); err != nil {
		return err
	}
	b.overlay = &k8s.Unstructured{Object: obj}
	b.overlay.SetGroupVersionKind(b.gvk)
	return nil
}

// makeOverlayPatch clones the pre-parsed overlay for the given row. The clone
// keeps the cached b.overlay intact while name/namespace are rewritten per
// cluster.
func (b *vrsBuilder) makeOverlayPatch(row k8s.Object) *k8s.Unstructured {
	if b.overlay == nil {
		return nil
	}
	patch := b.overlay.DeepCopy()
	patch.SetName(row.GetName())
	patch.SetNamespace(row.GetNamespace())
	return patch
}

func (b *vrsBuilder) buildNodeSelector(nt v1.NodePoolType) map[string]string {
	nt = b.effectiveNodePoolType(nt)
	return map[string]string{
		b.labelKey(nodePoolLabelNamespace): b.cr.GetNamespace(),
		b.labelKey(nodePoolLabelType):      string(nt),
	}
}

func (b *vrsBuilder) buildToleration(nt v1.NodePoolType) *[]k8s.Toleration {
	nt = b.effectiveNodePoolType(nt)
	return &[]k8s.Toleration{
		{
			Key:      b.labelKey(nodePoolLabelNamespace),
			Operator: k8s.TolerationOpEqual,
			Value:    b.cr.GetNamespace(),
			Effect:   k8s.TaintEffectNoSchedule,
		},
		{
			Key:      b.labelKey(nodePoolLabelType),
			Operator: k8s.TolerationOpEqual,
			Value:    string(nt),
		},
	}
}

// effectiveNodePoolType returns the pool type to use after applying the
// agent → general fallback. The decision is a pure function of the
// pre-resolved capability — no I/O, no context.
func (b *vrsBuilder) effectiveNodePoolType(nt v1.NodePoolType) v1.NodePoolType {
	if nt == v1.NodePoolTypeValdAgent && !b.capability.HasAgentPool {
		return v1.NodePoolTypeGeneral
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
func (b *vrsBuilder) applyNodeAffinities(v *valdrelease.ValdRelease) {
	agentNS := b.buildNodeSelector(v1.NodePoolTypeValdAgent)
	agentTol := b.buildToleration(v1.NodePoolTypeValdAgent)
	generalNS := b.buildNodeSelector(v1.NodePoolTypeGeneral)
	generalTol := b.buildToleration(v1.NodePoolTypeGeneral)

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

// reflectPersistentVolume configures persistent storage on the agent when the
// CR requests it. CR-supplied StorageClass / AccessMode win; configured
// defaults (b.cfg.DefaultStorageClass / b.cfg.DefaultAccessMode) fill in
// anything the CR omits. The PV size formula is owned by agentPvSize.
func (b *vrsBuilder) reflectPersistentVolume(v *valdrelease.ValdRelease) {
	pv := b.cr.Spec.VectorEngine.Vald.Agent.PersistentVolume
	if pv == nil || !pv.Enabled {
		return
	}

	sc := pv.StorageClass
	if sc == "" {
		sc = b.cfg.DefaultStorageClass
	}
	am := pv.AccessMode
	if am == "" {
		am = b.cfg.DefaultAccessMode
	}
	memoryBytes := v.Spec.Agent.Resources.Requests.Memory().Value()
	size := agentPvSize(memoryBytes, b.cfg.PvBufferRatio, b.cfg.PvMinSizeBytes)
	v.Spec.Agent.SetPvEnable(sc, am, size)
}
