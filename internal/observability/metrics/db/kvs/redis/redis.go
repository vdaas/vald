//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package redis provides redis metrics functions
package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/trace"
)

type redisMetrics struct {
	queryTotal   metrics.Int64Measure
	queryLatency metrics.Float64Measure

	pipelineTotal   metrics.Int64Measure
	pipelineLatency metrics.Float64Measure

	cmdNameKey metrics.Key
	numCmdKey  metrics.Key

	mu sync.Mutex
	ms []metrics.MeasurementWithTags
}

type MetricsHook interface {
	metrics.Metric
	redis.Hook
}

type (
	startTimeKey         struct{}
	pipelineStartTimeKey struct{}
)

func New() (o MetricsHook, err error) {
	rms := new(redisMetrics)

	rms.queryTotal = *metrics.Int64(
		metrics.ValdOrg+"/db/kvs/redis/completed_query_total",
		"cumulative count of completed queries",
		metrics.UnitDimensionless,
	)

	rms.queryLatency = *metrics.Float64(
		metrics.ValdOrg+"/db/kvs/redis/query_latency",
		"query latency",
		metrics.UnitMilliseconds,
	)

	rms.pipelineTotal = *metrics.Int64(
		metrics.ValdOrg+"/db/kvs/redis/completed_pipeline_total",
		"cumulative count of completed pipeline",
		metrics.UnitDimensionless,
	)

	rms.pipelineLatency = *metrics.Float64(
		metrics.ValdOrg+"/db/kvs/redis/pipeline_latency",
		"pipeline latency",
		metrics.UnitMilliseconds,
	)

	rms.cmdNameKey, err = metrics.NewKey("redis_cmd_name")
	if err != nil {
		return nil, err
	}

	rms.numCmdKey, err = metrics.NewKey("redis_num_cmd")
	if err != nil {
		return nil, err
	}

	rms.ms = make([]metrics.MeasurementWithTags, 0)

	return rms, nil
}

func (rm *redisMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (rm *redisMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	rm.mu.Lock()
	defer func() {
		rm.ms = make([]metrics.MeasurementWithTags, 0)
		rm.mu.Unlock()
	}()

	return rm.ms, nil
}

func (rm *redisMetrics) View() []*metrics.View {
	queryKeys := []metrics.Key{
		rm.cmdNameKey,
	}

	pipelineKeys := []metrics.Key{
		rm.numCmdKey,
	}

	return []*metrics.View{
		{
			Name:        "db_kvs_redis_completed_query_total",
			Description: rm.queryTotal.Description(),
			TagKeys:     queryKeys,
			Measure:     &rm.queryTotal,
			Aggregation: metrics.Count(),
		},
		{
			Name:        "db_kvs_redis_query_latency",
			Description: rm.queryLatency.Description(),
			TagKeys:     queryKeys,
			Measure:     &rm.queryLatency,
			Aggregation: metrics.DefaultMillisecondsDistribution,
		},
		{
			Name:        "db_kvs_redis_completed_pipeline_total",
			Description: rm.pipelineTotal.Description(),
			TagKeys:     pipelineKeys,
			Measure:     &rm.pipelineTotal,
			Aggregation: metrics.Count(),
		},
		{
			Name:        "db_kvs_redis_pipeline_latency",
			Description: rm.pipelineLatency.Description(),
			TagKeys:     pipelineKeys,
			Measure:     &rm.pipelineLatency,
			Aggregation: metrics.DefaultMillisecondsDistribution,
		},
	}
}

func (rm *redisMetrics) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/db/kvs/redis")
	if span != nil {
		span.AddAttributes(
			trace.StringAttribute("cmd", cmd.Name()),
		)
	}

	return context.WithValue(ctx, startTimeKey{}, time.Now()), nil
}

func (rm *redisMetrics) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.FromContext(ctx)
	if span != nil {
		if err := cmd.Err(); err != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}

		span.End()
	}

	startTime, _ := ctx.Value(startTimeKey{}).(time.Time)
	latencyMillis := float64(time.Since(startTime)) / float64(time.Millisecond)

	tags := map[metrics.Key]string{
		rm.cmdNameKey: cmd.Name(),
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.ms = append(
		rm.ms,
		metrics.MeasurementWithTags{
			Measurement: rm.queryTotal.M(1),
			Tags:        tags,
		},
		metrics.MeasurementWithTags{
			Measurement: rm.queryLatency.M(latencyMillis),
			Tags:        tags,
		},
	)

	return nil
}

func (rm *redisMetrics) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	cmdStrs := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		cmdStrs = append(cmdStrs, cmd.Name())
	}

	ctx, span := trace.StartSpan(ctx, "vald/internal/db/kvs/redis/pipeline")
	if span != nil {
		span.AddAttributes(
			trace.Int64Attribute("num_cmd", int64(len(cmds))),
			trace.StringAttribute("cmds", fmt.Sprintf("%v", cmdStrs)),
		)
	}

	return context.WithValue(ctx, pipelineStartTimeKey{}, time.Now()), nil
}

func (rm *redisMetrics) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := trace.FromContext(ctx)
	if span != nil {
		var errs error
		for _, cmd := range cmds {
			if err := cmd.Err(); err != nil {
				errs = errors.Wrap(errs, err.Error())
			}
		}

		if errs != nil {
			span.SetStatus(trace.StatusCodeUnknown(errs.Error()))
		}

		span.End()
	}

	startTime, _ := ctx.Value(pipelineStartTimeKey{}).(time.Time)
	latencyMillis := float64(time.Since(startTime)) / float64(time.Millisecond)

	tags := map[metrics.Key]string{
		rm.numCmdKey: strconv.Itoa(len(cmds)),
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.ms = append(
		rm.ms,
		metrics.MeasurementWithTags{
			Measurement: rm.pipelineTotal.M(1),
			Tags:        tags,
		},
		metrics.MeasurementWithTags{
			Measurement: rm.pipelineLatency.M(latencyMillis),
			Tags:        tags,
		},
	)

	return nil
}
