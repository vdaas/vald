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
	"io"
	"os"
	"path/filepath"
	"reflect"
	"syscall"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	ctxio "github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type Restorer interface {
	Start(ctx context.Context) (<-chan error, error)
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

	log.Infof("restoring directory %s started", r.dir)
	defer log.Infof("restoring directory %s finished", r.dir)

	pr, pw := io.Pipe()
	defer pr.Close()

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer pw.Close()

		sr, err := r.storage.Reader(ctx)
		if err != nil {
			return err
		}

		sr, err = ctxio.NewReadCloserWithContext(ctx, sr)
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

		switch header.Typeflag {
		case tar.TypeDir:
			_, err = os.Stat(target)
			if err != nil {
				err = os.MkdirAll(target, 0o700)
				if err != nil {
					return err
				}
			}
		case tar.TypeReg:
			if _, err := os.Stat(target); err == nil {
				log.Warn(errors.ErrFileAlreadyExists(target))
				return nil
			} else if !os.IsNotExist(err) {
				return err
			}

			f, err := file.Open(
				target,
				os.O_CREATE|os.O_RDWR,
				os.FileMode(header.Mode),
			)
			if err != nil {
				return err
			}

			fw, err := ctxio.NewWriterWithContext(ctx, f)
			if err != nil {
				return errors.Wrap(f.Close(), err.Error())
			}

			_, err = io.Copy(fw, tr)
			if err != nil {
				return errors.Wrap(f.Close(), err.Error())
			}

			err = f.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
