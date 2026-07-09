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

package portforward

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// mockClientSet wires the client-go fake clientset into the kclient.ClientSet
// interface used by the forwarder.
type mockClientSet struct {
	cs kubernetes.Interface
}

func (m *mockClientSet) GetClientSet() kubernetes.Interface {
	return m.cs
}

func (*mockClientSet) GetRESTConfig() *rest.Config {
	return &rest.Config{Host: "https://127.0.0.1:6443"}
}

func newTestForwarder(t *testing.T) Forwarder {
	t.Helper()
	pf, err := New(
		WithClient(&mockClientSet{cs: fake.NewClientset()}),
		WithNamespace("default"),
		WithServiceName("vald-lb-gateway"),
		WithAddress("127.0.0.1"),
		WithPorts(map[uint16]uint16{18081: 8081}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return pf
}

// TestPortforwardExtended_DoesNotMutateSharedClient guards against the SPDY
// transport being assigned to the caller-provided (possibly process-shared)
// http.Client: the session must build its own client instead.
func TestPortforwardExtended_DoesNotMutateSharedClient(t *testing.T) {
	t.Parallel()

	shared := &http.Client{Timeout: 3 * time.Second}

	cancel, _, err := portforwardExtended(
		context.Background(),
		&mockClientSet{cs: fake.NewClientset()},
		"default", "vald-agent-0",
		nil, // no addresses so the call fails before dialing
		map[uint16]uint16{18081: 8081},
		shared,
	)
	if cancel != nil {
		cancel()
	}
	if !errors.Is(err, errors.ErrPortForwardAddressNotFound) {
		t.Fatalf("portforwardExtended() error = %v, want %v", err, errors.ErrPortForwardAddressNotFound)
	}
	if shared.Transport != nil {
		t.Errorf("shared http.Client transport mutated to %T, want nil", shared.Transport)
	}
	if shared.Timeout != 3*time.Second {
		t.Errorf("shared http.Client timeout mutated to %v", shared.Timeout)
	}
}

// TestPortForward_StopBeforeStart guards Stop against nil panics when Start
// was never called.
func TestPortForward_StopBeforeStart(t *testing.T) {
	t.Parallel()

	pf := newTestForwarder(t)
	if err := pf.Stop(); err != nil {
		t.Errorf("Stop() before Start error = %v, want nil", err)
	}
	if err := pf.Stop(); err != nil {
		t.Errorf("second Stop() error = %v, want nil", err)
	}
}

// TestPortForward_StartStopIdempotent verifies that Start returns once the
// context expires (no busy-wait hang) and that Stop is idempotent afterwards
// (no double close of the error channel).
func TestPortForward_StartStopIdempotent(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	pf := newTestForwarder(t)
	ech, err := pf.Start(ctx)
	if ech == nil {
		t.Fatal("Start() ech = nil, want non-nil channel")
	}
	if err == nil {
		t.Error("Start() error = nil, want context error because no pod ever becomes healthy")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		// First Stop drains the goroutines, second must be a no-op.
		_ = pf.Stop()
		if serr := pf.Stop(); serr != nil {
			t.Errorf("second Stop() error = %v, want nil", serr)
		}
	}()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("Stop() did not return: goroutines leaked or deadlocked")
	}
}
