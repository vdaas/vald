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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *AgentSidecar) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				Mode:               "sidecar",
				WatchDir:           "/var/index",
				AutoBackupDuration: "10ms",
				PostStopTimeout:    "5m",
				Filename:           "vald-ngt-1",
				FilenameSuffix:     "tar.gz",
				BlobStorage: &Blob{
					StorageType: "s3",
				},
				Compress: &CompressCore{
					CompressAlgorithm: GOB.String(),
				},
				RestoreBackoff: &Backoff{
					InitialDuration: "10ms",
				},
				Client: &Client{
					Net: new(Net),
				},
			}
			return test{
				name:   "return AgentSidecar when all of object are set",
				fields: fields,
				want: want{
					want: &AgentSidecar{
						Mode: fields.Mode,
						// WatchDir:           "/var/index",
						// AutoBackupDuration: "10ms",
						// PostStopTimeout:    "5m",
						// Filename:           "vald-ngt-1",
						// FilenameSuffix:     "tar.gz",
						// BlobStorage: &Blob{
						// 	StorageType: "s3",
						// },
						// Compress: &CompressCore{
						// 	CompressAlgorithm: GOB.String(),
						// },
						// RestoreBackoff: &Backoff{
						// 	InitialDuration: "10ms",
						// },
						// Client: &Client{
						// 	Net: new(Net),
						// },
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
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
