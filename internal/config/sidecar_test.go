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

// Package config providers configuration type and load configuration logic
package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestAgentSidecar_Bind(t *testing.T) {
	type fields struct {
		Mode               string
		WatchDir           string
		AutoBackupDuration string
		PostStopTimeout    string
		Filename           string
		FilenameSuffix     string
		BlobStorage        *Blob
		Compress           *CompressCore
		RestoreBackoff     *Backoff
		Client             *Client
	}
	type want struct {
		want *AgentSidecar
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *AgentSidecar) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *AgentSidecar) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			mode := "sidecar"
			watchDir := "sidecar"
			autoBackupDuration := "10ms"
			postStopTimeout := "5m"
			filename := "vald-ngt-1"
			filenameSuffix := "tar.gz"
			blobStorageType := "s3"
			compressAlgorithm := GOB.String()
			backoffInitialDuration := "10ms"
			return test{
				name: "return AgentSidecar when all of object are set",
				fields: fields{
					Mode:               mode,
					WatchDir:           watchDir,
					AutoBackupDuration: autoBackupDuration,
					PostStopTimeout:    postStopTimeout,
					Filename:           filename,
					FilenameSuffix:     filenameSuffix,
					BlobStorage: &Blob{
						StorageType: blobStorageType,
					},
					Compress: &CompressCore{
						CompressAlgorithm: compressAlgorithm,
					},
					RestoreBackoff: &Backoff{
						InitialDuration: backoffInitialDuration,
					},
					Client: &Client{
						Net: new(Net),
					},
				},
				want: want{
					want: &AgentSidecar{
						Mode:               mode,
						WatchDir:           watchDir,
						AutoBackupDuration: autoBackupDuration,
						PostStopTimeout:    postStopTimeout,
						Filename:           filename,
						FilenameSuffix:     filenameSuffix,
						BlobStorage: &Blob{
							StorageType: blobStorageType,
							S3:          new(S3Config),
						},
						Compress: &CompressCore{
							CompressAlgorithm: compressAlgorithm,
						},
						RestoreBackoff: &Backoff{
							InitialDuration: backoffInitialDuration,
						},
						Client: &Client{
							Net: new(Net),
						},
					},
				},
			}
		}(),
		func() test {
			mode := "sidecar"
			watchDir := "sidecar"
			autoBackupDuration := "10ms"
			postStopTimeout := "5m"
			filename := "vald-ngt-1"
			filenameSuffix := "tar.gz"
			return test{
				name: "return AgentSidecar when all of object are not set",
				fields: fields{
					Mode:               mode,
					WatchDir:           watchDir,
					AutoBackupDuration: autoBackupDuration,
					PostStopTimeout:    postStopTimeout,
					Filename:           filename,
					FilenameSuffix:     filenameSuffix,
				},
				want: want{
					want: &AgentSidecar{
						Mode:               mode,
						WatchDir:           watchDir,
						AutoBackupDuration: autoBackupDuration,
						PostStopTimeout:    postStopTimeout,
						Filename:           filename,
						FilenameSuffix:     filenameSuffix,
						BlobStorage:        new(Blob),
						Compress:           new(CompressCore),
						RestoreBackoff:     new(Backoff),
						Client:             new(Client),
					},
				},
			}
		}(),
		func() test {
			mode := "sidecar"
			watchDir := "sidecar"
			autoBackupDuration := "10ms"
			postStopTimeout := "5m"
			filename := "vald-ngt-1"
			filenameSuffix := "tar.gz"
			m := map[string]string{
				"MODE":                 mode,
				"WATCH_DIR":            watchDir,
				"AUTO_BACKUP_DURATION": autoBackupDuration,
				"POST_STOP_TIMEOUT":    postStopTimeout,
				"FILENAME":             filename,
				"FILENAME_SUFFIX":      filenameSuffix,
			}
			return test{
				name: "return AgentSidecar when the data is loaded from the environment variable",
				fields: fields{
					Mode:               "_MODE_",
					WatchDir:           "_WATCH_DIR_",
					AutoBackupDuration: "_AUTO_BACKUP_DURATION_",
					PostStopTimeout:    "_POST_STOP_TIMEOUT_",
					Filename:           "_FILENAME_",
					FilenameSuffix:     "_FILENAME_SUFFIX_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range m {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &AgentSidecar{
						Mode:               mode,
						WatchDir:           watchDir,
						AutoBackupDuration: autoBackupDuration,
						PostStopTimeout:    postStopTimeout,
						Filename:           filename,
						FilenameSuffix:     filenameSuffix,
						BlobStorage:        new(Blob),
						Compress:           new(CompressCore),
						RestoreBackoff:     new(Backoff),
						Client:             new(Client),
					},
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &AgentSidecar{
				Mode:               test.fields.Mode,
				WatchDir:           test.fields.WatchDir,
				AutoBackupDuration: test.fields.AutoBackupDuration,
				PostStopTimeout:    test.fields.PostStopTimeout,
				Filename:           test.fields.Filename,
				FilenameSuffix:     test.fields.FilenameSuffix,
				BlobStorage:        test.fields.BlobStorage,
				Compress:           test.fields.Compress,
				RestoreBackoff:     test.fields.RestoreBackoff,
				Client:             test.fields.Client,
			}

			got := s.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
