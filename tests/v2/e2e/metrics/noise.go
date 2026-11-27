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
	"math"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
)

type (
	// NoiseConfig configures environmental noise detection thresholds and sample guards.
	NoiseConfig struct {
		Threshold      float64 `json:"threshold"`
		MinDurationSec float64 `json:"min_duration_sec"`
		MinTotal       uint64  `json:"min_total"`
	}

	metricPair struct {
		name string
		c, b float64
	}
)

const (
	// DefaultNoiseThreshold is the default relative deviation threshold (30%).
	DefaultNoiseThreshold = 0.30
	// DefaultMinDurationSec is the default minimum execution duration in seconds for QPS/Duration checks.
	DefaultMinDurationSec = 0.1
	// DefaultMinTotal is the default minimum sample count required before evaluating noise.
	DefaultMinTotal = uint64(2)

	noiseDeltaEpsilon       = 1e-9
	maxNoiseMetricsPerScope = 4
)

// DefaultNoiseConfig returns the default noise detection configuration.
func DefaultNoiseConfig() NoiseConfig {
	return NoiseConfig{
		Threshold:      DefaultNoiseThreshold,
		MinDurationSec: DefaultMinDurationSec,
		MinTotal:       DefaultMinTotal,
	}
}

// CheckNoise compares current against baseline across Report.Summary and matched executions.
func CheckNoise(current, baseline *Report, cfg NoiseConfig) (retry bool, reason string, err error) {
	if !isFinitePositive(cfg.Threshold) {
		return false, "", errors.Errorf("threshold must be finite and > 0, got %v", cfg.Threshold)
	}
	if !isFiniteNonNegative(cfg.MinDurationSec) {
		return false, "", errors.Errorf("min-duration must be finite and >= 0, got %v", cfg.MinDurationSec)
	}
	if current == nil {
		return false, "", errors.New("current report is nil")
	}
	if baseline == nil {
		return false, "baseline report is nil", nil
	}

	var (
		compared  int
		maxDelta  float64
		worstDesc string
	)
	check := func(scope string, curr, base *GlobalSnapshot) {
		// Require minimum sample count so tiny operations do not trigger false noise retries.
		if curr == nil || base == nil || curr.Total < cfg.MinTotal || base.Total < cfg.MinTotal {
			return
		}
		pairs := make([]metricPair, 0, maxNoiseMetricsPerScope)
		cDur, bDur := curr.Duration().Seconds(), base.Duration().Seconds()
		// Sub-threshold durations amplify timer jitter in QPS; skip QPS/Duration comparisons when too short.
		if cDur >= cfg.MinDurationSec && bDur >= cfg.MinDurationSec {
			pairs = append(pairs,
				metricPair{"qps", curr.QPS(), base.QPS()},
				metricPair{"duration_sec", cDur, bDur},
			)
		}
		pairs = append(pairs,
			metricPair{"latency.mean_ns", curr.LatencyMeanNS(), base.LatencyMeanNS()},
			metricPair{"latency.p99_ns", curr.LatencyP99NS(), base.LatencyP99NS()},
		)
		for _, p := range pairs {
			if !isFinitePositive(p.b) || !isFiniteNonNegative(p.c) {
				continue
			}
			compared++
			delta := math.Abs(p.c-p.b) / p.b
			// Tolerance absorbs IEEE-754 rounding when a delta sits exactly on cfg.Threshold.
			if delta > cfg.Threshold+noiseDeltaEpsilon && delta > maxDelta {
				maxDelta = delta
				worstDesc = fmt.Sprintf(
					"%s %s delta %.1f%% > %.1f%%",
					scope,
					p.name,
					delta*percentFactor,
					cfg.Threshold*percentFactor,
				)
			}
		}
	}

	check(summaryScope, current.Summary, baseline.Summary)
	currIdx := buildExecutionIndex(current)
	baseIdx := buildExecutionIndex(baseline)
	for _, fe := range currIdx.ordered {
		if baseExec, ok := baseIdx.match(fe); ok {
			check(fe.Label, fe.Execution.Summary, baseExec.Summary)
		}
	}

	if compared == 0 {
		return false, "no comparable metrics met min_total/min_duration guards", nil
	}
	if worstDesc != "" {
		return true, worstDesc, nil
	}
	return false, "", nil
}

// CheckNoiseBaselineFile compares an in-memory current Report against baselinePath,
// skipping cleanly if baselinePath is missing or empty.
func CheckNoiseBaselineFile(
	current *Report, baselinePath string, cfg NoiseConfig,
) (retry bool, reason string, err error) {
	if current == nil {
		return false, "", errors.New("current report is nil")
	}
	if strings.TrimSpace(baselinePath) == "" {
		return false, "baseline not specified", nil
	}
	fi, exist, err := statPath(baselinePath)
	if err != nil {
		return false, "", errors.Wrapf(err, "failed to stat baseline report %s", baselinePath)
	}
	if !exist {
		return false, "baseline file not found", nil
	}
	if fi.IsDir() {
		return false, "", errors.Errorf("baseline path %s is a directory", baselinePath)
	}
	if fi.Size() == 0 {
		return false, "baseline file is empty", nil
	}
	baseline, err := LoadReport(baselinePath)
	if err != nil {
		return false, "", err
	}
	return CheckNoise(current, baseline, cfg)
}

// CheckNoiseFile loads currentPath and optional baselinePath, skipping cleanly if baselinePath is missing or empty.
func CheckNoiseFile(
	currentPath, baselinePath string, cfg NoiseConfig,
) (retry bool, reason string, err error) {
	current, err := LoadReport(currentPath)
	if err != nil {
		return false, "", err
	}
	return CheckNoiseBaselineFile(current, baselinePath, cfg)
}
