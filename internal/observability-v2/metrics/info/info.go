package info

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"go.opentelemetry.io/otel/attribute"
)

type info struct {
	name        string
	description string
	kvs         map[string]string
}

// New creates new general info metric according to the provided struct.
func New(name, description string, i interface{}) metrics.Metric {
	return &info{
		name: name,
		kvs:  labelKVs(i),
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

func (i *info) Register(m metrics.Meter) error {
	info, err := m.AsyncInt64().Gauge(
		i.name,
		metrics.WithDescription(i.description),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			info,
		},
		func(ctx context.Context) {
			attrs := make([]attribute.KeyValue, 0, len(i.kvs))
			for key, val := range i.kvs {
				attrs = append(attrs, attribute.String(key, val))
			}
			info.Observe(ctx, 1, attrs...)
		},
	)
}
