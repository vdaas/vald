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
package service

import (
	"testing"
	"time"
)

// saveIndexDuration drives the regular save ticker in Start; when it stays 0 Start
// turns it into math.MaxInt64 and only the forced saveIndexDurationLimit pass runs.
func TestNew_saveIndexDuration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		opts    []Option
		want    time.Duration
		wantErr bool
	}{
		{name: "default is unset", want: 0},
		{name: "option is applied", opts: []Option{WithSaveIndexDuration("10m")}, want: 10 * time.Minute},
		{name: "empty keeps unset", opts: []Option{WithSaveIndexDuration("")}, want: 0},
		{name: "invalid duration fails", opts: []Option{WithSaveIndexDuration("invalid")}, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			got, err := New(tc.opts...)
			if tc.wantErr {
				if err == nil {
					tt.Fatal("New() error = nil, want error")
				}
				return
			}
			if err != nil {
				tt.Fatalf("New() error = %v", err)
			}
			idx, ok := got.(*index)
			if !ok {
				tt.Fatalf("New() returned %T, want *index", got)
			}
			if idx.saveIndexDuration != tc.want {
				tt.Errorf("saveIndexDuration = %v, want %v", idx.saveIndexDuration, tc.want)
			}
			if idx.saveIndexDurationLimit != 3*time.Hour {
				tt.Errorf("saveIndexDurationLimit = %v, want default 3h", idx.saveIndexDurationLimit)
			}
		})
	}
}
