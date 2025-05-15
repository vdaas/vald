//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// It is a custom implementation similar to sync/errgroup.
package errgroup

import (
	"context"
	"runtime"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/semaphore"
)

// Group is a collection of goroutines working on subtasks that are part of the same overall task.
// A zero Group is valid; it has no limit on the number of active goroutines and does not cancel on error.
type Group interface {
	// Go starts the given function either in a new goroutine or inline (if limit == 1).
	Go(func() error)
	// SetLimit sets the maximum number of active goroutines.
	SetLimit(limit int)
	// TryGo attempts to start the given function, returning true if it was started.
	TryGo(func() error) bool
	// Wait blocks until all tasks started with Go have completed, returning the first non-nil error (if any).
	Wait() error
}

// group is the concrete implementation of Group.
type group struct {
	egctx  context.Context
	cancel context.CancelCauseFunc

	wg sync.WaitGroup

	// limit controls how many tasks can run concurrently.
	limit atomic.Int64
	// sem is used to limit concurrent goroutines when limit > 1.
	sem *semaphore.Weighted

	cancelOnce sync.Once
	mu         sync.RWMutex
	emap       map[string]struct{}
	errs       []error
}

var (
	instance Group
	once     sync.Once
)

// New creates a new Group and returns it along with its derived Context.
func New(ctx context.Context) (Group, context.Context) {
	g := &group{
		emap: make(map[string]struct{}),
	}
	// Create a context that can be canceled with a cause.
	g.egctx, g.cancel = context.WithCancelCause(ctx)
	return g, g.egctx
}

// WithContext returns a new Group and an associated Context derived from ctx.
// The derived Context is canceled the first time a function passed to Go returns a non-nil error
// or the first time Wait returns, whichever occurs first.
func WithContext(ctx context.Context) (Group, context.Context) {
	return New(ctx)
}

// Init initializes the global errgroup instance.
func Init(ctx context.Context) (egctx context.Context) {
	egctx = ctx
	once.Do(func() {
		instance, egctx = New(ctx)
	})
	return
}

// Get returns the global errgroup instance, initializing it if necessary.
func Get() Group {
	if instance == nil {
		Init(context.Background())
	}
	return instance
}

// Go is a package-level helper that calls the Go method on the global instance.
func Go(f func() error) {
	if instance == nil {
		Init(context.Background())
	}
	instance.Go(f)
}

// TryGo is a package-level helper that calls the TryGo method on the global instance.
func TryGo(f func() error) bool {
	if instance == nil {
		Init(context.Background())
	}
	return instance.TryGo(f)
}

// SetLimit sets the maximum number of active goroutines in the group.
// A negative value indicates no limit.
// This must not be modified while any tasks are active.
func (g *group) SetLimit(limit int) {
	g.limit.Store(int64(limit))
	if limit < 0 {
		return
	}
	// For concurrent execution, initialize or resize the semaphore.
	if g.sem == nil {
		g.sem = semaphore.NewWeighted(int64(limit))
	} else {
		g.sem.Resize(int64(limit))
	}
}

// exec executes the provided function inline (synchronously) when limit == 1.
// It wraps the call with wait group operations and reuses executeTask for error handling.
// Performance Note: Inline execution avoids the overhead of goroutine scheduling and context switching.
// run schedules the provided function to run in a new goroutine (asynchronously).
// It wraps the call with wait group operations and reuses executeTask for error handling.
func (g *group) exec(f func() error) {
	// Execute the task function.
	err := f()
	if err != nil {
		// If the error is not due to cancellation or deadline, yield and record it.
		if errors.IsNot(err, context.Canceled, context.DeadlineExceeded) {
			runtime.Gosched()
			g.appendErr(err)
		}
		// Cancel the context with the encountered error.
		g.doCancel(err)
	}
}

// run schedules the provided function to run in a new goroutine (asynchronously).
// It wraps the call with wait group operations and reuses executeTask for error handling.
func (g *group) run(f func() error) {
	g.wg.Add(1)
	go func() {
		// done() will call wg.Done() and release the semaphore slot.
		defer g.done()
		g.exec(f)
	}()
}

// Go calls the given function either in a new goroutine or inline based on the limit.
// For limit == 1, the function is executed inline to avoid unnecessary goroutine creation.
// For limit > 1, the function is scheduled in a new goroutine after acquiring the semaphore.
func (g *group) Go(f func() error) {
	if f == nil {
		return
	}
	// Check if we should execute inline (serial execution).
	if g.limit.Load() == 1 {
		g.exec(f)
		return
	}
	// In concurrent mode, acquire the semaphore before launching a new goroutine.
	if g.sem != nil {
		err := g.sem.Acquire(g.egctx, 1)
		if err != nil {
			// Handle errors from semaphore acquisition if not due to cancellation or deadline.
			if errors.IsNot(err, context.Canceled, context.DeadlineExceeded) {
				g.appendErr(err)
			}
			return
		}
	}
	g.run(f)
}

// TryGo attempts to run the function, starting it in a new goroutine or inline if limit == 1,
// only if the number of active tasks is below the configured limit.
// Returns true if the function was executed, false otherwise.
func (g *group) TryGo(f func() error) bool {
	if f == nil {
		return false
	}
	// In concurrent mode, try to acquire the semaphore without blocking.
	if g.sem != nil && !g.sem.TryAcquire(1) {
		return false
	}
	g.run(f)
	return true
}

// appendErr appends the error to the group's error list if it has not been recorded before.
func (g *group) appendErr(err error) {
	g.mu.RLock()
	_, ok := g.emap[err.Error()]
	g.mu.RUnlock()
	if !ok {
		g.mu.Lock()
		g.errs = append(g.errs, err)
		g.emap[err.Error()] = struct{}{}
		g.mu.Unlock()
	}
}

// done releases the semaphore (if used) and marks the task as done in the wait group.
func (g *group) done() {
	defer g.wg.Done()
	if g.sem != nil {
		g.sem.Release(1)
	}
}

// doCancel cancels the group's context with the provided error.
// It ensures that cancellation is performed only once.
func (g *group) doCancel(err error) {
	g.cancelOnce.Do(func() {
		if g.cancel != nil {
			g.cancel(err)
		}
	})
}

// Wait blocks until all tasks started with Go have completed.
// It returns the first non-nil error (if any) from the executed tasks.
func (g *group) Wait() (err error) {
	g.wg.Wait()
	// After all tasks complete, cancel the context to propagate cancellation if needed.
	g.doCancel(context.Canceled)
	g.mu.RLock()
	defer g.mu.RUnlock()
	switch len(g.errs) {
	case 0:
		return nil
	case 1:
		return g.errs[0]
	default:
		return errors.Join(g.errs...)
	}
}

// Wait is a package-level helper that calls the Wait method on the global instance.
func Wait() error {
	if instance == nil {
		return nil
	}
	return instance.Wait()
}
