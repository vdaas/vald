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
	"context"
	"maps"

	json "github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/metadata"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease"
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
	capability nodePoolCapability
}

func newVrsBuilder(
	cr *v1.ValdOperatorRelease, cfg *config.Config, capability nodePoolCapability,
) *vrsBuilder {
	return &vrsBuilder{
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

		// The row is constructed locally (not fetched), so its GVK is set here
		// as part of construction: it flows through mergeOverlay's unstructured
		// round-trip into every built item and is what Syncer keys on.
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

		// Resolve the agent's NodePool spec, applying the general-pool
		// fallback rule.
		agentPool := resolveAgentNodePool(infra)
		row.SetScaledResources(agentPool.NodeCount, agentPool.MachineResource, resourceParams)

		// Reflect optional settings (PV, node affinities) after resources are confirmed.
		b.reflectPersistentVolume(row)
		b.applyNodeAffinities(row)

		// Merge the per-infra labels on top of the existing ones so the
		// managed-by / managed-resource labels set above survive. The labels
		// are constant within an infra, so this happens once before the
		// unstructured conversion below.
		row.SetLabels(mergeLabels(row.GetLabels(), b.buildLabels(i)))

		// Convert the fully-built row once per infra: between the clusters of
		// one infra the row differs only by metadata.name, which is stamped
		// onto the unstructured form inside the loop. The conversion is a JSON
		// round trip of the whole tree and dominates the per-cluster cost when
		// repeated (measured ~49% of mergeOverlay).
		current, err := resource.ToUnstructured(row)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert row to unstructured")
		}

		for _, cluster := range infra.Clusters {
			name, err := b.makeName(b.cr.GetNamespace(), cluster.Name)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to make name for cluster %s", cluster.Name)
			}
			// Renaming current in place between iterations is safe: mergeOverlay
			// materializes its merged map into an independent typed object
			// before returning, so nothing built for a previous cluster keeps a
			// reference into current's maps.
			current.SetName(name)
			u, err := b.mergeOverlay(current)
			if err != nil {
				return nil, errors.Wrap(err, "failed to merge overlay")
			}

			b.list.Items = append(b.list.Items, *u)
		}
	}

	return resource.ObjectsOf(b.list.Items), nil
}

// fetchExistingVrs lists the ValdRelease objects currently present in the
// namespace so the syncer can prune the ones the current build no longer
// produces. The GVK the client strips on decode is restored by
// resource.ListObjects from the client scheme, with a static fallback to
// valdrelease.GVK so every returned item is unconditionally identifiable.
func fetchExistingVrs(ctx context.Context, c k8s.Client, namespace string) ([]k8s.Object, error) {
	exists := &valdrelease.ValdReleaseList{}
	if _, err := resource.ListObjects(ctx, c, exists, k8s.InNamespace(namespace)); err != nil {
		return nil, errors.Wrap(err, "failed to list existing ValdRelease objects")
	}
	// The scheme-based restoration above is best-effort: it silently skips
	// when the scheme cannot resolve the type (unregistered, or the same Go
	// type registered under multiple GVKs — the typical path when the CRD
	// gains a new API version). An item left with an empty GVK would give the
	// syncer a "///ns/name" prune key that never matches the desired set, and
	// every existing ValdRelease would be deleted. Stamping the static GVK
	// here restores the old unconditional guarantee against that deletion
	// accident; the empty check makes it a no-op whenever the automatic
	// restoration already succeeded.
	for i := range exists.Items {
		if exists.Items[i].GroupVersionKind().Empty() {
			exists.Items[i].SetGroupVersionKind(valdrelease.GVK)
		}
	}
	return resource.ObjectsOf(exists.Items), nil
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

func (b *vrsBuilder) buildDefaults() *valdrelease.Defaults {
	vald := b.cr.Spec.VectorEngine.Vald
	return &valdrelease.Defaults{
		Logging: b.buildLogging(vald.Defaults.LogLevel),
	}
}

func (b *vrsBuilder) buildAgent() *valdrelease.Agent {
	input := b.cr.Spec.VectorEngine.Vald.Agent
	return &valdrelease.Agent{
		Logging: b.buildLogging(input.LogLevel),
		Kind:    ptr(valdrelease.AgentKindStatefulSet),
		RollingUpdate: &valdrelease.AgentRollingUpdate{
			MaxUnavailable: new(b.cfg.AgentMaxUnavailable),
			MaxSurge:       new(b.cfg.AgentMaxSurge),
		},
		Ngt: b.buildAgentNgt(),
	}
}

func (b *vrsBuilder) buildAgentNgt() *valdrelease.AgentNgt {
	input := b.cr.Spec.VectorEngine.Vald.Agent.Ngt
	return &valdrelease.AgentNgt{
		Dimension:          new(input.Dimension),
		DistanceType:       new(valdrelease.AgentNgtDistanceType(input.DistanceType)),
		ObjectType:         new(valdrelease.AgentNgtObjectType(input.ObjectType)),
		SearchEdgeSize:     new(input.SearchEdgeSize),
		CreationEdgeSize:   new(input.CreationEdgeSize),
		EnableInMemoryMode: new(b.cfg.AgentEnableInMemoryMode),
	}
}

func (b *vrsBuilder) buildGateway() *valdrelease.Gateway {
	return &valdrelease.Gateway{
		Lb: b.buildLb(),
	}
}

func (b *vrsBuilder) buildLb() *valdrelease.GatewayLb {
	inputGw := b.cr.Spec.VectorEngine.Vald.Gateway

	return &valdrelease.GatewayLb{
		Logging: b.buildLogging(inputGw.LogLevel),
		Hpa: &valdrelease.Hpa{
			TargetCPUUtilizationPercentage: new(b.cfg.GatewayHpaTargetCPUUtilization),
		},
		GatewayConfig: &valdrelease.GatewayLbGatewayConfig{
			IndexReplica: new(inputGw.IndexReplica),
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
func (b *vrsBuilder) buildIngress(in *v1.GatewayIngress) *valdrelease.GatewayLbIngress {
	base := &valdrelease.GatewayLbIngress{
		DefaultBackend: &valdrelease.GatewayLbIngressDefaultBackend{Enabled: new(false)},
		PathType:       new(b.getGatewayIngressPathType()),
		ServicePort:    new(b.cfg.GatewayIngressServicePort),
	}
	if len(b.cfg.GatewayIngressAnnotations) > 0 {
		ann := toAnyMap(mergeLabels(nil, b.cfg.GatewayIngressAnnotations))
		base.Annotations = &ann
	}
	if !b.cfg.EnableIngress || in == nil || !in.Enabled {
		return base
	}
	base.Enabled = new(true)
	base.Host = new(in.Host)
	return base
}

func (b *vrsBuilder) getGatewayIngressPathType() string {
	switch b.cfg.GatewayIngressPathType {
	case string(k8s.PathTypeExact):
		return string(k8s.PathTypeExact)
	case string(k8s.PathTypeImplementationSpecific):
		return string(k8s.PathTypeImplementationSpecific)
	default:
		return string(k8s.PathTypePrefix)
	}
}

// getGatewayServiceType resolves the gateway service type with the priority
// CR spec > operator config > NodePort default.
func (b *vrsBuilder) getGatewayServiceType(st string) *valdrelease.GatewayLbServiceType {
	if st == "" {
		st = b.cfg.GatewayServiceType
	}
	switch st {
	case string(k8s.ServiceTypeClusterIP):
		return ptr(valdrelease.GatewayLbServiceTypeClusterIP)
	case string(k8s.ServiceTypeLoadBalancer):
		return ptr(valdrelease.GatewayLbServiceTypeLoadBalancer)
	default:
		return ptr(valdrelease.GatewayLbServiceTypeNodePort)
	}
}

func (b *vrsBuilder) buildDiscoverer() *valdrelease.Discoverer {
	input := b.cr.Spec.VectorEngine.Vald.Discoverer

	return &valdrelease.Discoverer{
		Logging: b.buildLogging(input.LogLevel),
		ClusterRole: &valdrelease.DiscovererClusterRole{
			Name: new(b.cr.Namespace),
		},
		ClusterRoleBinding: &valdrelease.DiscovererClusterRoleBinding{
			Name: new(b.cr.Namespace),
		},
		Kind: new(valdrelease.DiscovererKind(input.Kind)),
	}
}

func (b *vrsBuilder) buildManager() *valdrelease.Manager {
	indexer := b.cr.Spec.VectorEngine.Vald.Indexer

	m := &valdrelease.Manager{
		Index: &valdrelease.ManagerIndex{
			Logging: b.buildLogging(indexer.LogLevel),
			Enabled: new(indexer.Manager),
		},
	}

	if indexer.Manager {
		m.Index.Indexer = &valdrelease.ManagerIndexIndexer{
			AutoIndexDurationLimit:     new(b.cfg.IndexerAutoIndexDurationLimit),
			AutoSaveIndexDurationLimit: new(b.cfg.IndexerAutoSaveIndexDurationLimit),
			AutoIndexCheckDuration:     new(indexer.IndexDuration),
			AutoSaveIndexWaitDuration:  new(indexer.SaveDuration),
			Concurrency:                new(indexer.Concurrency),
		}
		return m
	}

	m.Index.Creator = &valdrelease.ManagerIndexCreator{
		Enabled:     new(true),
		Schedule:    new(indexer.IndexSchedule),
		Suspend:     new(indexer.IndexSuspend),
		Concurrency: new(indexer.Concurrency),
	}
	m.Index.Saver = &valdrelease.ManagerIndexSaver{
		Enabled:  new(true),
		Schedule: new(indexer.SaveSchedule),
		Suspend:  new(indexer.SaveSuspend),
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

func (b *vrsBuilder) buildLogging(ll string) *valdrelease.Logging {
	vald := b.cr.Spec.VectorEngine.Vald
	l := ll
	if l == "" {
		if vald.Defaults.LogLevel != "" {
			l = vald.Defaults.LogLevel
		} else {
			l = b.cfg.VrsLogLevel
		}
	}
	return &valdrelease.Logging{
		Level:  l,
		Format: b.cfg.VrsLogFormat,
		Logger: b.cfg.VrsLogger,
	}
}

// ptr returns a pointer to v. The generated ValdRelease types use pointers for
// nearly every field, so the builder wraps its scalar inputs with this.
func ptr[T any](v T) *T { return new(v) }

// toAnyMap converts a string map into the free-form map[string]interface{}
// shape the generated schema uses for annotation-like fields.
func toAnyMap(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// nsPtr clones m into an independent NodeSelector and returns its address, so
// multiple components can share the same source map without aliasing.
func nsPtr(m map[string]string) *valdrelease.NodeSelector {
	ns := maps.Clone(m)
	return &ns
}

func (b *vrsBuilder) labelKey(suffix string) string {
	return labelKey(b.cfg.NodePoolLabelPrefix, suffix)
}

// mergeOverlay layers the default VRS template, the built row (already
// converted to its unstructured form by Build) and the CR-supplied overlay
// patch on top of each other and decodes the merged result back into a
// ValdRelease. current is treated as read-only.
//
// Merge semantics: each successive layer wins over the previous one for every
// key that is present, including bool false, numeric zero, and empty strings.
// Merging happens at the unstructured (map[string]any) level to avoid the
// limitations of reflection-based mergers with types that have unexported
// fields (resource.Quantity) or zero-value booleans (*bool → false).
func (b *vrsBuilder) mergeOverlay(current *k8s.Unstructured) (*valdrelease.ValdRelease, error) {
	// Merge at the unstructured map level: base → current → overlay.
	// Each subsequent layer wins for every key it provides.
	//
	// The cached default VRS map is passed directly as the merge base:
	// deepMergeMap never mutates its inputs (it writes only into clones), so
	// the previous defensive DeepCopy of the whole template tree is not needed.
	// TestVrsBuilder_Build_DoesNotMutateDefaultVrsCache pins this invariant.
	merged := deepMergeMap(b.cfg.DefaultVrs.Us.Object, current.Object)
	if patch := b.makeOverlayPatch(current); patch != nil {
		merged = deepMergeMap(merged, patch.Object)
	}

	// Stamping name/namespace on the merged result is equivalent to the old
	// behaviour of stamping the base copy: `current` always carries the row's
	// metadata.name (and namespace when non-empty), so those keys already won
	// the merge; Set* additionally removes a template-supplied value when the
	// row's is empty, exactly as SetName/SetNamespace on the base copy did.
	// Writing into merged["metadata"] cannot touch the cached template: the
	// metadata map here is either a fresh clone made by deepMergeMap (both
	// layers carry metadata) or current's own map (template has none) — and in
	// the latter case Set* re-writes the very values current already carries,
	// leaving even that shared map unchanged.
	mergedU := &k8s.Unstructured{Object: merged}
	mergedU.SetName(current.GetName())
	mergedU.SetNamespace(current.GetNamespace())

	raw, err := json.Marshal(merged)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal merged ValdRelease")
	}
	var vr valdrelease.ValdRelease
	if err := json.Unmarshal(raw, &vr); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal merged ValdRelease")
	}
	return &vr, nil
}

// deepMergeMap returns a new map that is the result of deeply merging src into
// dst. src wins for every key it provides; nested maps are merged recursively.
// All other value types (including bool false and numeric zero) are taken from
// src verbatim when present.
func deepMergeMap(dst, src map[string]any) map[string]any {
	out := maps.Clone(dst)
	for k, sv := range src {
		if sm, ok := sv.(map[string]any); ok {
			if dm, ok := out[k].(map[string]any); ok {
				out[k] = deepMergeMap(dm, sm)
				continue
			}
		}
		out[k] = sv
	}
	return out
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
	// The overlay patch is constructed locally from CR-supplied JSON (not
	// fetched), so its GVK is pinned here: the patch is the last merge layer,
	// which keeps a stray apiVersion/kind in the raw overlay from leaking
	// into the built items.
	b.overlay = &k8s.Unstructured{Object: obj}
	b.overlay.SetGroupVersionKind(valdrelease.GVK)
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

	v.Spec.Agent.NodeSelector = nsPtr(agentNS)
	v.Spec.Agent.Tolerations = agentTol

	v.Spec.Gateway.Lb.NodeSelector = nsPtr(generalNS)
	v.Spec.Gateway.Lb.Tolerations = generalTol

	v.Spec.Discoverer.NodeSelector = nsPtr(generalNS)
	v.Spec.Discoverer.Tolerations = generalTol

	v.Spec.Manager.Index.NodeSelector = nsPtr(generalNS)
	v.Spec.Manager.Index.Tolerations = generalTol

	if v.Spec.Manager.Index.Saver != nil {
		v.Spec.Manager.Index.Saver.NodeSelector = nsPtr(generalNS)
		v.Spec.Manager.Index.Saver.Tolerations = generalTol
	}
	if v.Spec.Manager.Index.Creator != nil {
		v.Spec.Manager.Index.Creator.NodeSelector = nsPtr(generalNS)
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
