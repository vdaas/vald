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
package assets

import (
	"testing"

	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

type Dataset = assets.Dataset

func Data(name string) func(testing.TB) Dataset {
	return func(tb testing.TB) Dataset {
		tb.Helper()
		fn := assets.Data(name)
		if fn == nil {
			tb.Fatalf("not supported dataset: %s", name)
		}
		dataset, err := fn()
		if err != nil {
			tb.Fatalf("failed to load dataset: %s, err: %v", name, err)
		}
		return dataset
	}
}
