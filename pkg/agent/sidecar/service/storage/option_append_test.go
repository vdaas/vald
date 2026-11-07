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
package storage

import (
	"testing"

	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage"
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage/urlopener"
)

// TestOptions_appendOnce pins that the slice options append each given option
// exactly once, also on the first call (when the slice is still nil).
func TestOptions_appendOnce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		count func(*bs) int
		name  string
		calls []Option
		want  int
	}{
		{
			name:  "WithCloudStorageOpts first call",
			calls: []Option{WithCloudStorageOpts(cloudstorage.WithURL("gs://a"))},
			count: func(b *bs) int { return len(b.cloudStorageOpts) },
			want:  1,
		},
		{
			name: "WithCloudStorageOpts two calls",
			calls: []Option{
				WithCloudStorageOpts(cloudstorage.WithURL("gs://a")),
				WithCloudStorageOpts(cloudstorage.WithURL("gs://b")),
			},
			count: func(b *bs) int { return len(b.cloudStorageOpts) },
			want:  2,
		},
		{
			name:  "WithCloudStorageURLOpenerOpts first call",
			calls: []Option{WithCloudStorageURLOpenerOpts(urlopener.WithCredentialsJSON("{}"))},
			count: func(b *bs) int { return len(b.cloudStorageURLOpenerOpts) },
			want:  1,
		},
		{
			name: "WithCloudStorageURLOpenerOpts two calls",
			calls: []Option{
				WithCloudStorageURLOpenerOpts(urlopener.WithCredentialsJSON("{}")),
				WithCloudStorageURLOpenerOpts(urlopener.WithCredentialsFile("f")),
			},
			count: func(b *bs) int { return len(b.cloudStorageURLOpenerOpts) },
			want:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := new(bs)
			for _, opt := range tt.calls {
				if err := opt(b); err != nil {
					t.Fatalf("option returned error: %v", err)
				}
			}
			if got := tt.count(b); got != tt.want {
				t.Errorf("len = %d, want %d", got, tt.want)
			}
		})
	}
}
