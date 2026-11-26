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

package types

import (
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	// Any is a type alias of anypb.Any representing a serialized protocol buffer message.
	Any = anypb.Any
	// FloatValue is a type alias of wrapperspb.FloatValue representing a wrapped float32 value.
	FloatValue = wrapperspb.FloatValue
)

// Float returns a wrapped FloatValue containing v.
func Float(v float32) *FloatValue {
	return wrapperspb.Float(v)
}
