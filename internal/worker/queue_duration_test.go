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

package worker

import (
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

// time.NewTicker panics for a non-positive interval, so the option must reject it
// and leave the previously configured (default) value untouched.
func TestWithQueueCheckDuration_NonPositive(t *testing.T) {
	t.Parallel()
	const initial = 7 * time.Second
	tests := []struct {
		name    string
		dur     string
		want    time.Duration
		wantErr bool
	}{
		{name: "zero is rejected", dur: "0s", want: initial, wantErr: true},
		{name: "negative is rejected", dur: "-1s", want: initial, wantErr: true},
		{name: "positive is applied", dur: "2s", want: 2 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := new(queue)
			o.qcdur = initial
			err := WithQueueCheckDuration(tt.dur)(o)
			if tt.wantErr {
				var e *errors.ErrInvalidOption
				if !errors.As(err, &e) {
					t.Errorf("error = %v, want ErrInvalidOption", err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if o.qcdur != tt.want {
				t.Errorf("qcdur = %v, want %v", o.qcdur, tt.want)
			}
		})
	}
}
