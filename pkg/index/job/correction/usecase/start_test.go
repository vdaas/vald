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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/index/job/correction/service"
)

type fakeListener struct{}

func (*fakeListener) ListenAndServe(context.Context) <-chan error {
	return nil
}

func (*fakeListener) Shutdown(context.Context) error {
	return nil
}

// fakeJob embeds the service interface so that only the methods used by
// run.Start need to be implemented.
type fakeJob struct {
	service.Corrector
}

func (*fakeJob) StartClient(context.Context) (<-chan error, error) {
	return make(chan error), nil
}

func (*fakeJob) Start(context.Context) error {
	return nil
}

// TestRun_Start_FindProcessError verifies that Start does not terminate the
// process (log.Fatal) when the job cannot find its own process to SIGTERM.
// With the old log.Fatalf this test binary would exit instead of reaching the
// assertions below.
func TestRun_Start_FindProcessError(t *testing.T) {
	// not parallel: os.FindProcess is a package level variable.
	called := make(chan struct{}, 1)
	orig := os.FindProcess
	t.Cleanup(func() { os.FindProcess = orig })
	os.FindProcess = func(int) (*os.Process, error) {
		select {
		case called <- struct{}{}:
		default:
		}
		return nil, errors.New("find process failed")
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	r := &run{
		eg:        eg,
		server:    new(fakeListener),
		corrector: new(fakeJob),
	}
	if _, err := r.Start(ctx); err != nil {
		t.Fatalf("Start returned an error: %v", err)
	}

	select {
	case <-called:
	case <-time.After(10 * time.Second):
		t.Fatal("os.FindProcess was not called by the job goroutine")
	}

	cancel()
	done := make(chan error, 1)
	go func() {
		done <- eg.Wait()
	}()
	select {
	case err := <-done:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Errorf("unexpected error from errgroup: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("goroutines started by Start did not finish")
	}
}
