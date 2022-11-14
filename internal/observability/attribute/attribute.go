// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package attribute

import "go.opentelemetry.io/otel/attribute"

type (
	KeyValue = attribute.KeyValue
	Key      = attribute.Key
)

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
