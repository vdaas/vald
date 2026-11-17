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

package watch

import (
	"context"
	"reflect"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// reinitBackoff is the wait time before retrying to re-initialize the watcher after a failure.
const reinitBackoff = time.Second

// Watcher is an interface that represents a file monitor.
type Watcher interface {
	Start(ctx context.Context) (<-chan error, error)
	Add(dirs ...string) error
	Remove(dirs ...string) error
	Stop(ctx context.Context) error
}

type watch struct {
	eg       errgroup.Group
	w        *fsnotify.Watcher
	dirs     map[string]struct{}
	onChange func(ctx context.Context, name string) error
	onCreate func(ctx context.Context, name string) error
	onRename func(ctx context.Context, name string) error
	onDelete func(ctx context.Context, name string) error
	onWrite  func(ctx context.Context, name string) error
	onChmod  func(ctx context.Context, name string) error
	mu       sync.RWMutex
	// stopped is set by Stop so that the Start goroutine neither re-initializes nor keeps watching afterwards.
	stopped bool
}

// New returns Watcher implementation.
func New(opts ...Option) (Watcher, error) {
	w := new(watch)
	for _, opt := range append(defaultOptions, opts...) {
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
	dirs := make([]string, 0, len(w.dirs))
	for d := range w.dirs {
		dirs = append(dirs, d)
	}
	w.mu.RUnlock()
	for _, dir := range dirs {
		log.Infof("Watching: %s", dir)

		err = watcher.Add(dir)
		if err != nil {
			return nil, errors.Join(err, watcher.Close())
		}
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		// Stop ran while the new watcher was being prepared: do not install it.
		return w, watcher.Close()
	}
	if w.w != nil {
		err = w.w.Close()
		if err != nil {
			return nil, errors.Join(err, watcher.Close())
		}
	}
	w.w = watcher

	return w, nil
}

// Start starts watching all named files or directories. If an error occurs, returns the error.
// And performs the processing corresponding to the file change event, and returns an error via channel if an error occurs in them.
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
			w.mu.RLock()
			cur, stopped := w.w, w.stopped
			w.mu.RUnlock()
			if cur == nil || stopped {
				// Stop has been called.
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case event, ok = <-cur.Events:
				if ok {
					if w.onChange != nil {
						err = w.onChange(ctx, event.Name)
						if err != nil {
							handleErr(ctx, err)
						}
						err = nil
					}
					switch {
					case event.Op&fsnotify.Write == fsnotify.Write && w.onWrite != nil:
						log.Debugf("File %s modified. Trigger onWrite hook.", event.Name)
						err = w.onWrite(ctx, event.Name)
					case event.Op&fsnotify.Create == fsnotify.Create && w.onCreate != nil:
						log.Debugf("File %s created. Trigger onCreate hook.", event.Name)
						err = w.onCreate(ctx, event.Name)
					case event.Op&fsnotify.Remove == fsnotify.Remove && w.onDelete != nil:
						log.Debugf("File %s deleted. Trigger onDelete hook.", event.Name)
						err = w.onDelete(ctx, event.Name)
					case event.Op&fsnotify.Rename == fsnotify.Rename && w.onRename != nil:
						log.Debugf("File %s renamed. Trigger onRename hook.", event.Name)
						err = w.onRename(ctx, event.Name)
					case event.Op&fsnotify.Chmod == fsnotify.Chmod && w.onChmod != nil:
						log.Debugf("Permission of file %s changed. Trigger onChmod hook.", event.Name)
						err = w.onChmod(ctx, event.Name)
					}
				}
			case err, ok = <-cur.Errors:
			}
			if !ok {
				if _, err = w.init(); err != nil {
					handleErr(ctx, err)
					// the closed watcher is still installed, so wait before retrying to avoid busy-spinning.
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(reinitBackoff):
					}
					continue
				}
			}
			if err != nil {
				handleErr(ctx, err)
			}
		}
	}))
	return ech, nil
}

// Add starts watching all named files or directories. If an error occurs, returns the error.
func (w *watch) Add(dirs ...string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, dir := range dirs {
		if w.w != nil {
			err = w.w.Add(dir)
			if err != nil {
				return err
			}
			w.dirs[dir] = struct{}{}
		}
	}
	return nil
}

// Remove stops watching all named files or directories. If an error occurs, returns the error.
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

// Stop stops watching all named files or directories. If an error occurs, returns the error.
func (w *watch) Stop(context.Context) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.stopped = true
	for dir := range w.dirs {
		delete(w.dirs, dir)
		if w.w != nil {
			// keep going on failure (e.g. the directory was deleted) so that the watcher is always closed below.
			if rerr := w.w.Remove(dir); rerr != nil {
				err = errors.Join(err, rerr)
			}
		}
	}
	if w.w != nil {
		// closing wakes the goroutine started by Start; w.w is cleared so that it terminates instead of re-initializing.
		err = errors.Join(err, w.w.Close())
		w.w = nil
	}
	return err
}
