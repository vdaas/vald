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

// Package service manages the main logic of server.
package service

import (
	"context"
	"encoding/gob"
	"fmt"
	"io/fs"
	"math"
	"path/filepath"
	"slices"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	kvald "github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	testdata "github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/internal/test/mock/k8s"
	"github.com/vdaas/vald/internal/test/testify"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

var defaultConfig = config.NGT{
	Dimension:           100,
	DistanceType:        "l2",
	ObjectType:          "float",
	BulkInsertChunkSize: 10,
	CreationEdgeSize:    20,
	SearchEdgeSize:      10,
	EnableProactiveGC:   false,
	EnableCopyOnWrite:   false,
	KVSDB: &config.KVSDB{
		Concurrency: 10,
	},
	BrokenIndexHistoryLimit: 1,
	ErrorBufferLimit:        100,
}

type index struct {
	uuid string
	vec  []float32
}

func TestNew(t *testing.T) {
	type args struct {
		cfg  *config.NGT
		opts []Option
	}
	type want struct {
		err error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			tmpDir := t.TempDir()
			brokenDir := filepath.Join(tmpDir, brokenIndexDirName)
			return test{
				name: "New creates `origin` and `broken` directory with default options",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					dirs, err := file.ListInDir(tmpDir)
					if err != nil {
						return err
					}

					// extract folder name from dir path into a map
					dirSet := make(map[string]struct{}, len(dirs))
					for _, dir := range dirs {
						// extract folder name from dir path
						dir = dir[len(tmpDir)+1:]
						dirSet[dir] = struct{}{}
					}

					// check if the dirs set contains folder names origin, backup and broken.
					if _, ok := dirSet[originIndexDirName]; !ok {
						return fmt.Errorf("failed to create origin dir")
					}
					if _, ok := dirSet[brokenIndexDirName]; !ok {
						return fmt.Errorf("failed to create broken dir")
					}

					// check if the broken index directory is empty
					files, err := file.ListInDir(brokenDir)
					if err != nil {
						return err
					}
					if len(files) != 0 {
						return fmt.Errorf("broken index directory is not empty")
					}
					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			originDir := filepath.Join(tmpDir, originIndexDirName)
			testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: "New migrates index files into `origin`",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
					// copy testdata index files into tmpDir which is a old index directory
					// this should be moved to origin directory by the migration process
					file.CopyDir(t.Context(), testIndexDir, tmpDir)
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					files, err := file.ListInDir(tmpDir)
					if err != nil {
						return err
					}

					// extract folder name from dir path into a map
					dirSet := make(map[string]struct{}, len(files))
					for _, dir := range files {
						// extract folder name from dir path
						dirSet[filepath.Base(dir)] = struct{}{}
					}

					// check if the dirs set contains folder names origin, backup and broken.
					if _, ok := dirSet[originIndexDirName]; !ok {
						return fmt.Errorf("failed to create origin dir")
					}
					if _, ok := dirSet[brokenIndexDirName]; !ok {
						return fmt.Errorf("failed to create broken dir")
					}

					// check if the origin index directory has index files
					files, err = file.ListInDir(originDir)
					if err != nil {
						return err
					}
					if len(files) == 0 {
						return fmt.Errorf("migration failed to move index files")
					}
					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			originDir := filepath.Join(tmpDir, originIndexDirName)
			testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: "New migrates does not migrate index files if origin directory already exists",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
					// copy testdata index files into tmpDir which is a old index directory
					err := file.CopyDir(t.Context(), testIndexDir, tmpDir)
					if err != nil {
						t.Errorf("failed to copy testdata index files: %v", err)
					}

					// copy testdata index files into tmpDir which is a old index directory
					err = file.MkdirAll(originDir, fs.ModePerm)
					if err != nil {
						t.Errorf("failed to create origin directory: %v", err)
					}
					err = file.CopyDir(t.Context(), testIndexDir, originDir)
					if err != nil {
						t.Errorf("failed to copy testdata index files: %v", err)
					}
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					files, err := file.ListInDir(tmpDir)
					if err != nil {
						return err
					}

					metadataExists := false
					for _, file := range files {
						if filepath.Base(file) == "metadata.json" {
							metadataExists = true
						}
					}
					if !metadataExists {
						return fmt.Errorf("migration should not happen")
					}
					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			brokenDir := filepath.Join(tmpDir, brokenIndexDirName)
			config := defaultConfig
			config.EnableCopyOnWrite = true
			return test{
				name: "New creates `origin`, `backup` and `broken` directory with CoW enabled",
				args: args{
					cfg: &config,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					dirs, err := file.ListInDir(tmpDir)
					if err != nil {
						return err
					}

					// extract folder name from dir path into a map
					dirSet := make(map[string]struct{}, len(dirs))
					for _, dir := range dirs {
						// extract folder name from dir path
						dir = dir[len(tmpDir)+1:]
						dirSet[dir] = struct{}{}
					}

					// check if the dirs set contains folder names origin, backup and broken.
					if _, ok := dirSet[originIndexDirName]; !ok {
						return fmt.Errorf("failed to create origin dir")
					}
					if _, ok := dirSet[oldIndexDirName]; !ok {
						return fmt.Errorf("failed to create backup dir")
					}
					if _, ok := dirSet[brokenIndexDirName]; !ok {
						return fmt.Errorf("failed to create broken dir")
					}

					// check if the broken index directory is empty
					files, err := file.ListInDir(brokenDir)
					if err != nil {
						return err
					}
					if len(files) != 0 {
						return fmt.Errorf("broken index directory is not empty")
					}

					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			originDir := filepath.Join(tmpDir, originIndexDirName)
			brokenDir := filepath.Join(tmpDir, brokenIndexDirName)
			testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			config := defaultConfig
			config.BrokenIndexHistoryLimit = 1
			return test{
				name: "New succeeds to backup broken index",
				args: args{
					cfg: &config,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
					if err := file.MkdirAll(originDir, fs.ModePerm); err != nil {
						t.Errorf("failed to create origin dir: %v", err)
					}
					file.CopyDir(t.Context(), testIndexDir, originDir)
					// remove metadata.json to make it broken
					if err := os.Remove(filepath.Join(originDir, "metadata.json")); err != nil {
						t.Errorf("failed to remove index file: %v", err)
					}
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					files, err := file.ListInDir(brokenDir)
					if err != nil {
						return err
					}
					if len(files) != 1 {
						return fmt.Errorf("only one generation should be in broken dir")
					}

					broken, err := file.ListInDir(files[0])
					if err != nil {
						return err
					}
					if len(broken) == 0 {
						return fmt.Errorf("failed to move broken index files")
					}
					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			originDir := filepath.Join(tmpDir, originIndexDirName)
			brokenDir := filepath.Join(tmpDir, brokenIndexDirName)
			testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			config := defaultConfig
			config.BrokenIndexHistoryLimit = 1
			return test{
				name: "New succeeds to rotate broken index backup when the number of generations exceeds the limit",
				args: args{
					cfg: &config,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
					if err := file.MkdirAll(originDir, fs.ModePerm); err != nil {
						t.Errorf("failed to create origin dir: %v", err)
					}
					file.CopyDir(t.Context(), testIndexDir, originDir)
					// remove metadata.json to make it broken
					if err := os.Remove(filepath.Join(originDir, "metadata.json")); err != nil {
						t.Errorf("failed to remove index file: %v", err)
					}

					if err := file.MkdirAll(brokenDir, fs.ModePerm); err != nil {
						t.Errorf("failed to create broken dir: %v", err)
					}
					gen1 := filepath.Join(brokenDir, fmt.Sprint(time.Now().UnixNano()))
					if err := file.MkdirAll(gen1, fs.ModePerm); err != nil {
						t.Errorf("failed to create gen1 dir: %v", err)
					}
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					files, err := file.ListInDir(brokenDir)
					if err != nil {
						return err
					}
					if len(files) != 1 {
						return fmt.Errorf("only one generation should be in broken dir")
					}

					broken, err := file.ListInDir(files[0])
					if err != nil {
						return err
					}
					if len(broken) == 0 {
						return fmt.Errorf("failed to move broken index files")
					}
					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			originDir := filepath.Join(tmpDir, originIndexDirName)
			brokenDir := filepath.Join(tmpDir, brokenIndexDirName)
			testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			config := defaultConfig
			config.BrokenIndexHistoryLimit = 0
			return test{
				name: "New does not backup when history limit is 0",
				args: args{
					cfg: &config,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
					if err := file.MkdirAll(originDir, fs.ModePerm); err != nil {
						t.Errorf("failed to create origin dir: %v", err)
					}
					file.CopyDir(t.Context(), testIndexDir, originDir)
					// remove metadata.json to make it broken
					if err := os.Remove(filepath.Join(originDir, "metadata.json")); err != nil {
						t.Errorf("failed to remove index file: %v", err)
					}
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					files, err := file.ListInDir(brokenDir)
					if err != nil {
						return err
					}
					if len(files) != 0 {
						return fmt.Errorf("backup should not happen")
					}
					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			indexDir := filepath.Join(tmpDir, "foo") // this does not exists when this test starts
			brokenDir := filepath.Join(indexDir, brokenIndexDirName)
			config := defaultConfig
			return test{
				name: "New creates `origin` and `backup` directory even when index path does not exist",
				args: args{
					cfg: &config,
					opts: []Option{
						WithIndexPath(indexDir),
					},
				},
				want: want{
					err: nil,
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					dirs, err := file.ListInDir(indexDir)
					if err != nil {
						return err
					}

					// extract folder name from dir path into a map
					dirSet := make(map[string]struct{}, len(dirs))
					for _, dir := range dirs {
						// extract folder name from dir path
						dirSet[filepath.Base(dir)] = struct{}{}
					}

					// check if the dirs set contains folder names origin, backup and broken.
					if _, ok := dirSet[originIndexDirName]; !ok {
						return fmt.Errorf("failed to create origin dir")
					}
					if _, ok := dirSet[brokenIndexDirName]; !ok {
						return fmt.Errorf("failed to create broken dir")
					}

					// check if the broken index directory is empty
					files, err := file.ListInDir(brokenDir)
					if err != nil {
						return err
					}
					if len(files) != 0 {
						return fmt.Errorf("broken index directory is not empty")
					}

					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			originDir := filepath.Join(tmpDir, originIndexDirName)
			backupDir := filepath.Join(tmpDir, oldIndexDirName)
			brokenDir := filepath.Join(tmpDir, brokenIndexDirName)
			testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			config := defaultConfig
			config.BrokenIndexHistoryLimit = 1
			config.EnableCopyOnWrite = true
			return test{
				name: "New backup broken index when CoW is enabled and failed to load primary index",
				args: args{
					cfg: &config,
					opts: []Option{
						WithIndexPath(tmpDir),
					},
				},
				want: want{
					err: nil,
				},
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
					if err := file.MkdirAll(originDir, fs.ModePerm); err != nil {
						t.Errorf("failed to create origin dir: %v", err)
					}
					if err := file.CopyDir(t.Context(), testIndexDir, originDir); err != nil {
						t.Errorf("failed to copy test index: %v", err)
					}
					// remove metadata.json to make it broken
					if err := os.Remove(filepath.Join(originDir, "metadata.json")); err != nil {
						t.Errorf("failed to remove index file: %v", err)
					}

					if err := file.MkdirAll(backupDir, fs.ModePerm); err != nil {
						t.Errorf("failed to create backup dir: %v", err)
					}
					if err := file.CopyDir(t.Context(), testIndexDir, backupDir); err != nil {
						t.Errorf("failed to copy test index: %v", err)
					}
				},
				checkFunc: func(w want, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					files, err := file.ListInDir(brokenDir)
					if err != nil {
						return err
					}
					if len(files) != 1 {
						return fmt.Errorf("only one generation should be in broken dir but there's %v", len(files))
					}

					broken, err := file.ListInDir(files[0])
					if err != nil {
						return err
					}
					if len(broken) == 0 {
						return fmt.Errorf("failed to move broken index files")
					}

					files, err = file.ListInDir(originDir)
					if err != nil {
						return err
					}
					if len(files) != 0 {
						return fmt.Errorf("failed to move origin index files to broken directory")
					}

					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			_, err := New(test.args.cfg, test.args.opts...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_needsBackup(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		need bool
	}
	type test struct {
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
		want       want
	}
	defaultCheckFunc := func(w want, need bool) error {
		if need != w.need {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", need, w.need)
		}
		return nil
	}
	tests := []test{
		func() test {
			tmpDir := t.TempDir()
			validIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: "returns false when it's an initial state",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: false,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(t.Context(), validIndexDir, tmpDir); err != nil {
						t.Errorf("failed to copy index files: %v", err)
					}

					// remove .json and .kvsdb files to simulate an initial state
					files, err := file.ListInDir(tmpDir)
					if err != nil {
						t.Errorf("failed to list index files: %v", err)
					}
					for _, file := range files {
						if strings.HasSuffix(file, ".json") || strings.HasSuffix(file, ".kvsdb") {
							if err := os.Remove(file); err != nil {
								t.Errorf("failed to remove index file: %v", err)
							}
						}
					}
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			validIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: "returns true when there's index files but no metadata.json",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: true,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(t.Context(), validIndexDir, tmpDir); err != nil {
						t.Errorf("failed to copy index files: %v", err)
					}

					// remove metadata.json
					metafile := filepath.Join(tmpDir, "metadata.json")
					if err := os.Remove(metafile); err != nil {
						t.Errorf("failed to remove metadata.json: %v", err)
					}
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			validIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: "returns true when meta.IsInvalid is true",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: true,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(t.Context(), validIndexDir, tmpDir); err != nil {
						t.Errorf("failed to copy index files: %v", err)
					}

					// change IsInvalid in metadata.json
					metafile := filepath.Join(tmpDir, "metadata.json")
					meta, err := metadata.Load(metafile)
					if err != nil {
						t.Errorf("failed to load metadata.json: %v", err)
					}
					meta.IsInvalid = true
					meta.NGT.IndexCount = 0
					if err := metadata.Store(metafile, meta); err != nil {
						t.Errorf("failed to store metadata.json: %v", err)
					}
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			validIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: "returns true when meta.IsInvalid is true",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: true,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(t.Context(), validIndexDir, tmpDir); err != nil {
						t.Errorf("failed to copy index files: %v", err)
					}

					// change NGT.IndexCount in metadata.json
					metafile := filepath.Join(tmpDir, "metadata.json")
					meta, err := metadata.Load(metafile)
					if err != nil {
						t.Errorf("failed to load metadata.json: %v", err)
					}
					meta.IsInvalid = false
					meta.NGT.IndexCount = 100
					if err := metadata.Store(metafile, meta); err != nil {
						t.Errorf("failed to store metadata.json: %v", err)
					}
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			validIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: "returns false when NGT.IndexCount is 0",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: false,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(t.Context(), validIndexDir, tmpDir); err != nil {
						t.Errorf("failed to copy index files: %v", err)
					}

					// change NGT.IndexCount in metadata.json
					metafile := filepath.Join(tmpDir, "metadata.json")
					meta, err := metadata.Load(metafile)
					if err != nil {
						t.Errorf("failed to load metadata.json: %v", err)
					}
					meta.IsInvalid = false
					meta.NGT.IndexCount = 0
					if err := metadata.Store(metafile, meta); err != nil {
						t.Errorf("failed to store metadata.json: %v", err)
					}
				},
			}
		}(),
	}
	// An unreadable metadata.json next to index data (crash while writing it) must be
	// treated as a broken index so that rebuild backs it up instead of deleting it.
	for _, corrupt := range []struct {
		modify func(t *testing.T, metafile string)
		name   string
	}{
		{
			name: "returns true when metadata.json is empty",
			modify: func(t *testing.T, metafile string) {
				t.Helper()
				if err := os.Truncate(metafile, 0); err != nil {
					t.Errorf("failed to truncate metadata.json: %v", err)
				}
			},
		},
		{
			name: "returns true when metadata.json is truncated",
			modify: func(t *testing.T, metafile string) {
				t.Helper()
				b, err := os.ReadFile(metafile)
				if err != nil {
					t.Errorf("failed to read metadata.json: %v", err)
					return
				}
				if len(b) < 2 {
					t.Errorf("metadata.json is too short to truncate: %d bytes", len(b))
					return
				}
				if err := os.Truncate(metafile, int64(len(b)/2)); err != nil {
					t.Errorf("failed to truncate metadata.json: %v", err)
				}
			},
		},
		{
			name: "returns true when metadata.json is null",
			modify: func(t *testing.T, metafile string) {
				t.Helper()
				if err := os.WriteFile(metafile, []byte("null"), fs.ModePerm); err != nil {
					t.Errorf("failed to write metadata.json: %v", err)
				}
			},
		},
	} {
		tests = append(tests, func() test {
			tmpDir := t.TempDir()
			validIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
			return test{
				name: corrupt.name,
				args: args{
					path: tmpDir,
				},
				want: want{
					need: true,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(t.Context(), validIndexDir, tmpDir); err != nil {
						t.Errorf("failed to copy index files: %v", err)
					}
					corrupt.modify(t, filepath.Join(tmpDir, metadata.AgentMetadataFileName))
				},
			}
		}())
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			need := needsBackup(test.args.path)
			if err := checkFunc(test.want, need); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNGT_GetObject(t *testing.T) {
	t.Parallel()

	type test struct {
		testfunc func(t *testing.T)
		name     string
	}

	tests := []test{
		{
			name:     "returns vector and timestamp when vector is found in vqueue",
			testfunc: testReturnFromVq,
		},
		{
			name:     "returns vector and timestamp when vector is found in kvs",
			testfunc: testReturnFromKvs,
		},
		{
			name:     "returns error when vector is not found in vector queue or kvs",
			testfunc: testNotFoundInBothVqAndKvs,
		},
		{
			name:     "returns error when vector is not found in vq found in kvs but also in delete queue",
			testfunc: testFoundInBothIvqAndDvq,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			test.testfunc(tt)
		})
	}
}

func testReturnFromVq(t *testing.T) {
	ngt, err := New(&defaultConfig)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	err = ngt.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, now)
	require.NoError(t, err)

	vec, ts, err := ngt.GetObject("test-uuid")
	require.NoError(t, err)
	require.Equal(t, []float32{1.0, 2.0, 3.0}, vec)
	require.Equal(t, now, ts)
}

func testReturnFromKvs(t *testing.T) {
	config := defaultConfig
	config.Dimension = 3
	ngt, err := New(&config)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	err = ngt.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, now)
	require.NoError(t, err)

	err = ngt.CreateIndex(t.Context(), 10)
	require.NoError(t, err)

	buflen := ngt.InsertVQueueBufferLen()
	require.Equal(t, buflen, uint64(0))

	vec, ts, err := ngt.GetObject("test-uuid")
	require.NoError(t, err)
	require.Equal(t, []float32{1.0, 2.0, 3.0}, vec)
	require.Equal(t, now, ts)
}

func testNotFoundInBothVqAndKvs(t *testing.T) {
	ngt, err := New(&defaultConfig)
	require.NoError(t, err)

	_, _, err = ngt.GetObject("test-uuid")
	want := errors.ErrObjectIDNotFound("test-uuid")
	require.Equal(t, err.Error(), want.Error())
}

func testFoundInBothIvqAndDvq(t *testing.T) {
	config := defaultConfig
	config.Dimension = 3
	ngt, err := New(&config)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	err = ngt.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, now)
	require.NoError(t, err)

	err = ngt.CreateIndex(t.Context(), 10)
	require.NoError(t, err)

	buflen := ngt.InsertVQueueBufferLen()
	require.Equal(t, buflen, uint64(0))

	err = ngt.Delete("test-uuid")
	require.NoError(t, err)

	_, _, err = ngt.GetObject("test-uuid")
	want := errors.ErrObjectIDNotFound("test-uuid")
	require.Equal(t, err.Error(), want.Error())
}

func Test_ngt_CreateIndex(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg  *config.NGT
		opts []Option
	}
	type test struct {
		want error
		name string
		args args
	}

	setup := func(t *testing.T) string {
		tmpDir := t.TempDir()
		testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
		err := file.CopyDir(t.Context(), testIndexDir, tmpDir)
		require.NoError(t, err)

		return tmpDir
	}

	tests := []test{
		func() test {
			tmpDir := setup(t)
			return test{
				name: "CreateIndex returns ErrUncommittedIndexNotFound when there is nothing uncommitted",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
						WithIsReadReplica(false),
					},
				},
				want: errors.ErrUncommittedIndexNotFound,
			}
		}(),
		func() test {
			tmpDir := setup(t)
			return test{
				name: "CreateIndex returns ErrWriteOperationToReadReplica when try to create index to read replica",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
						WithIsReadReplica(true),
					},
				},
				want: errors.ErrWriteOperationToReadReplica,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ngt, err := New(test.args.cfg, test.args.opts...)
			require.NoError(tt, err)

			err = ngt.CreateIndex(t.Context(), test.args.cfg.DefaultPoolSize)
			require.Equal(tt, test.want, err)
		})
	}
}

func Test_ngt_SaveIndex(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg  *config.NGT
		opts []Option
	}
	type test struct {
		want error
		name string
		args args
	}

	setup := func(t *testing.T) string {
		tmpDir := t.TempDir()
		testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
		err := file.CopyDir(t.Context(), testIndexDir, tmpDir)
		require.NoError(t, err)

		return tmpDir
	}

	tests := []test{
		func() test {
			tmpDir := setup(t)
			return test{
				name: "CreateIndex successes when there is nothing to save",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
						WithIsReadReplica(false),
					},
				},
				want: nil,
			}
		}(),
		func() test {
			tmpDir := setup(t)
			return test{
				name: "SaveIndex returns ErrWriteOperationToReadReplica when try to save index to read replica",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
						WithIsReadReplica(true),
					},
				},
				want: errors.ErrWriteOperationToReadReplica,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ngt, err := New(test.args.cfg, test.args.opts...)
			require.NoError(tt, err)

			err = ngt.SaveIndex(t.Context())
			require.Equal(tt, test.want, err)
		})
	}
}

func Test_ngt_Close(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg  *config.NGT
		opts []Option
	}
	type test struct {
		want error
		name string
		args args
	}

	setup := func(t *testing.T) string {
		tmpDir := t.TempDir()
		testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
		err := file.CopyDir(t.Context(), testIndexDir, tmpDir)
		require.NoError(t, err)

		return tmpDir
	}

	tests := []test{
		func() test {
			tmpDir := setup(t)
			return test{
				name: "Close successes when it is a read replica because it skips all the Close operations",
				args: args{
					cfg: &defaultConfig,
					opts: []Option{
						WithIndexPath(tmpDir),
						WithIsReadReplica(true),
					},
				},
				want: nil,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ngt, err := New(test.args.cfg, test.args.opts...)
			require.NoError(tt, err)

			err = ngt.Close(t.Context())
			require.Equal(tt, test.want, err)
		})
	}
}

func Test_ngt_Close_waitsForInflightCreateIndex(t *testing.T) {
	t.Parallel()
	cfg := defaultConfig
	cfg.Dimension = 3

	nn, err := New(&cfg, WithIndexPath(t.TempDir()))
	require.NoError(t, err)
	n, ok := nn.(*ngt)
	require.True(t, ok)

	require.NoError(t, n.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, time.Now().UnixNano()))
	require.NoError(t, n.CreateIndex(t.Context(), 10))
	require.NoError(t, n.SaveIndex(t.Context()))
	require.Equal(t, n.NumberOfCreateIndexExecution(), n.lastNumberOfCreateIndexExecution())

	// Simulate an auto-index CreateIndex in its graph build phase: the batch is already popped
	// from the vqueue and nocie is not bumped yet.
	n.indexing.Store(true)

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- n.Close(ctx)
	}()
	select {
	case err := <-done:
		t.Fatalf("Close returned while CreateIndex was still in progress: %v", err)
	case <-time.After(300 * time.Millisecond):
	}

	// The in-flight CreateIndex finishes.
	n.nocie.Add(1)
	n.indexing.Store(false)

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		cancel()
		t.Fatal("Close did not return after the in-flight CreateIndex finished")
	}
	require.Equal(t, n.NumberOfCreateIndexExecution(), n.lastNumberOfCreateIndexExecution(),
		"Close must save the generation produced by the in-flight CreateIndex")
}

func TestCreateIndex_exportsIndexInfoOutsideIndexingWindow(t *testing.T) {
	t.Parallel()
	cfg := defaultConfig
	cfg.Dimension = 3
	cfg.EnableExportIndexInfoToK8s = true
	cfg.PodName = "test-pod"

	var (
		n                    NGT
		exported             atomic.Bool
		indexingDuringExport atomic.Bool
	)
	mock := &k8s.PatcherMock{}
	mock.On("ApplyPodAnnotations",
		testify.Anything,
		testify.Anything,
		testify.Anything,
		testify.Anything,
	).Run(func(testify.Arguments) {
		exported.Store(true)
		indexingDuringExport.Store(n.IsIndexing())
	}).Return(nil)

	var err error
	n, err = New(&cfg, WithPatcher(mock))
	require.NoError(t, err)
	require.NoError(t, n.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, time.Now().UnixNano()))
	require.NoError(t, n.CreateIndex(t.Context(), 10))

	require.True(t, exported.Load(), "index info must be exported after CreateIndex")
	require.False(t, indexingDuringExport.Load(), "search must not be blocked by the k8s export")
}

func TestExportIndexInfo(t *testing.T) {
	t.Parallel()
	config := defaultConfig
	config.Dimension = 3
	config.EnableExportIndexInfoToK8s = true
	config.PodName = "test-pod"

	type test struct {
		testfunc func(t *testing.T)
		name     string
	}

	tests := []test{
		{
			name: "export after create index one vector",
			testfunc: func(t *testing.T) {
				mock := &k8s.PatcherMock{}
				mock.On("ApplyPodAnnotations",
					testify.Anything,
					testify.Anything,
					testify.Anything,
					testify.Anything,
				).Return(nil)

				ngt, err := New(&config, WithPatcher(mock))
				require.NoError(t, err)

				now := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, now)
				require.NoError(t, err)

				err = ngt.CreateIndex(t.Context(), 10)
				require.NoError(t, err)

				// expected entries
				expected := map[string]string{
					kvald.IndexCountAnnotationsKey:                     "1",
					kvald.UncommittedAnnotationsKey:                    "0",
					kvald.UnsavedCreateIndexExecutionNumAnnotationsKey: "1",
					kvald.UnsavedProcessedVQAnnotationsKey:             "1",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 1)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expected)
			},
		},
		{
			name: "export after create index multiple vectors",
			testfunc: func(t *testing.T) {
				mock := &k8s.PatcherMock{}
				mock.On("ApplyPodAnnotations",
					testify.Anything,
					testify.Anything,
					testify.Anything,
					testify.Anything,
				).Return(nil)

				ngt, err := New(&config, WithPatcher(mock))
				require.NoError(t, err)

				time1 := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, time1)
				require.NoError(t, err)

				time2 := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid2", []float32{1.0, 2.0, 3.0}, time2)
				require.NoError(t, err)

				err = ngt.CreateIndex(t.Context(), 10)
				require.NoError(t, err)

				// expected entries
				expected := map[string]string{
					kvald.IndexCountAnnotationsKey:                     "2",
					kvald.UncommittedAnnotationsKey:                    "0",
					kvald.UnsavedCreateIndexExecutionNumAnnotationsKey: "1",
					kvald.UnsavedProcessedVQAnnotationsKey:             "2",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 1)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expected)
			},
		},
		{
			name: "export after create index multiple times",
			testfunc: func(t *testing.T) {
				mock := &k8s.PatcherMock{}
				mock.On("ApplyPodAnnotations",
					testify.Anything,
					testify.Anything,
					testify.Anything,
					testify.Anything,
				).Return(nil)

				ngt, err := New(&config, WithPatcher(mock))
				require.NoError(t, err)

				time1 := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, time1)
				require.NoError(t, err)

				err = ngt.CreateIndex(t.Context(), 10)
				require.NoError(t, err)

				time2 := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid2", []float32{1.0, 2.0, 3.0}, time2)
				require.NoError(t, err)

				err = ngt.CreateIndex(t.Context(), 10)
				require.NoError(t, err)

				// expected entries
				expected := map[string]string{
					kvald.IndexCountAnnotationsKey:                     "2",
					kvald.UncommittedAnnotationsKey:                    "0",
					kvald.UnsavedCreateIndexExecutionNumAnnotationsKey: "2",
					kvald.UnsavedProcessedVQAnnotationsKey:             "2",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 2)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expected)
			},
		},
		{
			name: "export after create index multiple vectors and save index",
			testfunc: func(t *testing.T) {
				mock := &k8s.PatcherMock{}
				mock.On("ApplyPodAnnotations",
					testify.Anything,
					testify.Anything,
					testify.Anything,
					testify.Anything,
				).Return(nil)

				tmpdir := t.TempDir()

				ngt, err := New(&config,
					WithIndexPath(tmpdir),
					WithPatcher(mock),
				)
				require.NoError(t, err)

				time1 := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, time1)
				require.NoError(t, err)

				time2 := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid2", []float32{1.0, 2.0, 3.0}, time2)
				require.NoError(t, err)

				ctx := t.Context()
				err = ngt.CreateIndex(ctx, 10)
				require.NoError(t, err)

				// set time in context for testing
				saveIndexTime := time.Now()
				ctx = context.WithValue(ctx, saveIndexTimeKey, saveIndexTime)

				err = ngt.SaveIndex(ctx)
				require.NoError(t, err)

				// expected entries
				expectedAfterCreate := map[string]string{
					kvald.IndexCountAnnotationsKey:                     "2",
					kvald.UncommittedAnnotationsKey:                    "0",
					kvald.UnsavedCreateIndexExecutionNumAnnotationsKey: "1",
					kvald.UnsavedProcessedVQAnnotationsKey:             "2",
				}
				expectedAfterSave := map[string]string{
					kvald.LastTimeSaveIndexTimestampAnnotationsKey:     saveIndexTime.UTC().Format(kvald.TimeFormat),
					kvald.UnsavedCreateIndexExecutionNumAnnotationsKey: "0",
					kvald.UnsavedProcessedVQAnnotationsKey:             "0",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 2)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expectedAfterCreate)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expectedAfterSave)
			},
		},
		{
			name: "export after inserting vectors",
			testfunc: func(t *testing.T) {
				mock := &k8s.PatcherMock{}
				mock.On("ApplyPodAnnotations",
					testify.Anything,
					testify.Anything,
					testify.Anything,
					testify.Anything,
				).Return(nil)

				tmpdir := t.TempDir()

				n, err := New(&config,
					WithIndexPath(tmpdir),
					WithPatcher(mock),
				)
				require.NoError(t, err)

				time1 := time.Now().UnixNano()
				err = n.InsertWithTime("test-uuid", []float32{1.0, 2.0, 3.0}, time1)
				require.NoError(t, err)

				time2 := time.Now().UnixNano()
				err = n.InsertWithTime("test-uuid2", []float32{1.0, 2.0, 3.0}, time2)
				require.NoError(t, err)

				ctx := t.Context()
				ExportMetricsOnTick(n.(*ngt), ctx)

				// expected entries
				expectedAfterInsert := map[string]string{
					kvald.IndexCountAnnotationsKey:  "0",
					kvald.UncommittedAnnotationsKey: "2",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 1)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expectedAfterInsert)
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			test.testfunc(tt)
		})
	}
}

func Test_ngt_InsertUpsert(t *testing.T) {
	if testing.Short() {
		t.Skip("The execution of this test takes a lot of time, so it is not performed during the short test\ttest: Test_ngt_InsertUpsert")
		return
	}
	type args struct {
		indices  []index
		poolSize uint32
		bulkSize int
	}
	type fields struct {
		tmpPath           atomic.Value
		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		fmap              map[string]uint32
		svcCfg            *config.NGT
		basePath          string
		oldPath           string
		path              string
		svcOpts           []Option
		sdur              time.Duration
		nocie             uint64
		dim               int
		alen              int
		lim               time.Duration
		nogce             uint64
		dur               time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		kvsdbConcurrency  int
		idelay            time.Duration
		lastNocie         uint64
		backupGen         uint64
		smu               sync.Mutex
		cimu              sync.Mutex
		fmu               sync.Mutex
		cowmu             sync.Mutex
		poolSize          uint32
		radius            float32
		epsilon           float32
		inMem             bool
		enableCopyOnWrite bool
		dcd               bool
		enableProactiveGC bool
	}
	type want struct {
		err error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		fields     fields
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	var (
		// default NGT configuration for test
		kvsdbCfg  = &config.KVSDB{}
		vqueueCfg = &config.VQueue{}
	)
	tests := []test{
		func() test {
			count := 10000000
			return test{
				name: fmt.Sprintf("insert & upsert %d random and 11 digits added to each vector element", count),
				args: args{
					indices: createRandomData(count, &createRandomDataConfig{
						additionaldigits: 11,
					}),
					poolSize: uint32(count / 10),
					bulkSize: count / 10,
				},
				fields: fields{
					svcCfg: &config.NGT{
						Dimension:    128,
						DistanceType: core.Cosine.String(),
						ObjectType:   core.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []Option{
						WithEnableInMemoryMode(true),
					},
				},
			}
		}(),
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx := tt.Context()

			// defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			n, err := New(test.fields.svcCfg, append(test.fields.svcOpts, WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}
			for _, idx := range test.args.indices {
				err = n.Insert(idx.uuid, idx.vec)
				if err := checkFunc(test.want, err); err != nil {
					tt.Errorf("error = %v", err)
				}

			}

			log.Warn("start create index operation")
			err = n.CreateAndSaveIndex(ctx, test.args.poolSize)
			if err != nil {
				tt.Errorf("error creating index: %v", err)
			}
			log.Warn("start update operation")
			for i := range 100 {
				idx := i
				eg.Go(safety.RecoverFunc(func() error {
					log.Warnf("started %d-1", idx)
					for _, idx := range test.args.indices[:len(test.args.indices)/3] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-1", idx)
					return nil
				}))

				eg.Go(safety.RecoverFunc(func() error {
					log.Warnf("started %d-2", idx)
					for _, idx := range test.args.indices[len(test.args.indices)/3 : 2*len(test.args.indices)/3] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-2", idx)
					return nil
				}))

				eg.Go(safety.RecoverFunc(func() error {
					log.Warnf("started %d-3", idx)
					for _, idx := range test.args.indices[2*len(test.args.indices)/3:] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-3", idx)
					return nil
				}))
			}
			eg.Wait()

			log.Warn("start final create index operation")
			err = n.CreateAndSaveIndex(ctx, test.args.poolSize)
			if err != nil {
				tt.Errorf("error creating index: %v", err)
			}
		})
	}
}

// NOTE: After moving this implementation to the e2e package, remove this test function.
func Test_ngt_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("The execution of this test takes a lot of time, so it is not performed during the short test\ttest: Test_ngt_E2E")
		return
	}
	type args struct {
		client   grpc.Client
		addr     string
		requests []*payload.Upsert_MultiRequest
	}
	type want struct {
		err error
	}
	type test struct {
		want       want
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	multiUpsertRequestGenFunc := func(indices []index, chunk int) (res []*payload.Upsert_MultiRequest) {
		reqs := make([]*payload.Upsert_Request, 0, chunk)
		for i := range indices {
			if len(reqs) == chunk-1 {
				res = append(res, &payload.Upsert_MultiRequest{
					Requests: reqs,
				})
				reqs = make([]*payload.Upsert_Request, 0, chunk)
			} else {
				reqs = append(reqs, &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     indices[i].uuid,
						Vector: indices[i].vec,
					},
					Config: &payload.Upsert_Config{
						SkipStrictExistCheck: true,
					},
				})
			}
		}
		if len(reqs) > 0 {
			res = append(res, &payload.Upsert_MultiRequest{
				Requests: reqs,
			})
		}
		return res
	}

	tests := []test{
		{
			name: "insert & upsert 100 random",
			args: args{
				requests: multiUpsertRequestGenFunc(
					createRandomData(500000, new(createRandomDataConfig)),
					50,
				),
				addr:   "127.0.0.1:8080",
				client: grpc.New("NGT Test Client", grpc.WithInsecure(true)),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx := tt.Context()

			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			defer test.args.client.Close(ctx)

			for i := range 2 {
				for _, req := range test.args.requests {
					_, err := test.args.client.Do(ctx, test.args.addr,
						func(ctx context.Context, conn *grpc.ClientConn, opts ...grpc.CallOption) (any, error) {
							return vald.NewFromConn(conn).MultiUpsert(ctx, req)
						})
					if err != nil {
						t.Error(err)
					}
				}
				log.Info("%d step: finished all requests", i+1)
				time.Sleep(3 * time.Second)
			}
		})
	}
}

type createRandomDataConfig struct {
	additionaldigits int
}

func (cfg *createRandomDataConfig) verify() *createRandomDataConfig {
	if cfg == nil {
		cfg = new(createRandomDataConfig)
	}
	if cfg.additionaldigits < 0 {
		cfg.additionaldigits = 0
	}
	return cfg
}

func createRandomData(num int, cfg *createRandomDataConfig) []index {
	cfg = cfg.verify()

	var ad float32 = 1.0
	for i := 0; i < cfg.additionaldigits; i++ {
		ad = ad * 0.1
	}

	result := make([]index, 0)
	f32s, _ := vector.GenF32Vec(vector.NegativeUniform, num, 128)

	for idx, vec := range f32s {
		for i := range vec {
			if f := vec[i] * ad; f == 0.0 {
				if vec[i] > 0.0 {
					vec[i] = math.MaxFloat32
				} else if vec[i] < 0.0 {
					vec[i] = math.SmallestNonzeroFloat32
				}
				continue
			}
			vec[i] = vec[i] * ad
		}
		result = append(result, index{
			uuid: fmt.Sprintf("%s_%s-%s:%d:%d,%d", uuid.New().String(), uuid.New().String(), uuid.New().String(), idx, idx/100, idx%100),
			vec:  vec,
		})
	}

	return result
}

// doneDuringDeletePhaseCtx is done from the start (so vqueue ranges stop after their first
// entry), but its Err reports nil on the first call, i.e. it looks live when createIndexLocked
// starts and is cancelled while the delete phase runs.
//
//nolint:containedctx // test fake overriding Done/Err of a parent ctx
type doneDuringDeletePhaseCtx struct {
	context.Context
	done     chan struct{}
	errCalls atomic.Int32
}

func (c *doneDuringDeletePhaseCtx) Done() <-chan struct{} { return c.done }

func (c *doneDuringDeletePhaseCtx) Err() error {
	if c.errCalls.Add(1) == 1 {
		return nil
	}
	return context.Canceled
}

// Test_ngt_createIndexLocked_doesNotApplyInsertsAfterPartialDeletes pins that a run whose ctx
// is done never applies an insert while the delete of the same Update may still be queued.
func Test_ngt_createIndexLocked_doesNotApplyInsertsAfterPartialDeletes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ctx         func(t *testing.T) context.Context
		name        string
		wantIVQ     int
		wantDVQ     int
		wantCreated bool
	}{
		{
			name: "ctx done before the run starts leaves the vqueue untouched",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			wantIVQ:     2,
			wantDVQ:     2,
			wantCreated: false,
		},
		{
			name: "ctx done during the delete phase keeps every insert queued",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				done := make(chan struct{})
				close(done)
				return &doneDuringDeletePhaseCtx{Context: t.Context(), done: done}
			},
			wantIVQ: 2,
			wantDVQ: 1,
			// the delete already applied is still finished and counted so that it gets saved
			wantCreated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := defaultConfig
			cfg.Dimension = 3

			nn, err := New(&cfg)
			require.NoError(t, err)
			n, ok := nn.(*ngt)
			require.True(t, ok)

			base := time.Now().UnixNano()
			uuids := []string{"uuid-a", "uuid-b"}
			for _, uuid := range uuids {
				require.NoError(t, n.InsertWithTime(uuid, []float32{1, 2, 3}, base))
			}
			require.NoError(t, n.CreateIndex(t.Context(), 10))
			oids := make(map[string]uint32, len(uuids))
			for _, uuid := range uuids {
				oid, _, ok := n.kvs.Get(uuid)
				require.True(t, ok)
				oids[uuid] = oid
			}
			// Update pairs: delete at t, insert at t+1.
			newVec := []float32{4, 5, 6}
			for _, uuid := range uuids {
				require.NoError(t, n.vq.PushDelete(uuid, base+10))
				require.NoError(t, n.vq.PushInsert(uuid, newVec, base+11))
			}
			nocie := n.NumberOfCreateIndexExecution()

			created, err := n.createIndexLocked(tt.ctx(t), 10)
			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, tt.wantCreated, created)
			require.Equal(t, tt.wantIVQ, n.vq.IVQLen(), "no insert may be applied")
			require.Equal(t, tt.wantDVQ, n.vq.DVQLen())
			if tt.wantCreated {
				require.Equal(t, nocie+1, n.NumberOfCreateIndexExecution())
			} else {
				require.Equal(t, nocie, n.NumberOfCreateIndexExecution())
			}
			for _, uuid := range uuids {
				// a uuid whose delete is still queued must keep its old mapping
				if _, _, dts, _ := n.vq.GetVectorWithTimestamp(uuid); dts != 0 {
					oid, _, ok := n.kvs.Get(uuid)
					require.True(t, ok)
					require.Equal(t, oids[uuid], oid)
				}
			}

			// The next run with a live ctx applies the rest.
			require.NoError(t, n.CreateIndex(t.Context(), 10))
			require.Zero(t, n.vq.IVQLen())
			require.Zero(t, n.vq.DVQLen())
			require.Equal(t, uint64(len(uuids)), n.Len())
			for _, uuid := range uuids {
				vec, ts, err := n.GetObject(uuid)
				require.NoError(t, err)
				require.Equal(t, newVec, vec)
				require.Equal(t, base+11, ts)
			}
		})
	}
}

// Test_ngt_copyNGT_keepsCreateIndexWaiters pins that RegenerateIndexes' copyNGT does not reset
// the waiter count of CreateIndex callers blocked on this instance; resetting it made the
// waiter's decrement wrap wfci, disabling CreateIndex dedup and hanging Close's wait.
func Test_ngt_copyNGT_keepsCreateIndexWaiters(t *testing.T) {
	t.Parallel()

	n := new(ngt)
	atomic.StoreUint64(&n.wfci, 1)
	n.copyNGT(new(ngt))
	require.Equal(t, uint64(1), atomic.LoadUint64(&n.wfci))

	// the waiter leaves its wait loop
	atomic.AddUint64(&n.wfci, ^uint64(0))
	require.Zero(t, atomic.LoadUint64(&n.wfci))

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	require.NoError(t, n.waitForInflightIndexing(ctx), "Close must not wait for a waiter that already left")
}

func writeGob(t *testing.T, path string, v any) {
	t.Helper()
	fi, err := file.Open(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	require.NoError(t, err)
	if err := gob.NewEncoder(fi).Encode(v); err != nil {
		t.Fatal(errors.Join(err, fi.Close()))
	}
	require.NoError(t, fi.Close())
}

func closeKVS(t *testing.T, n *ngt) {
	t.Helper()
	if n.kvs != nil {
		if err := n.kvs.Close(); err != nil {
			t.Error(err)
		}
	}
}

// Test_ngt_initNGT_loadTimeout verifies that initNGT returns on a load timeout
// instead of rebuilding the index while the abandoned load is still running.
// The timestamp kvsdb file is a FIFO with no writer, so loadKVS blocks in open
// and load deterministically times out.
func Test_ngt_initNGT_loadTimeout(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	require.NoError(t, metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		NGT: &metadata.NGT{IndexCount: 1},
	}))
	writeGob(t, file.Join(dir, kvsFileName), map[string]uint32{"a": 1})
	fifo := file.Join(dir, kvsTimestampFileName)
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))

	eg, _ := errgroup.New(t.Context())
	n := &ngt{
		path:             dir,
		fmap:             make(map[string]int64),
		kvsdbConcurrency: 1,
		eg:               eg,
		minLit:           100 * time.Millisecond,
		maxLit:           100 * time.Millisecond,
		litFactor:        time.Millisecond,
	}

	err := n.initNGT(core.WithDimension(2))
	if !errors.Is(err, errors.ErrIndexLoadTimeout) {
		t.Errorf("initNGT() err = %v, want ErrIndexLoadTimeout", err)
	}

	// Release the abandoned loadKVS: an O_RDWR open of a FIFO does not block on
	// Linux and serves the reader whether it is already blocked in open or
	// opens later; the fd stays open until the abandoned load has finished.
	w, err := file.Open(fifo, os.O_RDWR, 0)
	if err != nil {
		// a fallback removed the FIFO, so the abandoned load cannot be released.
		t.Fatal(err)
	}
	if err := gob.NewEncoder(w).Encode(map[string]int64{}); err != nil {
		t.Error(err)
	}
	done := make(chan error, 1)
	go func() { done <- eg.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Error(err)
		}
		closeKVS(t, n)
		if n.core != nil {
			n.core.Close()
		}
	case <-time.After(30 * time.Second):
		t.Error("abandoned load did not finish")
	}
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

// saveTestIndex stores an NGT index of one vector at dir together with its
// metadata (claiming indexCount entries) and a kvsdb with the given uuids.
func saveTestIndex(t *testing.T, dir string, indexCount uint64, uuids map[string]uint32) {
	t.Helper()
	require.NoError(t, file.MkdirAll(dir, os.ModePerm))
	c, err := core.New(core.WithIndexPath(dir), core.WithDimension(2))
	require.NoError(t, err)
	_, err = c.Insert([]float32{1, 2})
	require.NoError(t, errors.Join(err, c.CreateIndex(1), c.SaveIndex()))
	c.Close()
	require.NoError(t, metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		NGT: &metadata.NGT{IndexCount: indexCount},
	}))
	ts := make(map[string]int64, len(uuids))
	for k := range uuids {
		ts[k] = 1
	}
	writeGob(t, file.Join(dir, kvsFileName), uuids)
	writeGob(t, file.Join(dir, kvsTimestampFileName), ts)
}

// Test_ngt_initNGT_reloadsTooFewPrimary verifies that, in CoW mode, when the
// primary kvsdb is too small for its metadata and the secondary has even fewer
// entries, initNGT starts with the reloaded primary. The reload reports the
// same ErrIndicesAreTooFewComparedToMetadata as the first primary load; treating
// it as a failure fell through to a fresh empty index and kvsdb.
func Test_ngt_initNGT_reloadsTooFewPrimary(t *testing.T) {
	t.Parallel()
	base := t.TempDir()
	eg, _ := errgroup.New(t.Context())
	n := &ngt{
		basePath:          base,
		path:              file.Join(base, originIndexDirName),
		oldPath:           file.Join(base, oldIndexDirName),
		brokenPath:        file.Join(base, brokenIndexDirName),
		fmap:              make(map[string]int64),
		kvsdbConcurrency:  1,
		eg:                eg,
		minLit:            time.Minute,
		maxLit:            time.Minute,
		litFactor:         time.Millisecond,
		enableCopyOnWrite: true,
	}
	n.tmpPath.Store(t.TempDir())
	// primary: 2 entries for 10 indices (too few); secondary: 1 entry, consistent.
	saveTestIndex(t, n.path, 10, map[string]uint32{"a": 1, "b": 2})
	saveTestIndex(t, n.oldPath, 1, map[string]uint32{"a": 1})

	require.NoError(t, n.initNGT(core.WithDimension(2)))
	defer closeKVS(t, n)
	require.NotNil(t, n.core)
	defer n.core.Close()
	require.Equal(t, uint64(2), n.kvs.Len(), "the reloaded primary kvsdb must be kept")
	for _, uuid := range []string{"a", "b"} {
		_, _, ok := n.kvs.Get(uuid)
		require.True(t, ok, "kvs entry of %s is missing", uuid)
	}
}

// Test_ngt_initNGT_keepsTooFewPrimaryWithoutSecondary verifies that, in CoW
// mode, when the primary kvsdb is too small for its metadata and neither the
// secondary nor the non-CoW index can be loaded, initNGT starts with the
// reloaded primary instead of discarding it for a fresh empty index.
func Test_ngt_initNGT_keepsTooFewPrimaryWithoutSecondary(t *testing.T) {
	t.Parallel()
	base := t.TempDir()
	eg, _ := errgroup.New(t.Context())
	n := &ngt{
		basePath:          base,
		path:              file.Join(base, originIndexDirName),
		oldPath:           file.Join(base, oldIndexDirName),
		brokenPath:        file.Join(base, brokenIndexDirName),
		fmap:              make(map[string]int64),
		kvsdbConcurrency:  1,
		eg:                eg,
		minLit:            time.Minute,
		maxLit:            time.Minute,
		litFactor:         time.Millisecond,
		enableCopyOnWrite: true,
	}
	n.tmpPath.Store(t.TempDir())
	// primary: 2 entries for 10 indices (too few); no secondary, and base holds
	// no metadata, so both fallback loads fail.
	saveTestIndex(t, n.path, 10, map[string]uint32{"a": 1, "b": 2})

	require.NoError(t, n.initNGT(core.WithDimension(2)))
	defer closeKVS(t, n)
	require.NotNil(t, n.core)
	defer n.core.Close()
	require.Equal(t, uint64(2), n.kvs.Len(), "the reloaded primary kvsdb must be kept")
	for _, uuid := range []string{"a", "b"} {
		_, _, ok := n.kvs.Get(uuid)
		require.True(t, ok, "kvs entry of %s is missing", uuid)
	}
}

// closeCountCore counts Close calls on the wrapped core.
type closeCountCore struct {
	core.NGT
	closes atomic.Int32
}

func (c *closeCountCore) Close() {
	c.closes.Add(1)
	c.NGT.Close()
}

// Test_ngt_RegenerateIndexes_closesOldCore verifies that Flush frees the
// replaced core instead of leaking the whole C-side index on every call.
func Test_ngt_RegenerateIndexes_closesOldCore(t *testing.T) {
	t.Parallel()
	cfg := defaultConfig
	cfg.Dimension = 3
	nn, err := New(&cfg)
	require.NoError(t, err)
	n, ok := nn.(*ngt)
	require.True(t, ok)
	old := &closeCountCore{NGT: n.core}
	n.core = old

	require.NoError(t, n.RegenerateIndexes(t.Context()))
	require.Equal(t, int32(1), old.closes.Load(), "the replaced core must be closed once")
	require.NotSame(t, old, n.core)
	require.NotNil(t, n.core)
	require.NoError(t, n.Close(t.Context()))
	require.Equal(t, int32(1), old.closes.Load(), "Close must not close the replaced core again")
}

// Test_ngt_multiple_flatErrors verifies that insertMultiple and deleteMultiple
// join the per-ID errors once into a flat error (one entry per failed ID)
// instead of nesting one level per failure, and that every per-ID error the
// gRPC handlers look up with errors.Is is still reported.
func Test_ngt_multiple_flatErrors(t *testing.T) {
	t.Parallel()
	const n = 4
	uuids := make([]string, 0, n)
	for i := range n {
		uuids = append(uuids, fmt.Sprintf("flat-err-uuid-%d", i))
	}

	newNGT := func(t *testing.T, stored []string) *ngt {
		t.Helper()
		vq, err := vqueue.New()
		if err != nil {
			t.Fatal(err)
		}
		kv := kvs.New()
		now := time.Now().UnixNano()
		for i, uuid := range stored {
			kv.Set(uuid, uint32(i+1), now)
		}
		return &ngt{kvs: kv, vq: vq}
	}

	tests := []struct {
		run     func(t *testing.T) error
		want    func(uuid string) error
		name    string
		ordered bool
	}{
		{
			name: "insertMultiple reports ErrUUIDAlreadyExists for every stored id",
			run: func(t *testing.T) error {
				t.Helper()
				vecs := make(map[string][]float32, n)
				for _, uuid := range uuids {
					vecs[uuid] = []float32{1}
				}
				return newNGT(t, uuids).insertMultiple(vecs, time.Now().UnixNano(), true)
			},
			want: errors.ErrUUIDAlreadyExists,
		},
		{
			name: "deleteMultiple reports ErrObjectIDNotFound for every missing id newest first",
			run: func(t *testing.T) error {
				t.Helper()
				return newNGT(t, nil).deleteMultiple(uuids, time.Now().UnixNano(), true)
			},
			want:    errors.ErrObjectIDNotFound,
			ordered: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.run(t)
			if err == nil {
				t.Fatal("got nil error, want joined per-ID errors")
			}
			u, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf("got %T, want a joined error", err)
			}
			got := u.Unwrap()
			if len(got) != n {
				t.Fatalf("got %d flat entries, want %d: %v", len(got), n, err)
			}
			for i, uuid := range uuids {
				want := tc.want(uuid)
				if !errors.Is(err, want) {
					t.Errorf("errors.Is(err, %v) = false, err = %v", want, err)
				}
				if tc.ordered && got[n-1-i].Error() != want.Error() {
					t.Errorf("entry %d = %v, want %v (newest first)", n-1-i, got[n-1-i], want)
				}
			}
		})
	}
}

// Test_ngt_toSearchResponse_dropsNaNDistance verifies that results whose
// distance is NaN are dropped while the order of the finite results is kept,
// and that an all-NaN result set on a non-empty index is reported as
// ErrEmptySearchResult.
func Test_ngt_toSearchResponse_dropsNaNDistance(t *testing.T) {
	t.Parallel()
	nan := float32(math.NaN())
	// kvs oid -> uuid: 1 -> "a", 2 -> "b", 3 -> "c"
	stored := []string{"a", "b", "c"}
	tests := []struct {
		wantErr   error
		name      string
		stored    []string
		sr        []algorithm.SearchResult
		wantIDs   []string
		wantDists []float32
		wantNil   bool
	}{
		{
			name:   "mixed NaN and finite keeps only finite in the original order",
			stored: stored,
			sr: []algorithm.SearchResult{
				{ID: 3, Distance: 0.3},
				{ID: 1, Distance: nan},
				{ID: 2, Distance: 0.5},
				{ID: 4, Distance: nan},
			},
			wantIDs:   []string{"c", "b"},
			wantDists: []float32{0.3, 0.5},
		},
		{
			name:   "all NaN on a non-empty index returns ErrEmptySearchResult",
			stored: stored,
			sr: []algorithm.SearchResult{
				{ID: 1, Distance: nan},
				{ID: 2, Distance: nan},
			},
			wantErr: errors.ErrEmptySearchResult,
			wantNil: true,
		},
		{
			name:    "all NaN on an empty index returns nil, nil",
			sr:      []algorithm.SearchResult{{ID: 1, Distance: nan}},
			wantNil: true,
		},
		{
			name:    "empty input on an empty index returns nil, nil",
			wantNil: true,
		},
		{
			name:    "empty input on a non-empty index returns ErrEmptySearchResult",
			stored:  stored,
			wantErr: errors.ErrEmptySearchResult,
			wantNil: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			kv := kvs.New()
			now := time.Now().UnixNano()
			for i, uuid := range tc.stored {
				kv.Set(uuid, uint32(i+1), now)
			}
			n := &ngt{kvs: kv}
			res, err := n.toSearchResponse(tc.sr)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("toSearchResponse() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantNil {
				if res != nil {
					t.Fatalf("toSearchResponse() = %v, want nil", res)
				}
				return
			}
			var gotIDs []string
			var gotDists []float32
			for _, r := range res.GetResults() {
				gotIDs = append(gotIDs, r.GetId())
				gotDists = append(gotDists, r.GetDistance())
			}
			if !slices.Equal(gotIDs, tc.wantIDs) || !slices.Equal(gotDists, tc.wantDists) {
				t.Errorf("toSearchResponse() = ids %v dists %v, want ids %v dists %v", gotIDs, gotDists, tc.wantIDs, tc.wantDists)
			}
		})
	}
}

// Test_ngt_createIndexLocked_newestWins pins that CreateIndex applies a queued delete or
// insert only when it is newer than the indexed object. A delete older than the indexed
// object used to remove it, and a stale Update (delete at t, insert at t+1 <= indexed ts)
// used to repoint the uuid to the stale vector. An insert applied over an indexed object
// must remove the replaced object instead of orphaning it in the core.
func Test_ngt_createIndexLocked_newestWins(t *testing.T) {
	t.Parallel()

	const uuid = "uuid"
	oldVec := []float32{1, 2, 3}
	newVec := []float32{4, 5, 6}
	tests := []struct {
		op         func(t *testing.T, n *ngt, base int64)
		wantTS     func(base int64) int64
		name       string
		wantVec    []float32
		wantExists bool
		// wantOldRemoved reports whether the initially indexed object is gone from the core.
		wantOldRemoved bool
	}{
		{
			name: "delete older than the indexed object is skipped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.DeleteWithTime(uuid, base-10))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "stale update is skipped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base-5))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "update whose insert ties the indexed timestamp is skipped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base-1))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "delete at the indexed timestamp is applied",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.DeleteWithTime(uuid, base))
			},
			wantOldRemoved: true,
		},
		{
			name: "newer delete is applied",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.DeleteWithTime(uuid, base+10))
			},
			wantOldRemoved: true,
		},
		{
			name: "newer update is applied",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base+10))
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 11 },
			wantOldRemoved: true,
		},
		{
			name: "UpdateTimestamp of a queued update keeps the new vector",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base+10))
				require.NoError(t, n.UpdateTimestamp(uuid, base+20, false))
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 20 },
			wantOldRemoved: true,
		},
		{
			// UpdateTimestamp used to raise the kvs timestamp past the queued delete, so
			// CreateIndex skipped it and served the superseded indexed vector.
			name: "UpdateTimestamp after a queued update and a newer delete does not resurrect",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base+10))
				require.NoError(t, n.DeleteWithTime(uuid, base+20))
				require.Error(t, n.UpdateTimestamp(uuid, base+30, false))
			},
			wantOldRemoved: true,
		},
		{
			name: "newer insert over the indexed object replaces it",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.vq.PushInsert(uuid, newVec, base+5))
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 5 },
			wantOldRemoved: true,
		},
		{
			name: "older insert over the indexed object is dropped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.vq.PushInsert(uuid, newVec, base-5))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := defaultConfig
			cfg.Dimension = len(oldVec)

			nn, err := New(&cfg)
			require.NoError(t, err)
			n, ok := nn.(*ngt)
			require.True(t, ok)

			base := time.Now().UnixNano()
			require.NoError(t, n.InsertWithTime(uuid, oldVec, base))
			require.NoError(t, n.CreateIndex(t.Context(), 10))
			oldOid, _, ok := n.kvs.Get(uuid)
			require.True(t, ok)

			tt.op(t, n, base)
			require.NoError(t, n.CreateIndex(t.Context(), 10))

			require.Zero(t, n.vq.IVQLen())
			require.Zero(t, n.vq.DVQLen())
			coreVec, err := n.core.GetVector(uint(oldOid))
			if tt.wantOldRemoved {
				// NGT reuses removed object IDs, so an insert applied after the old object was
				// removed may get oldOid back. Only require that oldOid no longer holds oldVec.
				if err == nil {
					require.NotEqual(t, oldVec, coreVec, "the replaced or deleted object must be removed from the core")
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, oldVec, coreVec)
			}

			oid, kts, ok := n.kvs.Get(uuid)
			require.Equal(t, tt.wantExists, ok)
			if !tt.wantExists {
				require.Zero(t, n.Len())
				_, _, err = n.GetObject(uuid)
				require.Error(t, err)
				return
			}
			require.Equal(t, uint64(1), n.Len())
			if !tt.wantOldRemoved {
				require.Equal(t, oldOid, oid, "a skipped operation must not repoint the uuid")
			}
			require.Equal(t, tt.wantTS(base), kts)
			vec, ts, err := n.GetObject(uuid)
			require.NoError(t, err)
			require.Equal(t, tt.wantVec, vec)
			require.Equal(t, tt.wantTS(base), ts)
		})
	}
}

// Test_ngt_saveIndex_cancelledCtxFullSnapshot verifies that a save started
// with a cancelled context still persists the whole kvs. kvs.Range stops
// silently on cancellation, so the snapshot used to be partial (here: empty)
// while kvsdb was truncated and the metadata IndexCount matched the truncated
// count, hiding the loss from the load-time IndexCount/2 > kvs.Len() check.
// This is the path of the Start loop's SaveIndex on SIGTERM and of a SaveIndex
// RPC whose client cancels.
func Test_ngt_saveIndex_cancelledCtxFullSnapshot(t *testing.T) {
	t.Parallel()
	const num = 100

	svc, err := New(&defaultConfig, WithIndexPath(t.TempDir()))
	if err != nil {
		t.Fatalf("failed to create ngt service: %v", err)
	}
	n, ok := svc.(*ngt)
	if !ok {
		t.Fatalf("New returned %T, want *ngt", svc)
	}
	defer func() {
		if cerr := n.Close(context.WithoutCancel(t.Context())); cerr != nil {
			t.Errorf("Close() err = %v", cerr)
		}
	}()

	for i := range num {
		vec := make([]float32, defaultConfig.Dimension)
		for j := range vec {
			vec[j] = float32(i*defaultConfig.Dimension + j + 1)
		}
		if err = n.Insert("uuid-"+strconv.Itoa(i), vec); err != nil {
			t.Fatalf("Insert(%d) err = %v", i, err)
		}
	}
	if err = n.CreateIndex(t.Context(), 4); err != nil {
		t.Fatalf("CreateIndex() err = %v", err)
	}
	if got := n.kvs.Len(); got != num {
		t.Fatalf("setup: kvs.Len() = %d, want %d", got, num)
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err = n.saveIndex(ctx); err != nil {
		t.Fatalf("saveIndex() err = %v, want nil", err)
	}

	readGob := func(name string, v any) {
		t.Helper()
		fi, oerr := file.Open(file.Join(n.path, name), os.O_RDONLY, os.ModePerm)
		if oerr != nil {
			t.Fatal(oerr)
		}
		if derr := gob.NewDecoder(fi).Decode(v); derr != nil {
			t.Fatal(errors.Join(derr, fi.Close()))
		}
		if cerr := fi.Close(); cerr != nil {
			t.Fatal(cerr)
		}
	}
	m := make(map[string]uint32)
	readGob(kvsFileName, &m)
	if len(m) != num {
		t.Errorf("kvsdb entries = %d, want %d", len(m), num)
	}
	mt := make(map[string]int64)
	readGob(kvsTimestampFileName, &mt)
	if len(mt) != num {
		t.Errorf("kvs timestamp entries = %d, want %d", len(mt), num)
	}
	meta, err := metadata.Load(file.Join(n.path, metadata.AgentMetadataFileName))
	if err != nil {
		t.Fatal(err)
	}
	if meta.NGT == nil || meta.NGT.IndexCount != num {
		t.Errorf("metadata = %+v, want NGT.IndexCount %d", meta.NGT, num)
	}
}

// update deletes at t and inserts at t+1. t == MaxInt64 used to wrap the
// insert to MinInt64, which the queue drops behind the MaxInt64 delete: the
// update reported success and the object disappeared. It must be rejected
// before the delete; ordinary timestamps keep the delete-at-t / insert-at-t+1.
func Test_ngt_update_maxTimestamp(t *testing.T) {
	t.Parallel()
	const (
		uuid = "update-ts-uuid"
		oldT = int64(100)
	)
	oldVec, newVec := []float32{9}, []float32{1}
	tests := []struct {
		run     func(n *ngt, ts int64) error
		name    string
		ts      int64
		wantErr bool
	}{
		{name: "update rejects MaxInt64", run: func(n *ngt, ts int64) error { return n.update(uuid, newVec, ts) }, ts: math.MaxInt64, wantErr: true},
		{name: "update keeps a normal timestamp", run: func(n *ngt, ts int64) error { return n.update(uuid, newVec, ts) }, ts: 200},
		{name: "update keeps MaxInt64-1", run: func(n *ngt, ts int64) error { return n.update(uuid, newVec, ts) }, ts: math.MaxInt64 - 1},
		{name: "updateMultiple rejects MaxInt64", run: func(n *ngt, ts int64) error {
			return n.updateMultiple(map[string][]float32{uuid: newVec}, ts)
		}, ts: math.MaxInt64, wantErr: true},
		{name: "updateMultiple keeps a normal timestamp", run: func(n *ngt, ts int64) error {
			return n.updateMultiple(map[string][]float32{uuid: newVec}, ts)
		}, ts: 200},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			kv := kvs.New()
			t.Cleanup(func() {
				if cerr := kv.Close(); cerr != nil {
					t.Error(cerr)
				}
			})
			if err = vq.PushInsert(uuid, oldVec, oldT); err != nil {
				t.Fatal(err)
			}
			n := &ngt{kvs: kv, vq: vq, dim: 1}

			err = test.run(n, test.ts)
			vec, its, dts, exists := vq.GetVectorWithTimestamp(uuid)
			if test.wantErr {
				if !errors.Is(err, errors.ErrInvalidTimestamp(test.ts)) {
					t.Fatalf("got error %v, want %v", err, errors.ErrInvalidTimestamp(test.ts))
				}
				if !exists || !slices.Equal(vec, oldVec) || its != oldT || dts != 0 {
					t.Errorf("rejected update changed the queue: vec=%v its=%d dts=%d exists=%t", vec, its, dts, exists)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !exists || !slices.Equal(vec, newVec) || dts != test.ts || its != test.ts+1 {
				t.Errorf("got vec=%v its=%d dts=%d exists=%t, want vec=%v its=%d dts=%d", vec, its, dts, exists, newVec, test.ts+1, test.ts)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_newNGT(t *testing.T) {
// 	type args struct {
// 		cfg  *config.NGT
// 		opts []Option
// 	}
// 	type want struct {
// 		wantN *ngt
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *ngt, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN *ngt, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			gotN, err := newNGT(test.args.cfg, test.args.opts...)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_copyNGT(t *testing.T) {
// 	type args struct {
// 		src *ngt
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           src:ngt{},
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           src:ngt{},
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			n.copyNGT(test.args.src)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_migrate(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		path string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
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
// 		           path:"",
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
// 		           path:"",
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
// 			err := migrate(test.args.ctx, test.args.path)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_prepareFolders(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.prepareFolders(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_load(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		path string
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           path:"",
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           path:"",
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.load(test.args.ctx, test.args.path, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_backupBroken(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.backupBroken(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_rebuild(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		path string
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           path:"",
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           path:"",
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.rebuild(test.args.ctx, test.args.path, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_initNGT(t *testing.T) {
// 	type args struct {
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.initNGT(test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_loadKVS(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		path    string
// 		timeout time.Duration
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           path:"",
// 		           timeout:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           path:"",
// 		           timeout:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.loadKVS(test.args.ctx, test.args.path, test.args.timeout)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want <-chan error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error) error {
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_Search(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		vec      []float32
// 		size     uint32
// 		epsilon  float32
// 		radius   float32
// 		edgeSize int32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           vec:nil,
// 		           size:0,
// 		           epsilon:0,
// 		           radius:0,
// 		           edgeSize:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vec:nil,
// 		           size:0,
// 		           epsilon:0,
// 		           radius:0,
// 		           edgeSize:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotRes, err := n.Search(test.args.ctx, test.args.vec, test.args.size, test.args.epsilon, test.args.radius, test.args.edgeSize)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_SearchByID(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		uuid     string
// 		size     uint32
// 		epsilon  float32
// 		radius   float32
// 		edgeSize int32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantVec []float32
// 		wantDst *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []float32, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotDst *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotDst, w.wantDst) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDst, w.wantDst)
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
// 		           uuid:"",
// 		           size:0,
// 		           epsilon:0,
// 		           radius:0,
// 		           edgeSize:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           size:0,
// 		           epsilon:0,
// 		           radius:0,
// 		           edgeSize:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotVec, gotDst, err := n.SearchByID(test.args.ctx, test.args.uuid, test.args.size, test.args.epsilon, test.args.radius, test.args.edgeSize)
// 			if err := checkFunc(test.want, gotVec, gotDst, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_LinearSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		vec  []float32
// 		size uint32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           vec:nil,
// 		           size:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vec:nil,
// 		           size:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotRes, err := n.LinearSearch(test.args.ctx, test.args.vec, test.args.size)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_LinearSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		uuid string
// 		size uint32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantVec []float32
// 		wantDst *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []float32, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotDst *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotDst, w.wantDst) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDst, w.wantDst)
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
// 		           uuid:"",
// 		           size:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           size:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotVec, gotDst, err := n.LinearSearchByID(test.args.ctx, test.args.uuid, test.args.size)
// 			if err := checkFunc(test.want, gotVec, gotDst, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_Insert(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.Insert(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_InsertWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.InsertWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_insert(t *testing.T) {
// 	type args struct {
// 		uuid       string
// 		vec        []float32
// 		t          int64
// 		validation bool
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.insert(test.args.uuid, test.args.vec, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_InsertMultiple(t *testing.T) {
// 	type args struct {
// 		vecs map[string][]float32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           vecs:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vecs:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.InsertMultiple(test.args.vecs)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_InsertMultipleWithTime(t *testing.T) {
// 	type args struct {
// 		vecs map[string][]float32
// 		t    int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           vecs:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vecs:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.InsertMultipleWithTime(test.args.vecs, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_insertMultiple(t *testing.T) {
// 	type args struct {
// 		vecs       map[string][]float32
// 		now        int64
// 		validation bool
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           vecs:nil,
// 		           now:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vecs:nil,
// 		           now:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.insertMultiple(test.args.vecs, test.args.now, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_Update(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.Update(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_UpdateWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.UpdateWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_update(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.update(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_UpdateMultiple(t *testing.T) {
// 	type args struct {
// 		vecs map[string][]float32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           vecs:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vecs:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.UpdateMultiple(test.args.vecs)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_UpdateMultipleWithTime(t *testing.T) {
// 	type args struct {
// 		vecs map[string][]float32
// 		t    int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           vecs:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vecs:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.UpdateMultipleWithTime(test.args.vecs, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_updateMultiple(t *testing.T) {
// 	type args struct {
// 		vecs map[string][]float32
// 		t    int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           vecs:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           vecs:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.updateMultiple(test.args.vecs, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_UpdateTimestamp(t *testing.T) {
// 	type args struct {
// 		uuid  string
// 		ts    int64
// 		force bool
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.UpdateTimestamp(test.args.uuid, test.args.ts, test.args.force)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_Delete(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.Delete(test.args.uuid)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_DeleteWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		t    int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.DeleteWithTime(test.args.uuid, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_delete(t *testing.T) {
// 	type args struct {
// 		uuid       string
// 		t          int64
// 		validation bool
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           t:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           t:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.delete(test.args.uuid, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_DeleteMultiple(t *testing.T) {
// 	type args struct {
// 		uuids []string
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuids:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuids:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.DeleteMultiple(test.args.uuids...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_DeleteMultipleWithTime(t *testing.T) {
// 	type args struct {
// 		uuids []string
// 		t     int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuids:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuids:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.DeleteMultipleWithTime(test.args.uuids, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_deleteMultiple(t *testing.T) {
// 	type args struct {
// 		uuids      []string
// 		now        int64
// 		validation bool
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuids:nil,
// 		           now:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuids:nil,
// 		           now:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.deleteMultiple(test.args.uuids, test.args.now, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_RegenerateIndexes(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.RegenerateIndexes(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_loadStatistics(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.loadStatistics(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_removeInvalidIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			n.removeInvalidIndex(test.args.ctx)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_saveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.saveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_CreateAndSaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		poolSize uint32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           poolSize:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           poolSize:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.CreateAndSaveIndex(test.args.ctx, test.args.poolSize)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_moveAndSwitchSavedData(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.moveAndSwitchSavedData(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_mktmp(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.mktmp()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_Exists(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantOid uint32
// 		wantOk  bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotOid uint32, gotOk bool) error {
// 		if !reflect.DeepEqual(gotOid, w.wantOid) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOid, w.wantOid)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotOid, gotOk := n.Exists(test.args.uuid)
// 			if err := checkFunc(test.want, gotOid, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_GetObject(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantVec       []float32
// 		wantTimestamp int64
// 		err           error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []float32, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotTimestamp int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotVec, gotTimestamp, err := n.GetObject(test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, gotTimestamp, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_readyForUpdate(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		ts   int64
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           uuid:"",
// 		           vec:nil,
// 		           ts:0,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           uuid:"",
// 		           vec:nil,
// 		           ts:0,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.readyForUpdate(test.args.uuid, test.args.vec, test.args.ts)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_IsSaving(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.IsSaving()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_IsIndexing(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.IsIndexing()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_IsFlushing(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.IsFlushing()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_UUIDs(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantUuids []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotUuids []string) error {
// 		if !reflect.DeepEqual(gotUuids, w.wantUuids) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuids, w.wantUuids)
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotUuids := n.UUIDs(test.args.ctx)
// 			if err := checkFunc(test.want, gotUuids); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_NumberOfCreateIndexExecution(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.NumberOfCreateIndexExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_NumberOfProactiveGCExecution(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.NumberOfProactiveGCExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_lastNumberOfCreateIndexExecution(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.lastNumberOfCreateIndexExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_gc(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			n.gc()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_Len(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.Len()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_InsertVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.InsertVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_DeleteVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.DeleteVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_GetDimensionSize(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.GetDimensionSize()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_BrokenIndexCount(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.BrokenIndexCount()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_ListObjectFunc(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(uuid string, oid uint32, ts int64) bool
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			n.ListObjectFunc(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_IndexStatistics(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantStats *payload.Info_Index_Statistics
// 		err       error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Statistics, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotStats *payload.Info_Index_Statistics, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotStats, w.wantStats) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStats, w.wantStats)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotStats, err := n.IndexStatistics()
// 			if err := checkFunc(test.want, gotStats, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_IsStatisticsEnabled(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got := n.IsStatisticsEnabled()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_IndexProperty(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		want *payload.Info_Index_Property
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Property, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Info_Index_Property, err error) error {
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
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			got, err := n.IndexProperty()
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_toSearchResponse(t *testing.T) {
// 	type args struct {
// 		sr []algorithm.SearchResult
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           sr:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           sr:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotRes, err := n.toSearchResponse(test.args.sr)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_uncommittedEntry(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantK string
// 		wantV string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotK string, gotV string) error {
// 		if !reflect.DeepEqual(gotK, w.wantK) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotK, w.wantK)
// 		}
// 		if !reflect.DeepEqual(gotV, w.wantV) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotK, gotV := n.uncommittedEntry()
// 			if err := checkFunc(test.want, gotK, gotV); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_processedVqEntries(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantK string
// 		wantV string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotK string, gotV string) error {
// 		if !reflect.DeepEqual(gotK, w.wantK) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotK, w.wantK)
// 		}
// 		if !reflect.DeepEqual(gotV, w.wantV) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotK, gotV := n.processedVqEntries()
// 			if err := checkFunc(test.want, gotK, gotV); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_unsavedNumberOfCreateIndexExecutionEntry(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantK string
// 		wantV string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotK string, gotV string) error {
// 		if !reflect.DeepEqual(gotK, w.wantK) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotK, w.wantK)
// 		}
// 		if !reflect.DeepEqual(gotV, w.wantV) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotK, gotV := n.unsavedNumberOfCreateIndexExecutionEntry()
// 			if err := checkFunc(test.want, gotK, gotV); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_lastTimeSaveIndexTimestampEntry(t *testing.T) {
// 	type args struct {
// 		timestamp time.Time
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantK string
// 		wantV string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotK string, gotV string) error {
// 		if !reflect.DeepEqual(gotK, w.wantK) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotK, w.wantK)
// 		}
// 		if !reflect.DeepEqual(gotV, w.wantV) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           timestamp:time.Time{},
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           timestamp:time.Time{},
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotK, gotV := n.lastTimeSaveIndexTimestampEntry(test.args.timestamp)
// 			if err := checkFunc(test.want, gotK, gotV); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_indexCountEntry(t *testing.T) {
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
// 	}
// 	type want struct {
// 		wantK string
// 		wantV string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotK string, gotV string) error {
// 		if !reflect.DeepEqual(gotK, w.wantK) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotK, w.wantK)
// 		}
// 		if !reflect.DeepEqual(gotV, w.wantV) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			gotK, gotV := n.indexCountEntry()
// 			if err := checkFunc(test.want, gotK, gotV); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_exportMetrics(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		entryFuncs []func() (k, v string)
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           entryFuncs:nil,
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           entryFuncs:nil,
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.exportMetrics(test.args.ctx, test.args.entryFuncs...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_exportMetricsOnTick(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.exportMetricsOnTick(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_exportMetricsOnCreateIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.exportMetricsOnCreateIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_exportMetricsOnSaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		tmpPath                 atomic.Value
// 		kvs                     kvs.BidiMap
// 		saving                  atomic.Value
// 		eg                      errgroup.Group
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		patcher                 client.Patcher
// 		fmap                    map[string]int64
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
// 		cfg                     *config.NGT
// 		basePath                string
// 		brokenPath              string
// 		oldPath                 string
// 		podName                 string
// 		path                    string
// 		podNamespace            string
// 		opts                    []Option
// 		nogce                   uint64
// 		kvsdbConcurrency        int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		exportIndexInfoDuration time.Duration
// 		historyLimit            int
// 		dim                     int
// 		nopvq                   atomic.Uint64
// 		nobic                   uint64
// 		idelay                  time.Duration
// 		wfci                    uint64
// 		nocie                   uint64
// 		lastNocie               uint64
// 		flushing                atomic.Bool
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		dcd                     bool
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		enableStatistics        bool
// 		inMem                   bool
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 		           core:nil,
// 		           tmpPath:nil,
// 		           kvs:nil,
// 		           saving:nil,
// 		           eg:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           patcher:nil,
// 		           fmap:nil,
// 		           statisticsCache:nil,
// 		           cfg:nil,
// 		           basePath:"",
// 		           brokenPath:"",
// 		           oldPath:"",
// 		           podName:"",
// 		           path:"",
// 		           podNamespace:"",
// 		           opts:nil,
// 		           nogce:0,
// 		           kvsdbConcurrency:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           exportIndexInfoDuration:nil,
// 		           historyLimit:0,
// 		           dim:0,
// 		           nopvq:nil,
// 		           nobic:0,
// 		           idelay:nil,
// 		           wfci:0,
// 		           nocie:0,
// 		           lastNocie:0,
// 		           flushing:nil,
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           dcd:false,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           enableStatistics:false,
// 		           inMem:false,
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
// 			n := &ngt{
// 				core:                    test.fields.core,
// 				tmpPath:                 test.fields.tmpPath,
// 				kvs:                     test.fields.kvs,
// 				saving:                  test.fields.saving,
// 				eg:                      test.fields.eg,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				patcher:                 test.fields.patcher,
// 				fmap:                    test.fields.fmap,
// 				statisticsCache:         test.fields.statisticsCache,
// 				cfg:                     test.fields.cfg,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				oldPath:                 test.fields.oldPath,
// 				podName:                 test.fields.podName,
// 				path:                    test.fields.path,
// 				podNamespace:            test.fields.podNamespace,
// 				opts:                    test.fields.opts,
// 				nogce:                   test.fields.nogce,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				historyLimit:            test.fields.historyLimit,
// 				dim:                     test.fields.dim,
// 				nopvq:                   test.fields.nopvq,
// 				nobic:                   test.fields.nobic,
// 				idelay:                  test.fields.idelay,
// 				wfci:                    test.fields.wfci,
// 				nocie:                   test.fields.nocie,
// 				lastNocie:               test.fields.lastNocie,
// 				flushing:                test.fields.flushing,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				dcd:                     test.fields.dcd,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				enableStatistics:        test.fields.enableStatistics,
// 				inMem:                   test.fields.inMem,
// 			}
//
// 			err := n.exportMetricsOnSaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
