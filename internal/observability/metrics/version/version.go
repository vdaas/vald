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

// Package version provides version info metrics functions
package version

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/observability/metrics"
)

var (
	reps = strings.NewReplacer("_", " ", ",omitempty", "")
)

type version struct {
	info metrics.Int64Measure
	kvs  map[metrics.Key]string
}

func New() (metrics.Metric, error) {
	kvs, err := labelKVs()
	if err != nil {
		return nil, err
	}

	return &version{
		info: *metrics.Int64(metrics.ValdOrg+"/version/info", "app version info", metrics.UnitDimensionless),
		kvs:  kvs,
	}, nil
}

func labelKVs() (map[metrics.Key]string, error) {
	d := info.Get()
	rt, rv := reflect.TypeOf(d), reflect.ValueOf(d)
	info := make(map[metrics.Key]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		v := rv.Field(i).Interface()
		value, ok := v.(string)
		if !ok {
			ss, ok := v.([]string)
			if ok {
				k, err := metrics.NewKey(reps.Replace(rt.Field(i).Tag.Get("json")))
				if err != nil {
					return nil, err
				}
				// tags must be less than 255 characters
				info[k] = fmt.Sprintf("%.255s", fmt.Sprintf("%v", ss))
			}
			continue
		}
		if value != "" {
			k, err := metrics.NewKey(reps.Replace(rt.Field(i).Tag.Get("json")))
			if err != nil {
				return nil, err
			}
			info[k] = value
		}
	}

	return info, nil
}

func (v *version) MeasurementsCount() int {
	cnt := 0
	rv := reflect.ValueOf(*v)
	for i := 0; i < rv.NumField(); i++ {
		if metrics.IsMeasureType(rv.Field(i).Type()) {
			cnt++
		}
	}
	return cnt
}

func (v *version) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (v *version) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{
		metrics.MeasurementWithTags{
			Measurement: v.info.M(int64(1)),
			Tags:        v.kvs,
		},
	}, nil
}

func (v *version) View() []*metrics.View {
	keys := make([]metrics.Key, 0, len(v.kvs))
	for k := range v.kvs {
		keys = append(keys, k)
	}

	return []*metrics.View{
		&metrics.View{
			Name:        "app_version_info",
			Description: "app version info",
			TagKeys:     keys,
			Measure:     &v.info,
			Aggregation: metrics.LastValue(),
		},
	}
}
