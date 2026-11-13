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

package restorer

import (
	"archive/tar"
	"context"
	"io/fs"
	"reflect"
	"syscall"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type Restorer interface {
	Start(ctx context.Context) (<-chan error, error)
	PreStop(ctx context.Context) error
}

type restorer struct {
	eg             errgroup.Group
	storage        storage.Storage
	bo             backoff.Backoff
	dir            string
	backoffOpts    []backoff.Option
	backoffEnabled bool
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

		var ok bool
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err, ok = <-sech:
				if !ok {
					// closed channels are always ready; stop selecting them
					// or this loop busy-spins until ctx ends.
					sech = nil
				}
			case err, ok = <-rech:
				if !ok {
					rech = nil
				}
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

	restore := func(ctx context.Context) (any, bool, error) {
		err := r.restore(ctx)
		if err != nil {
			// Only the final result is reported on ech (below): an attempt that a
			// later retry recovers from must not fail the init container.
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
			select {
			case <-ctx.Done():
			case ech <- err:
			}
		}

		err = p.Signal(syscall.SIGTERM) // TODO: #403
		if err != nil && !errors.Is(err, os.ErrProcessDone) {
			select {
			case <-ctx.Done():
			case ech <- err:
			}
			return err
		}
		return nil
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

	// written holds the regular files this attempt created (including one that
	// may be partially written). On failure they are removed so that a retry
	// does not hit them and mistake them for pre-existing data.
	var written []string
	defer func() {
		if err == nil {
			return
		}
		for _, path := range written {
			if rerr := os.Remove(path); rerr != nil && !errors.Is(rerr, fs.ErrNotExist) {
				log.Warnf("failed to remove %s left by a failed restore: %v", path, rerr)
			}
		}
	}()

	var (
		pr io.ReadCloser
		pw io.WriteCloser
	)
	// rpw is kept unwrapped because only *io.PipeWriter has CloseWithError.
	rpr, rpw := io.Pipe()
	pr, err = io.NewReadCloserWithContext(ctx, rpr)
	if err != nil {
		log.Warn(err)
		return err
	}
	defer pr.Close()
	pw, err = io.NewWriteCloserWithContext(ctx, rpw)
	if err != nil {
		log.Warn(err)
		return err
	}
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer func() {
			// Propagate the read/copy error to the pipe reader. A plain Close()
			// surfaces a clean io.EOF, so a truncated archive would be accepted
			// as complete. On success err is nil and this behaves like Close().
			if e := rpw.CloseWithError(err); e != nil {
				log.Errorf("error on closing pipe writer: %s", e)
			}
			// restore receives the error through the pipe. It is not returned to
			// the errgroup as well: that may be the process-wide group, where a
			// failed attempt would cancel the daemon context and end the restore
			// retry loop.
			err = nil
		}()
		sr, err := r.storage.Reader(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		sr, err = io.NewReadCloserWithContext(ctx, sr)
		if err != nil {
			log.Warn(err)
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
			log.Warn(err)
			if errors.Is(err, io.EOF) {
				log.Infof("finished to read and copy data from storage to io pipe by returning io.EOF error: %v", err)
				return nil
			}
			// restore() closed its end of the pipe because it has already decided
			// the result (e.g. the target file already exists). That is a normal
			// stop, not a failure: returning it would record an error on the
			// shared errgroup and fail the daemon shutdown.
			if errors.Is(err, io.ErrClosedPipe) {
				return nil
			}
			return err
		}

		return nil
	}))

	// When another task of the errgroup fails, the errgroup cancels ctx with that
	// error as the cause, and the loop may see the cancellation first.
	// Report the cause so the failure is not masked as "context canceled".
	canceled := func(err error) error {
		if cause := context.Cause(ctx); cause != nil && errors.Is(err, context.Canceled) {
			return cause
		}
		return err
	}

	tr := tar.NewReader(pr)

	for {
		select {
		case <-ctx.Done():
			return canceled(ctx.Err())
		default:
		}

		header, err := tr.Next()
		if err != nil {
			err = canceled(err)
			log.Warn(err)
			if errors.Is(err, io.EOF) {
				log.Infof("finished to restore directory %s finished by returning io.EOF error: %v", r.dir, err)
				return nil
			}
			return err
		}

		target := file.Join(r.dir, header.Name)
		log.Debug("restoring: ", target)
		// Guard against tar-slip: file.Join already cleans the path, so a ".."
		// substring check is ineffective (clean resolves ".." away, e.g.
		// "../../etc/x" -> "/etc/x"). Require the resolved target to stay within
		// r.dir and skip any entry that escapes it.
		cleanDir := file.Join(r.dir)
		if target != cleanDir && !strings.HasPrefix(target, cleanDir+string(os.PathSeparator)) {
			log.Warnf("skipping tar entry %q: resolved path %q escapes restore dir %q (possible tar-slip)", header.Name, target, cleanDir)
			continue
		}
		switch header.Typeflag {
		case tar.TypeDir:
			err = file.MkdirAll(target, fs.ModePerm)
			if err != nil {
				log.Warn(err)
				return err
			}
		case tar.TypeReg:
			//nolint:gosec // header.Mode is a POSIX file mode within uint32 range; any truncation only affects restored file permissions, not escalation.
			_, err = file.WriteFile(ctx, target, tr, fs.FileMode(header.Mode))
			if err != nil {
				log.Warn(err)
				if errors.Is(err, errors.ErrFileAlreadyExists(target)) {
					return nil
				}
				// WriteFile only reaches open/copy for a target that did not
				// exist, so the target (possibly partial) belongs to this attempt.
				written = append(written, target)
				return err
			}
			written = append(written, target)
		}
	}
}
