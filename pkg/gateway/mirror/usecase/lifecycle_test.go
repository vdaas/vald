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
	"testing"
	"time"

	client "github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// fakeMirrorClient records the context it was started with and whether it
// was stopped.
type fakeMirrorClient struct {
	client.Client
	ctx     context.Context //nolint:containedctx // records the Start context under test
	stopped bool
}

func (f *fakeMirrorClient) Start(ctx context.Context) (<-chan error, error) {
	f.ctx = ctx
	// never closed: the error forwarding loop of Start receives from it.
	return make(chan error), nil
}

func (f *fakeMirrorClient) Stop(context.Context) error {
	f.stopped = true
	return nil
}

type fakeServer struct{}

func (fakeServer) ListenAndServe(context.Context) <-chan error { return nil }
func (fakeServer) Shutdown(context.Context) error              { return nil }

// Test_run_outboundClientOutlivesShutdownSignal pins that the shutdown signal
// (cancellation of the Start context) does not stop the outbound client, so
// the server can still drain requests that use it; PostStop stops it.
func Test_run_outboundClientOutlivesShutdownSignal(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)
	mc := new(fakeMirrorClient)
	r := &run{
		eg:     eg,
		server: fakeServer{},
		client: mc,
	}
	sigctx, signal := context.WithCancel(egctx)
	if _, err := r.Start(sigctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	signal()
	if mc.ctx == nil {
		t.Fatal("client was not started")
	}
	if err := mc.ctx.Err(); err != nil {
		t.Errorf("client context is done after the shutdown signal: %v", err)
	}
	if err := r.Stop(ctx); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if err := mc.ctx.Err(); err != nil {
		t.Errorf("client context is done before PostStop: %v", err)
	}
	if err := r.PostStop(ctx); err != nil {
		t.Fatalf("PostStop() error = %v", err)
	}
	if !mc.stopped {
		t.Error("client was not stopped by PostStop")
	}
	if mc.ctx.Err() == nil {
		t.Error("client context is not cancelled after PostStop")
	}
	done := make(chan error, 1)
	go func() { done <- eg.Wait() }()
	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("timed out waiting for the error forwarding goroutine")
	}
}
