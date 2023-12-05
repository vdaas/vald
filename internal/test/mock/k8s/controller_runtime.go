// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type SubResourceWriterMock struct {
	client.SubResourceWriter

	UpdateFunc func(context.Context, client.Object, ...client.SubResourceUpdateOption) error
}

func (sm *SubResourceWriterMock) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	return sm.UpdateFunc(ctx, obj, opts...)
}

type ClientMock struct {
	client.Client

	StatusFunc      func() client.SubResourceWriter
	GetFunc         func(context.Context, client.ObjectKey, client.Object, ...client.GetOption) error
	CreateFunc      func(context.Context, client.Object, ...client.CreateOption) error
	DeleteFunc      func(context.Context, client.Object, ...client.DeleteOption) error
	DeleteAllOfFunc func(context.Context, client.Object, ...client.DeleteAllOfOption) error
}

func (cm *ClientMock) Status() client.SubResourceWriter {
	return cm.StatusFunc()
}

func (cm *ClientMock) Get(ctx context.Context, objKey client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	return cm.GetFunc(ctx, objKey, obj, opts...)
}

func (cm *ClientMock) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return cm.CreateFunc(ctx, obj, opts...)
}

func (cm *ClientMock) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return cm.DeleteFunc(ctx, obj, opts...)
}

func (cm *ClientMock) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return cm.DeleteAllOfFunc(ctx, obj, opts...)
}

type ManagerMock struct {
	manager.Manager

	GetClientFunc func() client.Client
}

func (mm *ManagerMock) GetClient() client.Client {
	return mm.GetClientFunc()
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
