//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
package runner

import "context"

type runnerMock struct {
	PreStartFunc func(ctx context.Context) error
	StartFunc    func(ctx context.Context) (<-chan error, error)
	PreStopFunc  func(ctx context.Context) error
	StopFunc     func(ctx context.Context) error
	PostStopFunc func(ctx context.Context) error
}

func (m *runnerMock) PreStart(ctx context.Context) error {
	return m.PreStartFunc(ctx)
}

func (m *runnerMock) Start(ctx context.Context) (<-chan error, error) {
	return m.StartFunc(ctx)
}

func (m *runnerMock) PreStop(ctx context.Context) error {
	return m.PreStopFunc(ctx)
}

func (m *runnerMock) Stop(ctx context.Context) error {
	return m.StopFunc(ctx)
}

func (m *runnerMock) PostStop(ctx context.Context) error {
	return m.PostStopFunc(ctx)
}
