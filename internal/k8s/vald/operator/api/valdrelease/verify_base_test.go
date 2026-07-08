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

package valdrelease

import (
	"testing"

	"github.com/vdaas/vald/internal/k8s/resource"
)

// TestVerifyBaseEmbedding guards the resource.Base first-field offset-0
// contract for every Base-embedding type in this package.
func TestVerifyBaseEmbedding(t *testing.T) {
	t.Parallel()

	if err := resource.VerifyBase[ValdRelease](); err != nil {
		t.Errorf("ValdRelease violates the Base embedding contract: %v", err)
	}
}
