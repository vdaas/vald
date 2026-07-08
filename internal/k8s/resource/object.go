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

	"github.com/vdaas/vald/internal/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	ctrl "sigs.k8s.io/controller-runtime"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// This file provides generic helpers for operating on arbitrary runtime
// objects (such as custom resources) through a scheme-aware client, so that
// controllers do not have to depend on controller-runtime helpers directly.

type (
	// ObjectAPI is the scheme-aware Kubernetes client used to operate on
	// arbitrary runtime objects including custom resources.
	ObjectAPI = kclient.Client
	// ObjectWatchAPI is an ObjectAPI that can additionally watch objects.
	ObjectWatchAPI = kclient.WithWatch
	// ObjectListType is the list counterpart of Object handled by ObjectAPI.
	ObjectListType = kclient.ObjectList
	// OperationResult describes the mutation applied by CreateOrUpdateObject.
	OperationResult = controllerutil.OperationResult
	// ListOption configures ListObjects calls.
	ListOption = kclient.ListOption
)

const (
	// OperationResultNone means that the resource has not been changed.
	OperationResultNone = controllerutil.OperationResultNone
	// OperationResultCreated means that a new resource is created.
	OperationResultCreated = controllerutil.OperationResultCreated
	// OperationResultUpdated means that an existing resource is updated.
	OperationResultUpdated = controllerutil.OperationResultUpdated
)

// GetObject fetches the object identified by name and namespace into obj and
// returns it. Pass an empty namespace for cluster-scoped objects.
func GetObject[T Object](
	ctx context.Context, c ObjectAPI, name, namespace string, obj T,
) (T, error) {
	if err := c.Get(ctx, kclient.ObjectKey{Name: name, Namespace: namespace}, obj); err != nil {
		return obj, err
	}
	return obj, nil
}

// RefreshObject re-fetches obj from the cluster using its own name and
// namespace as the key.
func RefreshObject[T Object](ctx context.Context, c ObjectAPI, obj T) (T, error) {
	if err := c.Get(ctx, kclient.ObjectKeyFromObject(obj), obj); err != nil {
		return obj, err
	}
	return obj, nil
}

// ListObjects fills list with the objects matching opts and returns it.
func ListObjects[L ObjectListType](
	ctx context.Context, c ObjectAPI, list L, opts ...ListOption,
) (L, error) {
	if err := c.List(ctx, list, opts...); err != nil {
		return list, err
	}
	return list, nil
}

// UpdateObjectStatus persists the status subresource of obj.
func UpdateObjectStatus[T Object](ctx context.Context, c ObjectAPI, obj T) error {
	return c.Status().Update(ctx, obj)
}

// DeleteObject removes obj from the cluster.
func DeleteObject[T Object](ctx context.Context, c ObjectAPI, obj T) error {
	return c.Delete(ctx, obj)
}

// SetControllerReference sets owner as the controlling OwnerReference of obj.
func SetControllerReference(owner, obj Object, scheme *runtime.Scheme) error {
	return controllerutil.SetControllerReference(owner, obj, scheme)
}

// CreateOrUpdateObject creates obj when it does not exist, otherwise applies
// mutate and updates it, returning the operation that was performed.
func CreateOrUpdateObject[T Object](
	ctx context.Context, c ObjectAPI, obj T, mutate func() error,
) (OperationResult, error) {
	return ctrl.CreateOrUpdate(ctx, c, obj, mutate)
}

// IgnoreNotFound returns nil when err is a Kubernetes NotFound error,
// otherwise it returns err unchanged.
func IgnoreNotFound(err error) error {
	return kclient.IgnoreNotFound(err)
}

// WatchObjects starts a watch for the list kind L matching opts and returns
// the watch interface. The caller owns the returned watcher and must Stop it.
func WatchObjects[L ObjectListType](
	ctx context.Context, c ObjectWatchAPI, list L, opts ...ListOption,
) (watch.Interface, error) {
	return c.Watch(ctx, list, opts...)
}

// LabelSelector builds a labels.Selector from a single requirement, e.g.
// LabelSelector("app", selection.Equals, []string{"vald-agent"}).
func LabelSelector(key string, op selection.Operator, vals []string) (labels.Selector, error) {
	requirement, err := labels.NewRequirement(key, op, vals)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create requirement on creating label selector")
	}
	return labels.NewSelector().Add(*requirement), nil
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
