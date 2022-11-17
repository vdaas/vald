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
package servers

import (
	"context"
)

type mockServer struct {
	NameFunc           func() string
	IsRunningFunc      func() bool
	ListenAndServeFunc func(context.Context, chan<- error) error
	ShutdownFunc       func(context.Context) error
}

func (ms *mockServer) Name() string {
	return ms.NameFunc()
}

func (ms *mockServer) IsRunning() bool {
	return ms.IsRunningFunc()
}

func (ms *mockServer) ListenAndServe(ctx context.Context, errCh chan<- error) error {
	return ms.ListenAndServeFunc(ctx, errCh)
}

func (ms *mockServer) Shutdown(ctx context.Context) error {
	return ms.ShutdownFunc(ctx)
}
