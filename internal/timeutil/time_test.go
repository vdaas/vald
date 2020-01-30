//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

package timeutil

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		t       string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "returns time.Nanosecond and nil when t is 1ns",
			t:       "1ns",
			want:    time.Nanosecond,
			wantErr: false,
		},
		{
			name:    "returns time.Millisecond and nil when t is 1ms",
			t:       "1ms",
			want:    time.Millisecond,
			wantErr: false,
		},
		{
			name:    "returns time.Second and nil when t is 1s",
			t:       "1s",
			want:    time.Second,
			wantErr: false,
		},
		{
			name:    "returns tme.Minute and nil when t is 1m",
			t:       "1m",
			want:    time.Minute,
			wantErr: false,
		},
		{
			name:    "returns time.Hour and nil when t is 1h",
			t:       "1h",
			want:    time.Hour,
			wantErr: false,
		},
		{
			name:    "returns 0 and nil when t is empty",
			t:       "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "returns 0 and incorrect string error when t is invalid",
			t:       "dummystring",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
