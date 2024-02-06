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
package info

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

type info struct {
	name        string
	description string
	kvs         map[string]string
}

// New creates new general info metric according to the provided struct.
func New(name, description string, i interface{}) metrics.Metric {
	return &info{
		name:        name,
		description: description,
		kvs:         labelKVs(i),
	}
}

func labelKVs(i interface{}) map[string]string {
	rt, rv := reflect.TypeOf(i), reflect.ValueOf(i)
	kvs := make(map[string]string, rt.NumField())
	for k := 0; k < rt.NumField(); k++ {
		keyName := rt.Field(k).Tag.Get("info")
		if keyName == "" {
			continue
		}

		var value string

		switch v := rv.Field(k); v.Kind() {
		case reflect.String:
			value = v.String()
		case reflect.Slice, reflect.Array:
			switch v.Interface().(type) {
			case []string:
				value = fmt.Sprintf("%.255s", fmt.Sprintf("%v", v.Interface()))
			case []rune:
				value = v.Convert(reflect.TypeOf("")).String()
			default:
				continue
			}
		case reflect.Bool:
			value = strconv.FormatBool(v.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(v.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value = strconv.FormatUint(v.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(v.Float(), 'E', -1, 64)
		default:
			continue
		}

		kvs[keyName] = value
	}

	return kvs
}

func (i *info) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        i.name,
				Description: i.description,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

func (i *info) Register(m metrics.Meter) error {
	info, err := m.Int64ObservableGauge(
		i.name,
		metrics.WithDescription(i.description),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrs := make([]attribute.KeyValue, 0, len(i.kvs))
			for key, val := range i.kvs {
				attrs = append(attrs, attribute.String(key, val))
			}
			o.ObserveInt64(info, 1, api.WithAttributes(attrs...))
			return nil
		}, info,
	)
	return err
}
