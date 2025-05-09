// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package transport

import (
	"context"
	"net/http"
)

type roundTripMock struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (rm *roundTripMock) RoundTrip(req *http.Request) (*http.Response, error) {
	return rm.RoundTripFunc(req)
}

type backoffMock struct {
	DoFunc    func(context.Context, func(context.Context) (any, bool, error)) (any, error)
	CloseFunc func()
}

func (bm *backoffMock) Do(
	ctx context.Context, fn func(context.Context) (any, bool, error),
) (any, error) {
	return bm.DoFunc(ctx, fn)
}

func (bm *backoffMock) Close() {
	bm.CloseFunc()
}
