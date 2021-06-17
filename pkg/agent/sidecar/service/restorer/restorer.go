//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package restorer provides restorer service
package restorer

import (
	"archive/tar"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type Restorer interface {
	Start(ctx context.Context) (<-chan error, error)
	PreStop(ctx context.Context) error
}

type restorer struct {
	dir string
	eg  errgroup.Group

	storage storage.Storage

	backoffEnabled bool
	backoffOpts    []backoff.Option
	bo             backoff.Backoff
}

func New(opts ...Option) (Restorer, error) {
	r := new(restorer)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if r.backoffEnabled {
		r.bo = backoff.New(r.backoffOpts...)
	}

	return r, nil
}

func (r *restorer) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 2)

	var sech, rech <-chan error
	var err error

	sech, err = r.storage.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	rech, err = r.startRestore(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if r.backoffEnabled {
			defer r.bo.Close()
		}

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-sech:
			case err = <-rech:
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

func (r *restorer) PreStop(ctx context.Context) error {
	if r.storage != nil {
		return r.storage.Stop(ctx)
	}
	return nil
}

func (r *restorer) startRestore(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)

	// TODO: related to #403.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return ech, err
	}

	restore := func(ctx context.Context) (interface{}, bool, error) {
		err := r.restore(ctx)
		if err != nil {
			log.Errorf("restoring failed: %s", err)
			return nil, true, err
		}

		return nil, false, nil
	}

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		if r.backoffEnabled {
			_, err = r.bo.Do(ctx, restore)
		} else {
			_, _, err = restore(ctx)
		}

		if err != nil {
			log.Errorf("couldn't restore: %s", err)
		}

		return p.Signal(syscall.SIGTERM) // TODO: #403
	}))

	return ech, nil
}

func (r *restorer) restore(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/restorer/Restorer.restore")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	log.Infof("started to restore directory %s", r.dir)

	pr, pw := io.Pipe()
	defer pr.Close()

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer pw.Close()

		sr, err := r.storage.Reader(ctx)
		if err != nil {
			return err
		}

		sr, err = io.NewReadCloserWithContext(ctx, sr)
		if err != nil {
			return err
		}
		defer func() {
			e := sr.Close()
			if e != nil {
				log.Errorf("error on closing blob-storage reader: %s", e)
			}
		}()

		_, err = io.Copy(pw, sr)
		if err != nil {
			return err
		}

		return nil
	}))

	tr := tar.NewReader(pr)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		target := filepath.Join(r.dir, header.Name)

		log.Debug("restoring: ", target)

		if strings.Contains(target, "..") {
			log.Warn(errors.ErrPathNotAllowed(target))
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			exist, fi, err := file.ExistsWithDetail(target)
			if !exist || err != nil || fi == nil || fi != nil && !fi.IsDir() {
				if exist {
					os.RemoveAll(target)
				}
				err = os.MkdirAll(target, 0o700)
				if err != nil {
					return err
				}
			}
		case tar.TypeReg:
			err = copyFile(ctx, target, tr, fs.FileMode(header.Mode))
			if err != nil {
				if errors.Is(err, errors.ErrFileAlreadyExists(target)) {
					log.Warn(err)
					return nil
				}
				return err
			}
		}
	}

	log.Infof("finished to restore directory %s finished", r.dir)

	return nil
}

func copyFile(ctx context.Context, target string, tr io.Reader, mode fs.FileMode) (err error) {
	exist, fi, err := file.ExistsWithDetail(target)
	switch {
	case err == nil, exist:
		return errors.ErrFileAlreadyExists(target)
	case err != nil && !os.IsNotExist(err):
		return err
	case fi != nil && fi.Size() != 0:
		return errors.ErrFileAlreadyExists(target)
	}

	f, err := file.Open(
		target,
		os.O_CREATE|os.O_RDWR,
		os.FileMode(mode),
	)
	if err != nil {
		return err
	}
	defer func() {
		if f != nil {
			derr := f.Close()
			if derr != nil {
				err = errors.Wrap(err, derr.Error())
			}
		}
	}()

	_, err = io.CopyWithContext(f, tr)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
