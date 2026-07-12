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

package resource

import (
	"context"
	"maps"
	"reflect"
	"strconv"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// defaultGenerationLabel is the fallback label key that records the owner
// generation on managed resources when the caller omits it.
const defaultGenerationLabel = "managed-generation"

// prunedResult marks an object that Sync deleted because the desired set no
// longer contains it.
const prunedResult OperationResult = "pruned"

// SyncResults maps an object identity key (group/version/kind/namespace/name)
// to the operation Sync applied to that object.
type SyncResults map[string]OperationResult

// Syncer applies a desired set of objects to the cluster on behalf of an
// owner object and prunes owned objects that are no longer desired. It is
// owner-type agnostic: any Object (typically a custom resource) can own the
// managed set.
type Syncer struct {
	api    k8s.Client
	scheme *runtime.Scheme
	// generationLabel is the label key that records the owner generation on
	// every managed resource.
	generationLabel string
}

// NewSyncer constructs a Syncer bound to the given client and scheme. The
// scheme is used for controllerutil.SetControllerReference; the client for the
// CreateOrUpdate / Delete calls. generationLabel overrides the label key used
// to record the owner generation; empty falls back to the default.
func NewSyncer(api k8s.Client, scheme *runtime.Scheme, generationLabel string) *Syncer {
	if generationLabel == "" {
		generationLabel = defaultGenerationLabel
	}
	return &Syncer{api: api, scheme: scheme, generationLabel: generationLabel}
}

// Sync applies the desired objects with owner-reference + managed-generation
// labels, then prunes objects returned by existing that are controlled by
// owner but absent from the desired set. existing may be nil to skip pruning.
// The returned map carries one entry per Create/Update/Prune.
func (s *Syncer) Sync(
	ctx context.Context,
	owner k8s.Object,
	desired []k8s.Object,
	existing func(ctx context.Context) ([]k8s.Object, error),
) (SyncResults, error) {
	opes := make(SyncResults, max(len(desired), 1))
	if len(desired) == 0 {
		// keep pruning below: an empty desired set must still delete every
		// owned object that is no longer produced.
		opes["no_resources"] = OperationResultNone
	}

	generation := strconv.FormatInt(owner.GetGeneration(), 10)

	for _, obj := range desired {
		key := syncKey(obj)

		labels := obj.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}
		labels[s.generationLabel] = generation
		obj.SetLabels(labels)

		if err := controllerutil.SetControllerReference(owner, obj, s.scheme); err != nil {
			return opes, errors.Wrap(err, "failed to set controller reference")
		}

		// ctrl.CreateOrUpdate overwrites obj with the live cluster state
		// before calling mutate, so snapshot the desired state first and
		// re-apply it inside mutate; otherwise the desired spec is discarded
		// and existing objects never converge.
		want, ok := obj.DeepCopyObject().(k8s.Object)
		if !ok {
			return opes, errors.Errorf("failed to snapshot desired state of %s", key)
		}
		ope, err := ctrl.CreateOrUpdate(ctx, s.api, obj, func() error {
			return applyDesiredState(obj, want)
		})
		if err != nil {
			return opes, errors.Wrap(err, "failed to create or update resource")
		}
		opes[key] = ope
	}

	if existing == nil {
		return opes, nil
	}
	prune, err := s.prune(ctx, owner, opes, existing)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prune old resources")
	}
	maps.Copy(opes, prune)
	return opes, nil
}

// applyDesiredState re-applies the desired state captured before
// ctrl.CreateOrUpdate's Get onto the freshly fetched obj: every non-metadata
// payload entry (spec, data, ...) plus labels, annotations and owner
// references. Live identity (name, namespace, uid, resourceVersion,
// creationTimestamp) and status are preserved so the update targets the
// current revision.
func applyDesiredState(obj, desired k8s.Object) error {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		d, ok := desired.(*unstructured.Unstructured)
		if !ok {
			return errors.Errorf("desired object type %T does not match fetched object type %T", desired, obj)
		}
		if o.Object == nil {
			o.Object = make(map[string]any, len(d.Object))
		}
		for k := range o.Object {
			if k == "metadata" || k == "status" {
				continue
			}
			if _, ok := d.Object[k]; !ok {
				delete(o.Object, k)
			}
		}
		for k, v := range d.Object {
			switch k {
			case "metadata", "status":
			default:
				o.Object[k] = v
			}
		}
	} else {
		ov, dv := reflect.ValueOf(obj), reflect.ValueOf(desired)
		if ov.Kind() != reflect.Pointer || dv.Kind() != reflect.Pointer ||
			ov.IsNil() || dv.IsNil() || ov.Type() != dv.Type() {
			return errors.Errorf("cannot apply desired state of %T onto %T", desired, obj)
		}
		oe, de := ov.Elem(), dv.Elem()
		if oe.Kind() != reflect.Struct {
			return errors.Errorf("unsupported object kind %s for %T", oe.Kind(), obj)
		}
		for i := range oe.NumField() {
			f := oe.Type().Field(i)
			if !f.IsExported() {
				continue
			}
			switch f.Name {
			case "TypeMeta", "ObjectMeta", "Status":
				continue
			}
			oe.Field(i).Set(de.Field(i))
		}
	}
	obj.SetLabels(desired.GetLabels())
	obj.SetAnnotations(desired.GetAnnotations())
	obj.SetOwnerReferences(desired.GetOwnerReferences())
	return nil
}

// prune deletes every object from existing that is controlled by owner but
// has no entry in applied (i.e. the current desired set no longer produces
// it).
func (s *Syncer) prune(
	ctx context.Context,
	owner k8s.Object,
	applied SyncResults,
	existing func(ctx context.Context) ([]k8s.Object, error),
) (SyncResults, error) {
	opes := SyncResults{}

	items, err := existing(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list existing resources")
	}

	for _, obj := range items {
		if !metav1.IsControlledBy(obj, owner) {
			continue
		}
		key := syncKey(obj)
		if _, keep := applied[key]; !keep {
			if err := s.api.Delete(ctx, obj); err != nil && !apierrors.IsNotFound(err) {
				return nil, errors.Wrapf(err, "failed to prune %s", key)
			}
			opes[key] = prunedResult
		}
	}
	return opes, nil
}

// syncKey returns a string that uniquely identifies an object across
// group/version/kind/namespace/name. The group + version are included so
// same-named/same-kind objects in different API groups don't collide.
func syncKey(obj k8s.Object) string {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return gvk.Group + "/" + gvk.Version + "/" + gvk.Kind + "/" + obj.GetNamespace() + "/" + obj.GetName()
}
