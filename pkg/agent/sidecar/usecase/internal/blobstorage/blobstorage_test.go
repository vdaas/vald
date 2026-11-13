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
package blobstorage

import (
	"testing"

	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/http/client"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/sidecar/config"
)

func newConfig(t *testing.T, algorithm string) *config.Data {
	t.Helper()

	// AgentSidecar.Bind only binds Client and BlobStorage when they are already
	// non-nil, and New dereferences Client.Net/Transport and BlobStorage.S3/CloudStorage.
	cfg := &config.Data{
		AgentSidecar: &iconf.AgentSidecar{
			Filename:       "backup",
			FilenameSuffix: ".tar.gz",
			Client:         new(iconf.Client),
			BlobStorage:    &iconf.Blob{StorageType: "s3", Bucket: "bucket"},
			Compress:       &iconf.CompressCore{CompressAlgorithm: algorithm},
		},
	}
	cfg.AgentSidecar = cfg.AgentSidecar.Bind()
	return cfg
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		algorithm  string
		clientOpts []client.Option
		wantErr    bool
	}{
		{name: "initcontainer wiring without extra client options", algorithm: "gzip"},
		{name: "sidecar wiring with the transport backoff", algorithm: "gzip", clientOpts: []client.Option{client.WithBackoffOpts(new(iconf.Backoff).Opts()...)}},
		{name: "unknown compressor is rejected", algorithm: "unknown", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bs, err := New(newConfig(t, tt.algorithm), errgroup.Get(), tt.clientOpts...)
			if tt.wantErr {
				if err == nil {
					t.Fatal("New() error = nil, want an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			info := bs.StorageInfo()
			if info.BucketName != "bucket" || info.Filename != "backup.tar.gz" {
				t.Errorf("StorageInfo() = %+v, want bucket %q and filename %q", info, "bucket", "backup.tar.gz")
			}
		})
	}
}
