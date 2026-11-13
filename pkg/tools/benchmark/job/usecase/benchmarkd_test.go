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
package usecase

import (
	"context"
	"slices"
	"testing"
)

type stopOrderObs struct{ calls *[]string }

func (*stopOrderObs) PreStart(context.Context) error { return nil }

func (*stopOrderObs) Start(context.Context) <-chan error { return nil }

func (o *stopOrderObs) Stop(context.Context) error {
	*o.calls = append(*o.calls, "observability")
	return nil
}

type stopOrderSrv struct{ calls *[]string }

func (*stopOrderSrv) ListenAndServe(context.Context) <-chan error { return nil }

func (s *stopOrderSrv) Shutdown(context.Context) error {
	*s.calls = append(*s.calls, "server")
	return nil
}

// Stop used to stop observability before shutting the server down, so the
// telemetry of the requests drained during Shutdown was dropped.
func TestRun_Stop_shutsDownServerBeforeObservability(t *testing.T) {
	t.Parallel()
	var calls []string
	r := &run{
		server:        &stopOrderSrv{calls: &calls},
		observability: &stopOrderObs{calls: &calls},
	}
	if err := r.Stop(t.Context()); err != nil {
		t.Fatalf("Stop() = %v, want nil", err)
	}
	if want := []string{"server", "observability"}; !slices.Equal(calls, want) {
		t.Errorf("Stop() call order = %v, want %v", calls, want)
	}
}
