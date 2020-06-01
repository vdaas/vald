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

// Package observer provides storage observer
package observer

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file/watch"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type StorageObserver interface {
	Start(ctx context.Context) (<-chan error, error)
}

type observer struct {
	w             watch.Watcher
	dir           string
	eg            errgroup.Group
	checkDuration time.Duration

	storage storage.Storage

	ch chan struct{}
}

func New(opts ...Option) (so StorageObserver, err error) {
	o := new(observer)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(o); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	o.w, err = watch.New(
		watch.WithDirs(o.dir),
		watch.WithErrGroup(o.eg),
		watch.WithOnWrite(o.onWrite),
		watch.WithOnCreate(o.onCreate),
	)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (o *observer) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)

	var wech, tech, sech, bech <-chan error
	var err error

	wech, err = o.w.Start(ctx)
	if err != nil {
		return nil, err
	}

	tech, err = o.startTicker(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	sech, err = o.storage.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	bech, err = o.startBackupLoop(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-wech:
			case err = <-tech:
			case err = <-sech:
			case err = <-bech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))

	return ech, nil
}

func (o *observer) startTicker(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)
	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		ct := time.NewTicker(o.checkDuration)
		defer ct.Stop()

		finalize := func() (err error) {
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				return err
			}
			return nil
		}

		for {
			select {
			case <-ctx.Done():
				return finalize()
			case <-ct.C:
				err = o.requestBackup(ctx)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
				}
			}
			if err != nil {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
			}
		}
	}))

	return ech, nil
}

func (o *observer) startBackupLoop(ctx context.Context) (<-chan error, error) {
	o.ch = make(chan struct{}, 1)

	ech := make(chan error, 100)
	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		finalize := func() (err error) {
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				return err
			}
			return nil
		}

		for {
			select {
			case <-ctx.Done():
				return finalize()
			case <-o.ch:
				err = o.backup(ctx)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
				}
			}
		}
	}))

	return ech, nil
}

func (o *observer) onWrite(ctx context.Context, name string) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.onWrite")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	return o.requestBackup(ctx)
}

func (o *observer) onCreate(ctx context.Context, name string) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.onCreate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	return o.requestBackup(ctx)
}

func (o *observer) requestBackup(ctx context.Context) error {
	select {
	case o.ch <- struct{}{}:
	default:
		log.Debug("cannot request backup: channel is full")
	}

	return nil
}

func (o *observer) backup(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.backup")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	pr, pw := io.Pipe()
	defer func() {
		e := pr.Close()
		if e != nil {
			log.Errorf("error on closing pipe reader: %s", e)
		}
	}()

	sw, err := o.storage.Writer(ctx)
	if err != nil {
		return err
	}
	defer func() {
		e := sw.Close()
		if e != nil {
			log.Errorf("error on closing blob-storage writer: %s", e)
		}
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer wg.Done()
		defer func() {
			e := pw.Close()
			if e != nil {
				log.Errorf("error on closing pipe writer: %s", e)
			}
		}()

		tw := tar.NewWriter(pw)
		defer func() {
			e := tw.Close()
			if e != nil {
				log.Errorf("error on closing tar writer: %s", e)
			}
		}()

		err = filepath.Walk(o.dir, func(file string, fi os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(o.dir, file)
			if err != nil {
				return err
			}

			header.Name = filepath.ToSlash(rel)

			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}

			log.Debug("writing: ", file)

			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}

				_, err = io.Copy(tw, data)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}))

	_, err = io.Copy(sw, pr)
	if err != nil {
		return err
	}

	wg.Wait()

	return nil
}
