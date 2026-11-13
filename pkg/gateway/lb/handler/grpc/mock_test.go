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
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Fixture values shared by the handler tests of this package.
const (
	fixtureUUID = "uuid"
	loopbackIP  = "127.0.0.1"
)

// regGateway is a minimal service.Gateway stub. Unset hooks fall through to the
// embedded nil interface, so a test touching an unexpected method fails loudly.
type regGateway struct {
	service.Gateway
	broadCast func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error
	doMulti   func(ctx context.Context, num int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error
	addrs     []string
	agents    int
	// addrsCalls counts Addrs calls: the handlers query it only while building error details.
	addrsCalls atomic.Int64
	// agentCalls counts GetAgentCount calls.
	agentCalls atomic.Int64
}

func (g *regGateway) Addrs(context.Context) []string {
	g.addrsCalls.Add(1)
	return g.addrs
}

func (g *regGateway) GetAgentCount(context.Context) int {
	g.agentCalls.Add(1)
	return g.agents
}

func (g *regGateway) BroadCast(
	ctx context.Context,
	_ service.BroadCastKind,
	f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error,
) error {
	return g.broadCast(ctx, f)
}

func (g *regGateway) DoMulti(
	ctx context.Context,
	num int,
	f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error,
) error {
	return g.doMulti(ctx, num, f)
}

// regClient is a minimal vald.Client stub.
type regClient struct {
	vald.Client
	update    func(ctx context.Context, in *payload.Update_Request) (*payload.Object_Location, error)
	search    func(ctx context.Context, in *payload.Search_Request) (*payload.Search_Response, error)
	insert    func(ctx context.Context, in *payload.Insert_Request) (*payload.Object_Location, error)
	remove    func(ctx context.Context, in *payload.Remove_Request) (*payload.Object_Location, error)
	exists    func(ctx context.Context, in *payload.Object_ID) (*payload.Object_ID, error)
	getObject func(ctx context.Context, in *payload.Object_VectorRequest) (*payload.Object_Vector, error)
	getTS     func(ctx context.Context, in *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error)
}

func (c *regClient) Insert(
	ctx context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.insert(ctx, in)
}

func (c *regClient) Remove(
	ctx context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.remove(ctx, in)
}

func (c *regClient) Exists(
	ctx context.Context, in *payload.Object_ID, _ ...grpc.CallOption,
) (*payload.Object_ID, error) {
	return c.exists(ctx, in)
}

func (c *regClient) GetObject(
	ctx context.Context, in *payload.Object_VectorRequest, _ ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return c.getObject(ctx, in)
}

func (c *regClient) GetTimestamp(
	ctx context.Context, in *payload.Object_TimestampRequest, _ ...grpc.CallOption,
) (*payload.Object_Timestamp, error) {
	return c.getTS(ctx, in)
}

func (c *regClient) Update(
	ctx context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return c.update(ctx, in)
}

func (c *regClient) Search(
	ctx context.Context, in *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return c.search(ctx, in)
}

// okClient answers every write/read the Insert/Remove/Update handlers issue, recording the
// requests it was forwarded.
type okClient struct {
	*regClient
	insertIn *payload.Insert_Request
	removeIn *payload.Remove_Request
}

func newOKClient(vec []float32) *okClient {
	const uuid = "uuid"
	c := new(okClient)
	loc := &payload.Object_Location{Uuid: uuid, Name: "agent-0", Ips: []string{loopbackIP}}
	c.regClient = &regClient{
		insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
			c.insertIn = in
			return loc, nil
		},
		remove: func(_ context.Context, in *payload.Remove_Request) (*payload.Object_Location, error) {
			c.removeIn = in
			return loc, nil
		},
		exists: func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{Id: uuid}, nil
		},
		getObject: func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
			return &payload.Object_Vector{Id: uuid, Vector: vec}, nil
		},
	}
	return c
}

// syncGateway runs every callback once, inline, against c.
func syncGateway(c vald.Client) *regGateway {
	return &regGateway{
		agents: 1,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", c)
		},
		doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", c)
		},
	}
}

// pwAgent answers every write with err (nil: success located on its name) and
// counts the inserts it was asked for.
type pwAgent struct {
	*regClient
	removeErr error
	insertErr error
	updateErr error
	name      string
	inserts   atomic.Int64
}

func newPWAgent(name string, removeErr, insertErr, updateErr error) *pwAgent {
	a := &pwAgent{name: name, removeErr: removeErr, insertErr: insertErr, updateErr: updateErr}
	loc := func(uuid string) *payload.Object_Location {
		return &payload.Object_Location{Uuid: uuid, Name: name, Ips: []string{name}}
	}
	a.regClient = &regClient{
		remove: func(_ context.Context, in *payload.Remove_Request) (*payload.Object_Location, error) {
			if a.removeErr != nil {
				return nil, a.removeErr
			}
			return loc(in.GetId().GetId()), nil
		},
		insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
			a.inserts.Add(1)
			if a.insertErr != nil {
				return nil, a.insertErr
			}
			return loc(in.GetVector().GetId()), nil
		},
		update: func(_ context.Context, in *payload.Update_Request) (*payload.Object_Location, error) {
			if a.updateErr != nil {
				return nil, a.updateErr
			}
			return loc(in.GetVector().GetId()), nil
		},
	}
	return a
}

func (a *pwAgent) UpdateTimestamp(
	_ context.Context, in *payload.Update_TimestampRequest, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	if a.updateErr != nil {
		return nil, a.updateErr
	}
	return &payload.Object_Location{Uuid: in.GetId(), Name: a.name, Ips: []string{a.name}}, nil
}

// pwGateway mimics the lb service gateway: BroadCast calls every agent and joins
// their errors in target order, DoMulti calls the agents in order until num of
// them succeeded (a failed agent is not counted) and only reports a missing
// connection, as OrderedRange logs the per-agent failures.
func pwGateway(agents ...*pwAgent) *regGateway {
	return &regGateway{
		agents: len(agents),
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			var errs []error
			for _, a := range agents {
				if err := f(ctx, a.name, a); err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		},
		doMulti: func(ctx context.Context, num int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			var cur int
			for _, a := range agents {
				if cur >= num {
					break
				}
				if err := f(ctx, a.name, a); err == nil {
					cur++
				}
			}
			return nil
		},
	}
}

func pwCode(t *testing.T, err error) codes.Code {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok || st == nil {
		t.Fatalf("not a status error: %v", err)
	}
	return st.Code()
}

// tsAgent mimics the NGT agent's timestamp handling for a single uuid: Insert takes
// the vector's timestamp and falls back to the config's only when it is 0
// (pkg/agent/core/ngt/handler/grpc/insert.go), and an insert older than a queued
// delete is silently dropped while still reporting success (vqueue.PushInsert).
type tsAgent struct {
	*regClient
	updateTS func(ctx context.Context, in *payload.Update_TimestampRequest) (*payload.Object_Location, error)

	mu       sync.Mutex
	delTS    int64
	stored   bool
	storedTS int64
	inserts  int
}

func newTSAgent(name, ip string, stored []float32, storedTS int64) *tsAgent {
	a := &tsAgent{stored: stored != nil, storedTS: storedTS}
	loc := &payload.Object_Location{Uuid: fixtureUUID, Name: name, Ips: []string{ip}}
	a.regClient = &regClient{
		insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
			ts := in.GetVector().GetTimestamp()
			if ts == 0 {
				ts = in.GetConfig().GetTimestamp()
			}
			a.mu.Lock()
			defer a.mu.Unlock()
			a.inserts++
			if a.delTS > ts {
				return loc, nil // dropped behind the newer delete
			}
			a.stored, a.storedTS = true, ts
			return loc, nil
		},
		remove: func(_ context.Context, in *payload.Remove_Request) (*payload.Object_Location, error) {
			a.mu.Lock()
			defer a.mu.Unlock()
			a.delTS = in.GetConfig().GetTimestamp()
			a.stored = false
			return loc, nil
		},
		exists: func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{Id: fixtureUUID}, nil
		},
		getObject: func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
			return &payload.Object_Vector{Id: fixtureUUID, Vector: stored, Timestamp: storedTS}, nil
		},
		// Like the agent's Update: delete at the config timestamp, insert at +1.
		update: func(_ context.Context, in *payload.Update_Request) (*payload.Object_Location, error) {
			if ts := in.GetConfig().GetTimestamp(); ts > 0 {
				a.mu.Lock()
				defer a.mu.Unlock()
				a.delTS = ts
				a.stored, a.storedTS = true, ts+1
			}
			return loc, nil
		},
	}
	a.updateTS = func(context.Context, *payload.Update_TimestampRequest) (*payload.Object_Location, error) {
		return loc, nil
	}
	return a
}

func (a *tsAgent) UpdateTimestamp(
	ctx context.Context, in *payload.Update_TimestampRequest, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return a.updateTS(ctx, in)
}

func (a *tsAgent) state() (stored bool, ts, delTS int64, inserts int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.stored, a.storedTS, a.delTS, a.inserts
}

// twoAgentGateway broadcasts to `first` only and runs DoMulti against `second`, so the
// handler sees a replica shortage it tops up on `second`.
func twoAgentGateway(first, second vald.Client) *regGateway {
	return &regGateway{
		agents: 2,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", first)
		},
		doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.2:8081", second)
		},
	}
}

func probeResponse() *payload.Search_Response {
	return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "found", Distance: 0.1}}}
}

// collectSearchMetrics returns the cumulative value per instrument name and
// attribute set currently reported by sm.
func collectSearchMetrics(
	t *testing.T, sm *lbmetrics.SearchMetrics,
) map[string]map[attribute.Distinct]int64 {
	t.Helper()
	views, err := sm.View()
	if err != nil {
		t.Fatalf("View() error = %v", err)
	}
	reader := sdkmetric.NewManualReader()
	opts := make([]sdkmetric.Option, 0, len(views)+1)
	opts = append(opts, sdkmetric.WithReader(reader))
	for _, vw := range views {
		opts = append(opts, sdkmetric.WithView(vw))
	}
	mp := sdkmetric.NewMeterProvider(opts...)
	t.Cleanup(func() {
		if err := mp.Shutdown(context.WithoutCancel(t.Context())); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	})
	if err := sm.Register(mp.Meter("test")); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(t.Context(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	got := make(map[string]map[attribute.Distinct]int64)
	for _, scm := range rm.ScopeMetrics {
		for _, m := range scm.Metrics {
			s, ok := m.Data.(metricdata.Sum[int64])
			if !ok {
				t.Errorf("%s: data = %T, want metricdata.Sum[int64]", m.Name, m.Data)
				continue
			}
			pts := make(map[attribute.Distinct]int64, len(s.DataPoints))
			for _, dp := range s.DataPoints {
				pts[dp.Attributes.Equivalent()] = dp.Value
			}
			got[m.Name] = pts
		}
	}
	return got
}
