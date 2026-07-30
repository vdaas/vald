//go:build e2e

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

// Package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"testing"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
)

const (
	// recallSumCounterName and recallSamplesCounterName are the
	// metrics.Metrics.CustomCounters names used to persist recall@k into an
	// execution's metrics.Collector. recall is a float64 in [0, 1], but
	// Collector custom counters only support monotonically increasing
	// int64/uint64 deltas (see metrics.CounterHandle.Add), so each recall
	// sample is scaled by recallScale (parts-per-million precision) before
	// being accumulated; recallSamplesCounterName tracks how many search
	// responses contributed to the sum so the mean can be recovered later
	// (see meanRecall).
	//
	// These names are opt-in: they only take effect for an
	// Execution/Operation/Strategy/Data whose metrics.custom_counters
	// includes them (see tests/v2/e2e/assets/agent_recall_qps.yaml).
	recallSumCounterName     = "recall_sum_ppm"
	recallSamplesCounterName = "recall_samples"
	recallScale              = 1_000_000.0
)

// recordRecall persists a single recall@k sample (as computed by
// calculateRecall/metrics.CalcRecall) into plan's own metrics.Collector, if
// metrics are enabled for plan. It is a no-op otherwise, mirroring how
// single() in grpc_test.go already guards latency recording on
// plan.Metrics.Enabled.
func recordRecall(plan *config.Execution, rc float64) {
	if plan == nil || plan.Metrics == nil || !plan.Metrics.Enabled || plan.Collector == nil {
		return
	}
	plan.Collector.IncCounter(recallSumCounterName, int64(rc*recallScale))
	plan.Collector.IncCounter(recallSamplesCounterName, 1)
}

// meanRecall recovers the mean recall@k recorded into col via recordRecall.
// ok is false if col is nil, the recall counters were never registered
// (metrics.custom_counters does not list them for this Collector), or no
// recall samples were recorded (e.g. col belongs to a non-search execution).
func meanRecall(col metrics.Collector) (mean float64, samples uint64, ok bool) {
	if col == nil {
		return 0, 0, false
	}
	sumH, err := col.CounterHandle(recallSumCounterName)
	if err != nil || sumH == nil {
		return 0, 0, false
	}
	samplesH, err := col.CounterHandle(recallSamplesCounterName)
	if err != nil || samplesH == nil {
		return 0, 0, false
	}
	samples = samplesH.Value()
	if samples == 0 {
		return 0, 0, false
	}
	mean = float64(sumH.Value()) / recallScale / float64(samples)
	return mean, samples, true
}

// qpsFromSnapshot derives search (or any other operation's) throughput from
// a metrics.GlobalSnapshot's request count and observed time span, without
// needing any dedicated QPS bookkeeping: every request already updates
// snap.Total/StartTime/LastUpdated via metrics.Collector.Record.
func qpsFromSnapshot(snap *metrics.GlobalSnapshot) (qps float64, ok bool) {
	if snap == nil || snap.Total == 0 {
		return 0, false
	}
	dur := snap.LastUpdated.Sub(snap.StartTime)
	if dur <= 0 {
		return 0, false
	}
	return float64(snap.Total) / dur.Seconds(), true
}

// logRecallAndQPS reports recall@k and QPS for col, using whichever of the
// two is available (e.g. an Insert execution has QPS but no recall samples,
// while a Search execution configured with the recall custom counters has
// both). It is safe to call for any execution's Collector, including ones
// where neither metric applies (in which case it logs nothing).
func logRecallAndQPS(t testing.TB, label string, col metrics.Collector) {
	t.Helper()
	if col == nil {
		return
	}
	snap := col.GlobalSnapshot()
	qps, qpsOK := qpsFromSnapshot(snap)
	mean, samples, recallOK := meanRecall(col)
	// In benchmark mode (BenchmarkE2EStrategy), surface the scenario metrics
	// on the sub-benchmark's own result line so they show up next to ns/op
	// in benchstat-compatible output. Note that this is called at the
	// operation/strategy grouping levels as well as at execution leaves:
	// grouping-level lines report collector aggregates merged from their
	// children (and their ns/op is total child wall time), so compare
	// benchstat lines only within the same tree depth.
	if qpsOK {
		test.ReportMetric(t, qps, "qps")
	}
	if recallOK {
		test.ReportMetric(t, mean, "recall@k")
	}
	switch {
	case recallOK && qpsOK:
		log.Infof("%s: recall@k=%.4f (samples=%d), QPS=%.2f (requests=%d)", label, mean, samples, qps, snap.Total)
	case recallOK:
		log.Infof("%s: recall@k=%.4f (samples=%d)", label, mean, samples)
	case qpsOK:
		log.Infof("%s: QPS=%.2f (requests=%d)", label, qps, snap.Total)
	}
}
