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

package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/util"
)

// ResourceSyncer renders a Builder's output, applies it to the cluster with
// owner-reference + managed-generation labels, and prunes any owned resources
// that the current Build no longer produces. It is the single place the
// controller talks to the cluster's write side for managed resources.
type ResourceSyncer struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// NewResourceSyncer constructs a ResourceSyncer bound to the given client and
// scheme. The scheme is used for SetControllerReference; the client for the
// CreateOrUpdate / List / Delete calls.
func NewResourceSyncer(c client.Client, scheme *runtime.Scheme) *ResourceSyncer {
	return &ResourceSyncer{Client: c, Scheme: scheme}
}

// Sync builds the desired objects via the Builder, applies them with
// owner-reference + managed-generation labels, and prunes any orphaned owned
// resources. The returned map carries one entry per Create/Update/Prune.
func (s *ResourceSyncer) Sync(ctx context.Context, b desired.Builder, owner *controllerv1.ValdOperatorRelease) (desired.OperationResults, error) {
	opes := desired.OperationResults{}
	builtObj, err := b.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build desired resource: %w", err)
	}
	if builtObj == nil {
		opes["no_resources"] = controllerutil.OperationResultNone
		return opes, nil
	}

	generation := strconv.FormatInt(owner.Generation, 10)

	items, err := util.ToObjectSlice(builtObj)
	if err != nil {
		return nil, err
	}
	for _, obj := range items {
		// Capture the key (which reads GVK) before CreateOrUpdate. The
		// client-side Get path inside CreateOrUpdate can strip TypeMeta on
		// some k8s client implementations, so reading GVK after that point
		// may return empty values.
		key := s.makeKey(obj)

		labels := obj.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}
		labels["managed-generation"] = generation
		obj.SetLabels(labels)

		if err := controllerutil.SetControllerReference(owner, obj, s.Scheme); err != nil {
			return opes, fmt.Errorf("failed to set controller reference: %w", err)
		}

		ope, err := ctrl.CreateOrUpdate(ctx, s.Client, obj, func() error { return nil })
		if err != nil {
			return opes, fmt.Errorf("failed to create or update resource: %w", err)
		}
		opes[key] = ope
	}

	prune, err := s.pruneOldResources(ctx, builtObj, opes, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to prune old resources: %w", err)
	}
	for k, v := range prune {
		opes[k] = v
	}
	return opes, nil
}

// pruneOldResources lists all currently-owned resources of the same kind as
// the Builder's output and deletes any not present in the current Sync.
// Listed items can lose their TypeMeta on some client implementations, so the
// item GVK is recovered from the parent List's GVK before keying.
func (s *ResourceSyncer) pruneOldResources(ctx context.Context, created client.ObjectList, results desired.OperationResults, vor *controllerv1.ValdOperatorRelease) (desired.OperationResults, error) {
	opes := desired.OperationResults{}
	// Capture the list GVK before calling List — some client implementations
	// reset TypeMeta on the populated list.
	itemGVK := itemGVKFromList(created.GetObjectKind().GroupVersionKind())
	exists := created.DeepCopyObject().(client.ObjectList)
	if err := s.Client.List(ctx, exists, client.InNamespace(vor.Namespace)); err != nil {
		return results, fmt.Errorf("failed to list existing resources: %w", err)
	}
	items, err := util.ToObjectSlice(exists)
	if err != nil {
		return nil, err
	}
	for _, obj := range items {
		if obj.GetObjectKind().GroupVersionKind().Empty() {
			obj.GetObjectKind().SetGroupVersionKind(itemGVK)
		}
		if !metav1.IsControlledBy(obj, vor) {
			continue
		}
		key := s.makeKey(obj)
		if _, keep := results[key]; !keep {
			if err := s.Client.Delete(ctx, obj); err != nil && !errors.IsNotFound(err) {
				return results, fmt.Errorf("failed to prune %s: %w", key, err)
			}
			opes[key] = "pruned"
		}
	}
	return opes, nil
}

// itemGVKFromList converts a list GVK to the GVK of its individual items by
// stripping the trailing "List" suffix from the Kind (the standard k8s
// convention: ConfigMapList → ConfigMap).
func itemGVKFromList(listGVK schema.GroupVersionKind) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   listGVK.Group,
		Version: listGVK.Version,
		Kind:    strings.TrimSuffix(listGVK.Kind, "List"),
	}
}

// makeKey returns a string that uniquely identifies an object across
// group/version/kind/namespace/name. The group + version are included so
// same-named/same-kind objects in different API groups don't collide.
func (s *ResourceSyncer) makeKey(obj client.Object) string {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return gvk.Group + "/" + gvk.Version + "/" + gvk.Kind + "/" + obj.GetNamespace() + "/" + obj.GetName()
}
