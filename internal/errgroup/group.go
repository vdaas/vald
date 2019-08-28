// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package errgroup provides server global wait group for graceful kill all goroutine
package errgroup

import (
	"context"
	"sync"

	"github.com/vdaas/vald/internal/errors"
)

type Group interface {
	Go(func() error)
	Wait() error
}

type group struct {
	cancel func()

	wg sync.WaitGroup

	cancelOnce sync.Once
	mu         sync.RWMutex
	emap       map[string]struct{}
	err        error
}

var (
	instance Group
	once     sync.Once
)

func New(ctx context.Context) (context.Context, Group) {
	egctx, cancel := context.WithCancel(ctx)
	return egctx, &group{
		emap:   make(map[string]struct{}),
		cancel: cancel,
	}
}

func Init(ctx context.Context) (egctx context.Context) {
	once.Do(func() {
		egctx, instance = New(ctx)
	})
	return
}

func Go(f func() error) {
	instance.Go(f)
}

func (g *group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.mu.Lock()
			g.emap[err.Error()] = struct{}{}
			g.mu.Unlock()
			g.cancelOnce.Do(func() {
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func Wait() error {
	return instance.Wait()
}

func (g *group) Wait() error {
	g.wg.Wait()
	g.cancelOnce.Do(func() {
		if g.cancel != nil {
			g.cancel()
		}
	})
	g.mu.RLock()
	for msg := range g.emap {
		g.err = errors.Wrap(g.err, msg)
	}
	g.mu.RUnlock()
	return g.err
}
