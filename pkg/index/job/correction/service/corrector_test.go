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
	"context"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/db/kvs/pogreb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// fakeAgent is an in-memory agent served over bufconn. Update mimics the agent:
// an identical vector only refreshes an older timestamp and returns
// AlreadyExists, otherwise the vector is stored at Config.Timestamp+1.
type fakeAgent struct {
	vald.UnimplementedValdServer

	objects   map[string]*payload.Object_Vector
	tsErr     error // returned by GetTimestamp when set
	streamErr error // returned after listing the objects when set
	insertErr error // returned by Insert when set
	removeErr error // returned by Remove when set

	updates   []*payload.Update_Request
	updateTSs []*payload.Update_TimestampRequest
	inserts   []*payload.Insert_Request

	mu        sync.Mutex
	tsErrOnce bool // when set, tsErr is returned only for the first GetTimestamp call
}

func (a *fakeAgent) StreamListObject(
	_ *payload.Object_List_Request, stream vald.Object_StreamListObjectServer,
) error {
	a.mu.Lock()
	ids := make([]string, 0, len(a.objects))
	for id := range a.objects {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	vecs := make([]*payload.Object_Vector, 0, len(ids))
	for _, id := range ids {
		vecs = append(vecs, a.objects[id].CloneVT())
	}
	serr := a.streamErr
	a.mu.Unlock()
	for _, vec := range vecs {
		if err := stream.Send(&payload.Object_List_Response{
			Payload: &payload.Object_List_Response_Vector{Vector: vec},
		}); err != nil {
			return err
		}
	}
	return serr
}

func (a *fakeAgent) GetTimestamp(
	_ context.Context, req *payload.Object_TimestampRequest,
) (*payload.Object_Timestamp, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.tsErr != nil {
		err := a.tsErr
		if a.tsErrOnce {
			a.tsErr = nil
		}
		return nil, err
	}
	obj, ok := a.objects[req.GetId().GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return &payload.Object_Timestamp{Id: obj.GetId(), Timestamp: obj.GetTimestamp()}, nil
}

func (a *fakeAgent) GetObject(
	_ context.Context, req *payload.Object_VectorRequest,
) (*payload.Object_Vector, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	obj, ok := a.objects[req.GetId().GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return obj.CloneVT(), nil
}

func (a *fakeAgent) Update(
	_ context.Context, req *payload.Update_Request,
) (*payload.Object_Location, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.updates = append(a.updates, req.CloneVT())
	id := req.GetVector().GetId()
	obj, ok := a.objects[id]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	ts := req.GetConfig().GetTimestamp()
	if slices.Equal(obj.GetVector(), req.GetVector().GetVector()) {
		if obj.GetTimestamp() < ts {
			obj.Timestamp = ts
		}
		return nil, status.Error(codes.AlreadyExists, "same data already exists")
	}
	a.objects[id] = &payload.Object_Vector{Id: id, Vector: slices.Clone(req.GetVector().GetVector()), Timestamp: ts + 1}
	return new(payload.Object_Location), nil
}

func (a *fakeAgent) UpdateTimestamp(
	_ context.Context, req *payload.Update_TimestampRequest,
) (*payload.Object_Location, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.updateTSs = append(a.updateTSs, req.CloneVT())
	obj, ok := a.objects[req.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	obj.Timestamp = req.GetTimestamp()
	return new(payload.Object_Location), nil
}

func (a *fakeAgent) Insert(
	_ context.Context, req *payload.Insert_Request,
) (*payload.Object_Location, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.inserts = append(a.inserts, req.CloneVT())
	if a.insertErr != nil {
		return nil, a.insertErr
	}
	a.objects[req.GetVector().GetId()] = req.GetVector().CloneVT()
	return new(payload.Object_Location), nil
}

func (a *fakeAgent) Remove(
	_ context.Context, req *payload.Remove_Request,
) (*payload.Object_Location, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.removeErr != nil {
		return nil, a.removeErr
	}
	id := req.GetId().GetId()
	if _, ok := a.objects[id]; !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	delete(a.objects, id)
	return new(payload.Object_Location), nil
}

func (a *fakeAgent) object(id string) *payload.Object_Vector {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.objects[id].CloneVT()
}

// fakeGRPCClient mirrors the error contract of the real client: OrderedRange
// and OrderedRangeConcurrent never return the callback errors.
type fakeGRPCClient struct {
	grpc.Client
	conns map[string]*grpc.ClientConn
	// unreachable addrs are skipped by OrderedRangeConcurrent without running
	// the callback and without an error, like the real client does for an addr
	// with no connection pool, a failed connection lookup or an open breaker.
	unreachable map[string]bool
	errs        []error // callback errors, swallowed like the real client does
	mu          sync.Mutex
}

func (f *fakeGRPCClient) OrderedRange(
	ctx context.Context,
	order []string,
	fn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	for _, addr := range order {
		if conn, ok := f.conns[addr]; ok {
			if err := fn(ctx, addr, conn); err != nil {
				f.mu.Lock()
				f.errs = append(f.errs, err)
				f.mu.Unlock()
			}
		}
	}
	return nil
}

// OrderedRangeConcurrent re-runs a failing callback like the real client's backoff
// (g.bo.Do in gRPCClient.do) before swallowing the final error.
func (f *fakeGRPCClient) OrderedRangeConcurrent(
	ctx context.Context,
	order []string,
	_ int,
	fn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	const attempts = 3
	order = slices.DeleteFunc(slices.Clone(order), func(addr string) bool { return f.unreachable[addr] })
	return f.OrderedRange(ctx, order, func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
		for range attempts {
			if err = fn(ctx, addr, conn, copts...); err == nil {
				return nil
			}
		}
		return err
	})
}

func (f *fakeGRPCClient) Do(
	ctx context.Context,
	addr string,
	fn func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error),
) (any, error) {
	conn, ok := f.conns[addr]
	if !ok {
		return nil, errors.ErrGRPCClientConnNotFound(addr)
	}
	return fn(ctx, conn)
}

type fakeDiscoverer struct {
	discoverer.Client
	client *fakeGRPCClient
	addrs  []string
}

func (f *fakeDiscoverer) GetClient() grpc.Client            { return f.client }
func (f *fakeDiscoverer) GetAddrs(context.Context) []string { return f.addrs }

type fakeGateway struct {
	vc.Client
	counts map[string]*payload.Info_Index_Count
}

func (f *fakeGateway) IndexDetail(
	context.Context, *payload.Empty, ...grpc.CallOption,
) (*payload.Info_Index_Detail, error) {
	return &payload.Info_Index_Detail{Counts: f.counts}, nil
}

type memKVS struct {
	pogreb.DB
	m  map[string][]byte
	mu sync.Mutex
}

func (m *memKVS) Get(key string) ([]byte, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.m[key]
	return v, ok, nil
}

func (m *memKVS) Set(key string, val []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = val
	return nil
}

// newTestCorrector serves every agent over bufconn. Agents are processed in
// the order given (stored counts are assigned in descending order). addrs is
// what the discoverer reports for shortage/oversupply correction.
func newTestCorrector(
	t *testing.T, indexReplica int, order []string, agents map[string]*fakeAgent, addrs []string,
) *correct {
	t.Helper()
	conns := make(map[string]*grpc.ClientConn, len(agents))
	counts := make(map[string]*payload.Info_Index_Count, len(agents))
	const maxStored = 1 << 10
	stored := uint32(maxStored)
	for _, addr := range order {
		lis := bufconn.Listen(1 << 20)
		srv := ggrpc.NewServer()
		vald.RegisterValdServer(srv, agents[addr])
		served := make(chan struct{})
		go func() {
			defer close(served)
			if err := srv.Serve(lis); err != nil {
				t.Logf("fake agent %s stopped: %v", addr, err)
			}
		}()
		conn, err := ggrpc.NewClient("passthrough:///"+addr,
			ggrpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
				return lis.DialContext(ctx)
			}),
			ggrpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("failed to create client for %s: %v", addr, err)
		}
		t.Cleanup(func() {
			if err := conn.Close(); err != nil {
				t.Logf("failed to close conn to %s: %v", addr, err)
			}
			srv.Stop()
			select {
			case <-served:
			case <-time.After(10 * time.Second):
				t.Errorf("fake agent %s did not stop", addr)
			}
		})
		conns[addr] = conn
		counts[addr] = &payload.Info_Index_Count{Stored: stored}
		stored--
	}
	return &correct{
		discoverer:            &fakeDiscoverer{client: &fakeGRPCClient{conns: conns}, addrs: addrs},
		gateway:               &fakeGateway{counts: counts},
		checkedList:           &memKVS{m: make(map[string][]byte)},
		indexReplica:          indexReplica,
		streamListConcurrency: 2,
	}
}

func startCorrector(t *testing.T, c *correct) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	return c.Start(ctx)
}

func newVec(ts int64, v ...float32) *payload.Object_Vector {
	return &payload.Object_Vector{Id: "x", Vector: v, Timestamp: ts}
}

func TestStart_ReportsFailures(t *testing.T) {
	t.Parallel()
	tests := []struct {
		agent     *fakeAgent
		name      string
		errSubstr string
		addrs     []string
		replica   int
		wantErr   bool
	}{
		{
			name:    "returns nil when the stream ends with EOF",
			replica: 1,
			agent:   &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}},
			addrs:   []string{"a0"},
		},
		{
			name:    "returns error when the stream ends with a real failure",
			replica: 1,
			agent: &fakeAgent{
				objects:   map[string]*payload.Object_Vector{"x": newVec(100, 1)},
				streamErr: status.Error(codes.Unavailable, "agent restarted"),
			},
			addrs:     []string{"a0"},
			wantErr:   true,
			errSubstr: "agent restarted",
		},
		{
			name:    "returns error when a worker fails",
			replica: 2,
			agent:   &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}},
			// no agent to insert the missing replica into
			addrs:     nil,
			wantErr:   true,
			errSubstr: errors.ErrNoAvailableAgentToInsert.Error(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			c := newTestCorrector(tt, tc.replica, []string{"a0"}, map[string]*fakeAgent{"a0": tc.agent}, tc.addrs)
			err := startCorrector(tt, c)
			if !tc.wantErr {
				if err != nil {
					tt.Errorf("Start() error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				tt.Fatal("Start() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), tc.errSubstr) {
				tt.Errorf("Start() error = %v, want it to contain %q", err, tc.errSubstr)
			}
		})
	}
}

func TestStart_CorrectsOriginTimestamp(t *testing.T) {
	t.Parallel()
	// a0 is processed first and holds an older copy than its replica a1.
	a0 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	a1 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(200, 2)}}
	c := newTestCorrector(t, 2, []string{"a0", "a1"}, map[string]*fakeAgent{"a0": a0, "a1": a1}, []string{"a0", "a1"})
	if err := startCorrector(t, c); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	// correctTimestamp only moves the timestamp forward (UpdateTimestamp); replacing the vector
	// content of a stale replica is a deferred design item.
	if got := a0.object("x"); got.GetTimestamp() != 200 {
		t.Errorf("origin replica = %v, want timestamp 200", got)
	}
	if got := a1.object("x"); got.GetTimestamp() != 200 || !slices.Equal(got.GetVector(), []float32{2}) {
		t.Errorf("latest replica changed: %v", got)
	}
	for addr, a := range map[string]*fakeAgent{"a0": a0, "a1": a1} {
		a.mu.Lock()
		n := len(a.inserts)
		a.mu.Unlock()
		if n != 0 {
			t.Errorf("no replica is missing, but %d inserts happened on %s", n, addr)
		}
	}
}

func TestStart_SkipsIdWithIncompleteReplicaInfo(t *testing.T) {
	t.Parallel()
	// x has 2 of 3 replicas (a0, a1), but a2 fails GetTimestamp, so whether a2
	// holds x is unknown and nothing may be inserted.
	a0 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	a1 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	a2 := &fakeAgent{
		objects: map[string]*payload.Object_Vector{},
		tsErr:   status.Error(codes.Unavailable, "temporarily unavailable"),
	}
	c := newTestCorrector(t, 3, []string{"a0", "a1", "a2"},
		map[string]*fakeAgent{"a0": a0, "a1": a1, "a2": a2}, []string{"a0", "a1", "a2"})
	// x is skipped, so it stays uncorrected and the job must not report success.
	err := startCorrector(t, c)
	if err == nil || !strings.Contains(err.Error(), "incomplete replica information") {
		t.Fatalf("Start() error = %v, want an incomplete replica information error", err)
	}
	for addr, a := range map[string]*fakeAgent{"a0": a0, "a1": a1, "a2": a2} {
		a.mu.Lock()
		n := len(a.inserts) + len(a.updates) + len(a.updateTSs)
		a.mu.Unlock()
		if n != 0 {
			t.Errorf("agent %s was mutated (%d calls) although the replica info of x is incomplete", addr, n)
		}
	}
}

// TestStart_UnreachedReplicaIsIncomplete pins that a replica agent whose
// callback never runs (no pool, conn not found, breaker open) makes the replica
// info incomplete instead of being counted as not holding the id.
func TestStart_UnreachedReplicaIsIncomplete(t *testing.T) {
	t.Parallel()
	// x is on a0 and a1; a2 cannot be reached during the replica lookup, so
	// whether a2 holds x is unknown and nothing may be inserted.
	a0 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	a1 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	a2 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	c := newTestCorrector(t, 3, []string{"a0", "a1", "a2"},
		map[string]*fakeAgent{"a0": a0, "a1": a1, "a2": a2}, []string{"a0", "a1", "a2"})
	fc, ok := c.discoverer.GetClient().(*fakeGRPCClient)
	if !ok {
		t.Fatalf("unexpected client type %T", c.discoverer.GetClient())
	}
	fc.unreachable = map[string]bool{"a2": true}
	err := startCorrector(t, c)
	if err == nil || !strings.Contains(err.Error(), "incomplete replica information") {
		t.Fatalf("Start() error = %v, want an incomplete replica information error", err)
	}
	for addr, a := range map[string]*fakeAgent{"a0": a0, "a1": a1, "a2": a2} {
		a.mu.Lock()
		n := len(a.inserts) + len(a.updates) + len(a.updateTSs)
		a.mu.Unlock()
		if n != 0 {
			t.Errorf("agent %s was mutated (%d calls) although a2 was never reached", addr, n)
		}
	}
}

// TestStart_UnreachedAgentFailsJob pins that an agent OrderedRange never runs the
// callback for (no connection pool) fails the job: its ids were never checked.
func TestStart_UnreachedAgentFailsJob(t *testing.T) {
	t.Parallel()
	a0 := &fakeAgent{objects: map[string]*payload.Object_Vector{}}
	a1 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	c := newTestCorrector(t, 1, []string{"a0", "a1"},
		map[string]*fakeAgent{"a0": a0, "a1": a1}, []string{"a0", "a1"})
	fc, ok := c.discoverer.GetClient().(*fakeGRPCClient)
	if !ok {
		t.Fatalf("unexpected client type %T", c.discoverer.GetClient())
	}
	// IndexDetail still reports a1, but the client has no pool for it.
	delete(fc.conns, "a1")
	err := startCorrector(t, c)
	if err == nil || !errors.Is(err, errors.ErrGRPCClientConnNotFound("a1")) {
		t.Fatalf("Start() error = %v, want a connection not found error for a1", err)
	}
	if errors.Is(err, errors.ErrGRPCClientConnNotFound("a0")) {
		t.Errorf("Start() error = %v, reports the processed agent a0 as unprocessed", err)
	}
}

// TestStart_RetriedReplicaInfoIsComplete pins that a GetTimestamp failure that the
// client backoff retries successfully does not mark the replica info incomplete.
func TestStart_RetriedReplicaInfoIsComplete(t *testing.T) {
	t.Parallel()
	// x is on a0 and a1 (a1 fails its first GetTimestamp only) but not on a2, so
	// with index_replica 3 one copy is missing and must be inserted.
	a0 := &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}}
	a1 := &fakeAgent{
		objects:   map[string]*payload.Object_Vector{"x": newVec(100, 1)},
		tsErr:     status.Error(codes.Unavailable, "temporarily unavailable"),
		tsErrOnce: true,
	}
	a2 := &fakeAgent{objects: map[string]*payload.Object_Vector{}}
	c := newTestCorrector(t, 3, []string{"a0", "a1", "a2"},
		map[string]*fakeAgent{"a0": a0, "a1": a1, "a2": a2}, []string{"a0", "a1", "a2"})
	if err := startCorrector(t, c); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	total := 0
	for _, a := range []*fakeAgent{a0, a1, a2} {
		a.mu.Lock()
		total += len(a.inserts)
		a.mu.Unlock()
	}
	if total != 1 {
		t.Errorf("inserts = %d, want 1: a retried GetTimestamp must not leave the replica info incomplete", total)
	}
}

// TestStart_ReportsReplicaCorrectionFailures pins that an agent failure while
// fixing a replica shortage or oversupply fails the job, while a shortage left
// only because there are fewer agents than index_replica does not.
func TestStart_ReportsReplicaCorrectionFailures(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a0, a1      *fakeAgent
		name        string
		errSubstr   string
		replica     int
		wantErr     bool
		serial      bool
		wantInserts int
	}{
		{
			name:        "failed insert of a missing replica is an error",
			replica:     2,
			a0:          &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}},
			a1:          &fakeAgent{objects: map[string]*payload.Object_Vector{}, insertErr: status.Error(codes.Internal, "insert refused")},
			wantErr:     true,
			errSubstr:   "insert refused",
			wantInserts: 1,
		},
		{
			// a failed insert must not cancel the stream: with one worker the
			// ids are corrected one after another, so every id after the first
			// is only tried if the earlier failure did not stop the stream.
			name:    "failed insert does not stop correction of later ids",
			replica: 2,
			serial:  true,
			a0: &fakeAgent{objects: map[string]*payload.Object_Vector{
				"a": {Id: "a", Vector: []float32{1}, Timestamp: 100},
				"b": {Id: "b", Vector: []float32{2}, Timestamp: 100},
				"c": {Id: "c", Vector: []float32{3}, Timestamp: 100},
				"d": {Id: "d", Vector: []float32{4}, Timestamp: 100},
			}},
			a1:          &fakeAgent{objects: map[string]*payload.Object_Vector{}, insertErr: status.Error(codes.Internal, "insert refused")},
			wantErr:     true,
			errSubstr:   "insert refused",
			wantInserts: 4,
		},
		{
			name:        "shortage caused only by too few agents is not an error",
			replica:     3,
			a0:          &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}},
			a1:          &fakeAgent{objects: map[string]*payload.Object_Vector{}},
			wantInserts: 1,
		},
		{
			name:      "failed delete of an oversupplied replica is an error",
			replica:   1,
			a0:        &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}, removeErr: status.Error(codes.Internal, "remove refused")},
			a1:        &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}, removeErr: status.Error(codes.Internal, "remove refused")},
			wantErr:   true,
			errSubstr: "remove refused",
		},
		{
			name:    "oversupply fixed by a successful delete is not an error",
			replica: 1,
			a0:      &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}},
			a1:      &fakeAgent{objects: map[string]*payload.Object_Vector{"x": newVec(100, 1)}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			c := newTestCorrector(tt, tc.replica, []string{"a0", "a1"},
				map[string]*fakeAgent{"a0": tc.a0, "a1": tc.a1}, []string{"a0", "a1"})
			if tc.serial {
				c.streamListConcurrency = 1
			}
			err := startCorrector(tt, c)
			if !tc.wantErr {
				if err != nil {
					tt.Errorf("Start() error = %v, want nil", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.errSubstr) {
				tt.Errorf("Start() error = %v, want it to contain %q", err, tc.errSubstr)
			}
			tc.a1.mu.Lock()
			n := len(tc.a1.inserts)
			tc.a1.mu.Unlock()
			if n != tc.wantInserts {
				tt.Errorf("inserts on a1 = %d, want %d", n, tc.wantInserts)
			}
		})
	}
}

func TestStart_BoundsReportedCorrectionFailures(t *testing.T) {
	t.Parallel()
	const extra = 3
	objs := make(map[string]*payload.Object_Vector, maxReportedFailures+extra)
	for i := range maxReportedFailures + extra {
		id := "id" + strconv.Itoa(i)
		objs[id] = &payload.Object_Vector{Id: id, Vector: []float32{float32(i)}, Timestamp: 100}
	}
	a0 := &fakeAgent{objects: objs}
	a1 := &fakeAgent{objects: map[string]*payload.Object_Vector{}, insertErr: status.Error(codes.Internal, "insert refused")}
	c := newTestCorrector(t, 2, []string{"a0", "a1"},
		map[string]*fakeAgent{"a0": a0, "a1": a1}, []string{"a0", "a1"})
	err := startCorrector(t, c)
	if err == nil {
		t.Fatal("Start() error = nil, want an error")
	}
	if got := strings.Count(err.Error(), "failed to insert object"); got != maxReportedFailures {
		t.Errorf("reported failures = %d, want %d", got, maxReportedFailures)
	}
	if want := strconv.Itoa(extra) + " more index ids could not be corrected"; !strings.Contains(err.Error(), want) {
		t.Errorf("Start() error = %v, want it to contain %q", err, want)
	}
}
