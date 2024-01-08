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
package mock

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type MockSubResourceWriter struct {
	client.SubResourceWriter
}

func (*MockSubResourceWriter) Update(context.Context, client.Object, ...client.SubResourceUpdateOption) error {
	return nil
}

type MockClient struct {
	client.Client
}

func (*MockClient) Status() client.SubResourceWriter {
	s := MockSubResourceWriter{
		SubResourceWriter: &MockSubResourceWriter{},
	}
	return s.SubResourceWriter
}

func (*MockClient) Get(context.Context, client.ObjectKey, client.Object, ...client.GetOption) error {
	return nil
}

func (*MockClient) Create(context.Context, client.Object, ...client.CreateOption) error {
	return nil
}

func (*MockClient) Delete(context.Context, client.Object, ...client.DeleteOption) error {
	return nil
}

func (*MockClient) DeleteAllOf(context.Context, client.Object, ...client.DeleteAllOfOption) error {
	return nil
}

type MockManager struct {
	manager.Manager
}

func (*MockManager) GetClient() client.Client {
	c := &MockClient{
		Client: &MockClient{},
	}
	return c.Client
}
