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

package metrics

import (
	"fmt"
	"html"
	"math"
	"slices"
	"unicode/utf8"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
)

type (
	flattenedExecution struct {
		Strategy   string
		Operation  string
		Label      string
		Execution  ExecutionReport
		Occurrence int
	}

	executionKey struct {
		Strategy  string
		Operation string
		Name      string
		Type      string
		Mode      string
		Index     int
	}

	executionIndex struct {
		exact    map[executionKey]int
		fallback map[executionKey]int
		ordered  []flattenedExecution
		used     []bool
	}

	boundedMarkdownBuilder struct {
		header    string
		out       strings.Builder
		maxBytes  int
		truncated bool
	}
)

var (
	//nolint:gochecknoglobals // precompiled single-pass replacer for Markdown table cells.
	tableCellReplacer = strings.NewReplacer("|", "\\|", "\n", " ", "\r", "")
	//nolint:gochecknoglobals // precompiled single-pass replacer for Markdown text cells.
	markdownCellReplacer = strings.NewReplacer("|", "\\|", "\n", " ", "\r", "", "`", "\\`")
	//nolint:gochecknoglobals // precompiled single-pass newline stripper for HTML summary text.
	htmlSummaryReplacer = strings.NewReplacer("\n", " ", "\r", "")
)

const (
	// BenchmarkDiffMarker is the HTML comment marker placed on line 1 of Benchmark Diff PR comments.
	BenchmarkDiffMarker = "<!-- e2e-v2-benchmark-diff -->"
	// DefaultMaxCommentBytes caps generated Markdown below GitHub's 65,536-character PR comment limit.
	DefaultMaxCommentBytes = 60 * 1024

	naValue            = "N/A"
	summaryScope       = "summary"
	truncationNotice   = "\n> ... truncated due to GitHub comment size limit; see workflow artifacts for the full report.\n"
	executionCollision = -1
	percentFactor      = 100.0
	zeroPercentEpsilon = 0.05
	nsPerMicrosecond   = 1e3
	msPerSecond        = 1e3
	nsPerMillisecond   = 1e6
	nsPerSecond        = 1e9
)

// GenerateBenchmarkDiff loads PR, Daily latest, and Daily nightly reports and formats a 3-way Markdown diff.
func GenerateBenchmarkDiff(prDir, latestDir, nightlyDir string, maxBytes int) (string, error) {
	if strings.TrimSpace(prDir) == "" {
		return "", errors.New("pr-dir is required")
	}
	pr, err := LoadReportsFromDir(prDir)
	if err != nil {
		return "", err
	}
	latest, err := loadOptionalReportsDir(latestDir)
	if err != nil {
		return "", err
	}
	nightly, err := loadOptionalReportsDir(nightlyDir)
	if err != nil {
		return "", err
	}
	return FormatBenchmarkDiff(pr, latest, nightly, maxBytes), nil
}

// FormatBenchmarkDiff formats a 3-way PR vs Daily latest & nightly Benchmark Diff Markdown comment.
func FormatBenchmarkDiff(pr, latest, nightly map[string]*Report, maxBytes int) string {
	scenarios := collectScenarioNames(pr, latest, nightly)
	var header strings.Builder
	fmt.Fprint(&header, BenchmarkDiffMarker, "\n\n## E2E V2 Benchmark Diff\n\n")
	if len(scenarios) == 0 {
		fmt.Fprint(&header, "_No benchmark reports found._\n")
		return newBoundedMarkdownBuilder(header.String(), maxBytes).finish()
	}

	fmt.Fprint(&header, "### Scenario Summary\n\n",
		"| Scenario | PR QPS | vs latest | vs nightly | PR Duration | vs latest | vs nightly | PR Latency (Min / Mean / Max / P99) |\n",
		"|---|---:|---:|---:|---:|---:|---:|---|\n")

	mb := newBoundedMarkdownBuilder(header.String(), maxBytes)
	for _, sc := range scenarios {
		if !mb.append(formatScenarioSummaryRow(sc, pr[sc], latest[sc], nightly[sc])) {
			return mb.finish()
		}
	}
	for _, sc := range scenarios {
		if !mb.append(formatScenarioDetails(sc, pr[sc], latest[sc], nightly[sc])) {
			return mb.finish()
		}
	}
	return mb.finish()
}

func loadOptionalReportsDir(dir string) (map[string]*Report, error) {
	if strings.TrimSpace(dir) == "" {
		return make(map[string]*Report), nil
	}
	fi, exist, err := statPath(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat baseline directory %s", dir)
	}
	if !exist {
		return make(map[string]*Report), nil
	}
	if !fi.IsDir() {
		return nil, errors.Errorf("baseline path %s is not a directory", dir)
	}
	return LoadReportsFromDir(dir)
}

func collectScenarioNames(maps ...map[string]*Report) []string {
	seen := make(map[string]struct{})
	for _, m := range maps {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	names := make([]string, 0, len(seen))
	for k := range seen {
		names = append(names, k)
	}
	slices.Sort(names)
	return names
}

func summarySnap(r *Report) *GlobalSnapshot {
	if r == nil {
		return nil
	}
	return r.Summary
}

func formatScenarioSummaryRow(scenario string, prRep, latestRep, nightlyRep *Report) string {
	prS, latS, nigS := summarySnap(prRep), summarySnap(latestRep), summarySnap(nightlyRep)
	prQPS, prDur, prLat := formatSnapshotSummary(prS)
	vsLatQ, vsLatD := formatQPSDurDelta(prS, latS)
	vsNigQ, vsNigD := formatQPSDurDelta(prS, nigS)
	return fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
		formatInlineCode(scenario), prQPS, vsLatQ, vsNigQ, prDur, vsLatD, vsNigD, prLat)
}

func formatScenarioDetails(scenario string, prRep, latestRep, nightlyRep *Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n<details>\n<summary><code>%s</code> — Summary &amp; Per-RPC Breakdown</summary>\n\n", escapeHTMLSummary(scenario))
	writeScenarioAggregateTable(&b, summarySnap(prRep), summarySnap(latestRep), summarySnap(nightlyRep))
	writeScenarioExecutionTable(&b, prRep, latestRep, nightlyRep)
	fmt.Fprint(&b, "</details>\n")
	return b.String()
}

func writeScenarioAggregateTable(b *strings.Builder, prS, latS, nigS *GlobalSnapshot) {
	fmt.Fprint(b,
		"| Metric | PR | Daily `latest` | vs `latest` | Daily `nightly` | vs `nightly` |\n",
		"|---|---:|---:|---:|---:|---:|\n",
	)

	rows := []struct {
		get   func(*GlobalSnapshot) (float64, bool)
		fmtV  func(float64) string
		label string
	}{
		{(*GlobalSnapshot).qps, formatQPSValue, "QPS"},
		{(*GlobalSnapshot).durationSec, formatDurationValue, "Duration"},
		{(*GlobalSnapshot).latencyMinNS, formatLatencyNS, "Latency Min"},
		{(*GlobalSnapshot).latencyMeanNS, formatLatencyNS, "Latency Mean"},
		{(*GlobalSnapshot).latencyMaxNS, formatLatencyNS, "Latency Max"},
		{(*GlobalSnapshot).latencyP99NS, formatLatencyNS, "Latency P99"},
	}
	for _, r := range rows {
		prStr, latStr, dLat, nigStr, dNig := naValue, naValue, naValue, naValue, naValue
		latV, latOK := r.get(latS)
		if latOK {
			latStr = r.fmtV(latV)
		}
		nigV, nigOK := r.get(nigS)
		if nigOK {
			nigStr = r.fmtV(nigV)
		}
		if prV, prOK := r.get(prS); prOK {
			prStr = r.fmtV(prV)
			if latOK {
				dLat = formatDeltaPercent(prV, latV)
			}
			if nigOK {
				dNig = formatDeltaPercent(prV, nigV)
			}
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s |\n", r.label, prStr, latStr, dLat, nigStr, dNig)
	}
	fmt.Fprint(b, "\n")
}

func writeScenarioExecutionTable(b *strings.Builder, prRep, latestRep, nightlyRep *Report) {
	prIdx := buildExecutionIndex(prRep)
	latestIdx := buildExecutionIndex(latestRep)
	nightlyIdx := buildExecutionIndex(nightlyRep)

	primaryOrdered := prIdx.ordered
	if len(primaryOrdered) == 0 {
		primaryOrdered = latestIdx.ordered
	}
	if len(primaryOrdered) == 0 {
		primaryOrdered = nightlyIdx.ordered
	}
	if len(primaryOrdered) == 0 {
		fmt.Fprint(b, "_No per-RPC executions recorded._\n")
		return
	}

	fmt.Fprint(
		b,
		"| Execution | Type | Mode | Total (Err) | PR QPS | vs latest | vs nightly | PR Duration | vs latest | vs nightly | PR Latency (Min / Mean / Max / P99) | vs latest (Mean / P99) | vs nightly (Mean / P99) |\n",
		"|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|---|---|\n",
	)
	for _, fe := range primaryOrdered {
		var s, latS, nigS *GlobalSnapshot
		if prEx, ok := prIdx.match(fe); ok {
			s = prEx.Summary
		}
		if latEx, ok := latestIdx.match(fe); ok {
			latS = latEx.Summary
		}
		if nigEx, ok := nightlyIdx.match(fe); ok {
			nigS = nigEx.Summary
		}
		totalStr := naValue
		if s != nil {
			totalStr = fmt.Sprintf("%d (%d)", s.Total, s.Errors)
		}
		prQPS, prDur, prLat := formatSnapshotSummary(s)
		qLat, dLat := formatQPSDurDelta(s, latS)
		qNig, dNig := formatQPSDurDelta(s, nigS)
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			escapeMarkdownCell(fe.Label), formatInlineCode(fe.Execution.Type), formatInlineCode(fe.Execution.Mode),
			totalStr, prQPS, qLat, qNig, prDur, dLat, dNig, prLat,
			formatLatencyMeanP99Delta(s, latS), formatLatencyMeanP99Delta(s, nigS))
	}
}

func formatOptionalValue(v float64, ok bool, fn func(float64) string) string {
	if !ok {
		return naValue
	}
	return fn(v)
}

func formatOptionalDelta(curr float64, currOK bool, base float64, baseOK bool) string {
	if !currOK || !baseOK {
		return naValue
	}
	return formatDeltaPercent(curr, base)
}

func formatSnapshotSummary(s *GlobalSnapshot) (qps, dur, lat string) {
	if s == nil {
		return naValue, naValue, naValue
	}
	qV, qOK := s.qps()
	dV, dOK := s.durationSec()
	return formatOptionalValue(qV, qOK, formatQPSValue),
		formatOptionalValue(dV, dOK, formatDurationValue),
		formatLatencyTuple(s)
}

func formatQPSDurDelta(curr, base *GlobalSnapshot) (qpsDelta, durDelta string) {
	if curr == nil || base == nil {
		return naValue, naValue
	}
	cQ, cQOK := curr.qps()
	bQ, bQOK := base.qps()
	cD, cDOK := curr.durationSec()
	bD, bDOK := base.durationSec()
	return formatOptionalDelta(cQ, cQOK, bQ, bQOK),
		formatOptionalDelta(cD, cDOK, bD, bDOK)
}

// buildExecutionIndex flattens scenario executions into lookup tables for 1-to-1 matching.
// When multiple executions share the same (Strategy, Operation, Name, Type, Mode, Index) tuple,
// exact marks the key as colliding so match falls back to per-tuple occurrence order.
func buildExecutionIndex(report *Report) executionIndex {
	if report == nil {
		return executionIndex{}
	}
	idx := executionIndex{
		exact:    make(map[executionKey]int),
		fallback: make(map[executionKey]int),
	}
	counts := make(map[executionKey]int)
	for _, st := range report.Strategies {
		for _, op := range st.Operations {
			for _, ex := range op.Executions {
				baseKey := executionKey{Strategy: st.Name, Operation: op.Name, Name: ex.Name, Type: ex.Type, Mode: ex.Mode}
				occ := counts[baseKey]
				counts[baseKey] = occ + 1

				pos := len(idx.ordered)
				idx.ordered = append(idx.ordered, flattenedExecution{
					Strategy: st.Name, Operation: op.Name, Label: formatExecutionLabel(st.Name, op.Name, ex),
					Execution: ex, Occurrence: occ,
				})

				ek := baseKey
				ek.Index = ex.Index
				if _, exists := idx.exact[ek]; exists {
					idx.exact[ek] = executionCollision
				} else {
					idx.exact[ek] = pos
				}
				fk := baseKey
				fk.Index = occ
				idx.fallback[fk] = pos
			}
		}
	}
	idx.used = make([]bool, len(idx.ordered))
	return idx
}

// match finds an unused execution matching fe, trying its exact Index first and falling back
// to Occurrence order only if exact match is absent, colliding, or already consumed (preventing double-pairing).
func (idx *executionIndex) match(fe flattenedExecution) (ExecutionReport, bool) {
	if idx == nil || len(idx.ordered) == 0 {
		return ExecutionReport{}, false
	}
	ek := executionKey{
		Strategy: fe.Strategy, Operation: fe.Operation,
		Name: fe.Execution.Name, Type: fe.Execution.Type, Mode: fe.Execution.Mode, Index: fe.Execution.Index,
	}
	if pos, ok := idx.exact[ek]; ok && pos >= 0 && !idx.used[pos] {
		idx.used[pos] = true
		return idx.ordered[pos].Execution, true
	}
	ek.Index = fe.Occurrence
	if pos, ok := idx.fallback[ek]; ok && pos >= 0 && !idx.used[pos] {
		idx.used[pos] = true
		return idx.ordered[pos].Execution, true
	}
	return ExecutionReport{}, false
}

func formatExecutionLabel(strategy, operation string, ex ExecutionReport) string {
	execName := ex.Name
	if execName == "" {
		execName = ex.Type
	}
	if operation == "" {
		return fmt.Sprintf("%s / [%d] %s", strategy, ex.Index, execName)
	}
	return fmt.Sprintf("%s / %s / [%d] %s", strategy, operation, ex.Index, execName)
}

func formatQPSValue(v float64) string {
	return fmt.Sprintf("%.2f", finiteOrZero(v))
}

func formatDurationValue(sec float64) string {
	sec = finiteOrZero(sec)
	if sec <= 0 {
		return "0s"
	}
	if sec < 1 {
		return fmt.Sprintf("%.1fms", sec*msPerSecond)
	}
	return fmt.Sprintf("%.2fs", sec)
}

func formatLatencyTuple(s *GlobalSnapshot) string {
	minV, minOK := s.latencyMinNS()
	meanV, meanOK := s.latencyMeanNS()
	maxV, maxOK := s.latencyMaxNS()
	p99V, p99OK := s.latencyP99NS()
	if !minOK && !meanOK && !maxOK && !p99OK {
		return naValue
	}
	return fmt.Sprintf("%s / %s / %s / %s",
		formatOptionalValue(minV, minOK, formatLatencyNS),
		formatOptionalValue(meanV, meanOK, formatLatencyNS),
		formatOptionalValue(maxV, maxOK, formatLatencyNS),
		formatOptionalValue(p99V, p99OK, formatLatencyNS))
}

func formatLatencyMeanP99Delta(curr, base *GlobalSnapshot) string {
	if curr == nil || base == nil {
		return naValue
	}
	cMean, cMeanOK := curr.latencyMeanNS()
	bMean, bMeanOK := base.latencyMeanNS()
	cP99, cP99OK := curr.latencyP99NS()
	bP99, bP99OK := base.latencyP99NS()
	dMean := formatOptionalDelta(cMean, cMeanOK, bMean, bMeanOK)
	dP99 := formatOptionalDelta(cP99, cP99OK, bP99, bP99OK)
	if dMean == naValue && dP99 == naValue {
		return naValue
	}
	return fmt.Sprintf("%s / %s", dMean, dP99)
}

func formatLatencyNS(ns float64) string {
	ns = finiteOrZero(ns)
	if ns <= 0 {
		return "0ns"
	}
	switch {
	case ns < nsPerMicrosecond:
		return fmt.Sprintf("%.0fns", ns)
	case ns < nsPerMillisecond:
		return fmt.Sprintf("%.2fµs", ns/nsPerMicrosecond)
	case ns < nsPerSecond:
		return fmt.Sprintf("%.2fms", ns/nsPerMillisecond)
	default:
		return fmt.Sprintf("%.2fs", ns/nsPerSecond)
	}
}

func formatDeltaPercent(curr, base float64) string {
	if !isFinitePositive(base) || !isFiniteNonNegative(curr) {
		return naValue
	}
	pct := ((curr - base) / base) * percentFactor
	if math.IsNaN(pct) || math.IsInf(pct, 0) {
		return naValue
	}
	if math.Abs(pct) < zeroPercentEpsilon {
		return "0.0%"
	}
	return fmt.Sprintf("%+.1f%%", pct)
}

func sanitizeTableCell(s string) string {
	return tableCellReplacer.Replace(s)
}

func escapeMarkdownCell(s string) string {
	return markdownCellReplacer.Replace(s)
}

func escapeHTMLSummary(s string) string {
	return html.EscapeString(htmlSummaryReplacer.Replace(s))
}

func formatInlineCode(s string) string {
	cell := sanitizeTableCell(s)
	if strings.Contains(cell, "`") {
		escaped := strings.ReplaceAll(html.EscapeString(cell), "`", "&#96;")
		return "<code>" + escaped + "</code>"
	}
	return "`" + cell + "`"
}

func newBoundedMarkdownBuilder(header string, maxBytes int) *boundedMarkdownBuilder {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxCommentBytes
	}
	b := &boundedMarkdownBuilder{header: header, maxBytes: maxBytes}
	fmt.Fprint(&b.out, header)
	return b
}

func (b *boundedMarkdownBuilder) append(block string) bool {
	if b.truncated {
		return false
	}
	if b.out.Len()+len(block)+len(truncationNotice) > b.maxBytes {
		b.truncated = true
		return false
	}
	fmt.Fprint(&b.out, block)
	return true
}

func (b *boundedMarkdownBuilder) finish() string {
	if b.truncated && b.out.Len()+len(truncationNotice) <= b.maxBytes {
		fmt.Fprint(&b.out, truncationNotice)
	}
	res := b.out.String()
	if len(res) > b.maxBytes {
		if len(b.header)+len(truncationNotice) <= b.maxBytes {
			return b.header + truncationNotice
		}
		return truncateUTF8(res, b.maxBytes)
	}
	return res
}

func truncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	if maxBytes <= 0 {
		return ""
	}
	end := maxBytes
	for end > 0 && !utf8.RuneStart(s[end]) {
		end--
	}
	return s[:end]
}
