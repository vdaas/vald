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
package format

import "testing"

func TestString(t *testing.T) {
	type test struct {
		name   string
		format Format
		want   string
	}

	tests := []test{
		{
			name:   "returns raw",
			format: RAW,
			want:   "raw",
		},

		{
			name:   "returns json",
			format: JSON,
			want:   "json",
		},

		{
			name:   "returns ltsv",
			format: LTSV,
			want:   "ltsv",
		},

		{
			name:   "returns unknown",
			format: Format(100),
			want:   "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.format.String()
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func TestAtof(t *testing.T) {
	type test struct {
		name string
		str  string
		want Format
	}

	tests := []test{
		{
			name: "returns RAW when str is raw",
			str:  "raw",
			want: RAW,
		},

		{
			name: "returns RAW when str is RAw",
			str:  "RAw",
			want: RAW,
		},

		{
			name: "returns JSON when str is json",
			str:  "json",
			want: JSON,
		},

		{
			name: "returns JSON when str is JSOn",
			str:  "JSOn",
			want: JSON,
		},

		{
			name: "returns LTSV when str is ltsv",
			str:  "ltsv",
			want: LTSV,
		},

		{
			name: "returns LTSV when str is LTSv",
			str:  "LTSv",
			want: LTSV,
		},

		{
			name: "returns Unknown when str is Vald",
			str:  "Vald",
			want: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atof(tt.str)
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}
