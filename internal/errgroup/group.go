//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
