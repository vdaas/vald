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

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/sync/errgroup"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
)

func mustNewGroup(t *testing.T) errgroup.Group {
	t.Helper()
	eg, _ := errgroup.New(t.Context())
	return eg
}

func TestWithErrGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setsEG bool
	}{
		{
			name:   "replaces the errgroup when a non-nil group is given",
			setsEG: true,
		},
		{
			name:   "keeps the previous errgroup when nil is given",
			setsEG: false,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			preset := mustNewGroup(tt)
			o := &operator{eg: preset}

			var given errgroup.Group
			if test.setsEG {
				given = mustNewGroup(tt)
			}

			err := WithErrGroup(given)(o)
			require.NoError(tt, err)

			if test.setsEG {
				require.Same(tt, given, o.eg)
			} else {
				require.Same(tt, preset, o.eg)
			}
		})
	}
}

func TestWithReadReplicaEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		enabled bool
	}{
		{name: "enables read replica", enabled: true},
		{name: "disables read replica", enabled: false},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			o := new(operator)
			err := WithReadReplicaEnabled(test.enabled)(o)
			require.NoError(tt, err)
			require.Equal(tt, test.enabled, o.readReplicaEnabled)
		})
	}
}

func TestWithReadReplicaLabelKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
	}{
		{name: "sets a non-empty key", key: "vald-readreplica-id"},
		{name: "sets an empty key", key: ""},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			o := new(operator)
			err := WithReadReplicaLabelKey(test.key)(o)
			require.NoError(tt, err)
			require.Equal(tt, test.key, o.readReplicaLabelKey)
		})
	}
}

func TestWithRotationJobConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		concurrency uint
		wantErr     bool
	}{
		{
			name:        "sets the concurrency when greater than 0",
			concurrency: 5,
		},
		{
			name:        "returns ErrCriticalOption when 0 is given",
			concurrency: 0,
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			o := new(operator)
			err := WithRotationJobConcurrency(test.concurrency)(o)

			if test.wantErr {
				var critical *errors.ErrCriticalOption
				require.ErrorAs(tt, err, &critical)
				require.Zero(tt, o.rotationJobConcurrency, "operator state must be left untouched on failure")
				return
			}
			require.NoError(tt, err)
			require.Equal(tt, test.concurrency, o.rotationJobConcurrency)
		})
	}
}

func TestWithK8sClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		givenNil  bool
		wantsSame bool
	}{
		{
			name:      "sets the client when non-nil",
			wantsSame: true,
		},
		{
			name:     "keeps the previous client when nil is given",
			givenNil: true,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			preset := &mock.ValdK8sClientMock{}
			o := &operator{client: preset}

			var given client.Client
			if !test.givenNil {
				given = &mock.ValdK8sClientMock{}
			}

			err := WithK8sClient(given)(o)
			require.NoError(tt, err)

			if test.wantsSame {
				require.Same(tt, given, o.client)
			} else {
				require.Same(tt, preset, o.client)
			}
		})
	}
}

func TestDefaultOpts(t *testing.T) {
	t.Parallel()

	o := new(operator)
	for _, opt := range defaultOpts {
		require.NoError(t, opt(o))
	}

	require.NotNil(t, o.eg, "WithErrGroup(errgroup.Get()) must set a non-nil errgroup")
	require.Equal(t, uint(1), o.rotationJobConcurrency, "WithRotationJobConcurrency(1) must be applied by default")
}

// NOT IMPLEMENTED BELOW
