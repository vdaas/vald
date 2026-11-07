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

package grpc

import (
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func TestWithSearchQuorumGrace(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dur  string
		want time.Duration
	}{
		{name: "empty keeps it disabled", dur: ""},
		{name: "invalid keeps it disabled", dur: "soon"},
		{name: "zero keeps it disabled", dur: "0s"},
		{name: "negative keeps it disabled", dur: "-10ms"},
		{name: "positive enables it", dur: "10ms", want: 10 * time.Millisecond},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := new(server)
			WithSearchQuorumGrace(tc.dur)(s)
			if s.searchQuorumGrace != tc.want {
				t.Errorf("searchQuorumGrace = %v, want %v", s.searchQuorumGrace, tc.want)
			}
		})
	}
}

func Test_server_searchQuorum(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		grace   time.Duration
		replica int
		agents  int
		want    int
	}{
		{name: "disabled waits for every agent", replica: 3, agents: 6},
		{name: "no replication waits for every agent", grace: time.Millisecond, replica: 1, agents: 6},
		{name: "agents not above replica waits for every agent", grace: time.Millisecond, replica: 3, agents: 3},
		{name: "6 agents with replica 3 need 4", grace: time.Millisecond, replica: 3, agents: 6, want: 4},
		{name: "20 agents with replica 2 need 19", grace: time.Millisecond, replica: 2, agents: 20, want: 19},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := &server{searchQuorumGrace: tc.grace, replica: tc.replica}
			if got := s.searchQuorum(tc.agents); got != tc.want {
				t.Errorf("searchQuorum(%d) = %d, want %d", tc.agents, got, tc.want)
			}
		})
	}
}

// Test_server_Search_quorumKeepsQuorumResults verifies that an early-terminated
// search succeeds with the results of the answered agents and does not wait
// for the straggler.
func Test_server_Search_quorumKeepsQuorumResults(t *testing.T) {
	t.Parallel()
	const num = 10
	s := newStragglerServer(6, 3, time.Millisecond, 5*time.Second)
	WithSearchQuorumGrace("10ms")(s)
	start := time.Now()
	res, err := s.Search(t.Context(), &payload.Search_Request{
		Vector: make([]float32, 8),
		Config: &payload.Search_Config{Num: num, Timeout: int64(10 * time.Second)},
	})
	if err != nil {
		t.Fatalf("Search() err = %v", err)
	}
	if got := len(res.GetResults()); got != num {
		t.Errorf("Search() returned %d results, want %d", got, num)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Search() took %v, want it to stop waiting for the 5s straggler", elapsed)
	}
}
