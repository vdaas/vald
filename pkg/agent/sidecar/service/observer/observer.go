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

// Package observer provides storage observer
package observer

import (
	"archive/tar"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/file/watch"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type StorageObserver interface {
	Start(ctx context.Context) (<-chan error, error)
	PostStop(ctx context.Context) error
}

type observer struct {
	w             watch.Watcher
	dir           string
	eg            errgroup.Group
	checkDuration time.Duration

	metadataPath string

	postStopTimeout time.Duration

	watchEnabled  bool
	tickerEnabled bool

	storage storage.Storage

	ch chan struct{}

	hooks []Hook
}

func New(opts ...Option) (so StorageObserver, err error) {
	o := new(observer)
	for _, opt := range append(defaultOptions, opts...) {
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

	if o.watchEnabled {
		wech, err = o.w.Start(ctx)
		if err != nil {
			return nil, err
		}
	}

	if o.tickerEnabled {
		tech, err = o.startTicker(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
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

func (o *observer) PostStop(ctx context.Context) (err error) {
	defer o.storage.Stop(ctx)

	finalize := func() (err error) {
		err = ctx.Err()
		if err != nil && err != context.Canceled {
			return err
		}
		return nil
	}

	backup := func() error {
		metadata, err := metadata.Load(o.metadataPath)
		if err != nil {
			log.Warn("cannot read metadata of the backup files:", err)
			return err
		}

		if metadata.IsInvalid {
			log.Warn("backup skipped because the files are invalid")
			return nil
		}

		return o.backup(ctx)
	}

	t := time.Now()
	ch := make(chan struct{}, 1)
	defer close(ch)

	f := func(ctx context.Context, name string) error {
		if name == o.metadataPath {
			ch <- struct{}{}
			return nil
		}

		t = time.Now()
		return nil
	}

	o.w, err = watch.New(
		watch.WithDirs(o.dir),
		watch.WithErrGroup(o.eg),
		watch.WithOnWrite(f),
		watch.WithOnCreate(f),
	)
	if err != nil {
		return err
	}

	wctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err = o.w.Start(wctx)
	if err != nil {
		return err
	}
	defer func() {
		e := o.w.Stop(wctx)
		if e != nil {
			log.Error("an error occurred when watcher stopped:", e)
		}
	}()

	ticker := time.NewTicker(func() time.Duration {
		if o.postStopTimeout/5 < 2*time.Second {
			return o.postStopTimeout / 5
		}
		return 2 * time.Second
	}())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return finalize()
		case <-ch:
			return backup()
		case <-ticker.C:
			if time.Since(t) > o.postStopTimeout {
				return backup()
			}
		}
	}
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
				metadata, err := metadata.Load(o.metadataPath)
				if err != nil {
					log.Warn("cannot read metadata of the backup files:", err)
					ech <- err
					continue
				}

				if metadata.IsInvalid {
					log.Warn("backup skipped because the files are invalid")
					continue
				}

				err = o.requestBackup(ctx)
				if err != nil {
					ech <- err
					log.Error("failed to request backup:", err)
					err = nil
				}
			}
			if err != nil {
				log.Error("an error occurred on observer loop:", err)
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
					log.Error("an error occurred during backup:", err)
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

	if name != o.metadataPath {
		return nil
	}

	ok, err := o.isValidMetadata()
	if err != nil {
		log.Warn("cannot read metadata of the backup files:", err)
		return err
	}

	if ok {
		return o.requestBackup(ctx)
	}

	return o.terminate()
}

func (o *observer) onCreate(ctx context.Context, name string) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.onCreate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if name != o.metadataPath {
		return nil
	}

	ok, err := o.isValidMetadata()
	if err != nil {
		log.Warn("cannot read metadata of the backup files:", err)
		return err
	}

	if ok {
		return o.requestBackup(ctx)
	}

	return o.terminate()
}

func (o *observer) isValidMetadata() (bool, error) {
	metadata, err := metadata.Load(o.metadataPath)
	if err != nil {
		return false, err
	}

	return !metadata.IsInvalid, nil
}

func (*observer) terminate() error {
	log.Error("the process will be terminated because the files are invalid")

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}

	return p.Signal(syscall.SIGTERM)
}

func (o *observer) requestBackup(context.Context) error {
	select {
	case o.ch <- struct{}{}:
	default:
		log.Debug("cannot request backup: channel is full")
	}

	return nil
}

func (o *observer) backup(ctx context.Context) (err error) {
	bi := &BackupInfo{
		StorageInfo: o.storage.StorageInfo(),
	}
	bi.StartTime = time.Now()

	for _, hook := range o.hooks {
		ctx, err = hook.BeforeProcess(ctx, bi)
		if err != nil {
			return err
		}
	}

	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.backup")
	if span != nil {
		span.SetAttributes(
			attribute.String("storage_type", bi.StorageInfo.Type),
			attribute.String("bucket_name", bi.BucketName),
			attribute.String("filename", bi.Filename),
		)
	}
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	log.Infof("started to backup directory %s", o.dir)

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

		return filepath.Walk(o.dir, func(path string, fi os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(fi, path)
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(o.dir, path)
			if err != nil {
				return err
			}

			header.Name = filepath.ToSlash(rel)

			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}

			log.Debug("writing: ", path)

			if fi.IsDir() {
				return nil
			}

			data, err := file.Open(path, os.O_RDONLY, fs.ModePerm)
			if err != nil {
				return err
			}
			defer func() {
				e := data.Close()
				if e != nil {
					log.Errorf("failed to close %s: %s", path, e)
				}
			}()

			d, err := io.NewReaderWithContext(ctx, data)
			if err != nil {
				return err
			}

			_, err = io.Copy(tw, d)
			if err != nil {
				return err
			}
			return nil
		})
	}))

	prr, err := io.NewReaderWithContext(ctx, pr)
	if err != nil {
		return err
	}

	bi.Bytes, err = io.Copy(sw, prr)
	if err != nil {
		return err
	}

	wg.Wait()

	bi.EndTime = time.Now()
	for _, hook := range o.hooks {
		err = hook.AfterProcess(ctx, bi)
		if err != nil {
			return err
		}
	}

	log.Infof("finished to backup directory %s", o.dir)

	return nil
}
