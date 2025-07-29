//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package grpc

import (
	"context"
	"fmt"
	"slices"
	"sync/atomic"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

func (s *server) IndexInfo(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexInfoRPCName), apiName+"/"+vald.IndexInfoRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var (
		stored, uncommitted atomic.Uint32
		indexing, saving    atomic.Bool
	)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexInfoRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			info, err := vc.IndexInfo(sctx, new(payload.Empty), copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexInfoRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexInfoRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexInfoRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexInfoRPCName + ".BroadCase/" + target,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexInfoRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexInfoRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexInfoRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexInfoRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		log.Debug(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexDetailRPCName), apiName+"/"+vald.IndexDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var (
		mu     sync.Mutex
		detail = &payload.Info_Index_Detail{
			Counts:     make(map[string]*payload.Info_Index_Count),
			Replica:    uint32(s.replica),
			LiveAgents: uint32(s.gateway.GetAgentCount(ctx)),
		}
	)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexDetailRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			info, err := vc.IndexInfo(sctx, new(payload.Empty), copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexDetailRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexDetailRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexDetailRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexDetailRPCName + ".BroadCase/" + target,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			if info != nil {
				mu.Lock()
				detail.Counts[target] = info
				mu.Unlock()
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexDetailRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexDetailRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexDetailRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexDetailRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		log.Debug(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return detail, nil
}

func (s *server) IndexStatistics(
	ctx context.Context, req *payload.Empty,
) (vec *payload.Info_Index_Statistics, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexStatisticsRPCName), apiName+"/"+vald.IndexStatisticsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	details, err := s.IndexStatisticsDetail(ctx, req)
	if err != nil || details == nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexStatisticsRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexStatisticsRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexStatisticsRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexStatisticsRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		log.Debug(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return mergeInfoIndexStatistics(details.GetDetails()), nil
}

func (s *server) IndexStatisticsDetail(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_StatisticsDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexStatisticsDetailRPCName), apiName+"/"+vald.IndexStatisticsDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var (
		mu     sync.Mutex
		detail = &payload.Info_Index_StatisticsDetail{
			Details: make(map[string]*payload.Info_Index_Statistics, s.gateway.GetAgentCount(ctx)),
		}
	)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexStatisticsDetailRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			var stats *payload.Info_Index_Statistics
			stats, err = vc.IndexStatistics(sctx, new(payload.Empty), copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexStatisticsDetailRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexStatisticsDetailRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexStatisticsDetailRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexStatisticsDetailRPCName + ".BroadCase/" + target,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			if stats != nil {
				mu.Lock()
				detail.Details[target] = stats
				mu.Unlock()
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexStatisticsDetailRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexStatisticsDetailRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexStatisticsDetailRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexStatisticsDetailRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		log.Debug(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
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

func sumHistograms(hist1, hist2 []uint64) []uint64 {
	if len(hist1) < len(hist2) {
		hist1, hist2 = hist2, hist1
	}
	for i := range hist2 {
		hist1[i] += hist2[i]
	}
	return hist1
}

func mergeInfoIndexStatistics(
	stats map[string]*payload.Info_Index_Statistics,
) (merged *payload.Info_Index_Statistics) {
	merged = new(payload.Info_Index_Statistics)

	if len(stats) == 0 {
		return merged
	}

	var indegrees, outdegrees []int32
	var indegreeCounts [][]int64
	var outdegreeHistograms, indegreeHistograms [][]uint64
	merged.Valid = true

	for _, stat := range stats {
		if !stat.Valid {
			continue
		}
		indegrees = append(indegrees, stat.MedianIndegree)
		outdegrees = append(outdegrees, stat.MedianOutdegree)

		indegreeCounts = append(indegreeCounts, stat.IndegreeCount)
		outdegreeHistograms = append(outdegreeHistograms, stat.OutdegreeHistogram)
		indegreeHistograms = append(indegreeHistograms, stat.IndegreeHistogram)

		if stat.MaxNumberOfIndegree > merged.MaxNumberOfIndegree {
			merged.MaxNumberOfIndegree = stat.MaxNumberOfIndegree
		}
		if stat.MaxNumberOfOutdegree > merged.MaxNumberOfOutdegree {
			merged.MaxNumberOfOutdegree = stat.MaxNumberOfOutdegree
		}
		if stat.MinNumberOfIndegree < merged.MinNumberOfIndegree || merged.MinNumberOfIndegree == 0 {
			merged.MinNumberOfIndegree = stat.MinNumberOfIndegree
		}
		if stat.MinNumberOfOutdegree < merged.MinNumberOfOutdegree || merged.MinNumberOfOutdegree == 0 {
			merged.MinNumberOfOutdegree = stat.MinNumberOfOutdegree
		}
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

	merged.MedianIndegree = calculateMedian(indegrees)
	merged.MedianOutdegree = calculateMedian(outdegrees)
	merged.IndegreeCount = make([]int64, len(indegreeCounts[0]))
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

	for _, hist := range outdegreeHistograms {
		merged.OutdegreeHistogram = sumHistograms(merged.OutdegreeHistogram, hist)
	}

	for _, hist := range indegreeHistograms {
		merged.IndegreeHistogram = sumHistograms(merged.IndegreeHistogram, hist)
	}

	merged.ModeIndegree /= uint64(len(stats))
	merged.ModeOutdegree /= uint64(len(stats))
	merged.VarianceOfIndegree /= float64(len(stats))
	merged.VarianceOfOutdegree /= float64(len(stats))
	merged.MeanEdgeLength /= float64(len(stats))
	merged.MeanEdgeLengthFor10Edges /= float64(len(stats))
	merged.MeanIndegreeDistanceFor10Edges /= float64(len(stats))
	merged.MeanNumberOfEdgesPerNode /= float64(len(stats))
	merged.C1Indegree /= float64(len(stats))
	merged.C5Indegree /= float64(len(stats))
	merged.C95Outdegree /= float64(len(stats))
	merged.C99Outdegree /= float64(len(stats))

	return merged
}

func (s *server) IndexProperty(
	ctx context.Context, _ *payload.Empty,
) (detail *payload.Info_Index_PropertyDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexPropertyRPCName), apiName+"/"+vald.IndexPropertyRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var mu sync.Mutex
	detail = &payload.Info_Index_PropertyDetail{
		Details: make(map[string]*payload.Info_Index_Property, s.gateway.GetAgentCount(ctx)),
	}

	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexStatisticsDetailRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			var prop *payload.Info_Index_PropertyDetail
			prop, err = vc.IndexProperty(sctx, new(payload.Empty), copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled), errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName + ".BroadCast/" + target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded), errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName + ".BroadCast/" + target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexPropertyRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName + ".BroadCast/" + target,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st != nil && code != codes.Canceled && code != codes.DeadlineExceeded && code != codes.InvalidArgument && code != codes.NotFound && code != codes.OK && code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if prop != nil {
				mu.Lock()
				for key, value := range prop.Details {
					detail.Details[target+"-"+key] = value
				}
				mu.Unlock()
			}
			return nil
		})
		return nil
	}))
	for i := 0; i < s.gateway.GetAgentCount(ctx); i++ {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			log.Errorf("!!! IndexProperty API canceled: %v", err)
			break
		case err = <-ech:
			log.Errorf("!!! IndexProperty API error: %v", err)
		}
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexPropertyRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexPropertyRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexPropertyRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		log.Debug(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return detail, nil
}
