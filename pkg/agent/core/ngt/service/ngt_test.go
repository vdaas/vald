//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package service manages the main logic of server.
package service

import (
	"context"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	kvald "github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
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
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
					file.CopyDir(context.Background(), testIndexDir, tmpDir)
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
					err := file.CopyDir(context.Background(), testIndexDir, tmpDir)
					if err != nil {
						t.Errorf("failed to copy testdata index files: %v", err)
					}

					// copy testdata index files into tmpDir which is a old index directory
					err = file.MkdirAll(originDir, fs.ModePerm)
					if err != nil {
						t.Errorf("failed to create origin directory: %v", err)
					}
					err = file.CopyDir(context.Background(), testIndexDir, originDir)
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
					file.CopyDir(context.Background(), testIndexDir, originDir)
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
					file.CopyDir(context.Background(), testIndexDir, originDir)
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
					file.CopyDir(context.Background(), testIndexDir, originDir)
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
					if err := file.CopyDir(context.Background(), testIndexDir, originDir); err != nil {
						t.Errorf("failed to copy test index: %v", err)
					}
					// remove metadata.json to make it broken
					if err := os.Remove(filepath.Join(originDir, "metadata.json")); err != nil {
						t.Errorf("failed to remove index file: %v", err)
					}

					if err := file.MkdirAll(backupDir, fs.ModePerm); err != nil {
						t.Errorf("failed to create backup dir: %v", err)
					}
					if err := file.CopyDir(context.Background(), testIndexDir, backupDir); err != nil {
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
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
				name: "returns false when it's an initaial state",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: false,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(context.Background(), validIndexDir, tmpDir); err != nil {
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
					if err := file.CopyDir(context.Background(), validIndexDir, tmpDir); err != nil {
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
				name: "returns true when mets.IsInvalid is true",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: true,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(context.Background(), validIndexDir, tmpDir); err != nil {
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
				name: "returns true when mets.IsInvalid is true",
				args: args{
					path: tmpDir,
				},
				want: want{
					need: true,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.CopyDir(context.Background(), validIndexDir, tmpDir); err != nil {
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
					if err := file.CopyDir(context.Background(), validIndexDir, tmpDir); err != nil {
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
		name     string
		testfunc func(t *testing.T)
	}

	tests := []test{
		{
			"returns vector and timestamp when vector is found in vqueue",
			testReturnFromVq,
		},
		{
			"returns vector and timestamp when vector is found in kvs",
			testReturnFromKvs,
		},
		{
			"returns error when vector is not found in vector queue or kvs",
			testNotFoundInBothVqAndKvs,
		},
		{
			"returns error when vector is not found in vq found in kvs but also in delete queue",
			testFoundInBothIvqAndDvq,
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

	err = ngt.CreateIndex(context.Background(), 10)
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

	err = ngt.CreateIndex(context.Background(), 10)
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
		name string
		args args
		want error
	}

	setup := func(t *testing.T) string {
		tmpDir := t.TempDir()
		testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
		err := file.CopyDir(context.Background(), testIndexDir, tmpDir)
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

			err = ngt.CreateIndex(context.Background(), test.args.cfg.DefaultPoolSize)
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
		name string
		args args
		want error
	}

	setup := func(t *testing.T) string {
		tmpDir := t.TempDir()
		testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
		err := file.CopyDir(context.Background(), testIndexDir, tmpDir)
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

			err = ngt.SaveIndex(context.Background())
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
		name string
		args args
		want error
	}

	setup := func(t *testing.T) string {
		tmpDir := t.TempDir()
		testIndexDir := testdata.GetTestdataPath(testdata.ValidIndex)
		err := file.CopyDir(context.Background(), testIndexDir, tmpDir)
		require.NoError(t, err)

		return tmpDir
	}

	tests := []test{
		func() test {
			tmpDir := setup(t)
			return test{
				name: "Close returns ErrUncommittedIndexNotFound when it is not a read replica and try to Create Index because nothing has committed",
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

			err = ngt.Close(context.Background())
			require.Equal(tt, test.want, err)
		})
	}
}

func TestExportIndexInfo(t *testing.T) {
	t.Parallel()
	config := defaultConfig
	config.Dimension = 3
	config.EnableExportIndexInfoToK8s = true
	config.PodName = "test-pod"

	type test struct {
		name     string
		testfunc func(t *testing.T)
	}

	tests := []test{
		{
			"export after create index one vector",
			func(t *testing.T) {
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

				err = ngt.CreateIndex(context.Background(), 10)
				require.NoError(t, err)

				// expected entries
				expected := map[string]string{
					indexCountAnnotationsKey:                     "1",
					uncommittedAnnotationsKey:                    "0",
					unsavedCreateIndexExecutionNumAnnotationsKey: "1",
					unsavedProcessedVqAnnotationsKey:             "1",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 1)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expected)
			},
		},
		{
			"export after create index multiple vectors",
			func(t *testing.T) {
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

				err = ngt.CreateIndex(context.Background(), 10)
				require.NoError(t, err)

				// expected entries
				expected := map[string]string{
					indexCountAnnotationsKey:                     "2",
					uncommittedAnnotationsKey:                    "0",
					unsavedCreateIndexExecutionNumAnnotationsKey: "1",
					unsavedProcessedVqAnnotationsKey:             "2",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 1)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expected)
			},
		},
		{
			"export after create index multiple times",
			func(t *testing.T) {
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

				err = ngt.CreateIndex(context.Background(), 10)
				require.NoError(t, err)

				time2 := time.Now().UnixNano()
				err = ngt.InsertWithTime("test-uuid2", []float32{1.0, 2.0, 3.0}, time2)
				require.NoError(t, err)

				err = ngt.CreateIndex(context.Background(), 10)
				require.NoError(t, err)

				// expected entries
				expected := map[string]string{
					indexCountAnnotationsKey:                     "2",
					uncommittedAnnotationsKey:                    "0",
					unsavedCreateIndexExecutionNumAnnotationsKey: "2",
					unsavedProcessedVqAnnotationsKey:             "2",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 2)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expected)
			},
		},
		{
			"export after create index multiple vectors and save index",
			func(t *testing.T) {
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

				ctx := context.Background()
				err = ngt.CreateIndex(ctx, 10)
				require.NoError(t, err)

				// set time in context for testing
				saveIndexTime := time.Now()
				ctx = context.WithValue(ctx, saveIndexTimeKey, saveIndexTime)

				err = ngt.SaveIndex(ctx)
				require.NoError(t, err)

				// expected entries
				expectedAfterCreate := map[string]string{
					indexCountAnnotationsKey:                     "2",
					uncommittedAnnotationsKey:                    "0",
					unsavedCreateIndexExecutionNumAnnotationsKey: "1",
					unsavedProcessedVqAnnotationsKey:             "2",
				}
				expectedAfterSave := map[string]string{
					lastTimeSaveIndexTimestampAnnotationsKey:     saveIndexTime.UTC().Format(kvald.TimeFormat),
					unsavedCreateIndexExecutionNumAnnotationsKey: "0",
					unsavedProcessedVqAnnotationsKey:             "0",
				}
				// check mock called result
				mock.AssertExpectations(t)
				mock.AssertNumberOfCalls(t, "ApplyPodAnnotations", 2)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expectedAfterCreate)
				mock.AssertCalled(t, "ApplyPodAnnotations", testify.Anything, config.PodName, config.PodNamespace, expectedAfterSave)
			},
		},
		{
			"export after inserting vectors",
			func(t *testing.T) {
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

				ctx := context.Background()
				ExportMetricsOnTick(n.(*ngt), ctx)

				// expected entries
				expectedAfterInsert := map[string]string{
					indexCountAnnotationsKey:  "0",
					uncommittedAnnotationsKey: "2",
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
		idxes    []index
		poolSize uint32
		bulkSize int
	}
	type fields struct {
		svcCfg  *config.NGT
		svcOpts []Option

		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
					idxes: createRandomData(count, &createRandomDataConfig{
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

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
			for _, idx := range test.args.idxes {
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
			for i := 0; i < 100; i++ {
				idx := i
				eg.Go(safety.RecoverFunc(func() error {
					log.Warnf("started %d-1", idx)
					for _, idx := range test.args.idxes[:len(test.args.idxes)/3] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-1", idx)
					return nil
				}))

				eg.Go(safety.RecoverFunc(func() error {
					log.Warnf("started %d-2", idx)
					for _, idx := range test.args.idxes[len(test.args.idxes)/3 : 2*len(test.args.idxes)/3] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-2", idx)
					return nil
				}))

				eg.Go(safety.RecoverFunc(func() error {
					log.Warnf("started %d-3", idx)
					for _, idx := range test.args.idxes[2*len(test.args.idxes)/3:] {
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
		requests []*payload.Upsert_MultiRequest

		addr   string
		client grpc.Client
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		beforeFunc func(args)
		afterFunc  func(args)
	}
	multiUpsertRequestGenFunc := func(idxes []index, chunk int) (res []*payload.Upsert_MultiRequest) {
		reqs := make([]*payload.Upsert_Request, 0, chunk)
		for i := 0; i < len(idxes); i++ {
			if len(reqs) == chunk-1 {
				res = append(res, &payload.Upsert_MultiRequest{
					Requests: reqs,
				})
				reqs = make([]*payload.Upsert_Request, 0, chunk)
			} else {
				reqs = append(reqs, &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     idxes[i].uuid,
						Vector: idxes[i].vec,
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
				client: grpc.New(grpc.WithInsecure(true)),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			defer test.args.client.Close(ctx)

			for i := 0; i < 2; i++ {
				for _, req := range test.args.requests {
					_, err := test.args.client.Do(ctx, test.args.addr,
						func(ctx context.Context, conn *grpc.ClientConn, opts ...grpc.CallOption) (any, error) {
							return vald.NewValdClient(conn).MultiUpsert(ctx, req)
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		ctx     context.Context
// 		vec     []float32
// 		size    uint32
// 		epsilon float32
// 		radius  float32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
// 			}
//
// 			gotRes, err := n.Search(test.args.ctx, test.args.vec, test.args.size, test.args.epsilon, test.args.radius)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ngt_SearchByID(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		uuid    string
// 		size    uint32
// 		epsilon float32
// 		radius  float32
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		       },
// 		       fields: fields {
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           },
// 		           fields: fields {
// 		           core:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
// 			}
//
// 			gotVec, gotDst, err := n.SearchByID(test.args.ctx, test.args.uuid, test.args.size, test.args.epsilon, test.args.radius)
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// func Test_ngt_Delete(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// func Test_ngt_removeInvalidIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
// 			}
//
// 			err := n.readyForUpdate(test.args.uuid, test.args.vec)
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// func Test_ngt_toSearchResponse(t *testing.T) {
// 	type args struct {
// 		sr []algorithm.SearchResult
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// func Test_ngt_exportMetricsOnTick(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		core                    core.NGT
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
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
// 		eg                      errgroup.Group
// 		kvs                     kvs.BidiMap
// 		fmap                    map[string]int64
// 		vq                      vqueue.Queue
// 		indexing                atomic.Value
// 		flushing                atomic.Bool
// 		saving                  atomic.Value
// 		lastNocie               uint64
// 		nocie                   uint64
// 		nogce                   uint64
// 		wfci                    uint64
// 		nobic                   uint64
// 		nopvq                   atomic.Uint64
// 		cfg                     *config.NGT
// 		opts                    []Option
// 		inMem                   bool
// 		dim                     int
// 		alen                    int
// 		lim                     time.Duration
// 		dur                     time.Duration
// 		sdur                    time.Duration
// 		minLit                  time.Duration
// 		maxLit                  time.Duration
// 		litFactor               time.Duration
// 		enableProactiveGC       bool
// 		enableCopyOnWrite       bool
// 		podName                 string
// 		podNamespace            string
// 		path                    string
// 		tmpPath                 atomic.Value
// 		oldPath                 string
// 		basePath                string
// 		brokenPath              string
// 		poolSize                uint32
// 		radius                  float32
// 		epsilon                 float32
// 		idelay                  time.Duration
// 		dcd                     bool
// 		kvsdbConcurrency        int
// 		historyLimit            int
// 		isReadReplica           bool
// 		enableExportIndexInfo   bool
// 		exportIndexInfoDuration time.Duration
// 		patcher                 client.Patcher
// 		enableStatistics        bool
// 		statisticsCache         atomic.Pointer[payload.Info_Index_Statistics]
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 		           eg:nil,
// 		           kvs:nil,
// 		           fmap:nil,
// 		           vq:nil,
// 		           indexing:nil,
// 		           flushing:nil,
// 		           saving:nil,
// 		           lastNocie:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           nobic:0,
// 		           nopvq:nil,
// 		           cfg:nil,
// 		           opts:nil,
// 		           inMem:false,
// 		           dim:0,
// 		           alen:0,
// 		           lim:nil,
// 		           dur:nil,
// 		           sdur:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           enableProactiveGC:false,
// 		           enableCopyOnWrite:false,
// 		           podName:"",
// 		           podNamespace:"",
// 		           path:"",
// 		           tmpPath:nil,
// 		           oldPath:"",
// 		           basePath:"",
// 		           brokenPath:"",
// 		           poolSize:0,
// 		           radius:0,
// 		           epsilon:0,
// 		           idelay:nil,
// 		           dcd:false,
// 		           kvsdbConcurrency:0,
// 		           historyLimit:0,
// 		           isReadReplica:false,
// 		           enableExportIndexInfo:false,
// 		           exportIndexInfoDuration:nil,
// 		           patcher:nil,
// 		           enableStatistics:false,
// 		           statisticsCache:nil,
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
// 				eg:                      test.fields.eg,
// 				kvs:                     test.fields.kvs,
// 				fmap:                    test.fields.fmap,
// 				vq:                      test.fields.vq,
// 				indexing:                test.fields.indexing,
// 				flushing:                test.fields.flushing,
// 				saving:                  test.fields.saving,
// 				lastNocie:               test.fields.lastNocie,
// 				nocie:                   test.fields.nocie,
// 				nogce:                   test.fields.nogce,
// 				wfci:                    test.fields.wfci,
// 				nobic:                   test.fields.nobic,
// 				nopvq:                   test.fields.nopvq,
// 				cfg:                     test.fields.cfg,
// 				opts:                    test.fields.opts,
// 				inMem:                   test.fields.inMem,
// 				dim:                     test.fields.dim,
// 				alen:                    test.fields.alen,
// 				lim:                     test.fields.lim,
// 				dur:                     test.fields.dur,
// 				sdur:                    test.fields.sdur,
// 				minLit:                  test.fields.minLit,
// 				maxLit:                  test.fields.maxLit,
// 				litFactor:               test.fields.litFactor,
// 				enableProactiveGC:       test.fields.enableProactiveGC,
// 				enableCopyOnWrite:       test.fields.enableCopyOnWrite,
// 				podName:                 test.fields.podName,
// 				podNamespace:            test.fields.podNamespace,
// 				path:                    test.fields.path,
// 				tmpPath:                 test.fields.tmpPath,
// 				oldPath:                 test.fields.oldPath,
// 				basePath:                test.fields.basePath,
// 				brokenPath:              test.fields.brokenPath,
// 				poolSize:                test.fields.poolSize,
// 				radius:                  test.fields.radius,
// 				epsilon:                 test.fields.epsilon,
// 				idelay:                  test.fields.idelay,
// 				dcd:                     test.fields.dcd,
// 				kvsdbConcurrency:        test.fields.kvsdbConcurrency,
// 				historyLimit:            test.fields.historyLimit,
// 				isReadReplica:           test.fields.isReadReplica,
// 				enableExportIndexInfo:   test.fields.enableExportIndexInfo,
// 				exportIndexInfoDuration: test.fields.exportIndexInfoDuration,
// 				patcher:                 test.fields.patcher,
// 				enableStatistics:        test.fields.enableStatistics,
// 				statisticsCache:         test.fields.statisticsCache,
// 			}
//
// 			err := n.exportMetricsOnSaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
