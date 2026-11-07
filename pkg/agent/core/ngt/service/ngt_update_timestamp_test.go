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
