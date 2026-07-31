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

package resource

import (
	"context"

	json "github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"k8s.io/apimachinery/pkg/api/equality"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// This file provides generic helpers for operating on arbitrary runtime
// objects (such as custom resources) through a scheme-aware client, so that
// controllers do not have to depend on controller-runtime helpers directly.
//
// Object access goes through k8s.Client (the same controller-runtime
// client.Client alias reconcilers get from mgr.GetClient()), rather than a
// separately named type, so there is exactly one "scheme-aware client"
// abstraction across internal/k8s.

type (
	ObjectListType  = kclient.ObjectList
	OperationResult = controllerutil.OperationResult
	ListOption      = kclient.ListOption
)

const (
	OperationResultNone    = controllerutil.OperationResultNone
	OperationResultCreated = controllerutil.OperationResultCreated
	OperationResultUpdated = controllerutil.OperationResultUpdated
)

// IgnoreNotFound returns nil on a NotFound error and err otherwise.
// Delegates to sigs.k8s.io/controller-runtime/pkg/client.IgnoreNotFound.
func IgnoreNotFound(err error) error {
	return kclient.IgnoreNotFound(err)
}

// SemanticDeepEqual reports whether a1 and a2 are semantically equal
// (treating an empty slice/map as equal to nil, unlike reflect.DeepEqual).
// Delegates to k8s.io/apimachinery/pkg/api/equality.Semantic.DeepEqual.
func SemanticDeepEqual(a1, a2 any) bool {
	return equality.Semantic.DeepEqual(a1, a2)
}

// listKindSuffix is the conventional suffix Kubernetes list kinds append to
// their item kind (e.g. ValdReleaseList -> ValdRelease).
const listKindSuffix = "List"

// restoreObjectGVK repopulates obj's TypeMeta from the scheme when the client
// stripped it (controller-runtime's Get/List clear TypeMeta while decoding).
// A TypeMeta that is already populated is never overwritten, and restoration
// is best-effort: types the scheme cannot resolve (unregistered or ambiguous)
// keep their empty TypeMeta rather than failing the successful fetch.
func restoreObjectGVK(scheme *runtime.Scheme, obj runtime.Object) {
	if scheme == nil || obj == nil || !obj.GetObjectKind().GroupVersionKind().Empty() {
		return
	}
	gvk, err := apiutil.GVKForObject(obj, scheme)
	if err != nil {
		// Deliberate skip: an unresolved GVK just leaves the object exactly
		// as the client returned it (the pre-restoration behavior). Debug
		// level keeps the reconcile hot path quiet while still recording the
		// skip instead of silently discarding the error.
		log.Debugf("skipped GVK restoration for %T: %v", obj, err)
		return
	}
	obj.GetObjectKind().SetGroupVersionKind(gvk)
}

// restoreListGVK restores the list's own GVK from the scheme and then stamps
// the item GVK (the list kind minus the "List" suffix) onto every item whose
// TypeMeta is empty. Lists whose kind cannot be resolved or does not follow
// the "<Kind>List" convention leave their items untouched.
func restoreListGVK(scheme *runtime.Scheme, list ObjectListType) {
	if list == nil {
		return
	}
	if scheme != nil && list.GetObjectKind().GroupVersionKind().Empty() {
		lgvk, err := apiutil.GVKForObject(list, scheme)
		if err != nil {
			// An unresolvable list GVK (unregistered type, or the same Go
			// type registered under multiple GVKs) leaves every item with an
			// empty TypeMeta, which identity-keyed consumers (e.g. Syncer
			// prune keys) cannot match — warn so the root cause is visible.
			log.Warnf("failed to restore GVK of list %T from the scheme, item GVKs stay empty: %v", list, err)
			return
		}
		list.GetObjectKind().SetGroupVersionKind(lgvk)
	}
	gvk := list.GetObjectKind().GroupVersionKind()
	if gvk.Empty() || !strings.HasSuffix(gvk.Kind, listKindSuffix) || gvk.Kind == listKindSuffix {
		return
	}
	gvk.Kind = strings.TrimSuffix(gvk.Kind, listKindSuffix)
	items, err := apimeta.ExtractList(list)
	if err != nil {
		// Deliberate skip: a list without an extractable Items field keeps
		// the pre-restoration behavior for its items. Debug level: this is a
		// static property of the list type, not an actionable runtime fault.
		log.Debugf("skipped item GVK restoration for %T: %v", list, err)
		return
	}
	for _, item := range items {
		if item != nil && item.GetObjectKind().GroupVersionKind().Empty() {
			item.GetObjectKind().SetGroupVersionKind(gvk)
		}
	}
}

// GetObject fetches the object identified by name and namespace into obj and
// returns it with its GVK restored from the client scheme. Pass an empty
// namespace for cluster-scoped objects.
func GetObject[T Object](
	ctx context.Context, c k8s.Client, name, namespace string, obj T,
) (T, error) {
	if c == nil {
		return obj, errors.New("k8s client is nil")
	}
	if err := c.Get(ctx, kclient.ObjectKey{Name: name, Namespace: namespace}, obj); err != nil {
		return obj, err
	}
	restoreObjectGVK(c.Scheme(), obj)
	return obj, nil
}

// RefreshObject re-fetches obj from the cluster in place and returns it with
// its GVK restored from the client scheme.
func RefreshObject[T Object](ctx context.Context, c k8s.Client, obj T) (T, error) {
	if c == nil {
		return obj, errors.New("k8s client is nil")
	}
	if err := c.Get(ctx, kclient.ObjectKeyFromObject(obj), obj); err != nil {
		return obj, err
	}
	restoreObjectGVK(c.Scheme(), obj)
	return obj, nil
}

// ListObjects fetches every object matching opts into list and returns it
// with the list GVK and each item's GVK restored from the client scheme.
func ListObjects[L ObjectListType](
	ctx context.Context, c k8s.Client, list L, opts ...ListOption,
) (L, error) {
	if c == nil {
		return list, errors.New("k8s client is nil")
	}
	if err := c.List(ctx, list, opts...); err != nil {
		return list, err
	}
	restoreListGVK(c.Scheme(), list)
	return list, nil
}

// ToUnstructured converts obj into its unstructured representation through a
// JSON round-trip, so the result matches the object's JSON wire format
// exactly (including omitempty handling).
func ToUnstructured(obj Object) (*unstructured.Unstructured, error) {
	raw, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal object")
	}
	us := new(unstructured.Unstructured)
	if err := json.Unmarshal(raw, us); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal object into unstructured")
	}
	return us, nil
}

// ObjectsOf returns the addresses of items as an []Object slice, so typed
// list items can be handed to Object-based helpers without per-caller loops.
func ObjectsOf[T any, PT Objectable[T]](items []T) []Object {
	out := make([]Object, len(items))
	for i := range items {
		out[i] = PT(&items[i])
	}
	return out
}
