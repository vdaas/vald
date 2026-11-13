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
	"bytes"
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

// doneCountCtx counts Done() calls. A select loop evaluates ctx.Done() on every
// iteration, so the count is a deterministic measure of how often it spins.
type doneCountCtx struct {
	context.Context //nolint:containedctx // test wrapper that counts Done() calls
	n               atomic.Int64
}

func (c *doneCountCtx) Done() <-chan struct{} {
	c.n.Add(1)
	return c.Context.Done()
}

// closedStorage is a storage.Storage whose Start channel is already closed and
// whose Reader reports io.EOF, so a restore finishes immediately with no data.
type closedStorage struct {
	storage.Storage
}

func (closedStorage) Start(context.Context) (<-chan error, error) {
	ch := make(chan error)
	close(ch)
	return ch, nil
}

func (closedStorage) Reader(context.Context) (io.ReadCloser, error) {
	return nil, io.EOF
}

// TestRestorer_Start_ClosedChannelsDoNotSpin pins that the relay loop stops
// selecting the storage and restore error channels once they are closed. The
// restore channel is closed when the restore completes, which is while ctx is
// still live (startRestore then SIGTERMs the process, and ctx is cancelled only
// later). A closed channel is always ready, so without nil-ing them the loop
// busy-spins; the spin shows up as an unbounded number of ctx.Done()
// evaluations.
//
// This test must not run in parallel: startRestore signals the whole process
// with SIGTERM, which is intercepted here with signal.Notify.
func TestRestorer_Start_ClosedChannelsDoNotSpin(t *testing.T) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)
	defer signal.Stop(sigc)

	eg, egctx := errgroup.New(context.Background())
	cctx, cancel := context.WithCancel(egctx)
	defer cancel()
	ctx := &doneCountCtx{Context: cctx}

	r := &restorer{eg: eg, storage: closedStorage{}}
	if _, err := r.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	// the restore goroutine has finished (and closed the restore channel) once
	// it has delivered its SIGTERM.
	select {
	case <-sigc:
	case <-time.After(5 * time.Second):
		t.Fatal("restore did not complete")
	}

	// both channels are closed and ctx is live: the loop must be parked now.
	time.Sleep(100 * time.Millisecond)
	before := ctx.n.Load()
	time.Sleep(100 * time.Millisecond)
	if n := ctx.n.Load() - before; n > 100 {
		t.Errorf("ctx.Done() evaluated %d times in 100ms with closed channels; the relay loop is busy-spinning", n)
	}

	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
			t.Logf("eg.Wait() = %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("relay goroutine did not exit on context cancel")
	}
}

// failNStorage is a storage.Storage whose Reader fails the first `fails` calls
// and then reports io.EOF, so a later restore attempt succeeds with no data.
type failNStorage struct {
	storage.Storage
	err   error
	calls *atomic.Int64
	fails int64
}

func (s failNStorage) Reader(context.Context) (io.ReadCloser, error) {
	if s.calls.Add(1) <= s.fails {
		return nil, s.err
	}
	return nil, io.EOF
}

// TestRestorer_startRestore_reportsOnlyFinalResult pins that startRestore
// delivers only the final restore result on its error channel. A failed attempt
// that a later backoff retry recovers from must not be reported (it would make
// the init container exit non-zero), and with backoff disabled a failure must
// be reported once, not twice.
//
// This test must not run in parallel: startRestore signals the whole process
// with SIGTERM, which is intercepted here with signal.Notify.
func TestRestorer_startRestore_reportsOnlyFinalResult(t *testing.T) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)
	defer signal.Stop(sigc)

	errBoom := errors.New("transient storage failure")

	tests := []struct {
		name           string
		fails          int64
		backoffEnabled bool
		wantErrs       int
		wantCalls      int64
	}{
		{name: "failure recovered by a backoff retry is not reported", fails: 1, backoffEnabled: true, wantErrs: 0, wantCalls: 2},
		{name: "failure without backoff is reported once", fails: 1, backoffEnabled: false, wantErrs: 1, wantCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, egctx := errgroup.New(ctx)

			var calls atomic.Int64
			r := &restorer{
				eg:             eg,
				dir:            t.TempDir(),
				storage:        failNStorage{err: errBoom, calls: &calls, fails: tt.fails},
				backoffEnabled: tt.backoffEnabled,
			}
			if tt.backoffEnabled {
				r.bo = backoff.New(
					backoff.WithInitialDuration("1ms"),
					backoff.WithMaximumDuration("5ms"),
					backoff.WithJitterLimit("1ms"),
					backoff.WithRetryCount(5),
					backoff.WithBackOffTimeLimit("5s"),
					backoff.WithDisableErrorLog(),
				)
				defer r.bo.Close()
			}

			ech, err := r.startRestore(egctx)
			if err != nil {
				t.Fatalf("startRestore() error = %v, want nil", err)
			}

			var got []error
		drain:
			for {
				select {
				case e, ok := <-ech:
					if !ok {
						break drain
					}
					got = append(got, e)
				case <-ctx.Done():
					t.Fatal("restore error channel was not closed")
				}
			}
			// ech is closed after the SIGTERM is sent; wait for its delivery so it
			// cannot reach the process after signal.Stop.
			select {
			case <-sigc:
			case <-ctx.Done():
				t.Fatal("restore did not signal completion")
			}

			if len(got) != tt.wantErrs {
				t.Errorf("startRestore() reported %d errors %v, want %d", len(got), got, tt.wantErrs)
			}
			for _, e := range got {
				if !errors.Is(e, errBoom) {
					t.Errorf("startRestore() reported %v, want %v", e, errBoom)
				}
			}
			if n := calls.Load(); n != tt.wantCalls {
				t.Errorf("storage.Reader called %d times, want %d", n, tt.wantCalls)
			}

			cancel()
			done := make(chan struct{})
			go func() {
				defer close(done)
				if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
					t.Logf("eg.Wait() = %v", err)
				}
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("restore goroutines did not exit on context cancel")
			}
		})
	}
}

// Names of the files in the test tar archives.
const (
	fileA = "a.txt"
	fileB = "b.txt"
)

// errAfterReader yields data and then returns err (io.EOF when err is nil).
type errAfterReader struct {
	err  error
	data []byte
}

func (r *errAfterReader) Read(p []byte) (int, error) {
	if len(r.data) > 0 {
		n := copy(p, r.data)
		r.data = r.data[n:]
		return n, nil
	}
	if r.err == nil {
		return 0, io.EOF
	}
	return 0, r.err
}

func (*errAfterReader) Close() error { return nil }

type readerStorage struct {
	storage.Storage
	rc io.ReadCloser
}

func (s readerStorage) Reader(context.Context) (io.ReadCloser, error) {
	return s.rc, nil
}

// TestRestorer_restore_propagatesReadError pins that a read error from the blob
// reader reaches the tar consumer. The archive is cut right after a complete
// entry; with a plain pw.Close() the consumer saw a clean io.EOF at the entry
// boundary and restore() reported success for a truncated archive.
func TestRestorer_restore_propagatesReadError(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("blob read failed")

	archive := func(t *testing.T, trailer bool) []byte {
		t.Helper()
		content := []byte("hello")
		buf := new(bytes.Buffer)
		tw := tar.NewWriter(buf)
		if err := tw.WriteHeader(&tar.Header{
			Name:     "f.txt",
			Mode:     0o644,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(content); err != nil {
			t.Fatal(err)
		}
		if trailer {
			// Close writes the end-of-archive marker.
			if err := tw.Close(); err != nil {
				t.Fatal(err)
			}
		} else if err := tw.Flush(); err != nil {
			t.Fatal(err)
		}
		return buf.Bytes()
	}

	tests := []struct {
		readErr error
		wantErr error
		name    string
		trailer bool
	}{
		{name: "read error after a complete entry is propagated", trailer: false, readErr: errBoom, wantErr: errBoom},
		{name: "complete archive without error succeeds", trailer: true, readErr: nil, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, egctx := errgroup.New(ctx)

			r := &restorer{
				eg:  eg,
				dir: t.TempDir(),
				storage: readerStorage{
					rc: &errAfterReader{data: archive(t, tt.trailer), err: tt.readErr},
				},
			}

			// restore is bounded by ctx: both the pipe reads and the tar loop honour it.
			err := r.restore(egctx)
			if ctx.Err() != nil {
				t.Fatalf("restore did not finish in time: %v", ctx.Err())
			}
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("restore() error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("restore() error = %v, want %v", err, tt.wantErr)
			}

			// The copy goroutine hands its error to restore() through the pipe
			// only. Returning it to the errgroup, which is the process-wide group
			// in the init container, cancelled the daemon context and ended the
			// restore retry loop after a single failed attempt.
			waited := make(chan error, 1)
			go func() { waited <- eg.Wait() }()
			select {
			case err := <-waited:
				if err != nil {
					t.Errorf("eg.Wait() = %v, want nil", err)
				}
			case <-ctx.Done():
				t.Fatalf("eg.Wait did not return: %v", ctx.Err())
			}
			cancel()
		})
	}
}

// tarEntry is one regular file of a test archive; content may be shorter than
// size to cut the archive in the middle of the entry.
type tarEntry struct {
	name    string
	content []byte
	size    int64
}

func buildTar(t *testing.T, entries []tarEntry, trailer bool) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for _, e := range entries {
		if err := tw.WriteHeader(&tar.Header{
			Name:     e.name,
			Mode:     0o644,
			Size:     e.size,
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(e.content); err != nil {
			t.Fatal(err)
		}
	}
	if trailer {
		if err := tw.Close(); err != nil {
			t.Fatal(err)
		}
	} else if n := len(entries); n > 0 && int64(len(entries[n-1].content)) == entries[n-1].size {
		// Pad the last complete entry; a cut entry is left as is (its bytes
		// are already in buf, and Flush would reject the missing remainder).
		if err := tw.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	return buf.Bytes()
}

func runRestore(t *testing.T, dir string, data []byte, readErr error) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)
	r := &restorer{
		eg:      eg,
		dir:     dir,
		storage: readerStorage{rc: &errAfterReader{data: data, err: readErr}},
	}
	err := r.restore(egctx)
	if ctx.Err() != nil {
		t.Fatalf("restore did not finish in time: %v", ctx.Err())
	}
	cancel()
	if werr := eg.Wait(); werr != nil && !errors.Is(werr, context.Canceled) {
		t.Logf("eg.Wait() = %v", werr)
	}
	return err
}

// TestRestorer_restore_failedAttemptCleansUp pins that a failed attempt removes
// the files it wrote, including a partially written one. Previously the
// leftovers stayed in the directory and the retry stopped at the first one
// (ErrFileAlreadyExists treated as "already restored"), reporting success for
// an incomplete restore.
func TestRestorer_restore_failedAttemptCleansUp(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("blob read failed")
	full1 := []byte("first-file")
	full2 := []byte("second-file-content")

	tests := []struct {
		name    string
		partial []tarEntry
	}{
		{
			name:    "complete entry then read error",
			partial: []tarEntry{{name: fileA, content: full1, size: int64(len(full1))}},
		},
		{
			name: "complete entry then error inside the next entry",
			partial: []tarEntry{
				{name: fileA, content: full1, size: int64(len(full1))},
				{name: fileB, content: full2[:4], size: int64(len(full2))},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()

			// The error may surface as errBoom or as the cancellation it causes,
			// depending on which reader sees it first; only failure matters here.
			if err := runRestore(t, dir, buildTar(t, tt.partial, false), errBoom); err == nil {
				t.Fatal("first restore() succeeded, want an error")
			}
			for _, name := range []string{fileA, fileB} {
				if file.Exists(file.Join(dir, name)) {
					t.Errorf("%s left behind by the failed attempt", name)
				}
			}

			complete := buildTar(t, []tarEntry{
				{name: fileA, content: full1, size: int64(len(full1))},
				{name: fileB, content: full2, size: int64(len(full2))},
			}, true)
			if err := runRestore(t, dir, complete, nil); err != nil {
				t.Fatalf("retry restore() error = %v, want nil", err)
			}
			for name, want := range map[string][]byte{fileA: full1, fileB: full2} {
				got, err := file.ReadFile(file.Join(dir, name))
				if err != nil {
					t.Fatalf("read %s after retry: %v", name, err)
				}
				if !bytes.Equal(got, want) {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
		})
	}
}

// TestRestorer_restore_keepsPreExistingData pins that real pre-existing data is
// still left untouched: restore stops at the existing file and reports success.
func TestRestorer_restore_keepsPreExistingData(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	existing := []byte("local-data")
	if _, err := file.WriteFile(t.Context(), file.Join(dir, fileA), bytes.NewReader(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	remote := []byte("remote")
	data := buildTar(t, []tarEntry{{name: fileA, content: remote, size: int64(len(remote))}}, true)
	if err := runRestore(t, dir, data, nil); err != nil {
		t.Fatalf("restore() error = %v, want nil", err)
	}
	got, err := file.ReadFile(file.Join(dir, fileA))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, existing) {
		t.Errorf("a.txt = %q, want pre-existing %q", got, existing)
	}
}

// TestRestorer_restore_existingFileStopsCopyCleanly pins that restore() stopping
// early on an already restored file is not an error for the copy goroutine. It
// closes the pipe while the goroutine is still writing the rest of the archive;
// returning the resulting io.ErrClosedPipe recorded it on the errgroup, which is
// the process-wide group in the init container, so every restart over existing
// index files exited non-zero.
func TestRestorer_restore_existingFileStopsCopyCleanly(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)

	dir := t.TempDir()
	if _, err := file.OverWriteFile(ctx, file.Join(dir, "f.txt"), bytes.NewReader([]byte("old")), 0o600); err != nil {
		t.Fatal(err)
	}
	// The large second entry keeps the copy goroutine blocked on the pipe after
	// restore() has stopped at the first one.
	big := make([]byte, 256<<10)
	data := buildTar(t, []tarEntry{
		{name: "f.txt", content: []byte("new"), size: 3},
		{name: "g.bin", content: big, size: int64(len(big))},
	}, true)

	r := &restorer{
		eg:      eg,
		dir:     dir,
		storage: readerStorage{rc: &errAfterReader{data: data}},
	}
	if err := r.restore(egctx); err != nil {
		t.Fatalf("restore() error = %v, want nil", err)
	}

	waited := make(chan error, 1)
	go func() { waited <- eg.Wait() }()
	select {
	case err := <-waited:
		if err != nil {
			t.Errorf("eg.Wait() = %v, want nil: an early stop of restore is not a copy failure", err)
		}
	case <-ctx.Done():
		t.Fatalf("eg.Wait did not return: %v", ctx.Err())
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Restorer
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Restorer, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Restorer, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_restorer_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg             errgroup.Group
// 		storage        storage.Storage
// 		bo             backoff.Backoff
// 		dir            string
// 		backoffOpts    []backoff.Option
// 		backoffEnabled bool
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			r := &restorer{
// 				eg:             test.fields.eg,
// 				storage:        test.fields.storage,
// 				bo:             test.fields.bo,
// 				dir:            test.fields.dir,
// 				backoffOpts:    test.fields.backoffOpts,
// 				backoffEnabled: test.fields.backoffEnabled,
// 			}
//
// 			got, err := r.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_restorer_PreStop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg             errgroup.Group
// 		storage        storage.Storage
// 		bo             backoff.Backoff
// 		dir            string
// 		backoffOpts    []backoff.Option
// 		backoffEnabled bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			r := &restorer{
// 				eg:             test.fields.eg,
// 				storage:        test.fields.storage,
// 				bo:             test.fields.bo,
// 				dir:            test.fields.dir,
// 				backoffOpts:    test.fields.backoffOpts,
// 				backoffEnabled: test.fields.backoffEnabled,
// 			}
//
// 			err := r.PreStop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_restorer_startRestore(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg             errgroup.Group
// 		storage        storage.Storage
// 		bo             backoff.Backoff
// 		dir            string
// 		backoffOpts    []backoff.Option
// 		backoffEnabled bool
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			r := &restorer{
// 				eg:             test.fields.eg,
// 				storage:        test.fields.storage,
// 				bo:             test.fields.bo,
// 				dir:            test.fields.dir,
// 				backoffOpts:    test.fields.backoffOpts,
// 				backoffEnabled: test.fields.backoffEnabled,
// 			}
//
// 			got, err := r.startRestore(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_restorer_restore(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg             errgroup.Group
// 		storage        storage.Storage
// 		bo             backoff.Backoff
// 		dir            string
// 		backoffOpts    []backoff.Option
// 		backoffEnabled bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           storage:nil,
// 		           bo:nil,
// 		           dir:"",
// 		           backoffOpts:nil,
// 		           backoffEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			r := &restorer{
// 				eg:             test.fields.eg,
// 				storage:        test.fields.storage,
// 				bo:             test.fields.bo,
// 				dir:            test.fields.dir,
// 				backoffOpts:    test.fields.backoffOpts,
// 				backoffEnabled: test.fields.backoffEnabled,
// 			}
//
// 			err := r.restore(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
