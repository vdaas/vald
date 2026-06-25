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

func Test_orca_selectAddrs(t *testing.T) {
	t.Parallel()
	loads := map[string]resourceLoad{
		"agent-1": {cpu: 0.2, memory: 0.3},
		"agent-2": {cpu: 0.6, memory: 0.5},
		"agent-3": {cpu: 0.95, memory: 0.7},
	}
	tests := []struct {
		snapshot *resourceLoadSnapshot
		read     ORCAPolicy
		write    ORCAPolicy
		want     []string
		name     string
		kind     BroadCastKind
		addrs    []string
	}{
		{
			name: "read policy limits fanout and rejects overloaded agents",
			snapshot: &resourceLoadSnapshot{
				collectedAt: time.Now(),
				loads:       loads,
			},
			read:  ORCAPolicy{MinFanout: 1, MaxFanout: 2, CPUThreshold: 0.9, MemoryThreshold: 0.9},
			write: ORCAPolicy{MinFanout: 3, CPUThreshold: 0.5, MemoryThreshold: 0.9},
			kind:  READ,
			addrs: []string{"agent-3", "agent-2", "agent-1"},
			want:  []string{"agent-1", "agent-2"},
		},
		{
			name: "write policy backfills to minimum fanout",
			snapshot: &resourceLoadSnapshot{
				collectedAt: time.Now(),
				loads:       loads,
			},
			read:  ORCAPolicy{MinFanout: 1, MaxFanout: 1, CPUThreshold: 0.9, MemoryThreshold: 0.9},
			write: ORCAPolicy{MinFanout: 2, CPUThreshold: 0.5, MemoryThreshold: 0.9},
			kind:  WRITE,
			addrs: []string{"agent-3", "agent-2", "agent-1"},
			want:  []string{"agent-1", "agent-2"},
		},
		{
			name: "unknown agent remains eligible",
			snapshot: &resourceLoadSnapshot{
				collectedAt: time.Now(),
				loads:       loads,
			},
			read:  ORCAPolicy{MinFanout: 1, MaxFanout: 2, CPUThreshold: 0.5, MemoryThreshold: 0.9},
			write: ORCAPolicy{MinFanout: 1},
			kind:  READ,
			addrs: []string{"agent-3", "agent-unknown", "agent-1"},
			want:  []string{"agent-1", "agent-unknown"},
		},
		{
			name: "stale report disables control",
			snapshot: &resourceLoadSnapshot{
				collectedAt: time.Now().Add(-time.Minute),
				loads:       loads,
			},
			read:  ORCAPolicy{MinFanout: 1, MaxFanout: 1},
			write: ORCAPolicy{MinFanout: 1, MaxFanout: 1},
			kind:  READ,
			addrs: []string{"agent-1", "agent-2"},
			want:  nil,
		},
		{
			name:  "missing report disables control",
			read:  ORCAPolicy{MinFanout: 1, MaxFanout: 1},
			write: ORCAPolicy{MinFanout: 1, MaxFanout: 1},
			kind:  READ,
			addrs: []string{"agent-1", "agent-2"},
			want:  nil,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
			o := newORCA(time.Second, 3*time.Second, test.read, test.write)
			if test.snapshot != nil {
				o.snapshot.Store(test.snapshot)
			}
			if got := o.selectAddrs(test.kind, test.addrs); !reflect.DeepEqual(got, test.want) {
				t.Errorf("selectAddrs() = %v, want %v", got, test.want)
			}
		})
	}
}
