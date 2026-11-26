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

package metrics

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/attribute"
	"go.opentelemetry.io/otel/metric"
)

type stubInt64SumMetric struct {
	viewErr error
	regErr  error
	name    string
	val     int64
}

func (s *stubInt64SumMetric) View() ([]View, error) {
	if s.viewErr != nil {
		return nil, s.viewErr
	}
	return nil, nil
}

func (s *stubInt64SumMetric) Register(m Meter) error {
	if s.regErr != nil {
		return s.regErr
	}
	ctr, err := m.Int64ObservableCounter(s.name)
	if err != nil {
		return err
	}
	_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(ctr, s.val, metric.WithAttributes(attribute.String("k", "v")))
		return nil
	}, ctr)
	return err
}

func TestCollectInt64Sums(t *testing.T) {
	t.Parallel()
	attrSet := attribute.NewSet(attribute.String("k", "v"))
	key := attrSet.Equivalent()

	got, err := CollectInt64Sums(t.Context(), &stubInt64SumMetric{name: "test_counter", val: 42})
	if err != nil {
		t.Fatalf("CollectInt64Sums() error = %v", err)
	}
	if got["test_counter"][key] != 42 {
		t.Errorf("test_counter = %d, want 42", got["test_counter"][key])
	}

	viewErr := errors.New("view boom")
	if _, err := CollectInt64Sums(t.Context(), &stubInt64SumMetric{viewErr: viewErr}); !errors.Is(err, viewErr) {
		t.Errorf("err = %v, want %v", err, viewErr)
	}

	regErr := errors.New("register boom")
	if _, err := CollectInt64Sums(t.Context(), &stubInt64SumMetric{regErr: regErr}); !errors.Is(err, regErr) {
		t.Errorf("err = %v, want %v", err, regErr)
	}
}
