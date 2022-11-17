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
package middleware

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

func TestWithErrorGroup(t *testing.T) {
	t.Parallel()
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
					return errors.Errorf("invalid param was set")
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
					return errors.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			opt := WithTimeout(test.dur)
			if err := test.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()
	type test struct {
		name      string
		eg        errgroup.Group
		checkFunc func(TimeoutOption) error
	}

	tests := []test{
		func() test {
			eg, _ := errgroup.New(context.Background())

			return test{
				name: "set success",
				eg:   eg,
				checkFunc: func(opt TimeoutOption) error {
					got := new(timeout)
					opt(got)

					if !reflect.DeepEqual(got.eg, eg) {
						return errors.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			opt := WithErrorGroup(test.eg)
			if err := test.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
