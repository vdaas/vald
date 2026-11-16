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

// Package service
package service

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// Test_index_createIndex_queuesOnError pins that an agent whose CreateIndex fails with anything
// but FailedPrecondition is still queued for the next save tick, because the agent may have
// built a new index generation before returning the error.
func Test_index_createIndex_queuesOnError(t *testing.T) {
	t.Parallel()
	const addr = "agent-0"
	tests := []struct {
		createErr  error
		name       string
		wantQueued bool
		wantErr    bool
	}{
		{
			name:       "successful CreateIndex queues the agent",
			wantQueued: true,
		},
		{
			name:       "Internal error (e.g. statistics or export failure after build) queues the agent",
			createErr:  status.Error(codes.Internal, "failed to export index info"),
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:       "Canceled after the delete phase queues the agent",
			createErr:  status.Error(codes.Canceled, "canceled"),
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:       "Unavailable (lost reply) queues the agent",
			createErr:  status.Error(codes.Unavailable, "connection lost"),
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:      "FailedPrecondition (nothing to index) does not queue the agent",
			createErr: status.Error(codes.FailedPrecondition, "uncommitted index is empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := &fakeIndexAgent{createErr: tt.createErr}
			a.uncommitted.Store(100)
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				createIndexConcurrency: 1,
			}

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			if err := idx.createIndex(ctx, false); (err != nil) != tt.wantErr {
				t.Errorf("createIndex error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := a.creates.Load(); got != 1 {
				t.Errorf("CreateIndex calls = %d, want 1", got)
			}
			if _, got := idx.shouldSaveList.Load(addr); got != tt.wantQueued {
				t.Errorf("queued after createIndex = %v, want %v", got, tt.wantQueued)
			}
		})
	}
}

// Test_index_saveOnTick_forcedSaveRetriedWhenBusy pins that a forced-save tick that finds a
// save already running is rescheduled after saveIndexDuration instead of being dropped until
// the next saveIndexDurationLimit tick, and that a regular tick is still simply dropped.
func Test_index_saveOnTick_forcedSaveRetriedWhenBusy(t *testing.T) {
	t.Parallel()
	const (
		addr     = "agent-0"
		retry    = 10 * time.Millisecond
		longTick = time.Hour
		// waitFire bounds the wait for a rescheduled tick; quietFor is how long a ticker that
		// was not rescheduled is observed (it cannot fire before longTick).
		waitFire = 10 * time.Second
		quietFor = 100 * time.Millisecond
	)
	tests := []struct {
		name      string
		busy      bool
		force     bool
		wantFire  bool
		wantSaves uint32
	}{
		{
			name:     "busy forced tick is rescheduled after saveIndexDuration",
			busy:     true,
			force:    true,
			wantFire: true,
		},
		{
			name:  "busy regular tick is dropped",
			busy:  true,
			force: false,
		},
		{
			name:      "idle forced tick saves and resets to saveIndexDurationLimit",
			force:     true,
			wantSaves: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := new(fakeIndexAgent)
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				saveIndexConcurrency:   1,
				saveIndexDuration:      retry,
				saveIndexDurationLimit: longTick,
			}
			if !tt.force {
				idx.saveIndexDuration = longTick
			}
			tick := time.NewTicker(longTick)
			defer tick.Stop()
			var mu sync.Mutex
			if tt.busy {
				mu.Lock()
				defer mu.Unlock()
			}
			ech := make(chan error, 1)

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			idx.saveOnTick(ctx, &mu, tick, tt.force, ech)

			if got := a.saves.Load(); got != tt.wantSaves {
				t.Errorf("SaveIndex calls = %d, want %d", got, tt.wantSaves)
			}
			select {
			case err := <-ech:
				t.Errorf("unexpected error reported: %v", err)
			default:
			}
			wait := quietFor
			if tt.wantFire {
				wait = waitFire
			}
			timer := time.NewTimer(wait)
			defer timer.Stop()
			select {
			case <-tick.C:
				if !tt.wantFire {
					t.Error("ticker fired, want it to keep its long period")
				}
			case <-timer.C:
				if tt.wantFire {
					t.Errorf("ticker did not fire within %v, want the forced save rescheduled after %v", waitFire, retry)
				}
			}
		})
	}
}

// fakeIndexAgent serves IndexInfo with a live uncommitted count and counts CreateIndex calls.
type fakeIndexAgent struct {
	core.UnimplementedAgentServer
	vald.UnimplementedIndexServer

	saveErr     error // returned by SaveIndex
	createErr   error // returned by CreateIndex
	uncommitted atomic.Uint32
	creates     atomic.Uint32
	saves       atomic.Uint32
}

func (a *fakeIndexAgent) IndexInfo(
	context.Context, *payload.Empty,
) (*payload.Info_Index_Count, error) {
	return &payload.Info_Index_Count{Uncommitted: a.uncommitted.Load()}, nil
}

func (a *fakeIndexAgent) CreateIndex(
	context.Context, *payload.Control_CreateIndexRequest,
) (*payload.Empty, error) {
	a.creates.Add(1)
	if a.createErr != nil {
		return nil, a.createErr
	}
	a.uncommitted.Store(0)
	return new(payload.Empty), nil
}

// fakeIndexGRPCClient ranges over fixed conns sequentially.
type fakeIndexGRPCClient struct {
	grpc.Client
	conns map[string]*grpc.ClientConn
	order []string
}

func (f *fakeIndexGRPCClient) RangeConcurrent(
	ctx context.Context,
	_ int,
	fn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	return f.OrderedRangeConcurrent(ctx, f.order, 1, fn)
}

func (f *fakeIndexGRPCClient) OrderedRangeConcurrent(
	ctx context.Context,
	order []string,
	_ int,
	fn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	for _, addr := range order {
		if conn, ok := f.conns[addr]; ok {
			if err := fn(ctx, addr, conn); err != nil {
				return err
			}
		}
	}
	return nil
}

type fakeIndexDiscoverer struct {
	discoverer.Client
	client *fakeIndexGRPCClient
}

func (f *fakeIndexDiscoverer) GetClient() grpc.Client            { return f.client }
func (f *fakeIndexDiscoverer) GetAddrs(context.Context) []string { return f.client.order }

//nolint:unparam // addr is the passthrough target name and log label of the served agent; every current caller serves a single agent
func serveFakeIndexAgent(t *testing.T, addr string, a *fakeIndexAgent) *grpc.ClientConn {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := ggrpc.NewServer()
	core.RegisterAgentServer(srv, a)
	vald.RegisterIndexServer(srv, a)
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
	return conn
}

// Test_index_createIndex_refreshesInfosBeforeSkip pins that the per-agent skip decision uses
// the agent's current uncommitted count, not the snapshot left by the previous run.
func Test_index_createIndex_refreshesInfosBeforeSkip(t *testing.T) {
	t.Parallel()
	const addr = "agent-0"
	tests := []struct {
		name            string
		stale           uint32 // snapshot left by the previous run
		live            uint32 // what the agent reports now
		lowIndexSkip    bool
		wantCreateCalls uint32
	}{
		{
			name:            "stale zero but live backlog creates the index",
			stale:           0,
			live:            100,
			lowIndexSkip:    true,
			wantCreateCalls: 1,
		},
		{
			name:            "stale low but live backlog over minUncommitted creates the index",
			stale:           1,
			live:            100,
			lowIndexSkip:    true,
			wantCreateCalls: 1,
		},
		{
			name:            "stale backlog but live zero skips the agent",
			stale:           100,
			live:            0,
			lowIndexSkip:    false,
			wantCreateCalls: 0,
		},
		{
			name:            "stale backlog but live low skips the agent on a low-index-skip run",
			stale:           100,
			live:            1,
			lowIndexSkip:    true,
			wantCreateCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := new(fakeIndexAgent)
			a.uncommitted.Store(tt.live)
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				createIndexConcurrency: 1,
				minUncommitted:         10,
			}
			idx.indexInfos.Store(addr, &payload.Info_Index_Count{Uncommitted: tt.stale})

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			if err := idx.createIndex(ctx, tt.lowIndexSkip); err != nil {
				t.Fatalf("createIndex error = %v", err)
			}
			if got := a.creates.Load(); got != tt.wantCreateCalls {
				t.Errorf("CreateIndex calls = %d, want %d", got, tt.wantCreateCalls)
			}
			// The post-run refresh still records the state after indexing.
			info, ok := idx.indexInfos.Load(addr)
			if !ok {
				t.Fatal("indexInfos entry is missing after createIndex")
			}
			if got, want := info.GetUncommitted(), a.uncommitted.Load(); got != want {
				t.Errorf("indexInfos uncommitted = %d, want %d", got, want)
			}
		})
	}
}

func (a *fakeIndexAgent) SaveIndex(context.Context, *payload.Empty) (*payload.Empty, error) {
	a.saves.Add(1)
	if a.saveErr != nil {
		return nil, a.saveErr
	}
	return new(payload.Empty), nil
}

// Test_index_saveIndex_requeueOnError pins that an agent whose CreateIndex
// result is queued for saving stays queued when SaveIndex fails, so the next
// non-forced save tick retries it instead of waiting for the next CreateIndex
// or the force-save limit.
func Test_index_saveIndex_requeueOnError(t *testing.T) {
	t.Parallel()
	const addr = "agent-0"
	tests := []struct {
		saveErr    error
		name       string
		wantSaves  uint32
		queued     bool
		force      bool
		wantQueued bool
		wantErr    bool
	}{
		{
			name:       "failed save of a queued agent keeps it queued",
			queued:     true,
			saveErr:    status.Error(codes.Internal, "save failed"),
			wantSaves:  1,
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:       "failed forced save of a queued agent keeps it queued",
			queued:     true,
			force:      true,
			saveErr:    status.Error(codes.Unavailable, "agent unavailable"),
			wantSaves:  1,
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:      "successful save dequeues the agent",
			queued:    true,
			wantSaves: 1,
		},
		{
			name:      "skipped save (FailedPrecondition) dequeues the agent",
			queued:    true,
			saveErr:   status.Error(codes.FailedPrecondition, "nothing to save"),
			wantSaves: 1,
		},
		{
			name:      "failed forced save of an agent that was not queued does not queue it",
			force:     true,
			saveErr:   status.Error(codes.Internal, "save failed"),
			wantSaves: 1,
			wantErr:   true,
		},
		{
			name: "not queued and not forced does not call SaveIndex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := &fakeIndexAgent{saveErr: tt.saveErr}
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				saveIndexConcurrency: 1,
			}
			if tt.queued {
				idx.shouldSaveList.Store(addr, true)
			}

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			if err := idx.saveIndex(ctx, tt.force); (err != nil) != tt.wantErr {
				t.Errorf("saveIndex error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := a.saves.Load(); got != tt.wantSaves {
				t.Errorf("SaveIndex calls = %d, want %d", got, tt.wantSaves)
			}
			if _, got := idx.shouldSaveList.Load(addr); got != tt.wantQueued {
				t.Errorf("queued after saveIndex = %v, want %v", got, tt.wantQueued)
			}
		})
	}
}

// saveIndexConcurrency is passed to OrderedRangeConcurrent, where 0 means an
// unbounded fan-out of SaveIndex to every agent; it must default to a bound.
func TestNew_saveIndexConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts []Option
		want int
	}{
		{name: "default is serial", want: 1},
		{name: "option is applied", opts: []Option{WithSaveConcurrency(4)}, want: 4},
		{name: "zero keeps default", opts: []Option{WithSaveConcurrency(0)}, want: 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			got, err := New(tc.opts...)
			if err != nil {
				tt.Fatalf("New() error = %v", err)
			}
			idx, ok := got.(*index)
			if !ok {
				tt.Fatalf("New() returned %T, want *index", got)
			}
			if idx.saveIndexConcurrency != tc.want {
				tt.Errorf("saveIndexConcurrency = %d, want %d", idx.saveIndexConcurrency, tc.want)
			}
		})
	}
}

// saveIndexDuration drives the regular save ticker in Start; when it stays 0 Start
// turns it into math.MaxInt64 and only the forced saveIndexDurationLimit pass runs.
func TestNew_saveIndexDuration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		opts    []Option
		want    time.Duration
		wantErr bool
	}{
		{name: "default is unset", want: 0},
		{name: "option is applied", opts: []Option{WithSaveIndexDuration("10m")}, want: 10 * time.Minute},
		{name: "empty keeps unset", opts: []Option{WithSaveIndexDuration("")}, want: 0},
		{name: "invalid duration fails", opts: []Option{WithSaveIndexDuration("invalid")}, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			got, err := New(tc.opts...)
			if tc.wantErr {
				if err == nil {
					tt.Fatal("New() error = nil, want error")
				}
				return
			}
			if err != nil {
				tt.Fatalf("New() error = %v", err)
			}
			idx, ok := got.(*index)
			if !ok {
				tt.Fatalf("New() returned %T, want *index", got)
			}
			if idx.saveIndexDuration != tc.want {
				tt.Errorf("saveIndexDuration = %v, want %v", idx.saveIndexDuration, tc.want)
			}
			if idx.saveIndexDurationLimit != 3*time.Hour {
				tt.Errorf("saveIndexDurationLimit = %v, want default 3h", idx.saveIndexDurationLimit)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		wantIdx Indexer
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Indexer, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotIdx Indexer, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotIdx, w.wantIdx) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIdx, w.wantIdx)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotIdx, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, gotIdx, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got, err := idx.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_createIndex(t *testing.T) {
// 	type args struct {
// 		ctx                context.Context
// 		enableLowIndexSkip bool
// 	}
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           enableLowIndexSkip:false,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           enableLowIndexSkip:false,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			err := idx.createIndex(test.args.ctx, test.args.enableLowIndexSkip)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_saveIndex(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		force bool
// 	}
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           force:false,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           force:false,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			err := idx.saveIndex(test.args.ctx, test.args.force)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_loadInfos(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			err := idx.loadInfos(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_IsIndexing(t *testing.T) {
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got := idx.IsIndexing()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_IsSaving(t *testing.T) {
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got := idx.IsSaving()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_NumberOfUUIDs(t *testing.T) {
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want uint32
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint32) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got := idx.NumberOfUUIDs()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_NumberOfUncommittedUUIDs(t *testing.T) {
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want uint32
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint32) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got := idx.NumberOfUncommittedUUIDs()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_LoadIndexDetail(t *testing.T) {
// 	type fields struct {
// 		eg                     errgroup.Group
// 		client                 discoverer.Client
// 		shouldSaveList         sync.Map[string, bool]
// 		indexInfos             sync.Map[string, *payload.Info_Index_Count]
// 		createIndexConcurrency int
// 		saveIndexDuration      time.Duration
// 		saveIndexDurationLimit time.Duration
// 		indexDurationLimit     time.Duration
// 		indexDuration          time.Duration
// 		saveIndexConcurrency   int
// 		creationPoolSize       uint32
// 		indexing               atomic.Bool
// 		saving                 atomic.Bool
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		wantDetail *payload.Info_Index_Detail
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Detail) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotDetail *payload.Info_Index_Detail) error {
// 		if !reflect.DeepEqual(gotDetail, w.wantDetail) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDetail, w.wantDetail)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           client:nil,
// 		           shouldSaveList:nil,
// 		           indexInfos:nil,
// 		           createIndexConcurrency:0,
// 		           saveIndexDuration:nil,
// 		           saveIndexDurationLimit:nil,
// 		           indexDurationLimit:nil,
// 		           indexDuration:nil,
// 		           saveIndexConcurrency:0,
// 		           creationPoolSize:0,
// 		           indexing:nil,
// 		           saving:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				eg:                     test.fields.eg,
// 				client:                 test.fields.client,
// 				shouldSaveList:         test.fields.shouldSaveList,
// 				indexInfos:             test.fields.indexInfos,
// 				createIndexConcurrency: test.fields.createIndexConcurrency,
// 				saveIndexDuration:      test.fields.saveIndexDuration,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				indexDuration:          test.fields.indexDuration,
// 				saveIndexConcurrency:   test.fields.saveIndexConcurrency,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexing:               test.fields.indexing,
// 				saving:                 test.fields.saving,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			gotDetail := idx.LoadIndexDetail()
// 			if err := checkFunc(test.want, gotDetail); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
