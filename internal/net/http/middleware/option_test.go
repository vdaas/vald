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
package middleware

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
)

func TestWithErrorGroup(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(TimeoutOption) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt TimeoutOption) error {
				got := new(timeout)
				opt(got)

				if got.dur != 10*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			dur:  "ok",
			checkFunc: func(opt TimeoutOption) error {
				got := new(timeout)
				opt(got)

				if got.dur != 3*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTimeout(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type test struct {
		name      string
		eg        errgroup.Group
		checkFunc func(TimeoutOption) error
	}

	tests := []test{
		func() test {
			eg := errgroup.Get()

			return test{
				name: "set success",
				eg:   eg,
				checkFunc: func(opt TimeoutOption) error {
					got := new(timeout)
					opt(got)

					if !reflect.DeepEqual(got.eg, eg) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithErrorGroup(tt.eg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
