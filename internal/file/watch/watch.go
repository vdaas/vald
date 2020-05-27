//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package watch

import (
	"context"
	"reflect"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type Watcher interface {
	Start(ctx context.Context) (<-chan error, error)
	Add(dirs ...string) error
	Remove(dirs ...string) error
	Stop(ctx context.Context) error
}

type watch struct {
	w        *fsnotify.Watcher
	eg       errgroup.Group
	dirs     map[string]struct{}
	mu       sync.RWMutex
	onChange func(ctx context.Context, name string) error
	onCreate func(ctx context.Context, name string) error
	onRename func(ctx context.Context, name string) error
	onDelete func(ctx context.Context, name string) error
	onWrite  func(ctx context.Context, name string) error
	onChmod  func(ctx context.Context, name string) error
	onError  func(ctx context.Context, err error) error
}

func New(opts ...Option) (Watcher, error) {
	w := new(watch)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(w); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if len(w.dirs) == 0 {
		return nil, errors.ErrWatchDirNotFound
	}
	return w.init()
}

func (w *watch) init() (*watch, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w.mu.RLock()
	dirs := w.dirs
	w.mu.RUnlock()
	for dir := range dirs {
		log.Debug("Adding watcher target: ", dir)

		err = watcher.Add(dir)
		if err != nil {
			return nil, err
		}
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.w != nil {
		err = w.w.Close()
		if err != nil {
			return nil, err
		}
	}
	w.w = watcher

	return w, nil
}

func (w *watch) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 10)
	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		var (
			event fsnotify.Event
			ok    bool
		)
		handleErr := func(ctx context.Context, err error) {
			log.Error(err)
			select {
			case <-ctx.Done():
			case ech <- err:
			}
		}
		for {
			ok = true
			err = nil
			select {
			case <-ctx.Done():
				return ctx.Err()
			case event, ok = <-w.w.Events:
				if ok {
					log.Debug("change detected file: ", event.Name)
					if w.onChange != nil {
						err = w.onChange(ctx, event.Name)
						if err != nil {
							handleErr(ctx, err)
						}
						err = nil
					}
					switch {
					case event.Op&fsnotify.Write == fsnotify.Write && w.onWrite != nil:
						log.Debug("Modified file: ", event.Name)
						err = w.onWrite(ctx, event.Name)
					case event.Op&fsnotify.Create == fsnotify.Create && w.onCreate != nil:
						log.Debug("Created file: ", event.Name)
						err = w.onCreate(ctx, event.Name)
					case event.Op&fsnotify.Remove == fsnotify.Remove && w.onDelete != nil:
						log.Debug("Removed file: ", event.Name)
						err = w.onDelete(ctx, event.Name)
					case event.Op&fsnotify.Rename == fsnotify.Rename && w.onRename != nil:
						log.Debug("Renamed file: ", event.Name)
						err = w.onRename(ctx, event.Name)
					case event.Op&fsnotify.Chmod == fsnotify.Chmod && w.onChmod != nil:
						log.Debug("File changed permission: ", event.Name)
						err = w.onChmod(ctx, event.Name)
					}
				}
			case err, ok = <-w.w.Errors:
			}
			if !ok {
				w, err = w.init()
			}
			if err != nil {
				handleErr(ctx, err)
			}
		}
	}))
	return ech, nil
}

func (w *watch) Add(dirs ...string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, dir := range dirs {
		w.dirs[dir] = struct{}{}
		if w.w != nil {
			err = w.w.Add(dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *watch) Remove(dirs ...string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, dir := range dirs {
		delete(w.dirs, dir)
		if w.w != nil {
			err = w.w.Remove(dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *watch) Stop(ctx context.Context) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for dir := range w.dirs {
		delete(w.dirs, dir)
		if w.w != nil {
			err = w.w.Remove(dir)
			if err != nil {
				return err
			}
		}
	}
	if w.w != nil {
		return w.w.Close()
	}
	return nil
}
