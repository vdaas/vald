//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package errgroup provides server global wait group for graceful kill all goroutine
package errgroup

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/sync"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/semaphore"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid, has no limit on the number of active goroutines,
// and does not cancel on error.
type Group interface {
	Go(func() error)
	SetLimit(limit int)
	TryGo(func() error) bool
	Wait() error
}

type group struct {
	egctx  context.Context
	cancel context.CancelCauseFunc

	wg sync.WaitGroup

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

func New(ctx context.Context) (Group, context.Context) {
	g := &group{emap: make(map[string]struct{})}
	g.egctx, g.cancel = context.WithCancelCause(ctx)
	return g, g.egctx
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (Group, context.Context) {
	return New(ctx)
}

func Init(ctx context.Context) (egctx context.Context) {
	egctx = ctx
	once.Do(func() {
		instance, egctx = New(ctx)
	})
	return
}

func Get() Group {
	if instance == nil {
		Init(context.Background())
	}
	return instance
}

func Go(f func() error) {
	instance.Go(f)
}

func TryGo(f func() error) bool {
	return instance.TryGo(f)
}

// SetLimit limits the number of active goroutines in this group to at most n.
// A negative value indicates no limit.
//
// Any subsequent call to the Go method will block until it can add an active
// goroutine without exceeding the configured limit.
//
// The limit must not be modified while any goroutines in the group are active.
func (g *group) SetLimit(limit int) {
	if limit < 0 {
		g.sem = nil
		return
	}

	if g.sem == nil {
		g.sem = semaphore.NewWeighted(int64(limit))
	} else {
		g.sem.Resize(int64(limit))
	}
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext. The error will be returned by Wait.
func (g *group) Go(f func() error) {
	if f == nil {
		return
	}
	if g.sem != nil {
		err := g.sem.Acquire(g.egctx, 1)
		if err != nil {
			if !errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				g.appendErr(err)
			}
			return
		}
	}
	g.run(f)
}

// TryGo calls the given function in a new goroutine only if the number of
// active goroutines in the group is currently below the configured limit.
//
// The return value reports whether the goroutine was started.
func (g *group) TryGo(f func() error) bool {
	if f == nil {
		return false
	}
	if g.sem != nil && !g.sem.TryAcquire(1) {
		return false
	}
	g.run(f)
	return true
}

func (g *group) run(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.done()
		err := f()
		if err != nil {
			if !errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				runtime.Gosched()
				g.appendErr(err)
			}
			g.doCancel(err)
			return
		}
	}()
}

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

func (g *group) done() {
	defer g.wg.Done()
	if g.sem != nil {
		g.sem.Release(1)
	}
}

func (g *group) doCancel(err error) {
	g.cancelOnce.Do(func() {
		if g.cancel != nil {
			g.cancel(err)
		}
	})
}

func Wait() error {
	return instance.Wait()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *group) Wait() (err error) {
	g.wg.Wait()
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
