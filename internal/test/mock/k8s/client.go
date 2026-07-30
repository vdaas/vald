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

	"github.com/stretchr/testify/mock"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ValdK8sClientMock is a testify-based mock of client.Client, the unified
// interface replacing the former StandaloneClient/ClientSet split. The
// embedded k8s.Client promotes the entire controller-runtime CRUD surface
// (Get/List/Create/Update/Delete/Patch/DeleteAllOf/Apply/Status/SubResource/
// Scheme/RESTMapper/GroupVersionKindFor/IsObjectNamespaced): WithRaw backs it
// with the same fake client Raw() returns, so every promoted call and Raw()
// agree. Only the client.Client additions (GetClientSet/GetRESTConfig) are
// overridden below via testify, because callers set expectations on those
// directly rather than through a fake client.
type ValdK8sClientMock struct {
	k8s.Client
	mock.Mock
}

var _ client.Client = (*ValdK8sClientMock)(nil)

// GetClientSet returns the value configured via WithClientSet. Tests that
// exercise code paths reading it must set it via WithClientSet before use.
func (m *ValdK8sClientMock) GetClientSet() kubernetes.Interface {
	args := m.Called()
	cs, _ := args.Get(0).(kubernetes.Interface)
	return cs
}

// WithClientSet sets the value GetClientSet() returns and returns m for
// chaining.
func (m *ValdK8sClientMock) WithClientSet(cs kubernetes.Interface) *ValdK8sClientMock {
	m.On("GetClientSet").Return(cs)
	return m
}

// GetRESTConfig returns the value configured via WithRESTConfig. Tests that
// exercise code paths reading it must set it via WithRESTConfig before use.
func (m *ValdK8sClientMock) GetRESTConfig() *rest.Config {
	args := m.Called()
	cfg, _ := args.Get(0).(*rest.Config)
	return cfg
}

// WithRESTConfig sets the value GetRESTConfig() returns and returns m for
// chaining.
func (m *ValdK8sClientMock) WithRESTConfig(cfg *rest.Config) *ValdK8sClientMock {
	m.On("GetRESTConfig").Return(cfg)
	return m
}

// Raw returns the controller-runtime client backing calls made through
// resource.Client (built via resource.NewClientOf). Tests that exercise those
// paths must set it via WithRaw before use; other tests that only stub this
// mock's own methods (Get/List/Apply) never call it.
func (m *ValdK8sClientMock) Raw() crclient.WithWatch {
	args := m.Called()
	raw, _ := args.Get(0).(crclient.WithWatch)
	return raw
}

// WithRaw sets the value Raw() returns, backs the embedded k8s.Client with
// the same client so promoted methods (Create/Update/Delete/...) agree with
// Raw(), and returns m for chaining.
func (m *ValdK8sClientMock) WithRaw(raw crclient.WithWatch) *ValdK8sClientMock {
	m.On("Raw").Return(raw)
	m.Client = raw
	return m
}

type PatcherMock struct {
	mock.Mock
}

var _ client.Patcher = (*PatcherMock)(nil)

func (m *PatcherMock) ApplyPodAnnotations(
	ctx context.Context, name, namespace string, entries map[string]string,
) error {
	args := m.Called(ctx, name, namespace, entries)
	return args.Error(0)
}
