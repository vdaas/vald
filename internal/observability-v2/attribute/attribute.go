package attribute

import "go.opentelemetry.io/otel/attribute"

type KeyValue = attribute.KeyValue

func Bool(k string, v bool) KeyValue {
	return attribute.Bool(k, v)
}

func String(k, v string) KeyValue {
	return attribute.String(k, v)
}

func Int64(k string, v int64) KeyValue {
	return attribute.Int64(k, v)
}

func Float64(k string, v float64) KeyValue {
	return attribute.Float64(k, v)
}
