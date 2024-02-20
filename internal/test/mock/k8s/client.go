// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/stretchr/testify/mock"
	"github.com/vdaas/vald/internal/k8s/client"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ValdK8sClientMock struct {
	mock.Mock
}

var _ client.Client = (*ValdK8sClientMock)(nil)

func (m *ValdK8sClientMock) Get(ctx context.Context, name string, namespace string, obj client.Object, opts ...crclient.GetOption) error {
	args := m.Called(ctx, name, namespace, obj, opts)
	return args.Error(0)
}

func (m *ValdK8sClientMock) List(ctx context.Context, list crclient.ObjectList, opts ...client.ListOption) error {
	args := m.Called(ctx, list, opts)
	return args.Error(0)
}

func (m *ValdK8sClientMock) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *ValdK8sClientMock) Delete(ctx context.Context, obj client.Object, opts ...crclient.DeleteOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *ValdK8sClientMock) Update(ctx context.Context, obj client.Object, opts ...crclient.UpdateOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func (m *ValdK8sClientMock) Patch(ctx context.Context, obj client.Object, patch crclient.Patch, opts ...crclient.PatchOption) error {
	args := m.Called(ctx, obj, patch, opts)
	return args.Error(0)
}

func (m *ValdK8sClientMock) Watch(ctx context.Context, obj crclient.ObjectList, opts ...client.ListOption) (watch.Interface, error) {
	args := m.Called(ctx, obj, opts)
	return args.Get(0).(watch.Interface), args.Error(1)
}

func (m *ValdK8sClientMock) LabelSelector(key string, op selection.Operator, vals []string) (labels.Selector, error) {
	args := m.Called(key, op, vals)
	return args.Get(0).(labels.Selector), args.Error(1)
}
