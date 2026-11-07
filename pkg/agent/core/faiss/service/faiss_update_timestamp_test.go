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
	"math"
	"slices"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

// update deletes at t and inserts at t+1. t == MaxInt64 used to wrap the
// insert to MinInt64, which the queue drops behind the MaxInt64 delete: the
// update reported success and the object disappeared. It must be rejected
// before the delete; ordinary timestamps keep the delete-at-t / insert-at-t+1.
func Test_faiss_update_maxTimestamp(t *testing.T) {
	t.Parallel()
	const (
		uuid = "update-ts-uuid"
		oldT = int64(100)
	)
	oldVec, newVec := []float32{9}, []float32{1}
	tests := []struct {
		name    string
		ts      int64
		wantErr bool
	}{
		{name: "rejects MaxInt64", ts: math.MaxInt64, wantErr: true},
		{name: "keeps a normal timestamp", ts: 200},
		{name: "keeps MaxInt64-1", ts: math.MaxInt64 - 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = vq.PushInsert(uuid, oldVec, oldT); err != nil {
				t.Fatal(err)
			}
			f := &faiss{dim: 1, kvs: kvs.New(kvs.WithConcurrency(1)), vq: vq}
			defer closeKVS(t, f)

			err = f.update(uuid, newVec, test.ts)
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
