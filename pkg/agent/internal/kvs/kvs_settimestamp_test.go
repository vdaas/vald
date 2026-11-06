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
package kvs

import "testing"

// Test_bidi_SetTimestamp pins that SetTimestamp only changes the timestamp of a pair that
// still maps key to val, and never creates or repoints a pair the way Set does.
func Test_bidi_SetTimestamp(t *testing.T) {
	t.Parallel()

	const (
		key = "uuid"
		val = uint32(7)
		ts  = int64(10)
	)
	tests := []struct {
		setup  func(b BidiMap)
		name   string
		wantOk bool
		// wantVal and wantTs are the pair of key afterwards; wantVal == 0 means no pair.
		wantVal uint32
		wantTs  int64
	}{
		{
			name:    "updates the timestamp of the mapped pair",
			setup:   func(b BidiMap) { b.Set(key, val, 1) },
			wantOk:  true,
			wantVal: val,
			wantTs:  ts,
		},
		{
			name:  "does not re-create a deleted pair",
			setup: func(b BidiMap) { b.Set(key, val, 1); b.Delete(key) },
		},
		{
			name:    "does not touch a pair repointed to another val",
			setup:   func(b BidiMap) { b.Set(key, val, 1); b.Set(key, val+1, 2) },
			wantVal: val + 1,
			wantTs:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := New()
			tt.setup(b)

			if got := b.SetTimestamp(key, val, ts); got != tt.wantOk {
				t.Fatalf("SetTimestamp() = %v, want %v", got, tt.wantOk)
			}
			gotVal, gotTs, ok := b.Get(key)
			if ok != (tt.wantVal != 0) || gotVal != tt.wantVal || gotTs != tt.wantTs {
				t.Errorf("Get() = (%d, %d, %v), want (%d, %d)", gotVal, gotTs, ok, tt.wantVal, tt.wantTs)
			}
			wantLen := uint64(0)
			if tt.wantVal != 0 {
				wantLen = 1
			}
			if b.Len() != wantLen {
				t.Errorf("Len() = %d, want %d", b.Len(), wantLen)
			}
			// the inverse pair of val exists only while key still maps to it.
			gotKey, invTs, ok := b.GetInverse(val)
			if ok != (tt.wantVal == val) {
				t.Fatalf("GetInverse(%d) = (%q, %d, %v), want exists %v", val, gotKey, invTs, ok, tt.wantVal == val)
			}
			if ok && (gotKey != key || invTs != ts) {
				t.Errorf("GetInverse(%d) = (%q, %d), want (%q, %d)", val, gotKey, invTs, key, ts)
			}
		})
	}
}
