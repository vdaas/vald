//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package client

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"k8s.io/client-go/rest"
)

// TestFallbackToInCluster runs outside a real Kubernetes pod (as every unit
// test environment does), so rest.InClusterConfig always fails with
// rest.ErrNotInCluster; that determinism is what these cases rely on.
func TestFallbackToInCluster(t *testing.T) {
	t.Run("returns the in-cluster error alone when there is no preceding error", func(t *testing.T) {
		c, err := fallbackToInCluster(nil)
		if c != nil {
			t.Errorf("fallbackToInCluster(nil) client = %v, want nil", c)
		}
		if !errors.Is(err, rest.ErrNotInCluster) {
			t.Errorf("fallbackToInCluster(nil) error = %v, want to wrap rest.ErrNotInCluster", err)
		}
	})

	t.Run("joins the preceding error with the in-cluster error", func(t *testing.T) {
		origErr := errors.New("kubeconfig load failed")
		c, err := fallbackToInCluster(origErr)
		if c != nil {
			t.Errorf("client = %v, want nil", c)
		}
		if !errors.Is(err, origErr) {
			t.Errorf("error = %v, want to wrap origErr %v", err, origErr)
		}
		if !errors.Is(err, rest.ErrNotInCluster) {
			t.Errorf("error = %v, want to also wrap rest.ErrNotInCluster", err)
		}
	})
}
