//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package label provides custom label metrics functions
package label

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type custom struct {
	info metrics.Int64Measure
	kvs  map[metrics.Key]string
}

func New(m map[string]string) (metrics.Metric, error) {
	kvs, err := labelKVs(m)
	if err != nil {
		return nil, err
	}

	return &custom{
		info: *metrics.Int64(metrics.ValdOrg+"/custom/label_info", "custom label info", metrics.UnitDimensionless),
		kvs:  kvs,
	}, nil
}

func labelKVs(m map[string]string) (map[metrics.Key]string, error) {
	info := make(map[metrics.Key]string, len(m))
	for k, v := range m {
		k, err := metrics.NewKey(k)
		if err != nil {
			return nil, err
		}
		// tags must be less than 255 characters
		info[k] = fmt.Sprintf("%.255s", fmt.Sprintf("%v", v))
	}

	return info, nil
}

func (c *custom) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (c *custom) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{
		metrics.MeasurementWithTags{
			Measurement: c.info.M(int64(1)),
			Tags:        c.kvs,
		},
	}, nil
}

func (c *custom) View() []*metrics.View {
	keys := make([]metrics.Key, 0, len(c.kvs))
	for k := range c.kvs {
		keys = append(keys, k)
	}

	return []*metrics.View{
		&metrics.View{
			Name:        "app_custom_label_info",
			Description: "custom label info",
			TagKeys:     keys,
			Measure:     &c.info,
			Aggregation: metrics.LastValue(),
		},
	}
}
