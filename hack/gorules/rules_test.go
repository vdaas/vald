// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package gorules

import (
	"os"
	"testing"

	"github.com/quasilyte/go-ruleguard/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRules(t *testing.T) {
	if err := analyzer.Analyzer.Flags.Set("rules", "rules.go"); err != nil {
		t.Fatalf("set rules flag: %v", err)
	}
	analysistest.Run(
		t,
		os.Getenv("GOPATH"),
		analyzer.Analyzer,
		"github.com/vdaas/vald/hack/gorules/testdata",
	)
}
