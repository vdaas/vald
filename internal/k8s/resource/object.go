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

	json "github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// This file provides generic helpers for operating on arbitrary runtime
// objects (such as custom resources) through a scheme-aware client, so that
// controllers do not have to depend on controller-runtime helpers directly.

type (
	ObjectAPI       = kclient.Client
	ObjectListType  = kclient.ObjectList
	OperationResult = controllerutil.OperationResult
	ListOption      = kclient.ListOption
)

const (
	OperationResultNone    = controllerutil.OperationResultNone
	OperationResultCreated = controllerutil.OperationResultCreated
	OperationResultUpdated = controllerutil.OperationResultUpdated
)

var (
	IgnoreNotFound    = kclient.IgnoreNotFound
	SemanticDeepEqual = equality.Semantic.DeepEqual
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

func RefreshObject[T Object](ctx context.Context, c ObjectAPI, obj T) (T, error) {
	if err := c.Get(ctx, kclient.ObjectKeyFromObject(obj), obj); err != nil {
		return obj, err
	}
	return obj, nil
}

func ListObjects[L ObjectListType](
	ctx context.Context, c ObjectAPI, list L, opts ...ListOption,
) (L, error) {
	if err := c.List(ctx, list, opts...); err != nil {
		return list, err
	}
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
