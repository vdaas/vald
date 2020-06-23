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

	backoffOpts []backoff.Option
}

func New(opts ...Option) (Restorer, error) {
	r := new(restorer)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
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

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		b := backoff.New(r.backoffOpts...)
		defer b.Close()

		_, err = b.Do(ctx, func() (interface{}, error) {
			err := r.restore(ctx)
			if err != nil {
				log.Errorf("restoring failed: %s", err)

				if errors.IsErrBlobNoSuchBucket(err) ||
					errors.IsErrBlobNoSuchKey(err) {
					return nil, nil
				}

				return nil, err
			}

			return nil, nil
		})
		if err != nil {
			return p.Signal(syscall.SIGKILL) // TODO: #403
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

	pr, pw := io.Pipe()
	defer pr.Close()

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer pw.Close()

		sr, err := r.storage.Reader(ctx)
		if err != nil {
			return err
		}

		_, err = io.Copy(pw, sr)
		if err != nil {
			return err
		}

		return nil
	}))

	tr := tar.NewReader(pr)

	for {
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
				err = os.MkdirAll(target, 0700)
				if err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(
				target,
				os.O_CREATE|os.O_RDWR,
				os.FileMode(header.Mode),
			)
			if err != nil {
				return err
			}

			_, err = io.Copy(f, tr)
			if err != nil {
				return err
			}

			err = f.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
