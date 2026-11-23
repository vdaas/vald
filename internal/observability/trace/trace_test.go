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

package trace

import (
	"context"
	"testing"
)

func TestEnable(t *testing.T) {
	t.Cleanup(func() { enabled.Store(false) })

	ctx := context.Background()
	if Enabled() {
		t.Fatal("tracing must be disabled by default")
	}
	if got, span := StartSpan(ctx, "disabled"); span != nil || got != ctx {
		t.Fatalf("disabled StartSpan must return the input ctx and a nil span, got span=%v", span)
	}

	Enable()
	if !Enabled() {
		t.Fatal("Enable must turn tracing on")
	}
	got, span := StartSpan(ctx, "enabled")
	if span == nil {
		t.Fatal("enabled StartSpan must return a span")
	}
	End(span)
	if got == ctx {
		t.Fatal("enabled StartSpan must return a derived ctx")
	}
}
