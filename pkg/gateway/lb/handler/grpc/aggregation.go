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
	"cmp"
	"context"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type Aggregator interface {
	Send(ctx context.Context, data *payload.Search_Response)
	GetNum() int  // get top-k number
	GetFnum() int // get forwarding top-k number calculated by search ratio
	Result() *payload.Search_Response
	hasResults() bool
	uniqueCount() int
}

// DistPayload pairs a search result with its distance widened to float64.
// float32 -> float64 is exact, so comparing the float64 values orders results
// exactly like the float32 originals (and like the former big.Float
// comparison: -0 == +0, -Inf < finite < +Inf) without allocating.
type DistPayload struct {
	raw      *payload.Object_Distance
	distance float64
}

// rankable reports whether dist may take part in result ranking.
// NaN policy: a result whose distance is NaN is skipped (never aggregated, never
// marked visited). NaN has no rank, so keeping it would break the sorted
// invariant the aggregators rely on, while the remaining results stay valid.
// The old big.NewFloat(NaN) comparison panicked; the recovered panic dropped the
// agent's remaining results and turned the broadcast into an error.
func rankable(dist *payload.Object_Distance) bool {
	return dist != nil && !math.IsNaN(float64(dist.GetDistance()))
}

// cmpDist orders DistPayloads by ascending distance. Callers must only pass
// rankable (non-NaN) payloads.
func cmpDist(l, r DistPayload) int {
	return cmp.Compare(l.distance, r.distance)
}

const (
	defaultAggrFailedCap  = 8
	defaultAggrReachedCap = 8
	defaultAggrSeenCap    = 64
)

type searchTargetTracker struct {
	reached         map[string]struct{}
	failed          map[string]codes.Code
	fastReached     atomic.Int64
	hasSlow         atomic.Bool
	answered        atomic.Int64
	qtimer          atomic.Pointer[time.Timer]
	wave1Incomplete atomic.Bool
	mu              sync.Mutex
}

var (
	aggrFailedPool = sync.Pool[map[string]codes.Code]{ //nolint:gochecknoglobals // failure status maps are reused per search request
		New: func() map[string]codes.Code {
			return make(map[string]codes.Code, defaultAggrFailedCap)
		},
	}
	searchTrackerPool = sync.Pool[*searchTargetTracker]{ //nolint:gochecknoglobals // target tracking state is reused per search request
		New: func() *searchTargetTracker {
			return &searchTargetTracker{
				reached: make(map[string]struct{}, defaultAggrReachedCap),
			}
		},
	}
	aggrSeenPool = sync.Pool[map[string]struct{}]{ //nolint:gochecknoglobals // dedup id sets are reused per search response
		New: func() map[string]struct{} {
			return make(map[string]struct{}, defaultAggrSeenCap)
		},
	}
)

func putAggrSeen(seen map[string]struct{}) {
	if seen != nil && len(seen) <= maxAggrPrealloc {
		clear(seen)
		aggrSeenPool.Put(seen)
	}
}

func forwardSearchConfig(bcfg *payload.Search_Config, fnum int) *payload.Search_Config {
	// Forwarding Config to Agent: clone so the caller's config is never aliased or mutated.
	fcfg := bcfg.CloneVT()
	fcfg.Num = uint32(max(fnum, 0)) //nolint:gosec // max(fnum, 0) guarantees non-negative value within int->uint32 search top-k range
	fcfg.MinNum = 0
	return fcfg
}

func releaseSearchTargetTracker(tr *searchTargetTracker) {
	if tr.failed != nil {
		if len(tr.failed) <= maxAggrPrealloc {
			clear(tr.failed)
			aggrFailedPool.Put(tr.failed)
		}
		tr.failed = nil
	}
	if len(tr.reached) > maxAggrPrealloc {
		tr.reached = make(map[string]struct{}, defaultAggrReachedCap)
	} else {
		clear(tr.reached)
	}
	tr.fastReached.Store(0)
	tr.hasSlow.Store(false)
	tr.answered.Store(0)
	tr.qtimer.Store(nil)
	tr.wave1Incomplete.Store(false)
	searchTrackerPool.Put(tr)
}

func recordSearchTarget(tr *searchTargetTracker, target string, agentErr error) {
	if agentErr == nil && !tr.hasSlow.Load() {
		tr.fastReached.Add(1)
		return
	}
	tr.hasSlow.Store(true)
	tr.mu.Lock()
	tr.reached[target] = struct{}{}
	if code, ok := searchAgentFailure(agentErr); ok {
		if tr.failed == nil {
			tr.failed = aggrFailedPool.Get()
		}
		tr.failed[target] = code
	} else if tr.failed != nil {
		delete(tr.failed, target)
	}
	tr.mu.Unlock()
}

func isTolerableSearchMiss(err error) bool {
	if errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
		return false
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		return false
	}
	code := st.Code()
	return code == codes.NotFound || code == codes.InvalidArgument
}

func (s *server) searchSubsetCount(rpc string, agents int, bcfg *payload.Search_Config) int {
	const searchExtraReplicaCover = 1
	if s.replica <= 1 || agents <= s.replica || bcfg.GetRatio() != nil {
		return 0
	}
	switch rpc {
	case vald.LinearSearchRPCName, vald.LinearSearchByIDRPCName:
		return agents - s.replica + 1
	case vald.SearchRPCName, vald.SearchByIDRPCName:
		if s.searchQuorumGrace <= 0 && agents-s.replica+searchExtraReplicaCover < agents {
			return agents - s.replica + searchExtraReplicaCover
		}
	}
	return 0
}

func (s *server) executeSearchBroadcast(
	ctx, bctx context.Context,
	addrs []string,
	subsetCount, agents, num, minNum int,
	quorum int,
	aggr Aggregator,
	tr *searchTargetTracker,
	call func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (expectedAgents int, quorumReached bool, err error) {
	var (
		bcastCtx = bctx
		bcancel  context.CancelCauseFunc
	)
	if subsetCount == 0 && quorum > 0 {
		// bcastCtx scopes only the broadcast, so the quorum early termination below stops
		// the remaining agent calls while the aggregator keeps every result it got.
		bcastCtx, bcancel = context.WithCancelCause(bctx)
	}
	cb := func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		sctx := ctx
		var sspan trace.Span
		if trace.Enabled() {
			sctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/aggregationSearch/"+target)
			defer trace.End(sspan)
		}
		r, rerr := call(sctx, vc, copts...)
		if rerr != nil {
			berr := handleSearchBroadCastError(sspan, target, rerr)
			if subsetCount > 0 && !isTolerableSearchMiss(rerr) {
				tr.wave1Incomplete.Store(true)
			}
			recordSearchTarget(tr, target, rerr)
			return berr
		}
		if r == nil || len(r.GetResults()) == 0 {
			if sspan != nil {
				nferr := status.WrapWithNotFound(fmt.Sprintf("failed to process search request from %s", target),
					errors.ErrEmptySearchResult,
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					})
				errhandler.RecordSpanError(sspan, codes.NotFound, nferr)
			}
			recordSearchTarget(tr, target, nil)
			return nil
		}
		aggr.Send(sctx, r)
		recordSearchTarget(tr, target, nil)
		// Count an agent only after all its results were sent, so terminating
		// early never drops results of an agent that counts toward the quorum.
		if bcancel != nil && int(tr.answered.Add(1)) == quorum {
			tr.qtimer.Store(time.AfterFunc(s.searchQuorumGrace, func() {
				bcancel(errSearchQuorumReached)
			}))
		}
		return nil
	}

	bctxPolicy := grpc.WithSearchBroadcastPolicy(bcastCtx)
	if subsetCount > 0 {
		expectedAgents = subsetCount
		offset := s.gateway.NextRotateOffset(len(addrs))
		var sliced bool
		sliced, err = s.gateway.BroadCastSlice(bctxPolicy, service.READ, addrs, offset, subsetCount, cb)
		if !sliced {
			expectedAgents = int(tr.fastReached.Load())
			if tr.hasSlow.Load() {
				tr.mu.Lock()
				expectedAgents += len(tr.reached) + len(tr.failed)
				tr.mu.Unlock()
			}
			return expectedAgents, false, err
		}
		if ctx.Err() == nil && bctx.Err() == nil && (err == nil || !errors.Is(err, errGRPCClientConnNotFoundStar)) {
			var needFallback bool
			if !tr.hasSlow.Load() {
				needFallback = err != nil || tr.wave1Incomplete.Load() || int(tr.fastReached.Load()) < subsetCount
			} else {
				tr.mu.Lock()
				reachedCount := int(tr.fastReached.Load()) + len(tr.reached)
				needFallback = err != nil || tr.wave1Incomplete.Load() || len(tr.failed) > 0 || reachedCount < subsetCount
				tr.mu.Unlock()
			}
			if !needFallback {
				if (num != 0 && !aggr.hasResults()) || (0 < minNum && aggr.uniqueCount() < minNum) {
					needFallback = true
				}
			}
			if needFallback {
				expectedAgents = len(addrs)
				remOffset := (offset + subsetCount) % len(addrs)
				_, err2 := s.gateway.BroadCastSlice(bctxPolicy, service.READ_PRIMARY, addrs, remOffset, len(addrs)-subsetCount, cb)
				if err == nil {
					err = err2
				} else if err2 != nil {
					err = errors.Join(err, err2)
				}
			}
		}
		return expectedAgents, false, err
	}

	expectedAgents = agents
	// Per-agent failures are tolerated and summarized once per request below, so
	// the client must not log every failed agent/attempt at Warn/Error.
	// Replicas cover a failed agent, so an overloaded one (ResourceExhausted) is not
	// retried inside the search deadline (write paths keep backoff-and-wait).
	err = s.gateway.BroadCast(bctxPolicy, service.READ, cb)
	if t := tr.qtimer.Load(); t != nil {
		t.Stop()
	}
	quorumReached = bcancel != nil && errors.Is(context.Cause(bcastCtx), errSearchQuorumReached)
	if bcancel != nil {
		bcancel(nil)
	}
	if quorumReached && errors.Is(err, context.Canceled) {
		// the agents still running were canceled on purpose, not failed
		err = nil
	}
	return expectedAgents, quorumReached, err
}

func (s *server) aggregationSearch(
	ctx context.Context,
	rpc string, // agent RPC name call invokes, used as the metrics label
	agents int, // number of agents the broadcast should reach
	aggr Aggregator,
	bcfg *payload.Search_Config, // Base Config of Request
	call func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (res *payload.Search_Response, attrs []attribute.KeyValue, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "aggregationSearch"), apiName+"/aggregationSearch")
		defer trace.End(span)
	}

	if bcfg == nil {
		return nil, nil, errors.ErrInvalidSearchConfig("bcfg is nil in aggregationSearch")
	}

	num := aggr.GetNum()
	minNum := int(bcfg.GetMinNum())

	var timeout time.Duration
	if to := bcfg.GetTimeout(); to != 0 {
		timeout = time.Duration(to)
	} else {
		timeout = s.timeout
	}

	tr := searchTrackerPool.Get()
	if tr == nil {
		tr = &searchTargetTracker{
			reached: make(map[string]struct{}, defaultAggrReachedCap),
		}
	}
	parentCtx := ctx
	bctx := parentCtx
	var cancel context.CancelFunc
	if dl, ok := parentCtx.Deadline(); !ok || time.Until(dl) > timeout {
		bctx, cancel = context.WithTimeout(parentCtx, timeout)
	}

	var (
		addrs       []string
		subsetCount int
	)
	if s.replica > 1 && bcfg.GetRatio() == nil &&
		(s.searchQuorumGrace <= 0 || rpc == vald.LinearSearchRPCName || rpc == vald.LinearSearchByIDRPCName) {
		addrs = s.gateway.Addrs(bctx)
		subsetCount = s.searchSubsetCount(rpc, len(addrs), bcfg)
	}

	quorum := 0
	if subsetCount == 0 {
		quorum = s.searchQuorum(agents)
	}
	expectedAgents, quorumReached, err := s.executeSearchBroadcast(
		parentCtx, bctx, addrs, subsetCount, agents, num, minNum, quorum, aggr, tr, call,
	)
	bctxErr := bctx.Err()
	if cancel != nil {
		cancel()
	}
	var unreached, failures int
	if !tr.hasSlow.Load() {
		if bctxErr == nil && !quorumReached {
			unreached = max(expectedAgents-int(tr.fastReached.Load()), 0)
		}
		failures = unreached
	} else {
		tr.mu.Lock()
		// The client never runs the callback for a target whose circuit breaker is
		// open or that has no connection, so such agents are neither in failed nor
		// in reached. They are counted as unreached agents with code Unavailable.
		// Once the request ctx has ended, targets not reached yet are deadline or
		// cancel misses, which are not agent failures, so none are counted then.
		if bctxErr == nil && !quorumReached {
			reachedCount := int(tr.fastReached.Load()) + len(tr.reached)
			unreached = max(expectedAgents-reachedCount, 0)
		}
		failures = len(tr.failed) + unreached
		for _, code := range tr.failed {
			s.searchMetrics.RecordAgentFailure(rpc, code)
		}
		tr.mu.Unlock()
	}
	releaseSearchTargetTracker(tr)
	for range unreached {
		s.searchMetrics.RecordAgentFailure(rpc, codes.Unavailable)
	}
	var bcErr error // broadcast error of a request that still has results
	if err != nil && !aggr.hasResults() && errors.Is(err, errGRPCClientConnNotFoundStar) {
		err = status.WrapWithInternal("search API connection not found", err,
			s.errDetails(parentCtx, "search", bcfg.GetRequestId(), bcfg)...)
		if span != nil {
			attrs = trace.StatusCodeInternal(err.Error())
			errhandler.RecordSpanAttrs(span, attrs, err)
		}
		logPartialSearch(rpc, failures, s.replica, bcErr, err)
		return nil, attrs, err
	}
	res, attrs, bcErr, err = s.finalizeAggregationResponse(parentCtx, bctxErr, span, rpc, bcfg, aggr.Result(), num, minNum, failures, err)
	logPartialSearch(rpc, failures, s.replica, bcErr, err)
	return res, attrs, err
}

func (s *server) finalizeAggregationResponse(
	ctx context.Context,
	bctxErr error,
	span trace.Span,
	rpc string,
	bcfg *payload.Search_Config,
	res *payload.Search_Response,
	num, minNum, failures int,
	err error,
) (_ *payload.Search_Response, attrs []attribute.KeyValue, bcErr, reqErr error) {
	if num != 0 && len(res.GetResults()) > num {
		res.Results = res.GetResults()[:num]
	}

	if errors.Is(bctxErr, context.DeadlineExceeded) {
		if len(res.GetResults()) == 0 {
			reqErr = status.WrapWithDeadlineExceeded(
				"error search result length is 0 due to the timeoutage limit",
				errors.ErrEmptySearchResult,
				append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...,
			)
			if span != nil {
				attrs = trace.StatusCodeDeadlineExceeded(reqErr.Error())
				errhandler.RecordSpanAttrs(span, attrs, reqErr)
			}
			return nil, attrs, nil, reqErr
		}
		if 0 < minNum && len(res.GetResults()) < minNum {
			reqErr = status.WrapWithDeadlineExceeded(
				fmt.Sprintf("error search result length is not enough due to the timeoutage limit, required: %d, found: %d", minNum, len(res.GetResults())),
				errors.ErrInsuffcientSearchResult,
				append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...,
			)
			if span != nil {
				attrs = trace.StatusCodeDeadlineExceeded(reqErr.Error())
				errhandler.RecordSpanAttrs(span, attrs, reqErr)
			}
			return nil, attrs, nil, reqErr
		}
	}

	if err != nil {
		st, msg, perr := status.ParseError(err, codes.Internal,
			"failed to parse search gRPC error response",
			append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...)
		if span != nil {
			attrs = trace.FromGRPCStatus(st.Code(), msg)
			errhandler.RecordSpanAttrs(span, attrs, perr)
		}
		if len(res.GetResults()) == 0 {
			// with agent failures the deferred summary logs err at Warn
			if !searchSummaryLogsRequestErr(failures, perr) {
				log.Warn(perr)
			}
			return nil, attrs, nil, perr
		}
		// some agents failed but the others returned results: the severity is
		// decided below once the request is known to succeed.
		bcErr = perr
	}
	if num != 0 && len(res.GetResults()) == 0 {
		if err == nil {
			err = errors.ErrEmptySearchResult
		}
		reqErr = status.WrapWithNotFound("error search result length is 0", err,
			append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...)
		if span != nil {
			attrs = trace.StatusCodeNotFound(reqErr.Error())
			errhandler.RecordSpanAttrs(span, attrs, reqErr)
		}
		return nil, attrs, bcErr, reqErr
	}

	if 0 < minNum && len(res.GetResults()) < minNum {
		reqErr = status.WrapWithNotFound(
			fmt.Sprintf("error search result length is not enough required: %d, found: %d", minNum, len(res.GetResults())),
			errors.ErrInsuffcientSearchResult,
			append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...,
		)
		// with agent failures the deferred summary logs bcErr with err at Warn
		if bcErr != nil && !searchSummaryLogsRequestErr(failures, reqErr) {
			log.Warn(bcErr)
		}
		if span != nil {
			attrs = trace.StatusCodeNotFound(reqErr.Error())
			errhandler.RecordSpanAttrs(span, attrs, reqErr)
		}
		return nil, attrs, bcErr, reqErr
	}
	if failures > 0 {
		s.searchMetrics.RecordPartialSuccess(rpc, searchRecallAtRisk(failures, s.replica))
	}
	res.RequestId = bcfg.GetRequestId()
	return res, attrs, bcErr, nil
}

// searchRecallAtRisk reports whether failures agents failing may have made some
// vectors unreachable. Each vector is stored on replica agents, so while fewer
// than replica agents fail every vector still has a live replica and the partial
// result keeps the recall of a full broadcast.
func searchRecallAtRisk(failures, replica int) bool {
	return failures >= replica
}

// quorumError is the cancel cause of a search broadcast ended early by the quorum.
type quorumError string

func (e quorumError) Error() string { return string(e) }

const errSearchQuorumReached quorumError = "search quorum reached"

// searchQuorum returns how many agents must answer with results before a search
// may stop waiting for the rest: agents - replica + 1 answered agents include at
// least one replica of every index replica set. 0 means wait for every agent:
// the early termination is disabled, or replication leaves nothing to skip.
func (s *server) searchQuorum(agents int) int {
	if s.searchQuorumGrace <= 0 || s.replica <= 1 || agents <= s.replica {
		return 0
	}
	return agents - s.replica + 1
}

// searchAgentFailure reports whether an agent's search error counts as an agent
// failure, and with which code. Every error is a failure except the
// Canceled/DeadlineExceeded/NotFound misses, which are expected when the request
// deadline cuts the broadcast short or an agent has no neighbor, and
// InvalidArgument, which the request caused (e.g. a vector dimension the index
// does not have) on otherwise healthy agents. A non-status
// error is counted as Unknown. Whether the error fails the BroadCast is decided
// separately by handleSearchBroadCastError.
func searchAgentFailure(err error) (codes.Code, bool) {
	if err == nil ||
		errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
		return codes.OK, false
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		return codes.Unknown, true
	}
	if code := st.Code(); code != codes.OK && code != codes.Canceled &&
		code != codes.DeadlineExceeded && code != codes.NotFound &&
		code != codes.InvalidArgument {
		return code, true
	}
	return codes.OK, false
}

// partialSearchLogLevel returns the level logPartialSearch logs a request at,
// or level.Unknown when there is nothing to log. failures counts the failed and
// unreached agents, bcErr is the BroadCast error of a request that still returns
// results and reqErr the error the request returns.
//   - A failing request with agent failures is Warn; this summary is then its
//     only Warn line (see searchSummaryLogsRequestErr).
//   - A successful request is Warn when the failures may have cost recall
//     (failures >= replica) or when BroadCast failed without any agent failure
//     (e.g. a local ctx error), and Debug while the replicas cover the failures.
func partialSearchLogLevel(failures, replica int, bcErr, reqErr error) level.Level {
	switch {
	case reqErr != nil:
		if searchSummaryLogsRequestErr(failures, reqErr) {
			return level.WARN
		}
		return level.Unknown
	case failures == 0 && bcErr == nil:
		return level.Unknown
	case failures == 0, searchRecallAtRisk(failures, replica):
		return level.WARN
	}
	return level.DEBUG
}

// searchSummaryLogsRequestErr reports whether logPartialSearch logs the failing
// request's error reqErr at Warn, in which case aggregationSearch must not Warn
// the same request again, so a request that failed with agent failures is Warn-logged once.
// A failing request without agent failures (a deadline, no result) is not logged here.
func searchSummaryLogsRequestErr(failures int, reqErr error) bool {
	return reqErr != nil && failures > 0
}

// logPartialSearch logs the agent failures of a search request once, at the
// level partialSearchLogLevel decides.
func logPartialSearch(rpc string, failures, replica int, bcErr, reqErr error) {
	const (
		partial = "%s returned partial results: %d agent(s) failed or unreached, index replica %d: %v"
		failed  = "%s failed with %d agent(s) failed or unreached, index replica %d: %v"
		// a failing request whose BroadCast error was not returned as reqErr
		failedBC = failed + ", broadcast error: %v"
	)
	lv := partialSearchLogLevel(failures, replica, bcErr, reqErr)
	if lv == level.WARN && reqErr != nil && bcErr != nil {
		log.Warnf(failedBC, rpc, failures, replica, reqErr, bcErr)
	} else if lv == level.WARN && reqErr != nil {
		log.Warnf(failed, rpc, failures, replica, reqErr)
	} else if lv == level.WARN {
		log.Warnf(partial, rpc, failures, replica, bcErr)
	} else if lv == level.DEBUG {
		log.Debugf(partial, rpc, failures, replica, bcErr)
	}
}

// vald standard algorithm.
type valdStdAggr struct {
	// maxDist holds the math.Float64bits of the current distance bound. newStd
	// initializes it; Send updates it atomically under mu and reads it lock-free
	// on entry to prune responses whose closest candidate already exceeds the bound.
	maxDist atomic.Uint64
	seen    map[string]struct{}
	result  []*payload.Object_Distance
	mu      sync.Mutex
	num     int
	fnum    int
	closed  atomic.Bool
}

// handleSearchBroadCastError classifies a per-target BroadCast search RPC error:
// it records the error on span and returns a non-nil error only for codes that
// must fail the whole BroadCast (Internal/Unavailable/ResourceExhausted/
// InvalidArgument/Aborted); every
// other outcome (Canceled/DeadlineExceeded/NotFound/etc.) is a tolerable
// per-target miss and returns nil. Errors are logged at Debug only: the request
// summary (logPartialSearch) decides whether the agent failures warrant a Warn.
// Extracted from the previously duplicated initial-call and retry-call error
// handling in aggregationSearch.
func handleSearchBroadCastError(sspan trace.Span, target string, err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		if sspan != nil {
			sspan.RecordError(err)
			sspan.SetAttributes(trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix +
					"/vald.v1.search.BroadCast/" +
					target + " canceled: " + err.Error())...)
			sspan.SetStatus(trace.StatusError, err.Error())
		}
		return nil
	case errors.Is(err, context.DeadlineExceeded):
		if sspan != nil {
			sspan.RecordError(err)
			sspan.SetAttributes(trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix +
					"/vald.v1.search.BroadCast/" +
					target + " deadline_exceeded: " + err.Error())...)
			sspan.SetStatus(trace.StatusError, err.Error())
		}
		return nil
	default:
		st, ok := status.FromError(err)
		if !ok {
			log.Debug(err)
			return nil
		}
		if sspan != nil {
			sspan.RecordError(err)
			sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			sspan.SetStatus(trace.StatusError, err.Error())
		}
		switch st.Code() {
		// InvalidArgument (e.g. a dimension the agents' index does not have) and
		// Aborted (indexing/flushing in progress) are returned too: when no agent
		// produces results, the client must get that code, not an empty-result NotFound.
		case codes.Internal,
			codes.Unavailable,
			codes.ResourceExhausted,
			codes.InvalidArgument,
			codes.Aborted:
			log.Debug(err)
			return err
		case codes.NotFound,
			codes.Canceled,
			codes.DeadlineExceeded:
			return nil
		}
		log.Debug(err)
	}
	return nil
}

// maxAggrPrealloc caps the capacity the aggregators preallocate from the
// client-supplied num (and replica count). num is not bounded by the server, so
// sizing buffers from it directly lets a single small request allocate
// gigabytes before any agent replies. The capacities are only hints: the
// buffers still grow by append when more results arrive.
const maxAggrPrealloc = 1 << 12

// preallocLen returns n clamped to [0, maxAggrPrealloc] for use as a capacity hint.
func preallocLen(n int) int {
	return max(min(n, maxAggrPrealloc), 0)
}

// newStd returns a valdStdAggr configured for the standard aggregation algorithm.
// It sets the top-k size (`num`), the forwarding top-k size (`fnum`), and preallocates
// the result slice by num+1 (by num*replica when num is not positive), capped at maxAggrPrealloc.
func newStd(num, fnum, replica int) Aggregator {
	rhint := preallocLen(num * replica)
	if num > 0 {
		rhint = preallocLen(num + 1)
	}
	vsa := &valdStdAggr{
		num:    num,
		fnum:   fnum,
		result: make([]*payload.Object_Distance, 0, rhint),
	}
	vsa.maxDist.Store(math.Float64bits(math.MaxFloat64))
	return vsa
}

func (v *valdStdAggr) Send(ctx context.Context, data *payload.Search_Response) {
	if v.closed.Load() || ctx.Err() != nil {
		return
	}
	result := data.GetResults()
	if v.fnum > 0 && len(result) > v.fnum {
		result = result[:v.fnum]
	}
	if len(result) == 0 {
		return
	}
	if rankable(result[0]) && float64(result[0].GetDistance()) >= math.Float64frombits(v.maxDist.Load()) {
		return
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed.Load() {
		return
	}
	if v.seen == nil {
		v.seen = aggrSeenPool.Get()
	}
	fmax := math.Float64frombits(v.maxDist.Load())
	var dirtyMax bool
	for _, dist := range result {
		if !rankable(dist) {
			continue
		}
		distance := float64(dist.GetDistance())
		if distance >= fmax {
			break
		}
		id := dist.GetId()
		if _, ok := v.seen[id]; ok {
			continue
		}
		v.seen[id] = struct{}{}
		rl := len(v.result)
		var i int
		if rl == 0 || distance >= float64(v.result[rl-1].GetDistance()) {
			i = rl
		} else {
			i, _ = slices.BinarySearchFunc(v.result, distance, func(elem *payload.Object_Distance, target float64) int {
				if float64(elem.GetDistance()) <= target {
					return -1
				}
				return 1
			})
		}
		if v.num != 0 && rl >= v.num {
			if i < v.num {
				copy(v.result[i+1:v.num], v.result[i:v.num-1])
				v.result[i] = dist
				if distEnd := float64(v.result[v.num-1].GetDistance()); distEnd < fmax {
					fmax = distEnd
					dirtyMax = true
				}
			}
			continue
		}
		v.result = append(v.result, nil)
		if i < rl {
			copy(v.result[i+1:], v.result[i:rl])
		}
		v.result[i] = dist
		rl++
		if v.num != 0 && rl >= v.num {
			if distEnd := float64(v.result[rl-1].GetDistance()); distEnd < fmax {
				fmax = distEnd
				dirtyMax = true
			}
		}
	}
	if dirtyMax {
		v.maxDist.Store(math.Float64bits(fmax))
	}
}

func (v *valdStdAggr) Result() *payload.Search_Response {
	v.closed.Store(true)
	v.mu.Lock()
	putAggrSeen(v.seen)
	v.seen = nil
	res := v.result
	v.result = nil
	v.mu.Unlock()
	if v.num != 0 && len(res) > v.num {
		res = res[:v.num]
	}
	return &payload.Search_Response{
		Results: res,
	}
}

func (v *valdStdAggr) hasResults() bool {
	v.mu.Lock()
	ok := len(v.result) > 0
	v.mu.Unlock()
	return ok
}

func (v *valdStdAggr) uniqueCount() int {
	v.mu.Lock()
	n := len(v.result)
	v.mu.Unlock()
	return n
}

func (v *valdStdAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdStdAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

// pairing heap.
type valdPairingHeapAggr struct {
	maxDist atomic.Uint64
	nodes   []pairingNode
	visited map[string]struct{} // guarded by mu
	result  []*payload.Object_Distance
	root    int32
	num     int
	fnum    int
	mu      sync.Mutex
	closed  atomic.Bool
	filled  atomic.Bool
}

func newPairingHeap(num, fnum, _ int) Aggregator {
	return &valdPairingHeapAggr{
		root: pairingNull,
		num:  num,
		fnum: fnum,
	}
}

func (v *valdPairingHeapAggr) Send(ctx context.Context, data *payload.Search_Response) {
	if v.closed.Load() || ctx.Err() != nil {
		return
	}
	result := data.GetResults()
	if v.fnum > 0 && len(result) > v.fnum {
		result = result[:v.fnum]
	}
	if len(result) == 0 {
		return
	}
	if v.filled.Load() && rankable(result[0]) && float64(result[0].GetDistance()) >= math.Float64frombits(v.maxDist.Load()) {
		return
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed.Load() {
		return
	}
	if v.visited == nil {
		v.visited = aggrSeenPool.Get()
	}
	if v.nodes == nil {
		if s := pairingArenaPool.Get(); s != nil {
			v.nodes = s[:0]
		} else {
			v.nodes = make([]pairingNode, 0, defaultPairingArenaCap)
		}
	}
	filled := v.filled.Load()
	fmax := math.Float64frombits(v.maxDist.Load())
	var dirtyMax, dirtyFill bool
	for _, dist := range result {
		if !rankable(dist) {
			continue
		}
		distance := float64(dist.GetDistance())
		if filled && distance >= fmax {
			break
		}
		id := dist.GetId()
		if _, already := v.visited[id]; already {
			continue
		}
		v.visited[id] = struct{}{}
		if v.num > 0 {
			if len(v.nodes) < v.num {
				idx := int32(len(v.nodes)) //nolint:gosec // bounded by v.num
				v.nodes = append(v.nodes, pairingNode{
					dp:      DistPayload{raw: dist, distance: distance},
					child:   pairingNull,
					sibling: pairingNull,
				})
				v.root = mergeMaxNodes(v.nodes, v.root, idx)
				if len(v.nodes) == v.num {
					filled = true
					fmax = v.nodes[v.root].dp.distance
					dirtyMax = true
					dirtyFill = true
				}
				continue
			}
			oldRoot := v.root
			sub := mergePairsMaxNodes(v.nodes, v.nodes[oldRoot].child)
			v.nodes[oldRoot] = pairingNode{
				dp:      DistPayload{raw: dist, distance: distance},
				child:   pairingNull,
				sibling: pairingNull,
			}
			v.root = mergeMaxNodes(v.nodes, sub, oldRoot)
			if newMax := v.nodes[v.root].dp.distance; newMax < fmax {
				fmax = newMax
				dirtyMax = true
			}
			continue
		}
		idx := int32(len(v.nodes)) //nolint:gosec // arena length stays within int32 search result bounds
		v.nodes = append(v.nodes, pairingNode{
			dp:      DistPayload{raw: dist, distance: distance},
			child:   pairingNull,
			sibling: pairingNull,
		})
		v.root = mergeMinNodes(v.nodes, v.root, idx)
	}
	if dirtyMax {
		v.maxDist.Store(math.Float64bits(fmax))
	}
	if dirtyFill {
		v.filled.Store(true)
	}
}

func (v *valdPairingHeapAggr) Result() *payload.Search_Response {
	v.closed.Store(true)
	v.mu.Lock()
	putAggrSeen(v.visited)
	v.visited = nil
	res := v.result
	v.result = nil
	if v.num > 0 {
		k := len(v.nodes)
		if k > 0 {
			if cap(res) >= k {
				res = res[:k]
			} else {
				res = make([]*payload.Object_Distance, k)
			}
			for i := k - 1; i >= 0; i-- {
				res[i] = v.nodes[v.root].dp.raw
				v.root = mergePairsMaxNodes(v.nodes, v.nodes[v.root].child)
			}
		}
	} else {
		k := len(v.nodes)
		res = make([]*payload.Object_Distance, k)
		for i := range k {
			res[i] = v.nodes[v.root].dp.raw
			v.root = mergePairsMinNodes(v.nodes, v.nodes[v.root].child)
		}
	}
	if v.nodes != nil {
		if cap(v.nodes) <= maxAggrPrealloc {
			clear(v.nodes)
			pairingArenaPool.Put(v.nodes[:0])
		}
		v.nodes = nil
	}
	v.root = pairingNull
	v.mu.Unlock()
	return &payload.Search_Response{
		Results: res,
	}
}

func (v *valdPairingHeapAggr) hasResults() bool {
	v.mu.Lock()
	ok := len(v.nodes) > 0
	v.mu.Unlock()
	return ok
}

func (v *valdPairingHeapAggr) uniqueCount() int {
	v.mu.Lock()
	n := len(v.visited)
	v.mu.Unlock()
	return n
}

func (v *valdPairingHeapAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdPairingHeapAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

// plane sort.
type valdSliceAggr struct {
	result []DistPayload
	num    int
	fnum   int
	mu     sync.Mutex
}

func newSlice(num, fnum, replica int) Aggregator {
	return &valdSliceAggr{
		num:    num,
		fnum:   fnum,
		result: make([]DistPayload, 0, preallocLen(num*replica)),
	}
}

func (v *valdSliceAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if v.fnum > 0 && len(result) > v.fnum {
		result = result[:v.fnum]
	}
	if len(result) == 0 || ctx.Err() != nil {
		return
	}
	v.mu.Lock()
	for _, dist := range result {
		if !rankable(dist) {
			continue
		}
		v.result = append(v.result, DistPayload{raw: dist, distance: float64(dist.GetDistance())})
	}
	v.mu.Unlock()
}

func (v *valdSliceAggr) buildResultLocked() *payload.Search_Response {
	v.result = dedupSortByID(v.result)
	top := v.result
	if v.num != 0 && len(top) > v.num {
		top = top[:v.num]
	}
	res := &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(top)),
	}
	for _, r := range top {
		res.Results = append(res.Results, r.raw)
	}
	return res
}

func (v *valdSliceAggr) Result() *payload.Search_Response {
	v.mu.Lock()
	res := v.buildResultLocked()
	v.mu.Unlock()
	return res
}

func (v *valdSliceAggr) hasResults() bool {
	v.mu.Lock()
	ok := len(v.result) > 0
	v.mu.Unlock()
	return ok
}

func (v *valdSliceAggr) uniqueCount() int {
	v.mu.Lock()
	seen := aggrSeenPool.Get()
	for _, p := range v.result {
		seen[p.raw.GetId()] = struct{}{}
	}
	n := len(seen)
	putAggrSeen(seen)
	v.mu.Unlock()
	return n
}

func (v *valdSliceAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdSliceAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

// plane sort.
type valdPoolSliceAggr struct {
	valdSliceAggr
}

var (
	//nolint:gochecknoglobals // poolDist reuses distance slices across search responses
	poolDist = sync.Pool[[]DistPayload]{
		New: func() []DistPayload {
			return make([]DistPayload, 0, poolLen.Load())
		},
	}
	poolLen atomic.Uint64
)

func newPoolSlice(num, fnum, replica int) Aggregator {
	l := uint64(preallocLen(num * replica)) //nolint:gosec // preallocLen clamps to [0, maxAggrPrealloc]
	if poolLen.Load() < l {
		poolLen.Store(l)
	}
	res := poolDist.Get()
	if res != nil {
		res = res[:0]
	} else {
		res = make([]DistPayload, 0, l)
	}
	return &valdPoolSliceAggr{
		result: res,
		num:    num,
		fnum:   fnum,
	}
}

func (v *valdPoolSliceAggr) Result() *payload.Search_Response {
	v.mu.Lock()
	res := v.buildResultLocked()
	// do not let the pool retain the DistPayloads (and the responses they point to).
	if v.result != nil {
		if cap(v.result) <= maxAggrPrealloc {
			clear(v.result)
			poolDist.Put(v.result[:0])
		}
		v.result = nil
	}
	v.mu.Unlock()
	return res
}

// dedupSortByID keeps the first-arrived entry of each ID (as the std and
// pairing heap aggregators do through their visited maps), then stable-sorts
// by ascending distance. It works in place; the vacated tail is cleared so it
// does not keep the dropped entries reachable.
func dedupSortByID(x []DistPayload) []DistPayload {
	if len(x) <= 1 {
		return x
	}
	seen := aggrSeenPool.Get()
	out := x[:0]
	for _, p := range x {
		id := p.raw.GetId()
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			out = append(out, p)
		}
	}
	putAggrSeen(seen)
	clear(x[len(out):])
	slices.SortStableFunc(out, cmpDist)
	return out
}
