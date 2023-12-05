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

	"github.com/vdaas/vald/internal/k8s"
)

type ControllerMock struct {
	StartFunc      func(ctx context.Context) (<-chan error, error)
	GetManagerFunc func() k8s.Manager
}

func (m *ControllerMock) Start(ctx context.Context) (<-chan error, error) {
	return m.StartFunc(ctx)
}

func (m *ControllerMock) GetManager() k8s.Manager {
	return m.GetManagerFunc()
}
