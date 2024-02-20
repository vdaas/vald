// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package version

import (
	"context"
	"fmt"
	"reflect"

	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/strings"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	name        = "app_version_info"
	description = "app version info"
)

var reps = strings.NewReplacer("_", " ", ",omitempty", "")

type version struct {
	kvs map[string]string
}

func New(labels ...string) metrics.Metric {
	return &version{
		kvs: labelKVs(labels...),
	}
}

func labelKVs(labels ...string) map[string]string {
	labelMap := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		labelMap[reps.Replace(label)] = struct{}{}
	}

	d := info.Get()
	rt, rv := reflect.TypeOf(d), reflect.ValueOf(d)
	info := make(map[string]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		keyName := reps.Replace(rt.Field(i).Tag.Get("json"))
		if _, ok := labelMap[keyName]; !ok {
			continue
		}

		v := rv.Field(i).Interface()
		value, ok := v.(string)
		if !ok {
			ss, ok := v.([]string)
			if ok {
				// tags must be less than 255 characters
				info[keyName] = fmt.Sprintf("%.255s", fmt.Sprintf("%v", ss))
			}
			continue
		}
		if value != "" {
			info[keyName] = value
		}
	}

	return info
}

func (*version) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        name,
				Description: description,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

func (v *version) Register(m metrics.Meter) (err error) {
	info, err := m.Int64ObservableGauge(
		name,
		metrics.WithDescription(description),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrs := make([]attribute.KeyValue, 0, len(v.kvs))
			for key, val := range v.kvs {
				attrs = append(attrs, attribute.String(key, val))
			}
			o.ObserveInt64(info, 1, api.WithAttributes(attrs...))
			return nil
		}, info,
	)
	return err
}
