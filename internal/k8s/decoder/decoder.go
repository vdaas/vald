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
package decoder

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"
)

// Decoder represents a type alias of conversion.Decoder.
type Decoder = conversion.Decoder

// NewDecoder creates a Decoder given the runtime.Scheme.
// It will return an error when NewDecoder method failed.
func NewDecoder(scheme *runtime.Scheme) (*Decoder, error) {
	return conversion.NewDecoder(scheme)
}
