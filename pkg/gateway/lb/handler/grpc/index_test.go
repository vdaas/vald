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
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// mergeInfoIndexStatistics must not depend on map iteration order: a minimum
// of 0 is a real value (not "unset"), and IndegreeCount is sized from the
// longest input. Every assignment of the inputs to keys is tried repeatedly so
// that Go's randomized map order visits the agents in many orders.
func Test_mergeInfoIndexStatistics_orderIndependent(t *testing.T) {
	t.Parallel()
	// fresh inputs per call: sumHistograms may reuse an input histogram slice.
	inputs := func() []*payload.Info_Index_Statistics {
		return []*payload.Info_Index_Statistics{
			{Valid: true, MinNumberOfIndegree: 0, MinNumberOfOutdegree: 3, IndegreeCount: []int64{2, 4}},
			{Valid: true, MinNumberOfIndegree: 5, MinNumberOfOutdegree: 1, IndegreeCount: []int64{4, 6, 9}},
			{Valid: true, MinNumberOfIndegree: 2, MinNumberOfOutdegree: 0, IndegreeCount: []int64{6}},
		}
	}
	const (
		wantMinIn  = uint64(0)
		wantMinOut = uint64(0)
	)
	// index 0: (2+4+6)/3, index 1: (4+6)/2, index 2: 9/1
	wantCount := []int64{4, 5, 9}

	perms := [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	keys := []string{"a", "b", "c"}
	for _, perm := range perms {
		for range 50 {
			in := inputs()
			stats := make(map[string]*payload.Info_Index_Statistics, len(perm))
			for i, p := range perm {
				stats[keys[i]] = in[p]
			}
			got := mergeInfoIndexStatistics(stats)
			if got.GetMinNumberOfIndegree() != wantMinIn {
				t.Fatalf("perm %v: MinNumberOfIndegree = %d, want %d", perm, got.GetMinNumberOfIndegree(), wantMinIn)
			}
			if got.GetMinNumberOfOutdegree() != wantMinOut {
				t.Fatalf("perm %v: MinNumberOfOutdegree = %d, want %d", perm, got.GetMinNumberOfOutdegree(), wantMinOut)
			}
			if !slices.Equal(got.GetIndegreeCount(), wantCount) {
				t.Fatalf("perm %v: IndegreeCount = %v, want %v", perm, got.GetIndegreeCount(), wantCount)
			}
		}
	}
}

// The averaged fields used to be divided by every entry, invalid ones included,
// while the medians are taken over the valid entries only.
func Test_mergeInfoIndexStatistics_averagesOverValidEntries(t *testing.T) {
	t.Parallel()
	valid := func(v float64) *payload.Info_Index_Statistics {
		return &payload.Info_Index_Statistics{
			Valid:                          true,
			ModeIndegree:                   uint64(v),
			ModeOutdegree:                  uint64(v),
			VarianceOfIndegree:             v,
			VarianceOfOutdegree:            v,
			MeanEdgeLength:                 v,
			MeanEdgeLengthFor10Edges:       v,
			MeanIndegreeDistanceFor10Edges: v,
			MeanNumberOfEdgesPerNode:       v,
			C1Indegree:                     v,
			C5Indegree:                     v,
			C95Outdegree:                   v,
			C99Outdegree:                   v,
		}
	}
	tests := []struct {
		stats map[string]*payload.Info_Index_Statistics
		name  string
		want  float64
	}{
		{
			name:  "one valid one invalid",
			stats: map[string]*payload.Info_Index_Statistics{"a": valid(8), "b": {Valid: false}},
			want:  8,
		},
		{
			name: "two valid two invalid",
			stats: map[string]*payload.Info_Index_Statistics{
				"a": valid(4), "b": valid(8), "c": {Valid: false}, "d": {Valid: false},
			},
			want: 6,
		},
		{
			name:  "all valid",
			stats: map[string]*payload.Info_Index_Statistics{"a": valid(4), "b": valid(8)},
			want:  6,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := mergeInfoIndexStatistics(test.stats)
			for name, v := range map[string]float64{
				"ModeIndegree":                   float64(got.GetModeIndegree()),
				"ModeOutdegree":                  float64(got.GetModeOutdegree()),
				"VarianceOfIndegree":             got.GetVarianceOfIndegree(),
				"VarianceOfOutdegree":            got.GetVarianceOfOutdegree(),
				"MeanEdgeLength":                 got.GetMeanEdgeLength(),
				"MeanEdgeLengthFor10Edges":       got.GetMeanEdgeLengthFor10Edges(),
				"MeanIndegreeDistanceFor10Edges": got.GetMeanIndegreeDistanceFor10Edges(),
				"MeanNumberOfEdgesPerNode":       got.GetMeanNumberOfEdgesPerNode(),
				"C1Indegree":                     got.GetC1Indegree(),
				"C5Indegree":                     got.GetC5Indegree(),
				"C95Outdegree":                   got.GetC95Outdegree(),
				"C99Outdegree":                   got.GetC99Outdegree(),
			} {
				if v != test.want {
					t.Errorf("%s = %v, want %v", name, v, test.want)
				}
			}
		})
	}
}

// The ResourceType detail of an index broadcast error names the "BroadCast" stage (was "BroadCase").
func Test_server_handleIndexBroadCastError_resourceType(t *testing.T) {
	t.Parallel()
	const target = "127.0.0.1:8081"
	s := &server{name: "gateway", ip: loopbackIP}

	err := s.handleIndexBroadCastError(nil, target, vald.IndexInfoRPCName, status.Error(codes.Internal, "boom"))
	if err == nil {
		t.Fatal("expected the Internal error to be propagated")
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		t.Fatalf("expected a gRPC status, got: %v", err)
	}
	for _, d := range st.Details() {
		ri, ok := d.(*errdetails.ResourceInfo)
		if !ok {
			continue
		}
		want := "." + "BroadCast/" + target
		if !strings.HasSuffix(ri.GetResourceType(), want) {
			t.Errorf("ResourceType = %q, want suffix %q", ri.GetResourceType(), want)
		}
		return
	}
	t.Errorf("no ResourceInfo detail in %v", st.Details())
}

// When every entry is invalid there is nothing to merge; this used to index into an empty slice.
func Test_mergeInfoIndexStatistics_noValidEntry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		stats     map[string]*payload.Info_Index_Statistics
		name      string
		wantValid bool
	}{
		{name: "nil", stats: nil},
		{name: "all invalid", stats: map[string]*payload.Info_Index_Statistics{
			"a": {Valid: false, NumberOfNodes: 3},
			"b": {Valid: false},
		}},
		{name: "one valid one invalid", stats: map[string]*payload.Info_Index_Statistics{
			"a": {Valid: true, NumberOfNodes: 3},
			"b": {Valid: false},
		}, wantValid: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := mergeInfoIndexStatistics(test.stats)
			if got == nil {
				t.Fatal("got nil statistics")
			}
			if got.GetValid() != test.wantValid {
				t.Errorf("Valid = %t, want %t", got.GetValid(), test.wantValid)
			}
			if !test.wantValid && got.GetNumberOfNodes() != 0 {
				t.Errorf("NumberOfNodes = %d, want 0 for nothing to merge", got.GetNumberOfNodes())
			}
		})
	}
}

// propClient answers IndexProperty with a canned result.
type propClient struct {
	*regClient
	err error
}

func (c propClient) IndexProperty(
	context.Context, *payload.Empty, ...grpc.CallOption,
) (*payload.Info_Index_PropertyDetail, error) {
	if c.err != nil {
		return nil, c.err
	}
	return &payload.Info_Index_PropertyDetail{
		Details: map[string]*payload.Info_Index_Property{"p": {Dimension: 3}},
	}, nil
}

// IndexProperty shares handleIndexBroadCastError with the other Index* handlers: a failing agent is
// propagated, a NotFound/canceled one is tolerated, a healthy one is merged under "<target>-<key>".
func Test_server_IndexProperty_broadCastErrorHandling(t *testing.T) {
	t.Parallel()
	const target = "127.0.0.1:8081"
	tests := []struct {
		agentErr error
		name     string
		wantCode codes.Code // codes.OK: no error expected
		wantKeys int
	}{
		{name: "healthy agent", wantKeys: 1},
		{name: "NotFound is tolerated", agentErr: status.Error(codes.NotFound, "nf")},
		{name: "Canceled is tolerated", agentErr: context.Canceled},
		{name: "Internal is propagated", agentErr: status.Error(codes.Internal, "boom"), wantCode: codes.Internal},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gw := &regGateway{
				agents: 1,
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					return f(ctx, target, propClient{regClient: new(regClient), err: test.agentErr})
				},
			}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: gw, eg: eg, name: "gateway", ip: loopbackIP}

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			detail, err := s.IndexProperty(ctx, new(payload.Empty))
			if test.wantCode != codes.OK {
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error: %v, want code %v", err, test.wantCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := len(detail.GetDetails()); got != test.wantKeys {
				t.Errorf("got %d property entries, want %d", got, test.wantKeys)
			}
			if test.wantKeys > 0 && detail.GetDetails()[target+"-p"] == nil {
				t.Errorf("missing %q in %v", target+"-p", detail.GetDetails())
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_IndexInfo(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		in1 *payload.Empty
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantVec *payload.Info_Index_Count
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Count, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Info_Index_Count, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
// 		           in1:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           in1:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotVec, err := s.IndexInfo(test.args.ctx, test.args.in1)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_IndexDetail(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		in1 *payload.Empty
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantVec *payload.Info_Index_Detail
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Detail, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Info_Index_Detail, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
// 		           in1:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           in1:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotVec, err := s.IndexDetail(test.args.ctx, test.args.in1)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_IndexStatistics(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Empty
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantVec *payload.Info_Index_Statistics
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Statistics, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Info_Index_Statistics, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotVec, err := s.IndexStatistics(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_IndexStatisticsDetail(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		in1 *payload.Empty
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantVec *payload.Info_Index_StatisticsDetail
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_StatisticsDetail, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Info_Index_StatisticsDetail, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
// 		           in1:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           in1:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotVec, err := s.IndexStatisticsDetail(test.args.ctx, test.args.in1)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_calculateMedian(t *testing.T) {
// 	type args struct {
// 		data []int32
// 	}
// 	type want struct {
// 		want int32
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int32) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got int32) error {
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
// 		           data:nil,
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
// 		           data:nil,
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
// 			got := calculateMedian(test.args.data)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_sumHistograms(t *testing.T) {
// 	type args struct {
// 		hist1 []uint64
// 		hist2 []uint64
// 	}
// 	type want struct {
// 		want []uint64
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []uint64) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []uint64) error {
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
// 		           hist1:nil,
// 		           hist2:nil,
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
// 		           hist1:nil,
// 		           hist2:nil,
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
// 			got := sumHistograms(test.args.hist1, test.args.hist2)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mergeInfoIndexStatistics(t *testing.T) {
// 	type args struct {
// 		stats map[string]*payload.Info_Index_Statistics
// 	}
// 	type want struct {
// 		wantMerged *payload.Info_Index_Statistics
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Statistics) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotMerged *payload.Info_Index_Statistics) error {
// 		if !reflect.DeepEqual(gotMerged, w.wantMerged) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotMerged, w.wantMerged)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stats:nil,
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
// 		           stats:nil,
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
// 			gotMerged := mergeInfoIndexStatistics(test.args.stats)
// 			if err := checkFunc(test.want, gotMerged); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_IndexProperty(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		in1 *payload.Empty
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantDetail *payload.Info_Index_PropertyDetail
// 		err        error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_PropertyDetail, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotDetail *payload.Info_Index_PropertyDetail, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		       args: args {
// 		           ctx:nil,
// 		           in1:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 		           in1:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
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
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotDetail, err := s.IndexProperty(test.args.ctx, test.args.in1)
// 			if err := checkFunc(test.want, gotDetail, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
