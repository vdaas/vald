//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package service

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/operator/vald/config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type test struct {
		wantErr error
		cfg     *config.Operator
		name    string
	}

	tests := []test{
		{
			name:    "returns ErrInvalidConfig when config is nil",
			cfg:     nil,
			wantErr: errors.ErrInvalidConfig,
		},
		{
			name: "returns error when the default vrs template cannot be loaded",
			cfg: &config.Operator{
				Name: "vald-operator",
				Vrs: &config.Vrs{
					DefaultVrsPath: filepath.Join(t.TempDir(), "missing.yaml"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tc.cfg)
			if err == nil {
				t.Fatal("New() error = nil, want error")
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Errorf("New() error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// ctrlMock is a minimal k8s.Controller stub whose Start returns a
// pre-populated error channel.
type ctrlMock struct {
	dech chan error
}

func (c *ctrlMock) Start(context.Context) (<-chan error, error) {
	return c.dech, nil
}

func (*ctrlMock) GetManager() k8s.Manager {
	return nil
}

// TestOperator_Start_CtxCancelUnblocksErrorForward pins the ctx guard on the
// ech forward: with no reader on ech (capacity 2), the third controller error
// blocks the send, and only a ctx-guarded select lets the goroutine exit on
// cancellation instead of leaking.
func TestOperator_Start_CtxCancelUnblocksErrorForward(t *testing.T) {
	t.Parallel()

	dech := make(chan error, 3)
	for range 3 {
		dech <- errors.New("controller error")
	}

	eg, egctx := errgroup.New(context.Background())
	ctx, cancel := context.WithCancel(egctx)
	defer cancel()

	o := &operator{ctrl: &ctrlMock{dech: dech}, eg: eg}
	ech, err := o.Start(ctx)
	if err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	// Wait until the forwarding goroutine has filled ech and is blocked on
	// the third send.
	deadline := time.Now().Add(5 * time.Second)
	for len(ech) < cap(ech) {
		if time.Now().After(deadline) {
			t.Fatal("error forwarding did not fill the channel buffer")
		}
		time.Sleep(time.Millisecond)
	}

	cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = eg.Wait()
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("forwarding goroutine did not exit on context cancel; the ech send must be ctx-guarded")
	}
}
