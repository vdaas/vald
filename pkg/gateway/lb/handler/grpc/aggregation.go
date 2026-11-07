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
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type Aggregator interface {
	Start(ctx context.Context)
	Send(ctx context.Context, data *payload.Search_Response)
	GetNum() int  // get top-k number
	GetFnum() int // get forwarding top-k number calculated by search ratio
	Result() *payload.Search_Response
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
func cmpDist(l, r *DistPayload) int {
	return cmp.Compare(l.distance, r.distance)
}

func (s *server) aggregationSearch(
	ctx context.Context,
	rpc string, // agent RPC name f calls, used as the metrics label
	agents int, // number of agents the broadcast should reach
	aggr Aggregator,
	bcfg *payload.Search_Config, // Base Config of Request
	f func(ctx context.Context,
		fcfg *payload.Search_Config, // Forwarding Config to Agent
		vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (res *payload.Search_Response, attrs []attribute.KeyValue, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "aggregationSearch"), apiName+"/aggregationSearch")
	defer trace.End(span)

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

	fcfg := bcfg.CloneVT() // Forwarding Config to Agent, this config need to modify like below so it should be cloned
	fcfg.Num = uint32(aggr.GetFnum())
	fcfg.MinNum = 0

	// failed holds the code of every target whose agent call failed (see
	// searchAgentFailure), whether or not that failure fails the BroadCast. Like
	// service.Gateway.BroadCast, only the last outcome per target is kept because
	// the client backoff may re-invoke the callback for a target, so a target is
	// counted at most once per request. reached holds every target the callback
	// ran for.
	var (
		fmu     sync.Mutex
		failed  = make(map[string]codes.Code)
		reached = make(map[string]struct{})
	)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	aggr.Start(ctx)
	// bctx scopes only the broadcast, so the quorum early termination below stops
	// the remaining agent calls while the aggregator keeps every result it got.
	bctx, bcancel := context.WithCancelCause(ctx)
	defer bcancel(nil)
	quorum := s.searchQuorum(agents)
	var (
		answered atomic.Int64
		qtimer   atomic.Pointer[time.Timer]
	)
	defer func() {
		if t := qtimer.Load(); t != nil {
			t.Stop()
		}
	}()
	// Per-agent failures are tolerated and summarized once per request below, so
	// the client must not log every failed agent/attempt at Warn/Error.
	// Replicas cover a failed agent, so an overloaded one (ResourceExhausted) is not
	// retried inside the search deadline (write paths keep backoff-and-wait).
	err = s.gateway.BroadCast(grpc.WithQuietErrorLog(grpc.WithoutResourceExhaustedRetry(bctx)), service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		var agentErr error // error of the last agent call, nil when it answered
		defer func() {
			fmu.Lock()
			defer fmu.Unlock()
			reached[target] = struct{}{}
			if code, ok := searchAgentFailure(target, agentErr); ok {
				failed[target] = code
				return
			}
			delete(failed, target)
		}()
		sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/aggregationSearch/"+target)
		defer trace.End(sspan)
		r, err := f(sctx, fcfg, vc, copts...)
		if err != nil {
			agentErr = err
			return handleSearchBroadCastError(sspan, target, err)
		}
		if r == nil || len(r.GetResults()) == 0 {
			if sctx.Err() == nil {
				r, err = f(sctx, fcfg, vc, copts...)
				if err != nil {
					agentErr = err
					return handleSearchBroadCastError(sspan, target, err)
				}
			}
			if r == nil || len(r.GetResults()) == 0 {
				err = status.WrapWithNotFound(fmt.Sprintf("failed to process search request from %s", target),
					errors.ErrEmptySearchResult,
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					})
				errhandler.RecordSpanError(sspan, codes.NotFound, err)
				return nil
			}
		}
		aggr.Send(sctx, r)
		// Count an agent only after all its results were sent, so terminating
		// early never drops results of an agent that counts toward the quorum.
		if quorum > 0 && int(answered.Add(1)) == quorum {
			qtimer.Store(time.AfterFunc(s.searchQuorumGrace, func() {
				bcancel(errSearchQuorumReached)
			}))
		}
		return nil
	})
	quorumReached := errors.Is(context.Cause(bctx), errSearchQuorumReached)
	if quorumReached && errors.Is(err, context.Canceled) {
		// the agents still running were canceled on purpose, not failed
		err = nil
	}
	ctxErr := ctx.Err()
	cancel()
	fmu.Lock()
	// The client never runs the callback for a target whose circuit breaker is
	// open or that has no connection, so such agents are neither in failed nor
	// in reached. They are counted as unreached agents with code Unavailable.
	// Once the request ctx has ended, targets not reached yet are deadline or
	// cancel misses, which are not agent failures, so none are counted then.
	var unreached int
	if ctxErr == nil && !quorumReached {
		unreached = max(agents-len(reached), 0)
	}
	failures := len(failed) + unreached
	for _, code := range failed {
		s.searchMetrics.RecordAgentFailure(rpc, code)
	}
	fmu.Unlock()
	for range unreached {
		s.searchMetrics.RecordAgentFailure(rpc, codes.Unavailable)
	}
	var bcErr error // broadcast error of a request that still has results
	// The per-agent errors are logged at Debug; this one line per request
	// decides whether the agent failures deserve a Warn.
	defer func() {
		logPartialSearch(rpc, failures, s.replica, bcErr, err)
	}()
	if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		err = status.WrapWithInternal("search API connection not found", err,
			s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg)...)
		attrs = trace.StatusCodeInternal(err.Error())
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, attrs, err
	}
	res = aggr.Result()
	if num != 0 && len(res.GetResults()) > num {
		res.Results = res.GetResults()[:num]
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		if len(res.GetResults()) == 0 {
			err = status.WrapWithDeadlineExceeded(
				"error search result length is 0 due to the timeoutage limit",
				errors.ErrEmptySearchResult,
				append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			errhandler.RecordSpanAttrs(span, attrs, err)
			return nil, attrs, err
		}
		if 0 < minNum && len(res.GetResults()) < minNum {
			err = status.WrapWithDeadlineExceeded(
				fmt.Sprintf("error search result length is not enough due to the timeoutage limit, required: %d, found: %d", minNum, len(res.GetResults())),
				errors.ErrInsuffcientSearchResult,
				append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			errhandler.RecordSpanAttrs(span, attrs, err)
			return nil, attrs, err
		}
	}

	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse search gRPC error response",
			append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...)
		attrs = trace.FromGRPCStatus(st.Code(), msg)
		errhandler.RecordSpanAttrs(span, attrs, err)
		if len(res.GetResults()) == 0 {
			// with agent failures the deferred summary logs err at Warn
			if !searchSummaryLogsRequestErr(failures, err) {
				log.Warn(err)
			}
			return nil, attrs, err
		}
		// some agents failed but the others returned results: the severity is
		// decided below once the request is known to succeed.
		bcErr = err
	}
	if num != 0 && len(res.GetResults()) == 0 {
		if err == nil {
			err = errors.ErrEmptySearchResult
		}
		err = status.WrapWithNotFound("error search result length is 0", err,
			append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...)
		attrs = trace.StatusCodeNotFound(err.Error())
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, attrs, err
	}

	if 0 < minNum && len(res.GetResults()) < minNum {
		err = status.WrapWithNotFound(
			fmt.Sprintf("error search result length is not enough required: %d, found: %d", minNum, len(res.GetResults())),
			errors.ErrInsuffcientSearchResult,
			append(s.errDetails(ctx, "search", bcfg.GetRequestId(), bcfg), info.Get())...,
		)
		// with agent failures the deferred summary logs bcErr with err at Warn
		if bcErr != nil && !searchSummaryLogsRequestErr(failures, err) {
			log.Warn(bcErr)
		}
		attrs = trace.StatusCodeNotFound(err.Error())
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, attrs, err
	}
	if failures > 0 {
		s.searchMetrics.RecordPartialSuccess(rpc, searchRecallAtRisk(failures, s.replica))
	}
	res.RequestId = bcfg.GetRequestId()
	return res, attrs, nil
}

// searchRecallAtRisk reports whether failures agents failing may have made some
// vectors unreachable. Each vector is stored on replica agents, so while fewer
// than replica agents fail every vector still has a live replica and the partial
// result keeps the recall of a full broadcast.
func searchRecallAtRisk(failures, replica int) bool {
	return failures >= replica
}

// searchAgentFailure reports whether an agent's search error counts as an agent
// failure, and with which code. Every error is a failure except the
// Canceled/DeadlineExceeded/NotFound misses, which are expected when the request
// deadline cuts the broadcast short or an agent has no neighbor, and
// InvalidArgument, which the request caused (e.g. a vector dimension the index
// does not have) on otherwise healthy agents. A non-status
// error is counted as Unknown. Whether the error fails the BroadCast is decided
// separately by handleSearchBroadCastError.
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

func searchAgentFailure(target string, err error) (codes.Code, bool) {
	if err == nil ||
		errors.IsAny(err, context.Canceled, context.DeadlineExceeded) ||
		errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)) ||
		errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)) {
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
// the same request again: exactly one Warn line is logged per failed request.
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
	// maxDist holds the math.Float64bits of the current distance bound. Only
	// the Start goroutine stores it; Send reads it concurrently.
	maxDist atomic.Uint64
	eg      errgroup.Group
	dch     chan DistPayload
	cancel  context.CancelFunc
	result  []*payload.Object_Distance
	wg      sync.WaitGroup
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
	case errors.Is(err, context.Canceled),
		errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
		if sspan != nil {
			sspan.RecordError(err)
			sspan.SetAttributes(trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix +
					"/vald.v1.search.BroadCast/" +
					target + " canceled: " + err.Error())...)
			sspan.SetStatus(trace.StatusError, err.Error())
		}
		return nil
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
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
		case codes.NotFound:
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
// It sets the top-k size (`num`), the forwarding top-k size (`fnum`), and allocates internal buffers sized proportionally to `replica`
// (capped at maxAggrPrealloc).
// The returned Aggregator is initialized (including its error group) and ready to Start and receive results.
func newStd(num, fnum, replica int) Aggregator {
	hint := preallocLen(num * replica)
	vsa := &valdStdAggr{
		num:    num,
		fnum:   fnum,
		dch:    make(chan DistPayload, hint),
		result: make([]*payload.Object_Distance, 0, hint),
		eg:     errgroup.Get(),
	}
	vsa.maxDist.Store(math.Float64bits(math.MaxFloat64))
	vsa.closed.Store(false)
	return vsa
}

func (v *valdStdAggr) Start(ctx context.Context) {
	ctx, v.cancel = context.WithCancel(ctx)
	// seen de-duplicates result IDs (replicas return the same IDs). It is only
	// touched by add, which runs solely on the goroutine started below, so a
	// plain map suffices. The first-dequeued entry of an ID wins, whether or not
	// it is kept.
	seen := make(map[string]struct{}, cap(v.result))
	add := func(distance float64, dist *payload.Object_Distance) {
		if _, ok := seen[dist.GetId()]; ok {
			return
		}
		seen[dist.GetId()] = struct{}{}
		rl := len(v.result) // result length
		fmax := math.Float64frombits(v.maxDist.Load())
		if rl >= v.num && distance >= fmax {
			return
		}
		// stable insertion point: after every element whose distance <= distance.
		i := rl
		for i > 0 && distance < float64(v.result[i-1].GetDistance()) {
			i--
		}
		v.result = slices.Insert(v.result, i, dist)
		rl = len(v.result)
		if rl > v.num && v.num != 0 {
			v.result = v.result[:v.num]
			rl = len(v.result)
		}
		if distEnd := float64(v.result[rl-1].GetDistance()); rl >= v.num &&
			distEnd < fmax {
			v.maxDist.Store(math.Float64bits(distEnd))
		}
	}
	v.wg.Add(1)
	v.eg.Go(func() error {
		defer v.wg.Done()
		for {
			select {
			case <-ctx.Done():
				// Do NOT close(v.dch): concurrent Send goroutines fanned out by
				// BroadCast may be between their v.closed check and their send,
				// and a send on a closed channel panics. That panic is a
				// runtime.Error which safety.RecoverFunc re-panics, crashing the
				// gateway. Instead flag closed (future Sends return early) and
				// drain what is already buffered without blocking; in-flight
				// Sends observe ctx.Done() via their own select and stop.
				v.closed.Store(true)
				for {
					select {
					case dist := <-v.dch:
						add(dist.distance, dist.raw)
					default:
						return nil
					}
				}
			case dist := <-v.dch:
				add(dist.distance, dist.raw)
			}
		}
	})
}

func (v *valdStdAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.fnum {
		result = result[:v.fnum]
	}
	for _, dist := range result {
		if rankable(dist) {
			fdist := float64(dist.GetDistance())
			if fdist >= math.Float64frombits(v.maxDist.Load()) || v.closed.Load() {
				return
			}
			select {
			case <-ctx.Done():
				return
			case v.dch <- DistPayload{raw: dist, distance: fdist}:
			}
		}
	}
}

func (v *valdStdAggr) Result() *payload.Search_Response {
	v.closed.Store(true)
	v.cancel()
	v.wg.Wait()
	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	return &payload.Search_Response{
		Results: v.result,
	}
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
	ph      *PairingHeap
	visited map[string]struct{} // guarded by mu
	result  []*payload.Object_Distance
	num     int
	fnum    int
	mu      sync.Mutex
}

func newPairingHeap(num, fnum, _ int) Aggregator {
	return &valdPairingHeapAggr{
		num:     num,
		fnum:    fnum,
		ph:      new(PairingHeap),
		visited: make(map[string]struct{}, preallocLen(num)),
		result:  make([]*payload.Object_Distance, 0, preallocLen(num)),
	}
}

func (*valdPairingHeapAggr) Start(_ context.Context) {}

func (v *valdPairingHeapAggr) Send(ctx context.Context, data *payload.Search_Response) {
	dps := collectDistPayloads(ctx, data, v.fnum)
	if len(dps) == 0 {
		return
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	for i := range dps {
		id := dps[i].raw.GetId()
		if _, already := v.visited[id]; already {
			continue
		}
		v.visited[id] = struct{}{}
		v.ph = v.ph.Insert(&dps[i])
	}
}

func (v *valdPairingHeapAggr) Result() *payload.Search_Response {
	for !v.ph.IsEmpty() && len(v.result) <= v.num {
		var minDist *DistPayload
		minDist, v.ph = v.ph.ExtractMin()
		if minDist != nil {
			v.result = append(v.result, minDist.raw)
		} else if v.ph == nil {
			break
		}
	}
	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	return &payload.Search_Response{
		Results: v.result,
	}
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
	result []*DistPayload
	num    int
	fnum   int
	mu     sync.Mutex
}

func newSlice(num, fnum, replica int) Aggregator {
	return &valdSliceAggr{
		num:    num,
		fnum:   fnum,
		result: make([]*DistPayload, 0, preallocLen(num*replica)),
	}
}

func (*valdSliceAggr) Start(_ context.Context) {}

func (v *valdSliceAggr) Send(ctx context.Context, data *payload.Search_Response) {
	dps := collectDistPayloads(ctx, data, v.fnum)
	if len(dps) == 0 {
		return
	}
	v.mu.Lock()
	for i := range dps {
		v.result = append(v.result, &dps[i])
	}
	v.mu.Unlock()
}

func (v *valdSliceAggr) Result() (res *payload.Search_Response) {
	v.result = dedupSortByID(v.result)

	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	res = &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(v.result)),
	}
	for _, r := range v.result {
		res.Results = append(res.GetResults(), r.raw)
	}
	return res
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
	result []*DistPayload
	num    int
	fnum   int
	mu     sync.Mutex
}

var (
	poolDist = sync.Pool{
		New: func() any {
			return make([]*DistPayload, 0, poolLen.Load())
		},
	}
	poolLen atomic.Uint64
)

func newPoolSlice(num, fnum, replica int) Aggregator {
	l := uint64(preallocLen(num * replica)) //nolint:gosec // preallocLen clamps to [0, maxAggrPrealloc]
	if poolLen.Load() < l {
		poolLen.Store(l)
	}
	return &valdPoolSliceAggr{
		num:    num,
		fnum:   fnum,
		result: poolDist.Get().([]*DistPayload),
	}
}

func (*valdPoolSliceAggr) Start(_ context.Context) {}

func (v *valdPoolSliceAggr) Send(ctx context.Context, data *payload.Search_Response) {
	dps := collectDistPayloads(ctx, data, v.fnum)
	if len(dps) == 0 {
		return
	}
	v.mu.Lock()
	for i := range dps {
		v.result = append(v.result, &dps[i])
	}
	v.mu.Unlock()
}

func (v *valdPoolSliceAggr) Result() (res *payload.Search_Response) {
	v.result = dedupSortByID(v.result)

	top := v.result
	if len(top) > v.num {
		top = top[:v.num]
	}
	res = &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(top)),
	}
	for _, r := range top {
		res.Results = append(res.GetResults(), r.raw)
	}
	// do not let the pool retain the DistPayloads (and the responses they point to).
	clear(v.result)
	// skipcq: SCC-SA6002
	poolDist.Put(v.result[:0])
	return res
}

func (v *valdPoolSliceAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdPoolSliceAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

// collectDistPayloads returns the rankable entries among the first fnum results
// of data as one slab, stopping at the first entry seen after ctx is done. The
// slice and pairing heap aggregators take pointers into it, so each agent
// response costs one allocation and one lock instead of one per result.
func collectDistPayloads(
	ctx context.Context, data *payload.Search_Response, fnum int,
) []DistPayload {
	result := data.GetResults()
	if len(result) > fnum {
		result = result[:fnum]
	}
	dps := make([]DistPayload, 0, len(result))
	for _, dist := range result {
		if !rankable(dist) {
			continue
		}
		select {
		case <-ctx.Done():
			return dps
		default:
		}
		dps = append(dps, DistPayload{raw: dist, distance: float64(dist.GetDistance())})
	}
	return dps
}

// dedupSortByID keeps the first-arrived entry of each ID (as the std and
// pairing heap aggregators do through their visited maps), then stable-sorts
// by ascending distance. It works in place; the vacated tail is cleared so it
// does not keep the dropped entries reachable.
func dedupSortByID(x []*DistPayload) []*DistPayload {
	seen := make(map[string]struct{}, len(x))
	out := x[:0]
	for _, p := range x {
		if _, ok := seen[p.raw.GetId()]; !ok {
			seen[p.raw.GetId()] = struct{}{}
			out = append(out, p)
		}
	}
	clear(x[len(out):])
	slices.SortStableFunc(out, cmpDist)
	return out
}
