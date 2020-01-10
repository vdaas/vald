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
			name:    "parse duration success (ns)",
			t:       "1ns",
			want:    time.Nanosecond,
			wantErr: false,
		},
		{
			name:    "parse duration success (ms)",
			t:       "1ms",
			want:    time.Millisecond,
			wantErr: false,
		},
		{
			name:    "parse duration success (s)",
			t:       "1s",
			want:    time.Second,
			wantErr: false,
		},
		{
			name:    "parse duration success (m)",
			t:       "1m",
			want:    time.Minute,
			wantErr: false,
		},
		{
			name:    "parse duration success (h)",
			t:       "1h",
			want:    time.Hour,
			wantErr: false,
		},
		{
			name:    "parse empty string success",
			t:       "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "parse incorrect string return error",
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
