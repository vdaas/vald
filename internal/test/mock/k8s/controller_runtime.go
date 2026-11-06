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
package k8s

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type SubResourceWriterMock struct {
	client.SubResourceWriter

	UpdateFunc func(context.Context, client.Object, ...client.SubResourceUpdateOption) error
}

func (sm *SubResourceWriterMock) Update(
	ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption,
) error {
	return sm.UpdateFunc(ctx, obj, opts...)
}

type ClientMock struct {
	client.Client

	StatusFunc      func() client.SubResourceWriter
	GetFunc         func(context.Context, client.ObjectKey, client.Object, ...client.GetOption) error
	CreateFunc      func(context.Context, client.Object, ...client.CreateOption) error
	DeleteFunc      func(context.Context, client.Object, ...client.DeleteOption) error
	DeleteAllOfFunc func(context.Context, client.Object, ...client.DeleteAllOfOption) error
	ListFunc        func(context.Context, client.ObjectList, ...client.ListOption) error
	SchemeFunc      func() *runtime.Scheme
}

// List calls ListFunc, or leaves list empty and succeeds when ListFunc is unset.
func (cm *ClientMock) List(
	ctx context.Context, list client.ObjectList, opts ...client.ListOption,
) error {
	if cm.ListFunc == nil {
		return nil
	}
	return cm.ListFunc(ctx, list, opts...)
}

func (cm *ClientMock) Status() client.SubResourceWriter {
	return cm.StatusFunc()
}

// Scheme returns the configured scheme, or nil when SchemeFunc is unset. A
// nil scheme makes the resource helpers skip GVK restoration instead of
// panicking through the embedded nil client.Client.
func (cm *ClientMock) Scheme() *runtime.Scheme {
	if cm.SchemeFunc == nil {
		return nil
	}
	return cm.SchemeFunc()
}

func (cm *ClientMock) Get(
	ctx context.Context, objKey client.ObjectKey, obj client.Object, opts ...client.GetOption,
) error {
	return cm.GetFunc(ctx, objKey, obj, opts...)
}

func (cm *ClientMock) Create(
	ctx context.Context, obj client.Object, opts ...client.CreateOption,
) error {
	return cm.CreateFunc(ctx, obj, opts...)
}

func (cm *ClientMock) Delete(
	ctx context.Context, obj client.Object, opts ...client.DeleteOption,
) error {
	return cm.DeleteFunc(ctx, obj, opts...)
}

func (cm *ClientMock) DeleteAllOf(
	ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption,
) error {
	return cm.DeleteAllOfFunc(ctx, obj, opts...)
}

type ManagerMock struct {
	manager.Manager

	GetClientFunc func() client.Client
}

func (mm *ManagerMock) GetClient() client.Client {
	return mm.GetClientFunc()
}

// GetConfig returns an empty, non-nil *rest.Config so client.NewFromManager
// can build a clientset from this mock without dialing any cluster.
func (*ManagerMock) GetConfig() *rest.Config {
	return &rest.Config{}
}

// listManagerMock adds the scheme and field indexer that a list reconciler's
// NewReconciler registers into, on top of ManagerMock's client and config.
type listManagerMock struct {
	ManagerMock

	scheme *runtime.Scheme
}

func (m *listManagerMock) GetScheme() *runtime.Scheme {
	return m.scheme
}

func (*listManagerMock) GetFieldIndexer() client.FieldIndexer {
	return noopFieldIndexer{}
}

// noopFieldIndexer accepts every index registration without doing anything.
type noopFieldIndexer struct{}

func (noopFieldIndexer) IndexField(
	context.Context, client.Object, string, client.IndexerFunc,
) error {
	return nil
}

// NewListManagerMock returns a manager mock serving c as its client, with a
// fresh scheme and a no-op field indexer, so that tests outside internal/ can
// run a list reconciler's NewReconciler and Reconcile against a fake client.
func NewListManagerMock(c client.Client) manager.Manager {
	return &listManagerMock{
		GetClientFunc: func() client.Client {
			return c
		},
		scheme: runtime.NewScheme(),
	}
}

// NewDefaultManagerMock returns default empty mock object.
func NewDefaultManagerMock() manager.Manager {
	return &ManagerMock{
		GetClientFunc: func() client.Client {
			return &ClientMock{
				StatusFunc: func() client.SubResourceWriter {
					return &SubResourceWriterMock{
						UpdateFunc: func(_ context.Context, _ client.Object, _ ...client.SubResourceUpdateOption) error {
							return nil
						},
					}
				},
				GetFunc: func(_ context.Context, _ client.ObjectKey, _ client.Object, _ ...client.GetOption) error {
					return nil
				},
				CreateFunc: func(_ context.Context, _ client.Object, _ ...client.CreateOption) error {
					return nil
				},
				DeleteFunc: func(_ context.Context, _ client.Object, _ ...client.DeleteOption) error {
					return nil
				},
				DeleteAllOfFunc: func(_ context.Context, _ client.Object, _ ...client.DeleteAllOfOption) error {
					return nil
				},
			}
		},
	}
}

// NewHookManagerMock returns a manager mock whose client routes Create, Get
// and Status().Update to the given hooks, so that tests outside internal/ can
// inject failures without naming controller-runtime option types. A nil hook
// behaves like the corresponding NewDefaultManagerMock operation (success);
// Delete and DeleteAllOf always succeed.
func NewHookManagerMock(
	create func(context.Context, client.Object) error,
	get func(context.Context, client.ObjectKey, client.Object) error,
	statusUpdate func(context.Context, client.Object) error,
) manager.Manager {
	return &ManagerMock{
		GetClientFunc: func() client.Client {
			return &ClientMock{
				StatusFunc: func() client.SubResourceWriter {
					return &SubResourceWriterMock{
						UpdateFunc: func(ctx context.Context, obj client.Object, _ ...client.SubResourceUpdateOption) error {
							if statusUpdate == nil {
								return nil
							}
							return statusUpdate(ctx, obj)
						},
					}
				},
				GetFunc: func(ctx context.Context, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
					if get == nil {
						return nil
					}
					return get(ctx, key, obj)
				},
				CreateFunc: func(ctx context.Context, obj client.Object, _ ...client.CreateOption) error {
					if create == nil {
						return nil
					}
					return create(ctx, obj)
				},
				DeleteFunc: func(_ context.Context, _ client.Object, _ ...client.DeleteOption) error {
					return nil
				},
				DeleteAllOfFunc: func(_ context.Context, _ client.Object, _ ...client.DeleteAllOfOption) error {
					return nil
				},
			}
		},
	}
}

// NewAlreadyExistsError returns the Kubernetes AlreadyExists API error that a
// client's Create returns for an existing batch Job of the given name, so that
// tests outside internal/ can simulate it without importing apimachinery.
func NewAlreadyExistsError(name string) error {
	return apierrors.NewAlreadyExists(schema.GroupResource{Group: "batch", Resource: "jobs"}, name)
}

// NewNotFoundError returns the Kubernetes NotFound API error for a resource of
// the given name, so that tests outside internal/ can simulate it without
// importing apimachinery.
func NewNotFoundError(name string) error {
	return apierrors.NewNotFound(schema.GroupResource{Group: "vald.vdaas.org", Resource: "valdbenchmarkjobs"}, name)
}
