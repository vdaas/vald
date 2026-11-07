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

package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

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
