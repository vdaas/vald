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

// Package info provides general info metrics functions
package info

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type info struct {
	name     string
	fullname string
	info     metrics.Int64Measure
	kvs      map[metrics.Key]string
}

// New creates new general info metric according to the provided struct.
func New(name, fullname, description string, i interface{}) (metrics.Metric, error) {
	kvs, err := labelKVs(i)
	if err != nil {
		return nil, err
	}

	return &info{
		name:     name,
		fullname: fullname,
		info:     *metrics.Int64(metrics.ValdOrg+"/"+name, description, metrics.UnitDimensionless),
		kvs:      kvs,
	}, nil
}

func labelKVs(i interface{}) (map[metrics.Key]string, error) {
	rt, rv := reflect.TypeOf(i), reflect.ValueOf(i)
	kvs := make(map[metrics.Key]string, rt.NumField())
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

		mk, err := metrics.NewKey(keyName)
		if err != nil {
			return nil, err
		}

		kvs[mk] = value
	}

	return kvs, nil
}

func (i *info) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (i *info) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{
		{
			Measurement: i.info.M(int64(1)),
			Tags:        i.kvs,
		},
	}, nil
}

func (i *info) View() []*metrics.View {
	keys := make([]metrics.Key, 0, len(i.kvs))
	for k := range i.kvs {
		keys = append(keys, k)
	}

	return []*metrics.View{
		{
			Name:        i.fullname,
			Description: i.info.Description(),
			TagKeys:     keys,
			Measure:     &i.info,
			Aggregation: metrics.LastValue(),
		},
	}
}
