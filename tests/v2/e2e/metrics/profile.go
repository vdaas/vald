//go:build e2e

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

package profile

////////////////////////////////////////////////////////////////////////////////
// Struct Section
////////////////////////////////////////////////////////////////////////////////

// //////////////////////////////////////////////////////////////////////////////
// Const Section
// //////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Method Section
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Function Section
////////////////////////////////////////////////////////////////////////////////

// package profile は、高スループットなベンチマーク結果を
// スレッドセーフに集計するためのライブラリです。
//
// ライフサイクルは context.Context で管理されます。
// Summary() で中間スナップショットを取得し、
// Finalize() で最終結果を計算します。
//
// 成功/失敗の判定は、 WorkloadFunc が返す error が
// nil かどうかのみに基づきます。

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// --- Public Structs & Types ---

// MetricsOptions は、ベンチマーク実行時の設定を定義します。
type MetricsOptions struct {
	// Concurrency は、サマリーのメタデータとしてのみ使用されます。
	Concurrency int
	// StoreAllLatencies は、メモリ上にすべてのレイテンシ生データを保持するかどうか。
	// ** true に設定しない限り、Summary の Add() マージは失敗します。 **
	StoreAllLatencies bool
	// TimelineResolution は、タイムラインサンプリングの解像度（バケットサイズ）です。
	TimelineResolution time.Duration
	// TopK は、保持する低速なE2Eリクエストの数です。
	TopK int
	// Estimator は、使用する分位点推定量（Quantile Estimator）の種類です。
	Estimator EstimatorKind
	// TDigestCompression は、t-digest を使用する場合の圧縮レベルです。
	TDigestCompression int
	// Target は、サマリーのメタデータとして保存される識別子です。
	Target string
}

// DefaultMetricsOptions は、MetricsOptions のデフォルト値を返します。
func DefaultMetricsOptions() MetricsOptions {
	return MetricsOptions{
		Concurrency:        50,
		StoreAllLatencies:  true, // Add() のために true をデフォルトに
		TimelineResolution: time.Second,
		TopK:               10,
		Estimator:          EstimatorHistogram,
		TDigestCompression: 200,
		Target:             "unknown",
	}
}

// EstimatorKind は、分位点推定量の種類を定義します。
type EstimatorKind string

const (
	EstimatorHistogram EstimatorKind = "histogram"
	EstimatorP2        EstimatorKind = "p2"
	EstimatorTDigest   EstimatorKind = "tdigest"
)

// requestResult は、1回のリクエスト結果を保持します (private)。
type requestResult struct {
	Status    int
	Err       error
	QueuedAt  time.Time
	StartedAt time.Time
	EndedAt   time.Time
	QueueWait time.Duration
	Service   time.Duration
	E2E       time.Duration
}

// resultPool は requestResult の sync.Pool です。
var resultPool = sync.Pool{
	New: func() interface{} {
		return new(requestResult)
	},
}

// getResult は、sync.Pool から requestResult を取得し、リセットします (private)。
func getResult() *requestResult {
	res := resultPool.Get().(*requestResult)
	*res = requestResult{}
	return res
}

// putResult は、requestResult インスタンスを sync.Pool に返します (private)。
func putResult(res *requestResult) {
	resultPool.Put(res)
}

// WorkloadFunc は、ベンチマークワーカーが実行する作業を定義する関数型です。
type WorkloadFunc func() (status int, err error, dur time.Duration)

// Metrics は、ベンチマークの実行と集計を管理するメイン構造体です。
type Metrics struct {
	opts      MetricsOptions
	startTime time.Time
	agg       *aggregator
	ctx       context.Context    // 親コンテキスト (Run で使用)
	cancel    context.CancelFunc // 内部キャンセル用 (Finalize で使用)
}

// Init は、ベンチマークの実行ステージを初期化します。
// ctx: ベンチマーク全体のキャンセルに使用するコンテキスト。
func Init(ctx context.Context, opts MetricsOptions) *Metrics {
	var msStart runtime.MemStats
	runtime.ReadMemStats(&msStart)
	startTime := time.Now()

	internalCtx, internalCancel := context.WithCancel(context.Background())

	agg := newAggregator(internalCtx, opts, startTime, &msStart)

	// 親コンテキストのキャンセルを監視し、内部コンテキストをキャンセルする
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("profile: Parent context cancelled, stopping aggregator...")
			internalCancel()
			agg.closeStopChOnce()
		case <-internalCtx.Done():
			// 内部的に停止した場合
		}
	}()

	return &Metrics{
		opts:      opts,
		startTime: startTime,
		agg:       agg,
		ctx:       ctx,
		cancel:    internalCancel,
	}
}

// Run は、単一のワークロードを実行し、その結果を集計キューに非同期で送信します。
func (m *Metrics) Run(workload WorkloadFunc) error {
	select {
	case <-m.ctx.Done():
		return m.ctx.Err()
	case <-m.agg.ctx.Done(): // aggregator の内部コンテキストをチェック
		return fmt.Errorf("profile: aggregator already stopped")
	default:
		// 実行継続
	}

	m.agg.recordInflight(1) // Increment inflight

	enqueuedAt := time.Now()
	res := getResult()

	res.StartedAt = time.Now()
	status, err, dur := workload()
	res.EndedAt = time.Now()

	res.Status = status
	res.Err = err
	res.Service = dur
	res.QueuedAt = enqueuedAt
	res.QueueWait = res.StartedAt.Sub(enqueuedAt)
	res.E2E = res.EndedAt.Sub(res.StartedAt)

	// 結果を送信 (送信失敗時はプールに戻す)
	if !m.agg.append(m.ctx, res) {
		putResult(res)
		m.agg.recordInflight(-1) // Decrement if append fails
		select {
		case <-m.ctx.Done():
			return m.ctx.Err()
		default:
			// Check internal context as well
			select {
			case <-m.agg.ctx.Done():
				return fmt.Errorf("profile: aggregator stopped during result submission")
			default:
				// Should not happen if append returns false, but handle defensively
				return fmt.Errorf("profile: failed to append result for unknown reason")
			}
		}
	}
	m.agg.recordInflight(-1) // Decrement inflight after successful append
	return nil
}

// Summary は、現在の集計状態のスナップショット（計算前）を要求し、取得します。
func (m *Metrics) Summary() (*MetricsSummary, error) {
	// コレクターの内部コンテキストをチェック
	select {
	case <-m.agg.ctx.Done():
		return nil, fmt.Errorf("profile: aggregator already stopped or stopping")
	default:
		// Continue
	}

	// スナップショット要求
	select {
	case m.agg.summaryReqCh <- struct{}{}:
		// OK
	case <-m.agg.ctx.Done(): // Check internal context
		return nil, fmt.Errorf("profile: aggregator stopped while requesting summary")
	case <-m.ctx.Done(): // Check parent context
		return nil, fmt.Errorf("profile: parent context cancelled while requesting summary: %w", m.ctx.Err())
	}

	// スナップショット応答待機
	select {
	case snapshot := <-m.agg.summaryRespCh:
		snapshot.Concurrency = m.opts.Concurrency
		snapshot.EstimatorCfg = m.opts.Estimator
		return snapshot, nil
	case <-m.agg.ctx.Done(): // Check internal context
		return nil, fmt.Errorf("profile: aggregator stopped while waiting for summary snapshot")
	case <-m.ctx.Done(): // Check parent context
		return nil, fmt.Errorf("profile: parent context cancelled while waiting for summary snapshot: %w", m.ctx.Err())
	}
}

// Finalize は、アグリゲータに停止をシグナルし、
// すべての処理が完了するまで待機します。
// 最終的な計算を行い、結果のサマリーを返します。
func (m *Metrics) Finalize(
	totalDuration time.Duration, numRequests int, targetQPS int,
) *MetricsSummary {
	// 内部コンテキストをキャンセルしてコレクターに停止を伝える
	m.cancel()

	summary := m.agg.stopAndFinalize(m.startTime)

	summary.calculate(totalDuration, numRequests, targetQPS)

	summary.NumReq = numRequests
	summary.Concurrency = m.opts.Concurrency
	summary.EstimatorCfg = m.opts.Estimator
	summary.TargetQPS = targetQPS

	return summary
}

// StartTime はベンチマークの開始時刻を返します。
func (m *Metrics) StartTime() time.Time {
	return m.startTime
}

// --- Public JSON Export Structs ---

type ErrorDetail struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}
type Percentiles struct {
	Avg float64 `json:"avg_ms"`
	Min float64 `json:"min_ms"`
	P50 float64 `json:"p50_ms"`
	P90 float64 `json:"p90_ms"`
	P95 float64 `json:"p95_ms"`
	P99 float64 `json:"p99_ms"`
	Max float64 `json:"max_ms"`
}
type JSONHistogram struct {
	Buckets []JSONBucket `json:"buckets"`
	SumMS   float64      `json:"sum_ms"`
	Count   int64        `json:"count"`
	MinMS   float64      `json:"min_ms"`
	MaxMS   float64      `json:"max_ms"`
}
type JSONBucket struct {
	Le         string `json:"le"`
	Cumulative int64  `json:"cumulative"`
}
type JSONExport struct {
	Meta struct {
		Target      string        `json:"target"`
		TotalTimeMS float64       `json:"total_time_ms"`
		StartedAt   time.Time     `json:"started_at"`
		Resolution  time.Duration `json:"timeline_resolution"`
		Estimator   string        `json:"estimator"`
		Requests    int           `json:"requests"`
		Concurrency int           `json:"concurrency"`
		TargetQPS   int           `json:"target_qps"`
	} `json:"meta"`
	Summary struct {
		TotalRequests int64   `json:"total_requests"`
		QPS           float64 `json:"qps"`
		ErrorRate     float64 `json:"error_rate"`
		MaxInflight   int64   `json:"max_inflight"`
		GCDelta       uint32  `json:"gc_cycles"`
		AllocDiff     int64   `json:"alloc_diff_bytes"`
		SysDiff       int64   `json:"sys_diff_bytes"`
	} `json:"summary"`
	E2E     Percentiles `json:"e2e"`
	Service Percentiles `json:"service"`
	Queue   struct {
		AvgMS float64 `json:"avg_ms"`
		P90MS float64 `json:"p90_ms"`
	} `json:"queue"`
	StatusCounts        map[string]int           `json:"status_counts"`
	LatencyByStatus     map[string]Percentiles   `json:"latency_by_status"`
	LatencyByStatusSvc  map[string]Percentiles   `json:"service_by_status"`
	ServiceHistByStatus map[string]JSONHistogram `json:"service_hist_by_status"`
	Timeline            struct {
		Requests  []int `json:"requests"`
		Errors    []int `json:"errors"`
		Goroutine []int `json:"goroutine"`
	} `json:"timeline"`
	Errors  []ErrorDetail `json:"errors"`
	TopSlow []struct {
		E2EMS  float64 `json:"e2e_ms"`
		Status string  `json:"status"`
		Error  string  `json:"error,omitempty"`
	} `json:"top_slow"`
}

// MetricsSummary (public struct)
type MetricsSummary struct {
	opts MetricsOptions

	allResults []*requestResult

	StatusCounts     map[int]int
	AggregatedErrors []ErrorDetail
	errorMsgIndex    map[string]int
	successfulReqs   atomic.Int64
	failedReqs       atomic.Int64
	// MODIFIED: currentInflightInternal 追加
	currentInflightInternal atomic.Int64 // Internal counter for maxInflight CAS

	e2eEst          quantileEstimator
	svcEst          quantileEstimator
	byStatus        map[int]quantileEstimator
	byStatusSvc     map[int]quantileEstimator
	svcHistByStatus map[int]*histogram

	baseTime          time.Time
	res               time.Duration
	TimelineRequests  []int
	TimelineErrors    []int
	TimelineGoroutine []int

	topK int
	slow slowMinHeap

	maxInflight atomic.Int64

	TotalTime     time.Duration
	TotalRequests int64
	QPS           float64
	ErrorRate     float64
	QueueAvg      time.Duration
	QueueP90      time.Duration
	SvcAvg        time.Duration
	SvcP90        time.Duration

	msStart   runtime.MemStats
	GCDelta   uint32
	AllocDiff int64
	SysDiff   int64

	Target       string
	NumReq       int
	Concurrency  int
	TargetQPS    int
	EstimatorCfg EstimatorKind
}

// Add (public method)
func (s *MetricsSummary) Add(other *MetricsSummary) error {
	if !s.opts.StoreAllLatencies || !other.opts.StoreAllLatencies {
		return fmt.Errorf("profile: cannot merge summaries, StoreAllLatencies must be true for both")
	}

	s.allResults = append(s.allResults, other.allResults...)

	s.successfulReqs.Add(other.successfulReqs.Load())
	s.failedReqs.Add(other.failedReqs.Load())

	for code, count := range other.StatusCounts {
		s.StatusCounts[code] += count
	}

	for _, otherErr := range other.AggregatedErrors {
		if idx, ok := s.errorMsgIndex[otherErr.Message]; ok {
			s.AggregatedErrors[idx].Count += otherErr.Count
		} else {
			errCopy := otherErr
			s.AggregatedErrors = append(s.AggregatedErrors, errCopy)
			s.errorMsgIndex[otherErr.Message] = len(s.AggregatedErrors) - 1
		}
	}

	if other.baseTime.Before(s.baseTime) {
		s.baseTime = other.baseTime
	}
	for {
		curMax := s.maxInflight.Load()
		otherMax := other.maxInflight.Load()
		if otherMax > curMax {
			if s.maxInflight.CompareAndSwap(curMax, otherMax) {
				break
			}
		} else {
			break
		}
	}

	s.QPS = 0
	s.ErrorRate = 0
	s.TotalTime = 0
	s.TotalRequests = int64(len(s.allResults))
	s.e2eEst = nil
	s.svcEst = nil

	return nil
}

// calculate (private method)
func (s *MetricsSummary) calculate(totalDuration time.Duration, numRequests int, targetQPS int) {
	s.TotalTime = totalDuration
	s.NumReq = numRequests
	s.TargetQPS = targetQPS

	actualRequests := s.successfulReqs.Load() + s.failedReqs.Load()
	s.TotalRequests = actualRequests

	if !s.opts.StoreAllLatencies {
		// Estimators were updated in real-time
	}

	if s.opts.StoreAllLatencies {
		s.recalculateFromRaw()
		// recalculateFromRaw will update s.TotalRequests based on len(s.allResults)
	}

	calculatedQPS := 0.0
	if s.TotalTime > 0 && s.TotalRequests > 0 {
		calculatedQPS = float64(s.TotalRequests) / s.TotalTime.Seconds()
	}

	if targetQPS > 0 {
		s.QPS = float64(targetQPS)
	} else {
		s.QPS = calculatedQPS
	}

	denominator := float64(numRequests)
	if denominator <= 0 {
		denominator = float64(s.TotalRequests)
	}
	if denominator > 0 {
		s.ErrorRate = float64(s.failedReqs.Load()) / denominator
	} else {
		s.ErrorRate = 0
	}
}

// ExportJSON (public method)
func (s *MetricsSummary) ExportJSON() ([]byte, error) {
	out := JSONExport{}

	out.Meta.Target = s.Target
	out.Meta.TotalTimeMS = durMS(s.TotalTime)
	out.Meta.StartedAt = s.baseTime
	out.Meta.Resolution = s.res
	out.Meta.Estimator = string(s.EstimatorCfg)
	out.Meta.Requests = s.NumReq
	out.Meta.Concurrency = s.Concurrency
	out.Meta.TargetQPS = s.TargetQPS

	out.Summary.TotalRequests = s.TotalRequests
	out.Summary.QPS = s.QPS
	out.Summary.ErrorRate = s.ErrorRate
	out.Summary.MaxInflight = s.maxInflight.Load()
	out.Summary.GCDelta = s.GCDelta
	out.Summary.AllocDiff = s.AllocDiff
	out.Summary.SysDiff = s.SysDiff

	out.E2E = percentilesFrom(s.e2eEst)
	out.Service = percentilesFrom(s.svcEst)
	out.Queue.AvgMS = durMS(s.QueueAvg)
	out.Queue.P90MS = durMS(s.QueueP90)

	out.StatusCounts = make(map[string]int, len(s.StatusCounts))
	for k, v := range s.StatusCounts {
		key := strconv.Itoa(k)
		out.StatusCounts[key] = v
	}

	out.LatencyByStatus = make(map[string]Percentiles)
	out.LatencyByStatusSvc = make(map[string]Percentiles)
	for code, est := range s.byStatus {
		key := strconv.Itoa(code)
		out.LatencyByStatus[key] = percentilesFrom(est)
	}
	for code, est := range s.byStatusSvc {
		key := strconv.Itoa(code)
		out.LatencyByStatusSvc[key] = percentilesFrom(est)
	}

	out.ServiceHistByStatus = make(map[string]JSONHistogram)
	for code, h := range s.svcHistByStatus {
		key := strconv.Itoa(code)
		out.ServiceHistByStatus[key] = promHistFrom(h)
	}

	out.Timeline.Requests = append(out.Timeline.Requests, s.TimelineRequests...)
	out.Timeline.Errors = append(out.Timeline.Errors, s.TimelineErrors...)
	out.Timeline.Goroutine = append(out.Timeline.Goroutine, s.TimelineGoroutine...)
	out.Errors = append(out.Errors, s.AggregatedErrors...)

	if s.topK > 0 && len(s.slow) > 0 {
		tmp := make([]*slowItem, len(s.slow))
		copy(tmp, s.slow)
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].e2e > tmp[j].e2e })
		for _, it := range tmp {
			status := strconv.Itoa(it.status)
			out.TopSlow = append(out.TopSlow, struct {
				E2EMS  float64 `json:"e2e_ms"`
				Status string  `json:"status"`
				Error  string  `json:"error,omitempty"`
			}{
				E2EMS:  durMS(it.e2e),
				Status: status,
				Error:  it.errMsg,
			})
		}
	}
	return json.Marshal(out)
}

// PrintJSON (public method)
func (s *MetricsSummary) PrintJSON() {
	b, err := s.ExportJSON()
	if err != nil {
		fmt.Println(`{"error":"failed to export json","detail":"` + strings.ReplaceAll(err.Error(), `"`, `\"`) + `"}`)
		return
	}
	fmt.Println(string(b))
}

// --- Private Implementation ---

// aggregator は (private) 集計ロジックをカプセル化します。
type aggregator struct {
	summary       *MetricsSummary
	results       chan *requestResult
	gorSample     chan goroutineSample
	summaryReqCh  chan struct{}
	summaryRespCh chan *MetricsSummary

	wg       sync.WaitGroup
	stopOnce sync.Once
	stopCh   chan struct{}   // Finalize からの明示的な停止用
	ctx      context.Context // 外部キャンセル監視用
}

type goroutineSample struct {
	n int
	t time.Time
}

// newAggregator (private)
func newAggregator(
	ctx context.Context, opts MetricsOptions, startTime time.Time, msStart *runtime.MemStats,
) *aggregator {
	s := newMetricsSummary(opts, startTime, msStart)

	bufferSize := opts.Concurrency * 4
	if bufferSize < 1024 {
		bufferSize = 1024
	}

	a := &aggregator{
		summary:       s,
		results:       make(chan *requestResult, bufferSize),
		gorSample:     make(chan goroutineSample, 128),
		summaryReqCh:  make(chan struct{}),
		summaryRespCh: make(chan *MetricsSummary, 1),
		stopCh:        make(chan struct{}),
		ctx:           ctx, // Store context passed from Init
	}

	a.wg.Add(1)
	go a.collector()

	return a
}

// collector (private)
func (a *aggregator) collector() {
	defer a.wg.Done()

	resCh := a.results
	reqCh := a.summaryReqCh

	// gorSampleCh は現在未使用

	for {
		select {
		case res, ok := <-resCh:
			if !ok {
				resCh = nil
			} else {
				// Safety: Assign to temporary variable
				tempRes := res
				a.summary.append(tempRes, a.summary.opts.StoreAllLatencies)
				putResult(tempRes)
			}

		case <-reqCh: // Summary snapshot requested
			snapshot := a.summary.deepCopy()
			// Send the snapshot back (non-blocking due to buffered channel)
			// Need to handle potential block if RespCh is full or closed
			select {
			case a.summaryRespCh <- snapshot:
				// Successfully sent snapshot
			case <-a.stopCh:
				// Aggregator stopped while trying to send snapshot
			case <-a.ctx.Done():
				// Aggregator cancelled while trying to send snapshot
			}

		case <-a.stopCh: // Explicit stop signal from Finalize()
			fmt.Println("profile: collector received stop signal, draining...")
			goto DRAIN

		case <-a.ctx.Done(): // Context cancelled externally
			fmt.Println("profile: collector cancelled via context, draining...")
			a.closeStopChOnce() // Ensure stopCh is also closed
			goto DRAIN
		}

		// Exit condition (channels closed unexpectedly)
		if resCh == nil && reqCh == nil {
			fmt.Println("profile: collector channels closed unexpectedly.")
			a.closeStopChOnce()
			goto DRAIN
		}
	}

DRAIN:
	// Drain logic (after stopCh or ctx.Done())
	// Close reqCh safely (might already be closed by stopCh path)
	func() {
		defer func() { recover() }() // Ignore potential panic on closing closed channel
		close(reqCh)
	}()

	// Drain results channel
	if resCh != nil {
		// Non-blocking drain first
		for len(resCh) > 0 {
			res := <-resCh
			tempRes := res
			a.summary.append(tempRes, a.summary.opts.StoreAllLatencies)
			putResult(tempRes)
		}
		// Close and blocking drain for concurrent sends
		close(resCh)
		for res := range resCh {
			tempRes := res
			a.summary.append(tempRes, a.summary.opts.StoreAllLatencies)
			putResult(tempRes)
		}
	}

	fmt.Println("profile: collector finished draining.")
}

// append (private)
// ** MODIFIED: Checks ctx and stopCh **
func (a *aggregator) append(ctx context.Context, res *requestResult) bool {
	select {
	case a.results <- res:
		return true
	case <-a.stopCh: // Finalize が呼ばれた
		return false
	case <-ctx.Done(): // 親コンテキストがキャンセルされた
		return false
	case <-a.ctx.Done(): // 内部コンテキストもチェック (冗長かもしれないが安全のため)
		return false
	}
}

// recordInflight (private)
func (a *aggregator) recordInflight(delta int64) {
	newVal := a.summary.currentInflightInternal.Add(delta)
	if delta > 0 { // Increment
		for {
			curMax := a.summary.maxInflight.Load()
			if newVal > curMax {
				if a.summary.maxInflight.CompareAndSwap(curMax, newVal) {
					break
				}
			} else {
				break
			}
		}
	}
}

// stopAndFinalize (private)
func (a *aggregator) stopAndFinalize(startTime time.Time) *MetricsSummary {
	a.closeStopChOnce() // stopCh を閉じる

	a.wg.Wait() // コレクターの終了を待つ

	s := a.summary

	var msEnd runtime.MemStats
	runtime.ReadMemStats(&msEnd)
	if msEnd.NumGC >= s.msStart.NumGC {
		s.GCDelta = msEnd.NumGC - s.msStart.NumGC
	}
	s.AllocDiff = int64(msEnd.Alloc) - int64(s.msStart.Alloc)
	s.SysDiff = int64(msEnd.Sys) - int64(s.msStart.Sys)

	return s
}

// closeStopChOnce は stopCh を安全に一度だけ閉じます。
func (a *aggregator) closeStopChOnce() {
	a.stopOnce.Do(func() {
		close(a.stopCh)
	})
}

// newMetricsSummary (private)
func newMetricsSummary(
	opts MetricsOptions, start time.Time, msStart *runtime.MemStats,
) *MetricsSummary {
	s := &MetricsSummary{
		opts:             opts,
		StatusCounts:     make(map[int]int),
		AggregatedErrors: make([]ErrorDetail, 0),
		errorMsgIndex:    make(map[string]int),

		byStatus:        make(map[int]quantileEstimator),
		byStatusSvc:     make(map[int]quantileEstimator),
		svcHistByStatus: make(map[int]*histogram),

		baseTime: start,
		res:      opts.TimelineResolution,
		topK:     opts.TopK,
		Target:   opts.Target,
		msStart:  *msStart,
	}

	if opts.StoreAllLatencies {
		s.allResults = make([]*requestResult, 0)
	} else {
		s.e2eEst = newEstimator(opts.Estimator, opts.TDigestCompression)
		s.svcEst = newEstimator(opts.Estimator, opts.TDigestCompression)
	}
	return s
}

// deepCopy (private method on summary)
func (s *MetricsSummary) deepCopy() *MetricsSummary {
	c := new(MetricsSummary)
	*c = *s // Shallow copy first

	if s.opts.StoreAllLatencies && s.allResults != nil {
		c.allResults = make([]*requestResult, len(s.allResults))
		for i, r := range s.allResults {
			rc := *r
			c.allResults[i] = &rc
		}
	}

	c.StatusCounts = make(map[int]int, len(s.StatusCounts))
	for k, v := range s.StatusCounts {
		c.StatusCounts[k] = v
	}

	c.AggregatedErrors = make([]ErrorDetail, len(s.AggregatedErrors))
	copy(c.AggregatedErrors, s.AggregatedErrors)

	c.errorMsgIndex = make(map[string]int, len(s.errorMsgIndex))
	for k, v := range s.errorMsgIndex {
		c.errorMsgIndex[k] = v
	}

	// Omit estimator copies for snapshot
	c.e2eEst = nil
	c.svcEst = nil
	c.byStatus = make(map[int]quantileEstimator)
	c.byStatusSvc = make(map[int]quantileEstimator)
	c.svcHistByStatus = make(map[int]*histogram)

	c.TimelineRequests = append([]int(nil), s.TimelineRequests...)
	c.TimelineErrors = append([]int(nil), s.TimelineErrors...)
	c.TimelineGoroutine = append([]int(nil), s.TimelineGoroutine...)

	c.slow = make(slowMinHeap, len(s.slow))
	copy(c.slow, s.slow)

	c.successfulReqs.Store(s.successfulReqs.Load())
	c.failedReqs.Store(s.failedReqs.Load())
	c.maxInflight.Store(s.maxInflight.Load()) // Store current max snapshot

	// Reset calculated fields
	c.TotalTime = 0
	c.TotalRequests = 0
	c.QPS = 0
	c.ErrorRate = 0

	return c
}

// append (private method on summary)
func (s *MetricsSummary) append(res *requestResult, keepRaw bool) {
	s.StatusCounts[res.Status]++

	if res.Err != nil {
		s.failedReqs.Add(1)
		msg := res.Err.Error()
		if idx, ok := s.errorMsgIndex[msg]; ok {
			s.AggregatedErrors[idx].Count++
		} else {
			s.AggregatedErrors = append(s.AggregatedErrors, ErrorDetail{Message: msg, Count: 1})
			s.errorMsgIndex[msg] = len(s.AggregatedErrors) - 1
		}
	} else {
		s.successfulReqs.Add(1)
	}

	if keepRaw {
		resCopy := *res
		s.allResults = append(s.allResults, &resCopy)
	}

	if !keepRaw {
		s.e2eEst.Observe(res.E2E)
		s.svcEst.Observe(res.Service)
		getOrCreateEstimator(s.byStatus, res.Status, s.opts.Estimator, s.opts.TDigestCompression).Observe(res.E2E)
		getOrCreateEstimator(s.byStatusSvc, res.Status, s.opts.Estimator, s.opts.TDigestCompression).Observe(res.Service)
		getOrCreateHist(s.svcHistByStatus, res.Status).Observe(res.Service)
	}

	idx, rq, er, gr := ensureTimelineSlot(s.baseTime, s.res, s.TimelineRequests, s.TimelineErrors, s.TimelineGoroutine, res.EndedAt)
	rq[idx]++

	if res.Err != nil {
		er[idx]++
	}
	s.TimelineRequests, s.TimelineErrors, s.TimelineGoroutine = rq, er, gr
	s.recordTopK(res.E2E, res.Status, res.Err)
}

// recalculateFromRaw (private method on summary)
func (s *MetricsSummary) recalculateFromRaw() {
	if !s.opts.StoreAllLatencies {
		return
	}

	s.TotalRequests = int64(len(s.allResults))

	s.e2eEst = newEstimator(s.opts.Estimator, s.opts.TDigestCompression)
	s.svcEst = newEstimator(s.opts.Estimator, s.opts.TDigestCompression)
	s.byStatus = make(map[int]quantileEstimator)
	s.byStatusSvc = make(map[int]quantileEstimator)
	s.svcHistByStatus = make(map[int]*histogram)

	queueWaits := make([]time.Duration, len(s.allResults))

	for i, res := range s.allResults {
		s.e2eEst.Observe(res.E2E)
		s.svcEst.Observe(res.Service)

		getOrCreateEstimator(s.byStatus, res.Status, s.opts.Estimator, s.opts.TDigestCompression).Observe(res.E2E)
		getOrCreateEstimator(s.byStatusSvc, res.Status, s.opts.Estimator, s.opts.TDigestCompression).Observe(res.Service)
		getOrCreateHist(s.svcHistByStatus, res.Status).Observe(res.Service)

		queueWaits[i] = res.QueueWait
	}

	if len(queueWaits) > 0 {
		sort.Slice(queueWaits, func(i, j int) bool { return queueWaits[i] < queueWaits[j] })
		s.QueueAvg = avg(queueWaits)
		s.QueueP90 = pctFromSorted(queueWaits, 90)
	}

	if s.svcEst != nil && s.svcEst.Total() > 0 {
		s.SvcAvg = s.svcEst.Avg()
		s.SvcP90 = s.svcEst.Percentile(90)
	}
}

// setGoroutineSample (private method on summary)
func (s *MetricsSummary) setGoroutineSample(now time.Time, n int) {
	// (未使用)
}

// recordTopK (private method on summary)
func (s *MetricsSummary) recordTopK(e2e time.Duration, status int, err error) {
	if s.topK <= 0 {
		return
	}
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	if len(s.slow) < s.topK {
		heap.Push(&s.slow, &slowItem{e2e: e2e, status: status, errMsg: msg})
		return
	}
	if s.slow[0].e2e < e2e {
		heap.Pop(&s.slow)
		heap.Push(&s.slow, &slowItem{e2e: e2e, status: status, errMsg: msg})
	}
}

// --- All internal estimators and helpers ---

// quantileEstimator (private interface)
type quantileEstimator interface {
	Observe(d time.Duration)
	Percentile(p float64) time.Duration
	Avg() time.Duration
	Min() time.Duration
	Max() time.Duration
	Total() int64
}

// baseEstimator (private struct)
type baseEstimator struct {
	minV  float64
	maxV  float64
	sum   float64
	total int64
}

func newBaseEstimator() baseEstimator {
	return baseEstimator{minV: math.Inf(1), maxV: math.Inf(-1)}
}

func (b *baseEstimator) ObserveBase(x float64) {
	if x < b.minV {
		b.minV = x
	}
	if x > b.maxV {
		b.maxV = x
	}
	b.sum += x
	b.total++
}

func (b *baseEstimator) Avg() time.Duration {
	if b.total == 0 {
		return 0
	}
	return time.Duration(b.sum / float64(b.total))
}

func (b *baseEstimator) Min() time.Duration {
	if b.total == 0 {
		return 0
	}
	return time.Duration(b.minV)
}

func (b *baseEstimator) Max() time.Duration { return time.Duration(b.maxV) }
func (b *baseEstimator) Total() int64       { return b.total }

// histogram (private struct)
var defaultBuckets = []time.Duration{
	250 * time.Microsecond,
	500 * time.Microsecond,
	1 * time.Millisecond, 2 * time.Millisecond, 4 * time.Millisecond, 8 * time.Millisecond,
	16 * time.Millisecond, 32 * time.Millisecond, 64 * time.Millisecond, 128 * time.Millisecond,
	256 * time.Millisecond, 512 * time.Millisecond, 1 * time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second,
}

type histogram struct {
	Buckets []time.Duration
	Counts  []int64
	TotalN  int64
	Sum     time.Duration
	MinV    time.Duration
	MaxV    time.Duration
}

func newHistogram(b []time.Duration) *histogram {
	return &histogram{
		Buckets: append([]time.Duration(nil), b...),
		Counts:  make([]int64, len(b)+1),
		MinV:    time.Duration(math.MaxInt64),
	}
}

func (h *histogram) Observe(d time.Duration) {
	if d < h.MinV {
		h.MinV = d
	}
	if d > h.MaxV {
		h.MaxV = d
	}
	h.TotalN++
	h.Sum += d
	i := sort.Search(len(h.Buckets), func(i int) bool { return d <= h.Buckets[i] })
	h.Counts[i]++
}

func (h *histogram) Avg() time.Duration {
	if h.TotalN == 0 {
		return 0
	}
	return time.Duration(int64(h.Sum) / h.TotalN)
}

func (h *histogram) Percentile(p float64) time.Duration {
	if h.TotalN == 0 {
		return 0
	}
	target := int64(math.Ceil((p / 100.0) * float64(h.TotalN)))
	var cum int64
	for i, c := range h.Counts {
		cum += c
		if cum >= target {
			var low time.Duration
			if i == 0 {
				low = 0
			} else {
				low = h.Buckets[i-1]
			}
			var high time.Duration
			if i < len(h.Buckets) {
				high = h.Buckets[i]
			} else {
				high = h.MaxV
				if high < low {
					high = low
				}
			}
			return low + (high-low)/2
		}
	}
	return h.MaxV
}

// histEstimator (private struct)
type histEstimator struct{ h *histogram }

func newHistEstimator() *histEstimator                       { return &histEstimator{h: newHistogram(defaultBuckets)} }
func (he *histEstimator) Observe(d time.Duration)            { he.h.Observe(d) }
func (he *histEstimator) Percentile(p float64) time.Duration { return he.h.Percentile(p) }
func (he *histEstimator) Avg() time.Duration                 { return he.h.Avg() }
func (he *histEstimator) Min() time.Duration {
	if he.h.TotalN == 0 {
		return 0
	}
	return he.h.MinV
}
func (he *histEstimator) Max() time.Duration { return he.h.MaxV }
func (he *histEstimator) Total() int64       { return he.h.TotalN }

// p2Quant (private struct)
type p2Quant struct {
	p     float64
	count int
	init  []float64
	q     [5]float64
	n     [5]int
	ns    [5]float64
	dn    [5]float64
}

func newP2Quant(p float64) *p2Quant {
	return &p2Quant{p: p, init: make([]float64, 0, 5)}
}

func (pq *p2Quant) Observe(x float64) {
	pq.count++
	if len(pq.init) < 5 {
		pq.init = append(pq.init, x)
		if len(pq.init) == 5 {
			sort.Float64s(pq.init)
			for i := 0; i < 5; i++ {
				pq.q[i] = pq.init[i]
				pq.n[i] = i + 1
			}
			p := []float64{0, pq.p / 2, pq.p, (1 + pq.p) / 2, 1}
			for i := 0; i < 5; i++ {
				pq.ns[i] = 1 + p[i]*float64(pq.count-1)
			}
			pq.dn = [5]float64{0, pq.p / 2, pq.p, (1 + pq.p) / 2, 1}
		}
		return
	}
	var k int
	if x < pq.q[0] {
		pq.q[0] = x
		k = 0
	} else if x >= pq.q[4] {
		pq.q[4] = x
		k = 3
	} else {
		for i := 0; i < 4; i++ {
			if pq.q[i] <= x && x < pq.q[i+1] {
				k = i
				break
			}
		}
	}
	for i := k + 1; i < 5; i++ {
		pq.n[i]++
	}
	for i := 0; i < 5; i++ {
		pq.ns[i] += pq.dn[i]
	}
	for i := 1; i <= 3; i++ {
		d := pq.ns[i] - float64(pq.n[i])
		if (d >= 1 && pq.n[i+1]-pq.n[i] > 1) || (d <= -1 && pq.n[i]-pq.n[i-1] > 1) {
			s := 1
			if d < 0 {
				s = -1
			}
			qip := pq.q[i] + float64(s)*parabolic(pq.q[i-1], pq.q[i], pq.q[i+1], pq.n[i-1], pq.n[i], pq.n[i+1])
			if pq.q[i-1] < qip && qip < pq.q[i+1] {
				pq.q[i] = qip
			} else {
				pq.q[i] = pq.q[i] + float64(s)*(pq.q[i+s]-pq.q[i])/float64(pq.n[i+s]-pq.n[i])
			}
			pq.n[i] += s
		}
	}
}

func parabolic(qm1, q0, q1 float64, nm1, n0, n1 int) float64 {
	den1 := float64(n1 - n0)
	den2 := float64(n0 - nm1)
	if den1 == 0 || den2 == 0 || n1 == nm1 {
		return 0
	}
	a := float64(n0-nm1) * (q1 - q0) / den1
	b := float64(n1-n0) * (q0 - qm1) / den2
	return (a + b) / float64(n1-nm1)
}

func (pq *p2Quant) Value() float64 {
	if pq.count == 0 {
		return 0
	}
	if len(pq.init) < 5 {
		cp := append([]float64(nil), pq.init...)
		sort.Float64s(cp)
		rank := int(math.Ceil(pq.p * float64(len(cp))))
		if rank < 1 {
			rank = 1
		}
		if rank > len(cp) {
			rank = len(cp)
		}
		return cp[rank-1]
	}
	return pq.q[2]
}

// p2Estimator (private struct)
type p2Estimator struct {
	baseEstimator
	ps []float64
	qs []*p2Quant
}

func newP2Estimator() *p2Estimator {
	ps := []float64{0.50, 0.90, 0.95, 0.99}
	qs := make([]*p2Quant, 0, len(ps))
	for _, p := range ps {
		qs = append(qs, newP2Quant(p))
	}
	return &p2Estimator{
		baseEstimator: newBaseEstimator(),
		ps:            ps,
		qs:            qs,
	}
}

func (pe *p2Estimator) Observe(d time.Duration) {
	x := float64(d)
	pe.ObserveBase(x)
	for _, q := range pe.qs {
		q.Observe(x)
	}
}

func (pe *p2Estimator) Percentile(p float64) time.Duration {
	if pe.total == 0 {
		return 0
	}
	type kv struct{ p, v float64 }
	arr := make([]kv, 0, len(pe.qs))
	for i, q := range pe.qs {
		arr = append(arr, kv{p: pe.ps[i], v: q.Value()})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].p < arr[j].p })
	if p <= arr[0].p {
		lo := arr[0]
		alpha := p / (lo.p + 1e-9)
		minVal := float64(pe.Min())
		return time.Duration(minVal + alpha*(lo.v-minVal))
	}
	if p >= arr[len(arr)-1].p {
		hi := arr[len(arr)-1]
		alpha := (p - hi.p) / (1 - hi.p + 1e-9)
		return time.Duration(hi.v + alpha*(float64(pe.Max())-hi.v))
	}
	for i := 0; i < len(arr)-1; i++ {
		if arr[i].p <= p && p <= arr[i+1].p {
			r := (p - arr[i].p) / (arr[i+1].p - arr[i].p + 1e-9)
			return time.Duration(arr[i].v + r*(arr[i+1].v-arr[i].v))
		}
	}
	return time.Duration(arr[len(arr)-1].v)
}

// tDigestEstimator (private struct)
type (
	centroid         struct{ mean, weight float64 }
	tDigestEstimator struct {
		baseEstimator
		comp      int
		cents     []centroid
		needsSort bool
	}
)

func newTDigestEstimator(compression int) *tDigestEstimator {
	if compression <= 0 {
		compression = 200
	}
	return &tDigestEstimator{
		baseEstimator: newBaseEstimator(),
		comp:          compression,
	}
}

func (td *tDigestEstimator) Observe(d time.Duration) {
	x := float64(d)
	td.ObserveBase(x)
	td.cents = append(td.cents, centroid{mean: x, weight: 1})
	td.needsSort = true
	if len(td.cents) > td.comp*3/2 {
		td.compress()
	}
}

func (td *tDigestEstimator) Percentile(p float64) time.Duration {
	if td.total == 0 {
		return 0
	}
	td.ensureSorted()
	target := p / 100.0 * float64(td.total)
	var cum float64
	prev := centroid{}
	for i, c := range td.cents {
		if i == 0 {
			cum += c.weight / 2
			if target <= cum {
				return time.Duration(c.mean)
			}
			prev = c
			continue
		}
		cum += (prev.weight/2 + c.weight/2)
		if target <= cum {
			span := prev.weight/2 + c.weight/2
			alpha := (target - (cum - span)) / (span + 1e-9)
			return time.Duration(prev.mean + alpha*(c.mean-prev.mean))
		}
		prev = c
	}
	return time.Duration(td.maxV)
}

func (td *tDigestEstimator) ensureSorted() {
	if td.needsSort {
		sort.Slice(td.cents, func(i, j int) bool { return td.cents[i].mean < td.cents[j].mean })
		td.needsSort = false
	}
}

func (td *tDigestEstimator) compress() {
	td.ensureSorted()
	if len(td.cents) <= td.comp {
		return
	}
	newC := make([]centroid, 0, td.comp)
	cur := td.cents[0]
	for i := 1; i < len(td.cents); i++ {
		next := td.cents[i]
		if len(newC)+len(td.cents)-i+1 > td.comp {
			mergedW := cur.weight + next.weight
			cur.mean = (cur.mean*cur.weight + next.mean*next.weight) / mergedW
			cur.weight = mergedW
		} else {
			newC = append(newC, cur)
			cur = next
		}
	}
	newC = append(newC, cur)
	td.cents = newC
	td.needsSort = false
}

// newEstimator (private func)
func newEstimator(kind EstimatorKind, tdComp int) quantileEstimator {
	switch kind {
	case EstimatorP2:
		return newP2Estimator()
	case EstimatorTDigest:
		return newTDigestEstimator(tdComp)
	default:
		return newHistEstimator()
	}
}

// percentilesFrom (private func)
func percentilesFrom(est quantileEstimator) Percentiles {
	if est == nil || est.Total() == 0 {
		return Percentiles{}
	}
	return Percentiles{
		Avg: durMS(est.Avg()),
		Min: durMS(est.Min()),
		P50: durMS(est.Percentile(50)),
		P90: durMS(est.Percentile(90)),
		P95: durMS(est.Percentile(95)),
		P99: durMS(est.Percentile(99)),
		Max: durMS(est.Max()),
	}
}

// keysSorted (private func)
func keysSorted(m map[int]int) []int {
	ks := make([]int, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Ints(ks)
	return ks
}

// ensureTimelineSlot (private func)
func ensureTimelineSlot(
	base time.Time, res time.Duration, reqs, errs, g []int, t time.Time,
) (int, []int, []int, []int) {
	idx := int(t.Sub(base) / res)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(reqs) {
		n := idx + 1
		reqs = append(reqs, make([]int, n-len(reqs))...)
		errs = append(errs, make([]int, n-len(errs))...)
		g = append(g, make([]int, n-len(g))...)
	}
	return idx, reqs, errs, g
}

// getOrCreateEstimator (private func)
func getOrCreateEstimator(
	m map[int]quantileEstimator, code int, kind EstimatorKind, tdComp int,
) quantileEstimator {
	est, ok := m[code]
	if !ok {
		est = newEstimator(kind, tdComp)
		m[code] = est
	}
	return est
}

// getOrCreateHist (private func)
func getOrCreateHist(m map[int]*histogram, code int) *histogram {
	h, ok := m[code]
	if !ok {
		h = newHistogram(defaultBuckets)
		m[code] = h
	}
	return h
}

// slowItem, slowMinHeap (private types)
type slowItem struct {
	e2e    time.Duration
	status int
	errMsg string
}
type slowMinHeap []*slowItem

func (h slowMinHeap) Len() int            { return len(h) }
func (h slowMinHeap) Less(i, j int) bool  { return h[i].e2e < h[j].e2e }
func (h slowMinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *slowMinHeap) Push(x interface{}) { *h = append(*h, x.(*slowItem)) }
func (h *slowMinHeap) Pop() interface{} {
	old := *h
	it := old[len(old)-1]
	*h = old[:len(old)-1]
	return it
}

// promHistFrom (private func)
func promHistFrom(h *histogram) JSONHistogram {
	if h == nil || h.TotalN == 0 {
		return JSONHistogram{}
	}
	bk := make([]JSONBucket, 0, len(h.Buckets)+1)
	var cum int64
	for i, c := range h.Counts {
		cum += c
		var le string
		if i < len(h.Buckets) {
			le = h.Buckets[i].String()
		} else {
			le = "+Inf"
		}
		bk = append(bk, JSONBucket{Le: le, Cumulative: cum})
	}
	return JSONHistogram{
		Buckets: bk,
		SumMS:   durMS(h.Sum),
		Count:   h.TotalN,
		MinMS:   durMS(h.MinV),
		MaxMS:   durMS(h.MaxV),
	}
}

// durMS (private func)
func durMS(d time.Duration) float64 { return float64(d) / 1e6 }

// avg (private func)
func avg(xs []time.Duration) time.Duration {
	if len(xs) == 0 {
		return 0
	}
	var sum int64
	for _, d := range xs {
		sum += int64(d)
	}
	return time.Duration(sum / int64(len(xs)))
}

// pctFromSorted (private func)
func pctFromSorted(xs []time.Duration, p int) time.Duration {
	if len(xs) == 0 {
		return 0
	}
	idx := int(math.Ceil(float64(p)/100*float64(len(xs)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(xs) {
		idx = len(xs) - 1
	}
	return xs[idx]
}
