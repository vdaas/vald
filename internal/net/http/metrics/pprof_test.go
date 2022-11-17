//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package metrics

import "testing"

func TestNewPProfHandler(t *testing.T) {
	tests := []struct {
		name        string
		initialized bool
	}{
		{
			name:        "initialize success",
			initialized: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPProfHandler()
			if (got != nil) != tt.initialized {
				t.Error("NewPProfHandler() is wrong.")
			}
		})
	}
}
