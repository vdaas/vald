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
	"fmt"
	"slices"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

var emptyPayload = new(payload.Empty) //nolint:gochecknoglobals // shared immutable empty protobuf request

// handleIndexMethodError wraps a method-level index RPC error with the vald.v1.<rpc>
// resource context, records it on the span, and returns the wrapped error (the
// caller returns it alongside a nil payload). Extracted from the verbatim-
// duplicated method-level handling across the Index* handlers.
func (s *server) handleIndexMethodError(
	ctx context.Context, span trace.Span, rpcName string, err error,
) error {
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + rpcName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errGRPCClientConnNotFoundStar):
		err = status.WrapWithInternal(rpcName+" API connection not found", err, resInfo)
		if span != nil {
			attrs = trace.StatusCodeInternal(err.Error())
		}
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(rpcName+" API canceled", err, resInfo)
		if span != nil {
			attrs = trace.StatusCodeCancelled(err.Error())
		}
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(rpcName+" API deadline exceeded", err, resInfo)
		if span != nil {
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		}
	default:
		if span != nil {
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
	}
	log.Debug(err)
	if span != nil {
		errhandler.RecordSpanAttrs(span, attrs, err)
	}
	return err
}

func (s *server) handleIndexBroadCastError(
	sspan trace.Span, target, rpcName string, err error,
) error {
	var (
		attrs trace.Attributes
		st    *status.Status
		msg   string
		code  codes.Code
	)
	switch {
	case errors.Is(err, context.Canceled):
		if sspan != nil {
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix +
					"/vald.v1." + rpcName + ".BroadCast/" +
					target + " canceled: " + err.Error())
		}
		code = codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		if sspan != nil {
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix +
					"/vald.v1." + rpcName + ".BroadCast/" +
					target + " deadline_exceeded: " + err.Error())
		}
		code = codes.DeadlineExceeded
	default:
		st, msg, err = status.ParseError(err, codes.NotFound, "error "+rpcName+" API",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + rpcName + ".BroadCast/" + target,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
			})
		if st != nil {
			code = st.Code()
		} else {
			code = codes.NotFound
		}
		if sspan != nil {
			attrs = trace.FromGRPCStatus(code, msg)
		}
	}
	if sspan != nil {
		errhandler.RecordSpanAttrs(sspan, attrs, err)
	}
	if err != nil && st != nil &&
		code != codes.Canceled &&
		code != codes.DeadlineExceeded &&
		code != codes.InvalidArgument &&
		code != codes.NotFound &&
		code != codes.OK &&
		code != codes.Unimplemented {
		return err
	}
	return nil
}

func (s *server) broadcastIndex(
	ctx context.Context,
	rpcName string,
	fn func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error,
) error {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+rpcName), apiName+"/"+rpcName)
		defer trace.End(span)
	}
	err := s.gateway.BroadCast(ctx, service.READ, func(bctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		sctx := bctx
		var sspan trace.Span
		if trace.Enabled() {
			sctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(bctx, "BroadCast/"+target), apiName+"/"+rpcName+"/"+target)
			defer trace.End(sspan)
		}
		if rerr := fn(sctx, target, vc, copts...); rerr != nil {
			return s.handleIndexBroadCastError(sspan, target, rpcName, rerr)
		}
		return nil
	})
	if err == nil {
		err = ctx.Err()
	}
	if err != nil {
		return s.handleIndexMethodError(ctx, span, rpcName, err)
	}
	return nil
}

func (s *server) IndexInfo(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_Count, err error) {
	var (
		stored, uncommitted atomic.Uint32
		indexing, saving    atomic.Bool
	)
	if err = s.broadcastIndex(ctx, vald.IndexInfoRPCName, func(ctx context.Context, _ string, vc vald.Client, copts ...grpc.CallOption) error {
		info, rerr := vc.IndexInfo(ctx, emptyPayload, copts...)
		if rerr != nil {
			return rerr
		}
		if info != nil {
			stored.Add(info.GetStored())
			uncommitted.Add(info.GetUncommitted())
			if info.GetIndexing() {
				indexing.Store(true)
			}
			if info.GetSaving() {
				saving.Store(true)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &payload.Info_Index_Count{
		Stored:      stored.Load(),
		Uncommitted: uncommitted.Load(),
		Indexing:    indexing.Load(),
		Saving:      saving.Load(),
	}, nil
}

func (s *server) IndexDetail(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_Detail, err error) {
	var (
		mu     sync.Mutex
		detail = &payload.Info_Index_Detail{
			Counts:     make(map[string]*payload.Info_Index_Count),
			Replica:    uint32(s.replica),
			LiveAgents: uint32(s.gateway.GetAgentCount(ctx)),
		}
	)
	if err = s.broadcastIndex(ctx, vald.IndexDetailRPCName, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		info, rerr := vc.IndexInfo(ctx, emptyPayload, copts...)
		if rerr != nil {
			return rerr
		}
		if info != nil {
			mu.Lock()
			detail.Counts[target] = info
			mu.Unlock()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return detail, nil
}

func (s *server) IndexStatistics(
	ctx context.Context, req *payload.Empty,
) (vec *payload.Info_Index_Statistics, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexStatisticsRPCName), apiName+"/"+vald.IndexStatisticsRPCName)
		defer trace.End(span)
	}
	details, err := s.IndexStatisticsDetail(ctx, req)
	if err != nil || details == nil {
		return nil, s.handleIndexMethodError(ctx, span, vald.IndexStatisticsRPCName, err)
	}
	return mergeInfoIndexStatistics(details.GetDetails()), nil
}

func (s *server) IndexStatisticsDetail(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_StatisticsDetail, err error) {
	var (
		mu     sync.Mutex
		detail = &payload.Info_Index_StatisticsDetail{
			Details: make(map[string]*payload.Info_Index_Statistics, s.gateway.GetAgentCount(ctx)),
		}
	)
	if err = s.broadcastIndex(ctx, vald.IndexStatisticsDetailRPCName, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		stats, rerr := vc.IndexStatistics(ctx, emptyPayload, copts...)
		if rerr != nil {
			return rerr
		}
		if stats != nil {
			mu.Lock()
			detail.Details[target] = stats
			mu.Unlock()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return detail, nil
}

func calculateMedian(data []int32) int32 {
	slices.Sort(data)
	n := len(data)
	if n%2 == 0 {
		return (data[n/2-1] + data[n/2]) / 2
	}
	return data[n/2]
}

func sumHistograms(dst, src []uint64) []uint64 {
	if len(src) > len(dst) {
		grown := make([]uint64, len(src))
		copy(grown, dst)
		dst = grown
	}
	for i, v := range src {
		dst[i] += v
	}
	return dst
}

func mergeInfoIndexStatistics(
	stats map[string]*payload.Info_Index_Statistics,
) (merged *payload.Info_Index_Statistics) {
	merged = new(payload.Info_Index_Statistics)

	if len(stats) == 0 {
		return merged
	}

	indegrees := make([]int32, 0, len(stats))
	outdegrees := make([]int32, 0, len(stats))
	indegreeCounts := make([][]int64, 0, len(stats))
	outdegreeHistograms := make([][]uint64, 0, len(stats))
	indegreeHistograms := make([][]uint64, 0, len(stats))
	var maxOutHistLen, maxInHistLen int
	merged.Valid = true

	// first marks the first valid entry: 0 is a legitimate Min*NumberOf*degree,
	// so it cannot double as the "unset" sentinel without making the result
	// depend on map iteration order.
	first := true
	for _, stat := range stats {
		if !stat.Valid {
			continue
		}
		indegrees = append(indegrees, stat.MedianIndegree)
		outdegrees = append(outdegrees, stat.MedianOutdegree)

		indegreeCounts = append(indegreeCounts, stat.IndegreeCount)
		outdegreeHistograms = append(outdegreeHistograms, stat.OutdegreeHistogram)
		indegreeHistograms = append(indegreeHistograms, stat.IndegreeHistogram)
		maxOutHistLen = max(maxOutHistLen, len(stat.OutdegreeHistogram))
		maxInHistLen = max(maxInHistLen, len(stat.IndegreeHistogram))

		if stat.MaxNumberOfIndegree > merged.MaxNumberOfIndegree {
			merged.MaxNumberOfIndegree = stat.MaxNumberOfIndegree
		}
		if stat.MaxNumberOfOutdegree > merged.MaxNumberOfOutdegree {
			merged.MaxNumberOfOutdegree = stat.MaxNumberOfOutdegree
		}
		if first || stat.MinNumberOfIndegree < merged.MinNumberOfIndegree {
			merged.MinNumberOfIndegree = stat.MinNumberOfIndegree
		}
		if first || stat.MinNumberOfOutdegree < merged.MinNumberOfOutdegree {
			merged.MinNumberOfOutdegree = stat.MinNumberOfOutdegree
		}
		first = false
		merged.ModeIndegree += stat.ModeIndegree
		merged.ModeOutdegree += stat.ModeOutdegree
		merged.NodesSkippedFor10Edges += stat.NodesSkippedFor10Edges
		merged.NodesSkippedForIndegreeDistance += stat.NodesSkippedForIndegreeDistance
		merged.NumberOfEdges += stat.NumberOfEdges
		merged.NumberOfIndexedObjects += stat.NumberOfIndexedObjects
		merged.NumberOfNodes += stat.NumberOfNodes
		merged.NumberOfNodesWithoutEdges += stat.NumberOfNodesWithoutEdges
		merged.NumberOfNodesWithoutIndegree += stat.NumberOfNodesWithoutIndegree
		merged.NumberOfObjects += stat.NumberOfObjects
		merged.NumberOfRemovedObjects += stat.NumberOfRemovedObjects
		merged.SizeOfObjectRepository += stat.SizeOfObjectRepository
		merged.SizeOfRefinementObjectRepository += stat.SizeOfRefinementObjectRepository

		merged.VarianceOfIndegree += stat.VarianceOfIndegree
		merged.VarianceOfOutdegree += stat.VarianceOfOutdegree
		merged.MeanEdgeLength += stat.MeanEdgeLength
		merged.MeanEdgeLengthFor10Edges += stat.MeanEdgeLengthFor10Edges
		merged.MeanIndegreeDistanceFor10Edges += stat.MeanIndegreeDistanceFor10Edges
		merged.MeanNumberOfEdgesPerNode += stat.MeanNumberOfEdgesPerNode

		merged.C1Indegree += stat.C1Indegree
		merged.C5Indegree += stat.C5Indegree
		merged.C95Outdegree += stat.C95Outdegree
		merged.C99Outdegree += stat.C99Outdegree
	}

	// every entry was invalid: nothing to merge, same result as an empty input
	if len(indegrees) == 0 {
		return new(payload.Info_Index_Statistics)
	}

	merged.MedianIndegree = calculateMedian(indegrees)
	merged.MedianOutdegree = calculateMedian(outdegrees)
	// size from the longest input so the result does not depend on which agent
	// the map range visits first; shorter inputs are skipped per index below.
	var maxLen int
	for _, count := range indegreeCounts {
		maxLen = max(maxLen, len(count))
	}
	merged.IndegreeCount = make([]int64, maxLen)
	for i := range merged.IndegreeCount {
		var (
			alen int64
			sum  int64
		)
		for _, count := range indegreeCounts {
			if i < len(count) {
				alen++
				sum += count[i]
			}
		}
		merged.IndegreeCount[i] = sum / alen
	}

	if maxOutHistLen > 0 {
		merged.OutdegreeHistogram = make([]uint64, maxOutHistLen)
		for _, hist := range outdegreeHistograms {
			merged.OutdegreeHistogram = sumHistograms(merged.OutdegreeHistogram, hist)
		}
	}

	if maxInHistLen > 0 {
		merged.IndegreeHistogram = make([]uint64, maxInHistLen)
		for _, hist := range indegreeHistograms {
			merged.IndegreeHistogram = sumHistograms(merged.IndegreeHistogram, hist)
		}
	}

	// average over the valid entries only, the same population as the medians
	n := len(indegrees)
	merged.ModeIndegree /= uint64(n)
	merged.ModeOutdegree /= uint64(n)
	merged.VarianceOfIndegree /= float64(n)
	merged.VarianceOfOutdegree /= float64(n)
	merged.MeanEdgeLength /= float64(n)
	merged.MeanEdgeLengthFor10Edges /= float64(n)
	merged.MeanIndegreeDistanceFor10Edges /= float64(n)
	merged.MeanNumberOfEdgesPerNode /= float64(n)
	merged.C1Indegree /= float64(n)
	merged.C5Indegree /= float64(n)
	merged.C95Outdegree /= float64(n)
	merged.C99Outdegree /= float64(n)

	return merged
}

func (s *server) IndexProperty(
	ctx context.Context, _ *payload.Empty,
) (detail *payload.Info_Index_PropertyDetail, err error) {
	var mu sync.Mutex
	detail = &payload.Info_Index_PropertyDetail{
		Details: make(map[string]*payload.Info_Index_Property, s.gateway.GetAgentCount(ctx)),
	}
	if err = s.broadcastIndex(ctx, vald.IndexPropertyRPCName, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		prop, rerr := vc.IndexProperty(ctx, emptyPayload, copts...)
		if rerr != nil {
			return rerr
		}
		if prop != nil {
			mu.Lock()
			for key, value := range prop.Details {
				detail.Details[target+"-"+key] = value
			}
			mu.Unlock()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return detail, nil
}
