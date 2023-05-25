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
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
)

type Group interface {
	Go(func() error)
	Limitation(int)
	Wait() error
}

type group struct {
	egctx  context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	limitation       chan struct{}
	enableLimitation atomic.Bool
	cancelOnce       sync.Once
	mu               sync.RWMutex
	emap             map[string]struct{}
	errs             []error
	err              error
}

var (
	instance Group
	once     sync.Once
)

func New(ctx context.Context) (Group, context.Context) {
	egctx, cancel := context.WithCancel(ctx)
	g := &group{
		egctx:  egctx,
		emap:   make(map[string]struct{}),
		cancel: cancel,
	}
	g.enableLimitation.Store(false)
	return g, egctx
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

func (g *group) Limitation(limit int) {
	if limit > 0 {
		ch := make(chan struct{}, limit)
		g.closeLimitation()
		g.limitation = ch
		g.enableLimitation.Store(true)
	} else {
		g.enableLimitation.Store(false)
	}
}

func (g *group) closeLimitation() {
	if g.limitation != nil {
		for {
			select {
			case _, live := <-g.limitation:
				if !live {
					return
				}
			default:
				close(g.limitation)
				return
			}
		}
	}
}

func (g *group) Go(f func() error) {
	if f != nil {
		g.wg.Add(1)
		go func() {
			defer g.wg.Done()
			var err error
			if g.enableLimitation.Load() {
				select {
				case <-g.egctx.Done():
					return
				case g.limitation <- struct{}{}:
				}
				err = f()
				select {
				case <-g.limitation:
				case <-g.egctx.Done():
				}
			} else {
				err = f()
			}
			if err != nil && !errors.Is(err, context.Canceled) {
				runtime.Gosched()
				g.mu.RLock()
				_, ok := g.emap[err.Error()]
				g.mu.RUnlock()
				if !ok {
					g.mu.Lock()
					g.errs = append(g.errs, err)
					g.emap[err.Error()] = struct{}{}
					g.mu.Unlock()
				}
				g.doCancel()
				return
			}
		}()
	}
}

func (g *group) doCancel() {
	g.cancelOnce.Do(func() {
		if g.cancel != nil {
			g.cancel()
		}
	})
}

func Wait() error {
	return instance.Wait()
}

func (g *group) Wait() error {
	g.wg.Wait()
	g.doCancel()
	g.closeLimitation()
	g.enableLimitation.Store(false)
	g.mu.RLock()
	switch len(g.errs) {
	case 0:
		g.mu.RUnlock()
		return nil
	case 1:
		g.err = g.errs[0]
	default:
		g.err = g.errs[0]
		for _, err := range g.errs[1:] {
			g.err = errors.Join(g.err, err)
		}
	}
	g.mu.RUnlock()
	return g.err
}
