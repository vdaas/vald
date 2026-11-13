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
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// startErrManager is a manager.Manager whose Start fails; every other method
// is unused by controller.Start with no resource controllers.
type startErrManager struct {
	manager.Manager
	err error
}

func (m *startErrManager) Start(context.Context) error {
	return m.err
}

func TestControllerStartErrorCarriesName(t *testing.T) {
	t.Parallel()
	errStart := errors.New("manager start failed")
	tests := []struct {
		name       string
		ctrlName   string
		wantPrefix string
	}{
		{
			name:       "error is prefixed with the controller name set by WithControllerName",
			ctrlName:   "vald benchmark scenario operator",
			wantPrefix: "controller vald benchmark scenario operator: ",
		},
		{
			name:     "error is unchanged when no controller name is set",
			ctrlName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, egctx := errgroup.New(ctx)
			c := &controller{eg: eg, mgr: &startErrManager{err: errStart}}
			if err := WithControllerName(tt.ctrlName)(c); err != nil {
				t.Fatalf("WithControllerName() error = %v", err)
			}
			ech, err := c.Start(egctx)
			if err != nil {
				t.Fatalf("Start() error = %v", err)
			}
			var got error
			select {
			case got = <-ech:
			case <-ctx.Done():
				t.Fatal("timed out waiting for the manager start error")
			}
			if !errors.Is(got, errStart) {
				t.Fatalf("got error %v, want it to wrap %v", got, errStart)
			}
			if tt.wantPrefix == "" {
				if got.Error() != errStart.Error() {
					t.Errorf("got error %q, want %q", got.Error(), errStart.Error())
				}
			} else if !strings.HasPrefix(got.Error(), tt.wantPrefix) {
				t.Errorf("got error %q, want prefix %q", got.Error(), tt.wantPrefix)
			}
			if err := eg.Wait(); err != nil {
				t.Errorf("errgroup Wait() error = %v", err)
			}
		})
	}
}
