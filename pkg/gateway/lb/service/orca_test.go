//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_orca_order(t *testing.T) {
	t.Parallel()
	loads := map[string]resourceLoad{
		"agent-1": {cpu: 0.2, memory: 0.3},
		"agent-2": {cpu: 0.6, memory: 0.5},
		"agent-3": {cpu: 0.95, memory: 0.7},
	}
	tests := []struct {
		snapshot *resourceLoadSnapshot
		want     []string
		name     string
		addrs    []string
	}{
		{
			name: "orders all known agents by load",
			snapshot: &resourceLoadSnapshot{
				collectedAt: time.Now(),
				loads:       loads,
			},
			addrs: []string{"agent-3", "agent-2", "agent-1"},
			want:  []string{"agent-1", "agent-2", "agent-3"},
		},
		{
			name: "keeps unknown agents as fallback",
			snapshot: &resourceLoadSnapshot{
				collectedAt: time.Now(),
				loads:       loads,
			},
			addrs: []string{"unknown-1", "agent-2", "unknown-2", "agent-1"},
			want:  []string{"agent-1", "agent-2", "unknown-1", "unknown-2"},
		},
		{
			name: "stale report disables control",
			snapshot: &resourceLoadSnapshot{
				collectedAt: time.Now().Add(-time.Minute),
				loads:       loads,
			},
			addrs: []string{"agent-1", "agent-2"},
			want:  nil,
		},
		{
			name:  "missing report disables control",
			addrs: []string{"agent-1", "agent-2"},
			want:  nil,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
			o := newORCA(time.Second, 3*time.Second, 3)
			if test.snapshot != nil {
				o.snapshot.Store(test.snapshot)
			}
			if got := o.order(test.addrs); !reflect.DeepEqual(got, test.want) {
				t.Errorf("order() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_orca_read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		addrs   []string
		replica int
		want    []string
	}{
		{
			name:    "selects N minus replica plus one agents",
			addrs:   []string{"agent-1", "agent-2", "agent-3", "agent-4", "agent-5"},
			replica: 3,
			want:    []string{"agent-1", "agent-2", "agent-3"},
		},
		{
			name:    "selects at least one agent",
			addrs:   []string{"agent-1", "agent-2"},
			replica: 3,
			want:    []string{"agent-1"},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
			loads := make(map[string]resourceLoad, len(test.addrs))
			for idx, addr := range test.addrs {
				loads[addr] = resourceLoad{cpu: float64(idx) / 10}
			}
			o := newORCA(time.Second, 3*time.Second, test.replica)
			o.snapshot.Store(&resourceLoadSnapshot{
				collectedAt: time.Now(),
				loads:       loads,
			})
			if got := o.read(test.addrs); !reflect.DeepEqual(got, test.want) {
				t.Errorf("read() = %v, want %v", got, test.want)
			}
		})
	}
}
